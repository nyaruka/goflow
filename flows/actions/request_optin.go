package actions

import (
	"context"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeRequestOptIn, func() flows.Action { return &RequestOptInAction{} })
}

// TypeRequestOptIn is the type for the send optin action
const TypeRequestOptIn string = "request_optin"

// RequestOptInAction can be used to send an optin to the contact if the channel supports that.
//
// An [event:optin_requested] event will be created if the optin was requested.
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
type RequestOptInAction struct {
	baseAction
	onlineAction

	OptIn *assets.OptInReference `json:"optin" validate:"required"`
}

// NewRequestOptIn creates a new request optin action
func NewRequestOptIn(uuid flows.ActionUUID, optIn *assets.OptInReference) *RequestOptInAction {
	return &RequestOptInAction{
		baseAction: newBaseAction(TypeRequestOptIn, uuid),
		OptIn:      optIn,
	}
}

// Execute creates the optin events
func (a *RequestOptInAction) Execute(ctx context.Context, run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	optIn := run.Session().Assets().OptIns().Get(a.OptIn.UUID)
	destinations := run.Contact().ResolveDestinations(false)

	if len(destinations) > 0 {
		ch := destinations[0].Channel
		urn := destinations[0].URN

		if ch.HasFeature(assets.ChannelFeatureOptIns) {
			logEvent(events.NewOptInRequested(optIn, ch, urn.URN()))
		}
	}

	return nil
}

func (a *RequestOptInAction) Inspect(result func(*flows.ResultInfo), dependency func(assets.Reference)) {
	dependency(a.OptIn)
}
