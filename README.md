# go-mozpref #

[![GoDoc](https://godoc.org/github.com/hansmi/go-mozpref/github?status.svg)](https://godoc.org/github.com/hansmi/go-mozpref/github)

go-mozpref is a Go library for reading and writing preference files as used by
[Mozilla Firefox][] and [Thunderbird][]. The parser consists of
a [Ragel][]-generated state machine.

[Mozilla Firefox]: https://www.mozilla.org/firefox/
[Thunderbird]: https://www.thunderbird.net/
[Ragel]: https://www.colm.net/open-source/ragel/


## Usage ##

```go
import mozpref "github.com/hansmi/go-mozpref"
```

Preferences can be read from any `io.Reader`:

```go
file, err := os.Open("prefs.js")
if err != nil {
	// Handle error
}

prefs, err := mozpref.ReadFrom(file)
if err != nil {
	// Handle error
}

for name, p := range prefs {
	fmt.Printf("%s = %s\n", name, p.Value)
}
```

Writing works likewise with any `io.Writer`:

```go
prefs := mozpref.PrefMap{
	"example": mozpref.Pref{
		Value: true,
		Flags: mozpref.Locked,
	},
}

_, err := prefs.WriteTo(os.Stdout)
if err != nil {
	// Handle error
}
```

Preferences can be marked as locked or sticky:

```
prefs["example"].Flags |= mozpref.Sticky
prefs["other"] = &mozpref.Pref{
	Value: "Hello World",
	Flags: mozpref.Locked | mozpref.Sticky,
}
```

Entries suitable for `user.js` files can be written using the
`mozpref.UserPref` flag.


## Versioning ##

go-mozpref follows [semver](https://semver.org/).


## License ##

This library is distributed under the BSD-style license found in the
[LICENSE](./LICENSE) file.
