package inspect

import (
	"reflect"
	"testing"

	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/assert"
)

type badTagsStruct1 struct {
	Valid1 string            `json:"valid1" engine:"localized"`
	Valid2 []string          `json:"valid2" engine:"localized,evaluated"`
	Valid3 map[string]string `json:"valid3" engine:"evaluated"`
	Valid4 string
	Bad1   int `engine:"evaluated"` // an int field can't be evaluated
	Bad2   int `engine:"localized"` // or localized
}

func (s badTagsStruct1) LocalizationUUID() utils.UUID {
	return utils.UUID("11e2c40c-ae26-448b-a3b2-4c275516bcc0")
}

type badTagsStruct2 struct {
	Bad string `engine:"localized"` // container struct doesn't implement localizable
}

func TestParseEngineTag(t *testing.T) {
	typ1 := reflect.TypeOf(badTagsStruct1{})
	typ2 := reflect.TypeOf(badTagsStruct2{})

	assertTags := func(fieldIndex int, name string, localized bool, evaluated bool) {
		f := typ1.Field(fieldIndex)

		assert.Equal(t, name, jsonNameTag(f))

		actualLocalized, actualEvaluated := parseEngineTag(typ1, f)
		assert.Equal(t, localized, actualLocalized)
		assert.Equal(t, evaluated, actualEvaluated)
	}

	assertTags(0, "valid1", true, false)
	assertTags(1, "valid2", true, true)
	assertTags(2, "valid3", false, true)
	assertTags(3, "", false, false)

	assert.Panics(t, func() { parseEngineTag(typ1, typ1.Field(4)) })
	assert.Panics(t, func() { parseEngineTag(typ1, typ1.Field(5)) })

	assert.Panics(t, func() { parseEngineTag(typ2, typ2.Field(0)) })
}

type nestedStruct struct {
	Foo string `json:"foo" engine:"localized,evaluated"`
}

type testTaggedStruct struct {
	nestedStruct
	Bar string `json:"bar"`
}

func (s testTaggedStruct) LocalizationUUID() utils.UUID {
	return utils.UUID("11e2c40c-ae26-448b-a3b2-4c275516bcc0")
}

func TestExtractEngineFields(t *testing.T) {
	typ := reflect.TypeOf(testTaggedStruct{})

	assert.Equal(t, []*engineField{
		&engineField{jsonName: "foo", localized: true, evaluated: true, index: []int{0, 0}},
		&engineField{jsonName: "bar", localized: false, evaluated: false, index: []int{1}},
	}, extractEngineFields(typ, typ))
}

func TestWalkFields(t *testing.T) {
	// can start with a struct
	v := reflect.ValueOf(testTaggedStruct{nestedStruct: nestedStruct{Foo: "Hello"}, Bar: "World"})

	values := make([]string, 0)
	walk(v, nil, func(sv reflect.Value, fv reflect.Value, ef *engineField) {
		values = append(values, fv.String())
	})

	assert.Equal(t, []string{"Hello", "World"}, values)

	// or a slice of structs
	v = reflect.ValueOf([]testTaggedStruct{
		testTaggedStruct{nestedStruct: nestedStruct{Foo: "Hello"}, Bar: "World"},
		testTaggedStruct{nestedStruct: nestedStruct{Foo: "Hola"}, Bar: "Mundo"},
	})

	values = make([]string, 0)
	walk(v, nil, func(sv reflect.Value, fv reflect.Value, ef *engineField) {
		values = append(values, fv.String())
	})

	assert.Equal(t, []string{"Hello", "World", "Hola", "Mundo"}, values)
}
