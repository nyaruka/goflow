package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeIVRForwarded, func() flows.Event { return &IVRForwardedEvent{} })
}

// TypeIVRForwarded is a constant for IVR forwarded events
const TypeIVRForwarded string = "ivr_forwarded"

// IVRForwardedEvent events are created when an action wants to forward an IVR call to another phone number.
//
//   {
//     "type": "ivr_forwarded",
//     "created_on": "2006-01-02T15:04:05Z",
//     "urn": "tel:+12065551212"
//   }
//
// @event ivr_forwarded
type IVRForwardedEvent struct {
	baseEvent

	URN urns.URN `json:"urn" validate:"required,urn"`
}

// NewIVRForwarded creates a new IVR forwarded event
func NewIVRForwarded(urn urns.URN) *IVRForwardedEvent {
	return &IVRForwardedEvent{
		baseEvent: newBaseEvent(TypeIVRForwarded),
		URN:       urn,
	}
}
