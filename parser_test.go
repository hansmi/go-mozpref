package mozpref

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustReadFile(t *testing.T, filename string, mustExist bool) []byte {
	content, err := ioutil.ReadFile(filename)
	if err != nil && (mustExist || !os.IsNotExist(err)) {
		t.Fatalf("Reading %q: %s", filename, err)
	}
	return content
}

func readExpectedPrefs(r io.Reader) (PrefMap, error) {
	dec := json.NewDecoder(r)
	dec.UseNumber()

	var decoded map[string]struct {
		Value interface{} `json:"value"`
		Flags uint        `json:"flags"`
	}

	if err := dec.Decode(&decoded); err != nil {
		return nil, fmt.Errorf("JSON error: %s", err)
	}

	if err := dec.Decode(nil); err != io.EOF {
		return nil, fmt.Errorf("Error when EOF expected: %s", err)
	}

	result := PrefMap{}

	for key, i := range decoded {
		pref := &Pref{
			Flags: i.Flags,
		}

		if num, ok := i.Value.(json.Number); ok {
			num64, err := num.Int64()
			if err != nil {
				return nil, err
			}

			pref.Value = int(num64)
		} else {
			pref.Value = i.Value
		}

		result[key] = pref
	}

	return result, nil
}

type parserTestData struct {
	name          string
	input         []byte
	expected      PrefMap
	expectedError string
}

func runParserTests(t *testing.T, tests []parserTestData) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefs, err := parse(tt.input)

			if len(tt.expectedError) > 0 {
				if assert.Error(t, err) {
					assert.EqualValues(t, tt.expectedError, err.Error())
				}

				return
			}

			if assert.NoError(t, err) {
				assert.EqualValues(t, tt.expected, prefs)

				checkRoundTrip(t, prefs)
			}
		})
	}
}

func TestParser(t *testing.T) {
	match, err := filepath.Glob("testdata/*.prefs")
	if err != nil {
		t.Fatal(err)
	}

	var tests = []parserTestData{}

	for _, i := range match {
		if filepath.Ext(i) != ".prefs" {
			t.Fatalf("Unexpected extension: %q", i)
		}

		basename := strings.TrimSuffix(i, ".prefs")
		datafile := basename + ".json"
		errorfile := basename + ".error"

		prefsContent := mustReadFile(t, i, true)
		dataContent := mustReadFile(t, datafile, false)
		errorContent := mustReadFile(t, errorfile, false)

		var prefs PrefMap

		if len(dataContent) > 0 {
			prefs, err = readExpectedPrefs(bytes.NewReader(dataContent))
			if err != nil {
				t.Fatalf("Error in %q: %s", datafile, err)
			}
		}

		tests = append(tests, parserTestData{
			name:          basename,
			input:         prefsContent,
			expected:      prefs,
			expectedError: strings.TrimSpace(string(errorContent)),
		})
	}

	runParserTests(t, tests)
}

func TestParserNewlines(t *testing.T) {
	tests := []parserTestData{
		parserTestData{
			input:    []byte("\r"),
			expected: PrefMap{},
		},
		parserTestData{
			input:    []byte("\n"),
			expected: PrefMap{},
		},
		parserTestData{
			input:    []byte("\r\n"),
			expected: PrefMap{},
		},
		parserTestData{
			input:         []byte("\n\n\n?"),
			expectedError: "Line 4: Syntax error",
		},
		parserTestData{
			input:         []byte("\r\r\r?"),
			expectedError: "Line 4: Syntax error",
		},
		parserTestData{
			input:         []byte("\r\n\r\n\r\n?"),
			expectedError: "Line 4: Syntax error",
		},

		// Mixed newline types
		parserTestData{
			input:         []byte("\n\r\n\n\r\r\n\r\r?"),
			expectedError: "Line 8: Syntax error",
		},
	}

	runParserTests(t, tests)
}
