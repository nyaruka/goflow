package utils_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

// test variable resolver
type resolver struct {
	defaultString string
}

func (r *resolver) Atomize() interface{} { return r.defaultString }
func (r *resolver) Resolve(key string) interface{} {
	return fmt.Errorf("No such key")
}

func TestToString(t *testing.T) {
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
		{nil, "", false},
		{fmt.Errorf("Error"), "", true},
		{"string1", "string1", false},
		{true, "true", false},
		{int(15), "15", false},
		{int32(15), "15", false},
		{int64(15), "15", false},
		{float32(15.5), "15.5", false},
		{float64(15.5), "15.5", false},
		{decimal.NewFromFloat(15.5), "15.5", false},
		{testResolver, "Resolver", false},
		{date1, "2017-06-23T15:30:00.000000Z", false},
		{[]time.Time{date1, date2}, "2017-06-23T15:30:00.000000Z, 2017-07-18T15:30:00.000000-05:00", false},
		{[]string{"one", "two", "three"}, "one, two, three", false},
		{[]bool{true, false, true}, "true, false, true", false},
		{[]decimal.Decimal{decimal.NewFromFloat(1.5), decimal.NewFromFloat(2.5)}, "1.5, 2.5", false},
		{[]int{5, -10, 15}, "5, -10, 15", false},
		{strMap, "{\"one\":\"1.0\"}", false},
		{struct{}{}, "", true},
	}

	env := utils.NewDefaultEnvironment()

	for _, test := range tests {
		result, err := utils.ToString(env, test.input)

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
		{int32(15), decimal.NewFromFloat(15), false},
		{int64(15), decimal.NewFromFloat(15), false},
		{float32(15.5), decimal.NewFromFloat(15.5), false},
		{float64(15.5), decimal.NewFromFloat(15.5), false},
		{"15.5", decimal.NewFromFloat(15.5), false},
		{"lO.5", decimal.NewFromFloat(10.5), false},
		{testResolver, decimal.NewFromFloat(155), false},
		{struct{}{}, decimal.NewFromFloat(0), true},
	}

	env := utils.NewDefaultEnvironment()

	for _, test := range tests {
		result, err := utils.ToDecimal(env, test.input)

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
		{int32(15), true, false},
		{int64(15), true, false},
		{float32(15.5), true, false},
		{float64(15.5), true, false},
		{"15.5", true, false},
		{"lO.5", true, false},
		{"", false, false},
		{testResolver, true, false},
		{utils.JSONFragment([]byte(`false`)), false, false},
		{utils.JSONFragment([]byte(`true`)), true, false},
		{utils.JSONFragment([]byte(`[]`)), false, false},
		{utils.JSONFragment([]byte(`15.5`)), true, false},
		{utils.JSONFragment([]byte(`0`)), false, false},
		{utils.JSONFragment([]byte(`[5]`)), true, false},
		{utils.JSONFragment([]byte("{\n}")), false, false},
		{utils.JSONFragment([]byte(`{"one": "two"}`)), true, false},
		{struct{}{}, false, true},
	}

	env := utils.NewDefaultEnvironment()

	for _, test := range tests {
		result, err := utils.ToBool(env, test.input)

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
		{int32(15), "15", false},
		{int64(15), "15", false},
		{float32(15.5), "15.5", false},
		{float64(15.5), "15.5", false},
		{decimal.NewFromFloat(15.5), "15.5", false},
		{testResolver, `"Resolver"`, false},
		{date1, `"2017-06-23T15:30:00.000000Z"`, false},
		{[]time.Time{date1, date2}, `["2017-06-23T15:30:00.000000Z","2017-07-18T15:30:00.000000-05:00"]`, false},
		{[]string{"one", "two", "three"}, `["one","two","three"]`, false},
		{[]bool{true, false, true}, `[true,false,true]`, false},
		{[]decimal.Decimal{decimal.NewFromFloat(1.5), decimal.NewFromFloat(2.5)}, `["1.5","2.5"]`, false},
		{[]int{5, -10, 15}, `[5,-10,15]`, false},
		{strMap, `{"one":"1.0"}`, false},
		{struct{}{}, "", true},
	}

	env := utils.NewDefaultEnvironment()

	for _, test := range tests {
		fragment, err := utils.ToJSON(env, test.input)
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
