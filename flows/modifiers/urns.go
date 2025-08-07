package modifiers

import (
	"fmt"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
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

// URNs modifies the URNs on a contact
type URNs struct {
	baseModifier

	URNs         []urns.URN       `json:"urns" validate:"required"`
	Modification URNsModification `json:"modification" validate:"required,eq=append|eq=remove|eq=set"`
}

// NewURNs creates a new URNs modifier
func NewURNs(urnz []urns.URN, modification URNsModification) *URNs {
	return &URNs{
		baseModifier: newBaseModifier(TypeURNs),
		URNs:         urnz,
		Modification: modification,
	}
}

// Apply applies this modification to the given contact
func (m *URNs) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) bool {
	modified := false

	// validate modifier URNs and throw away any invalid
	urnz := make([]urns.URN, 0, len(m.URNs))
	for _, urn := range m.URNs {
		urn := urn.Normalize()
		if err := urn.Validate(); err != nil {
			log(events.NewError(fmt.Sprintf("'%s' is not valid URN", urn)))
		} else {
			urnz = append(urnz, urn)
		}
	}

	switch m.Modification {
	case URNsAppend:
		for _, urn := range urnz {
			if len(contact.URNs()) >= flows.MaxContactURNs {
				log(events.NewError(fmt.Sprintf("contact has too many URNs, limit is %d", flows.MaxContactURNs)))
				break
			} else if contact.AddURN(urn) {
				modified = true
			}
		}
	case URNsRemove:
		for _, urn := range urnz {
			if contact.RemoveURN(urn) {
				modified = true
			}
		}
	case URNsSet:
		modified = contact.SetURNs(urnz)
	}

	if modified {
		log(events.NewContactURNsChanged(contact.URNs().RawURNs()))
		return true
	}
	return false
}

var _ flows.Modifier = (*URNs)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readURNs(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	m := &URNs{}
	return m, utils.UnmarshalAndValidate(data, m)
}
