package modifiers

import (
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
func (m *URNs) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventLogger) (bool, error) {
	modified := false

	// validate modifier URNs and throw away any invalid
	urnz := make([]urns.URN, 0, len(m.urnz))
	for _, urn := range m.urnz {
		urn := urn.Normalize()
		if err := urn.Validate(); err != nil {
			log(events.NewError(fmt.Sprintf("'%s' is not valid URN", urn)))
		} else {
			urnz = append(urnz, urn)
		}
	}

	switch m.modification {
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
		log(events.NewContactURNsChanged(contact.URNs().Encode()))
		return true, nil
	}
	return false, nil
}

var _ flows.Modifier = (*URNs)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type urnsEnvelope struct {
	utils.TypedEnvelope

	URNs         []urns.URN       `json:"urns" validate:"required"`
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
