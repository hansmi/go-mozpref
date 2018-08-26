package mozpref

import (
	"bytes"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrefWriteToError(t *testing.T) {
	var tests = []struct {
		name     string
		pref     Pref
		expected string
	}{
		{"int-overflow", Pref{
			Value: 1 << 31,
		}, "Integer overflow (2147483648 > 2147483647)"},
		{"int-underflow", Pref{
			Value: -(1<<31 + 1),
		}, "Integer underflow (-2147483649 < -2147483648)"},

		{"unsupported-float64", Pref{
			Value: float64(0.0),
		}, "Unsupported value 0 (type float64)"},

		{"unsupported-string-slice", Pref{
			Value: []string{},
		}, "Unsupported value [] (type []string)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			_, err := tt.pref.writeTo(buf, tt.name)

			if assert.Error(t, err) {
				assert.Equal(t, tt.expected, err.Error())
			}
		})
	}
}

func TestPrefWriteTo(t *testing.T) {
	var tests = []struct {
		name     string
		pref     Pref
		expected string
	}{
		{"string", Pref{
			Value: "Hello World",
		}, `pref("string", "Hello World");`},

		{"bool-false", Pref{
			Value: false,
		}, `pref("bool-false", false);`},
		{"bool-true", Pref{
			Value: true,
		}, `pref("bool-true", true);`},

		{"int-zero", Pref{
			Value: 0,
		}, `pref("int-zero", 0);`},
		{"int", Pref{
			Value: 1000,
		}, `pref("int", 1000);`},

		{"int32-min", Pref{
			Value: math.MinInt32,
		}, `pref("int32-min", -2147483648);`},
		{"int32-max", Pref{
			Value: math.MaxInt32,
		}, `pref("int32-max", 2147483647);`},

		{"int32", Pref{
			Value: int32(1),
		}, `pref("int32", 1);`},
		{"int32-neg", Pref{
			Value: int32(-1),
		}, `pref("int32-neg", -1);`},

		{"int64", Pref{
			Value: int64(1),
		}, `pref("int64", 1);`},
		{"int64-neg", Pref{
			Value: int64(-1),
		}, `pref("int64-neg", -1);`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			_, err := tt.pref.writeTo(buf, tt.name)

			if assert.NoError(t, err) {
				assert.Equal(t, tt.expected, buf.String())
			}
		})
	}
}

func TestPrefMapWriteError(t *testing.T) {
	var tests = []struct {
		name     string
		prefs    PrefMap
		expected string
	}{
		{"", PrefMap{
			"a-first": &Pref{
				Value: true,
			},
			"b-testname": &Pref{
				Value: map[string]string{},
			},
		}, `Pref "b-testname": Unsupported value map[] (type map[string]string)`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			_, err := tt.prefs.WriteTo(buf)

			if assert.Error(t, err) {
				assert.Equal(t, tt.expected, err.Error())
			}

		})
	}
}

func TestPrefMapWrite(t *testing.T) {
	var tests = []struct {
		name     string
		prefs    PrefMap
		expected string
	}{
		{"empty", PrefMap{}, ""},
		{"", PrefMap{
			"a-bool": &Pref{
				Value: true,
			},
			"b-int": &Pref{
				Value: -987,
			},
			"c-string": &Pref{
				Value: "Foobar",
			},
		}, "pref(\"a-bool\", true);\npref(\"b-int\", -987);\npref(\"c-string\", \"Foobar\");\n"},
		{"sticky", PrefMap{
			"a": &Pref{
				Value: false,
				Flags: Sticky,
			},
		}, "pref(\"a\", false, sticky);\n"},
		{"locked", PrefMap{
			"a": &Pref{
				Value: true,
				Flags: Locked,
			},
		}, "pref(\"a\", true, locked);\n"},
		{"sticky-locked", PrefMap{
			"a": &Pref{
				Value: -1,
				Flags: Locked | Sticky,
			},
		}, "pref(\"a\", -1, sticky, locked);\n"},
		{"user-sticky-locked", PrefMap{
			"a": &Pref{
				Value: -1,
				Flags: Locked | Sticky | UserPref,
			},
		}, "user_pref(\"a\", -1, sticky, locked);\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			_, err := tt.prefs.WriteTo(buf)

			if assert.NoError(t, err) {
				assert.EqualValues(t, tt.expected, buf.String())
			}
		})
	}
}
