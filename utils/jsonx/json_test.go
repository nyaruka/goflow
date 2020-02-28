package jsonx_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/utils/jsonx"

	"github.com/stretchr/testify/assert"
)

func TestMarshaling(t *testing.T) {
	j, err := jsonx.Marshal(nil)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`null`), j)

	// check that HTML entities aren't encoded
	j, err = jsonx.Marshal("Rwanda > Kigali & Ecuador")
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"Rwanda > Kigali & Ecuador"`), j)

	j, err = jsonx.Marshal(map[string]string{"foo": "bar"})
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"foo":"bar"}`), j)

	j, err = jsonx.MarshalPretty(map[string]string{"foo": "bar"})
	assert.NoError(t, err)
	assert.Equal(t, []byte("{\n    \"foo\": \"bar\"\n}"), j)

	j, err = jsonx.MarshalMerged(map[string]string{"foo": "bar"}, map[string]string{"zed": "xyz"})
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"foo":"bar","zed":"xyz"}`), j)
}

func TestUnmarshalArray(t *testing.T) {
	// test empty array
	msgs, err := jsonx.UnmarshalArray([]byte(`[]`))
	assert.NoError(t, err)
	assert.Equal(t, []json.RawMessage{}, msgs)
}

func TestUnmarshalWithLimit(t *testing.T) {
	data := []byte(`{"foo": "Hello"}`)
	buffer := ioutil.NopCloser(bytes.NewReader(data))

	// try with sufficiently large limit
	s := &struct {
		Foo string `json:"foo"`
	}{}
	err := jsonx.UnmarshalWithLimit(buffer, s, 1000)
	assert.NoError(t, err)
	assert.Equal(t, "Hello", s.Foo)

	// try with limit that's smaller than the input
	buffer = ioutil.NopCloser(bytes.NewReader(data))
	s = &struct {
		Foo string `json:"foo"`
	}{}
	err = jsonx.UnmarshalWithLimit(buffer, s, 5)
	assert.EqualError(t, err, "unexpected end of JSON input")
}

func TestDecodeGeneric(t *testing.T) {
	// parse a JSON object into a map
	data := []byte(`{"bool": true, "number": 123.34, "text": "hello", "object": {"foo": "bar"}, "array": [1, "x"]}`)
	vals, err := jsonx.DecodeGeneric(data)
	assert.NoError(t, err)

	asMap := vals.(map[string]interface{})
	assert.Equal(t, true, asMap["bool"])
	assert.Equal(t, json.Number("123.34"), asMap["number"])
	assert.Equal(t, "hello", asMap["text"])
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, asMap["object"])
	assert.Equal(t, []interface{}{json.Number("1"), "x"}, asMap["array"])

	// parse a JSON array into a slice
	data = []byte(`[{"foo": 123}, {"foo": 456}]`)
	vals, err = jsonx.DecodeGeneric(data)
	assert.NoError(t, err)

	asSlice := vals.([]interface{})
	assert.Equal(t, map[string]interface{}{"foo": json.Number("123")}, asSlice[0])
	assert.Equal(t, map[string]interface{}{"foo": json.Number("456")}, asSlice[1])
}
