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
		return &action, err

	case TypeAddToGroup:
		action := AddToGroupAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, err

	case TypeEmail:
		action := EmailAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, err

	case TypeFlow:
		action := FlowAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, err

	case TypeMsg:
		action := MsgAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, err

	case TypeRemoveFromGroup:
		action := RemoveFromGroupAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, err

	case TypeReply:
		action := ReplyAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, err

	case TypeSaveResult:
		action := SaveResultAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, err

	case TypeSaveToContact:
		action := SaveToContactAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, err

	case TypeSetPreferredChannel:
		action := PreferredChannelAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, err

	case TypeWebhook:
		action := WebhookAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, err

	default:
		return nil, fmt.Errorf("Unknown action type: %s", envelope.Type)
	}
}
