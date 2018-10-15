// Package jsonmuncher is a very performant streaming JSON parser.
package jsonmuncher

import (
	"io"
	"strconv"
	"unicode/utf8"
	"unsafe"
)

// JsonType represents the data type of a JsonValue.
type JsonType byte

const (
	// Null values are always 'null'.
	Null JsonType = iota
	// Bool values are 'true' or 'false'.
	Bool
	// Number values are double-precision floating point numbers.
	Number
	// String values are unicode strings.
	String
	// Array values are ordered collections of arbitrary JSON values.
	Array
	// Object values are maps from strings to arbitrary JSON values.
	Object
)

// JsonStatus represents the current read status of a JsonValue.
type JsonStatus byte

const (
	// Incomplete means there was a read or parse error while parsing the value.
	Incomplete JsonStatus = iota
	// Working means the value is currently in the process of being parsed.
	Working
	// Complete means the value has been parsed successfully in its entirety.
	Complete
)

// buffer is a read buffer for the JSON parser.
type buffer struct {
	data    []byte
	stream  io.Reader
	foffs   uint64
	err     error
	readerr error
	offs    uint32
	erroffs uint32
	depth   uint32
	escapes byte
	escape1 byte
	escape2 byte
	escape3 byte
	escape4 byte
	curr    byte
}

// JsonValue represents a JSON value. This is the primary structure used in this
// library.
type JsonValue struct {
	// buffer is a pointer to the read buffer.
	buffer *buffer
	// numval is the parsed value, assuming this is a Number.
	numval float64
	// depth is the nesting depth of this value.
	depth uint32
	// Type is the data type of this value.
	Type JsonType
	// Status is the read status of this value.
	Status JsonStatus
	// boolval is the parsed value, assuming this is a Bool. If this is an
	// Object or Array, whether the first element has been parsed yet.
	boolval bool
	// keynext (assuming this is an Object) is true if the next thing to read is
	// a key, false if it's a value.
	keynext bool
}

// noescape prevents escape to the heap (unsafe, use with caution)
// only used to avoid heap allocations when we know they're not necessary
//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

// feedq, feed, and next should all be the same function, but they've been
// separated for the sake of performance. feedq and next can be inlined, and so
// will run much more quickly than they otherwise would. These functions should
// be called like:
//   _ = feedq(buf) && feed(buf)
//   next(buf)

// feedq checks if the next byte is beyond the end of the buffer. Can be
// inlined; separated from `feed' to make inlining possible.
func feedq(buf *buffer) bool {
	return int(buf.offs) >= len(buf.data)
}

// feed feeds the buffer with the next chunk, assuming feedq is true. Cannot be
// inlined, because of the call to Read().
func feed(buf *buffer) bool {
	buf.foffs += uint64(len(buf.data))
	var erroffs, readoffs int
	var readerr error
	for erroffs < len(buf.data) && readerr == nil {
		readoffs, readerr = buf.stream.Read(buf.data[erroffs:])
		erroffs += readoffs
	}
	buf.readerr = readerr
	buf.erroffs = uint32(erroffs)
	buf.err = nil
	buf.offs = 0
	return false
}

// next consumes the next byte from the input stream, storing it in the
// lookahead. Can be inlined; separated from the above to make inlining
// possible.
func next(buf *buffer, e ...byte) {
	if buf.offs < buf.erroffs {
		buf.curr = buf.data[buf.offs]
		buf.offs++
	} else {
		buf.curr = 0
		buf.err = buf.readerr
	}
}

// foffs calculates a file offset based on the buffer.
func foffs(buf *buffer) uint64 {
	return buf.foffs - uint64(len(buf.data)) + uint64(buf.offs) - 1
}

// newErrUnexpected is a slightly easier way to make an ErrUnexpectedChar.
func newErrUnexpected(buf *buffer, e ...byte) ErrUnexpectedChar {
	if buf.err == io.EOF {
		return newErrUnexpectedEOF(1+foffs(buf), e...)
	}
	return newErrUnexpectedChar(foffs(buf), buf.curr, e...)
}

