package mozpref

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

var errStringSyntax = errors.New("Invalid string syntax")

// Unquote interprets s as a single or double quoted string literal and returns
// the quoted string value.
//
// Derived from strconv.Unquote which is
// Copyright (c) 2009 The Go Authors. All rights reserved.
// Licence text: https://golang.org/LICENSE
func unquote(s string) (string, error) {
	n := len(s)

	if n < 2 {
		return "", errStringSyntax
	}

	quote := s[0]
	if quote != s[n-1] || (quote != '"' && quote != '\'') {
		return "", errStringSyntax
	}

	s = s[1 : n-1]

	if strings.IndexByte(s, '\n') >= 0 {
		return "", errStringSyntax
	}

	// Avoid allocation in trivial cases
	if strings.IndexByte(s, '\\') < 0 && strings.IndexByte(s, quote) < 0 {
		return s, nil
	}

	var runeTmp [utf8.UTFMax]byte

	// Try to avoid unnecessary allocations
	buf := make([]byte, 0, 3*len(s)/2)

	for len(s) > 0 {
		c, multibyte, ss, err := strconv.UnquoteChar(s, quote)
		if err != nil {
			return "", err
		}
		s = ss
		if c < utf8.RuneSelf || !multibyte {
			buf = append(buf, byte(c))
		} else {
			n := utf8.EncodeRune(runeTmp[:], c)
			buf = append(buf, runeTmp[:n]...)
		}
	}

	return string(buf), nil
}
