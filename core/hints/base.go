package hints

import (
	"fmt"

	"github.com/nyaruka/goflow/utils"
)

// Hint tells the caller what type of input the flow is expecting
type Hint interface {
	utils.Typed
}

var registeredTypes = map[string](func() Hint){}

// RegisterType registers a new type of wait
func registerType(name string, initFunc func() Hint) {
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

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// Read reads a hint from the given JSON
func Read(data []byte) (Hint, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, fmt.Errorf("unknown type: '%s'", typeName)
	}

	hint := f()
	return hint, utils.UnmarshalAndValidate(data, hint)
}
