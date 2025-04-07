package flows_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/assert"
)

func TestLocals(t *testing.T) {
	l1 := flows.NewLocals()
	assert.True(t, l1.IsZero())

	l1.Set("foo", "bar")
	l1.Set("zed", "123")
	l1.Set("tmp", "xyz")

	assert.Equal(t, "bar", l1.Get("foo"))
	assert.Equal(t, "123", l1.Get("zed"))
	assert.Equal(t, "xyz", l1.Get("tmp"))
	assert.False(t, l1.IsZero())

	l1.Clear("tmp")

	marshaled, err := json.Marshal(l1)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"foo":"bar","zed":"123"}`, string(marshaled))

	var l2 flows.Locals
	err = json.Unmarshal(marshaled, &l2)
	assert.NoError(t, err)

	assert.Equal(t, "bar", l2.Get("foo"))
	assert.Equal(t, "123", l2.Get("zed"))
}

func TestLocalRefValidation(t *testing.T) {
	type testStruct struct {
		Valid1   string `json:"valid1"   validate:"local_ref"`
		Valid2   string `json:"valid2"   validate:"local_ref"`
		Invalid1 string `json:"invalid1" validate:"local_ref"`
		Invalid2 string `json:"invalid2" validate:"local_ref"`
	}

	obj := testStruct{
		Valid1:   "color123",
		Valid2:   "_llm_output",
		Invalid1: "1234567890123456789012345678901234567890123456789012345678901234567890", // too long
		Invalid2: "1foo",                                                                   // starts with a number
	}
	err := utils.Validate(obj)
	assert.EqualError(t, err, "field 'invalid1' is not a valid local variable reference, field 'invalid2' is not a valid local variable reference")
}
