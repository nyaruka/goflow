package types

import (
	"math"
	"regexp"
	"strings"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// only parse numbers like 123 or 123.456 or .456
var decimalRegexp = regexp.MustCompile(`^-?(([0-9]+)|([0-9]+\.[0-9]+)|(\.[0-9]+))$`)

func init() {
	decimal.MarshalJSONWithoutQuotes = true
}

// XNumber is a whole or fractional number.
//
//   @(1234) -> 1234
//   @(1234.5678) -> 1234.5678
//   @(format_number(1234.5670)) -> 1,234.567
//   @(json(1234.5678)) -> 1234.5678
//
// @type number
type XNumber struct {
	native decimal.Decimal
}

// NewXNumber creates a new XNumber
func NewXNumber(value decimal.Decimal) XNumber {
	return XNumber{native: value}
}

// NewXNumberFromInt creates a new XNumber from the given int
func NewXNumberFromInt(value int) XNumber {
	return NewXNumber(decimal.New(int64(value), 0))
}

// NewXNumberFromInt64 creates a new XNumber from the given int
func NewXNumberFromInt64(value int64) XNumber {
	return NewXNumber(decimal.New(value, 0))
}

// RequireXNumberFromString creates a new XNumber from the given string or panics (used for tests)
func RequireXNumberFromString(value string) XNumber {
	num, err := newXNumberFromString(value)
	if err != nil {
		panic(errors.Wrapf(err, "error parsing '%s' as number", value))
	}
	return num
}

// Describe returns a representation of this type for error messages
func (x XNumber) Describe() string { return x.Render() }

// Truthy determines truthiness for this type
func (x XNumber) Truthy() bool {
	return !x.Equals(XNumberZero)
}

// Render returns the canonical text representation
func (x XNumber) Render() string { return x.Native().String() }

// Format returns the pretty text representation
func (x XNumber) Format(env envs.Environment) string {
	return x.FormatCustom(env.NumberFormat(), -1, true)
}

// FormatCustom provides customised formatting
func (x XNumber) FormatCustom(format *envs.NumberFormat, places int, groupDigits bool) string {
	var formatted string

	if places >= 0 {
		formatted = x.Native().StringFixed(int32(places))
	} else {
		formatted = x.Native().String()
	}

	parts := strings.Split(formatted, ".")

	// add thousands separators
	if groupDigits {
		sb := strings.Builder{}
		for i, r := range parts[0] {
			sb.WriteRune(r)

			d := (len(parts[0]) - 1) - i
			if d%3 == 0 && d > 0 {
				sb.WriteString(format.DigitGroupingSymbol)
			}
		}
		parts[0] = sb.String()
	}

	return strings.Join(parts, format.DecimalSymbol)
}

// String returns the native string representation of this type
func (x XNumber) String() string { return `XNumber(` + x.Render() + `)` }

// Native returns the native value of this type
func (x XNumber) Native() decimal.Decimal { return x.native }

// Equals determines equality for this type
func (x XNumber) Equals(o XValue) bool {
	other := o.(XNumber)

	return x.Native().Equals(other.Native())
}

// Compare compares this number to another
func (x XNumber) Compare(o XValue) int {
	other := o.(XNumber)

	return x.Native().Cmp(other.Native())
}

// MarshalJSON is called when a struct containing this type is marshaled
func (x XNumber) MarshalJSON() ([]byte, error) {
	return x.Native().MarshalJSON()
}

// UnmarshalJSON is called when a struct containing this type is unmarshaled
func (x *XNumber) UnmarshalJSON(data []byte) error {
	nativePtr := &x.native
	return nativePtr.UnmarshalJSON(data)
}

// XNumberZero is the zero number value
var XNumberZero = NewXNumber(decimal.Zero)
var _ XValue = XNumberZero

// parses a number from a string
func newXNumberFromString(s string) (XNumber, error) {
	s = strings.TrimSpace(s)

	if !decimalRegexp.MatchString(s) {
		return XNumberZero, errors.New("not a valid number format")
	}

	// we can assume anything that matched our regex is parseable
	d := decimal.RequireFromString(s)

	return NewXNumber(d), nil
}

// ToXNumber converts the given value to a number or returns an error if that isn't possible
func ToXNumber(env envs.Environment, x XValue) (XNumber, XError) {
	if !utils.IsNil(x) {
		switch typed := x.(type) {
		case XError:
			return XNumberZero, typed
		case XNumber:
			return typed, nil
		case XText:
			parsed, err := newXNumberFromString(typed.Native())
			if err == nil {
				return parsed, nil
			}
		case *XObject:
			if typed.hasDefault() {
				return ToXNumber(env, typed.Default())
			}
		}
	}

	return XNumberZero, NewXErrorf("unable to convert %s to a number", Describe(x))
}

// ToInteger tries to convert the passed in value to an integer or returns an error if that isn't possible
func ToInteger(env envs.Environment, x XValue) (int, XError) {
	number, err := ToXNumber(env, x)
	if err != nil {
		return 0, err
	}

	intPart := number.Native().IntPart()

	if intPart < math.MinInt32 || intPart > math.MaxInt32 {
		return 0, NewXErrorf("number value %s is out of range for an integer", number.Render())
	}

	return int(intPart), nil
}
