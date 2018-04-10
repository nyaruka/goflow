package types

import (
	"time"

	"github.com/nyaruka/goflow/utils"
)

// XDate is a date
type XDate time.Time

// NewXDate creates a new date
func NewXDate(value time.Time) XDate {
	return XDate(value)
}

// Reduce returns the primitive version of this type (i.e. itself)
func (x XDate) Reduce() XPrimitive { return x }

// ToString converts this type to a string
func (x XDate) ToString() XString { return XString(utils.DateToISO(x.Native())) }

// ToBool converts this type to a bool
func (x XDate) ToBool() XBool { return XBool(!x.Native().IsZero()) }

// ToJSON converts this type to JSON
func (x XDate) ToJSON() XString { return RequireMarshalToXString(utils.DateToISO(x.Native())) }

// Native returns the native value of this type
func (x XDate) Native() time.Time { return time.Time(x) }

// Compare compares this date to another
func (x XDate) Compare(other XDate) int {
	switch {
	case x.Native().Before(other.Native()):
		return -1
	case x.Native().After(other.Native()):
		return 1
	default:
		return 0
	}
}

// MarshalJSON is called when a struct containing this type is marshaled
func (x XDate) MarshalJSON() ([]byte, error) {
	nativePtr := (time.Time)(x)
	return nativePtr.MarshalJSON()
}

// UnmarshalJSON is called when a struct containing this type is unmarshaled
func (x *XDate) UnmarshalJSON(data []byte) error {
	nativePtr := (*time.Time)(x)
	return nativePtr.UnmarshalJSON(data)
}

// XDateZero is the zero time value
var XDateZero = NewXDate(time.Time{})
var _ XPrimitive = XDateZero
