package jwt

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"unicode"
)

const reclen = 170
const recs = 20000

// Lifted and adapted from strings_test.go.
// The original 1MB single-row input is unrealistic workload for Fields.
// Replace it with smaller records of differing content.
var makeFieldsInput = func() []string {
	ret := make([]string, recs)
	for r := range ret {
		x := make([]byte, reclen)
		// Input is ~5% space, ~5% 2-byte UTF-8, rest ASCII non-space.
		for i := range x {
			switch rand.Intn(20) {
			case 0:
				x[i] = ' '
			case 1:
				if i > 0 && x[i-1] == 'x' {
					copy(x[i-1:], "χ")
					break
				}
				fallthrough
			default:
				x[i] = 'x'
			}
		}
		ret[r] = string(x)
	}
	return ret
}

var fieldsInput = makeFieldsInput()

func BenchmarkFields(b *testing.B) {
	bs := []struct {
		name string
		f    func(s string)
	}{
		{"Fields", func(s string) { strings.Fields(s) }},
		{"FieldsFuncUnicodeIsSpace", func(s string) { strings.FieldsFunc(s, unicode.IsSpace) }},
		{"FieldsFuncLatin1Switch", func(s string) { strings.FieldsFunc(s, isLatin1SpaceSwitch) }},
		{"FieldsFuncLatin1If", func(s string) { strings.FieldsFunc(s, isLatin1SpaceIf) }}, // Why so slow?

		{"FieldsFuncAltUnicodeIsSpace", func(s string) { FieldsFuncAlt(s, unicode.IsSpace) }},
		{"FieldsFuncAltLatin1Switch", func(s string) { FieldsFuncAlt(s, isLatin1SpaceSwitch) }},
		{"FieldsFuncAltLatin1If", func(s string) { FieldsFuncAlt(s, isLatin1SpaceIf) }}, // Why so slow?

		// This routine is significantly faster than the others, but it does
		// a simpler job. It's here to provide a bound of sorts on possible speedup.
		{"Split", func(s string) { strings.Split(s, " ") }},
	}
	for _, v := range bs {
		bench := func(b *testing.B) {
			f := v.f
			b.SetBytes(reclen)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				f(fieldsInput[i%recs])
			}
		}
		b.Run(v.name, bench)

	}
}

// Like unicode.IsSpace, but only supports Latin-1 spaces. On the off chance
// that they differ in inlining.
func isLatin1SpaceSwitch(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
		return true
	}
	return false
}

// Equivalent to isLatin1SpaceSwitch, but apparently generates much slower code.
// 3095 ns/op vs 2215 ns/op on my machine.
// go version devel +214be5b302 Sun Mar 26 04:40:20 2017 +0000 linux/amd64
func isLatin1SpaceIf(r rune) bool {
	return r == '\t' ||
		r == '\n' ||
		r == '\v' ||
		r == '\f' ||
		r == '\r' ||
		r == ' ' ||
		r == 0x85 ||
		r == 0xA0
}

func FieldsFuncAlt(s string, f func(rune) bool) []string {
	type span struct {
		from, to int
	}
	spans := make([]span, 0, 16)

	// First find the fields.
	from := 0
	inField := false
	for i, r := range s {
		wasInField := inField
		inField = !f(r)
		if wasInField && !inField {
			spans = append(spans, span{from: from, to: i})
		}
		if inField {
			if !wasInField {
				from = i
			}
		}
	}
	if inField {
		spans = append(spans, span{from: from, to: len(s)})
	}

	// Now copy them.
	a := make([]string, len(spans))
	for i, sp := range spans {
		a[i] = s[sp.from:sp.to]
	}
	return a
}

func TestFieldsFuncAlt(t *testing.T) {
	for _, v := range fieldsInput {
		got := FieldsFuncAlt(v, unicode.IsSpace)
		want := strings.Fields(v)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("FieldsFuncAlt(%q)=%q, want=%q", v, got, want)
		}
	}
}
