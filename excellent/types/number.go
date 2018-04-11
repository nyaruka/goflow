package types

import (
	"github.com/shopspring/decimal"
)

func init() {
	decimal.MarshalJSONWithoutQuotes = true
}

// XNumber is any whole or fractional number
type XNumber decimal.Decimal

// NewXNumber creates a new XNumber
func NewXNumber(value decimal.Decimal) XNumber {
	return XNumber(value)
}

// NewXNumberFromInt creates a new XNumber from the given int
func NewXNumberFromInt(value int) XNumber {
	return XNumber(decimal.New(int64(value), 0))
}

// NewXNumberFromInt64 creates a new XNumber from the given int
func NewXNumberFromInt64(value int64) XNumber {
	return XNumber(decimal.New(value, 0))
}

// RequireXNumberFromString creates a new XNumber from the given string
func RequireXNumberFromString(value string) XNumber {
	return XNumber(decimal.RequireFromString(value))
}

// Reduce returns the primitive version of this type (i.e. itself)
func (x XNumber) Reduce() XPrimitive { return x }

// ToXString converts this type to a string
func (x XNumber) ToXString() XString { return XString(x.Native().String()) }

// ToXBool converts this type to a bool
func (x XNumber) ToXBool() XBool { return XBool(!x.Native().Equals(decimal.Zero)) }

// ToXJSON converts this type to JSON
func (x XNumber) ToXJSON() XString { return MustMarshalToXString(x.Native()) }

// Native returns the native value of this type
func (x XNumber) Native() decimal.Decimal { return decimal.Decimal(x) }

// Compare compares this number to another
func (x XNumber) Compare(other XNumber) int {
	return x.Native().Cmp(other.Native())
}

// MarshalJSON is called when a struct containing this type is marshaled
func (x XNumber) MarshalJSON() ([]byte, error) {
	nativePtr := (decimal.Decimal)(x)
	return nativePtr.MarshalJSON()
}

// UnmarshalJSON is called when a struct containing this type is unmarshaled
func (x *XNumber) UnmarshalJSON(data []byte) error {
	nativePtr := (*decimal.Decimal)(x)
	return nativePtr.UnmarshalJSON(data)
}

// XNumberZero is the zero number value
var XNumberZero = XNumber(decimal.Zero)
var _ XPrimitive = XNumberZero
