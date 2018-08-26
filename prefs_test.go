package mozpref

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFrom(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.WriteString("// Comment\npref('test', -100);")

	prefs, err := ReadFrom(buf)

	if assert.NoError(t, err) {
		assert.EqualValues(t, prefs, PrefMap{
			"test": &Pref{
				Value: -100,
			},
		})
	}
}

func TestFromAndToMap(t *testing.T) {
	input := map[string]interface{}{
		"Hello":   "World",
		"test":    false,
		"invalid": []string{},
	}

	for _, flags := range []uint{0, Sticky | Locked | UserPref} {
		prefs := FromMap(input, flags)

		assert.EqualValues(t, PrefMap{
			"Hello": &Pref{
				Value: "World",
				Flags: flags,
			},
			"test": &Pref{
				Value: false,
				Flags: flags,
			},
			"invalid": &Pref{
				Value: []string{},
				Flags: flags,
			},
		}, prefs)

		assert.EqualValues(t, input, prefs.ToMap())
	}
}
