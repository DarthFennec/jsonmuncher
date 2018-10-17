package jsonmuncher

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"
)

func assert(t *testing.T, k bool, args ...interface{}) {
	if k {
		t.Fatal(args...)
	}
}

// TODO add useful error messages instead of just numbers

func TestBasicParsing(t *testing.T) {
	json := "  [ [   true  ,null ], [ ] , false ]  "
	r := strings.NewReader(json)
	v1, e1 := Parse(r, 16)
	assert(t, v1.Type != Array || e1 != nil,
		"1", v1.Type, e1)
	v2, e2 := v1.NextValue()
	assert(t, v2.Type != Array || e2 != nil,
		"2", v2.Type, e2)
	v3, e3 := v2.NextValue()
	assert(t, v3.Type != Bool || e3 != nil,
		"3", v3.Type, e3)
	vb, eb := v3.ValueBool()
	assert(t, vb != true || eb != nil,
		"4", vb, eb)
	v3, e3 = v2.NextValue()
	assert(t, v3.Type != Null || e3 != nil,
		"5", v3.Type, e3)
	_, e3 = v2.NextValue()
	assert(t, e3 != EndOfValue,
		"6", e3)
	v2, e2 = v1.NextValue()
	assert(t, v2.Type != Array || e2 != nil,
		"7", v2.Type, e2)
	_, e3 = v2.NextValue()
	assert(t, e3 != EndOfValue,
		"8", e3)
	v2, e2 = v1.NextValue()
	assert(t, v2.Type != Bool || e2 != nil,
		"9", v2.Type, e2)
	vb, eb = v2.ValueBool()
	assert(t, vb != false || eb != nil,
		"10", vb, eb)
	_, e2 = v1.NextValue()
	assert(t, e2 != EndOfValue,
		"11", e2)
}

func TestObjectParsing(t *testing.T) {
	json := "  { \"full\" : {   \"foo\":  \"bar\"  ,\"baz\"  :\"ban\" }, \"empty\" : { }}  "
	var buf [8]byte
	r := strings.NewReader(json)
	v1, e1 := Parse(r, 16)
	assert(t, v1.Type != Object || e1 != nil,
		"1", v1.Type, e1)
	vk, ek := v1.NextKey()
	assert(t, vk.Type != String || ek != nil,
		"2", vk.Type, ek)
	s, eof := vk.Read(buf[:])
	assert(t, eof != io.EOF || string(buf[:s]) != "full",
		"3", string(buf[:]), s, eof)
	v2, e2 := v1.NextValue()
	assert(t, v2.Type != Object || e2 != nil,
		"4", v2.Type, e2)
	vk, ek = v2.NextKey()
	assert(t, vk.Type != String || ek != nil,
		"5", vk.Type, ek)
	s, eof = vk.Read(buf[:])
	assert(t, eof != io.EOF || string(buf[:s]) != "foo",
		"6", string(buf[:]), s, eof)
	v3, e3 := v2.NextValue()
	assert(t, v3.Type != String || e3 != nil,
		"7", v3.Type, e3)
	s, eof = v3.Read(buf[:])
	assert(t, eof != io.EOF || string(buf[:s]) != "bar",
		"8", string(buf[:]), s, eof)
	vk, ek = v2.NextKey()
	assert(t, vk.Type != String || ek != nil,
		"9", vk.Type, ek)
	s, eof = vk.Read(buf[:])
	assert(t, eof != io.EOF || string(buf[:s]) != "baz",
		"10", string(buf[:]), s, eof)
	v3, e3 = v2.NextValue()
	assert(t, v3.Type != String || e3 != nil,
		"11", v3.Type, e3)
	s, eof = v3.Read(buf[:])
	assert(t, eof != io.EOF || string(buf[:s]) != "ban",
		"12", string(buf[:]), s, eof)
	_, e3 = v2.NextValue()
	assert(t, e3 != EndOfValue,
		"13", e3)
	vk, ek = v1.NextKey()
	assert(t, vk.Type != String || ek != nil,
		"14", vk.Type, ek)
	s, eof = vk.Read(buf[:])
	assert(t, eof != io.EOF || string(buf[:s]) != "empty",
		"15", string(buf[:]), s, eof)
	v2, e2 = v1.NextValue()
	assert(t, v2.Type != Object,
		"16", v2.Type, e2)
	assert(t, e2 != nil,
		"17", v2.Type, e2)
	_, e3 = v2.NextValue()
	assert(t, e3 != EndOfValue,
		"18", e3)
	_, e2 = v1.NextValue()
	assert(t, e2 != EndOfValue,
		"19", e2)
}

