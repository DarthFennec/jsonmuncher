package jsonmuncher

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

// special errors
var EndOfValue error = errors.New("End of value reached")
var UnexpectedEOF error = errors.New("Unexpected end of file")

// the data type of a JSON value
type JsonType byte

const (
	Null JsonType = iota
	Bool
	Number
	String
	Array
	Object
)

// the status of a value reader
type JsonStatus byte

const (
	Incomplete JsonStatus = iota
	Working
	Complete
)

// a buffer for an io.Reader:
// - data: an array to hold the buffered data
// - stream: the io.Reader to read from
// - err: the error associated with the value in the lookahead, or nil
// - readerr: the last error from Read()
// - offs: the index in data of the next byte to read
// - erroffs: the offset of the last error from Read()
// - depth: the current depth of the read value (nested arrays, objects, strings)
// - escape[s123]: buffer for string unicode escapes (see readUnicode function)
// - curr: the lookahead value
type Buffer struct {
	data    []byte
	stream  io.Reader
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

// a JSON value:
// - buffer: a pointer to the buffer the value is reading from
// - numval: if numeric, the parsed float
// - depth: the depth of the value (nested arrays, objects, strings)
// - Type: the data type of the value
// - Status: the status of this value reader
// - boolval: if boolean, the parsed value; if object or array, whether the
//     first element has been read
// - keynext: if object, whether the next thing to read is a key
type JsonValue struct {
	buffer  *Buffer
	numval  float64
	depth   uint32
	Type    JsonType
	Status  JsonStatus
	boolval bool
	keynext bool
}

// feedq, feed, and next should all be the same function, but they've been
// separated for the sake of performance. feedq and next can be inlined, and so
// will run much more quickly than they otherwise would. These functions should
// be called like:
//   _ = feedq(buf) && feed(buf)
//   next(buf)

// check if the next byte is beyond the end of the buffer
// can be inlined; separated from `feed' to make inlining possible
func feedq(buf *Buffer) bool {
	return int(buf.offs) >= len(buf.data)
}

// assuming feedq is true, feed the buffer with the next chunk
// cannot be inlined, because of the call to Read()
func feed(buf *Buffer) bool {
	var erroffs, readoffs int
	var readerr error
	for erroffs < len(buf.data) && readerr == nil {
		readoffs, readerr = buf.stream.Read(buf.data[erroffs:])
		erroffs += readoffs
	}
	if readerr == io.EOF {
		buf.readerr = UnexpectedEOF
	} else {
		buf.readerr = readerr
	}
	buf.erroffs = uint32(erroffs)
	buf.err = nil
	buf.offs = 0
	return false
}

// consume the next byte from the input stream, storing it in the lookahead
// can be inlined; separated from the above to make inlining possible
func next(buf *Buffer) {
	buf.curr = buf.data[buf.offs]
	if buf.offs < buf.erroffs {
		buf.offs += 1
	} else {
		buf.err = buf.readerr
	}
}

// skip whitespace until the next significant character
func skipSpace(buf *Buffer) (byte, error) {
	if buf.err != nil {
		return 0, buf.err
	}
	c := buf.curr
	for {
		switch c {
		case ' ', '\t', '\r', '\n':
			_ = feedq(buf) && feed(buf)
			next(buf)
			if buf.err != nil {
				return 0, buf.err
			}
			c = buf.curr
		default:
			return c, nil
		}
	}
}

// read a boolean or null value from the stream
func readKeyword(buf *Buffer) (JsonValue, error) {
	var kw string
	var typ JsonType = Bool
	var val bool = false
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
	for i := 1; i < len(kw); i += 1 {
		_ = feedq(buf) && feed(buf)
		next(buf)
		if buf.err != nil {
			return JsonValue{}, buf.err
		} else if kw[i] != buf.curr {
			return JsonValue{}, errors.New(fmt.Sprintf("Expected 'null', 'true', or 'false', got %q%q", kw[:i], buf.curr))
		}
	}
	_ = feedq(buf) && feed(buf)
	next(buf)
	return JsonValue{buf, 0, buf.depth + 1, typ, Complete, val, false}, nil
}

// read a numeric value from the stream
func readNumber(buf *Buffer) (JsonValue, error) {
	var b [32]byte
	sl := b[:0]
	simple := true
L:
	for buf.err != UnexpectedEOF {
		if buf.err != nil {
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
		var val int64
		neg := false
		idx := 0
		if sl[0] == '-' {
			if len(sl) == 1 {
				return JsonValue{}, errors.New("Bad formatting in number: nonassociated sign")
			}
			neg = true
			idx = 1
		}
		for ; idx < len(sl); idx += 1 {
			if sl[idx] <= '9' && sl[idx] >= '0' {
				val = 10*val + int64(sl[idx]-'0')
			} else {
				return JsonValue{}, errors.New("Bad formatting in number: unexpected non-digit symbol")
			}
		}
		if neg {
			val = -val
		}
		return JsonValue{buf, float64(val), buf.depth + 1, Number, Complete, false, false}, nil
	} else {
		f, err := strconv.ParseFloat(string(sl), 64)
		if err != nil {
			return JsonValue{}, err
		}
		return JsonValue{buf, f, buf.depth + 1, Number, Complete, false, false}, nil
	}
}

// read a string, array, or object from the stream
func readStream(buf *Buffer) (JsonValue, error) {
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
	buf.depth += 1
	return JsonValue{buf, 0, buf.depth, typ, Working, false, typ == Object}, nil
}

// read any value from the stream
func readValue(buf *Buffer) (JsonValue, error) {
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
		return JsonValue{}, errors.New(fmt.Sprintf("Unexpected character %q", buf.curr))
	}
}

// parse a reader stream into a JSON value
func Parse(r io.Reader, size int) (JsonValue, error) {
	data := make([]byte, size)
	buf := Buffer{data, r, nil, nil, uint32(size), 0, 0, 0, 0, 0, 0, 0, 0}
	_ = feedq(&buf) && feed(&buf)
	next(&buf)
	if buf.err != nil {
		return JsonValue{}, buf.err
	}
	return readValue(&buf)
}

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

// Reader interface (string values only)
func (data *JsonValue) Read(b []byte) (int, error) {
	if data.Type != String || data.Status != Working {
		return 0, errors.New("Operation only permitted for in-progress string values")
	}
	i := 0
	if data.buffer.escapes > 0 {
		i = streamEscape(data.buffer, b, 0)
	}
	for ; i < len(b); i += 1 {
		if data.buffer.err != nil {
			data.Status = Incomplete
			return i, data.buffer.err
		}
		c := data.buffer.curr
		switch {
		case c == '"':
			_ = feedq(data.buffer) && feed(data.buffer)
			next(data.buffer)
			data.Status = Complete
			data.buffer.depth -= 1
			return i, io.EOF
		case c == '\\':
			_ = feedq(data.buffer) && feed(data.buffer)
			next(data.buffer)
			if data.buffer.err != nil {
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
				return i, errors.New("Escape in string literal must be followed by \", /, \\, b, f, n, r, t, or u")
			}
		case c <= '\x1F':
			data.Status = Incomplete
			return i, errors.New("Control characters are not allowed in string literals")
		default:
			_ = feedq(data.buffer) && feed(data.buffer)
			next(data.buffer)
			b[i] = c
		}
	}
	return len(b), nil
}

// read a unicode escape from a string value
func readUnicode(buf *Buffer) error {
	_ = feedq(buf) && feed(buf)
	next(buf)
	pt1, err := parseHex(buf)
	if err != nil {
		return err
	}
	var cp rune
	if pt1 >= 0xD800 && pt1 <= 0xDFFF {
		if buf.err != nil {
			return buf.err
		}
		if buf.curr != '\\' {
			return errors.New("Expected UTF-16 surrogate pair")
		}
		_ = feedq(buf) && feed(buf)
		next(buf)
		if buf.err != nil {
			return buf.err
		}
		if buf.curr != 'u' {
			return errors.New("Expected UTF-16 surrogate pair")
		}
		_ = feedq(buf) && feed(buf)
		next(buf)
		pt2, err := parseHex(buf)
		if err != nil {
			return err
		}
		if pt2 < 0xDC00 || pt2 > 0xDFFF {
			return errors.New("Expected UTF-16 surrogate pair")
		}
		cp = 0x10000 + rune(pt1-0xD800)<<10 + rune(pt2-0xDC00)
	} else {
		cp = rune(pt1)
	}
	var b [4]byte
	ct := utf8.EncodeRune(b[:], cp)
	buf.escapes = byte(5 - ct)
	j := 0
	for i := buf.escapes; i < 5; i += 1 {
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
		j += 1
	}
	return nil
}

// read the next four digits from the buffer and parse them as a 16 bit hex
func parseHex(buf *Buffer) (uint16, error) {
	var n uint16
	for i := 0; i < 4; i += 1 {
		if buf.err != nil {
			return 0, buf.err
		}
		switch {
		case buf.curr <= '9' && buf.curr >= '0':
			n = n<<4 + uint16(buf.curr-'0')
		case buf.curr <= 'F' && buf.curr >= 'A':
			n = n<<4 + uint16(buf.curr-'A'+10)
		case buf.curr <= 'f' && buf.curr >= 'a':
			n = n<<4 + uint16(buf.curr-'a'+10)
		default:
			return 0, errors.New(fmt.Sprintf("Unicode escapes must be four-digit hex values, not %q", buf.curr))
		}
		_ = feedq(buf) && feed(buf)
		next(buf)
	}
	return n, nil
}

// stream unicode escapes out of the buffer
func streamEscape(buf *Buffer, b []byte, i int) int {
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
		buf.escapes += 1
		i += 1
	}
	if buf.escapes >= 5 {
		buf.escapes = 0
	}
	return i
}

