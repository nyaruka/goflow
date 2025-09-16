package actions_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"testing"
	"text/template"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/services/classification/wit"
	"github.com/nyaruka/goflow/services/email/smtp"
	"github.com/nyaruka/goflow/services/webhooks"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/test/services"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/smtpx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaultContactJSON = []byte(`{
	"uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f",
	"name": "Ryan Lewis",
	"status": "active",
	"language": "eng",
	"timezone": "America/Guayaquil",
	"urns": [
		"tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d&id=123",
		"twitterid:54784326227#nyaruka"
	],
	"groups": [
		{"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Testers"},
		{"uuid": "0ec97956-c451-48a0-a180-1ce766623e31", "name": "Males"}
	],
	"fields": {
		"gender": {
			"text": "Male"
		}
	},
	"created_on": "2018-06-20T11:40:30.123456789-00:00"
}`)

func TestActionTypes(t *testing.T) {
	assetsJSON, err := os.ReadFile("testdata/_assets.json")
	require.NoError(t, err)

	typeNames := make([]string, 0)
	for typeName := range actions.RegisteredTypes() {
		typeNames = append(typeNames, typeName)
	}

	sort.Strings(typeNames)

	for _, typeName := range typeNames {
		testActionType(t, assetsJSON, typeName)
	}
}

func testActionType(t *testing.T, assetsJSON []byte, typeName string) {
	testPath := fmt.Sprintf("testdata/%s.json", typeName)
	testFile, err := os.ReadFile(testPath)
	require.NoError(t, err)

	tests := []struct {
		Description  string               `json:"description"`
		HTTPMocks    *httpx.MockRequestor `json:"http_mocks,omitempty"`
		SMTPError    string               `json:"smtp_error,omitempty"`
		Contact      json.RawMessage      `json:"contact,omitempty"`
		HasTicket    bool                 `json:"has_ticket,omitempty"`
		NoInput      bool                 `json:"no_input,omitempty"`
		RedactURNs   bool                 `json:"redact_urns,omitempty"`
		AsBatch      bool                 `json:"as_batch,omitempty"`
		Action       json.RawMessage      `json:"action"`
		Localization json.RawMessage      `json:"localization,omitempty"`
		InFlowType   flows.FlowType       `json:"in_flow_type,omitempty"`

		ReadError         string          `json:"read_error,omitempty"`
		DependenciesError string          `json:"dependencies_error,omitempty"`
		SkipValidation    bool            `json:"skip_validation,omitempty"`
		Events            json.RawMessage `json:"events,omitempty"`
		Webhook           json.RawMessage `json:"webhook,omitempty"`
		ContactAfter      json.RawMessage `json:"contact_after,omitempty"`
		LocalsAfter       json.RawMessage `json:"locals_after,omitempty"`
		Templates         []string        `json:"templates,omitempty"`
		LocalizedText     []string        `json:"localizables,omitempty"`
		Inspection        json.RawMessage `json:"inspection,omitempty"`
	}{}

	jsonx.MustUnmarshal(testFile, &tests)

	defer httpx.SetRequestor(httpx.DefaultRequestor)
	defer smtpx.SetSender(smtpx.DefaultSender)

	for i, tc := range tests {
		test.MockUniverse()

		var clonedMocks *httpx.MockRequestor
		if tc.HTTPMocks != nil {
			httpx.SetRequestor(tc.HTTPMocks)
			clonedMocks = tc.HTTPMocks.Clone()
		} else {
			httpx.SetRequestor(httpx.DefaultRequestor)
		}
		if tc.SMTPError != "" {
			smtpx.SetSender(smtpx.NewMockSender(errors.New(tc.SMTPError)))
		} else {
			smtpx.SetSender(smtpx.NewMockSender(nil))
		}

		testName := fmt.Sprintf("test '%s' for action type '%s'", tc.Description, typeName)

		// pick a suitable "holder" flow in our assets JSON
		flowIndex := 0
		flowUUID := assets.FlowUUID("bead76f5-dac4-4c9d-996c-c62b326e8c0a")
		if tc.InFlowType == flows.FlowTypeVoice {
			flowIndex = 1
			flowUUID = assets.FlowUUID("7a84463d-d209-4d3e-a0ff-79f977cd7bd0")
		}

		// inject the action into a suitable node's actions in that flow
		actionsPath := []string{"flows", fmt.Sprintf("[%d]", flowIndex), "nodes", "[0]", "actions"}
		actionsJson := fmt.Appendf(nil, "[%s]", string(tc.Action))
		testAssetsJSON := test.JSONReplace(assetsJSON, actionsPath, actionsJson)

		// if we have a localization section, inject that too
		if tc.Localization != nil {
			localizationPath := []string{"flows", fmt.Sprintf("[%d]", flowIndex), "localization"}
			testAssetsJSON = test.JSONReplace(testAssetsJSON, localizationPath, tc.Localization)
		}

		// create session assets
		sa, err := test.CreateSessionAssets(testAssetsJSON, "")
		require.NoError(t, err, "unable to create session assets in %s", testName)

		// now try to read the flow, and if we expect a read error, check that
		flow, err := sa.Flows().Get(flowUUID)
		if tc.ReadError != "" {
			rootErr := test.RootError(err)
			assert.EqualError(t, rootErr, tc.ReadError, "read error mismatch in %s", testName)
			continue
		} else {
			assert.NoError(t, err, "unexpected read error in %s", testName)
		}

		envBuilder := envs.NewBuilder().
			WithAllowedLanguages("eng", "spa").
			WithDefaultCountry("RW")

		if tc.RedactURNs {
			envBuilder.WithRedactionPolicy(envs.RedactionPolicyURNs)
		}

		env := envBuilder.Build()

		// optionally load our contact
		contactJSON := defaultContactJSON
		if tc.Contact != nil {
			contactJSON = tc.Contact
		}

		contact, err := flows.ReadContact(sa, contactJSON, assets.PanicOnMissing)
		require.NoError(t, err)

		if tc.HasTicket {
			topic := sa.Topics().Get("0d9a2c56-6fc2-4f27-93c5-a6322e26b740")
			contact.Tickets().Add(flows.NewTicket("7f44b065-ec28-4d7a-bbb4-0bda3b75b19d", flows.TicketStatusOpen, topic, nil))
			contact.ReevaluateQueryBasedGroups(env)
		}

		// and switch their language
		if tc.Localization != nil {
			contact.SetLanguage(i18n.Language("spa"))
		}

		var trigger flows.Trigger
		var call *flows.Call

		if tc.NoInput || tc.AsBatch {
			tb := triggers.NewBuilder(flow.Reference(false)).Manual().AsBatch()

			if flow.Type() == flows.FlowTypeVoice {
				channel := sa.Channels().Get("57f1078f-88aa-46f4-a59a-948a5739c03d")
				call = flows.NewCall("01978a2f-ad9a-7f2e-ad44-6e7547078cec", channel, urns.URN("tel:+12065551212"))
			}

			trigger = tb.Build()
		} else {
			msg := flows.NewMsgIn(
				urns.URN("tel:+12065551212"),
				nil,
				"Hi everybody",
				[]utils.Attachment{
					"image/jpeg:http://http://s3.amazon.com/bucket/test.jpg",
					"audio/mp3:http://s3.amazon.com/bucket/test.mp3",
				},
				"",
			)
			trigger = triggers.NewBuilder(flow.Reference(false)).MsgReceived(events.NewMsgReceived(msg)).Build()
		}

		// create an engine instance
		eng := engine.NewBuilder().
			WithLLMPrompts(map[string]*template.Template{
				"categorize": template.Must(template.New("").Parse("Categorize the following text into one of the following categories and only return that category or <CANT> if you can't: {{ .arg1 }}")),
			}).
			WithEmailServiceFactory(func(flows.SessionAssets) (flows.EmailService, error) {
				return smtp.NewService("smtp://nyaruka:pass123@mail.temba.io?from=flows@temba.io", nil)
			}).
			WithWebhookServiceFactory(webhooks.NewServiceFactory(http.DefaultClient, nil, nil, map[string]string{"User-Agent": "goflow-testing"}, 100000)).
			WithClassificationServiceFactory(func(c *flows.Classifier) (flows.ClassificationService, error) {
				if c.Type() == "wit" {
					return wit.NewService(http.DefaultClient, nil, c, "123456789"), nil
				}
				return nil, fmt.Errorf("no classification service available for %s", c.Reference())
			}).
			WithLLMServiceFactory(func(l *flows.LLM) (flows.LLMService, error) {
				return services.NewLLM(), nil
			}).
			WithAirtimeServiceFactory(func(flows.SessionAssets) (flows.AirtimeService, error) {
				return services.NewAirtime("RWF"), nil
			}).
			Build()

		// create session
		session, sprint, err := eng.NewSession(context.Background(), sa, env, contact, trigger, call)
		require.NoError(t, err)

		// check all http mocks were used
		if tc.HTTPMocks != nil {
			require.False(t, tc.HTTPMocks.HasUnused(), "unused HTTP mocks in %s", testName)
		}

		// clone test case and populate with actual values
		actual := tc
		actual.HTTPMocks = clonedMocks

		// re-marshal the action
		actual.Action, err = jsonx.Marshal(flow.Nodes()[0].Actions()[0])
		require.NoError(t, err)

		// and the events
		run := session.Runs()[0]
		actual.Events = jsonx.MustMarshal(sprint.Events()[1 : len(sprint.Events())-1]) // trim initial run_started and final run_ended

		if tc.Webhook != nil {
			actual.Webhook = jsonx.MustMarshal(run.Webhook())
		}
		if tc.ContactAfter != nil {
			actual.ContactAfter, _ = jsonx.Marshal(session.Contact())
		}
		if tc.LocalsAfter != nil {
			actual.LocalsAfter, _ = jsonx.Marshal(run.Locals())
		}
		if tc.Templates != nil {
			actual.Templates = flow.ExtractTemplates()
		}
		if tc.LocalizedText != nil {
			actual.LocalizedText = flow.ExtractLocalizables()
		}
		if tc.Inspection != nil {
			actual.Inspection, _ = jsonx.Marshal(flow.Inspect(sa))
		}

		if !test.UpdateSnapshots {
			// check the action marshaled correctly
			test.AssertEqualJSON(t, tc.Action, actual.Action, "marshal mismatch in %s", testName)

			// check events are what we expected
			test.AssertEqualJSON(t, tc.Events, actual.Events, "events mismatch in %s", testName)

			// check webhook is in expected state
			if tc.Webhook != nil {
				test.AssertEqualJSON(t, tc.Webhook, actual.Webhook, "webhook mismatch in %s", testName)
			}

			// check contact is in the expected state
			if tc.ContactAfter != nil {
				test.AssertEqualJSON(t, tc.ContactAfter, actual.ContactAfter, "contact mismatch in %s", testName)
			}

			// check locals match
			if tc.LocalsAfter != nil {
				test.AssertEqualJSON(t, tc.LocalsAfter, actual.LocalsAfter, "locals mismatch in %s", testName)
			}

			// check extracted templates
			if tc.Templates != nil {
				assert.Equal(t, tc.Templates, actual.Templates, "extracted templates mismatch in %s", testName)
			}

			// check extracted localized text
			if tc.LocalizedText != nil {
				assert.Equal(t, tc.LocalizedText, actual.LocalizedText, "extracted localized text mismatch in %s", testName)
			}

			// check inspection results
			if tc.Inspection != nil {
				test.AssertEqualJSON(t, tc.Inspection, actual.Inspection, "inspection mismatch in %s", testName)
			}
		} else {
			tests[i] = actual
		}
	}

	if test.UpdateSnapshots {
		actualJSON, err := jsonx.MarshalPretty(tests)
		require.NoError(t, err)

		err = os.WriteFile(testPath, actualJSON, 0666)
		require.NoError(t, err)
	}
}

func TestConstructors(t *testing.T) {
	actionUUID := flows.ActionUUID("ad154980-7bf7-4ab8-8728-545fd6378912")

	tests := []struct {
		action flows.Action
		json   string
	}{
		{
			actions.NewAddContactGroups(
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
			actions.NewAddContactURN(
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
			actions.NewAddInputLabels(
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
			actions.NewCallClassifier(
				actionUUID,
				assets.NewClassifierReference("0baee364-07a7-4c93-9778-9f55a35903bb", "Booking"),
				"@input.text",
				"Intent",
			),
			`{
			"type": "call_classifier",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"classifier": {
				"uuid": "0baee364-07a7-4c93-9778-9f55a35903bb",
				"name": "Booking"
			},
			"input": "@input.text",
			"result_name": "Intent"
		}`,
		},
		{
			actions.NewCallLLM(
				actionUUID,
				assets.NewLLMReference("0baee364-07a7-4c93-9778-9f55a35903bb", "GPT-4"),
				"Tell a joke about a person with this name",
				"@contact.name",
				"the_joke",
			),
			`{
			"type": "call_llm",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"llm": {
				"uuid": "0baee364-07a7-4c93-9778-9f55a35903bb",
				"name": "GPT-4"
			},
			"instructions": "Tell a joke about a person with this name",
			"input": "@contact.name",
			"output_local": "the_joke"
		}`,
		},
		{
			actions.NewCallResthook(
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
			actions.NewCallWebhook(
				actionUUID,
				"POST",
				"http://example.com/ping",
				map[string]string{
					"Authentication": "Token @fields.token",
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
				"Authentication": "Token @fields.token"
			},
			"body": "{\"contact_id\": 234}",
			"result_name": "Webhook Response"
		}`,
		},
		{
			actions.NewOpenTicket(
				actionUUID,
				assets.NewTopicReference("472a7a73-96cb-4736-b567-056d987cc5b4", "Weather"),
				"Where are my cookies?",
				assets.NewUserReference("0c78ef47-7d56-44d8-8f57-96e0f30e8f44", "Bob McTickets"),
			),
			`{
				"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
				"type": "open_ticket",
				"topic": {
					"uuid": "472a7a73-96cb-4736-b567-056d987cc5b4",
					"name": "Weather"
				},
				"note": "Where are my cookies?",
				"assignee": {
					"uuid": "0c78ef47-7d56-44d8-8f57-96e0f30e8f44",
					"name": "Bob McTickets"
				}
			}`,
		},
		{
			actions.NewPlayAudio(
				actionUUID,
				"http://uploads.temba.io/2353262.m4a",
			),
			`{
				"type": "play_audio",
				"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
				"audio_url": "http://uploads.temba.io/2353262.m4a"
			}`,
		},
		{
			actions.NewSayMsg(
				actionUUID,
				"Hi @contact.name, are you ready to complete today's survey?",
				"http://uploads.temba.io/2353262.m4a",
			),
			`{
			"type": "say_msg",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"audio_url": "http://uploads.temba.io/2353262.m4a",
			"text": "Hi @contact.name, are you ready to complete today's survey?"
		}`,
		},
		{
			actions.NewRemoveContactGroups(
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
			]
		}`,
		},
		{
			actions.NewSendBroadcast(
				actionUUID,
				"Hi there",
				[]string{"http://example.com/red.jpg"},
				[]string{"Red", "Blue"},
				[]*assets.GroupReference{
					assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"),
				},
				[]*flows.ContactReference{
					flows.NewContactReference(flows.ContactUUID("cbe87f5c-cda2-4f90-b5dd-0ac93a884950"), "Bob Smith"),
				},
				"fields.age > 20",
				[]urns.URN{"twitter:nyaruka"},
				nil,
			),
			`{
			"type": "send_broadcast",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"text": "Hi there",
			"attachments": ["http://example.com/red.jpg"],
			"quick_replies": ["Red", "Blue"],
			"groups": [
				{
					"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
					"name": "Testers"
				}
			],
            "contacts": [
				{
					"uuid": "cbe87f5c-cda2-4f90-b5dd-0ac93a884950",
					"name": "Bob Smith"
				}
			],
			"contact_query": "fields.age > 20",
			"urns": ["twitter:nyaruka"]
		}`,
		},
		{
			actions.NewSendEmail(
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
			actions.NewRequestOptIn(
				actionUUID,
				assets.NewOptInReference("248be71d-78e9-4d71-a6c4-9981d369e5cb", "Joke Of The Day"),
			),
			`{
			"type": "request_optin",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"optin": {
				"uuid": "248be71d-78e9-4d71-a6c4-9981d369e5cb",
				"name": "Joke Of The Day"
			}
		}`,
		},
		{
			actions.NewSendMsg(
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
			actions.NewSetContactChannel(
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
			actions.NewSetContactField(
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
			actions.NewSetContactLanguage(
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
			actions.NewSetContactName(
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
			actions.NewSetContactStatus(
				actionUUID,
				flows.ContactStatusBlocked,
			),
			`{
				"type": "set_contact_status",
				"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
				"status": "blocked"
			}`,
		},
		{
			actions.NewSetContactTimezone(
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
			actions.NewSetRunResult(
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
			actions.NewEnterFlow(
				actionUUID,
				assets.NewFlowReference(assets.FlowUUID("fece6eac-9127-4343-9269-56e88f391562"), "Parent"),
				true, // terminal
			),
			`{
			"type": "enter_flow",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"flow": {
				"uuid": "fece6eac-9127-4343-9269-56e88f391562",
				"name": "Parent"
			},
			"terminal": true
		}`,
		},
		{
			actions.NewStartSession(
				actionUUID,
				assets.NewFlowReference(assets.FlowUUID("fece6eac-9127-4343-9269-56e88f391562"), "Parent"),

				[]*assets.GroupReference{
					assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"),
				},
				[]*flows.ContactReference{
					flows.NewContactReference(flows.ContactUUID("cbe87f5c-cda2-4f90-b5dd-0ac93a884950"), "Bob Smith"),
				},
				"fields.age > 20",
				[]urns.URN{"twitter:nyaruka"},
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
			"groups": [
				{
					"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
					"name": "Testers"
				}
			],
            "contacts": [
				{
					"uuid": "cbe87f5c-cda2-4f90-b5dd-0ac93a884950",
					"name": "Bob Smith"
				}
			],
			"contact_query": "fields.age > 20",
			"urns": ["twitter:nyaruka"],
			"exclusions": {},
			"create_contact": true
		}`,
		},
	}

	for _, tc := range tests {
		// test marshaling the action
		actualJSON, err := jsonx.Marshal(tc.action)
		assert.NoError(t, err)

		test.AssertEqualJSON(t, []byte(tc.json), actualJSON, "new action produced unexpected JSON")
	}
}

