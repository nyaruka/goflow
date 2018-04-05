package types_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
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

func TestXValueRequiredConversions(t *testing.T) {
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
		}, {
			value:    types.NewXJSONArray([]byte(`[]`)),
			asJSON:   `[]`,
			asString: `[]`,
			asBool:   false,
		}, {
			value:    types.NewXJSONArray([]byte(`[5, "x"]`)),
			asJSON:   `[5, "x"]`,
			asString: `[5, "x"]`,
			asBool:   true,
		}, {
			value:    types.NewXJSONObject([]byte(`{}`)),
			asJSON:   `{}`,
			asString: `{}`,
			asBool:   false,
		}, {
			value:    types.NewXJSONObject([]byte(`{"foo":"World","bar":456}`)),
			asJSON:   `{"foo":"World","bar":456}`,
			asString: `{"foo":"World","bar":456}`,
			asBool:   true,
		},
	}
	for _, test := range tests {
		assert.Equal(t, types.NewXString(test.asJSON), types.ToXJSON(test.value), "ToXJSON failed for %+v", test.value)
		assert.Equal(t, types.NewXString(test.asString), types.ToXString(test.value), "ToXString failed for %+v", test.value)
		assert.Equal(t, types.NewXBool(test.asBool), types.ToXBool(test.value), "ToXBool failed for %+v", test.value)
	}
}

func TestToXNumber(t *testing.T) {
	var tests = []struct {
		value    types.XValue
		asNumber types.XNumber
		hasError bool
	}{
		{nil, types.XNumberZero, false},
		{types.NewXError(fmt.Errorf("Error")), types.XNumberZero, true},
		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(123), false},
		{types.NewXString("15.5"), types.RequireXNumberFromString("15.5"), false},
		{types.NewXString("lO.5"), types.RequireXNumberFromString("10.5"), false},
		{NewTestXObject("Hello", 123), types.XNumberZero, true},
		{NewTestXObject("123.45000", 123), types.RequireXNumberFromString("123.45"), false},
	}

	for _, test := range tests {
		result, err := types.ToXNumber(test.value)

		if test.hasError {
			assert.Error(t, err, "expected error for input '%s'", test.value)
		} else {
			assert.NoError(t, err, "unexpected error for input '%s'", test.value)
			assert.Equal(t, test.asNumber.Native(), result.Native(), "result mismatch for input '%+v'", test.value)
		}
	}
}

func TestToXTime(t *testing.T) {
	var tests = []struct {
		value    types.XValue
		asNumber types.XTime
		hasError bool
	}{
		{nil, types.XTimeZero, false},
		{types.NewXError(fmt.Errorf("Error")), types.XTimeZero, true},
		{types.NewXNumberFromInt(123), types.XTimeZero, true},
		{types.NewXString("2018-06-05"), types.NewXTime(time.Date(2018, 6, 5, 0, 0, 0, 0, time.UTC)), false},
		{types.NewXString("wha?"), types.XTimeZero, true},
		{NewTestXObject("Hello", 123), types.XTimeZero, true},
		{NewTestXObject("2018/6/5", 123), types.NewXTime(time.Date(2018, 6, 5, 0, 0, 0, 0, time.UTC)), false},
	}

	env := utils.NewDefaultEnvironment()

	for _, test := range tests {
		result, err := types.ToXTime(env, test.value)

		if test.hasError {
			assert.Error(t, err, "expected error for input '%s'", test.value)
		} else {
			assert.NoError(t, err, "unexpected error for input '%s'", test.value)
			assert.Equal(t, test.asNumber.Native(), result.Native(), "result mismatch for input '%+v'", test.value)
		}
	}
}

// Legacy...

// test variable resolver
type resolver struct {
	defaultString string
}

func (r *resolver) Atomize() interface{} { return r.defaultString }
func (r *resolver) Resolve(key string) interface{} {
	return fmt.Errorf("No such key")
}

func TestToString(t *testing.T) {
	chi, err := time.LoadLocation("America/Chicago")
	if err != nil {
		t.Fatal("Unable to load America/Chicago timezone")
	}

	date1 := time.Date(2017, 6, 23, 15, 30, 0, 0, time.UTC)
	date2 := time.Date(2017, 7, 18, 15, 30, 0, 0, chi)

	testResolver := &resolver{"Resolver"}

	var tests = []struct {
		input    interface{}
		expected string
		hasError bool
	}{
		{nil, "", false},
		{fmt.Errorf("Error"), "", true},
		{"string1", "string1", false},
		{true, "true", false},
		{int(15), "15", false},
		{decimal.NewFromFloat(15.5), "15.5", false},
		{testResolver, "Resolver", false},
		{date1, "2017-06-23T15:30:00.000000Z", false},
		{types.NewArray(date1, date2), "2017-06-23T15:30:00.000000Z, 2017-07-18T15:30:00.000000-05:00", false},
		{types.NewArray("one", "two", "three"), "one, two, three", false},
		{types.NewArray(true, false, true), "true, false, true", false},
		{types.NewArray(decimal.NewFromFloat(1.5), decimal.NewFromFloat(2.5)), "1.5, 2.5", false},
		{types.NewArray(5, -10, 15), "5, -10, 15", false},
		{struct{}{}, "", true},
	}

	env := utils.NewDefaultEnvironment()

	for _, test := range tests {
		result, err := types.ToString(env, test.input)

		if err != nil && !test.hasError {
			t.Errorf("Unexpected error calling ToString on '%v': %s", test.input, err)
		}

		if err == nil && test.hasError {
			t.Errorf("Did not receive expected error calling ToString on '%v': %s", test.input, err)
		}

		if result != test.expected {
			t.Errorf("Unexpected result calling ToString on '%v', got: %s expected: %s", test.input, result, test.expected)
		}
	}
}

