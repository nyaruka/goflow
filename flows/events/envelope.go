package events

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func EventFromEnvelope(envelope *utils.TypedEnvelope) (flows.Event, error) {
	switch envelope.Type {

	case ADD_TO_GROUP:
		event := AddToGroupEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case ERROR:
		event := ErrorEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case FLOW_ENTER:
		event := FlowEnterEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case FLOW_EXIT:
		event := FlowExitEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case FLOW_WAIT:
		event := FlowWaitEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case MSG_IN:
		event := MsgInEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case MSG_OUT:
		event := MsgOutEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case MSG_WAIT:
		event := MsgWaitEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case REMOVE_FROM_GROUP:
		event := RemoveFromGroupEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case SAVE_RESULT:
		event := SaveResultEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case SAVE_TO_CONTACT:
		event := SaveToContactEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case SET_LANGUAGE:
		event := SetLanguageEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	case WEBHOOK:
		event := WebhookEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, envelope.TraceError(err)

	default:
		return nil, fmt.Errorf("Unknown event type: %s", envelope.Type)
	}
}
