package engine

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// NewFlowEnvironment creates and returns a new FlowEnvironment given the passed in environment and flow map
func NewFlowEnvironment(env utils.Environment, flowList []flows.Flow) flows.FlowEnvironment {
	flowMap := make(map[flows.FlowUUID]flows.Flow, len(flowList))
	for _, f := range flowList {
		flowMap[f.UUID()] = f
	}

	return &flowEnvironment{env, flowMap}
}

type flowEnvironment struct {
	utils.Environment
	flows map[flows.FlowUUID]flows.Flow
}

func (e *flowEnvironment) GetFlow(uuid flows.FlowUUID) (flows.Flow, error) {
	flow, exists := e.flows[uuid]
	if exists {
		return flow, nil
	}
	return nil, fmt.Errorf("Unable to find flow with UUID: %s", uuid)
}
