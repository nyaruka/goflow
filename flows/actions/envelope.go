package actions

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func ActionFromEnvelope(envelope *utils.TypedEnvelope) (flows.Action, error) {
	switch envelope.Type {

	case TypeAddLabel:
		action := AddLabelAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case TypeAddToGroup:
		action := AddToGroupAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case TypeEmail:
		action := EmailAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case TypeFlow:
		action := FlowAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case TypeMsg:
		action := MsgAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case TypeRemoveFromGroup:
		action := RemoveFromGroupAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case TypeReply:
		action := ReplyAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case TypeSaveResult:
		action := SaveResultAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case TypeSaveToContact:
		action := SaveToContactAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case TypeSetPreferredChannel:
		action := PreferredChannelAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	case TypeWebhook:
		action := WebhookAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, envelope.TraceError(err)

	default:
		return nil, fmt.Errorf("Unknown action type: %s", envelope.Type)
	}
}
