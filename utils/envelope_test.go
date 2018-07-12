package utils_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

type testObject struct {
	Foo   string `json:"foo"`
	Other int    `json:"other"`
}

func (t *testObject) Type() string { return "second" }

func TestTypedEnvelope(t *testing.T) {
	// error if JSON is malformed
	e := &utils.TypedEnvelope{}
	err := json.Unmarshal([]byte(`{`), e)
	assert.EqualError(t, err, "unexpected end of JSON input")

	// error if we don't have a type field
	e = &utils.TypedEnvelope{}
	err = json.Unmarshal([]byte(`{"foo":"bar","other":1234}`), e)
	assert.EqualError(t, err, "field 'type' is required")

	e = &utils.TypedEnvelope{}
	err = json.Unmarshal([]byte(`{"type":"first","foo":"bar","other":1234}`), e)
	assert.NoError(t, err)
	assert.Equal(t, "first", e.Type)
	assert.Equal(t, `{"type":"first","foo":"bar","other":1234}`, string(e.Data))

	o := &testObject{Foo: "bar", Other: 6543}
	e, err = utils.EnvelopeFromTyped(o)
	assert.NoError(t, err)
	assert.Equal(t, "second", e.Type)
	assert.Equal(t, `{"foo":"bar","other":6543}`, string(e.Data))

	data, err := json.Marshal(e)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"second","foo":"bar","other":6543}`, string(data))

	// nil typed value marshals to nil envelope
	e, err = utils.EnvelopeFromTyped(nil)
	assert.NoError(t, err)
	assert.Equal(t, (*utils.TypedEnvelope)(nil), e)
}
