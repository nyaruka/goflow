package modifiers

import (
	"github.com/nyaruka/goflow/flows"
)

var registeredTypes = map[string](func() Modifier){}

// RegisterType registers a new type of modifier
func RegisterType(name string, initFunc func() Modifier) {
	registeredTypes[name] = initFunc
}

// Modifier is something which can modify a contact
type Modifier interface {
	// Apply applies this modification to the given contact
	Apply(flows.SessionAssets, *flows.Contact) flows.Event
}

// the base of all modifier types
type baseModifier struct {
	Type_ string `json:"type" validate:"required"`
}

func newBaseModifier(typeName string) baseModifier {
	return baseModifier{Type_: typeName}
}

// Type returns the type of this modifier
func (m *baseModifier) Type() string { return m.Type_ }
