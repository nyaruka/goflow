package actions

import (
	"context"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/core/events"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeRequestOptIn, func() flows.Action { return &RequestOptIn{} })
}

// TypeRequestOptIn is the type for the send optin action
const TypeRequestOptIn string = "request_optin"

// OptInUUID is the UUID of an opt-in
type OptInUUID uuids.UUID

// OptInReference is used to reference an opt-in
type OptInReference struct {
	UUID OptInUUID `json:"uuid" validate:"required,uuid"`
	Name string    `json:"name" validate:"max=64"`
}

// RequestOptIn was used to request an optin from the contact but is no longer supported.
//
// An [event:error] event will be created.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "request_optin",
//	  "optin": {
//	    "uuid": "248be71d-78e9-4d71-a6c4-9981d369e5cb",
//	    "name": "Joke Of The Day"
//	  }
//	}
//
// @action request_optin
type RequestOptIn struct {
	baseAction
	onlineAction

	OptIn *OptInReference `json:"optin" validate:"required"`
}

// NewRequestOptIn creates a new request optin action
func NewRequestOptIn(uuid flows.ActionUUID, optIn *OptInReference) *RequestOptIn {
	return &RequestOptIn{
		baseAction: newBaseAction(TypeRequestOptIn, uuid),
		OptIn:      optIn,
	}
}

// Execute logs an error event as this action is no longer supported
func (a *RequestOptIn) Execute(ctx context.Context, run flows.Run, step flows.Step, log events.EventLogger) error {
	log(events.NewActionUnsupportedError("opt-in requests are no longer supported"))

	return nil
}
