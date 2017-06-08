package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// test variable resolver
type resolver struct {
	defaultString string
}

func (r *resolver) Default() interface{} { return r.defaultString }
func (r *resolver) Resolve(key string) interface{} {
	return fmt.Errorf("No such key")
}

// test stringer
type stringer struct{}

func (s *stringer) String() string { return "Stringer" }

func TestToString(t *testing.T) {
	strMap := make(map[string]string)
	strMap["one"] = "1.0"

	chi, err := time.LoadLocation("America/Chicago")
	if err != nil {
		t.Fatal("Unable to load America/Chicago timezone")
	}

	date1 := time.Date(2017, 6, 23, 15, 30, 0, 0, time.UTC)
	date2 := time.Date(2017, 7, 18, 15, 30, 0, 0, chi)

	testStringer := &stringer{}
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
		{testStringer, "Stringer", false},
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

	env := NewDefaultEnvironment()

	for _, test := range tests {
		result, err := ToString(env, test.input)

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

	env := NewDefaultEnvironment()

	for _, test := range tests {
		result, err := ToDecimal(env, test.input)

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
