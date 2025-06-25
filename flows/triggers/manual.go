package triggers

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeManual, readManual)
}

// TypeManual is the type for manually triggered sessions
const TypeManual string = "manual"

// Manual is used when a session was triggered manually by a user
//
//	{
//	  "type": "manual",
//	  "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//	  "user": {"uuid": "0c78ef47-7d56-44d8-8f57-96e0f30e8f44", "name": "Bob"},
//	  "origin": "ui",
//	  "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//	}
//
// @trigger manual
type Manual struct {
	baseTrigger

	user   *flows.User
	origin string
}

// Context for manual triggers always has non-nil params
func (t *Manual) Context(env envs.Environment) map[string]types.XValue {
	c := t.context()
	c.user = flows.Context(env, t.user)
	c.origin = t.origin
	return c.asMap()
}

var _ flows.Trigger = (*Manual)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// ManualBuilder is a builder for manual type triggers
type ManualBuilder struct {
	t *Manual
}

// Manual returns a manual trigger builder
func (b *Builder) Manual() *ManualBuilder {
	return &ManualBuilder{
		t: &Manual{baseTrigger: newBaseTrigger(TypeManual, b.flow, false, nil)},
	}
}

// WithParams sets the params for the trigger
func (b *ManualBuilder) WithParams(params *types.XObject) *ManualBuilder {
	b.t.params = params
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
func (b *ManualBuilder) Build() *Manual {
	return b.t
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type manualEnvelope struct {
	baseEnvelope

	User   *assets.UserReference `json:"user,omitempty" validate:"omitempty"`
	Origin string                `json:"origin,omitempty"`
}

func readManual(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &manualEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	var user *flows.User
	if e.User != nil {
		user = sa.Users().Get(e.User.UUID)
		if user == nil {
			missing(e.User, nil)
		}
	}

	t := &Manual{
		user:   user,
		origin: e.Origin,
	}

	if err := t.unmarshal(sa, &e.baseEnvelope, missing); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *Manual) MarshalJSON() ([]byte, error) {
	var userRef *assets.UserReference
	if t.user != nil {
		userRef = t.user.Reference()
	}

	e := &manualEnvelope{
		User:   userRef,
		Origin: t.origin,
	}

	if err := t.marshal(&e.baseEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
