/*
   Each test should process a 333mb json record (based on a large Helm index)
   It should read two fields from each entry, nested in arrays within a map, as
   well as two meta fields from the root
*/
package benchmark

import (
	"encoding/json"
	"testing"

	"github.com/darthfennec/jsonmuncher"
	"github.com/json-iterator/go"
	"github.com/Jeffail/gabs"
	"github.com/a8m/djson"
	"github.com/antonholmquist/jason"
	"github.com/bitly/go-simplejson"
	"github.com/buger/jsonparser"
	jlexer "github.com/mailru/easyjson/jlexer"
	"github.com/mreiferson/go-ujson"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/ugorji/go/codec"
	"io/ioutil"
	"os"
)

/*
   github.com/darthfennec/jsonmuncher
*/
func BenchmarkJsonMuncherHuge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hugeFixture, _ := os.Open("./fixture_huge.json")
		data, _ := jsonmuncher.Parse(hugeFixture, 4096)
		_, apiver, _, _ := data.FindKey("apiVersion")
		apiver.Close()
		_, gen, _, _ := data.FindKey("generated")
		gen.Close()
		_, entries, _, _ := data.FindKey("entries")
		for {
			chartname, eof := entries.NextKey()
			if eof == jsonmuncher.EndOfValue {
				break
			}
			chartname.Close()
			chart, _ := entries.NextValue()
			for {
				entry, eof := chart.NextValue()
				if eof == jsonmuncher.EndOfValue {
					break
				}
				_, name, _, _ := entry.FindKey("name")
				name.Close()
				_, version, _, _ := entry.FindKey("version")
				version.Close()
				entry.Close()
			}
		}
	}
}

/*
   github.com/buger/jsonparser
*/
func BenchmarkJsonParserHuge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hugeFixture, _ := ioutil.ReadFile("./fixture_huge.json")
		jsonparser.GetString(hugeFixture, "apiVersion")
		jsonparser.GetString(hugeFixture, "generated")
		jsonparser.ObjectEach(hugeFixture, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			jsonparser.ArrayEach(value, func(value2 []byte, dataType jsonparser.ValueType, offset2 int, err error) {
				jsonparser.GetString(value2, "name")
				jsonparser.GetString(value2, "version")
				nothing()
			})
			return nil
		}, "entries")
	}
}

/*
   encoding/json
*/
func BenchmarkEncodingJsonStructHuge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hugeFixture, _ := ioutil.ReadFile("./fixture_huge.json")
		var data IndexFile
		json.Unmarshal(hugeFixture, &data)
		nothing(data.APIVersion, data.Generated)
		for chartname, entry := range data.Entries {
			nothing(chartname)
			for _, chart := range entry {
				nothing(chart.Name, chart.Version)
			}
		}
	}
}

func BenchmarkEncodingJsonInterfaceHuge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hugeFixture, _ := ioutil.ReadFile("./fixture_huge.json")
		var data interface{}
		json.Unmarshal(hugeFixture, &data)
		m := data.(map[string]interface{})
		nothing(m["apiVersion"].(string), m["generated"].(string))
		entries := m["entries"].(map[string]interface{})
		for chartname, entry := range entries {
			nothing(chartname)
			charts := entry.([]interface{})
			for _, c := range charts {
				chart := c.(map[string]interface{})
				nothing(chart["name"].(string), chart["version"].(string))
			}
		}
	}
}

/*
   github.com/json-iterator/go
*/
func BenchmarkJsonIteratorHuge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hugeFixture, _ := ioutil.ReadFile("./fixture_huge.json")
		var data IndexFile
		jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(hugeFixture, &data)
		nothing(data.APIVersion, data.Generated)
		for chartname, entry := range data.Entries {
			nothing(chartname)
			for _, chart := range entry {
				nothing(chart.Name, chart.Version)
			}
		}
	}
}

/*
   github.com/Jeffail/gabs
*/
func BenchmarkGabsHuge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hugeFixture, _ := ioutil.ReadFile("./fixture_huge.json")
		json, _ := gabs.ParseJSON(hugeFixture)
		nothing(json.Path("apiVersion"), json.Path("generated"))
		entries, _ := json.Path("entries").ChildrenMap()
		for chartname, entry := range entries {
			nothing(chartname)
			charts, _ := entry.Children()
			for _, chart := range charts {
				nothing(chart.Path("name"), chart.Path("version"))
			}
		}
	}
}

