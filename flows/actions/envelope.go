package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func ActionFromEnvelope(envelope *utils.TypedEnvelope) (flows.Action, error) {
	var action flows.Action

	switch envelope.Type {
	case TypeAddInputLabels:
		action = &AddInputLabelsAction{}
	case TypeAddContactGroups:
		action = &AddContactGroupsAction{}
	case TypeAddURN:
		action = &AddURNAction{}
	case TypeSendEmail:
		action = &EmailAction{}
	case TypeStartFlow:
		action = &StartFlowAction{}
	case TypeStartSession:
		action = &StartSessionAction{}
	case TypeSendMsg:
		action = &SendMsgAction{}
	case TypeRemoveFromGroup:
		action = &RemoveFromGroupAction{}
	case TypeReply:
		action = &ReplyAction{}
	case TypeSaveFlowResult:
		action = &SaveFlowResultAction{}
	case TypeSaveContactField:
		action = &SaveContactField{}
	case TypeSetPreferredChannel:
		action = &PreferredChannelAction{}
	case TypeUpdateContact:
		action = &UpdateContactAction{}
	case TypeCallWebhook:
		action = &WebhookAction{}
	default:
		return nil, fmt.Errorf("unknown action type: %s", envelope.Type)
	}

	return action, utils.UnmarshalAndValidate(envelope.Data, action, fmt.Sprintf("action[type=%s]", envelope.Type))
}
