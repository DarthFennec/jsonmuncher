/*
   Each test should process 190 byte http log like json record
   It should read multiple fields
*/
package benchmark

import (
	"encoding/json"
	"testing"

	"github.com/Jeffail/gabs"
	"github.com/a8m/djson"
	"github.com/antonholmquist/jason"
	"github.com/bcicen/jstream"
	"github.com/bitly/go-simplejson"
	"github.com/buger/jsonparser"
	"github.com/darthfennec/jsonmuncher"
	"github.com/francoispqt/gojay"
	"github.com/json-iterator/go"
	jlexer "github.com/mailru/easyjson/jlexer"
	"github.com/mreiferson/go-ujson"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/ugorji/go/codec"
	"io/ioutil"
	"os"
)

// Just for emulating field access, so it will not throw "evaluated but not used"
func nothing(_ ...interface{}) {}

/*
   github.com/darthfennec/jsonmuncher
*/
func BenchmarkJsonMuncherSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := os.Open("./fixture_small.json")
		data, _ := jsonmuncher.Parse(smallFixture, 256)
		_, st, _, _ := data.FindKey("st")
		st.ValueNum()
		_, uuid, _, _ := data.FindKey("uuid")
		uuid.Close()
		_, ua, _, _ := data.FindKey("ua")
		ua.Close()
		_, tz, _, _ := data.FindKey("tz")
		tz.ValueNum()
	}
}

/*
   github.com/buger/jsonparser
*/
func BenchmarkJsonParserSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := ioutil.ReadFile("./fixture_small.json")
		jsonparser.Get(smallFixture, "uuid")
		jsonparser.GetInt(smallFixture, "tz")
		jsonparser.Get(smallFixture, "ua")
		jsonparser.GetInt(smallFixture, "st")

		nothing()
	}
}

/*
   encoding/json
*/
func BenchmarkEncodingJsonStructSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := ioutil.ReadFile("./fixture_small.json")
		var data SmallPayload
		json.Unmarshal(smallFixture, &data)

		nothing(data.Uuid, data.Tz, data.Ua, data.St)
	}
}

func BenchmarkEncodingJsonInterfaceSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := ioutil.ReadFile("./fixture_small.json")
		var data interface{}
		json.Unmarshal(smallFixture, &data)
		m := data.(map[string]interface{})

		nothing(m["uuid"].(string), m["tz"].(float64), m["ua"].(string), m["st"].(float64))
	}
}

func BenchmarkEncodingJsonStreamStructSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := os.Open("./fixture_small.json")
		var data SmallPayload
		json.NewDecoder(smallFixture).Decode(&data)

		nothing(data.Uuid, data.Tz, data.Ua, data.St)
	}
}

func BenchmarkEncodingJsonStreamInterfaceSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := os.Open("./fixture_small.json")
		var data interface{}
		json.NewDecoder(smallFixture).Decode(&data)
		m := data.(map[string]interface{})

		nothing(m["uuid"].(string), m["tz"].(float64), m["ua"].(string), m["st"].(float64))
	}
}

/*
   github.com/bcicen/jstream
*/
func BenchmarkJstreamSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := os.Open("./fixture_small.json")
		decoder := jstream.NewDecoder(smallFixture, 0)
		for c := range decoder.Stream() {
			m := c.Value.(map[string]interface{})
			nothing(m["uuid"].(string), m["tz"].(float64), m["ua"].(string), m["st"].(float64))
		}
	}
}

/*
   github.com/francoispqt/gojay
*/
func BenchmarkGojaySmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := ioutil.ReadFile("./fixture_small.json")
		var data SmallPayload
		gojay.UnmarshalJSONObject(smallFixture, &data)

		nothing(data.Uuid, data.Tz, data.Ua, data.St)
	}
}

/*
   github.com/json-iterator/go
*/
func BenchmarkJsonIteratorSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := ioutil.ReadFile("./fixture_small.json")
		var data SmallPayload
		jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(smallFixture, &data)

		nothing(data.Uuid, data.Tz, data.Ua, data.St)
	}
}

/*
   github.com/Jeffail/gabs
*/
func BenchmarkGabsSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := ioutil.ReadFile("./fixture_small.json")
		json, _ := gabs.ParseJSON(smallFixture)

		nothing(
			json.Path("uuid").Data().(string),
			json.Path("tz").Data().(float64),
			json.Path("ua").Data().(string),
			json.Path("st").Data().(float64),
		)
	}
}

/*
   github.com/bitly/go-simplejson
*/
func BenchmarkGoSimpleJsonSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := ioutil.ReadFile("./fixture_small.json")
		json, _ := simplejson.NewJson(smallFixture)

		json.Get("uuid").String()
		json.Get("tz").Float64()
		json.Get("ua").String()
		json.Get("st").Float64()

		nothing()
	}
}

/*
   github.com/pquerna/ffjson
*/
func BenchmarkFFJsonSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := ioutil.ReadFile("./fixture_small.json")
		var data SmallPayload
		ffjson.Unmarshal(smallFixture, &data)

		nothing(data.Uuid, data.Tz, data.Ua, data.St)
	}
}

/*
   github.com/bitly/go-simplejson
*/
func BenchmarkJasonSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := ioutil.ReadFile("./fixture_small.json")
		json, _ := jason.NewObjectFromBytes(smallFixture)

		json.GetString("uuid")
		json.GetFloat64("tz")
		json.GetString("ua")
		json.GetFloat64("st")

		nothing()
	}
}

/*
   github.com/mreiferson/go-ujson
*/
func BenchmarkUjsonSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := ioutil.ReadFile("./fixture_small.json")
		json, _ := ujson.NewFromBytes(smallFixture)

		json.Get("uuid").String()
		json.Get("tz").Float64()
		json.Get("ua").String()
		json.Get("st").Float64()

		nothing()
	}
}

/*
   github.com/a8m/djson
*/
func BenchmarkDjsonSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := ioutil.ReadFile("./fixture_small.json")
		m, _ := djson.DecodeObject(smallFixture)
		nothing(m["uuid"].(string), m["tz"].(float64), m["ua"].(string), m["st"].(float64))
	}
}

/*
   github.com/ugorji/go/codec
*/
func BenchmarkUgorjiSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := ioutil.ReadFile("./fixture_small.json")
		decoder := codec.NewDecoderBytes(smallFixture, new(codec.JsonHandle))
		data := new(SmallPayload)
		data.CodecDecodeSelf(decoder)

		nothing(data.Uuid, data.Tz, data.Ua, data.St)
	}
}

/*
   github.com/mailru/easyjson
*/
func BenchmarkEasyJsonSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallFixture, _ := ioutil.ReadFile("./fixture_small.json")
		lexer := &jlexer.Lexer{Data: smallFixture}
		data := new(SmallPayload)
		data.UnmarshalEasyJSON(lexer)

		nothing(data.Uuid, data.Tz, data.Ua, data.St)
	}
}
