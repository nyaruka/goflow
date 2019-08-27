package utils_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestJSONMarshaling(t *testing.T) {
	j, err := utils.JSONMarshal(nil)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`null`), j)

	j, err = utils.JSONMarshal("Rwanda > Kigali")
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"Rwanda > Kigali"`), j)

	j, err = utils.JSONMarshal(map[string]string{"foo": "bar"})
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"foo":"bar"}`), j)

	j, err = utils.JSONMarshalPretty(map[string]string{"foo": "bar"})
	assert.NoError(t, err)
	assert.Equal(t, []byte("{\n    \"foo\": \"bar\"\n}"), j)
}

func TestUnmarshalArray(t *testing.T) {
	// test empty array
	msgs, err := utils.UnmarshalArray([]byte(`[]`))
	assert.NoError(t, err)
	assert.Equal(t, []json.RawMessage{}, msgs)
}

func TestUnmarshalAndValidateWithLimit(t *testing.T) {
	data := []byte(`{"foo": "Hello"}`)
	buffer := ioutil.NopCloser(bytes.NewReader(data))

	// try with sufficiently large limit
	s := &struct {
		Foo string `json:"foo"`
	}{}
	err := utils.UnmarshalAndValidateWithLimit(buffer, s, 1000)
	assert.NoError(t, err)
	assert.Equal(t, "Hello", s.Foo)

	// try with limit that's smaller than the input
	buffer = ioutil.NopCloser(bytes.NewReader(data))
	s = &struct {
		Foo string `json:"foo"`
	}{}
	err = utils.UnmarshalAndValidateWithLimit(buffer, s, 5)
	assert.EqualError(t, err, "unexpected end of JSON input")
}

func TestJSONDecodeGeneric(t *testing.T) {
	// parse a JSON object into a map
	data := []byte(`{"bool": true, "number": 123.34, "text": "hello", "object": {"foo": "bar"}, "array": [1, "x"]}`)
	vals, err := utils.JSONDecodeGeneric(data)
	assert.NoError(t, err)

	asMap := vals.(map[string]interface{})
	assert.Equal(t, true, asMap["bool"])
	assert.Equal(t, json.Number("123.34"), asMap["number"])
	assert.Equal(t, "hello", asMap["text"])
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, asMap["object"])
	assert.Equal(t, []interface{}{json.Number("1"), "x"}, asMap["array"])

	// parse a JSON array into a slice
	data = []byte(`[{"foo": 123}, {"foo": 456}]`)
	vals, err = utils.JSONDecodeGeneric(data)
	assert.NoError(t, err)

	asSlice := vals.([]interface{})
	assert.Equal(t, map[string]interface{}{"foo": json.Number("123")}, asSlice[0])
	assert.Equal(t, map[string]interface{}{"foo": json.Number("456")}, asSlice[1])
}

func TestReadTypeFromJSON(t *testing.T) {
	_, err := utils.ReadTypeFromJSON([]byte(`{}`))
	assert.EqualError(t, err, "field 'type' is required")

	_, err = utils.ReadTypeFromJSON([]byte(`{"type": ""}`))
	assert.EqualError(t, err, "field 'type' is required")

	typeName, err := utils.ReadTypeFromJSON([]byte(`{"thing": 2, "type": "foo"}`))
	assert.NoError(t, err)
	assert.Equal(t, "foo", typeName)
}

func TestExtractResponseJSON(t *testing.T) {
	// valid HTTP trace and response body
	trace := "HTTP/1.1 200 OK\r\n" +
		"Server: nginx\r\n" +
		"Content-Type: application/json\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"Date: Tue, 27 Aug 2019 16:26:18 GMT\r\n" +
		"\r\n" +
		"{\"fact\":\"Cats have nine lives thanks to a flexible spine and powerful leg and back muscles\",\"length\":83}"

	assert.Equal(t, json.RawMessage(`{"fact":"Cats have nine lives thanks to a flexible spine and powerful leg and back muscles","length":83}`), utils.ExtractResponseJSON([]byte(trace)))

	// not a valid trace with body
	assert.Nil(t, utils.ExtractResponseJSON([]byte("HTTP/1.1 200 OK\r\n{}")))

	// valid trace, but body not JSON
	trace = "HTTP/1.1 200 OK\r\n" +
		"Server: nginx\r\n" +
		"Date: Tue, 27 Aug 2019 16:26:18 GMT\r\n" +
		"\r\n" +
		"Cats have nine lives"

	assert.Equal(t, json.RawMessage(`"Cats have nine lives"`), utils.ExtractResponseJSON([]byte(trace)))
}
