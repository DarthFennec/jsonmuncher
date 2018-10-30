package jsonmuncher

import (
	"strings"
	"testing"
)

func BenchmarkIntParsing(b *testing.B) {
	json := "[0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9]"
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(json)
		array, err := Parse(r, 16)
		var elem JsonValue
		for err != EndOfValue {
			elem, err = array.NextValue()
			if err == nil {
				elem.ValueNum()
			}
		}
	}
}

func BenchmarkFloatParsing(b *testing.B) {
	json := "[0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9]"
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(json)
		array, err := Parse(r, 16)
		var elem JsonValue
		for err != EndOfValue {
			elem, err = array.NextValue()
			if err == nil {
				elem.ValueNum()
			}
		}
	}
}

func BenchmarkIntClosing(b *testing.B) {
	json := "[0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9]"
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(json)
		array, err := Parse(r, 16)
		var elem JsonValue
		for err != EndOfValue {
			elem, err = array.NextValue()
			elem.Close()
		}
	}
}

func BenchmarkFloatClosing(b *testing.B) {
	json := "[0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9,0.1,2.3,4.5,6.7,8.9]"
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(json)
		array, err := Parse(r, 16)
		var elem JsonValue
		for err != EndOfValue {
			elem, err = array.NextValue()
			elem.Close()
		}
	}
}
