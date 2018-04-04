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
		asBool   bool
	}{
		{
			value:    types.NewXString(""),
			asString: "",
			asJSON:   `""`,
			asBool:   false, // empty strings are false
		}, {
			value:    types.NewXString("FALSE"),
			asString: "FALSE",
			asJSON:   `"FALSE"`,
			asBool:   false, // because it's string value is "false"
		}, {
			value:    types.NewXString("hello \"bob\""),
			asString: "hello \"bob\"",
			asJSON:   `"hello \"bob\""`,
			asBool:   true,
		}, {
			value:    types.NewXNumberFromInt(0),
			asString: "0",
			asJSON:   `0`,
			asBool:   false, // because any decimal != 0 is true
		}, {
			value:    types.NewXNumberFromInt(123),
			asString: "123",
			asJSON:   `123`,
			asBool:   true, // because any decimal != 0 is true
		}, {
			value:    types.RequireXNumberFromString("123.00"),
			asString: "123",
			asJSON:   `123`,
			asBool:   true,
		}, {
			value:    types.RequireXNumberFromString("123.45"),
			asString: "123.45",
			asJSON:   `123.45`,
			asBool:   true,
		}, {
			value:    types.NewXBool(false),
			asString: "false",
			asJSON:   `false`,
			asBool:   false,
		}, {
			value:    types.NewXBool(true),
			asString: "true",
			asJSON:   `true`,
			asBool:   true,
		}, {
			value:    types.NewXTime(date1),
			asString: "2017-06-23T15:30:00.000000Z",
			asJSON:   `"2017-06-23T15:30:00.000000Z"`,
			asBool:   true,
		}, {
			value:    types.NewXTime(date2),
			asString: "2017-07-18T15:30:00.000000-05:00",
			asJSON:   `"2017-07-18T15:30:00.000000-05:00"`,
			asBool:   true,
		}, {
			value:    types.NewXError(fmt.Errorf("it failed")),
			asString: "it failed",
			asJSON:   `"it failed"`,
			asBool:   false,
		}, {
			value:    types.NewXArray(),
			asString: `[]`,
			asJSON:   `[]`,
			asBool:   false,
		}, {
			value:    types.NewXArray(types.NewXTime(date1), types.NewXTime(date2)),
			asString: `["2017-06-23T15:30:00.000000Z","2017-07-18T15:30:00.000000-05:00"]`,
			asJSON:   `["2017-06-23T15:30:00.000000Z","2017-07-18T15:30:00.000000-05:00"]`,
			asBool:   true,
		}, {
			value:    NewTestXObject("Hello", 123),
			asString: "Hello",
			asJSON:   `{"foo":"Hello","bar":123}`,
			asBool:   true,
		}, {
			value:    NewTestXObject("", 123),
			asString: "",
			asJSON:   `{"foo":"","bar":123}`,
			asBool:   false, // because it reduces to a string which itself is false
		}, {
			value:    types.NewXArray(NewTestXObject("Hello", 123), NewTestXObject("World", 456)),
			asString: `["Hello","World"]`,
			asJSON:   `[{"foo":"Hello","bar":123},{"foo":"World","bar":456}]`,
			asBool:   true,
		},
	}
	for _, test := range tests {
		assert.Equal(t, types.NewXString(test.asString), types.ToXString(test.value), "ToXString failed for %+v", test.value)
		assert.Equal(t, types.NewXString(test.asJSON), test.value.ToJSON(), "ToJSON failed for %+v", test.value)
		assert.Equal(t, types.NewXBool(test.asBool), types.ToXBool(test.value), "ToXBool failed for %+v", test.value)
	}
}
