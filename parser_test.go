package jsonmuncher

import (
	"io"
	"strconv"
	"strings"
	"testing"
)

// TODO add useful error messages instead of just numbers

func TestBasicParsing(t *testing.T) {
	json := "  [ [   true  ,null ], [ ] , false ]  "
	r := strings.NewReader(json)
	v1, e1 := Parse(r, 16)
	if v1.Type != Array || e1 != nil {
		t.Fatal("1", v1.Type, e1)
	}
	v2, e2 := v1.NextValue()
	if v2.Type != Array || e2 != nil {
		t.Fatal("2", v2.Type, e2)
	}
	v3, e3 := v2.NextValue()
	if v3.Type != Bool || e3 != nil {
		t.Fatal("3", v3.Type, e3)
	}
	vb, eb := v3.ValueBool()
	if vb != true || eb != nil {
		t.Fatal("4", vb, eb)
	}
	v3, e3 = v2.NextValue()
	if v3.Type != Null || e3 != nil {
		t.Fatal("5", v3.Type, e3)
	}
	_, e3 = v2.NextValue()
	if e3 != EndOfValue {
		t.Fatal("6", e3)
	}
	v2, e2 = v1.NextValue()
	if v2.Type != Array || e2 != nil {
		t.Fatal("7", v2.Type, e2)
	}
	_, e3 = v2.NextValue()
	if e3 != EndOfValue {
		t.Fatal("8", e3)
	}
	v2, e2 = v1.NextValue()
	if v2.Type != Bool || e2 != nil {
		t.Fatal("9", v2.Type, e2)
	}
	vb, eb = v2.ValueBool()
	if vb != false || eb != nil {
		t.Fatal("10", vb, eb)
	}
	_, e2 = v1.NextValue()
	if e2 != EndOfValue {
		t.Fatal("11", e2)
	}
}

func TestObjectParsing(t *testing.T) {
	json := "  { \"full\" : {   \"foo\":  \"bar\"  ,\"baz\"  :\"ban\" }, \"empty\" : { } , \"one\"  : \"two\"}  "
	var buf [8]byte
	r := strings.NewReader(json)
	v1, e1 := Parse(r, 16)
	if v1.Type != Object || e1 != nil {
		t.Fatal("1", v1.Type, e1)
	}
	vk, ek := v1.NextKey()
	if vk.Type != String || ek != nil {
		t.Fatal("2", vk.Type, ek)
	}
	s, eof := vk.Read(buf[:])
	if eof != io.EOF || string(buf[:s]) != "full" {
		t.Fatal("3", string(buf[:]), s, eof)
	}
	v2, e2 := v1.NextValue()
	if v2.Type != Object || e2 != nil {
		t.Fatal("4", v2.Type, e2)
	}
	vk, ek = v2.NextKey()
	if vk.Type != String || ek != nil {
		t.Fatal("5", vk.Type, ek)
	}
	s, eof = vk.Read(buf[:])
	if eof != io.EOF || string(buf[:s]) != "foo" {
		t.Fatal("6", string(buf[:]), s, eof)
	}
	v3, e3 := v2.NextValue()
	if v3.Type != String || e3 != nil {
		t.Fatal("7", v3.Type, e3)
	}
	s, eof = v3.Read(buf[:])
	if eof != io.EOF || string(buf[:s]) != "bar" {
		t.Fatal("8", string(buf[:]), s, eof)
	}
	vk, ek = v2.NextKey()
	if vk.Type != String || ek != nil {
		t.Fatal("9", vk.Type, ek)
	}
	s, eof = vk.Read(buf[:])
	if eof != io.EOF || string(buf[:s]) != "baz" {
		t.Fatal("10", string(buf[:]), s, eof)
	}
	v3, e3 = v2.NextValue()
	if v3.Type != String || e3 != nil {
		t.Fatal("11", v3.Type, e3)
	}
	s, eof = v3.Read(buf[:])
	if eof != io.EOF || string(buf[:s]) != "ban" {
		t.Fatal("12", string(buf[:]), s, eof)
	}
	_, e3 = v2.NextValue()
	if e3 != EndOfValue {
		t.Fatal("13", e3)
	}
	vk, ek = v1.NextKey()
	if vk.Type != String || ek != nil {
		t.Fatal("14", vk.Type, ek)
	}
	s, eof = vk.Read(buf[:])
	if eof != io.EOF || string(buf[:s]) != "empty" {
		t.Fatal("15", string(buf[:]), s, eof)
	}
	v2, e2 = v1.NextValue()
	if v2.Type != Object || e2 != nil {
		t.Fatal("16", v2.Type, e2)
	}
	_, e3 = v2.NextValue()
	if e3 != EndOfValue {
		t.Fatal("17", e3)
	}
	vk, ek = v1.NextKey()
	if vk.Type != String || ek != nil {
		t.Fatal("18", vk.Type, ek)
	}
	s, eof = vk.Read(buf[:])
	if eof != io.EOF || string(buf[:s]) != "one" {
		t.Fatal("19", string(buf[:]), s, eof)
	}
	v2, e2 = v1.NextValue()
	if v2.Type != String || e2 != nil {
		t.Fatal("20", v2.Type, e2)
	}
	s, eof = v2.Read(buf[:])
	if eof != io.EOF || string(buf[:s]) != "two" {
		t.Fatal("21", string(buf[:]), s, eof)
	}
	_, e2 = v1.NextValue()
	if e2 != EndOfValue {
		t.Fatal("22", e2)
	}
}

