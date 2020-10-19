package actions

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeForwardIVR, func() flows.Action { return &ForwardIVRAction{} })
}

// TypeForwardIVR is the type for the forward IVR action
const TypeForwardIVR string = "forward_ivr"

// ForwardIVRAction can be used to forward an IVR call to another number, perhaps a human agent. It will generate
// an [event:ivr_forwarded] event if successful.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "forward_ivr",
//     "urn": "tel:+12065551212"
//   }
//
// @action forward_ivr
type ForwardIVRAction struct {
	baseAction
	voiceAction

	URN urns.URN `json:"urn" validate:"required,urn"`
}

// NewForwardIVR creates a new say message action
func NewForwardIVR(uuid flows.ActionUUID, urn urns.URN) *ForwardIVRAction {
	return &ForwardIVRAction{
		baseAction: newBaseAction(TypeForwardIVR, uuid),
		URN:        urn,
	}
}

// Execute runs this action
func (a *ForwardIVRAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	logEvent(events.NewIVRForwarded(a.URN))

	return nil
}
