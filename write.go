package mozpref

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
)

// Serialize a single preference entry
func (p Pref) writeTo(w io.Writer, name string) (int64, error) {
	var spec string
	var err error

	if (p.Flags & UserPref) != 0 {
		spec = "user_pref("
	} else {
		spec = "pref("
	}

	parts := []string{
		spec,
		strconv.Quote(name),
		", ",
	}

	appendPart := func(v string) {
		parts = append(parts, v)
	}

	appendInteger := func(v int64) error {
		if v < math.MinInt32 {
			return fmt.Errorf("Integer underflow (%d < %d)", v, math.MinInt32)
		}

		if v > math.MaxInt32 {
			return fmt.Errorf("Integer overflow (%d > %d)", v, math.MaxInt32)
		}

		appendPart(strconv.FormatInt(v, 10))

		return nil
	}

	err = nil

	switch v := p.Value.(type) {
	case string:
		appendPart(strconv.Quote(v))

	case int:
		err = appendInteger(int64(v))

	case int32:
		err = appendInteger(int64(v))

	case int64:
		err = appendInteger(v)

	case bool:
		if v {
			appendPart("true")
		} else {
			appendPart("false")
		}

	default:
		return 0, fmt.Errorf("Unsupported value %v (type %T)", p.Value, p.Value)
	}

	if err != nil {
		return 0, err
	}

	if (p.Flags & Sticky) != 0 {
		parts = append(parts, ", sticky")
	}

	if (p.Flags & Locked) != 0 {
		parts = append(parts, ", locked")
	}

	parts = append(parts, ");")

	var written, n int

	for _, i := range parts {
		n, err = io.WriteString(w, i)
		if err != nil {
			return 0, err
		}

		written += n
	}

	return int64(written), nil
}

// WriteTo writes all preferences in map to an io.Writer using the standard
// format, sorted by key.
func (p PrefMap) WriteTo(w io.Writer) (int64, error) {
	keys := make([]string, 0, len(p))

	for i := range p {
		keys = append(keys, i)
	}

	sort.Strings(keys)

	var err error
	var written int64
	var n int
	var n64 int64

	for _, i := range keys {
		n64, err = p[i].writeTo(w, i)
		if err != nil {
			return 0, fmt.Errorf("Pref %q: %s", i, err)
		}

		written += n64

		n, err = io.WriteString(w, "\n")
		if err != nil {
			return 0, err
		}

		written += int64(n)
	}

	return written, err
}
