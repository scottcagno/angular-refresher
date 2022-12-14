package jwt

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

// timePrecision sets the precision of times and dates. This has an
// influence on the precision of times when comparing expiry or other
// related time fields. It is also the precision of times when serializing.
// For backwards compatibility the default precision is set to seconds, so
// that no fractional timestamps are generated.
const timePrecision = time.Second

// NumericDate represents a JSON numeric date value as referenced in
// the documentation at https://datatracker.ietf.org/doc/html/rfc7519#section-2.
type NumericDate struct {
	time.Time
}

// NewNumericDate constructs a new *NumericDate from a standard library
// time.Time struct. It will truncate the timestamp according to the precision
// specified in TimePrecision.
func NewNumericDate(t time.Time) *NumericDate {
	return &NumericDate{t.Truncate(timePrecision)}
}

// newNumericDateFromSeconds creates a new *NumericDate out of a float64
// representing a UNIX epoch with the float fraction representing non-integer
// seconds.
func newNumericDateFromSeconds(f float64) *NumericDate {
	round, frac := math.Modf(f)
	return NewNumericDate(time.Unix(int64(round), int64(frac*1e9)))
}

// MarshalJSON is an implementation of the json.RawMessage interface and
// serializes the UNIX epoch represented in NumericDate to a byte array, using
// the precision specified in timePrecision.
func (date NumericDate) MarshalJSON() (b []byte, err error) {
	trunc := date.Truncate(timePrecision)
	ns := float64(trunc.Nanosecond()) / float64(time.Second)
	output := append(
		[]byte(strconv.FormatInt(trunc.Unix(), 10)),
		[]byte(strconv.FormatFloat(ns, 'f', 0, 64))[1:]...,
	)
	return output, nil
}

func ErrParsingNumericDate(err error) error {
	return fmt.Errorf("could not parse NumericData: %w", err)
}

func ErrConvertingJSONNumber(err error) error {
	return fmt.Errorf("could not convert json number value to float: %w", err)
}

// UnmarshalJSON is an implementation of the json.RawMessage interface and
// de-serialises a NumericDate from a JSON representation, i.e. a json.Number.
// This number represents a UNIX epoch with either integer or non-integer seconds.
func (date *NumericDate) UnmarshalJSON(b []byte) (err error) {
	var number json.Number
	err = json.Unmarshal(b, &number)
	if err != nil {
		return ErrParsingNumericDate(err)
	}
	f, err := number.Float64()
	if err != nil {
		return ErrConvertingJSONNumber(err)
	}
	n := newNumericDateFromSeconds(f)
	*date = *n
	return nil
}

// ClaimStrings is basically just a slice of strings, but it can be either
// serialized from a string array or just a string. This type is necessary,
// since the "aud" claim can either be a single string or an array.
type ClaimStrings []string

func (s *ClaimStrings) UnmarshalJSON(data []byte) error {
	var value any
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	var aud []string
	switch v := value.(type) {
	case string:
		aud = append(aud, v)
	case []string:
		aud = v
	case []any:
		for _, vv := range v {
			vs, ok := vv.(string)
			if !ok {
				return &json.UnsupportedTypeError{Type: reflect.TypeOf(vv)}
			}
			aud = append(aud, vs)
		}
	case nil:
		return nil
	default:
		return &json.UnsupportedTypeError{Type: reflect.TypeOf(v)}
	}
	*s = aud
	return err
}

func (s ClaimStrings) MarshalJSON() (b []byte, err error) {
	return json.Marshal([]string(s))
}
