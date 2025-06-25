package modifiers

import (
	"errors"
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// ErrNoModifier is the error instance returned when a modifier is read but due to missing assets can't be returned
var ErrNoModifier = errors.New("no modifier to return because of missing assets")

type readFunc func(flows.SessionAssets, []byte, assets.MissingCallback) (flows.Modifier, error)

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

// Apply applies the given modifier to the given contact and re-evaluates query based groups if necessary
func Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, c *flows.Contact, mod flows.Modifier, logEvent flows.EventCallback) bool {
	modified := mod.Apply(eng, env, sa, c, logEvent)
	if modified {
		ReevaluateGroups(env, c, logEvent)
	}
	return modified
}

// ReevaluateGroups is a helper to re-evaluate groups and log any changes to membership
func ReevaluateGroups(env envs.Environment, contact *flows.Contact, log flows.EventCallback) {
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

// Read reads a modifier from the given JSON
func Read(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := RegisteredTypes[typeName]
	if f == nil {
		return nil, fmt.Errorf("unknown type: '%s'", typeName)
	}
	return f(sa, data, missing)
}
