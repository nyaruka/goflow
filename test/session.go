package test

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
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
                "name": "My Android Phone",
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
        ]
    },
    {
        "type": "flow",
        "url": "http://testserver/assets/flow/50c3706e-fedb-42c0-8eab-dda3335714b7",
        "content": {
            "uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
            "name": "Registration",
            "language": "eng",
            "type": "messaging",
            "revision": 123,
            "nodes": [
                {
                    "uuid": "72a1f5df-49f9-45df-94c9-d86f7ea064e5",
                    "actions": [
                        {
                            "uuid": "9487a60e-a6ef-4a88-b35d-894bfe074144",
                            "type": "start_flow",
                            "flow": {
                                "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
                                "name": "Collect Age"
                            }
                        }
                    ],
                    "exits": [
                        {
                            "uuid": "37d8813f-1402-4ad2-9cc2-e9054a96525b",
                            "destination_node_uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03"
                        }
                    ]
                },
                {
                    "uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
                    "actions": [
                        {
                            "uuid": "5508e6a7-26ce-4b3b-b32e-bb4e2e614f5d",
                            "type": "set_run_result",
                            "name": "Phone Number",
                            "value": "+12344563452"
                        },
                        {
                            "uuid": "72fea511-246f-49ad-846d-853b22ecc9c9",
                            "type": "set_run_result",
                            "name": "Favorite Color",
                            "value": "red",
                            "category": "Red"
                        },
                        {
                            "uuid": "06153fbd-3e2c-413a-b0df-ed15d631835a",
                            "type": "call_webhook",
                            "method": "GET",
                            "url": "http://localhost/?cmd=echo&content=%7B%22results%22%3A%5B%7B%22state%22%3A%22WA%22%7D%2C%7B%22state%22%3A%22IN%22%7D%5D%7D"
                        }
                    ],
                    "exits": [
                        {
                            "uuid": "d898f9a4-f0fc-4ac4-a639-c98c602bb511",
                            "destination_node_uuid": "c0781400-737f-4940-9a6c-1ec1c3df0325"
                        }
                    ]
                },
                {
                    "uuid": "c0781400-737f-4940-9a6c-1ec1c3df0325",
                    "actions": []
                }
            ]
        }
    },
    {
        "type": "flow",
        "url": "http://testserver/assets/flow/b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
        "content": {
            "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
            "name": "Collect Age",
            "language": "eng",
            "type": "messaging",
            "nodes": [{
                "uuid": "d9dba561-b5ee-4f62-ba44-60c4dc242b84",
                "actions": [
                    {
                        "uuid": "4ed673b3-bdcc-40f2-944b-6ad1c82eb3ee",
                        "type": "set_run_result",
                        "name": "Age",
                        "value": "23",
                        "category": "Youth"
                    },
                    {
                        "uuid": "7a0c3cec-ef84-41aa-bf2b-be8259038683",
                        "type": "set_contact_field",
                        "field": {
                            "key": "age",
                            "name": "Age"
                        },
                        "value": "@run.results.age"
                    }
                ]
            }]
        }
    },
    {
        "type": "flow",
        "url": "http://testserver/assets/flow/fece6eac-9127-4343-9269-56e88f391562",
        "content": {
            "uuid": "fece6eac-9127-4343-9269-56e88f391562",
            "name": "Parent",
            "language": "eng",
            "type": "messaging",
            "nodes": []
        }
    },
    {
        "type": "field",
        "url": "http://testserver/assets/field",
        "content": [
            {"key": "gender", "label": "Gender", "value_type": "text"},
            {"key": "age", "label": "Age", "value_type": "number"},
            {"key": "join_date", "label": "Join Date", "value_type": "datetime"},
            {"key": "activation_token", "label": "Activation Token", "value_type": "text"}
        ]
    },
    {
        "type": "group",
        "url": "http://testserver/assets/group",
        "content": [
            {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Testers"},
            {"uuid": "4f1f98fc-27a7-4a69-bbdb-24744ba739a9", "name": "Males"},
            {"uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a", "name": "Customers"}
        ]
    },
    {
        "type": "label",
        "url": "http://testserver/assets/label",
        "content": [
            {
                "uuid": "3f65d88a-95dc-4140-9451-943e94e06fea",
                "name": "Spam"
            }
        ]
    },
    {
        "type": "location_hierarchy",
        "url": "http://testserver/assets/location_hierarchy",
        "content": [
            {
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
        ]
    },
    {
        "type": "resthook",
        "url": "http://testserver/assets/resthook",
        "content": [
            {
                "slug": "new-registration", 
                "subscribers": [
                    "http://localhost/?cmd=success"
                ]
            }
        ]
    }
]`

var sessionTrigger = `{
    "type": "flow_action",
    "triggered_on": "2017-12-31T11:31:15.035757258-02:00",
    "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
    "contact": {
        "uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f",
        "id": 1234567,
        "name": "Ryan Lewis",
        "language": "eng",
        "timezone": "America/Guayaquil",
        "created_on": "2018-06-20T11:40:30.123456789-00:00",
        "urns": [
            "tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d", 
            "twitterid:54784326227#nyaruka",
            "mailto:foo@bar.com"
        ],
        "groups": [
            {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Testers"},
            {"uuid": "4f1f98fc-27a7-4a69-bbdb-24744ba739a9", "name": "Males"}
        ],
        "fields": {
            "gender": {
                "text": "Male"
            },
            "join_date": {
                "text": "2017-12-02", "datetime": "2017-12-02T00:00:00-02:00"
            },
            "activation_token": {
                "text": "AACC55"
            }
        }
    },
    "run": {
        "uuid": "4213ac47-93fd-48c4-af12-7da8218ef09d",
        "contact": {
            "uuid": "c59b0033-e748-4240-9d4c-e85eb6800151",
            "name": "Jasmine",
            "language": "spa",
            "urns": [],
            "fields": {
                "age": {
                    "text": "33 years", "number": 33
                },
                "gender": {
                    "text": "Female"
                }
            }
        },
        "flow": {
            "uuid": "fece6eac-9127-4343-9269-56e88f391562",
            "name": "Parent Flow"
        },
        "results": {
            "role": {
                "created_on": "2000-01-01T00:00:00.000000000-00:00",
                "input": "a reporter",
                "name": "Role",
                "node_uuid": "385cb848-5043-448e-9123-05cbcf26ad74",
                "value": "reporter",
                "category": "Reporter"
            }
        },
        "status": "active"
    },
    "environment": {
        "date_format": "YYYY-MM-DD",
        "languages": [
            "eng",
            "spa"
        ],
        "redaction_policy": "none",
        "time_format": "hh:mm",
        "timezone": "America/Guayaquil"
    },
    "params": {"source": "website","address": {"state": "WA"}}
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

// CreateTestSession creates a standard example session for testing
func CreateTestSession(testServerURL string, actionToAdd flows.Action) (flows.Session, error) {
	// different tests different ports for the test HTTP server
	sessionAssets = strings.Replace(sessionAssets, "http://localhost", testServerURL, -1)

	session, err := CreateSession(json.RawMessage(sessionAssets))
	if err != nil {
		return nil, err
	}

	// optional modify the main flow by adding the provided action to the final empty node
	if actionToAdd != nil {
		flow, _ := session.Assets().GetFlow(flows.FlowUUID("50c3706e-fedb-42c0-8eab-dda3335714b7"))
		flow.Nodes()[2].AddAction(actionToAdd)
	}

	// read our trigger
	triggerEnvelope := &utils.TypedEnvelope{}
	if err := triggerEnvelope.UnmarshalJSON(json.RawMessage(sessionTrigger)); err != nil {
		return nil, fmt.Errorf("error unmarsalling trigger: %s", err)
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

// CreateSession creates a session with the given assets
func CreateSession(sessionAssets json.RawMessage) (flows.Session, error) {
	// load our assets into a cache
	assetCache := assets.NewAssetCache(100, 5)
	err := assetCache.Include(sessionAssets)
	if err != nil {
		return nil, err
	}

	// create our engine session
	session := engine.NewSession(engine.NewMockAssetServer(assetCache), engine.NewDefaultConfig(), TestHTTPClient)
	return session, nil
}
