package test

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/triggers"

	"github.com/pkg/errors"
)

var sessionAssets = `{
    "channels": [
        {
            "uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d",
            "name": "My Android Phone",
            "address": "+17036975131",
            "schemes": ["tel"],
            "roles": ["send", "receive"],
            "country": "US"
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
    "classifiers": [
        {
            "uuid": "1c06c884-39dd-4ce4-ad9f-9a01cbe6c000",
            "name": "Booking",
            "type": "wit",
            "intents": ["book_flight", "book_hotel"]
        }
    ],
    "ticketers": [
        {
            "uuid": "19dc6346-9623-4fe4-be80-538d493ecdf5",
            "name": "Support Tickets",
            "type": "mailgun"
        }
    ],
    "flows": [
        {
            "uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
            "name": "Registration",
            "spec_version": "13.0",
            "language": "eng",
            "type": "messaging",
            "revision": 123,
            "nodes": [
                {
                    "uuid": "72a1f5df-49f9-45df-94c9-d86f7ea064e5",
                    "actions": [
                        {
                            "uuid": "9487a60e-a6ef-4a88-b35d-894bfe074144",
                            "type": "enter_flow",
                            "flow": {
                                "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
                                "name": "Collect Age"
                            }
                        }
                    ],
                    "exits": [
                        {
                            "uuid": "d7a36118-0a38-4b35-a7e4-ae89042f0d3c",
                            "destination_uuid": "3dcccbb4-d29c-41dd-a01f-16d814c9ab82"
                        }
                    ]
                },
                {
                    "uuid": "3dcccbb4-d29c-41dd-a01f-16d814c9ab82",
                    "router": {
                        "type": "switch",
                        "wait": {
                            "type": "msg"
                        },
                        "categories": [
                            {
                                "uuid": "37d8813f-1402-4ad2-9cc2-e9054a96525b",
                                "name": "All Responses",
                                "exit_uuid": "100f2d68-2481-4137-a0a3-177620ba3c5f"
                            }
                        ],
                        "operand": "@input.text",
                        "default_category_uuid": "37d8813f-1402-4ad2-9cc2-e9054a96525b"
                    },
                    "exits": [
                        {
                            "uuid": "100f2d68-2481-4137-a0a3-177620ba3c5f",
                            "destination_uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03"
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
                            "uuid": "821eef31-c6d2-45b1-8f6a-d396e4959bbf",
                            "type": "set_run_result",
                            "name": "2Factor",
                            "value": "34634624463525"
                        },
                        {
                            "uuid": "06153fbd-3e2c-413a-b0df-ed15d631835a",
                            "type": "call_webhook",
                            "method": "GET",
                            "url": "http://localhost/?content=%7B%22results%22%3A%5B%7B%22state%22%3A%22WA%22%7D%2C%7B%22state%22%3A%22IN%22%7D%5D%7D",
                            "result_name": "webhook"
                        },
                        {
                            "uuid": "bd821625-5254-40ca-be17-e9a4dc5bde99",
                            "type": "call_classifier",
                            "classifier": {
                                "uuid": "1c06c884-39dd-4ce4-ad9f-9a01cbe6c000",
                                "name": "Booking"
                            },
                            "input": "@input.text",
                            "result_name": "Intent"
                        }
                    ],
                    "exits": [
                        {
                            "uuid": "d898f9a4-f0fc-4ac4-a639-c98c602bb511",
                            "destination_uuid": "c0781400-737f-4940-9a6c-1ec1c3df0325"
                        }
                    ]
                },
                {
                    "uuid": "c0781400-737f-4940-9a6c-1ec1c3df0325",
                    "actions": [],
                    "exits": [
                        {
                            "uuid": "9fc5f8b4-2247-43db-b899-ab1ac50ba06c"
                        }
                    ]
                }
            ]
        },
        {
            "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
            "name": "Collect Age",
            "spec_version": "13.0",
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
                        "value": "@results.age.value"
                    }
                ],
                "exits": [
                    {
                        "uuid": "4ee148c8-4026-41da-9d4c-08cb4d60b0d7"
                    }
                ]
            }]
        },
        {
            "uuid": "fece6eac-9127-4343-9269-56e88f391562",
            "name": "Parent",
            "spec_version": "13.0",
            "language": "eng",
            "type": "messaging",
            "nodes": []
        },
        {
            "uuid": "aa71426e-13bd-4607-a4f5-77666ff9c4bf",
            "name": "Voice Test",
            "spec_version": "13.0",
            "language": "eng",
            "type": "voice",
            "nodes": []
        }
    ],
    "fields": [
        {"uuid": "d66a7823-eada-40e5-9a3a-57239d4690bf", "key": "gender", "name": "Gender", "type": "text"},
        {"uuid": "f1b5aea6-6586-41c7-9020-1a6326cc6565", "key": "age", "name": "Age", "type": "number"},
        {"uuid": "6c86d5ab-3fd9-4a5c-a5b6-48168b016747", "key": "join_date", "name": "Join Date", "type": "datetime"},
        {"uuid": "c88d2640-d124-438a-b666-5ec53a353dcd", "key": "activation_token", "name": "Activation Token", "type": "text"},
        {"uuid": "ab9c0631-d8cd-4e77-a5a2-66a8b077e385", "key": "state", "name": "State", "type": "state"},
        {"uuid": "3bfc3908-a402-48ea-841c-b73b5ef3a254", "key": "not_set", "name": "Not set", "type": "text"}
    ],
    "groups": [
        {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Testers"},
        {"uuid": "4f1f98fc-27a7-4a69-bbdb-24744ba739a9", "name": "Males"},
        {"uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a", "name": "Customers"}
    ],
    "labels": [
        {"uuid": "3f65d88a-95dc-4140-9451-943e94e06fea", "name": "Spam"}
    ],
    "locations": [
        {
            "name": "Rwanda",
            "aliases": ["Ruanda"],		
            "children": [
                {
                    "name": "Kigali City",
                    "aliases": ["Kigali", "Kigari"],
                    "children": [
                        {
                            "name": "Gasabo",
                            "children": [
                                {
                                    "name": "Gisozi"
                                },
                                {
                                    "name": "Ndera"
                                }
                            ]
                        },
                        {
                            "name": "Nyarugenge",
                            "children": []
                        }
                    ]
                }
            ]
        }
    ],
    "resthooks": [
        {
            "slug": "new-registration", 
            "subscribers": [
                "http://localhost/?cmd=success"
            ]
        }
    ],
    "users": [
        {
            "email": "bob@nyaruka.com",
            "name": "Bob"
        }
    ]
}`

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
            "tel:+12024561111?channel=57f1078f-88aa-46f4-a59a-948a5739c03d", 
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
        },
        "tickets": [
            {
                "uuid": "e5f5a9b0-1c08-4e56-8f5c-92e00bc3cf52",
                "subject": "Old ticket",
                "body": "I have a problem",
                "ticketer": {
                    "name": "Support Tickets",
                    "uuid": "19dc6346-9623-4fe4-be80-538d493ecdf5"
                }
            },
            {
                "uuid": "78d1fe0d-7e39-461e-81c3-a6a25f15ed69",
                "subject": "Question",
                "body": "What day is it?",
                "ticketer": {
                    "name": "Support Tickets",
                    "uuid": "19dc6346-9623-4fe4-be80-538d493ecdf5"
                },
                "assignee": {"email": "bob@nyaruka.com", "name": "Bob"}
            }
        ]
    },
    "run_summary": {
        "uuid": "4213ac47-93fd-48c4-af12-7da8218ef09d",
        "contact": {
            "uuid": "c59b0033-e748-4240-9d4c-e85eb6800151",
            "name": "Jasmine",
            "created_on": "2018-01-01T12:00:00.000000000-00:00",
            "language": "spa",
            "urns": [
                "tel:+12024562222"
            ],
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
        "date_format": "DD-MM-YYYY",
        "allowed_languages": [
            "eng", 
            "spa"
        ],
        "default_country": "US",
        "redaction_policy": "none",
        "time_format": "tt:mm",
        "timezone": "America/Guayaquil"
    },
    "params": {"source": "website","address": {"state": "WA"}}
}`

var sessionResume = `{
    "type": "msg",
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
    "resumed_on": "2017-12-31T11:35:10.035757258-02:00"
}`

var voiceSessionAssets = `{
    "channels": [
        {
            "uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d",
            "name": "My Android Phone",
            "address": "+17036975131",
            "schemes": ["tel"],
            "roles": ["send", "receive"],
            "country": "US"
        },
        {
            "uuid": "fd47a886-451b-46fb-bcb6-242a4046c0c0",
            "name": "Nexmo",
            "address": "+12024560010",
            "schemes": ["tel"],
            "roles": ["send", "receive", "call", "answer"]
        }
    ],
    "flows": [
        {
            "uuid": "aa71426e-13bd-4607-a4f5-77666ff9c4bf",
            "name": "Voice Test",
            "spec_version": "13.0",
            "language": "eng",
            "type": "voice",
            "nodes": [
                {
                    "uuid": "6da04a32-6c84-40d9-b614-3782fde7af80",
                    "actions": [],
                    "exits": [
                        {
                            "uuid": "9082b6ec-a65f-4677-8b3c-2f8de402ff13"
                        }
                    ]
                }
            ]
        }
    ],
    "fields": [
        {"uuid": "d66a7823-eada-40e5-9a3a-57239d4690bf", "key": "gender", "name": "Gender", "type": "text"}
    ],
    "groups": [
        {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Testers"},
        {"uuid": "4f1f98fc-27a7-4a69-bbdb-24744ba739a9", "name": "Males"},
        {"uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a", "name": "Customers"}
    ]
}`

var voiceSessionTrigger = `{
    "type": "channel",
    "triggered_on": "2017-12-31T11:31:15.035757258-02:00",
    "event": {
        "type": "incoming_call",
        "channel": {"uuid": "fd47a886-451b-46fb-bcb6-242a4046c0c0", "name": "Nexmo"}
    },
    "connection": {
        "channel": {"uuid": "fd47a886-451b-46fb-bcb6-242a4046c0c0", "name": "Nexmo"},
        "urn": "tel:+12065551212"
    },
    "flow": {"uuid": "aa71426e-13bd-4607-a4f5-77666ff9c4bf", "name": "Voice Test"},
    "contact": {
        "uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f",
        "id": 1234567,
        "name": "Ryan Lewis",
        "language": "eng",
        "timezone": "America/Guayaquil",
        "created_on": "2018-06-20T11:40:30.123456789-00:00",
        "urns": [
            "tel:+12024561111"
        ],
        "groups": [
            {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Testers"},
            {"uuid": "4f1f98fc-27a7-4a69-bbdb-24744ba739a9", "name": "Males"}
        ],
        "fields": {
            "gender": {
                "text": "Male"
            }
        }
    },
    "environment": {
        "date_format": "DD-MM-YYYY",
        "allowed_languages": [
            "eng", 
            "spa"
        ],
        "redaction_policy": "none",
        "time_format": "hh:mm",
        "timezone": "America/Guayaquil"
    }
}`

// CreateTestSession creates a standard example session for testing
func CreateTestSession(testServerURL string, redact envs.RedactionPolicy) (flows.Session, []flows.Event, error) {
	assetsJSON := json.RawMessage(sessionAssets)

	sa, err := CreateSessionAssets(assetsJSON, testServerURL)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error creating test session")
	}

	// read our trigger
	triggerJSON := json.RawMessage(sessionTrigger)
	triggerJSON = JSONReplace(triggerJSON, []string{"environment", "redaction_policy"}, []byte(fmt.Sprintf(`"%s"`, redact)))

	trigger, err := triggers.ReadTrigger(sa, triggerJSON, assets.PanicOnMissing)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error reading trigger")
	}

	eng := NewEngine()

	session, _, err := eng.NewSession(sa, trigger)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error starting test session")
	}

	// read our resume
	resume, err := resumes.ReadResume(sa, json.RawMessage(sessionResume), assets.PanicOnMissing)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error reading resume")
	}

	sprint, err := session.Resume(resume)
	return session, sprint.Events(), err
}

// CreateTestVoiceSession creates a standard example session for testing voice flows and actions
func CreateTestVoiceSession(testServerURL string) (flows.Session, []flows.Event, error) {
	assetsJSON := json.RawMessage(voiceSessionAssets)

	sa, err := CreateSessionAssets(assetsJSON, testServerURL)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error creating test voice session assets")
	}

	// read our trigger
	trigger, err := triggers.ReadTrigger(sa, json.RawMessage(voiceSessionTrigger), assets.PanicOnMissing)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error reading trigger")
	}

	eng := NewEngine()
	session, sprint, err := eng.NewSession(sa, trigger)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error starting test voice session")
	}

	return session, sprint.Events(), err
}

// CreateSessionAssets creates assets from given JSON
func CreateSessionAssets(assetsJSON json.RawMessage, testServerURL string) (flows.SessionAssets, error) {
	env := envs.NewBuilder().Build()

	// different tests different ports for the test HTTP server
	if testServerURL != "" {
		assetsJSON = json.RawMessage(strings.Replace(string(assetsJSON), "http://localhost", testServerURL, -1))
	}

	// read our assets into a source
	source, err := static.NewSource(assetsJSON)
	if err != nil {
		return nil, errors.Wrap(err, "error loading test assets")
	}

	// create our engine session
	sa, err := engine.NewSessionAssets(env, source, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating test session assets")
	}

	return sa, nil
}

// CreateSession creates a new session from the give assets
func CreateSession(assetsJSON json.RawMessage, flowUUID assets.FlowUUID) (flows.Session, flows.Sprint, error) {
	sa, err := CreateSessionAssets(assetsJSON, "")
	if err != nil {
		return nil, nil, err
	}

	flow, err := sa.Flows().Get(flowUUID)
	if err != nil {
		return nil, nil, err
	}

	env := envs.NewBuilder().Build()
	contact := flows.NewEmptyContact(sa, "Bob", envs.NilLanguage, nil)
	trigger := triggers.NewBuilder(env, flow.Reference(), contact).Manual().Build()
	eng := engine.NewBuilder().Build()

	return eng.NewSession(sa, trigger)
}

// ResumeSession resumes the given session with potentially different assets
func ResumeSession(session flows.Session, assetsJSON json.RawMessage, msgText string) (flows.Session, flows.Sprint, error) {
	// reload session with new assets
	sessionJSON, err := jsonx.Marshal(session)
	if err != nil {
		return nil, nil, err
	}

	sa, err := CreateSessionAssets(assetsJSON, "")
	if err != nil {
		return nil, nil, err
	}

	// re-use same engine instance
	eng := session.Engine()

	session, err = eng.ReadSession(sa, sessionJSON, assets.IgnoreMissing)
	if err != nil {
		return nil, nil, err
	}

	msg := flows.NewMsgIn(flows.MsgUUID(uuids.New()), urns.NilURN, nil, msgText, nil)

	sprint, err := session.Resume(resumes.NewMsg(session.Environment(), session.Contact(), msg))

	return session, sprint, err
}

// EventLog is a utility for testing things which take an event logger function
type EventLog struct {
	Events []flows.Event
}

// NewEventLog creates a new event log
func NewEventLog() *EventLog {
	return &EventLog{make([]flows.Event, 0)}
}

func (l *EventLog) Log(e flows.Event) {
	l.Events = append(l.Events, e)
}