func TestNumericParsing(t *testing.T) {
	json := "[-20.50e+1, 400E-2 ,12345678, -321, 75.5]"
	r := strings.NewReader(json)
	v1, e1 := Parse(r, 16)
	assert(t, v1.Type != Array || e1 != nil,
		"1", v1.Type, e1)
	v2, e2 := v1.NextValue()
	assert(t, v2.Type != Number || e2 != nil,
		"2", v2.Type, e2)
	vn, en := v2.ValueNum()
	assert(t, vn != -205 || en != nil,
		"3", vn, en)
	v2, e2 = v1.NextValue()
	assert(t, v2.Type != Number || e2 != nil,
		"4", v2.Type, e2)
	vn, en = v2.ValueNum()
	assert(t, vn != 4 || en != nil,
		"5", vn, en)
	v2, e2 = v1.NextValue()
	assert(t, v2.Type != Number || e2 != nil,
		"6", v2.Type, e2)
	vn, en = v2.ValueNum()
	assert(t, vn != 12345678 || en != nil,
		"7", vn, en)
	v2, e2 = v1.NextValue()
	assert(t, v2.Type != Number || e2 != nil,
		"8", v2.Type, e2)
	vn, en = v2.ValueNum()
	assert(t, vn != -321 || en != nil,
		"9", vn, en)
	v2, e2 = v1.NextValue()
	assert(t, v2.Type != Number || e2 != nil,
		"10", v2.Type, e2)
	vn, en = v2.ValueNum()
	assert(t, vn != 75.5 || en != nil,
		"11", vn, en)
	_, e2 = v1.NextValue()
	assert(t, e2 != EndOfValue,
		"12", e2)
}

func TestStringParsing(t *testing.T) {
	json := "[\"\", \" \\// \\\\ \\n \\t \\b \\r \\f \\\" \", \" (\\u256f\\u00b0\\u25a1\\u00b0\\uff09\\u256f\\ufe35 \\u253b\\u2501\\u253b \"]"
	var buf [48]byte
	r := strings.NewReader(json)
	v1, e1 := Parse(r, 16)
	assert(t, v1.Type != Array || e1 != nil,
		"1", v1.Type, e1)
	v2, e2 := v1.NextValue()
	assert(t, v2.Type != String || e2 != nil,
		"2", v2.Type, e2)
	s, eof := v2.Read(buf[:])
	assert(t, eof != io.EOF || string(buf[:s]) != "",
		"3", string(buf[:]), s, eof)
	v2, e2 = v1.NextValue()
	assert(t, v2.Type != String || e2 != nil,
		"4", v2.Type, e2)
	s, eof = v2.Read(buf[:])
	assert(t, eof != io.EOF || string(buf[:s]) != " // \\ \n \t \b \r \f \" ",
		"5", string(buf[:]), s, eof)
	v2, e2 = v1.NextValue()
	assert(t, v2.Type != String || e2 != nil,
		"6", v2.Type, e2)
	s, eof = v2.Read(buf[:])
	assert(t, eof != io.EOF || string(buf[:s]) != " (╯°□°）╯︵ ┻━┻ ",
		"7", string(buf[:]), s, eof)
	_, e2 = v1.NextValue()
	assert(t, e2 != EndOfValue,
		"8", e2)
}

func TestEscapeBuffer(t *testing.T) {
	json := "\"[\\uD83E\\uDDF8]\""
	var buf [16]byte
	var buf1 [1]byte
	r := strings.NewReader(json)
	v, e := Parse(r, 1)
	assert(t, v.Type != String || e != nil,
		"1", v.Type, e)
	s := 0
	_, eof := v.Read(buf1[:])
	for eof == nil {
		buf[s] = buf1[0]
		s++
		_, eof = v.Read(buf1[:])
	}
	assert(t, eof != io.EOF || string(buf[:s]) != "[🧸]",
		"2", string(buf[:]), s, eof)
}

