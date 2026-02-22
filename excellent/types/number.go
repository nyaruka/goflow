package types

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/shopspring/decimal"
)

// only parse numbers like 123 or 123.456 or .456
var decimalRegexp = regexp.MustCompile(`^-?(([0-9]+)|([0-9]+\.[0-9]+)|(\.[0-9]+))$`)

// MaxNumberDigits is the maximum number of significant digits in a number
const MaxNumberDigits = 36

func init() {
	decimal.MarshalJSONWithoutQuotes = true
}

// XNumber is a whole or fractional number.
//
//	@(1234) -> 1234
//	@(1234.5678) -> 1234.5678
//	@(format_number(1234.5670)) -> 1,234.567
//	@(json(1234.5678)) -> 1234.5678
//
// @type number
type XNumber struct {
	baseValue

	native decimal.Decimal
}

// newXNumber creates a new XNumber without range checking - for use with known-safe values
func newXNumber(value decimal.Decimal) *XNumber {
	return &XNumber{native: value}
}

// NewXNumber creates a new XNumber from the given decimal value, returning an error if the value
// is outside the range of values that can be persisted
func NewXNumber(value decimal.Decimal) XValue {
	if err := CheckDecimalRange(value); err != nil {
		return NewXErrorf("number value out of range")
	}
	return newXNumber(value)
}

// NewXNumberFromInt creates a new XNumber from the given int
func NewXNumberFromInt(value int) *XNumber {
	return newXNumber(decimal.New(int64(value), 0))
}

// NewXNumberFromInt64 creates a new XNumber from the given int
func NewXNumberFromInt64(value int64) *XNumber {
	return newXNumber(decimal.New(value, 0))
}

// RequireXNumberFromString creates a new XNumber from the given string or panics (used for tests)
func RequireXNumberFromString(value string) *XNumber {
	num, err := newXNumberFromString(value)
	if err != nil {
		panic(fmt.Errorf("error parsing '%s' as number: %w", value, err))
	}
	return num
}

// Describe returns a representation of this type for error messages
func (x *XNumber) Describe() string { return x.Render() }

// Truthy determines truthiness for this type
func (x *XNumber) Truthy() bool {
	return !x.Equals(XNumberZero)
}

// Render returns the canonical text representation
func (x *XNumber) Render() string { return x.Native().String() }

// Format returns the pretty text representation
func (x *XNumber) Format(env envs.Environment) string {
	return x.FormatCustom(env.NumberFormat(), -1, true)
}

// FormatCustom provides customised formatting
func (x *XNumber) FormatCustom(format *envs.NumberFormat, places int, groupDigits bool) string {
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
func (x *XNumber) String() string { return `XNumber(` + x.Render() + `)` }

// Native returns the native value of this type
func (x *XNumber) Native() decimal.Decimal { return x.native }

// Equals determines equality for this type
func (x *XNumber) Equals(o XValue) bool {
	other := o.(*XNumber)

	return x.Native().Equal(other.Native())
}

// Compare compares this number to another
func (x *XNumber) Compare(o XValue) int {
	other := o.(*XNumber)

	return x.Native().Cmp(other.Native())
}

// MarshalJSON is called when a struct containing this type is marshaled
func (x *XNumber) MarshalJSON() ([]byte, error) {
	return x.Native().MarshalJSON()
}

// UnmarshalJSON is called when a struct containing this type is unmarshaled
func (x *XNumber) UnmarshalJSON(data []byte) error {
	return jsonx.Unmarshal(data, &x.native)
}

// XNumberZero is the zero number value
var XNumberZero = newXNumber(decimal.Zero)
var _ XValue = XNumberZero

// CheckDecimalRange checks that the given decimal value is within the range of values that can be
// persisted to our database. It enforces two constraints:
//
//  1. The number of significant digits (excluding trailing zeros) must not exceed MaxNumberDigits (36).
//  2. The magnitude (adjusted exponent) must be within ±100.
//
// Trailing zeros in the coefficient are not counted as significant, so a number like
// 1234567895171680000000000000000000000000 (15 significant digits) is valid despite having 40 total digits.
func CheckDecimalRange(d decimal.Decimal) error {
	if d.IsZero() {
		return nil
	}

	// count significant digits by removing trailing zeros from the coefficient
	s := d.Coefficient().String()
	s = strings.TrimLeft(s, "-")
	s = strings.TrimRight(s, "0")
	if len(s) > MaxNumberDigits {
		return errors.New("number has too many digits")
	}

	adjExp := int64(d.Exponent()) + int64(d.NumDigits()) - 1
	if adjExp > 100 || adjExp < -100 {
		return errors.New("number value is out of permitted range")
	}

	return nil
}

// NewXNumberFromString parses a number from a string
func NewXNumberFromString(s string) (*XNumber, error) {
	return newXNumberFromString(s)
}

// parses a number from a string
func newXNumberFromString(s string) (*XNumber, error) {
	s = strings.TrimSpace(s)

	if !decimalRegexp.MatchString(s) {
		return XNumberZero, errors.New("not a valid number format")
	}

	// we can assume anything that matched our regex is parseable
	d := decimal.RequireFromString(s)

	if err := CheckDecimalRange(d); err != nil {
		return XNumberZero, err
	}

	return newXNumber(d), nil
}

// ToXNumber converts the given value to a number or returns an error if that isn't possible
func ToXNumber(env envs.Environment, x XValue) (*XNumber, *XError) {
	if !IsNil(x) {
		switch typed := x.(type) {
		case *XError:
			return XNumberZero, typed
		case *XNumber:
			return typed, nil
		case *XText:
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
func ToInteger(env envs.Environment, x XValue) (int, *XError) {
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
