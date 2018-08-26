package mozpref

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkRoundTrip(t *testing.T, prefs PrefMap) {
	buf := &bytes.Buffer{}

	n, err := prefs.WriteTo(buf)
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, buf.Len(), n)

	parsed, err := parse(buf.Bytes())
	if assert.NoError(t, err) {
		assert.EqualValues(t, prefs, parsed)
	}
}

func TestWriteRoundTrip(t *testing.T) {
	input := PrefMap{
		"val-false": &Pref{
			Value: false,
		},
		"val-true": &Pref{
			Value: true,
		},
		"val-sticky": &Pref{
			Value: 900,
			Flags: Sticky,
		},
		"val-locked": &Pref{
			Value: "Hello",
			Flags: Locked,
		},
		"val-sticky-locked": &Pref{
			Value: "World",
			Flags: Locked | Sticky,
		},
		"val-negative": &Pref{
			Value: -123,
		},
		"val-int": &Pref{
			Value: int(1 << 15),
		},
	}

	checkRoundTrip(t, input)
}
