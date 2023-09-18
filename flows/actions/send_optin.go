package actions

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeSendOptIn, func() flows.Action { return &SendOptInAction{} })
}

// TypeSendOptIn is the type for the send optin action
const TypeSendOptIn string = "send_optin"

// SendOptInAction can be used to send an optin to the contact if the channel supports that.
//
// An [event:optin_created] event will be created if the optin was sent.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "send_optin",
//	  "optin": {
//	    "uuid": "248be71d-78e9-4d71-a6c4-9981d369e5cb",
//	    "name": "Joke Of The Day"
//	  }
//	}
//
// @action send_optin
type SendOptInAction struct {
	baseAction
	onlineAction

	OptIn *assets.OptInReference `json:"optin" validate:"required,dive"`
}

// NewSendOptIn creates a new send optin action
func NewSendOptIn(uuid flows.ActionUUID, optIn *assets.OptInReference) *SendOptInAction {
	return &SendOptInAction{
		baseAction: newBaseAction(TypeSendOptIn, uuid),
		OptIn:      optIn,
	}
}

// Execute creates the optin events
func (a *SendOptInAction) Execute(run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	optIn := run.Session().Assets().OptIns().Get(a.OptIn.UUID)
	destinations := run.Contact().ResolveDestinations(false)

	if len(destinations) > 0 {
		ch := destinations[0].Channel
		urn := destinations[0].URN

		if ch.HasFeature(assets.ChannelFeatureOptIns) {
			logEvent(events.NewOptInCreated(optIn, ch, urn.URN()))
		}
	}

	return nil
}
