[![GoDoc](https://godoc.org/github.com/DarthFennec/jsonmuncher?status.svg)](https://godoc.org/github.com/DarthFennec/jsonmuncher)
[![Go Report Card](https://goreportcard.com/badge/github.com/DarthFennec/jsonmuncher)](https://goreportcard.com/report/github.com/DarthFennec/jsonmuncher)
[![License](https://img.shields.io/github/license/DarthFennec/jsonmuncher.svg)](https://github.com/DarthFennec/jsonmuncher/blob/master/LICENSE)

JSON Muncher
============

A highly efficient streaming JSON parser for Go.

But why though?
---------------

Do we really need yet another JSON parser? There are a dozen or so other
projects that solve this problem, would one of those not suffice?

Different situations call for different approaches. Each of these projects
exists because someone found the existing solutions lacking in some way, a gap
that needed filling. This might be related to speed, memory footprint, ease of
use, some combination of these, or something else entirely. I've found that none
of the existing parsers fill my particular gap.

JSON Muncher is designed to be as fast as possible, but its primary focus is
memory efficiency. It employs the following design concepts:

- **Interactive**. Each step of the parse is explicitly triggered by the caller.
  In some ways this might be seen as detrimental, as it means a little more code
  is often needed to parse a file than is required by other parsers. However, it
  also means you have finer control over how the parse progresses, which can
  allow for drastic efficiency improvements.
- **Streaming**. Rather than load the entire file into memory at once, JSON
  Muncher reads only what it needs from the input stream. This heavily reduces
  the memory footprint, and is especially helpful when parsing very large files.
- **No memory allocation** (almost). Allocating memory can be costly, in time as
  well as in space. JSON Muncher avoids allocations as much as possible, without
  sacrificing usability. However, sometimes allocations are necessary.
  Allocations are only made in these extremely limited cases:
  - Two or three allocations are made to initialize the input buffer. This only
    happens once, when the parse begins.
  - An error is allocated if there is a problem parsing the stream. This happens
    at most once (usually not at all).
  - If a numeric literal in the JSON stream is too long, a temporary buffer is
    allocated to store it during parsing. This only happens if the literal
    exceeds 32 characters in length, which is extremely unlikely in practice.

Performance
-----------

If you want to see the raw benchmarks or run them yourself,
[go here](https://github.com/darthfennec/jsonmuncher/tree/master/benchmark/).

### Speed

Library                                   | Small JSON     | Medium JSON     | Large JSON      | Huge JSON
:-----------------------------------------|---------------:|----------------:|----------------:|-------------------:
[`github.com/antonholmquist/jason`][]     | 24120 ns/op    | 63887 ns/op     | 1044236 ns/op   | 9854035703 ns/op
[`github.com/bcicen/jstream`][]           | 40813 ns/op    | 84099 ns/op     | 1969405 ns/op   | 5866067676 ns/op
[`github.com/bitly/go-simplejson`][]      | 12923 ns/op    | 59926 ns/op     | 1000596 ns/op   | 5876533891 ns/op
[`github.com/ugorji/go/codec`][]          | 9467 ns/op     | 56543 ns/op     | 798903 ns/op    | 7739377297 ns/op
[`github.com/jeffail/gabs`][]             | 12494 ns/op    | 54700 ns/op     | 870205 ns/op    | 5661634428 ns/op
[`github.com/mreiferson/go-ujson`][]      | 10737 ns/op    | 41953 ns/op     | 679211 ns/op    | 4026916938 ns/op
[`github.com/json-iterator/go`][]         | 10582 ns/op    | 30852 ns/op     | 410814 ns/op    | 3225396334 ns/op
[`github.com/a8m/djson`][]                | 9475 ns/op     | 33466 ns/op     | 531208 ns/op    | 3048507867 ns/op
[`encoding/json`][] (interface streaming) | 9651 ns/op     | 55407 ns/op     | 931944 ns/op    | 5966943523 ns/op
[`encoding/json`][] (struct streaming)    | 9257 ns/op     | 44838 ns/op     | 655320 ns/op    | 5735771777 ns/op
[`encoding/json`][] (interface)           | 12018 ns/op    | 55449 ns/op     | 842066 ns/op    | 5701351770 ns/op
[`encoding/json`][] (struct)              | 11285 ns/op    | 42425 ns/op     | 606856 ns/op    | 5435946960 ns/op
[`github.com/francoispqt/gojay`][]        | 7603 ns/op     | 15310 ns/op     | 153690 ns/op    | 2090216626 ns/op
[`github.com/pquerna/ffjson`][]           | 9163 ns/op     | 21394 ns/op     | 248859 ns/op    | 2415598919 ns/op
[`github.com/mailru/easyjson`][]          | 7948 ns/op     | 15691 ns/op     | 175398 ns/op    | 2051049192 ns/op
[`github.com/buger/jsonparser`][]         | 7322 ns/op     | 16174 ns/op     | 111023 ns/op    | 1135070002 ns/op
[`github.com/darthfennec/jsonmuncher`][]  | **5937 ns/op** | **13783 ns/op** | **94460 ns/op** | **761287513 ns/op**

### Memory

Library                                   | Small JSON   | Medium JSON   | Large JSON    | Huge JSON
:-----------------------------------------|-------------:|--------------:|--------------:|---------------:
[`github.com/antonholmquist/jason`][]     | 8333 B/op    | 22443 B/op    | 421071 B/op   | 4191166648 B/op
[`github.com/bcicen/jstream`][]           | 13289 B/op   | 14713 B/op    | 438465 B/op   | 1129458008 B/op
[`github.com/bitly/go-simplejson`][]      | 3337 B/op    | 20603 B/op    | 392635 B/op   | 2563080408 B/op
[`github.com/ugorji/go/codec`][]          | 2304 B/op    | 5789 B/op     | 57458 B/op    | 2667890632 B/op
[`github.com/jeffail/gabs`][]             | 2649 B/op    | 14440 B/op    | 265079 B/op   | 1517427480 B/op
[`github.com/mreiferson/go-ujson`][]      | 2633 B/op    | 15203 B/op    | 288540 B/op   | 1593388936 B/op
[`github.com/json-iterator/go`][]         | 2001 B/op    | 7615 B/op     | 118218 B/op   | 1839351112 B/op
[`github.com/a8m/djson`][]                | 2345 B/op    | 13659 B/op    | 261144 B/op   | 1489389136 B/op
[`encoding/json`][] (interface streaming) | 2217 B/op    | 17036 B/op    | 341692 B/op   | 2214184568 B/op
[`encoding/json`][] (struct streaming)    | 1608 B/op    | 7692 B/op     | 136168 B/op   | 2167391392 B/op
[`encoding/json`][] (interface)           | 2521 B/op    | 13964 B/op    | 261799 B/op   | 1489402024 B/op
[`encoding/json`][] (struct)              | 1912 B/op    | 4626 B/op     | 56264 B/op    | 1442501104 B/op
[`github.com/francoispqt/gojay`][]        | 1520 B/op    | 6474 B/op     | 102668 B/op   | 1911409520 B/op
[`github.com/pquerna/ffjson`][]           | 1752 B/op    | 4346 B/op     | 55977 B/op    | 1442499717 B/op
[`github.com/mailru/easyjson`][]          | 1304 B/op    | 3952 B/op     | 55096 B/op    | 1510499090 B/op
[`github.com/buger/jsonparser`][]         | 1168 B/op    | 3536 B/op     | 49616 B/op    | 360846475 B/op
[`github.com/darthfennec/jsonmuncher`][]  | **496 B/op** | **1264 B/op** | **4336 B/op** | **4336 B/op**

### Allocations

Library                                   | Small JSON      | Medium JSON     | Large JSON      | Huge JSON
:-----------------------------------------|----------------:|----------------:|----------------:|------------------:
[`github.com/antonholmquist/jason`][]     | 104 allocs/op   | 248 allocs/op   | 3284 allocs/op  | 49634480 allocs/op
[`github.com/bcicen/jstream`][]           | 40 allocs/op    | 172 allocs/op   | 5484 allocs/op  | 28533146 allocs/op
[`github.com/bitly/go-simplejson`][]      | 39 allocs/op    | 220 allocs/op   | 2845 allocs/op  | 28034667 allocs/op
[`github.com/ugorji/go/codec`][]          | 12 allocs/op    | 36 allocs/op    | 254 allocs/op   | 18503643 allocs/op
[`github.com/jeffail/gabs`][]             | 47 allocs/op    | 232 allocs/op   | 3041 allocs/op  | 29534794 allocs/op
[`github.com/mreiferson/go-ujson`][]      | 46 allocs/op    | 284 allocs/op   | 4021 allocs/op  | 34534906 allocs/op
[`github.com/json-iterator/go`][]         | 32 allocs/op    | 101 allocs/op   | 1379 allocs/op  | 17002141 allocs/op
[`github.com/a8m/djson`][]                | 34 allocs/op    | 201 allocs/op   | 2746 allocs/op  | 28034807 allocs/op
[`encoding/json`][] (interface streaming) | 38 allocs/op    | 217 allocs/op   | 2889 allocs/op  | 28034497 allocs/op
[`encoding/json`][] (struct streaming)    | 22 allocs/op    | 34 allocs/op    | 256 allocs/op   | 10502057 allocs/op
[`encoding/json`][] (interface)           | 39 allocs/op    | 213 allocs/op   | 2881 allocs/op  | 28034854 allocs/op
[`encoding/json`][] (struct)              | 23 allocs/op    | 30 allocs/op    | 248 allocs/op   | 10502039 allocs/op
[`github.com/francoispqt/gojay`][]        | 13 allocs/op    | 20 allocs/op    | 178 allocs/op   | 10002027 allocs/op
[`github.com/pquerna/ffjson`][]           | 21 allocs/op    | 25 allocs/op    | 243 allocs/op   | 10502030 allocs/op
[`github.com/mailru/easyjson`][]          | 15 allocs/op    | 19 allocs/op    | 232 allocs/op   | 9502023 allocs/op
[`github.com/buger/jsonparser`][]         | 7 allocs/op     | 7 allocs/op     | 7 allocs/op     | 1000012 allocs/op
[`github.com/darthfennec/jsonmuncher`][]  | **6 allocs/op** | **6 allocs/op** | **6 allocs/op** | **6 allocs/op**

[`github.com/darthfennec/jsonmuncher`]: https://github.com/darthfennec/jsonmuncher
[`github.com/buger/jsonparser`]: https://github.com/buger/jsonparser
[`encoding/json`]: https://golang.org/pkg/encoding/json
[`github.com/bcicen/jstream`]: https://github.com/bcicen/jstream
[`github.com/francoispqt/gojay`]: https://github.com/francoispqt/gojay
[`github.com/json-iterator/go`]: https://github.com/json-iterator/go
[`github.com/Jeffail/gabs`]: https://github.com/Jeffail/gabs
[`github.com/bitly/go-simplejson`]: https://github.com/bitly/go-simplejson
[`github.com/pquerna/ffjson`]: https://github.com/pquerna/ffjson
[`github.com/antonholmquist/jason`]: https://github.com/antonholmquist/jason
[`github.com/mreiferson/go-ujson`]: https://github.com/mreiferson/go-ujson
[`github.com/a8m/djson`]: https://github.com/a8m/djson
[`github.com/ugorji/go/codec`]: https://github.com/ugorji/go/tree/master/codec
[`github.com/mailru/easyjson`]: https://github.com/mailru/easyjson

API Reference
-------------

[GoDoc](https://godoc.org/github.com/DarthFennec/jsonmuncher)

The API lives in the `jsonmuncher` package.

### `Parse()`

``` go
func Parse(r io.Reader, size int) (JsonValue, error)
```

The `Parse` function starts parsing a stream. The stream is passed as the first
argument, and can be anything that implements the `io.Reader` interface. The
second argument is the size of the read buffer, in bytes.

This function will allocate a read buffer of the appropriate size, read a chunk
of the stream into it, and start parsing. It returns a `JsonValue`, which can be
used to further parse the stream.

The buffer size can be anything, but a power of two is recommended, to avoid
potential hardware-related slowness. `4096` (4KiB) or `8192` (8KiB) are
generally both good buffer sizes.

### `JsonValue`

This struct represents a value parsed from the stream. It has two exported
fields:

``` go
type JsonValue struct {
    Type   JsonType
    Status JsonStatus
}
```

- `Type` describes what kind of data the `JsonValue` contains. Its value might
  be `Null`, `Bool`, `Number`, `String`, `Array`, or `Object`.
- `Status` describes the read status of the `JsonValue`. Its value is one of the
  following:
  - `Incomplete`: There was a problem parsing the value, or some other error was
    encountered. As long as you aren't ignoring errors, you shouldn't see this.
  - `Working`: The value is partially parsed, and can be parsed further by
    reading more from the stream.
  - `Complete`: The value was parsed in its entirety, and can no longer read
    from the stream.

The rest of the API consists of methods on the `JsonValue` struct.

### `ValueBool()`

``` go
func (data *JsonValue) ValueBool() (bool, error)
```

If this `JsonValue` is a `Bool`, return the value (`true` or `false`).
Otherwise, return an error.

### `ValueNum()`

``` go
func (data *JsonValue) ValueNum() (float64, error)
```

If this `JsonValue` is a `Number`, return the value as a double-precision float.
Otherwise, return an error.

### `Read()`

``` go
func (data *JsonValue) Read(b []byte) (int, error)
```

An implementation of the `io.Reader` interface. If this `JsonValue` is a
`String`, write its contents into the argument slice, and return the number of
bytes written. Otherwise, return an error.

### `NextKey()`

``` go
func (data *JsonValue) NextKey() (JsonValue, error)
```

If this `JsonValue` is an `Object`, read to the next key and return it.
Otherwise, return an error.

### `NextValue()`

``` go
func (data *JsonValue) NextValue() (JsonValue, error)
```

If this `JsonValue` is an `Object` or `Array`, read to the next value and return
it. Otherwise, return an error.

In the case of an `Array`, each time `NextValue()` is called, the next item in
the array will be returned. In the case of an `Object`, alternate between
calling `NextKey()` and `NextValue()` to get each key/value pair in the object,
or else data will be skipped (if you miss calling a `NextKey()` or
`NextValue()`, that key or value is discarded).

When `NextKey()` or `NextValue()` is called but the object or array has been
read to the end, a `jsonmuncher.EndOfValue` error is returned.

### `Close()`

``` go
func (data *JsonValue) Close() error
```

An implementation of the `io.Closer` interface. When called, this "closes" the
`JsonValue` by simply discarding the remainder of the value from the stream. Use
this if you don't want the rest of the value's data, as it's a good deal faster
than parsing. This is effective when used with `String`s, `Number`s, `Object`s,
or `Array`s.

`NextKey()`, `NextValue()`, and `Close()` cannot be called on an `Object` or
`Array` that has a partially-parsed child; you must fully parse or `Close()` a
child value before continuing to read its parent. Otherwise, an error is
returned.

### `Compare()`

``` go
func (data *JsonValue) Compare(vals ...string) (string, bool, error)
```

A helper function, designed to check the value of a `String` without allocating
memory. Consumes the string in the process. Given one or more strings as
arguments, compare the value against each argument. Return `true` with the
matched string if there was an exact match, and return `false` otherwise.

### `FindKey()`

``` go
func (data *JsonValue) FindKey(keys ...string) (string, JsonValue, bool, error)
```

A helper function, designed to quickly search an `Object` without allocating
memory. Given one or more strings as arguments, check each key in the object and
compare it to each argument. When an exact match is found, read the
corresponding value from the object, and return `true` with the matched string
and the value. If the object is exhausted and no match is found, return `false`.

To compare against multiple strings while avoiding allocations, `Compare()` and
`FindKey()` both must sort their arguments. Both functions perform an in-place
sort, but this step is faster if the arguments are already ordered correctly.
For this reason, it's recommended that when these functions are used, arguments
are passed in alphabetical order. Also, if a slice of arguments is passed using
something like `jsonval.Compare(sliceval ...)`, the slice will be sorted
in-place by the function, so it may not be in the same order after the function
runs.
