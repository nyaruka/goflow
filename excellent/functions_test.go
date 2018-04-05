package excellent

import (
	"fmt"
	//"math"
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	//"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var errorArg = fmt.Errorf("I am error")
var la, _ = time.LoadLocation("America/Los_Angeles")

var funcTests = []struct {
	name     string
	args     []types.XValue
	expected types.XValue
	hasError bool
}{
	{"and", []types.XValue{types.XBoolTrue}, types.XBoolTrue, false},
	{"and", []types.XValue{types.XBoolFalse}, types.XBoolFalse, false},
	{"and", []types.XValue{types.XBoolTrue, types.XBoolFalse}, types.XBoolFalse, false},
	{"and", []types.XValue{}, types.XBoolFalse, true},

	{"char", []types.XValue{types.NewXNumberFromInt(33)}, types.NewXString("!"), false},
	{"char", []types.XValue{types.NewXNumberFromInt(128513)}, types.NewXString("😁"), false},
	{"char", []types.XValue{types.NewXString("not a number")}, types.XStringEmpty, true},
	{"char", []types.XValue{}, types.XStringEmpty, true},
}

func TestFunctions(t *testing.T) {
	env := utils.NewEnvironment(utils.DateFormatDayMonthYear, utils.TimeFormatHourMinuteSecond, time.UTC, utils.LanguageList{})

	for _, test := range funcTests {
		xFunc := XFUNCTIONS[test.name]
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic running function %s(%#v): %#v", test.name, test.args, r)
			}
		}()

		result := xFunc(env, test.args...)
		err, _ := result.(types.XError)

		if test.hasError {
			assert.Error(t, err, "expected error running function %s(%#v)", test.name, test.args)
		} else {
			assert.NoError(t, err, "unexpected error running function %s(%#v): %s", test.name, test.args, err)

			// and the match itself
			cmp, err := types.Compare(env, result, test.expected)
			if err != nil {
				t.Errorf("error while comparing expected: '%#v' with result: '%#v': %v for function %s(%#v)", test.expected, result, err, test.name, test.args)
			}

			if cmp != 0 {
				t.Errorf("unexpected value, expected '%v', got '%v' for function %s(%#v)", test.expected, result, test.name, test.args)
			}
		}
	}
}
