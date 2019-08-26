package hints

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

var registeredTypes = map[string](func() flows.Hint){}

// RegisterType registers a new type of wait
func registerType(name string, initFunc func() flows.Hint) {
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

// ReadHint reads a hint from the given JSON
func ReadHint(data json.RawMessage) (flows.Hint, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, errors.Errorf("unknown type: '%s'", typeName)
	}

	hint := f()
	return hint, utils.UnmarshalAndValidate(data, hint)
}