func TestReadAction(t *testing.T) {
	// error if no type field
	_, err := actions.Read([]byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = actions.Read([]byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}

func TestResthookPayload(t *testing.T) {
	test.MockUniverse()

	server := test.NewHTTPServer(49999, test.MockWebhooksHandler)
	defer server.Close()

	session, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	require.NoError(t, err)
	run := session.Runs()[0]

	log := test.NewEventLog()
	payload, _ := run.EvaluateTemplate(actions.ResthookPayload, log.Log)
	require.NoError(t, log.Error())

	pretty, err := jsonx.MarshalPretty(json.RawMessage(payload))
	require.NoError(t, err)

	test.AssertSnapshot(t, "resthook_payload", string(pretty))
}

func TestStartSessionLoopProtection(t *testing.T) {
	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(`{
		"flows": [
			{
				"uuid": "5472a1c3-63e1-484f-8485-cc8ecb16a058",
				"name": "Inception",
				"spec_version": "13.1.0",
				"language": "eng",
				"type": "messaging",
				"nodes": [
					{
						"uuid": "cc49453a-78ed-48a6-8b94-318b46517071",
						"actions": [
							{
								"uuid": "cdf981ae-a9cf-4c32-98f3-65bac07bf990",
								"type": "start_session",
								"flow": {
									"uuid": "5472a1c3-63e1-484f-8485-cc8ecb16a058", 
									"name": "Inception"
								},
								"contacts": [
									{
										"uuid": "51b41bf2-b2e1-439b-ab9b-dd4c9cef6266", 
										"name": "Bob"
									}
								]
							}
						],
						"exits": [
							{
								"uuid": "717ee506-7b2d-4a18-b142-eafed0c5e9d8"
							}
						]
					}
				]
			}
		]
	}`))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	flow := assets.NewFlowReference("5472a1c3-63e1-484f-8485-cc8ecb16a058", "Inception")
	contact := flows.NewEmptyContact(sa, "Bob", i18n.Language("eng"), nil)

	eng := engine.NewBuilder().Build()
	_, sprint, err := eng.NewSession(context.Background(), sa, env, contact, triggers.NewBuilder(flow).Manual().Build(), nil)
	require.NoError(t, err)

	sessions := make([]flows.Session, 0)
	var session flows.Session

	for {
		// look for a session triggered event
		var event *events.SessionTriggered
		for _, e := range sprint.Events() {
			if e.Type() == events.TypeSessionTriggered {
				event = e.(*events.SessionTriggered)
			}
		}

		// if it exists, trigger a new session
		if event != nil {
			trigger := triggers.NewBuilder(flow).FlowAction(event.History, event.RunSummary).Build()

			session, sprint, err = eng.NewSession(context.Background(), sa, env, contact, trigger, nil)
			require.NoError(t, err)

			sessions = append(sessions, session)
		} else {
			break
		}
	}

	assert.Equal(t, 5, len(sessions))

	// final session should have an error event
	penultimateEvent := sprint.Events()[len(sprint.Events())-2]
	assert.Equal(t, events.TypeError, penultimateEvent.Type())
	assert.Equal(t, "too many sessions have been spawned since the last time input was received", penultimateEvent.(*events.Error).Text)
	finalEvent := sprint.Events()[len(sprint.Events())-1]
	assert.Equal(t, events.TypeRunEnded, finalEvent.Type())
}

func TestStartSessionLoopProtectionWithInput(t *testing.T) {
	ctx := context.Background()
	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(`{
		"flows": [
			{
				"uuid": "5472a1c3-63e1-484f-8485-cc8ecb16a058",
				"name": "Inception",
				"spec_version": "13.1.0",
				"language": "eng",
				"type": "messaging",
				"nodes": [
					{
						"uuid": "bff26109-b7b4-465f-8c9e-ddbf465af5f1",
						"actions": [
							{
								"uuid": "b3779a48-351c-499f-a2ba-f497b3248659",
								"type": "send_msg",
								"text": "Say something"
							}
						],
						"exits": [
							{
								"uuid": "a5efd0ef-303f-4ae9-9f86-0c7ddaf3abf1",
								"destination_uuid": "4c83851e-f0bf-4c59-ba11-220ecccfcebb"
							}
						]
					},
					{
						"uuid": "4c83851e-f0bf-4c59-ba11-220ecccfcebb",
						"router": {
							"type": "switch",
							"wait": {
								"type": "msg"
							},
							"categories": [
								{
									"uuid": "37d8813f-1402-4ad2-9cc2-e9054a96525b",
									"name": "All Responses",
									"exit_uuid": "fc2fcd23-7c4a-44bd-a8c6-6c88e6ed09f8"
								}
							],
							"default_category_uuid": "37d8813f-1402-4ad2-9cc2-e9054a96525b",
							"result_name": "Name",
							"operand": "@input.text",
							"cases": []
						},
						"exits": [
							{
								"uuid": "fc2fcd23-7c4a-44bd-a8c6-6c88e6ed09f8",
								"destination_uuid": "cc49453a-78ed-48a6-8b94-318b46517071"
							}
						]
					},
					{
						"uuid": "cc49453a-78ed-48a6-8b94-318b46517071",
						"actions": [
							{
								"uuid": "cdf981ae-a9cf-4c32-98f3-65bac07bf990",
								"type": "start_session",
								"flow": {
									"uuid": "5472a1c3-63e1-484f-8485-cc8ecb16a058", 
									"name": "Inception"
								},
								"contacts": [
									{
										"uuid": "51b41bf2-b2e1-439b-ab9b-dd4c9cef6266", 
										"name": "Bob"
									}
								]
							}
						],
						"exits": [
							{
								"uuid": "717ee506-7b2d-4a18-b142-eafed0c5e9d8"
							}
						]
					}
				]
			}
		]
	}`))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	flow := assets.NewFlowReference("5472a1c3-63e1-484f-8485-cc8ecb16a058", "Inception")
	contact := flows.NewEmptyContact(sa, "Bob", i18n.Language("eng"), nil)

	eng := engine.NewBuilder().Build()
	session, sprint, err := eng.NewSession(context.Background(), sa, env, contact, triggers.NewBuilder(flow).Manual().Build(), nil)
	require.NoError(t, err)

	sessions := make([]flows.Session, 0)

	for {
		if len(sessions) == 10 {
			break
		}

		if session.Status() == flows.SessionStatusWaiting {
			resume := resumes.NewMsg(events.NewMsgReceived(flows.NewMsgIn(urns.NilURN, nil, "Hi there", nil, "SMS1234")))
			sprint, err = session.Resume(ctx, resume)
			require.NoError(t, err)
		}

		// look for a session triggered event
		var event *events.SessionTriggered
		for _, e := range sprint.Events() {
			if e.Type() == events.TypeSessionTriggered {
				event = e.(*events.SessionTriggered)
			}
		}

		// if it exists, trigger a new session
		if event != nil {
			trigger := triggers.NewBuilder(flow).FlowAction(event.History, event.RunSummary).Build()

			session, sprint, err = eng.NewSession(ctx, sa, env, contact, trigger, nil)
			require.NoError(t, err)

			sessions = append(sessions, session)
		} else {
			break
		}
	}

	assert.Equal(t, 10, len(sessions))
}
