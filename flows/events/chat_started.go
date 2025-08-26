package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeChatStarted, func() flows.Event { return &ChatStarted{} })
}

// TypeChatStarted is the type of our chat started event
const TypeChatStarted string = "chat_started"

// ChatStarted events are created when a contact initiates a chat session.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "chat_started",
//	  "created_on": "2019-01-02T15:04:05Z",
//	  "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Facebook"},
//	  "params": {"referrer_id": "acme"}
//	}
//
// @event chat_started
type ChatStarted struct {
	BaseEvent

	Channel *assets.ChannelReference `json:"channel"`
	Params  map[string]string        `json:"params,omitempty"`
}

// NewChatStarted returns a new chat started event
func NewChatStarted(channel *assets.ChannelReference, params map[string]string) *ChatStarted {
	return &ChatStarted{
		BaseEvent: NewBaseEvent(TypeChatStarted),
		Channel:   channel,
		Params:    params,
	}
}

var _ flows.Event = (*ChatStarted)(nil)
