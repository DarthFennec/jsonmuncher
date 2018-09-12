package jsonmuncher

import (
	"io"
	"testing"
	"strings"
)

// TODO add useful error messages instead of just numbers

func TestBasicParsing(t *testing.T) {
	json := "  [ [   true  ,null ], [ ] , false ]  "
	r := strings.NewReader(json)
	v1, e1 := Parse(r, 4096)
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
	v1, e1 := Parse(r, 4096)
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
	json := "[-20.50e+1, 400E-2 ,12345678, 75.5]"
	r := strings.NewReader(json)
	v1, e1 := Parse(r, 4096)
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
	if vn != 75.5 || en != nil {
		t.Fatal("9", vn, en)
	}
	_, e2 = v1.NextValue()
	if e2 != EndOfValue {
		t.Fatal("10", e2)
	}
}

func TestStringParsing(t *testing.T) {
	json := "[\"\", \" \\// \\\\ \\n \\t \\b \\r \\f \\\" \", \" (\\u256f\\u00b0\\u25a1\\u00b0\\uff09\\u256f\\ufe35 \\u253b\\u2501\\u253b \"]"
	var buf [48]byte
	r := strings.NewReader(json)
	v1, e1 := Parse(r, 4096)
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
	json := "\"\\ud83e\\uddf8\""
	var buf [8]byte
	r := strings.NewReader(json)
	v, e := Parse(r, 1)
	if v.Type != String || e != nil {
		t.Fatal("1", v.Type, e)
	}
	s, eof := v.Read(buf[:])
	if eof != io.EOF || string(buf[:s]) != "🧸" {
		t.Fatal("2", string(buf[:]), s, eof)
	}
}