// skipSpace skips whitespace until the next significant character.
func skipSpace(buf *buffer) (byte, error) {
	if buf.err != nil && buf.err != io.EOF {
		return 0, buf.err
	}
	c := buf.curr
	for {
		switch c {
		case ' ', '\t', '\r', '\n':
			_ = feedq(buf) && feed(buf)
			next(buf)
			if buf.err != nil && buf.err != io.EOF {
				return 0, buf.err
			}
			c = buf.curr
		default:
			return c, nil
		}
	}
}

// readKeyword reads a boolean or null value from the stream.
func readKeyword(buf *buffer) (JsonValue, error) {
	var kw string
	var typ = Bool
	var val = false
	switch buf.curr {
	case 'n':
		kw = "null"
		typ = Null
	case 't':
		kw = "true"
		val = true
	case 'f':
		kw = "false"
	}
	for i := 1; i < len(kw); i++ {
		_ = feedq(buf) && feed(buf)
		next(buf)
		if buf.err != nil && buf.err != io.EOF {
			return JsonValue{}, buf.err
		} else if kw[i] != buf.curr {
			return JsonValue{}, newErrUnexpected(buf, kw[i])
		}
	}
	_ = feedq(buf) && feed(buf)
	next(buf)
	return JsonValue{buf, 0, buf.depth + 1, typ, Complete, val, false}, nil
}

// readInt is a special case of readNumber, and is designed to parse integers.
// This parses integers in about half the time compared to strconv.ParseInt().
func readInt(buf *buffer, sl []byte) (JsonValue, error) {
	var val int64
	neg := false
	idx := 0
	if sl[0] == '-' {
		if len(sl) == 1 {
			return JsonValue{}, newErrUnexpected(buf,
				'0', '1', '2', '3', '4', '5', '6', '7', '8', '9')
		}
		neg = true
		idx = 1
	}
	for ; idx < len(sl); idx++ {
		if sl[idx] <= '9' && sl[idx] >= '0' {
			val = 10*val + int64(sl[idx]-'0')
		} else {
			offs := foffs(buf) + uint64(1+idx-len(sl))
			return JsonValue{}, newErrUnexpectedChar(offs, sl[idx],
				'0', '1', '2', '3', '4', '5', '6', '7', '8', '9')
		}
	}
	if neg {
		val = -val
	}
	return JsonValue{buf, float64(val), buf.depth + 1, Number, Complete, false, false}, nil
}

// readNumber reads a numeric value from the stream.
func readNumber(buf *buffer) (JsonValue, error) {
	var b [32]byte
	sl := b[:0]
	simple := true
L:
	for {
		if buf.err != nil {
			if buf.err == io.EOF {
				break L
			}
			return JsonValue{}, buf.err
		}
		switch buf.curr {
		case '+', '.', 'e', 'E':
			simple = false
			fallthrough
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			sl = append(sl, buf.curr)
			_ = feedq(buf) && feed(buf)
			next(buf)
		default:
			break L
		}
	}
	if simple && len(sl) < 19 {
		return readInt(buf, sl)
	}
	// strconv.ParseFloat takes a string, but we only have a []byte.
	// Converting to string requires a new alloc and a copy, plus an escape
	// to heap for the underlying array. This line does an unsafe cast and
	// sidesteps escape analysis, avoiding those expensive extra steps.
	f, err := strconv.ParseFloat(*(*string)(noescape(unsafe.Pointer(&sl))), 64)
	if err != nil {
		return JsonValue{}, err
	}
	return JsonValue{buf, f, buf.depth + 1, Number, Complete, false, false}, nil
}

// readStream reads a string, array, or object from the stream.
func readStream(buf *buffer) (JsonValue, error) {
	var typ JsonType
	switch buf.curr {
	case '"':
		typ = String
		_ = feedq(buf) && feed(buf)
		next(buf)
	case '{':
		typ = Object
	case '[':
		typ = Array
	}
	buf.depth++
	return JsonValue{buf, 0, buf.depth, typ, Working, false, typ == Object}, nil
}

// readValue reads any value from the stream.
func readValue(buf *buffer) (JsonValue, error) {
	_, err := skipSpace(buf)
	if err != nil {
		return JsonValue{}, err
	}
	switch buf.curr {
	case '{', '[', '"':
		return readStream(buf)
	case 'n', 't', 'f':
		return readKeyword(buf)
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return readNumber(buf)
	default:
		return JsonValue{}, newErrUnexpected(buf, '{', '[', '"', 'n', 't', 'f',
			'-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9')
	}
}

