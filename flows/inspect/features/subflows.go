package features

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
)

// FeatureSubflows means flow enters a child flow or triggers another flow session
const FeatureSubflows flows.Feature = "subflows"

func init() {
	registerType(FeatureSubflows, checkSubflows)
}

func checkSubflows(flow flows.Flow) bool {
	return hasActionTypes(flow, actions.TypeEnterFlow, actions.TypeStartSession)
}
