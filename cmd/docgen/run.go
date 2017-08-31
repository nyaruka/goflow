package main

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/utils"
)

var assetsDef = `
[
	{
		"type": "flow",
		"url": "http://testserver/assets/flow/50c3706e-fedb-42c0-8eab-dda3335714b7",
		"content": {
			"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
			"name": "ActionFlow",
			"nodes": [{
				"uuid": "72a1f5df-49f9-45df-94c9-d86f7ea064e5",
				"actions": [%s]
			}]
		}
	},ł
	{
		"type": "flow",
		"url": "http://testserver/assets/flow/b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
		"content": {
			"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
			"name": "Subflow",
			"nodes": [{
				"uuid": "d9dba561-b5ee-4f62-ba44-60c4dc242b84",
				"actions": []
			}]
		}
	]
]
`

var emptyDef = `
[
	{
		"type": "flow",
		"url": "http://testserver/assets/flow/50c3706e-fedb-42c0-8eab-dda3335714b7",
		"content": {
			"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
			"name": "EmptyFlow",
			"nodes": []
		}
	}
]
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

func createExampleSession(assetsDef string) (flows.Session, error) {
	// read our assets
	assetCache := engine.NewAssetCache()
	if err := assetCache.Include(json.RawMessage(assetsDef)); err != nil {
		return nil, err
	}

	// create our engine session
	assetURLs := map[engine.AssetItemType]string{"flow": "http://testserver/assets/flow"}
	session := engine.NewSession(assetCache, assetURLs)

	// create our contact
	contact, err := flows.ReadContact(session.Assets(), json.RawMessage(contactDef))
	if err != nil {
		return nil, err
	}

	session.SetContact(contact)

	// and start the example flow
	err = session.StartFlow(flows.FlowUUID("50c3706e-fedb-42c0-8eab-dda3335714b7"), nil)
	return session, err
}

func eventsForAction(actionJSON []byte) (json.RawMessage, error) {
	assetsDef := fmt.Sprintf(assetsDef, actionJSON)
	session, err := createExampleSession(assetsDef)
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
