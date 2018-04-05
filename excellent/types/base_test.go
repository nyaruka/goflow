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

var _ types.XValue = &TestXObject{}

func TestXValuesToStringAndJSON(t *testing.T) {
	chi, err := time.LoadLocation("America/Chicago")
	require.NoError(t, err)

	date1 := time.Date(2017, 6, 23, 15, 30, 0, 0, time.UTC)
	date2 := time.Date(2017, 7, 18, 15, 30, 0, 0, chi)

	tests := []struct {
		value    types.XValue
		asJSON   string
		asString string
		asBool   bool
	}{
		{
			value:    types.NewXString(""),
			asJSON:   `""`,
			asString: "",
			asBool:   false, // empty strings are false
		}, {
			value:    types.NewXString("FALSE"),
			asJSON:   `"FALSE"`,
			asString: "FALSE",
			asBool:   false, // because it's string value is "false"
		}, {
			value:    types.NewXString("hello \"bob\""),
			asJSON:   `"hello \"bob\""`,
			asString: "hello \"bob\"",
			asBool:   true,
		}, {
			value:    types.NewXNumberFromInt(0),
			asJSON:   `0`,
			asString: "0",
			asBool:   false, // because any decimal != 0 is true
		}, {
			value:    types.NewXNumberFromInt(123),
			asJSON:   `123`,
			asString: "123",
			asBool:   true, // because any decimal != 0 is true
		}, {
			value:    types.RequireXNumberFromString("123.00"),
			asJSON:   `123`,
			asString: "123",
			asBool:   true,
		}, {
			value:    types.RequireXNumberFromString("123.45"),
			asJSON:   `123.45`,
			asString: "123.45",
			asBool:   true,
		}, {
			value:    types.NewXBool(false),
			asJSON:   `false`,
			asString: "false",
			asBool:   false,
		}, {
			value:    types.NewXBool(true),
			asJSON:   `true`,
			asString: "true",
			asBool:   true,
		}, {
			value:    types.NewXTime(date1),
			asJSON:   `"2017-06-23T15:30:00.000000Z"`,
			asString: "2017-06-23T15:30:00.000000Z",
			asBool:   true,
		}, {
			value:    types.NewXTime(date2),
			asJSON:   `"2017-07-18T15:30:00.000000-05:00"`,
			asString: "2017-07-18T15:30:00.000000-05:00",
			asBool:   true,
		}, {
			value:    types.NewXError(fmt.Errorf("it failed")),
			asJSON:   `"it failed"`,
			asString: "it failed",
			asBool:   false,
		}, {
			value:    types.NewXArray(),
			asJSON:   `[]`,
			asString: `[]`,
			asBool:   false,
		}, {
			value:    types.NewXArray(types.NewXTime(date1), types.NewXTime(date2)),
			asJSON:   `["2017-06-23T15:30:00.000000Z","2017-07-18T15:30:00.000000-05:00"]`,
			asString: `["2017-06-23T15:30:00.000000Z","2017-07-18T15:30:00.000000-05:00"]`,
			asBool:   true,
		}, {
			value:    NewTestXObject("Hello", 123),
			asJSON:   `{"foo":"Hello","bar":123}`,
			asString: "Hello",
			asBool:   true,
		}, {
			value:    NewTestXObject("", 123),
			asJSON:   `{"foo":"","bar":123}`,
			asString: "",
			asBool:   false, // because it reduces to a string which itself is false
		}, {
			value:    types.NewXArray(NewTestXObject("Hello", 123), NewTestXObject("World", 456)),
			asJSON:   `[{"foo":"Hello","bar":123},{"foo":"World","bar":456}]`,
			asString: `["Hello","World"]`,
			asBool:   true,
		}, {
			value: types.NewXMap(map[string]types.XValue{
				"first":  NewTestXObject("Hello", 123),
				"second": NewTestXObject("World", 456),
			}),
			asJSON:   `{"first":{"foo":"Hello","bar":123},"second":{"foo":"World","bar":456}}`,
			asString: `{"first":"Hello","second":"World"}`,
			asBool:   true,
		},
	}
	for _, test := range tests {
		assert.Equal(t, types.NewXString(test.asJSON), test.value.ToJSON(), "ToJSON failed for %+v", test.value)
		assert.Equal(t, types.NewXString(test.asString), test.value.Reduce().ToString(), "ToString failed for %+v", test.value)
		assert.Equal(t, types.NewXBool(test.asBool), test.value.Reduce().ToBool(), "ToXBool failed for %+v", test.value)
	}
}
