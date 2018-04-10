package main

import (
	"encoding/json"
	"fmt"
	"github.com/nyaruka/goflow/flows/events"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/utils"
)

var sessionAssets = `[
    {
        "type": "channel",
        "url": "http://testserver/assets/channel",
        "content": [
            {
                "uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d",
                "name": "Android Channel",
                "address": "+12345671111",
                "schemes": ["tel"],
                "roles": ["send", "receive"]
            },
            {
                "uuid": "8e21f093-99aa-413b-b55b-758b54308fcb",
                "name": "Twitter Channel",
                "address": "nyaruka",
                "schemes": ["twitter"],
                "roles": ["send", "receive"]
            },
            {
                "uuid": "4bb288a0-7fca-4da1-abe8-59a593aff648",
                "name": "Facebook Channel",
                "address": "235326346322111",
                "schemes": ["facebook"],
                "roles": ["send", "receive"]
            }
        ],
        "is_set": true
    },
    {
        "type": "flow",
        "url": "http://testserver/assets/flow/50c3706e-fedb-42c0-8eab-dda3335714b7",
        "content": {
            "uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
            "name": "Action Flow",
            "nodes": [{
                "uuid": "72a1f5df-49f9-45df-94c9-d86f7ea064e5",
                "actions": [
                    {
                        "uuid": "5508e6a7-26ce-4b3b-b32e-bb4e2e614f5d",
                        "type": "set_run_result",
                        "name": "Phone Number",
                        "value": "+12344563452"
                    },
                    {
                        "uuid": "06153fbd-3e2c-413a-b0df-ed15d631835a",
                        "type": "call_webhook",
                        "method": "GET",
                        "url": "http://localhost:49999/?cmd=success"
                    }
                ]
            }]
        }
    },
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
    },
    {
        "type": "field",
        "url": "http://testserver/assets/field",
        "content": [
            {"key": "gender", "label": "Gender", "value_type": "text"},
            {"key": "activation_token", "label": "Activation Token", "value_type": "text"}
        ],
        "is_set": true
    },
    {
        "type": "group",
        "url": "http://testserver/assets/group",
        "content": [
            {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Testers"},
            {"uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a", "name": "Customers"}
        ],
        "is_set": true
    },
    {
        "type": "label",
        "url": "http://testserver/assets/label",
        "content": [
            {
                "uuid": "3f65d88a-95dc-4140-9451-943e94e06fea",
                "name": "Spam"
            }
        ],
        "is_set": true
    },
    {
        "type": "location_hierarchy",
        "url": "http://testserver/assets/location_hierarchy",
        "content": {
            "id": "2342",
            "name": "Rwanda",
            "aliases": ["Ruanda"],		
            "children": [
                {
                    "id": "234521",
                    "name": "Kigali City",
                    "aliases": ["Kigali", "Kigari"],
                    "children": [
                        {
                            "id": "57735322",
                            "name": "Gasabo",
                            "children": [
                                {
                                    "id": "575743222",
                                    "name": "Gisozi"
                                },
                                {
                                    "id": "457378732",
                                    "name": "Ndera"
                                }
                            ]
                        },
                        {
                            "id": "46547322",
                            "name": "Nyarugenge",
                            "children": []
                        }
                    ]
                }
            ]
        }
    }
]`

var sessionTrigger = `{
    "type": "manual",
    "triggered_on": "2017-12-31T11:31:15.035757258-02:00",
    "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Action Flow"},
    "contact": {
        "uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f",
        "name": "Ryan Lewis",
        "urns": ["tel:+12065551212", "mailto:foo@bar.com"],
        "groups": [
            {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Testers"}
        ],
        "fields": {
            "gender": {
                "text": "Male"
            },
            "activation_token": {
                "text": "AACC55"
            }
        }
    }
}`

var initialEvents = `[
    {
        "created_on": "2000-01-01T00:00:00.000000000-00:00",
        "msg": {
            "attachments": [
                "image/jpeg:http://s3.amazon.com/bucket/test.jpg",
                "audio/mp3:http://s3.amazon.com/bucket/test.mp3"
            ],
            "channel": {
                "name": "Nexmo",
                "uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d"
            },
            "text": "Hi there",
            "urn": "tel:+12065551212",
            "uuid": "9bf91c2b-ce58-4cef-aacc-281e03f69ab5"
        },
        "type": "msg_received"
    }
]`

func createExampleSession(actionToAdd flows.Action) (flows.Session, error) {
	// read our assets
	assetCache := engine.NewAssetCache(100, 5, "testing/1.0")
	if err := assetCache.Include(json.RawMessage(sessionAssets)); err != nil {
		return nil, err
	}

	// create our engine session
	session := engine.NewSession(assetCache, engine.NewMockAssetServer())

	// optional modify the main flow by adding the provided action
	if actionToAdd != nil {
		flow, _ := session.Assets().GetFlow(flows.FlowUUID("50c3706e-fedb-42c0-8eab-dda3335714b7"))
		flow.Nodes()[0].AddAction(actionToAdd)
	}

	// read our trigger
	triggerEnvelope := &utils.TypedEnvelope{}
	if err := triggerEnvelope.UnmarshalJSON(json.RawMessage(sessionTrigger)); err != nil {
		return nil, err
	}
	trigger, err := triggers.ReadTrigger(session, triggerEnvelope)
	if err != nil {
		return nil, fmt.Errorf("error reading trigger: %s", err)
	}

	// and the initial events
	eventEnvelopes := []*utils.TypedEnvelope{}
	if err := json.Unmarshal(json.RawMessage(initialEvents), &eventEnvelopes); err != nil {
		return nil, err
	}
	events, err := events.ReadEvents(eventEnvelopes)
	if err != nil {
		return nil, err
	}

	// and start the example flow
	err = session.Start(trigger, events)
	return session, err
}