func TestNumericParsing(t *testing.T) {
	json := "[-20.50e+1, 400E-2 ,12345678, -321, 75.5]"
	r := strings.NewReader(json)
	v1, e1 := Parse(r, 16)
	if v1.Type != Array || e1 != nil {
		t.Fatal("1", v1.Type, e1)
	}
	v2, e2 := v1.NextValue()
	if v2.Type != Number || e2 != nil {
		t.Fatal("2", v2.Type, e2)
	}
	vn, en := v2.ValueNum()
	if vn != -205 || en != nil {
		t.Fatal("3", vn, en)
	}
	v2, e2 = v1.NextValue()
	if v2.Type != Number || e2 != nil {
		t.Fatal("4", v2.Type, e2)
	}
	vn, en = v2.ValueNum()
	if vn != 4 || en != nil {
		t.Fatal("5", vn, en)
	}
	v2, e2 = v1.NextValue()
	if v2.Type != Number || e2 != nil {
		t.Fatal("6", v2.Type, e2)
	}
	vn, en = v2.ValueNum()
	if vn != 12345678 || en != nil {
		t.Fatal("7", vn, en)
	}
	v2, e2 = v1.NextValue()
	if v2.Type != Number || e2 != nil {
		t.Fatal("8", v2.Type, e2)
	}
	vn, en = v2.ValueNum()
	if vn != -321 || en != nil {
		t.Fatal("9", vn, en)
	}
	v2, e2 = v1.NextValue()
	if v2.Type != Number || e2 != nil {
		t.Fatal("10", v2.Type, e2)
	}
	vn, en = v2.ValueNum()
	if vn != 75.5 || en != nil {
		t.Fatal("11", vn, en)
	}
	_, e2 = v1.NextValue()
	if e2 != EndOfValue {
		t.Fatal("12", e2)
	}
}

