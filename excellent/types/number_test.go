package types_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/excellent/types"

	"github.com/stretchr/testify/assert"
)

func TestXNumber(t *testing.T) {
	// test creation
	assert.Equal(t, types.RequireXNumberFromString("123"), types.NewXNumberFromInt(123))
	assert.Equal(t, types.RequireXNumberFromString("123"), types.NewXNumberFromInt64(123))

	// test equality
	assert.True(t, types.NewXNumberFromInt(123).Equals(types.NewXNumberFromInt(123)))
	assert.False(t, types.NewXNumberFromInt(123).Equals(types.NewXNumberFromInt(124)))

	// test comparison
	assert.Equal(t, 0, types.NewXNumberFromInt(123).Compare(types.NewXNumberFromInt(123)))
	assert.Equal(t, -1, types.NewXNumberFromInt(123).Compare(types.NewXNumberFromInt(124)))
	assert.Equal(t, 1, types.NewXNumberFromInt(124).Compare(types.NewXNumberFromInt(123)))

	// unmarshal with quotes
	var num types.XNumber
	err := json.Unmarshal([]byte(`"23.45"`), &num)
	assert.NoError(t, err)
	assert.Equal(t, types.RequireXNumberFromString("23.45"), num)

	// unmarshal without quotes
	err = json.Unmarshal([]byte(`34.56`), &num)
	assert.NoError(t, err)
	assert.Equal(t, types.RequireXNumberFromString("34.56"), num)

	// marshal (doesn't use quotes)
	data, err := json.Marshal(types.RequireXNumberFromString("23.45"))
	assert.NoError(t, err)
	assert.Equal(t, []byte(`23.45`), data)
}

func TestToXNumberAndInteger(t *testing.T) {
	var tests = []struct {
		value     types.XValue
		asNumber  types.XNumber
		asInteger int
		hasError  bool
	}{
		{nil, types.XNumberZero, 0, true},
		{types.NewXErrorf("Error"), types.XNumberZero, 0, true},
		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(123), 123, false},
		{types.NewXText("15.5"), types.RequireXNumberFromString("15.5"), 15, false},
		{types.NewXText("12345678901234567890"), types.RequireXNumberFromString("12345678901234567890"), 0, true}, // out of int range
		{NewTestXObject("Hello", 123), types.XNumberZero, 0, true},
		{NewTestXObject("123.45000", 123), types.RequireXNumberFromString("123.45"), 123, false},
	}

	for _, test := range tests {
		number, err := types.ToXNumber(test.value)
		integer, err := types.ToInteger(test.value)

		if test.hasError {
			assert.Error(t, err, "expected error for input %T{%s}", test.value, test.value)
		} else {
			assert.NoError(t, err, "unexpected error for input %T{%s}", test.value, test.value)
			assert.Equal(t, test.asNumber.Native(), number.Native(), "number mismatch for input %T{%s}", test.value, test.value)
			assert.Equal(t, test.asInteger, integer, "integer mismatch for input %T{%s}", test.value, test.value)
		}
	}
}
