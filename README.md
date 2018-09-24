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

Library                                  | Small JSON  | Medium JSON | Large JSON    | Huge JSON
:----------------------------------------|------------:|------------:|--------------:|----------------:
[`github.com/antonholmquist/jason`][]    | 25239 ns/op | 63614 ns/op | 1022959 ns/op | 9590892930 ns/op
[`github.com/bitly/go-simplejson`][]     | 13217 ns/op | 59932 ns/op | 991333 ns/op  | 5972472324 ns/op
[`github.com/ugorji/go/codec`][]         | 9877 ns/op  | 57036 ns/op | 798633 ns/op  | 7792939090 ns/op
[`github.com/Jeffail/gabs`][]            | 13063 ns/op | 54059 ns/op | 855958 ns/op  | 5822377683 ns/op
[`github.com/mreiferson/go-ujson`][]     | 11632 ns/op | 42901 ns/op | 666410 ns/op  | 3945387376 ns/op
[`github.com/json-iterator/go`][]        | 11449 ns/op | 31346 ns/op | 410673 ns/op  | 3322275736 ns/op
[`github.com/a8m/djson`][]               | 10051 ns/op | 33364 ns/op | 523898 ns/op  | 2857508424 ns/op
[`encoding/json`][] (interface)          | 12756 ns/op | 52360 ns/op | 841988 ns/op  | 5700589742 ns/op
[`encoding/json`][] (struct)             | 12587 ns/op | 43188 ns/op | 612850 ns/op  | 5452741452 ns/op
[`github.com/pquerna/ffjson`][]          | 10190 ns/op | 22878 ns/op | 250692 ns/op  | 2475434876 ns/op
[`github.com/mailru/easyjson`][]         | 8815 ns/op  | 16699 ns/op | 174619 ns/op  | 2041085219 ns/op
[`github.com/buger/jsonparser`][]        | 8126 ns/op  | 16866 ns/op | 111169 ns/op  | 1147621332 ns/op
[`github.com/darthfennec/jsonmuncher`][] | **6339 ns/op** | **14731 ns/op** | **108270 ns/op** | **859714936 ns/op**

### Memory

Library                                  | Small JSON | Medium JSON | Large JSON  | Huge JSON
:----------------------------------------|-----------:|------------:|------------:|---------------:
[`github.com/antonholmquist/jason`][]    | 8333 B/op  | 22444 B/op  | 421090 B/op | 4191279992 B/op
[`github.com/bitly/go-simplejson`][]     | 3337 B/op  | 20603 B/op  | 392650 B/op | 2563132440 B/op
[`github.com/ugorji/go/codec`][]         | 2304 B/op  | 5789 B/op   | 57458 B/op  | 2667890904 B/op
[`github.com/Jeffail/gabs`][]            | 2649 B/op  | 14440 B/op  | 265119 B/op | 1517439336 B/op
[`github.com/mreiferson/go-ujson`][]     | 2633 B/op  | 15203 B/op  | 288516 B/op | 1593325728 B/op
[`github.com/json-iterator/go`][]        | 2001 B/op  | 7614 B/op   | 118215 B/op | 1839352848 B/op
[`github.com/a8m/djson`][]               | 2345 B/op  | 13659 B/op  | 261149 B/op | 1489454784 B/op
[`encoding/json`][] (interface)          | 2521 B/op  | 13964 B/op  | 261826 B/op | 1489375432 B/op
[`encoding/json`][] (struct)             | 1912 B/op  | 4626 B/op   | 56265 B/op  | 1442500544 B/op
[`github.com/pquerna/ffjson`][]          | 1752 B/op  | 4346 B/op   | 55976 B/op  | 1442499509 B/op
[`github.com/mailru/easyjson`][]         | 1304 B/op  | 3952 B/op   | 55096 B/op  | 1510498482 B/op
[`github.com/buger/jsonparser`][]        | 1168 B/op  | 3536 B/op   | 49616 B/op  | 360846811 B/op
[`github.com/darthfennec/jsonmuncher`][] | **496 B/op** | **1264 B/op** | **4336 B/op** | **4342 B/op**

### Allocations

Library                                  | Small JSON    | Medium JSON   | Large JSON     | Huge JSON
:----------------------------------------|--------------:|--------------:|---------------:|------------------:
[`github.com/antonholmquist/jason`][]    | 104 allocs/op | 248 allocs/op | 3284 allocs/op | 49634929 allocs/op
[`github.com/bitly/go-simplejson`][]     | 39 allocs/op  | 220 allocs/op | 2845 allocs/op | 28034845 allocs/op
[`github.com/ugorji/go/codec`][]         | 12 allocs/op  | 36 allocs/op  | 254 allocs/op  | 18503643 allocs/op
[`github.com/Jeffail/gabs`][]            | 47 allocs/op  | 232 allocs/op | 3041 allocs/op | 29534833 allocs/op
[`github.com/mreiferson/go-ujson`][]     | 46 allocs/op  | 284 allocs/op | 4021 allocs/op | 34534686 allocs/op
[`github.com/json-iterator/go`][]        | 32 allocs/op  | 101 allocs/op | 1379 allocs/op | 17002144 allocs/op
[`github.com/a8m/djson`][]               | 34 allocs/op  | 201 allocs/op | 2746 allocs/op | 28035037 allocs/op
[`encoding/json`][] (interface)          | 39 allocs/op  | 213 allocs/op | 2881 allocs/op | 28034766 allocs/op
[`encoding/json`][] (struct)             | 23 allocs/op  | 30 allocs/op  | 248 allocs/op  | 10502039 allocs/op
[`github.com/pquerna/ffjson`][]          | 21 allocs/op  | 25 allocs/op  | 243 allocs/op  | 10502029 allocs/op
[`github.com/mailru/easyjson`][]         | 15 allocs/op  | 19 allocs/op  | 232 allocs/op  | 9502021 allocs/op
[`github.com/buger/jsonparser`][]        | 7 allocs/op   | 7 allocs/op   | 7 allocs/op    | 1000012 allocs/op
[`github.com/darthfennec/jsonmuncher`][] | **6 allocs/op** | **6 allocs/op** | **6 allocs/op** | **6 allocs/op**

[`github.com/darthfennec/jsonmuncher`]: https://github.com/darthfennec/jsonmuncher
[`github.com/buger/jsonparser`]: https://github.com/buger/jsonparser
[`encoding/json`]: https://golang.org/pkg/encoding/json
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

The API lives in the `jsonmuncher` package.

If the end of the stream is unexpectedly reached at any point, a
`jsonmuncher.UnexpectedEOF` error is returned.

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
than parsing. This is effective when used with `String`s, `Object`s, or
`Array`s.

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
