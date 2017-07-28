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
	switch envelope.Type {

	case TypeAddToGroup:
		event := AddToGroupEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypeSendEmail:
		event := SendEmailEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypeError:
		event := ErrorEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypeFlowEntered:
		event := FlowEnteredEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypeFlowExited:
		event := FlowExitedEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypeFlowWait:
		event := FlowWaitEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypeMsgReceived:
		event := MsgReceivedEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypeSendMsg:
		event := SendMsgEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypeMsgWait:
		event := MsgWaitEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypeRemoveFromGroup:
		event := RemoveFromGroupEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypeSaveFlowResult:
		event := SaveFlowResultEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypeSaveContactField:
		event := SaveContactFieldEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypePreferredChannel:
		event := PreferredChannelEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypeSetEnvironment:
		event := SetEnvironmentEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypeSetExtra:
		event := SetExtraEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypeSetContact:
		event := SetContactEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypeUpdateContact:
		event := UpdateContactEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	case TypeWebhookCalled:
		event := WebhookCalledEvent{}
		return &event, utils.UnmarshalAndValidate(envelope.Data, &event, "event")

	default:
		return nil, fmt.Errorf("Unknown event type: %s", envelope.Type)
	}
}
