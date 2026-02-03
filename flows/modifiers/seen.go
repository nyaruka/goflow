package modifiers

import (
	"context"
	"time"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeSeen, readSeen)
}

// TypeSeen is the type of our seen modifier
const TypeSeen string = "seen"

// Seen modifies the last seen of a contact
type Seen struct {
	baseModifier

	seenOn time.Time
}

// NewSeen creates a new seen modifier
func NewSeen(seenOn time.Time) *Seen {
	return &Seen{
		baseModifier: newBaseModifier(TypeSeen),
		seenOn:       seenOn,
	}
}

// Apply applies this modification to the given contact
func (m *Seen) Apply(ctx context.Context, eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventLogger) (bool, error) {
	if contact.LastSeenOn() == nil || m.seenOn.After(*contact.LastSeenOn()) {
		contact.SetLastSeenOn(m.seenOn)
		log(events.NewContactLastSeenChanged(m.seenOn))
		return true, nil
	}

	return false, nil
}

var _ flows.Modifier = (*Seen)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type seenEnvelope struct {
	utils.TypedEnvelope

	SeenOn time.Time `json:"seen_on"`
}

func readSeen(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &seenEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	return NewSeen(e.SeenOn), nil
}

func (m *Seen) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&seenEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		SeenOn:        m.seenOn,
	})
}
