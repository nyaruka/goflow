package inspect

import (
	"reflect"
	"testing"

	"github.com/nyaruka/goflow/utils/uuids"

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

func (s badTagsStruct1) LocalizationUUID() uuids.UUID {
	return uuids.UUID("11e2c40c-ae26-448b-a3b2-4c275516bcc0")
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

type nestedFieldsStruct struct {
	Foo string `json:"foo" engine:"localized,evaluated"`
}

type taggedFieldsStruct struct {
	nestedFieldsStruct
	Bar string `json:"bar"`
}

func (s taggedFieldsStruct) LocalizationUUID() uuids.UUID {
	return uuids.UUID("11e2c40c-ae26-448b-a3b2-4c275516bcc0")
}

func TestExtractEngineFields(t *testing.T) {
	v := taggedFieldsStruct{nestedFieldsStruct: nestedFieldsStruct{Foo: "Hello"}, Bar: "World"}
	typ := reflect.TypeOf(v)

	fields := extractEngineFields(typ, typ)

	assert.Equal(t, "foo", fields[0].JSONName)
	assert.True(t, fields[0].Localized)
	assert.True(t, fields[0].Evaluated)
	assert.Equal(t, "Hello", fields[0].Getter(reflect.ValueOf(v)).String())

	assert.Equal(t, "bar", fields[1].JSONName)
	assert.False(t, fields[1].Localized)
	assert.False(t, fields[1].Evaluated)
	assert.Equal(t, "World", fields[1].Getter(reflect.ValueOf(v)).String())
}

type embeddedType struct {
	Foo string `json:"foo" engine:"localized,evaluated"`
}

type subType struct {
	Zed string `json:"zed"`
}

type containerStruct struct {
	embeddedType
	Bar   string    `json:"bar"`
	Sub   subType   `json:"sub"`
	Slice []subType `json:"slice"`
}

func (s containerStruct) LocalizationUUID() uuids.UUID {
	return uuids.UUID("11e2c40c-ae26-448b-a3b2-4c275516bcc0")
}

func TestWalk(t *testing.T) {
	// can start with a struct
	v := reflect.ValueOf(
		containerStruct{
			embeddedType: embeddedType{Foo: "Hello"},
			Bar:          "World",
			Sub:          subType{Zed: "Now"},
			Slice:        []subType{},
		})

	values := make([]interface{}, 0)
	walk(v, nil, func(sv reflect.Value, fv reflect.Value, ef *EngineField) {
		values = append(values, fv.Interface())
	})

	assert.Equal(t, []interface{}{"Hello", "World", subType{Zed: "Now"}, "Now", []subType{}}, values)

	// or a slice of structs
	v = reflect.ValueOf([]containerStruct{
		containerStruct{
			embeddedType: embeddedType{Foo: "Hello"},
			Bar:          "World",
			Sub:          subType{Zed: "Now"},
			Slice:        []subType{},
		},
		containerStruct{
			embeddedType: embeddedType{Foo: "Hola"},
			Bar:          "Mundo",
			Sub:          subType{Zed: "Ahora"},
			Slice:        []subType{},
		},
	})

	values = make([]interface{}, 0)
	walk(v, nil, func(sv reflect.Value, fv reflect.Value, ef *EngineField) {
		values = append(values, fv.Interface())
	})

	assert.Equal(t, []interface{}{
		"Hello",
		"World",
		subType{Zed: "Now"}, "Now",
		[]subType{},
		"Hola",
		"Mundo",
		subType{Zed: "Ahora"},
		"Ahora",
		[]subType{},
	}, values)
}

func TestWalkTypes(t *testing.T) {
	// can start with a struct
	typ := reflect.TypeOf(containerStruct{embeddedType: embeddedType{Foo: "Hello"}, Bar: "World"})

	paths := make([]string, 0)
	walkTypes(typ, "", func(path string, ef *EngineField) {
		paths = append(paths, path)
	})

	assert.Equal(t, []string{".foo", ".bar", ".sub", ".sub.zed", ".slice", ".slice[*].zed"}, paths)
}