func TestStringParsing(t *testing.T) {
	json := "[\"\", \" \\// \\\\ \\n \\t \\b \\r \\f \\\" \", \" (\\u256f\\u00b0\\u25a1\\u00b0\\uff09\\u256f\\ufe35 \\u253b\\u2501\\u253b \"]"
	var buf [48]byte
	r := strings.NewReader(json)
	v1, e1 := Parse(r, 16)
	if v1.Type != Array || e1 != nil {
		t.Fatal("1", v1.Type, e1)
	}
	v2, e2 := v1.NextValue()
	if v2.Type != String || e2 != nil {
		t.Fatal("2", v2.Type, e2)
	}
	s, eof := v2.Read(buf[:])
	if eof != io.EOF || string(buf[:s]) != "" {
		t.Fatal("3", string(buf[:]), s, eof)
	}
	v2, e2 = v1.NextValue()
	if v2.Type != String || e2 != nil {
		t.Fatal("4", v2.Type, e2)
	}
	s, eof = v2.Read(buf[:])
	if eof != io.EOF || string(buf[:s]) != " // \\ \n \t \b \r \f \" " {
		t.Fatal("5", string(buf[:]), s, eof)
	}
	v2, e2 = v1.NextValue()
	if v2.Type != String || e2 != nil {
		t.Fatal("6", v2.Type, e2)
	}
	s, eof = v2.Read(buf[:])
	if eof != io.EOF || string(buf[:s]) != " (╯°□°）╯︵ ┻━┻ " {
		t.Fatal("7", string(buf[:]), s, eof)
	}
	_, e2 = v1.NextValue()
	if e2 != EndOfValue {
		t.Fatal("8", e2)
	}
}

func TestEscapeBuffer(t *testing.T) {
	json := "\"[\\uD83E\\uDDF8]\""
	var buf [16]byte
	var buf1 [1]byte
	r := strings.NewReader(json)
	v, e := Parse(r, 1)
	if v.Type != String || e != nil {
		t.Fatal("1", v.Type, e)
	}
	s := 0
	_, eof := v.Read(buf1[:])
	for eof == nil {
		buf[s] = buf1[0]
		s += 1
		_, eof = v.Read(buf1[:])
	}
	if eof != io.EOF || string(buf[:s]) != "[🧸]" {
		t.Fatal("2", string(buf[:]), s, eof)
	}
}

func TestCloser(t *testing.T) {
	json := "{\"read\":[[\"\\\"skip\\\"\"]],\"close\":[[\"close\"]],\"skip\":[[\"read\"]]}"
	var buf [8]byte
	r := strings.NewReader(json)
	v1, e1 := Parse(r, 16)
	if v1.Type != Object || e1 != nil {
		t.Fatal("1", v1.Type, e1)
	}
	vk, ek := v1.NextKey()
	if vk.Type != String || ek != nil {
		t.Fatal("2", vk.Type, ek)
	}
	s, eof := vk.Read(buf[:])
	if eof != io.EOF || string(buf[:s]) != "read" {
		t.Fatal("3", string(buf[:]), s, eof)
	}
	vk, ek = v1.NextKey()
	if vk.Type != String || ek != nil {
		t.Fatal("4", vk.Type, ek)
	}
	eof = vk.Close()
	if vk.Status != Complete || eof != nil {
		t.Fatal("5", vk.Status, eof)
	}
	v2, e2 := v1.NextValue()
	if v2.Type != Array || e2 != nil {
		t.Fatal("6", v2.Type, e2)
	}
	eof = v2.Close()
	if v2.Status != Complete || eof != nil {
		t.Fatal("7", v2.Status, eof)
	}
	v2, e2 = v1.NextValue()
	if v2.Type != Array || e2 != nil {
		t.Fatal("8", v2.Type, e2)
	}
	v3, e3 := v2.NextValue()
	if v3.Type != Array || e3 != nil {
		t.Fatal("9", v3.Type, e3)
	}
	v4, e4 := v3.NextValue()
	if v4.Type != String || e4 != nil {
		t.Fatal("10", v4.Type, e4)
	}
	s, eof = v4.Read(buf[:])
	if eof != io.EOF || string(buf[:s]) != "read" {
		t.Fatal("11", string(buf[:]), s, eof)
	}
	_, e4 = v3.NextValue()
	if e4 != EndOfValue {
		t.Fatal("12", e4)
	}
	_, e3 = v2.NextValue()
	if e3 != EndOfValue {
		t.Fatal("13", e3)
	}
	_, e2 = v1.NextValue()
	if e2 != EndOfValue {
		t.Fatal("14", e2)
	}
}

// failure states

