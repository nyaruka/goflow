package types

import (
	"github.com/nyaruka/goflow/utils"
)

type XLazy struct {
	fetch   func() XValue
	value   XValue
	fetched bool
}

// NewXLazy creates a new XLazy
func NewXLazy(fetch func() XValue) *XLazy {
	return &XLazy{fetch: fetch}
}

// Describe returns a representation of this type for error messages
func (x *XLazy) Describe() string { return x.get().Describe() }

// Reduce returns the primitive version of this type (i.e. itself)
func (x *XLazy) Reduce(env utils.Environment) XPrimitive { return x.ToXJSON(env) }

// ToXJSON is called when this type is passed to @(json(...))
func (x *XLazy) ToXJSON(env utils.Environment) XText { return x.get().ToXJSON(env) }

func (x *XLazy) get() XValue {
	if !x.fetched {
		x.value = x.fetch()
		x.fetched = true
	}
	return x.value
}

var _ XValue = (*XLazy)(nil)