/*
   github.com/bitly/go-simplejson
*/
func BenchmarkGoSimpleJsonHuge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hugeFixture, _ := ioutil.ReadFile("./fixture_huge.json")
		json, _ := simplejson.NewJson(hugeFixture)
		nothing(json.Get("apiVersion"), json.Get("generated"))
		entries, _ := json.Get("entries").Map()
		for chartname, entry := range entries {
			nothing(chartname)
			charts := entry.([]interface{})
			for _, c := range charts {
				chart := c.(map[string]interface{})
				nothing(chart["name"].(string), chart["version"].(string))
			}
		}
	}
}

/*
   github.com/pquerna/ffjson
*/
func BenchmarkFFJsonHuge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hugeFixture, _ := ioutil.ReadFile("./fixture_huge.json")
		var data IndexFile
		ffjson.Unmarshal(hugeFixture, &data)
		nothing(data.APIVersion, data.Generated)
		for chartname, entry := range data.Entries {
			nothing(chartname)
			for _, chart := range entry {
				nothing(chart.Name, chart.Version)
			}
		}
	}
}

/*
   github.com/antonholmquist/jason
*/
func BenchmarkJasonHuge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hugeFixture, _ := ioutil.ReadFile("./fixture_huge.json")
		json, _ := jason.NewObjectFromBytes(hugeFixture)
		json.GetString("apiVersion")
		json.GetString("generated")
		entries, _ := json.GetObject("entries")
		for chartname, entry := range entries.Map() {
			nothing(chartname)
			charts, _ := entry.ObjectArray()
			for _, chart := range charts {
				chart.GetString("name")
				chart.GetString("version")
			}
		}
		nothing()
	}
}

/*
   github.com/mreiferson/go-ujson
*/
func BenchmarkUjsonHuge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hugeFixture, _ := ioutil.ReadFile("./fixture_huge.json")
		json, _ := ujson.NewFromBytes(hugeFixture)
		json.Get("apiVersion").String()
		json.Get("generated").String()
		entries := json.Get("entries").Map()
		for chartname, entry := range entries {
			nothing(chartname)
			charts := entry.(*[]interface{})
			for _, c := range *charts {
				chart := c.(map[string]interface{})
				nothing(chart["name"].(string), chart["version"].(string))
			}
		}
		nothing()
	}
}

/*
   github.com/a8m/djson
*/
func BenchmarkDjsonHuge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hugeFixture, _ := ioutil.ReadFile("./fixture_huge.json")
		m, _ := djson.DecodeObject(hugeFixture)
		nothing(m["apiVersion"].(string), m["generated"].(string))
		entries := m["entries"].(map[string]interface{})
		for chartname, entry := range entries {
			nothing(chartname)
			charts := entry.([]interface{})
			for _, c := range charts {
				chart := c.(map[string]interface{})
				nothing(chart["name"].(string), chart["version"].(string))
			}
		}
	}
}

/*
   github.com/ugorji/go/codec
*/
func BenchmarkUgorjiHuge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hugeFixture, _ := ioutil.ReadFile("./fixture_huge.json")
		decoder := codec.NewDecoderBytes(hugeFixture, new(codec.JsonHandle))
		data := new(IndexFile)
		json.Unmarshal(hugeFixture, &data)
		data.CodecDecodeSelf(decoder)
		nothing(data.APIVersion, data.Generated)
		for chartname, entry := range data.Entries {
			nothing(chartname)
			for _, chart := range entry {
				nothing(chart.Name, chart.Version)
			}
		}
	}
}

/*
   github.com/mailru/easyjson
*/
func BenchmarkEasyJsonHuge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hugeFixture, _ := ioutil.ReadFile("./fixture_huge.json")
		lexer := &jlexer.Lexer{Data: hugeFixture}
		data := new(IndexFile)
		data.UnmarshalEasyJSON(lexer)
		nothing(data.APIVersion, data.Generated)
		for chartname, entry := range data.Entries {
			nothing(chartname)
			for _, chart := range entry {
				nothing(chart.Name, chart.Version)
			}
		}
	}
}
