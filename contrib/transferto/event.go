package transferto

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	events.RegisterType(TypeAirtimeTransfered, func() flows.Event { return &AirtimeTransferedEvent{} })
}

// TypeAirtimeTransfered is the type of our airtime transferred event
const TypeAirtimeTransfered string = "airtime_transferred"

// AirtimeTransferedEvent events are created when airtime has been transferred to the contact
//
//   {
//     "type": "airtime_transferred",
//     "created_on": "2006-01-02T15:04:05Z"
//   }
//
// @event airtime_transferred
type AirtimeTransferedEvent struct {
	events.BaseEvent
}

// Type returns the type of this event
func (e *AirtimeTransferedEvent) Type() string { return TypeAirtimeTransfered }

// Validate validates our event is valid and has all the assets it needs
func (e *AirtimeTransferedEvent) Validate(assets flows.SessionAssets) error {
	return nil
}

// AllowedOrigin determines where this event type can originate
func (e *AirtimeTransferedEvent) AllowedOrigin() flows.EventOrigin { return flows.EventOriginEngine }

// Apply applies this event to the given run
func (e *AirtimeTransferedEvent) Apply(run flows.FlowRun) error {
	return nil
}
