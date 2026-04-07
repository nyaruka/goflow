package modifiers

import (
	"context"
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeRoutes, readRoutes)
}

// TypeRoutes is the type of our routes modifier
const TypeRoutes string = "routes"

// RoutesModification is the type of modification to make
type RoutesModification string

// the supported types of modification
const (
	RoutesAppend RoutesModification = "append"
	RoutesRemove RoutesModification = "remove"
	RoutesSet    RoutesModification = "set"
)

// Routes modifies the URNs on a contact while preserving channel affinity
type Routes struct {
	baseModifier

	routes       []flows.Route
	modification RoutesModification
}

// NewRoutes creates a new routes modifier
func NewRoutes(routes []flows.Route, modification RoutesModification) *Routes {
	return &Routes{
		baseModifier: newBaseModifier(TypeRoutes),
		routes:       routes,
		modification: modification,
	}
}

// Apply applies this modification to the given contact
func (m *Routes) Apply(ctx context.Context, eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventLogger) (bool, error) {
	modified := false

	valid := make([]flows.Route, 0, len(m.routes))
	for _, r := range m.routes {
		urn := r.URN.Normalize()

		// throw away invalid URNs
		if err := urn.Validate(); err != nil {
			log(events.NewError(fmt.Sprintf("'%s' is not valid URN", urn), ""))
			continue
		}

		// if adding or setting, try to claim the URN
		if (m.modification == RoutesAppend || m.modification == RoutesSet) && !contact.HasURN(urn) {
			claimed, err := eng.Options().ClaimURN(ctx, sa, contact, urn)
			if err != nil {
				return false, fmt.Errorf("error claiming URN %s: %w", urn, err)
			}
			if !claimed {
				log(events.NewError("URN is taken by another contact", events.ErrorCodeURNTaken, "urn", string(urn)))
				continue
			}
		}

		valid = append(valid, flows.Route{URN: urn, Channel: r.Channel})
	}

	switch m.modification {
	case RoutesAppend:
		for _, r := range valid {
			// only count budget for new URNs - updating an existing URN's channel doesn't grow the list
			if !contact.HasURN(r.URN) && len(contact.URNs()) >= flows.MaxContactURNs {
				log(events.NewError(fmt.Sprintf("Contact has too many URNs, limit is %d", flows.MaxContactURNs), ""))
				break
			}
			if contact.AddURN(r.URN, r.Channel) {
				modified = true
			}
		}
	case RoutesRemove:
		// remove is by URN identity, channel is ignored
		for _, r := range valid {
			if contact.RemoveURN(r.URN) {
				modified = true
			}
		}
	case RoutesSet:
		urnz := make([]urns.URN, len(valid))
		channels := make([]*flows.Channel, len(valid))
		for i, r := range valid {
			urnz[i] = r.URN
			channels[i] = r.Channel
		}
		modified = contact.SetURNsWithChannels(urnz, channels)
	}

	if modified {
		log(events.NewContactURNsChanged(contact.URNs().Encode()))
		return true, nil
	}
	return false, nil
}

var _ flows.Modifier = (*Routes)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type routeEnvelope struct {
	URN     urns.URN                 `json:"urn"     validate:"required"`
	Channel *assets.ChannelReference `json:"channel" validate:"required"`
}

type routesEnvelope struct {
	utils.TypedEnvelope

	Routes       []routeEnvelope    `json:"routes"       validate:"required,dive"`
	Modification RoutesModification `json:"modification" validate:"required,eq=append|eq=remove|eq=set"`
}

func readRoutes(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &routesEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	routes := make([]flows.Route, 0, len(e.Routes))
	for _, re := range e.Routes {
		channel := sa.Channels().Get(re.Channel.UUID)
		if channel == nil {
			missing(re.Channel, nil)
			continue
		}
		routes = append(routes, flows.Route{URN: re.URN, Channel: channel})
	}

	// if we had routes in the envelope but all their channels are missing, nothing to modify
	if len(e.Routes) > 0 && len(routes) == 0 {
		return nil, ErrNoModifier
	}

	return NewRoutes(routes, e.Modification), nil
}

func (m *Routes) MarshalJSON() ([]byte, error) {
	re := make([]routeEnvelope, len(m.routes))
	for i, r := range m.routes {
		re[i] = routeEnvelope{URN: r.URN, Channel: r.Channel.Reference()}
	}
	return jsonx.Marshal(&routesEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Routes:        re,
		Modification:  m.modification,
	})
}
