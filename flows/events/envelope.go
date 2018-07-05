package events

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// ReadEvents reads the events from the given envelopes
func ReadEvents(envelopes []*utils.TypedEnvelope) ([]flows.Event, error) {
	events := make([]flows.Event, len(envelopes))
	for e, envelope := range envelopes {
		event, err := EventFromEnvelope(envelope)
		if err != nil {
			return nil, err
		}
		event.SetFromCaller(true)
		events[e] = event
	}
	return events, nil
}

// EventFromEnvelope reads a single event from the given envelope
func EventFromEnvelope(envelope *utils.TypedEnvelope) (flows.Event, error) {
	var event flows.Event

	switch envelope.Type {
	case TypeBroadcastCreated:
		event = &BroadcastCreatedEvent{}
	case TypeContactChanged:
		event = &ContactChangedEvent{}
	case TypeContactChannelChanged:
		event = &ContactChannelChangedEvent{}
	case TypeContactFieldChanged:
		event = &ContactFieldChangedEvent{}
	case TypeContactGroupsAdded:
		event = &ContactGroupsAddedEvent{}
	case TypeContactGroupsRemoved:
		event = &ContactGroupsRemovedEvent{}
	case TypeContactLanguageChanged:
		event = &ContactLanguageChangedEvent{}
	case TypeContactNameChanged:
		event = &ContactNameChangedEvent{}
	case TypeContactTimezoneChanged:
		event = &ContactTimezoneChangedEvent{}
	case TypeContactURNAdded:
		event = &ContactURNAddedEvent{}
	case TypeEmailCreated:
		event = &EmailCreatedEvent{}
	case TypeEnvironmentChanged:
		event = &EnvironmentChangedEvent{}
	case TypeError:
		event = &ErrorEvent{}
	case TypeFlowTriggered:
		event = &FlowTriggeredEvent{}
	case TypeInputLabelsAdded:
		event = &InputLabelsAddedEvent{}
	case TypeMsgCreated:
		event = &MsgCreatedEvent{}
	case TypeMsgReceived:
		event = &MsgReceivedEvent{}
	case TypeMsgWait:
		event = &MsgWaitEvent{}
	case TypeNothingWait:
		event = &NothingWaitEvent{}
	case TypeResthookCalled:
		event = &ResthookCalledEvent{}
	case TypeRunExpired:
		event = &RunExpiredEvent{}
	case TypeRunResultChanged:
		event = &RunResultChangedEvent{}
	case TypeSessionTriggered:
		event = &SessionTriggeredEvent{}
	case TypeWaitTimedOut:
		event = &WaitTimedOutEvent{}
	case TypeWebhookCalled:
		event = &WebhookCalledEvent{}
	default:
		return nil, fmt.Errorf("unknown event type: %s", envelope.Type)
	}

	return event, utils.UnmarshalAndValidate(envelope.Data, event, fmt.Sprintf("event[type=%s]", envelope.Type))
}

// EventsToEnvelopes converts the given events to typed envelopes
func EventsToEnvelopes(events []flows.Event) ([]*utils.TypedEnvelope, error) {
	var err error
	envelopes := make([]*utils.TypedEnvelope, len(events))
	for e, event := range events {
		if envelopes[e], err = utils.EnvelopeFromTyped(event); err != nil {
			return nil, err
		}
	}
	return envelopes, nil
}
