package types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/nyaruka/goflow/utils"
)

// XValue is the base interface of all Excellent types
type XValue interface {
	ToJSON() XString
	Reduce() XPrimitive
}

// XPrimitive is the base interface of all Excellent primitive types
type XPrimitive interface {
	XValue

	ToString() XString
	ToBool() XBool
}

// XResolvable is the interface for types which can be keyed into, e.g. foo.bar
type XResolvable interface {
	Resolve(key string) XValue
}

// XLengthable is the interface for types which have a length
type XLengthable interface {
	Length() int
}

// XIndexable is the interface for types which can be indexed into, e.g. foo.0. Such objects
// also need to be lengthable so that the engine knows what is a valid index and what isn't.
type XIndexable interface {
	XLengthable

	Index(index int) XValue
}

// Compare returns the difference between two given values
func Compare(x1 XValue, x2 XValue) (int, error) {
	if utils.IsNil(x1) && utils.IsNil(x2) {
		return 0, nil
	} else if utils.IsNil(x1) || utils.IsNil(x2) {
		return 0, fmt.Errorf("can't compare non-nil and nil values: %v and %v", x1, x2)
	}

	x1 = x1.Reduce()
	x2 = x2.Reduce()

	if reflect.TypeOf(x1) != reflect.TypeOf(x2) {
		return 0, fmt.Errorf("can't compare different types of %#v and %#v", x1, x2)
	}

	// common types, do real comparisons
	switch typed := x1.(type) {
	case XError:
		return strings.Compare(typed.Error(), x2.(error).Error()), nil
	case XNumber:
		return typed.Native().Cmp(x2.(XNumber).Native()), nil
	case XBool:
		bool1 := typed.Native()
		bool2 := x2.(XBool).Native()

		switch {
		case !bool1 && bool2:
			return -1, nil
		case bool1 == bool2:
			return 0, nil
		case bool1 && !bool2:
			return 1, nil
		}
	case XDate:
		time1 := typed.Native()
		time2 := x2.(XDate).Native()

		switch {
		case time1.Before(time2):
			return -1, nil
		case time1.Equal(time2):
			return 0, nil
		case time1.After(time2):
			return 1, nil
		}
	case XString:
		return strings.Compare(typed.Native(), x2.(XString).Native()), nil
	}

	// TODO: find better fallback
	return strings.Compare(fmt.Sprintf("%v", x1), fmt.Sprintf("%v", x2)), nil
}
