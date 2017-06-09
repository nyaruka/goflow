package engine

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// NewFlowEnvironment creates and returns a new FlowEnvironment given the passed in environment and flow map
func NewFlowEnvironment(env utils.Environment, flowList []flows.Flow, runList []flows.FlowRun, contactList []*flows.Contact) flows.FlowEnvironment {
	flowMap := make(map[flows.FlowUUID]flows.Flow, len(flowList))
	for _, f := range flowList {
		flowMap[f.UUID()] = f
	}

	runMap := make(map[flows.RunUUID]flows.FlowRun, len(runList))
	for _, r := range runList {
		runMap[r.UUID()] = r
	}

	contactMap := make(map[flows.ContactUUID]*flows.Contact, len(contactList))
	for _, c := range contactList {
		contactMap[c.UUID()] = c
	}

	return &flowEnvironment{env, flowMap, runMap, contactMap}
}

type flowEnvironment struct {
	utils.Environment
	flows    map[flows.FlowUUID]flows.Flow
	runs     map[flows.RunUUID]flows.FlowRun
	contacts map[flows.ContactUUID]*flows.Contact
}

func (e *flowEnvironment) GetFlow(uuid flows.FlowUUID) (flows.Flow, error) {
	flow, exists := e.flows[uuid]
	if exists {
		return flow, nil
	}
	return nil, fmt.Errorf("unable to find flow with UUID: %s", uuid)
}

func (e *flowEnvironment) GetRun(uuid flows.RunUUID) (flows.FlowRun, error) {
	run, exists := e.runs[uuid]
	if exists {
		return run, nil
	}
	return nil, fmt.Errorf("unable to find run with UUID: %s", uuid)
}

func (e *flowEnvironment) GetContact(uuid flows.ContactUUID) (*flows.Contact, error) {
	contact, exists := e.contacts[uuid]
	if exists {
		return contact, nil
	}
	return nil, fmt.Errorf("unable to find contact with UUID: %s", uuid)
}
