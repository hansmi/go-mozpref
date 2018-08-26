package mozpref

import (
	"bytes"
	"io"
)

const (
	// Sticky preferences are retained in the configuration even when they
	// match the application default.
	Sticky = 1 << iota

	// Locked preferences can't be changed in the application user
	// interface.
	Locked

	// UserPref are used in "user.js"
	UserPref
)

// Pref holds the value of a preference and associated flags.
type Pref struct {
	// Value must be of type string, int or boolean
	Value interface{}

	// Bitfield with flags
	Flags uint
}

// PrefMap is a collection of preferences.
type PrefMap map[string]*Pref

// ReadFrom reads preferences from given reader and returns a map with the
// parsed values.
func ReadFrom(r io.Reader) (PrefMap, error) {
	buf := &bytes.Buffer{}

	if _, err := io.Copy(buf, r); err != nil {
		return nil, err
	}

	return parse(buf.Bytes())
}

// Parse reads preferences from a byte slice.
func Parse(b []byte) (PrefMap, error) {
	return parse(b)
}

// FromMap copies all entries in a map[string]interface{} into a PrefMap. All
// entries are given the same flags.
func FromMap(prefs map[string]interface{}, flags uint) PrefMap {
	result := PrefMap{}

	for key, value := range prefs {
		result[key] = &Pref{
			Value: value,
			Flags: flags,
		}
	}

	return result
}

// ToMap copies all entries in a PrefMap to a map[string]interface{}.
// Individual entry flags are ignored.
func (p PrefMap) ToMap() map[string]interface{} {
	result := make(map[string]interface{}, len(p))

	for key, i := range p {
		result[key] = i.Value
	}

	return result
}
