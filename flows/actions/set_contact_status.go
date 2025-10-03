package actions

import (
	"context"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/modifiers"
)

func init() {
	registerType(TypeSetContactStatus, func() flows.Action { return &SetContactStatus{} })
}

// TypeSetContactStatus is the type for the set contact status action
const TypeSetContactStatus string = "set_contact_status"

// SetContactStatus can be used to update the status of the contact, e.g. to block or unblock the contact.
// A [event:contact_status_changed] event will be created with the corresponding value.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "set_contact_status",
//	  "status": "blocked"
//	}
//
// @action set_contact_status
type SetContactStatus struct {
	baseAction
	universalAction

	Status flows.ContactStatus `json:"status" validate:"contact_status"`
}

// NewSetContactStatus creates a new set status action
func NewSetContactStatus(uuid flows.ActionUUID, status flows.ContactStatus) *SetContactStatus {
	return &SetContactStatus{
		baseAction: newBaseAction(TypeSetContactStatus, uuid),
		Status:     status,
	}
}

// Execute runs this action
func (a *SetContactStatus) Execute(ctx context.Context, run flows.Run, step flows.Step, logEvent flows.EventCallback) error {
	_, err := a.applyModifier(run, modifiers.NewStatus(a.Status), logEvent)
	return err
}
