package features

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
)

// FeatureMsgSends means flow sends a message or broadcast (including IVR msgs)
const FeatureMsgSends flows.Feature = "msg_sends"

func init() {
	registerType(FeatureMsgSends, checkMsgSends)
}

func checkMsgSends(flow flows.Flow) bool {
	return hasActionTypes(flow, actions.TypeSendMsg, actions.TypeSendBroadcast, actions.TypeSayMsg, actions.TypePlayAudio)
}
