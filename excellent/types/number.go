package types

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/random"
	"github.com/nyaruka/goflow/envs"
	"github.com/shopspring/decimal"
)

// only parse numbers like 123 or 123.456 or .456
var decimalRegexp = regexp.MustCompile(`^-?(([0-9]+)|([0-9]+\.[0-9]+)|(\.[0-9]+))$`)

// MaxNumberDigits is the maximum number of significant digits in a number
const MaxNumberDigits = 36

// maxRoundPlaces bounds the places argument to Round, RoundUp and RoundDown. Numbers are limited to
// MaxNumberDigits significant digits and a magnitude of ±1E100, so rounding to more (or fewer negative)
// places than this can never change a representable value - but shopspring's rescale allocates a 10^places
// big.Int, so a pathological places value like 2 billion would allocate hundreds of megabytes first.
// Clamping to this range leaves every valid result unchanged whilst keeping the allocation tiny.
const maxRoundPlaces = 1_000

// clampRoundPlaces bounds a rounding places argument to ±maxRoundPlaces.
func clampRoundPlaces(places int) int {
	return min(max(places, -maxRoundPlaces), maxRoundPlaces)
}

// maxExponentMagnitude is the largest magnitude allowed for an exponent passed to Pow. A number can't exceed
// ±1E100, so raising anything to a larger power can only ever produce an out-of-range result - and computing
// it first would burn CPU and memory that grow exponentially with the exponent. This is a templating engine,
// not a calculator, so exponents are held to the same magnitude limit as numbers themselves.
var maxExponentMagnitude = decimal.New(100, 0)

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

// RandomXNumber creates a new random XNumber in the range [0, 1)
func RandomXNumber() *XNumber {
	return newXNumber(random.Decimal())
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

// Add returns the sum of this number and the given number, or an error if the result is out of range
func (x *XNumber) Add(o *XNumber) (*XNumber, error) {
	return checkedXNumber(x.native.Add(o.native))
}

// Sub returns the difference of this number and the given number, or an error if the result is out of range
func (x *XNumber) Sub(o *XNumber) (*XNumber, error) {
	return checkedXNumber(x.native.Sub(o.native))
}

// Mul returns the product of this number and the given number, or an error if the result is out of range
func (x *XNumber) Mul(o *XNumber) (*XNumber, error) {
	return checkedXNumber(x.native.Mul(o.native))
}

// Div returns this number divided by the given number, or an error if the divisor is zero or the
// result is out of range
func (x *XNumber) Div(o *XNumber) (*XNumber, error) {
	if o.native.IsZero() {
		return nil, errors.New("division by zero")
	}
	return checkedXNumber(x.native.Div(o.native))
}

// Mod returns the remainder of the division of this number by the given number, or an error if the
// divisor is zero or the result is out of range
func (x *XNumber) Mod(o *XNumber) (*XNumber, error) {
	if o.native.IsZero() {
		return nil, errors.New("division by zero")
	}
	return checkedXNumber(x.native.Mod(o.native))
}

// Pow returns this number raised to the power of the given number, or an error if the exponent is too large
// or the result is out of range. The exponent is bounded because raising to a power grows the result - and
// shopspring's intermediate values - exponentially with the exponent, so a large exponent burns huge amounts
// of CPU and memory before the result's range is ever checked.
func (x *XNumber) Pow(o *XNumber) (*XNumber, error) {
	if o.native.Abs().GreaterThan(maxExponentMagnitude) {
		return nil, errors.New("number value out of range")
	}

	return checkedXNumber(x.native.Pow(o.native))
}

// Neg returns the negation of this number
func (x *XNumber) Neg() *XNumber {
	return newXNumber(x.native.Neg())
}

// Abs returns the absolute value of this number
func (x *XNumber) Abs() *XNumber {
	return newXNumber(x.native.Abs())
}

// Floor returns the nearest integer value less than or equal to this number
func (x *XNumber) Floor() *XNumber {
	return newXNumber(x.native.Floor())
}

// IntPart returns the integer component of this number as an int64
func (x *XNumber) IntPart() int64 {
	return x.native.IntPart()
}

// Round rounds this number to the given number of decimal places. If places < 0 it will round the
// integer part to the nearest 10^(-places). Returns an error if the result is out of range.
func (x *XNumber) Round(places int) (*XNumber, error) {
	places = clampRoundPlaces(places)
	return checkedXNumber(x.native.Round(int32(places)))
}

// RoundUp rounds this number up (towards positive infinity) to the given number of decimal places.
// Returns an error if the result is out of range.
func (x *XNumber) RoundUp(places int) (*XNumber, error) {
	places = clampRoundPlaces(places)
	if x.native.Round(int32(places)).Equal(x.native) {
		return x, nil
	}

	halfPrecision := decimal.New(5, -int32(places)-1)

	return checkedXNumber(x.native.Add(halfPrecision).Round(int32(places)))
}

// RoundDown rounds this number down (towards negative infinity) to the given number of decimal places.
// Returns an error if the result is out of range.
func (x *XNumber) RoundDown(places int) (*XNumber, error) {
	places = clampRoundPlaces(places)
	if x.native.Round(int32(places)).Equal(x.native) {
		return x, nil
	}

	halfPrecision := decimal.New(5, -int32(places)-1)

	return checkedXNumber(x.native.Sub(halfPrecision).Round(int32(places)))
}

// creates a new XNumber from the result of an arithmetic operation, checking that it is in range
func checkedXNumber(d decimal.Decimal) (*XNumber, error) {
	if err := CheckDecimalRange(d); err != nil {
		return nil, errors.New("number value out of range")
	}
	return newXNumber(d), nil
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
