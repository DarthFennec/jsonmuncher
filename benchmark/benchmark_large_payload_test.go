/*
   Each test should process 41kb json record (based on Discourse API)
   It should read 2 arrays, and for each item in array get few fields.
   Basically it means processing full JSON file.
*/
package benchmark

import (
	"encoding/json"
	"testing"

	"github.com/darthfennec/jsonmuncher"
	"github.com/json-iterator/go"
	"github.com/bcicen/jstream"
	"github.com/francoispqt/gojay"
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
func BenchmarkJsonMuncherLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := os.Open("./fixture_large.json")
		data, _ := jsonmuncher.Parse(largeFixture, 4096)
		_, users, _, _ := data.FindKey("users")
		for {
			user, eof := users.NextValue()
			if eof == jsonmuncher.EndOfValue {
				break
			}
			_, name, _, _ := user.FindKey("username")
			name.Close()
			user.Close()
		}
		_, tops, _, _ := data.FindKey("topics")
		_, topics, _, _ := tops.FindKey("topics")
		for {
			topic, eof := topics.NextValue()
			if eof == jsonmuncher.EndOfValue {
				break
			}
			_, id, _, _ := topic.FindKey("id")
			id.ValueNum()
			_, slug, _, _ := topic.FindKey("slug")
			slug.Close()
			topic.Close()
		}
	}
}

/*
   github.com/buger/jsonparser
*/
func BenchmarkJsonParserLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := ioutil.ReadFile("./fixture_large.json")
		jsonparser.ArrayEach(largeFixture, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			jsonparser.Get(value, "username")
			nothing()
		}, "users")

		jsonparser.ArrayEach(largeFixture, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			jsonparser.GetInt(value, "id")
			jsonparser.Get(value, "slug")
			nothing()
		}, "topics", "topics")
	}
}

/*
   encoding/json
*/
func BenchmarkEncodingJsonStructLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := ioutil.ReadFile("./fixture_large.json")
		var data LargePayload
		json.Unmarshal(largeFixture, &data)

		for _, u := range data.Users {
			nothing(u.Username)
		}

		for _, t := range data.Topics.Topics {
			nothing(t.Id, t.Slug)
		}
	}
}

func BenchmarkEncodingJsonInterfaceLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := ioutil.ReadFile("./fixture_large.json")
		var data interface{}
		json.Unmarshal(largeFixture, &data)
		m := data.(map[string]interface{})

		users := m["users"].([]interface{})
		for _, u := range users {
			nothing(u.(map[string]interface{})["username"].(string))
		}

		topics := m["topics"].(map[string]interface{})["topics"].([]interface{})
		for _, t := range topics {
			tI := t.(map[string]interface{})
			nothing(tI["id"].(float64), tI["slug"].(string))
		}
	}
}

func BenchmarkEncodingJsonStreamStructLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := os.Open("./fixture_large.json")
		var data LargePayload
		json.NewDecoder(largeFixture).Decode(&data)

		for _, u := range data.Users {
			nothing(u.Username)
		}

		for _, t := range data.Topics.Topics {
			nothing(t.Id, t.Slug)
		}
	}
}

func BenchmarkEncodingJsonStreamInterfaceLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := os.Open("./fixture_large.json")
		var data interface{}
		json.NewDecoder(largeFixture).Decode(&data)
		m := data.(map[string]interface{})

		users := m["users"].([]interface{})
		for _, u := range users {
			nothing(u.(map[string]interface{})["username"].(string))
		}

		topics := m["topics"].(map[string]interface{})["topics"].([]interface{})
		for _, t := range topics {
			tI := t.(map[string]interface{})
			nothing(tI["id"].(float64), tI["slug"].(string))
		}
	}
}

/*
   github.com/bcicen/jstream
*/
func BenchmarkJstreamLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := os.Open("./fixture_large.json")
		decoder := jstream.NewDecoder(largeFixture, 2)
		for c := range decoder.Stream() {
			switch c.Value.(type) {
			case map[string]interface{}:
				user := c.Value.(map[string]interface{})
				nothing(user["username"].(string))
			default:
				nothing()
			}
		}
		largeFixture, _ = os.Open("./fixture_large.json")
		decoder = jstream.NewDecoder(largeFixture, 3)
		for c := range decoder.Stream() {
			switch c.Value.(type) {
			case map[string]interface{}:
				topic := c.Value.(map[string]interface{})
				nothing(topic["id"].(float64), topic["slug"].(string))
			default:
				nothing()
			}
		}
	}
}

