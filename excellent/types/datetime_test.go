package types_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestXDateTimeMarshaling(t *testing.T) {
	var date types.XDateTime
	err := json.Unmarshal([]byte(`"2018-04-09T17:01:30Z"`), &date)
	assert.NoError(t, err)
	assert.Equal(t, types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), date)

	// marshal
	data, err := json.Marshal(types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)))
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"2018-04-09T17:01:30Z"`), data)
}

func TestToXDateTime(t *testing.T) {
	var tests = []struct {
		value    types.XValue
		asNumber types.XDateTime
		hasError bool
	}{
		{nil, types.XDateTimeZero, true},
		{types.NewXError(fmt.Errorf("Error")), types.XDateTimeZero, true},
		{types.NewXNumberFromInt(123), types.XDateTimeZero, true},
		{types.NewXText("2018-06-05"), types.NewXDateTime(time.Date(2018, 6, 5, 0, 0, 0, 0, time.UTC)), false},
		{types.NewXText("wha?"), types.XDateTimeZero, true},
		{NewTestXObject("Hello", 123), types.XDateTimeZero, true},
		{NewTestXObject("2018/6/5", 123), types.NewXDateTime(time.Date(2018, 6, 5, 0, 0, 0, 0, time.UTC)), false},
	}

	env := utils.NewDefaultEnvironment()

	for _, test := range tests {
		result, err := types.ToXDateTime(env, test.value)

		if test.hasError {
			assert.Error(t, err, "expected error for input %T{%s}", test.value, test.value)
		} else {
			assert.NoError(t, err, "unexpected error for input %T{%s}", test.value, test.value)
			assert.Equal(t, test.asNumber.Native(), result.Native(), "result mismatch for input %T{%s}", test.value, test.value)
		}
	}
}
