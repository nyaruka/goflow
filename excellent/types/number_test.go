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

func TestToXNumber(t *testing.T) {
	var tests = []struct {
		value    types.XValue
		asNumber types.XNumber
		hasError bool
	}{
		{nil, types.XNumberZero, true},
		{types.NewXErrorf("Error"), types.XNumberZero, true},
		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(123), false},
		{types.NewXText("15.5"), types.RequireXNumberFromString("15.5"), false},
		{types.NewXText("lO.5"), types.RequireXNumberFromString("10.5"), false},
		{NewTestXObject("Hello", 123), types.XNumberZero, true},
		{NewTestXObject("123.45000", 123), types.RequireXNumberFromString("123.45"), false},
	}

	for _, test := range tests {
		result, err := types.ToXNumber(test.value)

		if test.hasError {
			assert.Error(t, err, "expected error for input %T{%s}", test.value, test.value)
		} else {
			assert.NoError(t, err, "unexpected error for input %T{%s}", test.value, test.value)
			assert.Equal(t, test.asNumber.Native(), result.Native(), "result mismatch for input %T{%s}", test.value, test.value)
		}
	}
}