/*
   github.com/francoispqt/gojay
*/
func BenchmarkGojayLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := ioutil.ReadFile("./fixture_large.json")
		var data LargePayload
		gojay.UnmarshalJSONObject(largeFixture, &data)

		for _, u := range data.Users {
			nothing(u.Username)
		}

		for _, t := range data.Topics.Topics {
			nothing(t.Id, t.Slug)
		}
	}
}

/*
   github.com/json-iterator/go
*/
func BenchmarkJsonIteratorLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := ioutil.ReadFile("./fixture_large.json")
		var data LargePayload
		jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(largeFixture, &data)

		for _, u := range data.Users {
			nothing(u.Username)
		}

		for _, t := range data.Topics.Topics {
			nothing(t.Id, t.Slug)
		}
	}
}

/*
   github.com/Jeffail/gabs
*/
func BenchmarkGabsLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := ioutil.ReadFile("./fixture_large.json")
		json, _ := gabs.ParseJSON(largeFixture)
		users, _ := json.Path("users").Children()
		for _, user := range users {
			nothing(user.Path("username"))
		}
		topics, _ := json.Path("topics.topics").Children()
		for _, topic := range topics {
			nothing(topic.Path("id"))
			nothing(topic.Path("slug"))
		}
	}
}

/*
   github.com/bitly/go-simplejson
*/
func BenchmarkGoSimpleJsonLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := ioutil.ReadFile("./fixture_large.json")
		json, _ := simplejson.NewJson(largeFixture)
		users, _ := json.Get("users").Array()
		for _, user := range users {
			nothing(user.(map[string]interface{})["username"])
		}
		topics, _ := json.Get("topics").Get("topics").Array()
		for _, topic := range topics {
			nothing(topic.(map[string]interface{})["id"])
			nothing(topic.(map[string]interface{})["slug"])
		}
	}
}

/*
   github.com/pquerna/ffjson
*/
func BenchmarkFFJsonLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := ioutil.ReadFile("./fixture_large.json")
		var data LargePayload
		ffjson.Unmarshal(largeFixture, &data)

		for _, u := range data.Users {
			nothing(u.Username)
		}

		for _, t := range data.Topics.Topics {
			nothing(t.Id, t.Slug)
		}
	}
}

/*
   github.com/antonholmquist/jason
*/
func BenchmarkJasonLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := ioutil.ReadFile("./fixture_large.json")
		json, _ := jason.NewObjectFromBytes(largeFixture)
		users, _ := json.GetObjectArray("users")
		for _, user := range users {
			user.GetString("username")
		}
		topics, _ := json.GetObjectArray("topics.topics")
		for _, topic := range topics {
			topic.GetFloat64("id")
			topic.GetString("slug")
		}
		nothing()
	}
}

/*
   github.com/mreiferson/go-ujson
*/
func BenchmarkUjsonLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := ioutil.ReadFile("./fixture_large.json")
		json, _ := ujson.NewFromBytes(largeFixture)
		users := json.Get("users").Array()
		for _, user := range users {
			user.Get("url").String()
		}
		topics := json.Get("topics").Get("topics").Array()
		for _, topic := range topics {
			topic.Get("id").Float64()
			topic.Get("slug").String()
		}
		nothing()
	}
}

/*
   github.com/a8m/djson
*/
func BenchmarkDjsonLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := ioutil.ReadFile("./fixture_large.json")
		m, _ := djson.DecodeObject(largeFixture)
		users := m["users"].([]interface{})
		for _, u := range users {
			nothing(u.(map[string]interface{})["username"].(string))
		}

		topics := m["topics"].(map[string]interface{})["topics"].([]interface{})
		for _, t := range topics {
			tI := t.(map[string]interface{})
			nothing(tI["id"].(float64), tI["slug"].(string))
		}
	}
}

/*
   github.com/ugorji/go/codec
*/
func BenchmarkUgorjiLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := ioutil.ReadFile("./fixture_large.json")
		decoder := codec.NewDecoderBytes(largeFixture, new(codec.JsonHandle))
		data := new(LargePayload)
		json.Unmarshal(largeFixture, &data)
		data.CodecDecodeSelf(decoder)

		for _, u := range data.Users {
			nothing(u.Username)
		}

		for _, t := range data.Topics.Topics {
			nothing(t.Id, t.Slug)
		}
	}
}

/*
   github.com/mailru/easyjson
*/
func BenchmarkEasyJsonLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		largeFixture, _ := ioutil.ReadFile("./fixture_large.json")
		lexer := &jlexer.Lexer{Data: largeFixture}
		data := new(LargePayload)
		data.UnmarshalEasyJSON(lexer)

		for _, u := range data.Users {
			nothing(u.Username)
		}

		for _, t := range data.Topics.Topics {
			nothing(t.Id, t.Slug)
		}
	}
}
