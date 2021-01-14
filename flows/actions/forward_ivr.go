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
//     "phone": "+12065551212"
//   }
//
// @action forward_ivr
type ForwardIVRAction struct {
	baseAction
	voiceAction

	Phone string `json:"phone" validate:"required" engine:"evaluated"`
}

// NewForwardIVR creates a new say message action
func NewForwardIVR(uuid flows.ActionUUID, phone string) *ForwardIVRAction {
	return &ForwardIVRAction{
		baseAction: newBaseAction(TypeForwardIVR, uuid),
		Phone:      phone,
	}
}

// Execute runs this action
func (a *ForwardIVRAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	phone, err := run.EvaluateTemplate(a.Phone)
	if err != nil {
		logEvent(events.NewError(err))
		return nil
	}

	urn, err := urns.NewTelURNForCountry(phone, string(run.Contact().Country()))
	if err != nil {
		logEvent(events.NewError(err))
		return nil
	}

	logEvent(events.NewIVRForwarded(urn))

	return nil
}
