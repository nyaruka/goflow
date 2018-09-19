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

func TestJSONDecodeToMap(t *testing.T) {
	data := []byte(`{"bool": true, "number": 123.34, "text": "hello", "dict": {"foo": "bar"}, "list": [1, "x"]}`)
	asMap, err := utils.JSONDecodeToMap(data)
	assert.NoError(t, err)
	assert.Equal(t, true, asMap["bool"])
	assert.Equal(t, json.Number("123.34"), asMap["number"])
	assert.Equal(t, "hello", asMap["text"])
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, asMap["dict"])
	assert.Equal(t, []interface{}{json.Number("1"), "x"}, asMap["list"])
}