// Parse takes an io.Reader and begins to parse from it, returning a JsonValue.
// This function also takes a size (in bytes) to use when creating the read
// buffer.
func Parse(r io.Reader, size int) (JsonValue, error) {
	data := make([]byte, size)
	buf := buffer{data, r, 0, nil, nil, uint32(size), 0, 0, 0, 0, 0, 0, 0, 0}
	_ = feedq(&buf) && feed(&buf)
	next(&buf)
	return readValue(&buf)
}

// escapemap is a mapping from escape sequences to escaped character values.
var escapemap = [...]byte{
	'"':  '"',
	'/':  '/',
	'\\': '\\',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

// Read implements the io.Reader interface for JsonValues (specifically, String
// values). Reading from this interface provides the value of the string.
func (data *JsonValue) Read(b []byte) (int, error) {
	if data.Type != String {
		return 0, newErrTypeMismatch(data.Type, String)
	} else if data.Status == Complete {
		return 0, io.EOF
	} else if data.Status != Working {
		return 0, ErrIncomplete
	}
	i := 0
	if data.buffer.escapes > 0 {
		i = streamEscape(data.buffer, b, 0)
	}
	for ; i < len(b); i++ {
		if data.buffer.err != nil && data.buffer.err != io.EOF {
			data.Status = Incomplete
			return i, data.buffer.err
		}
		c := data.buffer.curr
		switch {
		case c == '"':
			_ = feedq(data.buffer) && feed(data.buffer)
			next(data.buffer)
			data.Status = Complete
			data.buffer.depth--
			return i, io.EOF
		case c == '\\':
			_ = feedq(data.buffer) && feed(data.buffer)
			next(data.buffer)
			if data.buffer.err != nil && data.buffer.err != io.EOF {
				data.Status = Incomplete
				return i, data.buffer.err
			}
			k := data.buffer.curr
			switch k {
			case 'u':
				err := readUnicode(data.buffer)
				if err != nil {
					data.Status = Incomplete
					return i, err
				}
				i = streamEscape(data.buffer, b, i) - 1
			case '"', '/', '\\', 'b', 'f', 'n', 'r', 't':
				_ = feedq(data.buffer) && feed(data.buffer)
				next(data.buffer)
				b[i] = escapemap[k]
			default:
				data.Status = Incomplete
				return i, newErrUnexpected(data.buffer,
					'"', '/', '\\', 'u', 'b', 'f', 'n', 'r', 't')
			}
		case c <= '\x1F':
			data.Status = Incomplete
			err := newErrUnexpected(data.buffer)
			if data.buffer.err == io.EOF {
				err.CustomMsg = "premature EOF while attempting to read string"
			} else {
				err.CustomMsg = "control characters are not allowed in string values"
			}
			return i, err
		default:
			_ = feedq(data.buffer) && feed(data.buffer)
			next(data.buffer)
			b[i] = c
		}
	}
	return len(b), nil
}

// readUnicode reads a unicode escape (\uXXXX) from within a string value. Also
// supports reading UTF-16 surrogate pairs.
func readUnicode(buf *buffer) error {
	_ = feedq(buf) && feed(buf)
	next(buf)
	pt1, _, _, _, _, err := parseHex(buf)
	if err != nil {
		return err
	}
	var cp rune
	if pt1 >= 0xD800 && pt1 <= 0xDFFF {
		if buf.err != nil && buf.err != io.EOF {
			return buf.err
		}
		if buf.curr != '\\' {
			return newErrUnexpected(buf, '\\')
		}
		_ = feedq(buf) && feed(buf)
		next(buf)
		if buf.err != nil && buf.err != io.EOF {
			return buf.err
		}
		if buf.curr != 'u' {
			return newErrUnexpected(buf, 'u')
		}
		_ = feedq(buf) && feed(buf)
		next(buf)
		pt2, cx, cy, _, _, err := parseHex(buf)
		if err != nil {
			return err
		}
		if pt2 < 0xD000 || pt2 > 0xDFFF {
			return newErrUnexpectedChar(foffs(buf)-4, cx, 'D', 'd')
		} else if pt2 < 0xDC00 {
			return newErrUnexpectedChar(foffs(buf)-3, cy, 'C', 'D', 'E', 'F', 'c', 'd', 'e', 'f')
		}
		cp = 0x10000 + rune(pt1-0xD800)<<10 + rune(pt2-0xDC00)
	} else {
		cp = rune(pt1)
	}
	var b [4]byte
	ct := utf8.EncodeRune(b[:], cp)
	buf.escapes = byte(5 - ct)
	j := 0
	for i := buf.escapes; i < 5; i++ {
		switch i {
		case 1:
			buf.escape1 = b[j]
		case 2:
			buf.escape2 = b[j]
		case 3:
			buf.escape3 = b[j]
		case 4:
			buf.escape4 = b[j]
		}
		j++
	}
	return nil
}

// parseHex reads the next four digits from the buffer, and parses them as a
// 16 bit hexadecimal value.
func parseHex(buf *buffer) (uint16, byte, byte, byte, byte, error) {
	var n uint16
	var cs [4]byte
	for i := 0; i < 4; i++ {
		if buf.err != nil && buf.err != io.EOF {
			return 0, 0, 0, 0, 0, buf.err
		}
		switch {
		case buf.curr <= '9' && buf.curr >= '0':
			n = n<<4 + uint16(buf.curr-'0')
		case buf.curr <= 'F' && buf.curr >= 'A':
			n = n<<4 + uint16(buf.curr-'A'+10)
		case buf.curr <= 'f' && buf.curr >= 'a':
			n = n<<4 + uint16(buf.curr-'a'+10)
		default:
			return 0, 0, 0, 0, 0, newErrUnexpected(buf,
				'A', 'B', 'C', 'D', 'E', 'F', 'a', 'b', 'c', 'd', 'e', 'f',
				'0', '1', '2', '3', '4', '5', '6', '7', '8', '9')
		}
		cs[i] = buf.curr
		_ = feedq(buf) && feed(buf)
		next(buf)
	}
	return n, cs[0], cs[1], cs[2], cs[3], nil
}

// streamEscape streams unicode escapes out of the buffer.
func streamEscape(buf *buffer, b []byte, i int) int {
	for buf.escapes < 5 && i < len(b) {
		switch buf.escapes {
		case 1:
			b[i] = buf.escape1
		case 2:
			b[i] = buf.escape2
		case 3:
			b[i] = buf.escape3
		case 4:
			b[i] = buf.escape4
		}
		buf.escapes++
		i++
	}
	if buf.escapes >= 5 {
		buf.escapes = 0
	}
	return i
}

// ValueNum returns the value of a Number.
func (data *JsonValue) ValueNum() (float64, error) {
	if data.Type == Number {
		return data.numval, nil
	}
	return 0, newErrTypeMismatch(data.Type, Number)
}

// ValueBool returns the value of a Bool.
func (data *JsonValue) ValueBool() (bool, error) {
	if data.Type == Bool {
		return data.boolval, nil
	}
	return false, newErrTypeMismatch(data.Type, Bool)
}

// readNext reads the next key or value from an Object or Array, respectively.
// This is a shared function, because the logic is the same in both cases.
func readNext(data *JsonValue, open byte, close byte) error {
	c, err := skipSpace(data.buffer)
	if err != nil {
		data.Status = Incomplete
		return err
	}
	if c == close {
		_ = feedq(data.buffer) && feed(data.buffer)
		next(data.buffer)
		data.Status = Complete
		data.buffer.depth--
		return EndOfValue
	}
	var expect byte = ','
	if data.boolval == false {
		expect = open
	}
	if c != expect {
		data.Status = Incomplete
		return newErrUnexpected(data.buffer, expect, close)
	}
	_ = feedq(data.buffer) && feed(data.buffer)
	next(data.buffer)
	if data.boolval == false {
		c, err = skipSpace(data.buffer)
		if err != nil {
			data.Status = Incomplete
			return err
		}
		if c == close {
			_ = feedq(data.buffer) && feed(data.buffer)
			next(data.buffer)
			data.Status = Complete
			data.buffer.depth--
			return EndOfValue
		}
	}
	return nil
}

// NextKey reads the next key from an Object, and returns it as a JsonValue. If
// the next part of the Object to parse is a value, that value is discarded, and
// the following key is returned. If the end of the Object is found, an
// EndOfValue error is returned.
func (data *JsonValue) NextKey() (JsonValue, error) {
	if data.Type != Object {
		return JsonValue{}, newErrTypeMismatch(data.Type, Object)
	} else if data.Status == Complete {
		return JsonValue{}, EndOfValue
	} else if data.Status != Working {
		return JsonValue{}, ErrIncomplete
	} else if data.depth != data.buffer.depth {
		return JsonValue{}, ErrWorkingChild
	}
	if data.keynext == false {
		val, err := objectNextValue(data)
		if err != nil {
			data.Status = Incomplete
			return JsonValue{}, err
		}
		err = val.Close()
		if err != nil {
			data.Status = Incomplete
			return JsonValue{}, err
		}
	}
	err := readNext(data, '{', '}')
	if err != nil {
		return JsonValue{}, err
	}
	_, err = skipSpace(data.buffer)
	if err != nil {
		data.Status = Incomplete
		return JsonValue{}, err
	}
	if data.buffer.curr != '"' {
		data.Status = Incomplete
		return JsonValue{}, newErrUnexpected(data.buffer, '"')
	}
	val, _ := readStream(data.buffer)
	data.boolval = true
	data.keynext = false
	return val, nil
}

// arrayNextValue reads the next value from an Object.
func objectNextValue(data *JsonValue) (JsonValue, error) {
	if data.depth != data.buffer.depth {
		return JsonValue{}, ErrWorkingChild
	}
	if data.keynext == true {
		key, err := data.NextKey()
		if err != nil {
			data.Status = Incomplete
			return JsonValue{}, err
		}
		err = key.Close()
		if err != nil {
			data.Status = Incomplete
			return JsonValue{}, err
		}
	}
	c, err := skipSpace(data.buffer)
	if err != nil {
		data.Status = Incomplete
		return JsonValue{}, err
	}
	if c != ':' {
		data.Status = Incomplete
		return JsonValue{}, newErrUnexpected(data.buffer, ':')
	}
	_ = feedq(data.buffer) && feed(data.buffer)
	next(data.buffer)
	val, err1 := readValue(data.buffer)
	if err1 != nil {
		data.Status = Incomplete
		return JsonValue{}, err1
	}
	data.keynext = true
	return val, nil
}

// arrayNextValue reads the next value from an Array.
func arrayNextValue(data *JsonValue) (JsonValue, error) {
	if data.depth != data.buffer.depth {
		return JsonValue{}, ErrWorkingChild
	}
	err := readNext(data, '[', ']')
	if err != nil {
		return JsonValue{}, err
	}
	val, err1 := readValue(data.buffer)
	if err1 != nil {
		data.Status = Incomplete
		return JsonValue{}, err1
	}
	data.boolval = true
	return val, nil
}

// NextValue reads the next value from an Object or Array, and returns it as a
// JsonValue. If this is used on an Object, and the next part of the Object to
// parse is a key, that key is discarded and the corresponding value is
// returned. If the end of the Object or Array is found, an EndOfValue error is
// returned.
func (data *JsonValue) NextValue() (JsonValue, error) {
	if data.Type == Array && data.Status == Working {
		return arrayNextValue(data)
	} else if data.Type == Object && data.Status == Working {
		return objectNextValue(data)
	} else if data.Type != Array && data.Type != Object {
		return JsonValue{}, newErrTypeMismatch(data.Type, Array, Object)
	} else if data.Status == Complete {
		return JsonValue{}, EndOfValue
	}
	return JsonValue{}, ErrIncomplete
}

// Close implements the io.Closer interface for JsonValues. Closing a JsonValue
// discards the remainder of that value from the stream. This is a fast way to
// ignore unimportant parts of the input to reach useful information.
func (data *JsonValue) Close() error {
	if data.Status == Complete {
		return nil
	} else if data.Status != Working {
		return ErrIncomplete
	} else if data.depth != data.buffer.depth {
		return ErrWorkingChild
	}
	if data.boolval == false {
		if data.Type == Object && data.buffer.curr == '{' ||
			data.Type == Array && data.buffer.curr == '[' {
			_ = feedq(data.buffer) && feed(data.buffer)
			next(data.buffer)
		}
	}
	instr := data.Type == String
	depth := 0
	for {
		if data.buffer.err != nil {
			data.Status = Incomplete
			if data.buffer.err == io.EOF {
				err := newErrUnexpected(data.buffer)
				err.CustomMsg = "premature EOF while attempting to close value"
				return err
			}
			return data.buffer.err
		}
		i := data.buffer.offs - 1
	InStr:
		if instr {
			for i < data.buffer.erroffs {
				switch data.buffer.data[i] {
				case '\\':
					i++
				case '"':
					if data.Type == String {
						data.buffer.offs = i + 1
						_ = feedq(data.buffer) && feed(data.buffer)
						next(data.buffer)
						data.Status = Complete
						data.buffer.depth--
						return nil
					}
					instr = false
					i++
					goto InStr
				}
				i++
			}
		} else {
			for i < data.buffer.erroffs {
				switch data.buffer.data[i] {
				case '{', '[':
					depth++
				case '}', ']':
					depth--
					if depth < 0 {
						data.buffer.offs = i + 1
						_ = feedq(data.buffer) && feed(data.buffer)
						next(data.buffer)
						data.Status = Complete
						data.buffer.depth--
						return nil
					}
				case '"':
					instr = true
					i++
					goto InStr
				}
				i++
			}
		}
		data.buffer.offs = data.buffer.erroffs
		if feedq(data.buffer) {
			feed(data.buffer)
			data.buffer.offs = i - uint32(len(data.buffer.data))
		}
		next(data.buffer)
	}
}

// simpleSort sorts the inputs in-place. Usually this is a short list, and may
// already be sorted (or mostly sorted), so a simple insertion sort is a good
// choice here.
func simpleSort(vals []string) {
	for i := 1; i < len(vals); i++ {
		for j := i; j > 0 && vals[j] < vals[j-1]; j-- {
			vals[j], vals[j-1] = vals[j-1], vals[j]
		}
	}
}

// compareRead reads a String value and does a comparison against the given
// slice of strings. The strings must be sorted for this to work properly.
func compareRead(data *JsonValue, vals []string) (string, bool, error) {
	var buf [16]byte
	x, y, z := 0, 0, 0
	for {
		l, err := data.Read(buf[:])
		if err != nil && err != io.EOF {
			return "", false, err
		}
		y = z
		z += l
		for {
			if len(vals[x]) >= z && vals[x][y:z] == string(buf[:l]) {
				break
			} else if x+1 >= len(vals) || vals[x][:y] != vals[x+1][:y] {
				return "", false, data.Close()
			}
			x++
		}
		if err == io.EOF && z == len(vals[x]) {
			return vals[x], true, nil
		} else if err == io.EOF {
			return "", false, nil
		}
	}
}

// Compare is a helper function, designed to read a String value and compare it
// against one or more arguments. If the String value matches none of the
// arguments, false is returned. Otherwise, true is returned along with the
// matched string.
func (data *JsonValue) Compare(vals ...string) (string, bool, error) {
	if len(vals) <= 0 {
		return "", false, ErrNoParamsSpecified
	}
	simpleSort(vals)
	return compareRead(data, vals)
}

// FindKey is a helper function, designed to find a specific value in an Object.
// It reads the object until it finds the first key that matches one of the
// provided arguments. If no keys match, false is returned. Otherwise, true is
// returned along with the matched string and the value associated with the
// matched key. Note that this will discard everything in the stream prior to
// the matched key/value pair.
func (data *JsonValue) FindKey(keys ...string) (string, JsonValue, bool, error) {
	if len(keys) <= 0 {
		return "", JsonValue{}, false, ErrNoParamsSpecified
	}
	simpleSort(keys)
	for {
		key, err := data.NextKey()
		if err == EndOfValue {
			return "", JsonValue{}, false, nil
		} else if err != nil {
			return "", JsonValue{}, false, err
		}
		k, match, err1 := compareRead(&key, keys)
		if err1 != nil {
			return "", JsonValue{}, false, err1
		} else if match {
			val, err := data.NextValue()
			if err != nil {
				return "", JsonValue{}, false, err
			}
			return k, val, true, nil
		}
	}
}
