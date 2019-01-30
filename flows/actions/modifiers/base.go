package modifiers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

type readFunc func(flows.SessionAssets, json.RawMessage) (flows.Modifier, error)

// RegisteredTypes is the registered modifier types
var RegisteredTypes = map[string]readFunc{}

// RegisterType registers a new type of modifier
func RegisterType(name string, f readFunc) {
	RegisteredTypes[name] = f
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

// helper to re-evaluate dynamic groups and log any changes to membership
func (m *baseModifier) reevaluateDynamicGroups(env utils.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	added, removed, errors := contact.ReevaluateDynamicGroups(env, assets.Groups())

	// add error event for each group we couldn't re-evaluate
	for _, err := range errors {
		log(events.NewErrorEvent(err))
	}

	// add groups changed event for the groups we were added/removed to/from
	if len(added) > 0 || len(removed) > 0 {
		log(events.NewContactGroupsChangedEvent(added, removed))
	}
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadModifier reads a modifier from the given JSON
func ReadModifier(assets flows.SessionAssets, data json.RawMessage) (flows.Modifier, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := RegisteredTypes[typeName]
	if f == nil {
		return nil, errors.Errorf("unknown type: '%s'", typeName)
	}
	return f(assets, data)
}
