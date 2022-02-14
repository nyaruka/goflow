package modifiers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

// ErrNoModifier is the error instance returned when a modifier is read but due to missing assets can't be returned
var ErrNoModifier = errors.New("no modifier to return because of missing assets")

type readFunc func(flows.SessionAssets, json.RawMessage, assets.MissingCallback) (flows.Modifier, error)

// RegisteredTypes is the registered modifier types
var RegisteredTypes = map[string]readFunc{}

// egisters a new type of modifier
func registerType(name string, f readFunc) {
	RegisteredTypes[name] = f
}

// base of all modifier types
type baseModifier struct {
	Type_ string `json:"type" validate:"required"`
}

// creates new base modifier
func newBaseModifier(typeName string) baseModifier {
	return baseModifier{Type_: typeName}
}

// Type returns the type of this modifier
func (m *baseModifier) Type() string { return m.Type_ }

// ReevaluateGroups is a helper to re-evaluate groups and log any changes to membership
func ReevaluateGroups(env envs.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	added, removed := contact.ReevaluateQueryBasedGroups(env)

	// make sure from all static groups are removed for non-active contacts
	if contact.Status() != flows.ContactStatusActive {
		for _, g := range contact.Groups().All() {
			if !g.UsesQuery() {
				removed = append(removed, g)
			}
		}
		contact.Groups().Clear()
	}

	// add groups changed event for the groups we were added/removed to/from
	if len(added) > 0 || len(removed) > 0 {
		log(events.NewContactGroupsChanged(added, removed))
	}
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadModifier reads a modifier from the given JSON
func ReadModifier(assets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Modifier, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := RegisteredTypes[typeName]
	if f == nil {
		return nil, errors.Errorf("unknown type: '%s'", typeName)
	}
	return f(assets, data, missing)
}
