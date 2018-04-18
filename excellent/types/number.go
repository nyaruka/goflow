package types

import (
	"github.com/shopspring/decimal"
)

func init() {
	decimal.MarshalJSONWithoutQuotes = true
}

// XNumber is any whole or fractional number
type XNumber struct {
	baseXPrimitive

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

// RequireXNumberFromString creates a new XNumber from the given string
func RequireXNumberFromString(value string) XNumber {
	return NewXNumber(decimal.RequireFromString(value))
}

// Reduce returns the primitive version of this type (i.e. itself)
func (x XNumber) Reduce() XPrimitive { return x }

// ToXText converts this type to text
func (x XNumber) ToXText() XText { return NewXText(x.Native().String()) }

// ToXBool converts this type to a bool
func (x XNumber) ToXBool() XBool { return NewXBool(!x.Equals(XNumberZero)) }

// ToXJSON is called when this type is passed to @(json(...))
func (x XNumber) ToXJSON() XText { return MustMarshalToXText(x.Native()) }

// Native returns the native value of this type
func (x XNumber) Native() decimal.Decimal { return x.native }

// Equals determines equality for this type
func (x XNumber) Equals(other XNumber) bool {
	return x.Native().Equals(other.Native())
}

// Compare compares this number to another
func (x XNumber) Compare(other XNumber) int {
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
var _ XPrimitive = XNumberZero