func TestCloser(t *testing.T) {
	json := "{\"read\":[[\"\\\"skip\\\"\"]],\"close\":[[\"close\"]],\"skip\":[[\"read\"]]}"
	var buf [8]byte
	r := strings.NewReader(json)
	v1, e1 := Parse(r, 16)
	assert(t, v1.Type != Object || e1 != nil,
		"1", v1.Type, e1)
	vk, ek := v1.NextKey()
	assert(t, vk.Type != String || ek != nil,
		"2", vk.Type, ek)
	s, eof := vk.Read(buf[:])
	assert(t, eof != io.EOF || string(buf[:s]) != "read",
		"3", string(buf[:]), s, eof)
	vk, ek = v1.NextKey()
	assert(t, vk.Type != String || ek != nil,
		"4", vk.Type, ek)
	eof = vk.Close()
	assert(t, vk.Status != Complete || eof != nil,
		"5", vk.Status, eof)
	v2, e2 := v1.NextValue()
	assert(t, v2.Type != Array || e2 != nil,
		"6", v2.Type, e2)
	eof = v2.Close()
	assert(t, v2.Status != Complete || eof != nil,
		"7", v2.Status, eof)
	v2, e2 = v1.NextValue()
	assert(t, v2.Type != Array || e2 != nil,
		"8", v2.Type, e2)
	v3, e3 := v2.NextValue()
	assert(t, v3.Type != Array || e3 != nil,
		"9", v3.Type, e3)
	v4, e4 := v3.NextValue()
	assert(t, v4.Type != String || e4 != nil,
		"10", v4.Type, e4)
	s, eof = v4.Read(buf[:])
	assert(t, eof != io.EOF || string(buf[:s]) != "read",
		"11", string(buf[:]), s, eof)
	_, e4 = v3.NextValue()
	assert(t, e4 != EndOfValue,
		"12", e4)
	_, e3 = v2.NextValue()
	assert(t, e3 != EndOfValue,
		"13", e3)
	_, e2 = v1.NextValue()
	assert(t, e2 != EndOfValue,
		"14", e2)
}

func TestCompareHelper(t *testing.T) {
	json := "[\"match1\",\"match2\",\"nomatch\"]"
	r := strings.NewReader(json)
	v1, e1 := Parse(r, 16)
	assert(t, v1.Type != Array || e1 != nil,
		"1", v1.Type, e1)
	v2, e2 := v1.NextValue()
	assert(t, v2.Type != String || e2 != nil,
		"2", v2.Type, e2)
	sk, mk, ek := v2.Compare("match1", "match2")
	assert(t, sk != "match1" || mk != true || ek != nil,
		"3", sk, mk, ek)
	v2, e2 = v1.NextValue()
	assert(t, v2.Type != String || e2 != nil,
		"4", v2.Type, e2)
	sk, mk, ek = v2.Compare()
	assert(t, sk != "" || mk != false || ek == nil || ek.Error() != "At least one argument should be provided",
		"5", sk, mk, ek)
	sk, mk, ek = v2.Compare("match1", "match2")
	assert(t, sk != "match2" || mk != true || ek != nil,
		"6", sk, mk, ek)
	v2, e2 = v1.NextValue()
	assert(t, v2.Type != String || e2 != nil,
		"7", v2.Type, e2)
	sk, mk, ek = v2.Compare("match1", "match2")
	assert(t, sk != "" || mk != false || ek != nil,
		"8", sk, mk, ek)
}

func TestFindKeyHelper(t *testing.T) {
	json := "{\"foo\":1,\"bar\":2,\"baz\":3,\"ban\":4}"
	r := strings.NewReader(json)
	v1, e1 := Parse(r, 16)
	assert(t, v1.Type != Object || e1 != nil,
		"1", v1.Type, e1)
	sk, vk, mk, ek := v1.FindKey("bar", "baz", "ban")
	assert(t, sk != "bar" || mk != true || ek != nil,
		"2", sk, mk, ek)
	vn, en := vk.ValueNum()
	assert(t, vn != 2 || en != nil,
		"3", vn, en)
	sk, vk, mk, ek = v1.FindKey()
	assert(t, sk != "" || mk != false || ek == nil || ek.Error() != "At least one argument should be provided",
		"4", sk, mk, ek)
	sk, vk, mk, ek = v1.FindKey("bar", "baz", "ban")
	assert(t, sk != "baz" || mk != true || ek != nil,
		"5", sk, mk, ek)
	vn, en = vk.ValueNum()
	assert(t, vn != 3 || en != nil,
		"6", vn, en)
	sk, vk, mk, ek = v1.FindKey("bank")
	assert(t, sk != "" || mk != false || ek != nil,
		"7", sk, mk, ek)
}

