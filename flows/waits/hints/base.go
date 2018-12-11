package hints

import (
	"github.com/nyaruka/goflow/flows"
)

var registeredTypes = map[string](func() flows.Hint){}

// RegisterType registers a new type of wait
func RegisterType(name string, initFunc func() flows.Hint) {
	registeredTypes[name] = initFunc
}

// the base of all hint types
type baseHint struct {
	Type_ string `json:"type" validate:"required"`
}

func newBaseHint(typeName string) baseHint {
	return baseHint{Type_: typeName}
}

// Type returns the type of this hint
func (h *baseHint) Type() string { return h.Type_ }
