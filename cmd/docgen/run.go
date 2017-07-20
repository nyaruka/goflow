package main

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/utils"
)

var flowDef = `
{
	"name": "ActionFlow",
	"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
	"nodes": [{
		"uuid": "72a1f5df-49f9-45df-94c9-d86f7ea064e5",
		"actions": [%s]
	}]
}
`

var emptyDef = `
{
	"name": "EmptyFlow",
	"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
	"nodes": []
}
`

var subflowDef = `
{
	"name": "Subflow",
	"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
	"nodes": [{
		"uuid": "d9dba561-b5ee-4f62-ba44-60c4dc242b84",
		"actions": []
	}]
}
`

var contactDef = `
{
	"name": "Ryan Lewis",
	"uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f",
	"urns": ["tel:%2B12065551212", "email:foo@bar.com"],
	"groups": [{
		"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
		"name": "Registered Users"
	}],
	"fields": {
		"activation_token": {
			"field_uuid": "ee46f9c4-b094-4e1b-ab0d-d4e65b4a99f1",
			"field_name": "Activation Token",
			"value": "XFW-JEV-9QE",
			"created_on": "2017-05-24T11:31:15.035757258-05:00"
		}
	}
}
`

func createExampleSession(flowDef string) (flows.Session, error) {
	// read our flow
	flow, err := definition.ReadFlow(json.RawMessage(flowDef))
	if err != nil {
		return nil, err
	}

	// and our subflow
	subflow, err := definition.ReadFlow(json.RawMessage(subflowDef))
	if err != nil {
		return nil, err
	}

	// create our contact
	contact, err := flows.ReadContact(json.RawMessage(contactDef))
	if err != nil {
		return nil, err
	}

	// start our flow
	env := engine.NewFlowEnvironment(utils.NewDefaultEnvironment(), []flows.Flow{subflow}, []flows.FlowRun{}, []*flows.Contact{})
	return engine.StartFlow(env, flow, contact, nil, nil, nil)
}

func eventsForAction(actionJSON []byte) (json.RawMessage, error) {
	flowDef := fmt.Sprintf(flowDef, actionJSON)
	session, err := createExampleSession(flowDef)
	if err != nil {
		return nil, err
	}

	eventLog := session.Log()
	eventJSON := make([]json.RawMessage, len(eventLog))
	for i, logEntry := range eventLog {
		typed, err := utils.EnvelopeFromTyped(logEntry.Event())
		if err != nil {
			return nil, err
		}
		eventJSON[i], err = json.MarshalIndent(typed, "", "    ")
		if err != nil {
			return nil, err
		}
	}
	if len(eventLog) == 1 {
		return eventJSON[0], err
	}
	js, err := json.MarshalIndent(eventJSON, "", "    ")
	if err != nil {
		return nil, err
	}
	return js, nil
}