// failure states

var eofErrors = []string{
	"Unexpected EOF at file offset 0, expected one of '{', '[', '\"', 'n', 't', 'f', '-', '0'-'9'",
	"Unexpected EOF at file offset 1, expected one of '{', '[', '\"', 'n', 't', 'f', '-', '0'-'9'",
	"Unexpected EOF at file offset 2, expected one of '{', '[', '\"', 'n', 't', 'f', '-', '0'-'9'",
	"Unexpected EOF at file offset 3, expected 'u'",
	"Unexpected EOF at file offset 4, expected 'l'",
	"Unexpected EOF at file offset 5, expected 'l'",
	"Unexpected EOF at file offset 6, expected one of ',', ']'",
	"Unexpected EOF at file offset 7, expected one of '{', '[', '\"', 'n', 't', 'f', '-', '0'-'9'",
	"Unexpected EOF at file offset 8, expected '\"'",
	"Unexpected EOF at file offset 9: premature EOF while attempting to close value",
	"Unexpected EOF at file offset 10: premature EOF while attempting to close value",
	"Unexpected EOF at file offset 11, expected ':'",
	"Unexpected EOF at file offset 12, expected one of '{', '[', '\"', 'n', 't', 'f', '-', '0'-'9'",
	"Unexpected EOF at file offset 13: premature EOF while attempting to read string",
	"Unexpected EOF at file offset 14: premature EOF while attempting to read string",
	"Unexpected EOF at file offset 15, expected one of ',', '}'",
	"Unexpected EOF at file offset 16, expected '\"'",
	"Unexpected EOF at file offset 17: premature EOF while attempting to read string",
	"Unexpected EOF at file offset 18: premature EOF while attempting to read string",
	"Unexpected EOF at file offset 19, expected ':'",
	"Unexpected EOF at file offset 20, expected one of '{', '[', '\"', 'n', 't', 'f', '-', '0'-'9'",
	"Unexpected EOF at file offset 21: premature EOF while attempting to close value",
	"Unexpected EOF at file offset 22: premature EOF while attempting to close value",
	"Unexpected EOF at file offset 23, expected one of ',', '}'",
	"Unexpected EOF at file offset 24, expected '\"'",
	"Unexpected EOF at file offset 25: premature EOF while attempting to read string",
	"Unexpected EOF at file offset 26: premature EOF while attempting to read string",
	"Unexpected EOF at file offset 27, expected one of '\"', '/', '\\\\', 'u', 'b', 'f', 'n', 'r', 't'",
	"Unexpected EOF at file offset 28, expected one of 'A'-'F', 'a'-'f', '0'-'9'",
	"Unexpected EOF at file offset 29, expected one of 'A'-'F', 'a'-'f', '0'-'9'",
	"Unexpected EOF at file offset 30, expected one of 'A'-'F', 'a'-'f', '0'-'9'",
	"Unexpected EOF at file offset 31, expected one of 'A'-'F', 'a'-'f', '0'-'9'",
	"Unexpected EOF at file offset 32, expected '\\\\'",
	"Unexpected EOF at file offset 33, expected 'u'",
	"Unexpected EOF at file offset 34, expected one of 'A'-'F', 'a'-'f', '0'-'9'",
	"Unexpected EOF at file offset 35, expected one of 'A'-'F', 'a'-'f', '0'-'9'",
	"Unexpected EOF at file offset 36, expected one of 'A'-'F', 'a'-'f', '0'-'9'",
	"Unexpected EOF at file offset 37, expected one of 'A'-'F', 'a'-'f', '0'-'9'",
	"Unexpected EOF at file offset 38: premature EOF while attempting to read string",
	"Unexpected EOF at file offset 39: premature EOF while attempting to read string",
	"Unexpected EOF at file offset 40, expected ':'",
	"Unexpected EOF at file offset 41, expected one of '{', '[', '\"', 'n', 't', 'f', '-', '0'-'9'",
	"Unexpected EOF at file offset 42: premature EOF while attempting to close value",
	"Unexpected EOF at file offset 43: premature EOF while attempting to close value",
	"Unexpected EOF at file offset 44: premature EOF while attempting to close value",
	"Unexpected EOF at file offset 45: premature EOF while attempting to close value",
	"Unexpected EOF at file offset 46: premature EOF while attempting to close value",
	"Unexpected EOF at file offset 47: premature EOF while attempting to close value",
	"Unexpected EOF at file offset 48: premature EOF while attempting to close value",
}

