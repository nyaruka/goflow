package events

import "github.com/nyaruka/goflow/flows"

// TypeAddURN is the type of our add URN event
const TypeAddURN string = "add_urn"

// AddURNEvent events will be created with the URN tht should be added to the current contact.
//
// ```
//   {
//     "type": "add_urn",
//     "created_on": "2006-01-02T15:04:05Z",
//     "urn": "tel:+12345678900"
//   }
// ```
//
// @event add_urn
type AddURNEvent struct {
	BaseEvent
	URN flows.URN `json:"urn" validate:"required"`
}

// NewAddURNEvent returns a new add URN event
func NewAddURNEvent(urn flows.URN) *AddURNEvent {
	return &AddURNEvent{URN: urn}
}

// Type returns the type of this event
func (e *AddURNEvent) Type() string { return TypeAddURN }

// Apply applies this event to the given run
func (e *AddURNEvent) Apply(run flows.FlowRun) error {
	run.Contact().AddURN(e.URN)
	return nil
}