func TestToDecimal(t *testing.T) {
	testResolver := &resolver{"155"}

	var tests = []struct {
		input    interface{}
		expected decimal.Decimal
		hasError bool
	}{
		{nil, decimal.NewFromFloat(0), false},
		{fmt.Errorf("Error"), decimal.Zero, true},
		{decimal.NewFromFloat(42), decimal.NewFromFloat(42), false},
		{int(15), decimal.NewFromFloat(15), false},
		{"15.5", decimal.NewFromFloat(15.5), false},
		{"lO.5", decimal.NewFromFloat(10.5), false},
		{testResolver, decimal.NewFromFloat(155), false},
		{struct{}{}, decimal.NewFromFloat(0), true},
	}

	env := utils.NewDefaultEnvironment()

	for _, test := range tests {
		result, err := types.ToDecimal(env, test.input)

		if err != nil && !test.hasError {
			t.Errorf("Unexpected error calling ToDecimal on '%v': %s", test.input, err)
		}

		if err == nil && test.hasError {
			t.Errorf("Did not receive expected error calling ToDecimal on '%v': %s", test.input, err)
		}

		if !result.Equals(test.expected) {
			t.Errorf("Unexpected result calling ToDecimal on '%v', got: %s expected: %s", test.input, result, test.expected)
		}
	}
}

func TestToBool(t *testing.T) {
	testResolver := &resolver{"155"}

	var tests = []struct {
		input    interface{}
		expected bool
		hasError bool
	}{
		{nil, false, false},
		{fmt.Errorf("Error"), false, true},
		{decimal.NewFromFloat(42), true, false},
		{int(0), false, false},
		{int(15), true, false},
		{"15.5", true, false},
		{"lO.5", true, false},
		{"", false, false},
		{testResolver, true, false},
		{types.JSONFragment([]byte(`false`)), false, false},
		{types.JSONFragment([]byte(`true`)), true, false},
		{types.JSONFragment([]byte(`[]`)), false, false},
		{types.JSONFragment([]byte(`15.5`)), true, false},
		{types.JSONFragment([]byte(`0`)), false, false},
		{types.JSONFragment([]byte(`[5]`)), true, false},
		{types.JSONFragment([]byte("{\n}")), false, false},
		{types.JSONFragment([]byte(`{"one": "two"}`)), true, false},
		{struct{}{}, false, true},
	}

	env := utils.NewDefaultEnvironment()

	for _, test := range tests {
		result, err := types.ToBool(env, test.input)

		if err != nil && !test.hasError {
			t.Errorf("Unexpected error calling ToBool on '%v': %s", test.input, err)
		}

		if err == nil && test.hasError {
			t.Errorf("Did not receive expected error calling ToBool on '%v': %s", test.input, err)
		}

		if result != test.expected {
			t.Errorf("Unexpected result calling ToBool on '%v', got: %t expected: %t", test.input, result, test.expected)
		}
	}
}

func TestToJSON(t *testing.T) {
	strMap := make(map[string]string)
	strMap["one"] = "1.0"

	chi, err := time.LoadLocation("America/Chicago")
	if err != nil {
		t.Fatal("Unable to load America/Chicago timezone")
	}

	date1 := time.Date(2017, 6, 23, 15, 30, 0, 0, time.UTC)
	date2 := time.Date(2017, 7, 18, 15, 30, 0, 0, chi)

	testResolver := &resolver{"Resolver"}

	var tests = []struct {
		input    interface{}
		expected string
		hasError bool
	}{
		{nil, "null", false},
		{fmt.Errorf("Error"), "", true},
		{"string1", `"string1"`, false},
		{true, "true", false},
		{int(15), "15", false},
		{decimal.NewFromFloat(15.5), "15.5", false},
		{testResolver, `"Resolver"`, false},
		{date1, `"2017-06-23T15:30:00.000000Z"`, false},
		{types.NewArray(date1, date2), `["2017-06-23T15:30:00Z","2017-07-18T15:30:00-05:00"]`, false},
		{types.NewArray("one", "two", "three"), `["one","two","three"]`, false},
		{types.NewArray(true, false, true), `[true,false,true]`, false},
		{types.NewArray(decimal.NewFromFloat(1.5), decimal.NewFromFloat(2.5)), `[1.5,2.5]`, false},
		{types.NewArray(5, -10, 15), `[5,-10,15]`, false},
		{struct{}{}, "", true},
	}

	env := utils.NewDefaultEnvironment()

	for _, test := range tests {
		fragment, err := types.ToJSON(env, test.input)
		result := string(fragment)

		if err != nil && !test.hasError {
			t.Errorf("Unexpected error calling ToJSON on '%v': %s", test.input, err)
		}

		if err == nil && test.hasError {
			t.Errorf("Did not receive expected error calling ToJSON on '%v': %s", test.input, err)
		}

		if result != test.expected {
			t.Errorf("Unexpected result calling ToJSON on '%v', got: %s expected: %s", test.input, result, test.expected)
		}
	}
}
