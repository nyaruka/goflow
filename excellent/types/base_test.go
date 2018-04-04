package types_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestXObject struct {
	types.BaseXObject

	foo string
	bar int
}

func NewTestXObject(foo string, bar int) *TestXObject {
	return &TestXObject{foo: foo, bar: bar}
}

// ToJSON converts this type to JSON
func (v *TestXObject) ToJSON() types.XString {
	e := struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}{
		Foo: v.foo,
		Bar: v.bar,
	}
	return types.RequireMarshalToXString(e)
}

func (v *TestXObject) Reduce() types.XPrimitive { return types.NewXString(v.foo) }

var _ types.XObject = &TestXObject{}

func TestXValuesToStringAndJSON(t *testing.T) {
	chi, err := time.LoadLocation("America/Chicago")
	require.NoError(t, err)

	date1 := time.Date(2017, 6, 23, 15, 30, 0, 0, time.UTC)
	date2 := time.Date(2017, 7, 18, 15, 30, 0, 0, chi)

	tests := []struct {
		value    types.XValue
		asString string
		asJSON   string
	}{
		{types.NewXString("hello \"bob\""), "hello \"bob\"", `"hello \"bob\""`},
		{types.NewXNumberFromInt(123), "123", `123`},
		{types.RequireXNumberFromString("123.00"), "123", `123`},
		{types.RequireXNumberFromString("123.45"), "123.45", `123.45`},
		{types.NewXBool(false), "false", `false`},
		{types.NewXBool(true), "true", `true`},
		{types.NewXTime(date1), "2017-06-23T15:30:00.000000Z", `"2017-06-23T15:30:00.000000Z"`},
		{types.NewXTime(date2), "2017-07-18T15:30:00.000000-05:00", `"2017-07-18T15:30:00.000000-05:00"`},
		{types.NewXError(fmt.Errorf("it failed")), "it failed", `"it failed"`},
		{
			types.NewXArray(types.NewXTime(date1), types.NewXTime(date2)),
			`["2017-06-23T15:30:00.000000Z","2017-07-18T15:30:00.000000-05:00"]`,
			`["2017-06-23T15:30:00.000000Z","2017-07-18T15:30:00.000000-05:00"]`,
		},
		{
			NewTestXObject("Hello", 123),
			"Hello",
			`{"foo":"Hello","bar":123}`,
		},
		{
			types.NewXArray(NewTestXObject("Hello", 123), NewTestXObject("World", 456)),
			`["Hello","World"]`,
			`[{"foo":"Hello","bar":123},{"foo":"World","bar":456}]`,
		},
	}
	for _, test := range tests {
		assert.Equal(t, types.NewXString(test.asString), types.ToXString(test.value))
		assert.Equal(t, types.NewXString(test.asJSON), test.value.ToJSON())
	}
}
