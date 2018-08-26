package mozpref

import (
  "fmt"
  "strconv"
)

type parsedPref struct {
  Name string
  *Pref
}

func newParsedPref(flags uint) parsedPref {
  return parsedPref{
    Pref: &Pref{
      Flags: flags,
    },
  }
}

%%{
# Reference implementation:
# https://searchfox.org/mozilla-central/source/modules/libpref/parser/src/lib.rs

  machine prefs;

  action storeValueStart {
    valueStart = p
  }

  action startPref {
    current = newParsedPref(0)
  }

  action startUserPref {
    current = newParsedPref(UserPref)
  }

  action startStickyPref {
    current = newParsedPref(Sticky)
  }

  action storePrefName {
    current.Name = strValue
  }

  action endPref {
    prefs[current.Name] = current.Pref
  }

  action storePrefValueString {
    current.Value = strValue
  }

  action storePrefValueInt {
    current.Value = intValue
  }

  action storePrefValueFalse {
    current.Value = false
  }

  action storePrefValueTrue {
    current.Value = true
  }

  action incrLineNumber {
    lineNumber++
  }

  action nextIsLinefeed { (p + 1) < pe && data[p + 1] == '\n' }

  newline = (
    '\n' |
    '\r\n' |
    ('\r' when !nextIsLinefeed)
  ) %incrLineNumber;

  any_count_line = any | newline;

  # Consume comment
  comment =
    ('/*' any_count_line* :>> '*/') |
    (('//' | '#') [^\r\n]* newline)
    ;

  # Single quote
  sliteralChar = [^'\\] | ('\\' any_count_line);
  sliteral = '\'' sliteralChar* '\'';

  # Double quote
  dliteralChar = [^"\\] | ('\\' any_count_line);
  dliteral = '"' dliteralChar* '"';

  string = (sliteral | dliteral) >storeValueStart %{
    strValue, err = unquote(string(data[valueStart:p]))
    if err != nil {
      goto fail
    }
  };

  intSign = '-' | '+';
  intLiteral = (intSign? digit+) >storeValueStart %{
    {
      var intValue64 int64

      intValue64, err = strconv.ParseInt(string(data[valueStart:p]), 10, 32)
      if err != nil {
        goto fail;
      }

      intValue = int(intValue64)
    }
  };

  whitespace = ([ \t\v\f] | newline | comment)*;

  prefValue =
    (string %storePrefValueString) |
    (intLiteral %storePrefValueInt) |
    ('false' %storePrefValueFalse) |
    ('true' %storePrefValueTrue);

  prefAttr =
    ',' whitespace
    (
      ('sticky' %{ current.Flags |= Sticky }) |
      ('locked' %{ current.Flags |= Locked })
    ) whitespace;

  pref =
    (
      ('pref' %startPref) |
      ('user_pref' %startUserPref) |
      ('sticky_pref' %startStickyPref)
    ) whitespace
    '(' whitespace
    (string %storePrefName) whitespace ',' whitespace
    prefValue whitespace
    prefAttr*
    ')' whitespace ';'
    %endPref
    ;

  main := |*
    whitespace;
    pref;
  *|;
}%%

%% write data;

func parse(data []byte) (PrefMap, error) {
  var err error
  var current parsedPref
  var intValue int
  var strValue string

  lineNumber := 1
  valueStart := 0

  prefs := make(PrefMap)

  cs := 0
  p := 0
  pe := len(data)
  eof := len(data)
  ts := 0
  te := 0
  act := 0

  _ = eof
  _ = te
  _ = ts
  _ = act

  %% write init;
  %% write exec;

  if cs == prefs_error || cs != prefs_first_final {
    err = fmt.Errorf("Syntax error")
  }

fail:
  if err != nil {
    err = fmt.Errorf("Line %d: %s", lineNumber, err)

    return nil, err
  }

  return prefs, err
}

// vim: set sw=2 sts=2 et syntax=ragel :