func TestEOFErrors(t *testing.T) {
	json := "[ null,{\"a\":\"a\",\"b\":\"b\",\"x\\uD83E\\uDDF8x\":\"y\\ny\"}]"
	f := func(tst string, trunc int) {
		data := json[:trunc]
		chk := func(ex error) {
			if ex != nil && ex != io.EOF {
				if ex.Error() == eofErrors[trunc] {
					panic("<found eof>")
				} else {
					t.Fatal(tst, ex.Error())
				}
			}
		}
		defer func() {
			rec := recover()
			if rec != "<found eof>" {
				panic(rec)
			}
		}()
		var buf [8]byte
		r := strings.NewReader(data)
		v1, e1 := Parse(r, 16)
		chk(e1)
		v2, e2 := v1.NextValue()
		chk(e2)
		v2, e2 = v1.NextValue()
		chk(e2)
		vk, ek := v2.NextValue()
		chk(ek)
		_, eof := vk.Read(buf[:])
		chk(eof)
		vk, ek = v2.NextKey()
		chk(ek)
		_, eof = vk.Read(buf[:])
		chk(eof)
		vk, ek = v2.NextKey()
		chk(ek)
		_, eof = vk.Read(buf[:])
		chk(eof)
		v3, e3 := v2.NextValue()
		chk(e3)
		eof = v3.Close()
		chk(eof)
		eof = v2.Close()
		chk(eof)
		eof = v1.Close()
		chk(eof)
		t.Fatal("No errors caught")
	}
	for i := 0; i < len(json); i++ {
		f(strconv.Itoa(i), i)
	}
}

func TestFindKeyEOFErrors(t *testing.T) {
	r := strings.NewReader("{foo")
	v1, e1 := Parse(r, 16)
	assert(t, v1.Type != Object || e1 != nil,
		"1", v1.Type, e1)
	sk, _, mk, ek := v1.FindKey("foo")
	assert(t, ek == nil,
		"2", sk, mk, ek)
	r = strings.NewReader("{\"foo")
	v1, e1 = Parse(r, 16)
	assert(t, v1.Type != Object || e1 != nil,
		"3", v1.Type, e1)
	sk, _, mk, ek = v1.FindKey("foo")
	assert(t, ek == nil,
		"4", sk, mk, ek)
	r = strings.NewReader("{\"foo\":1ee1}")
	v1, e1 = Parse(r, 16)
	assert(t, v1.Type != Object || e1 != nil,
		"5", v1.Type, e1)
	sk, _, mk, ek = v1.FindKey("foo")
	assert(t, ek == nil,
		"6", sk, mk, ek)
}

func TestBadKeywordErrors(t *testing.T) {
	r := strings.NewReader("nule")
	_, e1 := Parse(r, 16)
	assert(t, e1 == nil || e1.Error() != "Unexpected 'e' at file offset 3, expected 'l'",
		"1", e1)
	r = strings.NewReader("mull")
	_, e1 = Parse(r, 16)
	assert(t, e1 == nil || e1.Error() != "Unexpected 'm' at file offset 0, expected one of '{', '[', '\"', 'n', 't', 'f', '-', '0'-'9'",
		"2", e1)
}

func TestBadNumberErrors(t *testing.T) {
	r := strings.NewReader("-")
	_, e1 := Parse(r, 16)
	assert(t, e1 == nil || e1.Error() != "Unexpected EOF at file offset 1, expected one of '0'-'9'",
		"1", e1)
	r = strings.NewReader("1-2")
	_, e1 = Parse(r, 16)
	assert(t, e1 == nil || e1.Error() != "Unexpected '-' at file offset 1, expected one of '0'-'9'",
		"2", e1)
	r = strings.NewReader("1.e")
	_, e1 = Parse(r, 16)
	assert(t, e1 == nil || e1.Error() != "strconv.ParseFloat: parsing \"1.e\": invalid syntax",
		"3", e1)
}

