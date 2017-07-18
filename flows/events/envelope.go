package events

import (
	"encoding/json"
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
		err := json.Unmarshal(envelope.Data, &event)
		return &event, utils.ValidateAllUnlessErr(err, &event)

	case TypeSendEmail:
		event := SendEmailEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, utils.ValidateAllUnlessErr(err, &event)

	case TypeError:
		event := ErrorEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, utils.ValidateAllUnlessErr(err, &event)

	case TypeFlowEntered:
		event := FlowEnteredEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, utils.ValidateAllUnlessErr(err, &event)

	case TypeFlowExited:
		event := FlowExitedEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, utils.ValidateAllUnlessErr(err, &event)

	case TypeFlowWait:
		event := FlowWaitEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, utils.ValidateAllUnlessErr(err, &event)

	case TypeMsgReceived:
		event := MsgReceivedEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, utils.ValidateAllUnlessErr(err, &event)

	case TypeSendMsg:
		event := SendMsgEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, utils.ValidateAllUnlessErr(err, &event)

	case TypeMsgWait:
		event := MsgWaitEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, utils.ValidateAllUnlessErr(err, &event)

	case TypeRemoveFromGroup:
		event := RemoveFromGroupEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, utils.ValidateAllUnlessErr(err, &event)

	case TypeSaveFlowResult:
		event := SaveFlowResultEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, utils.ValidateAllUnlessErr(err, &event)

	case TypeSaveContactField:
		event := SaveContactFieldEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, utils.ValidateAllUnlessErr(err, &event)

	case TypeUpdateContact:
		event := UpdateContactEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, utils.ValidateAllUnlessErr(err, &event)

	case TypeWebhookCalled:
		event := WebhookCalledEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, utils.ValidateAllUnlessErr(err, &event)

	default:
		return nil, fmt.Errorf("Unknown event type: %s", envelope.Type)
	}
}
