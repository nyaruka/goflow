package triggers

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeChat, readChat)
}

// TypeChat is the type for sessions triggered by chat started events
const TypeChat string = "chat"

// Chat is used when a session was triggered by a chat started event
//
//	{
//	  "type": "chat",
//	  "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//	  "event": {
//	    "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	    "type": "chat_started",
//	    "created_on": "2006-01-02T15:04:05Z"
//	  },
//	  "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//	}
//
// @trigger chat
type Chat struct {
	baseTrigger

	event *events.ChatStarted
}

// Context for manual triggers always has non-nil params
func (t *Chat) Context(env envs.Environment) map[string]types.XValue {
	c := t.context()

	// expose as @trigger.params for compatibility with old channel trigger
	if t.event.Params != nil {
		params := make(map[string]types.XValue, len(t.event.Params))
		for k, v := range t.event.Params {
			params[k] = types.NewXText(v)
		}
		c.params = types.NewXObject(params)
	} else {
		c.params = types.XObjectEmpty
	}

	return c.asMap()
}

var _ flows.Trigger = (*Chat)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// ChatBuilder is a builder for chat type triggers
type ChatBuilder struct {
	t *Chat
}

// Chat returns a chat trigger builder
func (b *Builder) Chat(e *events.ChatStarted) *ChatBuilder {
	t := &Chat{
		baseTrigger: newBaseTrigger(TypeChat, b.flow, false, nil),
		event:       e,
	}

	return &ChatBuilder{t: t}
}

// Build builds the trigger
func (b *ChatBuilder) Build() *Chat {
	return b.t
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type chatEnvelope struct {
	baseEnvelope

	Event *events.ChatStarted `json:"event" validate:"required"`
}

func readChat(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &chatEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &Chat{
		event: e.Event,
	}

	if err := t.unmarshal(sa, &e.baseEnvelope, missing); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *Chat) MarshalJSON() ([]byte, error) {
	e := &chatEnvelope{
		Event: t.event,
	}

	if err := t.marshal(&e.baseEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
