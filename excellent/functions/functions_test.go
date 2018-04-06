package functions_test

import (
	//"math"
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	//"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var errorArg = types.NewXErrorf("I am error")
var la, _ = time.LoadLocation("America/Los_Angeles")

var xs = types.NewXString
var xn = types.RequireXNumberFromString
var newDecimal = types.RequireXNumberFromString

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

	{"char", []types.XValue{xn("33")}, xs("!"), false},
	{"char", []types.XValue{xn("128513")}, xs("😁"), false},
	{"char", []types.XValue{xs("not a number")}, types.XStringEmpty, true},
	{"char", []types.XValue{}, types.XStringEmpty, true},

	{"or", []types.XValue{types.XBoolTrue}, types.XBoolTrue, false},
	{"or", []types.XValue{types.XBoolFalse}, types.XBoolFalse, false},
	{"or", []types.XValue{types.XBoolTrue, types.XBoolFalse}, types.XBoolTrue, false},
	{"or", []types.XValue{}, types.XBoolFalse, true},

	{"if", []types.XValue{types.XBoolTrue, xs("10"), xs("20")}, xs("10"), false},
	{"if", []types.XValue{types.XBoolFalse, xs("10"), xs("20")}, xs("20"), false},
	{"if", []types.XValue{types.XBoolTrue, errorArg, xs("20")}, errorArg, true},
	{"if", []types.XValue{}, types.XBoolFalse, true},
	{"if", []types.XValue{errorArg, xs("10"), xs("20")}, xs("20"), false},

	{"round", []types.XValue{xs("10.5"), xs("0")}, xn("11"), false},
	{"round", []types.XValue{xs("10.5"), xs("1")}, xn("10.5"), false},
	{"round", []types.XValue{xs("10.51"), xs("1")}, xn("10.5"), false},
	{"round", []types.XValue{xs("10.56"), xs("1")}, xn("10.6"), false},
	{"round", []types.XValue{xs("12.56"), xs("-1")}, xn("10"), false},
	{"round", []types.XValue{xs("10.5")}, xn("11"), false},
	{"round", []types.XValue{xs("not_num"), xs("1")}, nil, true},
	{"round", []types.XValue{xs("10.5"), xs("not_num")}, nil, true},
	{"round", []types.XValue{xs("10.5"), xs("1"), xs("30")}, nil, true},

	{"round_up", []types.XValue{xs("10.5")}, xn("11"), false},
	{"round_up", []types.XValue{xs("10.2")}, xn("11"), false},
	{"round_up", []types.XValue{xs("not_num")}, nil, true},
	{"round_up", []types.XValue{}, nil, true},

	{"round_down", []types.XValue{xs("10.5")}, xn("10"), false},
	{"round_down", []types.XValue{xs("10.7")}, xn("10"), false},
	{"round_down", []types.XValue{xs("not_num")}, nil, true},
	{"round_down", []types.XValue{}, nil, true},

	{"max", []types.XValue{xs("10.5"), xs("11")}, xn("11"), false},
	{"max", []types.XValue{xs("10.2"), xs("9")}, xn("10.2"), false},
	{"max", []types.XValue{xs("not_num"), xs("9")}, nil, true},
	{"max", []types.XValue{xs("9"), xs("not_num")}, nil, true},
	{"max", []types.XValue{}, nil, true},
	{"min", []types.XValue{xs("10.5"), xs("11")}, xn("10.5"), false},
	{"min", []types.XValue{xs("10.2"), xs("9")}, xn("9"), false},
	{"min", []types.XValue{xs("not_num"), xs("9")}, nil, true},
	{"min", []types.XValue{xs("9"), xs("not_num")}, nil, true},
	{"min", []types.XValue{}, nil, true},

	{"mean", []types.XValue{xs("10"), xs("11")}, xn("10.5"), false},
	{"mean", []types.XValue{xs("10.2")}, xn("10.2"), false},
	{"mean", []types.XValue{xs("not_num")}, nil, true},
	{"mean", []types.XValue{xs("9"), xs("not_num")}, nil, true},
	{"mean", []types.XValue{}, nil, true},

	{"mod", []types.XValue{xs("10"), xs("3")}, xn("1"), false},
	{"mod", []types.XValue{xs("10"), xs("5")}, xn("0"), false},
	{"mod", []types.XValue{xs("not_num"), xs("3")}, nil, true},
	{"mod", []types.XValue{xs("9"), xs("not_num")}, nil, true},
	{"mod", []types.XValue{}, nil, true},

	{"read_code", []types.XValue{xs("123456")}, xs("1 2 3 , 4 5 6"), false},
	{"read_code", []types.XValue{xs("abcd")}, xs("a b c d"), false},
	{"read_code", []types.XValue{xs("12345678")}, xs("1 2 3 4 , 5 6 7 8"), false},
	{"read_code", []types.XValue{xs("12")}, xs("1 , 2"), false},
	{"read_code", []types.XValue{}, nil, true},

	{"split", []types.XValue{xs("1,2,3"), xs(",")}, types.NewXArray(xs("1"), xs("2"), xs("3")), false},
	{"split", []types.XValue{xs("1,2,3"), xs(".")}, types.NewXArray(xs("1,2,3")), false},
	{"split", []types.XValue{xs("1,2,3"), nil}, types.NewXArray(xs("1"), xs(","), xs("2"), xs(","), xs("3")), false},
	{"split", []types.XValue{}, nil, true},

	{"join", []types.XValue{types.NewXArray(xs("1"), xs("2"), xs("3")), xs(",")}, xs("1,2,3"), false},
	{"join", []types.XValue{types.NewXArray(), xs(",")}, xs(""), false},
	{"join", []types.XValue{types.NewXArray(xs("1")), xs(",")}, xs("1"), false},
	{"join", []types.XValue{xs("1,2,3"), nil}, nil, true},
	{"join", []types.XValue{types.NewXArray(xs("1,2,3")), nil}, xs("1,2,3"), false},
	{"join", []types.XValue{types.NewXArray(xs("1"))}, nil, true},

	{"title", []types.XValue{xs("hello")}, xs("Hello"), false},
	{"title", []types.XValue{xs("")}, xs(""), false},
	{"title", []types.XValue{nil}, xs(""), false},
	{"title", []types.XValue{}, nil, true},
}

func TestFunctions(t *testing.T) {
	env := utils.NewEnvironment(utils.DateFormatDayMonthYear, utils.TimeFormatHourMinuteSecond, time.UTC, utils.LanguageList{})

	for _, test := range funcTests {
		xFunc := functions.XFUNCTIONS[test.name]
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
			cmp, err := types.CompareXValues(result, test.expected)
			if err != nil {
				t.Errorf("error while comparing expected: '%#v' with result: '%#v': %v for function %s(%#v)", test.expected, result, err, test.name, test.args)
			}

			if cmp != 0 {
				t.Errorf("unexpected value, expected '%v', got '%v' for function %s(%#v)", test.expected, result, test.name, test.args)
			}
		}
	}
}
