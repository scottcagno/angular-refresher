package jwt

import (
	"crypto"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

var errKeyFuncError error = fmt.Errorf("error loading key")

var (
	jwtTestDefaultKey      *rsa.PublicKey
	jwtTestRSAPrivateKey   *rsa.PrivateKey
	jwtTestEC256PublicKey  crypto.PublicKey
	jwtTestEC256PrivateKey crypto.PrivateKey
	paddedKey              crypto.PublicKey
	defaultKeyFunc         KeyFunc = func(t *Token) (interface{}, error) { return jwtTestDefaultKey, nil }
	ecdsaKeyFunc           KeyFunc = func(t *Token) (interface{}, error) { return jwtTestEC256PublicKey, nil }
	paddedKeyFunc          KeyFunc = func(t *Token) (interface{}, error) { return paddedKey, nil }
	emptyKeyFunc           KeyFunc = func(t *Token) (interface{}, error) { return nil, nil }
	errorKeyFunc           KeyFunc = func(t *Token) (interface{}, error) { return nil, errKeyFuncError }
	nilKeyFunc             KeyFunc = nil
)

var jwtTestData = []struct {
	name          string
	tokenString   string
	keyfunc       KeyFunc
	claims        Claims
	valid         bool
	errors        uint32
	err           []error
	parser        *Parser
	signingMethod SigningMethod // The method to sign the JWT token for test purpose
}{
	{
		"basic",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		defaultKeyFunc,
		MapClaims{"foo": "bar"},
		true,
		0,
		nil,
		nil,
		SigningMethodRS256,
	},
	{
		"basic expired",
		"", // autogen
		defaultKeyFunc,
		MapClaims{"foo": "bar", "exp": float64(time.Now().Unix() - 100)},
		false,
		ValidationErrorExpired,
		[]error{ErrTokenExpired},
		nil,
		SigningMethodRS256,
	},
	{
		"basic nbf",
		"", // autogen
		defaultKeyFunc,
		MapClaims{"foo": "bar", "nbf": float64(time.Now().Unix() + 100)},
		false,
		ValidationErrorNotValidYet,
		[]error{ErrTokenNotValidYet},
		nil,
		SigningMethodRS256,
	},
	{
		"expired and nbf",
		"", // autogen
		defaultKeyFunc,
		MapClaims{"foo": "bar", "nbf": float64(time.Now().Unix() + 100), "exp": float64(time.Now().Unix() - 100)},
		false,
		ValidationErrorNotValidYet | ValidationErrorExpired,
		[]error{ErrTokenNotValidYet},
		nil,
		SigningMethodRS256,
	},
	{
		"basic invalid",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.EhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		defaultKeyFunc,
		MapClaims{"foo": "bar"},
		false,
		ValidationErrorSignatureInvalid,
		[]error{ErrTokenSignatureInvalid, rsa.ErrVerification},
		nil,
		SigningMethodRS256,
	},
	{
		"basic nokeyfunc",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		nilKeyFunc,
		MapClaims{"foo": "bar"},
		false,
		ValidationErrorUnverifiable,
		[]error{ErrTokenUnverifiable},
		nil,
		SigningMethodRS256,
	},
	{
		"basic nokey",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		emptyKeyFunc,
		MapClaims{"foo": "bar"},
		false,
		ValidationErrorSignatureInvalid,
		[]error{ErrTokenSignatureInvalid},
		nil,
		SigningMethodRS256,
	},
	{
		"basic errorkey",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		errorKeyFunc,
		MapClaims{"foo": "bar"},
		false,
		ValidationErrorUnverifiable,
		[]error{ErrTokenUnverifiable, errKeyFuncError},
		nil,
		SigningMethodRS256,
	},
	{
		"invalid signing method",
		"",
		defaultKeyFunc,
		MapClaims{"foo": "bar"},
		false,
		ValidationErrorSignatureInvalid,
		[]error{ErrTokenSignatureInvalid},
		&Parser{validMethods: []string{"HS256"}},
		SigningMethodRS256,
	},
	{
		"valid RSA signing method",
		"",
		defaultKeyFunc,
		MapClaims{"foo": "bar"},
		true,
		0,
		nil,
		&Parser{validMethods: []string{"RS256", "HS256"}},
		SigningMethodRS256,
	},
	{
		"JSON Number",
		"",
		defaultKeyFunc,
		MapClaims{"foo": json.Number("123.4")},
		true,
		0,
		nil,
		&Parser{useJSONNumber: true},
		SigningMethodRS256,
	},
	{
		"JSON Number - basic expired",
		"", // autogen
		defaultKeyFunc,
		MapClaims{"foo": "bar", "exp": json.Number(fmt.Sprintf("%v", time.Now().Unix()-100))},
		false,
		ValidationErrorExpired,
		[]error{ErrTokenExpired},
		&Parser{useJSONNumber: true},
		SigningMethodRS256,
	},
	{
		"JSON Number - basic nbf",
		"", // autogen
		defaultKeyFunc,
		MapClaims{"foo": "bar", "nbf": json.Number(fmt.Sprintf("%v", time.Now().Unix()+100))},
		false,
		ValidationErrorNotValidYet,
		[]error{ErrTokenNotValidYet},
		&Parser{useJSONNumber: true},
		SigningMethodRS256,
	},
	{
		"JSON Number - expired and nbf",
		"", // autogen
		defaultKeyFunc,
		MapClaims{
			"foo": "bar", "nbf": json.Number(fmt.Sprintf("%v", time.Now().Unix()+100)),
			"exp": json.Number(fmt.Sprintf("%v", time.Now().Unix()-100)),
		},
		false,
		ValidationErrorNotValidYet | ValidationErrorExpired,
		[]error{ErrTokenNotValidYet},
		&Parser{useJSONNumber: true},
		SigningMethodRS256,
	},
	{
		"SkipClaimsValidation during token parsing",
		"", // autogen
		defaultKeyFunc,
		MapClaims{"foo": "bar", "nbf": json.Number(fmt.Sprintf("%v", time.Now().Unix()+100))},
		true,
		0,
		nil,
		&Parser{useJSONNumber: true, skipClaimsValidation: true},
		SigningMethodRS256,
	},
	{
		"RFC7519 Claims",
		"",
		defaultKeyFunc,
		&RegisteredClaims{
			ExpiresAt: NewNumericDate(time.Now().Add(time.Second * 10)),
		},
		true,
		0,
		nil,
		&Parser{useJSONNumber: true},
		SigningMethodRS256,
	},
	{
		"RFC7519 Claims - single aud",
		"",
		defaultKeyFunc,
		&RegisteredClaims{
			Audience: ClaimStrings{"test"},
		},
		true,
		0,
		nil,
		&Parser{useJSONNumber: true},
		SigningMethodRS256,
	},
	{
		"RFC7519 Claims - multiple aud",
		"",
		defaultKeyFunc,
		&RegisteredClaims{
			Audience: ClaimStrings{"test", "test"},
		},
		true,
		0,
		nil,
		&Parser{useJSONNumber: true},
		SigningMethodRS256,
	},
	{
		"RFC7519 Claims - single aud with wrong type",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOjF9.8mAIDUfZNQT3TGm1QFIQp91OCpJpQpbB1-m9pA2mkHc", // { "aud": 1 }
		defaultKeyFunc,
		&RegisteredClaims{
			Audience: nil, // because of the unmarshal error, this will be empty
		},
		false,
		ValidationErrorMalformed,
		[]error{ErrTokenMalformed},
		&Parser{useJSONNumber: true},
		SigningMethodRS256,
	},
	{
		"RFC7519 Claims - multiple aud with wrong types",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsidGVzdCIsMV19.htEBUf7BVbfSmVoTFjXf3y6DLmDUuLy1vTJ14_EX7Ws", // { "aud": ["test", 1] }
		defaultKeyFunc,
		&RegisteredClaims{
			Audience: nil, // because of the unmarshal error, this will be empty
		},
		false,
		ValidationErrorMalformed,
		[]error{ErrTokenMalformed},
		&Parser{useJSONNumber: true},
		SigningMethodRS256,
	},
}

// signToken creates and returns a signed JWT token using signingMethod.
func signToken(claims Claims, signingMethod SigningMethod) string {
	var privateKey interface{}
	switch signingMethod {
	case SigningMethodRS256:
		privateKey = jwtTestRSAPrivateKey
	default:
		return ""
	}
	return MakeSampleToken(claims, signingMethod, privateKey)
}

func TestParser_Parse(t *testing.T) {

	// Iterate over test data set and run tests
	for _, data := range jwtTestData {
		t.Run(
			data.name, func(t *testing.T) {

				// If the token string is blank, use helper function to generate string
				if data.tokenString == "" {
					data.tokenString = signToken(data.claims, data.signingMethod)
				}

				// Parse the token
				var token *Token
				var ve *ValidationError
				var err error
				var parser = data.parser
				if parser == nil {
					parser = new(Parser)
				}
				// Figure out correct claims type
				switch data.claims.(type) {
				case MapClaims:
					token, err = parser.ParseWithClaims(data.tokenString, MapClaims{}, data.keyfunc)
				case *RegisteredClaims:
					token, err = parser.ParseWithClaims(data.tokenString, &RegisteredClaims{}, data.keyfunc)
				}

				// Verify result matches expectation
				if !reflect.DeepEqual(data.claims, token.Claims) {
					t.Errorf("[%v] Claims mismatch. Expecting: %v  Got: %v", data.name, data.claims, token.Claims)
				}

				if data.valid && err != nil {
					t.Errorf("[%v] Error while verifying token: %T:%v", data.name, err, err)
				}

				if !data.valid && err == nil {
					t.Errorf("[%v] Invalid token passed validation", data.name)
				}

				if (err == nil && !token.Valid) || (err != nil && token.Valid) {
					t.Errorf("[%v] Inconsistent behavior between returned error and token.Valid", data.name)
				}

				if data.errors != 0 {
					if err == nil {
						t.Errorf("[%v] Expecting error. Didn't get one.", data.name)
					} else {
						if errors.As(err, &ve) {
							// compare the bitfield part of the error
							if e := ve.Errors; e != data.errors {
								t.Errorf("[%v] Errors don't match expectation.  %v != %v", data.name, e, data.errors)
							}

							if err.Error() == errKeyFuncError.Error() && ve.Inner != errKeyFuncError {
								t.Errorf(
									"[%v] Inner error does not match expectation.  %v != %v", data.name, ve.Inner,
									errKeyFuncError,
								)
							}
						}
					}
				}

				if data.err != nil {
					if err == nil {
						t.Errorf("[%v] Expecting error(s). Didn't get one.", data.name)
					} else {
						var all = false
						for _, e := range data.err {
							all = errors.Is(err, e)
						}

						if !all {
							t.Errorf(
								"[%v] Errors don't match expectation.  %v should contain all of %v", data.name, err,
								data.err,
							)
						}
					}
				}

				if data.valid {
					if token.Signature == "" {
						t.Errorf("[%v] Signature is left unpopulated after parsing", data.name)
					}
					if !token.Valid {
						// The 'Valid' field should be set to true when invoking Parse()
						t.Errorf("[%v] Token.Valid field mismatch. Expecting true, got %v", data.name, token.Valid)
					}
				}
			},
		)
	}
}

func TestParser_ParseUnverified(t *testing.T) {

	// Iterate over test data set and run tests
	for _, data := range jwtTestData {
		// Skip test data, that intentionally contains malformed tokens, as they would lead to an error
		if data.errors&ValidationErrorMalformed != 0 {
			continue
		}

		t.Run(
			data.name, func(t *testing.T) {
				// If the token string is blank, use helper function to generate string
				if data.tokenString == "" {
					data.tokenString = signToken(data.claims, data.signingMethod)
				}

				// Parse the token
				var token *Token
				var err error
				var parser = data.parser
				if parser == nil {
					parser = new(Parser)
				}
				// Figure out correct claims type
				switch data.claims.(type) {
				case MapClaims:
					token, _, err = parser.ParseUnverified(data.tokenString, MapClaims{})
				case *RegisteredClaims:
					token, _, err = parser.ParseUnverified(data.tokenString, &RegisteredClaims{})
				}

				if err != nil {
					t.Errorf("[%v] Invalid token", data.name)
				}

				// Verify result matches expectation
				if !reflect.DeepEqual(data.claims, token.Claims) {
					t.Errorf("[%v] Claims mismatch. Expecting: %v  Got: %v", data.name, data.claims, token.Claims)
				}

				if data.valid && err != nil {
					t.Errorf("[%v] Error while verifying token: %T:%v", data.name, err, err)
				}
				if token.Valid {
					// The 'Valid' field should not be set to true when invoking ParseUnverified()
					t.Errorf("[%v] Token.Valid field mismatch. Expecting false, got %v", data.name, token.Valid)
				}
				if token.Signature != "" {
					// The signature was not validated, hence the 'Signature' field is not populated.
					t.Errorf("[%v] Token.Signature field mismatch. Expecting '', got %v", data.name, token.Signature)
				}
			},
		)
	}
}

var setPaddingTestData = []struct {
	name          string
	tokenString   string
	claims        Claims
	paddedDecode  bool
	strictDecode  bool
	signingMethod SigningMethod
	keyfunc       KeyFunc
	valid         bool
}{
	{
		name:          "Validated non-padded token with padding disabled",
		tokenString:   "",
		claims:        MapClaims{"foo": "paddedbar"},
		paddedDecode:  false,
		signingMethod: SigningMethodRS256,
		keyfunc:       defaultKeyFunc,
		valid:         true,
	},
	{
		name:          "Validated non-padded token with padding enabled",
		tokenString:   "",
		claims:        MapClaims{"foo": "paddedbar"},
		paddedDecode:  true,
		signingMethod: SigningMethodRS256,
		keyfunc:       defaultKeyFunc,
		valid:         true,
	},
	{
		name:          "Error for padded token with padding disabled",
		tokenString:   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJwYWRkZWRiYXIifQ==.20kGGJaYekGTRFf8b0TwhuETcR8lv5z2363X5jf7G1yTWVTwOmte5Ii8L8_OQbYwPoiVHmZY6iJPbt_DhCN42AeFY74BcsUhR-BVrYUVhKK0RppuzEcSlILDNeQsJDLEL035CPm1VO6Jrgk7enQPIctVxUesRgswP71OpGvJxy3j1k_J8p0WzZvRZTe1D_2Misa0UDGwnEIHhmr97fIpMSZjFxlcygQw8QN34IHLHIXMaTY1eiCf4CCr6rOS9wUeu7P3CPkmFq9XhxBT_LLCmIMhHnxP5x27FUJE_JZlfek0MmARcrhpsZS2sFhHAiWrjxjOE27jkDtv1nEwn65wMw==",
		claims:        MapClaims{"foo": "paddedbar"},
		paddedDecode:  false,
		signingMethod: SigningMethodRS256,
		keyfunc:       defaultKeyFunc,
		valid:         false,
	},
	{
		name:          "Validated padded token with padding enabled",
		tokenString:   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJwYWRkZWRiYXIifQ==.20kGGJaYekGTRFf8b0TwhuETcR8lv5z2363X5jf7G1yTWVTwOmte5Ii8L8_OQbYwPoiVHmZY6iJPbt_DhCN42AeFY74BcsUhR-BVrYUVhKK0RppuzEcSlILDNeQsJDLEL035CPm1VO6Jrgk7enQPIctVxUesRgswP71OpGvJxy3j1k_J8p0WzZvRZTe1D_2Misa0UDGwnEIHhmr97fIpMSZjFxlcygQw8QN34IHLHIXMaTY1eiCf4CCr6rOS9wUeu7P3CPkmFq9XhxBT_LLCmIMhHnxP5x27FUJE_JZlfek0MmARcrhpsZS2sFhHAiWrjxjOE27jkDtv1nEwn65wMw==",
		claims:        MapClaims{"foo": "paddedbar"},
		paddedDecode:  true,
		signingMethod: SigningMethodRS256,
		keyfunc:       defaultKeyFunc,
		valid:         true,
	},
	// DecodeStrict tests, DecodePaddingAllowed=false
	{
		name: "Validated non-padded token with padding disabled, non-strict decode, non-tweaked signature",
		tokenString: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJwYWRkZWRiYXIifQ.bI15h-7mN0f-2diX5I4ErgNQy1uM-rJS5Sz7O0iTWtWSBxY1h6wy8Ywxe5EZTEO6GiIfk7Lk-72Ex-c5aA40QKhPwWB9BJ8O_LfKpezUVBOn0jRItDnVdsk4ccl2zsOVkbA4U4QvdrSbOYMbwoRHzDXfTFpoeMWtn3ez0aENJ8dh4E1echHp5ByI9Pu2aBsvM1WVcMt_BySweCL3f4T7jNZeXDr7Txd00yUd2gdsHYPjXorOvsgaBKN5GLsWd1zIY5z-2gCC8CRSN-IJ4NNX5ifh7l-bOXE2q7szTqa9pvyE9y6TQJhNMSE2FotRce_TOPBWgGpQ-K2I7E8x7wZ8O" +
			"g",
		claims:        nil,
		paddedDecode:  false,
		strictDecode:  false,
		signingMethod: SigningMethodRS256,
		keyfunc:       defaultKeyFunc,
		valid:         true,
	},
	{
		name: "Validated non-padded token with padding disabled, non-strict decode, tweaked signature",
		tokenString: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJwYWRkZWRiYXIifQ.bI15h-7mN0f-2diX5I4ErgNQy1uM-rJS5Sz7O0iTWtWSBxY1h6wy8Ywxe5EZTEO6GiIfk7Lk-72Ex-c5aA40QKhPwWB9BJ8O_LfKpezUVBOn0jRItDnVdsk4ccl2zsOVkbA4U4QvdrSbOYMbwoRHzDXfTFpoeMWtn3ez0aENJ8dh4E1echHp5ByI9Pu2aBsvM1WVcMt_BySweCL3f4T7jNZeXDr7Txd00yUd2gdsHYPjXorOvsgaBKN5GLsWd1zIY5z-2gCC8CRSN-IJ4NNX5ifh7l-bOXE2q7szTqa9pvyE9y6TQJhNMSE2FotRce_TOPBWgGpQ-K2I7E8x7wZ8O" +
			"h",
		claims:        nil,
		paddedDecode:  false,
		strictDecode:  false,
		signingMethod: SigningMethodRS256,
		keyfunc:       defaultKeyFunc,
		valid:         true,
	},
	{
		name: "Validated non-padded token with padding disabled, strict decode, non-tweaked signature",
		tokenString: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJwYWRkZWRiYXIifQ.bI15h-7mN0f-2diX5I4ErgNQy1uM-rJS5Sz7O0iTWtWSBxY1h6wy8Ywxe5EZTEO6GiIfk7Lk-72Ex-c5aA40QKhPwWB9BJ8O_LfKpezUVBOn0jRItDnVdsk4ccl2zsOVkbA4U4QvdrSbOYMbwoRHzDXfTFpoeMWtn3ez0aENJ8dh4E1echHp5ByI9Pu2aBsvM1WVcMt_BySweCL3f4T7jNZeXDr7Txd00yUd2gdsHYPjXorOvsgaBKN5GLsWd1zIY5z-2gCC8CRSN-IJ4NNX5ifh7l-bOXE2q7szTqa9pvyE9y6TQJhNMSE2FotRce_TOPBWgGpQ-K2I7E8x7wZ8O" +
			"g",
		claims:        nil,
		paddedDecode:  false,
		strictDecode:  true,
		signingMethod: SigningMethodRS256,
		keyfunc:       defaultKeyFunc,
		valid:         true,
	},
	{
		name: "Error for non-padded token with padding disabled, strict decode, tweaked signature",
		tokenString: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJwYWRkZWRiYXIifQ.bI15h-7mN0f-2diX5I4ErgNQy1uM-rJS5Sz7O0iTWtWSBxY1h6wy8Ywxe5EZTEO6GiIfk7Lk-72Ex-c5aA40QKhPwWB9BJ8O_LfKpezUVBOn0jRItDnVdsk4ccl2zsOVkbA4U4QvdrSbOYMbwoRHzDXfTFpoeMWtn3ez0aENJ8dh4E1echHp5ByI9Pu2aBsvM1WVcMt_BySweCL3f4T7jNZeXDr7Txd00yUd2gdsHYPjXorOvsgaBKN5GLsWd1zIY5z-2gCC8CRSN-IJ4NNX5ifh7l-bOXE2q7szTqa9pvyE9y6TQJhNMSE2FotRce_TOPBWgGpQ-K2I7E8x7wZ8O" +
			"h",
		claims:        nil,
		paddedDecode:  false,
		strictDecode:  true,
		signingMethod: SigningMethodRS256,
		keyfunc:       defaultKeyFunc,
		valid:         false,
	},
}

// Extension of Parsing, this is to test out functionality specific to switching codecs with padding.
func TestSetPadding(t *testing.T) {
	for _, data := range setPaddingTestData {
		t.Run(
			data.name, func(t *testing.T) {

				// If the token string is blank, use helper function to generate string
				if data.tokenString == "" {
					data.tokenString = signToken(data.claims, data.signingMethod)
				}

				// Parse the token
				var token *Token
				var err error
				parser := new(Parser)

				// Figure out correct claims type
				token, err = parser.ParseWithClaims(data.tokenString, MapClaims{}, data.keyfunc)

				if (err == nil) != data.valid || token.Valid != data.valid {
					t.Errorf(
						"[%v] Error Parsing Token with decoding padding set to %v: %v",
						data.name,
						data.paddedDecode,
						err,
					)
				}

			},
		)
	}
}

func BenchmarkParseUnverified(b *testing.B) {

	// Iterate over test data set and run tests
	for _, data := range jwtTestData {
		// If the token string is blank, use helper function to generate string
		if data.tokenString == "" {
			data.tokenString = signToken(data.claims, data.signingMethod)
		}

		// Parse the token
		var parser = data.parser
		if parser == nil {
			parser = new(Parser)
		}
		// Figure out correct claims type
		switch data.claims.(type) {
		case MapClaims:
			b.Run(
				"map_claims", func(b *testing.B) {
					benchmarkParsing(b, parser, data.tokenString, MapClaims{})
				},
			)
		}
	}
}

// Helper method for benchmarking various parsing methods
func benchmarkParsing(b *testing.B, parser *Parser, tokenString string, claims Claims) {
	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				_, _, err := parser.ParseUnverified(tokenString, MapClaims{})
				if err != nil {
					b.Fatal(err)
				}
			}
		},
	)
}

// Helper method for benchmarking various signing methods
func benchmarkSigning(b *testing.B, method SigningMethod, key interface{}) {
	b.Helper()
	t := NewToken(method)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				if _, err := t.SignedString(key); err != nil {
					b.Fatal(err)
				}
			}
		},
	)
}

// MakeSampleToken creates and returns a encoded JWT token that has been signed with the specified cryptographic key.
func MakeSampleToken(c Claims, method SigningMethod, key interface{}) string {
	token := NewTokenWithClaims(method, c)
	s, e := token.SignedString(key)

	if e != nil {
		panic(e.Error())
	}

	return s
}
