package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func ActionFromEnvelope(envelope *utils.TypedEnvelope) (flows.Action, error) {
	switch envelope.Type {

	case TypeAddLabel:
		action := AddLabelAction{}
		return &action, utils.UnmarshalAndValidate(envelope.Data, &action, "action")

	case TypeAddToGroup:
		action := AddToGroupAction{}
		return &action, utils.UnmarshalAndValidate(envelope.Data, &action, "action")

	case TypeSendEmail:
		action := EmailAction{}
		return &action, utils.UnmarshalAndValidate(envelope.Data, &action, "action")

	case TypeStartFlow:
		action := StartFlowAction{}
		return &action, utils.UnmarshalAndValidate(envelope.Data, &action, "action")

	case TypeSendMsg:
		action := SendMsgAction{}
		return &action, utils.UnmarshalAndValidate(envelope.Data, &action, "action")

	case TypeRemoveFromGroup:
		action := RemoveFromGroupAction{}
		return &action, utils.UnmarshalAndValidate(envelope.Data, &action, "action")

	case TypeReply:
		action := ReplyAction{}
		return &action, utils.UnmarshalAndValidate(envelope.Data, &action, "action")

	case TypeSaveFlowResult:
		action := SaveFlowResultAction{}
		return &action, utils.UnmarshalAndValidate(envelope.Data, &action, "action")

	case TypeSaveContactField:
		action := SaveContactField{}
		return &action, utils.UnmarshalAndValidate(envelope.Data, &action, "action")

	case TypeSetPreferredChannel:
		action := PreferredChannelAction{}
		return &action, utils.UnmarshalAndValidate(envelope.Data, &action, "action")

	case TypeUpdateContact:
		action := UpdateContactAction{}
		return &action, utils.UnmarshalAndValidate(envelope.Data, &action, "action")

	case TypeCallWebhook:
		action := WebhookAction{}
		return &action, utils.UnmarshalAndValidate(envelope.Data, &action, "action")

	default:
		return nil, fmt.Errorf("Unknown action type: %s", envelope.Type)
	}
}
