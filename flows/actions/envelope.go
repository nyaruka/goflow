package actions

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func ActionFromEnvelope(envelope *utils.TypedEnvelope) (flows.Action, error) {
	switch envelope.Type {

	case MSG:
		action := MsgAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case ADD_LABEL:
		action := AddLabelAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case EMAIL:
		action := EmailAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case SET_PREFERRED_CHANNEL:
		action := PreferredChannelAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case ADD_TO_GROUP:
		action := AddToGroupAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case FLOW:
		action := FlowAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case SAVE_RESULT:
		action := SaveResultAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case SAVE_TO_CONTACT:
		action := SaveToContactAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case SET_LANGUAGE:
		action := SetLanguageAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case WEBHOOK:
		action := WebhookAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	default:
		return nil, fmt.Errorf("Unknown action type: %s", envelope.Type)
	}
}