func TestEOFErrors(t *testing.T) {
	json := "[ null,{\"a\":\"a\",\"b\":\"b\",\"x\\uD83E\\uDDF8x\":\"y\\ny\"}]"
	f := func(tst string, trunc int) {
		data := json[:trunc]
		chk := func(ex error) bool {
			res := false
			if ex != nil && ex != io.EOF {
				if ex == UnexpectedEOF {
					res = true
				} else {
					t.Fatal(tst, ex.Error())
				}
			}
			return res
		}
		var buf [8]byte
		r := strings.NewReader(data)
		v1, e1 := Parse(r, 16)
		if chk(e1) {
			return
		}
		v2, e2 := v1.NextValue()
		if chk(e2) {
			return
		}
		v2, e2 = v1.NextValue()
		if chk(e2) {
			return
		}
		vk, ek := v2.NextValue()
		if chk(ek) {
			return
		}
		_, eof := vk.Read(buf[:])
		if chk(eof) {
			return
		}
		vk, ek = v2.NextKey()
		if chk(ek) {
			return
		}
		_, eof = vk.Read(buf[:])
		if chk(eof) {
			return
		}
		vk, ek = v2.NextKey()
		if chk(ek) {
			return
		}
		_, eof = vk.Read(buf[:])
		if chk(eof) {
			return
		}
		v3, e3 := v2.NextValue()
		if chk(e3) {
			return
		}
		eof = v3.Close()
		if chk(eof) {
			return
		}
		eof = v2.Close()
		if chk(eof) {
			return
		}
		eof = v1.Close()
		if chk(eof) {
			return
		}
		t.Fatal("No errors caught")
	}
	for i := 0; i < len(json); i += 1 {
		f(strconv.Itoa(i), i)
	}
}

func TestBadKeywordErrors(t *testing.T) {
	r := strings.NewReader("nule")
	_, e1 := Parse(r, 16)
	if e1 == nil || e1.Error() != "Expected 'null', 'true', or 'false', got \"nul\"'e'" {
		t.Fatal("1", e1)
	}
	r = strings.NewReader("mull")
	_, e1 = Parse(r, 16)
	if e1 == nil || e1.Error() != "Unexpected character 'm'" {
		t.Fatal("2", e1)
	}
}

func TestBadNumberErrors(t *testing.T) {
	r := strings.NewReader("-")
	_, e1 := Parse(r, 16)
	if e1 == nil || e1.Error() != "Bad formatting in number: nonassociated sign" {
		t.Fatal("1", e1)
	}
	r = strings.NewReader("1-2")
	_, e1 = Parse(r, 16)
	if e1 == nil || e1.Error() != "Bad formatting in number: unexpected non-digit symbol" {
		t.Fatal("2", e1)
	}
	r = strings.NewReader("1.e")
	_, e1 = Parse(r, 16)
	if e1 == nil || e1.Error() != "strconv.ParseFloat: parsing \"1.e\": invalid syntax" {
		t.Fatal("3", e1)
	}
}

func TestBadStringErrors(t *testing.T) {
	r := strings.NewReader("\"\\w\"")
	var buf [8]byte
	v, _ := Parse(r, 16)
	_, e1 := v.Read(buf[:])
	if e1 == nil || e1.Error() != "Escape in string literal must be followed by \", /, \\, b, f, n, r, t, or u" {
		t.Fatal("1", e1)
	}
	r = strings.NewReader("\"\t\"")
	v, _ = Parse(r, 16)
	_, e1 = v.Read(buf[:])
	if e1 == nil || e1.Error() != "Control characters are not allowed in string literals" {
		t.Fatal("2", e1)
	}
	r = strings.NewReader("\"\\uD83Ex\"")
	v, _ = Parse(r, 16)
	_, e1 = v.Read(buf[:])
	if e1 == nil || e1.Error() != "Expected UTF-16 surrogate pair" {
		t.Fatal("3", e1)
	}
	r = strings.NewReader("\"\\uD83E\\n\"")
	v, _ = Parse(r, 16)
	_, e1 = v.Read(buf[:])
	if e1 == nil || e1.Error() != "Expected UTF-16 surrogate pair" {
		t.Fatal("3", e1)
	}
	r = strings.NewReader("\"\\uD83E\\u00B0\"")
	v, _ = Parse(r, 16)
	_, e1 = v.Read(buf[:])
	if e1 == nil || e1.Error() != "Expected UTF-16 surrogate pair" {
		t.Fatal("3", e1)
	}
	r = strings.NewReader("\"\\uD83Gx\"")
	v, _ = Parse(r, 16)
	_, e1 = v.Read(buf[:])
	if e1 == nil || e1.Error() != "Unicode escapes must be four-digit hex values, not 'G'" {
		t.Fatal("4", e1)
	}
}

