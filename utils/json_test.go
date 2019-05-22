package utils_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestIsValidJSON(t *testing.T) {
	assert.True(t, utils.IsValidJSON([]byte(`true`)))
	assert.True(t, utils.IsValidJSON([]byte(`false`)))
	assert.True(t, utils.IsValidJSON([]byte(`null`)))
	assert.True(t, utils.IsValidJSON([]byte(`"abc"`)))
	assert.True(t, utils.IsValidJSON([]byte(`123.456`)))
	assert.True(t, utils.IsValidJSON([]byte(`{}`)))
	assert.True(t, utils.IsValidJSON([]byte(`{"foo":"bar"}`)))
	assert.True(t, utils.IsValidJSON([]byte(`[]`)))
	assert.True(t, utils.IsValidJSON([]byte(`[1, "x"]`)))

	assert.False(t, utils.IsValidJSON(nil))
	assert.False(t, utils.IsValidJSON([]byte(`abc`)))
	assert.False(t, utils.IsValidJSON([]byte(`{`)))
	assert.False(t, utils.IsValidJSON([]byte(`{}xx`)))
	assert.False(t, utils.IsValidJSON([]byte(`{foo:"bar"}`)))
	assert.False(t, utils.IsValidJSON([]byte(`{0:"bar"}`)))
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

func TestGenericJSON(t *testing.T) {
	j, err := utils.ReadGenericJSON([]byte(`{"foo": {"bar": [{"doh": 123}]}}`))
	require.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": []interface{}{
				map[string]interface{}{
					"doh": json.Number("123"),
				},
			},
		},
	}, j.AsObject())

	objs := make([]map[string]interface{}, 0)
	j.WalkObjects(func(obj map[string]interface{}) {
		objs = append(objs, obj)
	})

	assert.Equal(t, []map[string]interface{}{
		j.AsObject(),
		map[string]interface{}{
			"bar": []interface{}{
				map[string]interface{}{
					"doh": json.Number("123"),
				},
			},
		},
		map[string]interface{}{
			"doh": json.Number("123"),
		},
	}, objs)

	marshaled, err := json.Marshal(j)
	assert.Equal(t, `{"foo":{"bar":[{"doh":123}]}}`, string(marshaled))
}
