package features

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
)

// FeatureWebhooks means flow calls a webhook or resthook
const FeatureWebhooks flows.Feature = "webhooks"

func init() {
	registerType(FeatureWebhooks, checkWebhooks)
}

func checkWebhooks(flow flows.Flow) bool {
	return hasActionTypes(flow, actions.TypeCallResthook, actions.TypeCallWebhook)
}
