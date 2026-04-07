package modifiers

import (
	"context"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeURNs, readURNs)
}

// TypeURNs is the type of our URNs modifier
const TypeURNs string = "urns"

// URNsModification is the type of modification to make
type URNsModification string

// the supported types of modification
const (
	URNsAppend URNsModification = "append"
	URNsRemove URNsModification = "remove"
	URNsSet    URNsModification = "set"
)

// URNs modifies the URNs on a contact.
//
// Deprecated: use the routes modifier instead, which takes (URN, channel) pairs and can set channel affinity in
// the same operation. This modifier is retained for backwards compatibility and delegates to the routes modifier
// internally.
type URNs struct {
	baseModifier

	urnz         []urns.URN
	modification URNsModification
}

// NewURNs creates a new URNs modifier
func NewURNs(urnz []urns.URN, modification URNsModification) *URNs {
	return &URNs{
		baseModifier: newBaseModifier(TypeURNs),
		urnz:         urnz,
		modification: modification,
	}
}

// Apply applies this modification to the given contact
func (m *URNs) Apply(ctx context.Context, eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventLogger) (bool, error) {
	// translate to the equivalent routes modifier - for set, preserve any existing channel affinity for URNs
	// that are still on the contact
	routes := make([]flows.Route, len(m.urnz))
	for i, u := range m.urnz {
		var ch *flows.Channel
		if m.modification == URNsSet {
			for _, existing := range contact.URNs() {
				if existing.Identity() == u.Identity() {
					ch = existing.Channel
					break
				}
			}
		}
		routes[i] = flows.Route{URN: u, Channel: ch}
	}

	var routesMod RoutesModification
	switch m.modification {
	case URNsAppend:
		routesMod = RoutesAppend
	case URNsRemove:
		routesMod = RoutesRemove
	case URNsSet:
		routesMod = RoutesSet
	}

	return NewRoutes(routes, routesMod).Apply(ctx, eng, env, sa, contact, log)
}

var _ flows.Modifier = (*URNs)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type urnsEnvelope struct {
	utils.TypedEnvelope

	URNs         []urns.URN       `json:"urns"         validate:"required"`
	Modification URNsModification `json:"modification" validate:"required,eq=append|eq=remove|eq=set"`
}

func readURNs(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &urnsEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	return NewURNs(e.URNs, e.Modification), nil
}

func (m *URNs) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&urnsEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		URNs:          m.urnz,
		Modification:  m.modification,
	})
}