// get the value (numeric values only)
func (data *JsonValue) ValueNum() (float64, error) {
	if data.Type == Number && data.Status == Complete {
		return data.numval, nil
	} else {
		return 0, errors.New("Operation only permitted for successfully parsed numeric values")
	}
}

// get the value (boolean values only)
func (data *JsonValue) ValueBool() (bool, error) {
	if data.Type == Bool && data.Status == Complete {
		return data.boolval, nil
	} else {
		return false, errors.New("Operation only permitted for successfully parsed boolean values")
	}
}

// get the next key (objects only)
func (data *JsonValue) NextKey() (JsonValue, error) {
	if data.Type != Object || data.Status != Working {
		return JsonValue{}, errors.New("Operation only permitted for in-progress object values")
	}
	if data.depth != data.buffer.depth {
		return JsonValue{}, errors.New("Cannot read from ancestor of working value")
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
	c, err := skipSpace(data.buffer)
	if err != nil {
		data.Status = Incomplete
		return JsonValue{}, err
	}
	if c == '}' {
		_ = feedq(data.buffer) && feed(data.buffer)
		next(data.buffer)
		data.Status = Complete
		data.buffer.depth -= 1
		return JsonValue{}, EndOfValue
	}
	var expect byte = ','
	if data.boolval == false {
		expect = '{'
	}
	if c != expect {
		data.Status = Incomplete
		return JsonValue{}, errors.New(fmt.Sprintf("Expected %q, got %q", expect, c))
	}
	_ = feedq(data.buffer) && feed(data.buffer)
	next(data.buffer)
	if data.boolval == false {
		c, err = skipSpace(data.buffer)
		if err != nil {
			data.Status = Incomplete
			return JsonValue{}, err
		}
		if c == '}' {
			_ = feedq(data.buffer) && feed(data.buffer)
			next(data.buffer)
			data.Status = Complete
			data.buffer.depth -= 1
			return JsonValue{}, EndOfValue
		}
	}
	val, err1 := readValue(data.buffer)
	if err1 != nil {
		data.Status = Incomplete
		return JsonValue{}, err1
	}
	if val.Type != String {
		data.Status = Incomplete
		return JsonValue{}, errors.New("Object keys must be string values")
	}
	data.boolval = true
	data.keynext = false
	return val, nil
}

// get the next value (objects only)
func objectNextValue(data *JsonValue) (JsonValue, error) {
	if data.depth != data.buffer.depth {
		return JsonValue{}, errors.New("Cannot read from ancestor of working value")
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
		return JsonValue{}, errors.New(fmt.Sprintf("Expected ':', got %q", c))
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

// get the next value (arrays only)
func arrayNextValue(data *JsonValue) (JsonValue, error) {
	if data.depth != data.buffer.depth {
		return JsonValue{}, errors.New("Cannot read from ancestor of working value")
	}
	c, err := skipSpace(data.buffer)
	if err != nil {
		data.Status = Incomplete
		return JsonValue{}, err
	}
	if c == ']' {
		_ = feedq(data.buffer) && feed(data.buffer)
		next(data.buffer)
		data.Status = Complete
		data.buffer.depth -= 1
		return JsonValue{}, EndOfValue
	}
	var expect byte = ','
	if data.boolval == false {
		expect = '['
	}
	if c != expect {
		data.Status = Incomplete
		return JsonValue{}, errors.New(fmt.Sprintf("Expected %q, got %q", expect, c))
	}
	_ = feedq(data.buffer) && feed(data.buffer)
	next(data.buffer)
	if data.boolval == false {
		c, err = skipSpace(data.buffer)
		if err != nil {
			data.Status = Incomplete
			return JsonValue{}, err
		}
		if c == ']' {
			_ = feedq(data.buffer) && feed(data.buffer)
			next(data.buffer)
			data.Status = Complete
			data.buffer.depth -= 1
			return JsonValue{}, EndOfValue
		}
	}
	val, err1 := readValue(data.buffer)
	if err1 != nil {
		data.Status = Incomplete
		return JsonValue{}, err1
	}
	data.boolval = true
	return val, nil
}

// get the next value (arrays and objects only)
func (data *JsonValue) NextValue() (JsonValue, error) {
	if data.Type == Array && data.Status == Working {
		return arrayNextValue(data)
	} else if data.Type == Object && data.Status == Working {
		return objectNextValue(data)
	} else {
		return JsonValue{}, errors.New("Operation only permitted for in-progress object and array values")
	}
}

// Closer interface, discard the remainder of this value from the stream
func (data *JsonValue) Close() error {
	if data.Status != Working {
		return nil
	}
	if data.depth != data.buffer.depth {
		return errors.New("Cannot close ancestor of working value")
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
	for depth >= 0 {
		if data.buffer.err != nil {
			data.Status = Incomplete
			return data.buffer.err
		}
		c := data.buffer.curr
		if instr {
			switch c {
			case '\\':
				_ = feedq(data.buffer) && feed(data.buffer)
				next(data.buffer)
				if data.buffer.err != nil {
					data.Status = Incomplete
					return data.buffer.err
				}
			case '"':
				instr = false
				depth -= 1
			}
		} else {
			switch c {
			case '}', ']':
				depth -= 1
			case '"':
				instr = true
				fallthrough
			case '{', '[':
				depth += 1
			}
		}
		_ = feedq(data.buffer) && feed(data.buffer)
		next(data.buffer)
	}
	data.Status = Complete
	data.buffer.depth -= 1
	return nil
}
