package types_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/types"

	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	var tests = []struct {
		x1       types.XValue
		x2       types.XValue
		result   int
		hasError bool
	}{
		{nil, nil, 0, false},
		{nil, types.NewXString(""), 0, true},
		{types.NewXError(fmt.Errorf("Error")), types.NewXError(fmt.Errorf("Error")), 0, false},
		{types.NewXError(fmt.Errorf("Error")), types.XTimeZero, 0, true}, // type mismatch
		{types.NewXString("bob"), types.NewXString("bob"), 0, false},
		{types.NewXString("bob"), types.NewXString("cat"), -1, false},
		{types.NewXString("bob"), types.NewXString("ann"), 1, false},
		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(123), 0, false},
		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(124), -1, false},
		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(122), 1, false},
	}

	for _, test := range tests {
		result, err := types.Compare(test.x1, test.x2)

		if test.hasError {
			assert.Error(t, err, "expected error for inputs '%s' and '%s'", test.x1, test.x2)
		} else {
			assert.NoError(t, err, "unexpected error for inputs '%s' and '%s'", test.x1, test.x2)
			assert.Equal(t, test.result, result, "result mismatch for inputs '%s' and '%s'", test.x1, test.x2)
		}
	}
}

func TestXNumberMarshaling(t *testing.T) {
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

func TestXTimeMarshaling(t *testing.T) {
	var date types.XTime
	err := json.Unmarshal([]byte(`"2018-04-09T17:01:30Z"`), &date)
	assert.NoError(t, err)
	assert.Equal(t, types.NewXTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), date)

	// marshal
	data, err := json.Marshal(types.NewXTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)))
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"2018-04-09T17:01:30Z"`), data)
}
