package types_test

import (
	"fmt"
	"testing"

	"github.com/nyaruka/goflow/excellent/types"

	"github.com/stretchr/testify/assert"
)

func TestCompareXValues(t *testing.T) {
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
		result, err := types.CompareXValues(test.x1, test.x2)

		if test.hasError {
			assert.Error(t, err, "expected error for inputs '%s' and '%s'", test.x1, test.x2)
		} else {
			assert.NoError(t, err, "unexpected error for inputs '%s' and '%s'", test.x1, test.x2)
			assert.Equal(t, test.result, result, "result mismatch for inputs '%s' and '%s'", test.x1, test.x2)
		}
	}
}
