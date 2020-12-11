package features

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/routers/waits"
)

// FeatureMsgWaits means flow has a router with a wait for input
const FeatureMsgWaits flows.Feature = "msg_waits"

func init() {
	registerType(FeatureMsgWaits, checkMsgWaits)
}

func checkMsgWaits(flow flows.Flow) bool {
	for _, n := range flow.Nodes() {
		if n.Router() != nil && n.Router().Wait() != nil && n.Router().Wait().Type() == waits.TypeMsg {
			return true
		}
	}
	return false
}
