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
	case TypeAddContactURN:
		action = &AddContactURNAction{}
	case TypeCallWebhook:
		action = &CallWebhookAction{}
	case TypeRemoveContactGroups:
		action = &RemoveContactGroupsAction{}
	case TypeSendBroadcast:
		action = &SendBroadcastAction{}
	case TypeSendEmail:
		action = &SendEmailAction{}
	case TypeSendMsg:
		action = &SendMsgAction{}
	case TypeSetContactChannel:
		action = &SetContactChannelAction{}
	case TypeSetContactField:
		action = &SetContactFieldAction{}
	case TypeSetContactLanguage:
		action = &SetContactLanguageAction{}
	case TypeSetContactName:
		action = &SetContactNameAction{}
	case TypeSetContactTimezone:
		action = &SetContactTimezoneAction{}
	case TypeSetRunResult:
		action = &SetRunResultAction{}
	case TypeStartFlow:
		action = &StartFlowAction{}
	case TypeStartSession:
		action = &StartSessionAction{}
	default:
		return nil, fmt.Errorf("unknown action type: %s", envelope.Type)
	}

	return action, utils.UnmarshalAndValidate(envelope.Data, action, fmt.Sprintf("action[type=%s]", envelope.Type))
}
