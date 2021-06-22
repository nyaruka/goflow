package triggers

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeManual, readManualTrigger)
}

// TypeManual is the type for manually triggered sessions
const TypeManual string = "manual"

// ManualTrigger is used when a session was triggered manually by a user
//
//   {
//     "type": "manual",
//     "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//     "contact": {
//       "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
//       "name": "Bob",
//       "created_on": "2018-01-01T12:00:00.000000Z"
//     },
//     "user": {"email": "bob@nyaruka.com", "name": "Bob"},
//     "origin": "ui",
//     "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//   }
//
// @trigger manual
type ManualTrigger struct {
	baseTrigger

	user   *flows.User
	origin string
}

// Context for manual triggers always has non-nil params
func (t *ManualTrigger) Context(env envs.Environment) map[string]types.XValue {
	c := t.context()
	c.user = flows.Context(env, t.user)
	c.origin = t.origin
	return c.asMap()
}

var _ flows.Trigger = (*ManualTrigger)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// ManualBuilder is a builder for manual type triggers
type ManualBuilder struct {
	t *ManualTrigger
}

// Manual returns a manual trigger builder
func (b *Builder) Manual() *ManualBuilder {
	return &ManualBuilder{
		t: &ManualTrigger{baseTrigger: newBaseTrigger(TypeManual, b.environment, b.flow, b.contact, nil, false, nil)},
	}
}

// WithParams sets the params for the trigger
func (b *ManualBuilder) WithParams(params *types.XObject) *ManualBuilder {
	b.t.params = params
	return b
}

// WithConnection sets the channel connection for the trigger
func (b *ManualBuilder) WithConnection(channel *assets.ChannelReference, urn urns.URN) *ManualBuilder {
	b.t.connection = flows.NewConnection(channel, urn)
	return b
}

// WithUser sets the user (e.g. an email address, login) for the trigger
func (b *ManualBuilder) WithUser(user *flows.User) *ManualBuilder {
	b.t.user = user
	return b
}

// WithOrigin sets the origin (e.g. ui, api) for the trigger
func (b *ManualBuilder) WithOrigin(origin string) *ManualBuilder {
	b.t.origin = origin
	return b
}

// AsBatch sets batch mode on for the trigger
func (b *ManualBuilder) AsBatch() *ManualBuilder {
	b.t.batch = true
	return b
}

// Build builds the trigger
func (b *ManualBuilder) Build() *ManualTrigger {
	return b.t
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type manualTriggerEnvelope struct {
	baseTriggerEnvelope
	User   *assets.UserReference `json:"user,omitempty" validate:"omitempty,dive"`
	Origin string                `json:"origin,omitempty"`
}

func readManualTrigger(sa flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &manualTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	var user *flows.User
	if e.User != nil {
		user = sa.Users().Get(e.User.Email)
		if user == nil {
			missing(e.User, nil)
		}
	}

	t := &ManualTrigger{
		user:   user,
		origin: e.Origin,
	}

	if err := t.unmarshal(sa, &e.baseTriggerEnvelope, missing); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *ManualTrigger) MarshalJSON() ([]byte, error) {
	var userRef *assets.UserReference
	if t.user != nil {
		userRef = t.user.Reference()
	}

	e := &manualTriggerEnvelope{
		User:   userRef,
		Origin: t.origin,
	}

	if err := t.marshal(&e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