func TestBadStringErrors(t *testing.T) {
	r := strings.NewReader("\"\\w\"")
	var buf [8]byte
	v, _ := Parse(r, 16)
	_, e1 := v.Read(buf[:])
	assert(t, e1 == nil || e1.Error() != "Unexpected 'w' at file offset 2, expected one of '\"', '/', '\\\\', 'u', 'b', 'f', 'n', 'r', 't'",
		"1", e1)
	r = strings.NewReader("\"\t\"")
	v, _ = Parse(r, 16)
	_, e1 = v.Read(buf[:])
	assert(t, e1 == nil || e1.Error() != "Unexpected '\\t' at file offset 1: control characters are not allowed in string values",
		"2", e1)
	r = strings.NewReader("\"\\uD83Ex\"")
	v, _ = Parse(r, 16)
	_, e1 = v.Read(buf[:])
	assert(t, e1 == nil || e1.Error() != "Unexpected 'x' at file offset 7, expected '\\\\'",
		"3", e1)
	r = strings.NewReader("\"\\uD83E\\n\"")
	v, _ = Parse(r, 16)
	_, e1 = v.Read(buf[:])
	assert(t, e1 == nil || e1.Error() != "Unexpected 'n' at file offset 8, expected 'u'",
		"4", e1)
	r = strings.NewReader("\"\\uD83E\\u00B0\"")
	v, _ = Parse(r, 16)
	_, e1 = v.Read(buf[:])
	assert(t, e1 == nil || e1.Error() != "Unexpected '0' at file offset 9, expected one of 'D', 'd'",
		"5", e1)
	r = strings.NewReader("\"\\uD83E\\uD0B0\"")
	v, _ = Parse(r, 16)
	_, e1 = v.Read(buf[:])
	assert(t, e1 == nil || e1.Error() != "Unexpected '0' at file offset 10, expected one of 'C'-'F', 'c'-'f'",
		"6", e1)
	r = strings.NewReader("\"\\uD83Gx\"")
	v, _ = Parse(r, 16)
	_, e1 = v.Read(buf[:])
	assert(t, e1 == nil || e1.Error() != "Unexpected 'G' at file offset 6, expected one of 'A'-'F', 'a'-'f', '0'-'9'",
		"7", e1)
}

func TestBadObjectErrors(t *testing.T) {
	r := strings.NewReader("{\"foo\",null}")
	v, _ := Parse(r, 16)
	_, e1 := v.NextValue()
	assert(t, e1 == nil || e1.Error() != "Unexpected ',' at file offset 6, expected ':'",
		"1", e1)
	r = strings.NewReader("{\"foo\":null:\"bar\":null}")
	v, _ = Parse(r, 16)
	v.NextValue()
	_, e1 = v.NextValue()
	assert(t, e1 == nil || e1.Error() != "Unexpected ':' at file offset 11, expected one of ',', '}'",
		"2", e1)
	r = strings.NewReader("[true:false]")
	v, _ = Parse(r, 16)
	v.NextValue()
	_, e1 = v.NextValue()
	assert(t, e1 == nil || e1.Error() != "Unexpected ':' at file offset 5, expected one of ',', ']'",
		"3", e1)
	r = strings.NewReader("{true:false}")
	v, _ = Parse(r, 16)
	_, e1 = v.NextKey()
	assert(t, e1 == nil || e1.Error() != "Unexpected 't' at file offset 1, expected '\"'",
		"4", e1)
}

