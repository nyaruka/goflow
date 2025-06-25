package triggers

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeMsg, readMsg)
}

// TypeMsg is the type for message triggered sessions
const TypeMsg string = "msg"

// Msg is used when a session was triggered by a message being received by the caller
//
//	{
//	  "type": "msg",
//	  "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//	  "event": {
//	    "type": "msg_received",
//	    "created_on": "2006-01-02T15:04:05Z",
//	    "msg": {
//	      "uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
//	      "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Twilio"},
//	      "urn": "tel:+12065551212",
//	      "text": "hi there",
//	      "attachments": ["https://s3.amazon.com/mybucket/attachment.jpg"]
//	    }
//	  },
//	  "keyword_match": {
//	    "type": "first_word",
//	    "keyword": "start"
//	  },
//	  "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//	}
//
// @trigger msg
type Msg struct {
	baseTrigger

	event *events.MsgReceived
	match *KeywordMatch
}

func (t *Msg) Event() flows.Event { return t.event }

// KeywordMatchType describes how the message matched a keyword
type KeywordMatchType string

// the different types of keyword match
const (
	KeywordMatchTypeFirstWord KeywordMatchType = "first_word"
	KeywordMatchTypeOnlyWord  KeywordMatchType = "only_word"
)

// KeywordMatch describes why the message triggered a session
type KeywordMatch struct {
	Type    KeywordMatchType `json:"type" validate:"required"`
	Keyword string           `json:"keyword" validate:"required"`
}

// NewKeywordMatch creates a new keyword match
func NewKeywordMatch(typeName KeywordMatchType, keyword string) *KeywordMatch {
	return &KeywordMatch{Type: typeName, Keyword: keyword}
}

// Initialize initializes the session
func (t *Msg) Initialize(session flows.Session) error {
	// update our input
	input := inputs.NewMsg(session, t.event.Msg, t.triggeredOn)
	session.SetInput(input)

	return t.baseTrigger.Initialize(session)
}

// Context for msg triggers additionally exposes the keyword match
func (t *Msg) Context(env envs.Environment) map[string]types.XValue {
	c := t.context()
	if t.match != nil {
		c.keyword = t.match.Keyword
	}
	return c.asMap()
}

var _ flows.Trigger = (*Msg)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// MsgBuilder is a builder for msg type triggers
type MsgBuilder struct {
	t *Msg
}

// Msg returns a msg trigger builder
func (b *Builder) Msg(e *events.MsgReceived) *MsgBuilder {
	return &MsgBuilder{
		t: &Msg{
			baseTrigger: newBaseTrigger(TypeMsg, b.flow, false, nil),
			event:       e,
		},
	}
}

// WithMatch sets the keyword match for the trigger
func (b *MsgBuilder) WithMatch(match *KeywordMatch) *MsgBuilder {
	b.t.match = match
	return b
}

// Build builds the trigger
func (b *MsgBuilder) Build() *Msg {
	return b.t
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type msgEnvelope struct {
	baseEnvelope

	Event *events.MsgReceived `json:"event"`         // TODO make required
	Msg   *flows.MsgIn        `json:"msg,omitempty"` // used by older sessions
	Match *KeywordMatch       `json:"keyword_match,omitempty" validate:"omitempty"`
}

func readMsg(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &msgEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &Msg{
		event: e.Event,
		match: e.Match,
	}

	// older triggers will have msg instead of event so convert that into an event
	if e.Msg != nil {
		t.event = &events.MsgReceived{
			BaseEvent: events.BaseEvent{Type_: events.TypeMsgReceived, CreatedOn_: e.baseEnvelope.TriggeredOn},
			Msg:       e.Msg,
		}
	}

	if err := t.unmarshal(sa, &e.baseEnvelope, missing); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *Msg) MarshalJSON() ([]byte, error) {
	e := &msgEnvelope{
		Event: t.event,
		Match: t.match,
	}

	if err := t.marshal(&e.baseEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
