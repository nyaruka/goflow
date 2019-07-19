package inspect

import (
	"reflect"
	"testing"

	"github.com/nyaruka/goflow/utils/uuids"
	"github.com/stretchr/testify/assert"
)

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
