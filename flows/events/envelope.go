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
		return &event, envelope.TraceError(err)

	case TypeEmail:
		event := EmailEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case TypeError:
		event := ErrorEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case TypeFlowEnter:
		event := FlowEnterEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case TypeFlowExit:
		event := FlowExitEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case TypeFlowWait:
		event := FlowWaitEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case TypeMsgIn:
		event := MsgInEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case TypeMsgOut:
		event := MsgOutEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case TypeMsgWait:
		event := MsgWaitEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case TypeRemoveFromGroup:
		event := RemoveFromGroupEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case TypeSaveResult:
		event := SaveResultEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case TypeSaveToContact:
		event := SaveToContactEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case TypeWebhook:
		event := WebhookEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	default:
		return nil, fmt.Errorf("Unknown event type: %s", envelope.Type)
	}
}