func TestAPIErrors(t *testing.T) {
	var buf [8]byte
	r := strings.NewReader("{\"foo\":[[null]]}")
	v, _ := Parse(r, 16)
	_, e1 := v.ValueNum()
	assert(t, e1 == nil || e1.Error() != "Method cannot be called on type Object, only on Number",
		"1", e1)
	_, e1 = v.ValueBool()
	assert(t, e1 == nil || e1.Error() != "Method cannot be called on type Object, only on Bool",
		"2", e1)
	_, e1 = v.Read(buf[:])
	assert(t, e1 == nil || e1.Error() != "Method cannot be called on type Object, only on String",
		"3", e1)
	k, _ := v.NextKey()
	_, e1 = k.NextKey()
	assert(t, e1 == nil || e1.Error() != "Method cannot be called on type String, only on Object",
		"4", e1)
	_, e1 = k.NextValue()
	assert(t, e1 == nil || e1.Error() != "Method cannot be called on type String, only on Array or Object",
		"5", e1)
	_, e1 = v.NextKey()
	assert(t, e1 == nil || e1.Error() != "Unable to consume when child element is partially read",
		"6", e1)
	_, e1 = v.NextValue()
	assert(t, e1 == nil || e1.Error() != "Unable to consume when child element is partially read",
		"7", e1)
	e1 = v.Close()
	assert(t, e1 == nil || e1.Error() != "Unable to consume when child element is partially read",
		"8", e1)
	k.Close()
	_, e1 = k.Read(buf[:])
	assert(t, e1 != io.EOF,
		"9", e1)
	a1, _ := v.NextValue()
	a2, _ := a1.NextValue()
	_, e1 = a1.NextValue()
	assert(t, e1 == nil || e1.Error() != "Unable to consume when child element is partially read",
		"10", e1)
	a2.Close()
	a1.Close()
	v.Close()
	_, e1 = v.NextKey()
	assert(t, e1 == nil || e1.Error() != "End of value reached",
		"11", e1)
	_, e1 = v.NextValue()
	assert(t, e1 == nil || e1.Error() != "End of value reached",
		"12", e1)
	e1 = v.Close()
	assert(t, e1 != nil,
		"13", e1)
}

func TestStreamErrors(t *testing.T) {
	var buf [8]byte
	streamerr := errors.New("Connection lost")
	f := func(json string, offs int) (JsonValue, error) {
		r := strings.NewReader(json)
		data := make([]byte, 16)
		buf := buffer{data, r, 0, nil, nil, uint32(16), 0, 0, 0, 0, 0, 0, 0, 0}
		_ = feedq(&buf) && feed(&buf)
		buf.readerr = streamerr
		buf.erroffs = uint32(offs)
		next(&buf)
		return readValue(&buf)
	}
	_, err := f(" ", 0)
	assert(t, err != streamerr,
		"1", err)
	_, err = f("  ", 1)
	assert(t, err != streamerr,
		"2", err)
	_, err = f("n ", 1)
	assert(t, err != streamerr,
		"3", err)
	_, err = f("- ", 1)
	assert(t, err != streamerr,
		"4", err)
	x, _ := f("{ ", 1)
	err = x.Close()
	assert(t, err != streamerr,
		"5", err)
	_, err = x.NextKey()
	assert(t, err.Error() != "Status incomplete denotes failed read",
		"6", err)
	_, err = x.NextValue()
	assert(t, err.Error() != "Status incomplete denotes failed read",
		"7", err)
	err = x.Close()
	assert(t, err.Error() != "Status incomplete denotes failed read",
		"8", err)
	x, _ = f("[ ", 1)
	_, err = x.NextValue()
	assert(t, err != streamerr,
		"9", err)
	x, _ = f("[null, ", 5)
	x.NextValue()
	_, err = x.NextValue()
	assert(t, err != streamerr,
		"10", err)
	x, _ = f("{\"\" ", 3)
	_, err = x.NextValue()
	assert(t, err != streamerr,
		"11", err)
	x, _ = f("{\"\":null ", 8)
	x.NextValue()
	_, err = x.NextKey()
	assert(t, err != streamerr,
		"12", err)
	x, _ = f("{\"\":null, ", 9)
	x.NextValue()
	_, err = x.NextKey()
	assert(t, err != streamerr,
		"13", err)
	x, _ = f("{ ", 1)
	_, err = x.NextKey()
	assert(t, err != streamerr,
		"14", err)
	x, _ = f("\" ", 1)
	_, err = x.Read(buf[:])
	assert(t, err != streamerr,
		"15", err)
	_, err = x.Read(buf[:])
	assert(t, err.Error() != "Status incomplete denotes failed read",
		"16", err)
	x, _ = f("\"\\ ", 2)
	_, err = x.Read(buf[:])
	assert(t, err != streamerr,
		"17", err)
	x, _ = f("\"\\uD83E ", 7)
	_, err = x.Read(buf[:])
	assert(t, err != streamerr,
		"18", err)
	x, _ = f("\"\\uD83E\\ ", 8)
	_, err = x.Read(buf[:])
	assert(t, err != streamerr,
		"19", err)
	x, _ = f("\"\\uD83E\\u ", 9)
	_, err = x.Read(buf[:])
	assert(t, err != streamerr,
		"20", err)
}
