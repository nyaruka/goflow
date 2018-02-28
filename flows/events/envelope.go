package events

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

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

func EventFromEnvelope(envelope *utils.TypedEnvelope) (flows.Event, error) {
	var event flows.Event

	switch envelope.Type {
	case TypeInputLabelsAdded:
		event = &InputLabelsAddedEvent{}
	case TypeContactGroupsAdded:
		event = &ContactGroupsAddedEvent{}
	case TypeContactURNAdded:
		event = &ContactURNAddedEvent{}
	case TypeEmailCreated:
		event = &EmailCreatedEvent{}
	case TypeError:
		event = &ErrorEvent{}
	case TypeFlowTriggered:
		event = &FlowTriggeredEvent{}
	case TypeSessionTriggered:
		event = &SessionTriggeredEvent{}
	case TypeRunExpired:
		event = &RunExpiredEvent{}
	case TypeMsgReceived:
		event = &MsgReceivedEvent{}
	case TypeBroadcastCreated:
		event = &BroadcastCreatedEvent{}
	case TypeMsgWait:
		event = &MsgWaitEvent{}
	case TypeNothingWait:
		event = &NothingWaitEvent{}
	case TypeContactGroupsRemoved:
		event = &ContactGroupsRemovedEvent{}
	case TypeRunResultChanged:
		event = &RunResultChangedEvent{}
	case TypeContactFieldChanged:
		event = &ContactFieldChangedEvent{}
	case TypeContactChannelChanged:
		event = &ContactChannelChangedEvent{}
	case TypeEnvironmentChanged:
		event = &EnvironmentChangedEvent{}
	case TypeContactChanged:
		event = &ContactChangedEvent{}
	case TypeContactPropertyChanged:
		event = &ContactPropertyChangedEvent{}
	case TypeWebhookCalled:
		event = &WebhookCalledEvent{}
	default:
		return nil, fmt.Errorf("Unknown event type: %s", envelope.Type)
	}

	return event, utils.UnmarshalAndValidate(envelope.Data, event, fmt.Sprintf("event[type=%s]", envelope.Type))
}
