package actions_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var actionUUID = flows.ActionUUID("ad154980-7bf7-4ab8-8728-545fd6378912")

var actionTests = []struct {
	action flows.Action
	json   string
}{
	{
		actions.NewAddContactGroupsAction(
			actionUUID,
			[]*assets.GroupReference{
				assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"),
				assets.NewVariableGroupReference("@(format_location(contact.fields.state)) Members"),
			},
		),
		`{
			"type": "add_contact_groups",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"groups": [
				{
					"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
					"name": "Testers"
				},
				{
					"name_match": "@(format_location(contact.fields.state)) Members"
				}
			]
		}`,
	},
	{
		actions.NewAddContactURNAction(
			actionUUID,
			"tel",
			"+234532626677",
		),
		`{
			"type": "add_contact_urn",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"scheme": "tel",
			"path": "+234532626677"
		}`,
	},
	{
		actions.NewAddInputLabelsAction(
			actionUUID,
			[]*assets.LabelReference{
				assets.NewLabelReference(assets.LabelUUID("3f65d88a-95dc-4140-9451-943e94e06fea"), "Spam"),
				assets.NewVariableLabelReference("@(format_location(contact.fields.state)) Messages"),
			},
		),
		`{
			"type": "add_input_labels",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"labels": [
				{
					"uuid": "3f65d88a-95dc-4140-9451-943e94e06fea",
					"name": "Spam"
				},
				{
					"name_match": "@(format_location(contact.fields.state)) Messages"
				}
			]
		}`,
	},
	{
		actions.NewCallResthookAction(
			actionUUID,
			"new-registration",
			"My Result",
		),
		`{
			"type": "call_resthook",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"resthook": "new-registration",
			"result_name": "My Result"
		}`,
	},
	{
		actions.NewCallWebhookAction(
			actionUUID,
			"POST",
			"http://example.com/ping",
			map[string]string{
				"Authentication": "Token @contact.fields.token",
			},
			`{"contact_id": 234}`, // body
			"Webhook Response",
		),
		`{
			"type": "call_webhook",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"method": "POST",
			"url": "http://example.com/ping",
			"headers": {
				"Authentication": "Token @contact.fields.token"
			},
			"body": "{\"contact_id\": 234}",
			"result_name": "Webhook Response"
		}`,
	},
	{
		actions.NewRemoveContactGroupsAction(
			actionUUID,
			[]*assets.GroupReference{
				assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"),
				assets.NewVariableGroupReference("@(format_location(contact.fields.state)) Members"),
			},
			false,
		),
		`{
			"type": "remove_contact_groups",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"groups": [
				{
					"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
					"name": "Testers"
				},
				{
					"name_match": "@(format_location(contact.fields.state)) Members"
				}
			],
			"all_groups": false
		}`,
	},
	{
		actions.NewSendBroadcastAction(
			actionUUID,
			"Hi there",
			[]string{"http://example.com/red.jpg"},
			[]string{"Red", "Blue"},
			[]urns.URN{"twitter:nyaruka"},
			[]*flows.ContactReference{
				flows.NewContactReference(flows.ContactUUID("cbe87f5c-cda2-4f90-b5dd-0ac93a884950"), "Bob Smith"),
			},
			[]*assets.GroupReference{
				assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"),
			},
			nil,
		),
		`{
			"type": "send_broadcast",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"text": "Hi there",
			"attachments": ["http://example.com/red.jpg"],
			"quick_replies": ["Red", "Blue"],
			"urns": ["twitter:nyaruka"],
            "contacts": [
				{
					"uuid": "cbe87f5c-cda2-4f90-b5dd-0ac93a884950",
					"name": "Bob Smith"
				}
			],
			"groups": [
				{
					"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
					"name": "Testers"
				}
			]
		}`,
	},
	{
		actions.NewSendEmailAction(
			actionUUID,
			[]string{"bob@example.com"},
			"Hi there",
			"So I was thinking...",
		),
		`{
			"type": "send_email",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"addresses": ["bob@example.com"],
			"subject": "Hi there",
			"body": "So I was thinking..."
		}`,
	},
	{
		actions.NewSendMsgAction(
			actionUUID,
			"Hi there",
			[]string{"http://example.com/red.jpg"},
			[]string{"Red", "Blue"},
			true,
		),
		`{
			"type": "send_msg",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"text": "Hi there",
			"attachments": ["http://example.com/red.jpg"],
			"quick_replies": ["Red", "Blue"],
			"all_urns": true
		}`,
	},
	{
		actions.NewSetContactChannelAction(
			actionUUID,
			assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), "My Android Phone"),
		),
		`{
			"type": "set_contact_channel",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"channel": {
				"uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d",
				"name": "My Android Phone"
			}
		}`,
	},
	{
		actions.NewSetContactFieldAction(
			actionUUID,
			assets.NewFieldReference("gender", "Gender"),
			"Male",
		),
		`{
			"type": "set_contact_field",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"field": {
				"key": "gender",
				"name": "Gender"
			},
			"value": "Male"
		}`,
	},
	{
		actions.NewSetContactLanguageAction(
			actionUUID,
			"eng",
		),
		`{
			"type": "set_contact_language",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"language": "eng"
		}`,
	},
	{
		actions.NewSetContactNameAction(
			actionUUID,
			"Bob",
		),
		`{
			"type": "set_contact_name",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"name": "Bob"
		}`,
	},
	{
		actions.NewSetContactTimezoneAction(
			actionUUID,
			"Africa/Kigali",
		),
		`{
			"type": "set_contact_timezone",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"timezone": "Africa/Kigali"
		}`,
	},
	{
		actions.NewSetRunResultAction(
			actionUUID,
			"Response 1",
			"yes",
			"Yes",
		),
		`{
			"type": "set_run_result",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"name": "Response 1",
			"value": "yes",
			"category": "Yes"
		}`,
	},
	{
		actions.NewStartFlowAction(
			actionUUID,
			assets.NewFlowReference(assets.FlowUUID("fece6eac-9127-4343-9269-56e88f391562"), "Parent"),
			true, // terminal
		),
		`{
			"type": "start_flow",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"flow": {
				"uuid": "fece6eac-9127-4343-9269-56e88f391562",
				"name": "Parent"
			},
			"terminal": true
		}`,
	},
	{
		actions.NewStartSessionAction(
			actionUUID,
			assets.NewFlowReference(assets.FlowUUID("fece6eac-9127-4343-9269-56e88f391562"), "Parent"),
			[]urns.URN{"twitter:nyaruka"},
			[]*flows.ContactReference{
				flows.NewContactReference(flows.ContactUUID("cbe87f5c-cda2-4f90-b5dd-0ac93a884950"), "Bob Smith"),
			},
			[]*assets.GroupReference{
				assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"),
			},
			nil,  // legacy vars
			true, // create new contact
		),
		`{
			"type": "start_session",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"flow": {
				"uuid": "fece6eac-9127-4343-9269-56e88f391562",
				"name": "Parent"
			},
			"urns": ["twitter:nyaruka"],
            "contacts": [
				{
					"uuid": "cbe87f5c-cda2-4f90-b5dd-0ac93a884950",
					"name": "Bob Smith"
				}
			],
			"groups": [
				{
					"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
					"name": "Testers"
				}
			],
			"create_contact": true
		}`,
	},
}

func TestActions(t *testing.T) {
	session, err := test.CreateTestSession("", nil)
	require.NoError(t, err)

	for _, tc := range actionTests {
		// test validating the action
		err := tc.action.Validate(session.Assets())
		assert.NoError(t, err)

		// test marshaling the action
		actualJSON, err := json.Marshal(tc.action)
		assert.NoError(t, err)

		test.AssertEqualJSON(t, json.RawMessage(tc.json), actualJSON, "new action produced unexpected JSON")
	}
}
