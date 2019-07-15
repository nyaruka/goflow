package inspect

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Valid1 string            `json:"valid1" engine:"localized"`
	Valid2 []string          `json:"valid2" engine:"localized,evaluated"`
	Valid3 map[string]string `json:"valid3" engine:"evaluated"`
	Valid4 string
	Bad1   int `engine:"evaluated"` // an int field can't be evaluated
	Bad2   int `engine:"localized"` // or localized
}

func TestParseEngineTag(t *testing.T) {
	typ := reflect.TypeOf(testStruct{})

	assertTags := func(fieldIndex int, name string, localized bool, evaluated bool) {
		f := typ.Field(fieldIndex)

		assert.Equal(t, name, jsonNameTag(f))

		actualLocalized, actualEvaluated := parseEngineTag(f)
		assert.Equal(t, localized, actualLocalized)
		assert.Equal(t, evaluated, actualEvaluated)
	}

	assertTags(0, "valid1", true, false)
	assertTags(1, "valid2", true, true)
	assertTags(2, "valid3", false, true)
	assertTags(3, "", false, false)

	assert.Panics(t, func() { parseEngineTag(typ.Field(4)) })
	assert.Panics(t, func() { parseEngineTag(typ.Field(5)) })
}
