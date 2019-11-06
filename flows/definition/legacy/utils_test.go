package legacy_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows/definition/legacy"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTranslations(t *testing.T) {
	// can unmarshall from a single string
	translations := make(legacy.Translations)
	json.Unmarshal([]byte(`"hello"`), &translations)
	assert.Equal(t, legacy.Translations{"base": "hello"}, translations)

	// or a map
	translations = make(legacy.Translations)
	json.Unmarshal([]byte(`{"eng": "hello", "fra": "bonjour"}`), &translations)
	assert.Equal(t, legacy.Translations{"eng": "hello", "fra": "bonjour"}, translations)

	// and back to JSON
	data, err := json.Marshal(translations)
	require.NoError(t, err)
	assert.Equal(t, []byte(`{"eng":"hello","fra":"bonjour"}`), data)

	translationSet := []legacy.Translations{
		{"eng": "Yes", "fra": "Oui"},
		{"eng": "No", "fra": "Non"},
		{"eng": "Maybe"},
		{"eng": "Never", "fra": "Jamas"},
	}
	assert.Equal(t, map[envs.Language][]string{
		"eng": {"Yes", "No", "Maybe", "Never"},
		"fra": {"Oui", "Non", "", "Jamas"},
	}, legacy.TransformTranslations(translationSet))
}

func TestStringOrNumber(t *testing.T) {
	// can unmarshall from a string
	var s legacy.StringOrNumber
	err := json.Unmarshal([]byte(`"123.45"`), &s)
	assert.NoError(t, err)
	assert.Equal(t, legacy.StringOrNumber("123.45"), s)

	// or a floating point (JSON number type)
	err = json.Unmarshal([]byte(`567.89`), &s)
	assert.NoError(t, err)
	assert.Equal(t, legacy.StringOrNumber("567.89"), s)

	err = json.Unmarshal([]byte(`-567.89`), &s)
	assert.NoError(t, err)
	assert.Equal(t, legacy.StringOrNumber("-567.89"), s)

	err = json.Unmarshal([]byte(`[]`), &s)
	assert.EqualError(t, err, "expected string or number, not [")
}

type testObject struct {
	Foo   string `json:"foo"`
	Other int    `json:"other"`
}

func (t *testObject) Type() string { return "second" }
func TestTypedEnvelope(t *testing.T) {
	// error if JSON is malformed
	e := &legacy.TypedEnvelope{}
	err := json.Unmarshal([]byte(`{`), e)
	assert.EqualError(t, err, "unexpected end of JSON input")

	// error if we don't have a type field
	e = &legacy.TypedEnvelope{}
	err = json.Unmarshal([]byte(`{"foo":"bar","other":1234}`), e)
	assert.EqualError(t, err, "field 'type' is required")

	e = &legacy.TypedEnvelope{}
	err = json.Unmarshal([]byte(`{"type":"first","foo":"bar","other":1234}`), e)
	assert.NoError(t, err)
	assert.Equal(t, "first", e.Type)
	assert.Equal(t, `{"type":"first","foo":"bar","other":1234}`, string(e.Data))
}

func TestURLJoin(t *testing.T) {
	assert.Equal(t, "http://myfiles.com/test.jpg", legacy.URLJoin("http://myfiles.com", "test.jpg"))
	assert.Equal(t, "http://myfiles.com/test.jpg", legacy.URLJoin("http://myfiles.com/", "test.jpg"))
	assert.Equal(t, "http://myfiles.com/test.jpg", legacy.URLJoin("http://myfiles.com", "/test.jpg"))
	assert.Equal(t, "http://myfiles.com/test.jpg", legacy.URLJoin("http://myfiles.com/", "/test.jpg"))
	assert.Equal(t, "http://myfiles.com/test.jpg", legacy.URLJoin("http://myfiles.com/", "http://myfiles.com/test.jpg"))
	assert.Equal(t, "https://myfiles.com/test.jpg", legacy.URLJoin("https://myfiles.com/", "https://myfiles.com/test.jpg"))
}
