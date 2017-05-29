package events

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func EventFromEnvelope(envelope *utils.TypedEnvelope) (flows.Event, error) {
	switch envelope.Type {

	case TypeAddToGroup:
		event := AddToGroupEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, err

	case TypeEmail:
		event := EmailEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, err

	case TypeError:
		event := ErrorEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, err

	case TypeFlowEnter:
		event := FlowEnterEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, err

	case TypeFlowExit:
		event := FlowExitEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, err

	case TypeFlowWait:
		event := FlowWaitEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, err

	case TypeMsgIn:
		event := MsgInEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, err

	case TypeMsgOut:
		event := MsgOutEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, err

	case TypeMsgWait:
		event := MsgWaitEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, err

	case TypeRemoveFromGroup:
		event := RemoveFromGroupEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, err

	case TypeSaveResult:
		event := SaveResultEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, err

	case TypeSaveToContact:
		event := SaveToContactEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, err

	case TypeWebhook:
		event := WebhookEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, err

	default:
		return nil, fmt.Errorf("Unknown event type: %s", envelope.Type)
	}
}