func TestBadObjectErrors(t *testing.T) {
	r := strings.NewReader("{\"foo\",null}")
	v, _ := Parse(r, 16)
	_, e1 := v.NextValue()
	if e1 == nil || e1.Error() != "Expected ':', got ','" {
		t.Fatal("1", e1)
	}
	r = strings.NewReader("{\"foo\":null:\"bar\":null}")
	v, _ = Parse(r, 16)
	v.NextValue()
	_, e1 = v.NextValue()
	if e1 == nil || e1.Error() != "Expected ',', got ':'" {
		t.Fatal("1", e1)
	}
	r = strings.NewReader("[true:false]")
	v, _ = Parse(r, 16)
	v.NextValue()
	_, e1 = v.NextValue()
	if e1 == nil || e1.Error() != "Expected ',', got ':'" {
		t.Fatal("2", e1)
	}
	r = strings.NewReader("{true:false}")
	v, _ = Parse(r, 16)
	_, e1 = v.NextKey()
	if e1 == nil || e1.Error() != "Object keys must be string values" {
		t.Fatal("3", e1)
	}
}

func TestAPIErrors(t *testing.T) {
	var buf [8]byte
	r := strings.NewReader("{\"foo\":[[null]]}")
	v, _ := Parse(r, 16)
	_, e1 := v.ValueNum()
	if e1 == nil || e1.Error() != "Operation only permitted for successfully parsed numeric values" {
		t.Fatal("1", e1)
	}
	_, e1 = v.ValueBool()
	if e1 == nil || e1.Error() != "Operation only permitted for successfully parsed boolean values" {
		t.Fatal("2", e1)
	}
	_, e1 = v.Read(buf[:])
	if e1 == nil || e1.Error() != "Operation only permitted for in-progress string values" {
		t.Fatal("3", e1)
	}
	k, _ := v.NextKey()
	_, e1 = k.NextKey()
	if e1 == nil || e1.Error() != "Operation only permitted for in-progress object values" {
		t.Fatal("4", e1)
	}
	_, e1 = k.NextValue()
	if e1 == nil || e1.Error() != "Operation only permitted for in-progress object and array values" {
		t.Fatal("5", e1)
	}
	_, e1 = v.NextKey()
	if e1 == nil || e1.Error() != "Cannot read from ancestor of working value" {
		t.Fatal("6", e1)
	}
	_, e1 = v.NextValue()
	if e1 == nil || e1.Error() != "Cannot read from ancestor of working value" {
		t.Fatal("7", e1)
	}
	e1 = v.Close()
	if e1 == nil || e1.Error() != "Cannot close ancestor of working value" {
		t.Fatal("8", e1)
	}
	k.Close()
	_, e1 = k.Read(buf[:])
	if e1 == nil || e1.Error() != "Operation only permitted for in-progress string values" {
		t.Fatal("9", e1)
	}
	a1, _ := v.NextValue()
	a2, _ := a1.NextValue()
	_, e1 = a1.NextValue()
	if e1 == nil || e1.Error() != "Cannot read from ancestor of working value" {
		t.Fatal("10", e1)
	}
	a2.Close()
	a1.Close()
	v.Close()
	_, e1 = v.NextKey()
	if e1 == nil || e1.Error() != "Operation only permitted for in-progress object values" {
		t.Fatal("11", e1)
	}
	_, e1 = v.NextValue()
	if e1 == nil || e1.Error() != "Operation only permitted for in-progress object and array values" {
		t.Fatal("12", e1)
	}
	e1 = v.Close()
	if e1 != nil {
		t.Fatal("13", e1)
	}
}
