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
		return &action, utils.ValidateUnlessErr(err, &action)

	case TypeAddToGroup:
		action := AddToGroupAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, utils.ValidateUnlessErr(err, &action)

	case TypeSendEmail:
		action := EmailAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, utils.ValidateUnlessErr(err, &action)

	case TypeStartFlow:
		action := StartFlowAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, utils.ValidateUnlessErr(err, &action)

	case TypeSendMsg:
		action := SendMsgAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, utils.ValidateUnlessErr(err, &action)

	case TypeRemoveFromGroup:
		action := RemoveFromGroupAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, utils.ValidateUnlessErr(err, &action)

	case TypeReply:
		action := ReplyAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, utils.ValidateUnlessErr(err, &action)

	case TypeSaveFlowResult:
		action := SaveFlowResultAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, utils.ValidateUnlessErr(err, &action)

	case TypeSaveContactField:
		action := SaveContactField{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, utils.ValidateUnlessErr(err, &action)

	case TypeSetPreferredChannel:
		action := PreferredChannelAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, utils.ValidateUnlessErr(err, &action)

	case TypeUpdateContact:
		action := UpdateContactAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, utils.ValidateUnlessErr(err, &action)

	case TypeCallWebhook:
		action := WebhookAction{}
		err := json.Unmarshal(envelope.Data, &action)
		return &action, utils.ValidateUnlessErr(err, &action)

	default:
		return nil, fmt.Errorf("Unknown action type: %s", envelope.Type)
	}
}
