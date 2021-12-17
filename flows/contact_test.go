package flows_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContact(t *testing.T) {
	source, err := static.NewSource([]byte(`{
		"channels": [
			{
				"uuid": "294a14d4-c998-41e5-a314-5941b97b89d7",
				"name": "My Android Phone",
				"address": "+17036975131",
				"schemes": ["tel"],
				"roles": ["send", "receive"],
				"country": "US"
			}
		],
		"ticketers": [
			{
				"uuid": "d605bb96-258d-4097-ad0a-080937db2212",
				"name": "Support Tickets",
				"type": "mailgun"
			}
		],
		"topics": [
			{
				"uuid": "472a7a73-96cb-4736-b567-056d987cc5b4",
				"name": "Weather"
			}
		]
	}`))
	require.NoError(t, err)

	env := envs.NewBuilder().Build()

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	android := sa.Channels().Get("294a14d4-c998-41e5-a314-5941b97b89d7")

	uuids.SetGenerator(uuids.NewSeededGenerator(1234))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	tz, _ := time.LoadLocation("America/Bogota")

	contact, err := flows.NewContact(
		sa,
		flows.ContactUUID(uuids.New()),
		flows.ContactID(12345),
		"Joe Bloggs",
		envs.Language("eng"),
		flows.ContactStatusActive,
		tz,
		time.Date(2017, 12, 15, 10, 0, 0, 0, time.UTC),
		nil,
		nil,
		nil,
		nil,
		nil,
		assets.PanicOnMissing,
	)
	require.NoError(t, err)

	assert.Equal(t, flows.URNList{}, contact.URNs())
	assert.Equal(t, flows.ContactStatusActive, contact.Status())
	assert.Nil(t, contact.LastSeenOn())
	assert.Nil(t, contact.PreferredChannel())

	contact.SetLastSeenOn(time.Date(2018, 12, 15, 10, 0, 0, 0, time.UTC))
	assert.Equal(t, time.Date(2018, 12, 15, 10, 0, 0, 0, time.UTC), *contact.LastSeenOn())

	contact.AddURN(urns.URN("tel:+12024561111?channel=294a14d4-c998-41e5-a314-5941b97b89d7"), nil)
	contact.AddURN(urns.URN("twitter:joey"), nil)
	contact.AddURN(urns.URN("whatsapp:235423721788"), nil)

	assert.Equal(t, "Joe Bloggs", contact.Name())
	assert.Equal(t, flows.ContactID(12345), contact.ID())
	assert.Equal(t, tz, contact.Timezone())
	assert.Equal(t, envs.Language("eng"), contact.Language())
	assert.Equal(t, android, contact.PreferredChannel())
	assert.Equal(t, envs.Country("US"), contact.Country())
	assert.Equal(t, "en-US", contact.Locale(env).ToBCP47())

	contact.SetStatus(flows.ContactStatusStopped)
	assert.Equal(t, flows.ContactStatusStopped, contact.Status())

	contact.SetStatus(flows.ContactStatusBlocked)
	assert.Equal(t, flows.ContactStatusBlocked, contact.Status())

	contact.SetStatus(flows.ContactStatusActive)
	assert.Equal(t, flows.ContactStatusActive, contact.Status())

	assert.True(t, contact.HasURN("tel:+12024561111"))      // has URN
	assert.True(t, contact.HasURN("tel:+120-2456-1111"))    // URN will be normalized
	assert.True(t, contact.HasURN("whatsapp:235423721788")) // has URN
	assert.False(t, contact.HasURN("tel:+16300000000"))     // doesn't have URN

	assert.False(t, contact.RemoveURN("tel:+16300000000"))      // doesn't have URN
	assert.True(t, contact.RemoveURN("whatsapp:235423721788"))  // did have URN
	assert.False(t, contact.RemoveURN("whatsapp:235423721788")) // no longer has URN

	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"discord":    nil,
		"ext":        nil,
		"facebook":   nil,
		"fcm":        nil,
		"freshchat":  nil,
		"instagram":  nil,
		"jiochat":    nil,
		"line":       nil,
		"mailto":     nil,
		"rocketchat": nil,
		"tel":        flows.NewContactURN(urns.URN("tel:+12024561111?channel=294a14d4-c998-41e5-a314-5941b97b89d7"), nil).ToXValue(env),
		"telegram":   nil,
		"twitter":    flows.NewContactURN(urns.URN("twitter:joey"), nil).ToXValue(env),
		"twitterid":  nil,
		"viber":      nil,
		"vk":         nil,
		"webchat":    nil,
		"wechat":     nil,
		"whatsapp":   nil,
	}), flows.ContextFunc(env, contact.URNs().MapContext))

	assert.Equal(t, 0, contact.Tickets().Count())

	mailgun := sa.Ticketers().Get("d605bb96-258d-4097-ad0a-080937db2212")
	weather := sa.Topics().Get("472a7a73-96cb-4736-b567-056d987cc5b4")
	ticket := flows.OpenTicket(mailgun, weather, "I have issues", nil)
	contact.Tickets().Add(ticket)

	assert.Equal(t, 1, contact.Tickets().Count())

	clone := contact.Clone()
	assert.Equal(t, "Joe Bloggs", clone.Name())
	assert.Equal(t, flows.ContactID(12345), clone.ID())
	assert.Equal(t, tz, clone.Timezone())
	assert.Equal(t, envs.Language("eng"), clone.Language())
	assert.Equal(t, android, contact.PreferredChannel())
	assert.Equal(t, 1, clone.Tickets().Count())

	// can also clone a null contact!
	mrNil := (*flows.Contact)(nil)
	assert.Nil(t, mrNil.Clone())

	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__":  types.NewXText("Joe Bloggs"),
		"channel":      flows.Context(env, android),
		"created_on":   types.NewXDateTime(contact.CreatedOn()),
		"last_seen_on": types.NewXDateTime(*contact.LastSeenOn()),
		"fields":       flows.Context(env, contact.Fields()),
		"first_name":   types.NewXText("Joe"),
		"groups":       contact.Groups().ToXValue(env),
		"id":           types.NewXText("12345"),
		"language":     types.NewXText("eng"),
		"name":         types.NewXText("Joe Bloggs"),
		"tickets":      contact.Tickets().ToXValue(env),
		"timezone":     types.NewXText("America/Bogota"),
		"urn":          contact.URNs()[0].ToXValue(env),
		"urns":         contact.URNs().ToXValue(env),
		"uuid":         types.NewXText(string(contact.UUID())),
	}), flows.Context(env, contact))

	assert.True(t, contact.ClearURNs()) // did have URNs
	assert.False(t, contact.ClearURNs())
	assert.Equal(t, flows.URNList{}, contact.URNs())

	marshaled, err := jsonx.Marshal(contact)
	require.NoError(t, err)

	fmt.Println(string(marshaled))

	unmarshaled, err := flows.ReadContact(sa, marshaled, assets.PanicOnMissing)
	require.NoError(t, err)

	assert.True(t, contact.Equal(unmarshaled))
}

func TestReadContact(t *testing.T) {
	source, err := static.NewSource([]byte(`{}`))
	require.NoError(t, err)

	env := envs.NewBuilder().Build()

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	// read minimal contact
	contact, err := flows.ReadContact(sa, []byte(`{"uuid": "a20f7948-e497-4a4a-be3c-b17f79f7ab7d", "created_on": "2020-07-22T13:50:30.123456789Z"}`), assets.PanicOnMissing)
	assert.NoError(t, err)
	assert.Equal(t, flows.ContactUUID("a20f7948-e497-4a4a-be3c-b17f79f7ab7d"), contact.UUID())
	assert.Equal(t, flows.ContactStatusActive, contact.Status())

	// read invalid contact
	_, err = flows.ReadContact(sa, []byte(`{"uuid": "a20f7948-e497-4a4a-be3c-b17f79f7ab7d", "status": "drunk", "created_on": "2020-07-22T13:50:30.123456789Z"}`), assets.PanicOnMissing)
	assert.EqualError(t, err, "unable to read contact: field 'status' is not a valid contact status")
}

func TestContactFormat(t *testing.T) {
	env := envs.NewBuilder().Build()
	sa, _ := engine.NewSessionAssets(env, static.NewEmptySource(), nil)

	// name takes precedence if set
	contact := flows.NewEmptyContact(sa, "Joe", envs.NilLanguage, nil)
	contact.AddURN(urns.URN("twitter:joey"), nil)
	assert.Equal(t, "Joe", contact.Format(env))

	// if not we fallback to URN
	contact, _ = flows.NewContact(
		sa,
		flows.ContactUUID(uuids.New()),
		flows.ContactID(1234),
		"",
		envs.NilLanguage,
		flows.ContactStatusActive,
		nil,
		time.Now(),
		nil,
		nil,
		nil,
		nil,
		nil,
		assets.PanicOnMissing,
	)
	contact.AddURN(urns.URN("twitter:joey"), nil)
	assert.Equal(t, "joey", contact.Format(env))

	anonEnv := envs.NewBuilder().WithRedactionPolicy(envs.RedactionPolicyURNs).Build()

	// unless URNs are redacted
	assert.Equal(t, "1234", contact.Format(anonEnv))

	// if we don't have name or URNs, then empty string
	contact = flows.NewEmptyContact(sa, "", envs.NilLanguage, nil)
	assert.Equal(t, "", contact.Format(env))
}

func TestContactSetPreferredChannel(t *testing.T) {
	env := envs.NewBuilder().Build()
	sa, _ := engine.NewSessionAssets(env, static.NewEmptySource(), nil)
	roles := []assets.ChannelRole{assets.ChannelRoleSend}
	receive_roles := []assets.ChannelRole{assets.ChannelRoleReceive}

	android := test.NewTelChannel("Android", "+250961111111", roles, nil, "RW", nil, false)
	android2 := test.NewTelChannel("Android", "+250961111112", receive_roles, nil, "RW", nil, false)
	twitter1 := test.NewChannel("Twitter", "nyaruka", []string{"twitter", "twitterid"}, roles, nil)
	twitter2 := test.NewChannel("Twitter", "nyaruka", []string{"twitter", "twitterid"}, roles, nil)

	contact := flows.NewEmptyContact(sa, "Joe", envs.NilLanguage, nil)
	contact.AddURN(urns.URN("twitter:joey"), nil)
	contact.AddURN(urns.URN("tel:+12345678999"), nil)
	contact.AddURN(urns.URN("tel:+18005555777"), nil)

	contact.UpdatePreferredChannel(android)

	// tel channels should be re-assigned to that channel, and moved to front of list
	assert.Equal(t, urns.URN("tel:+12345678999?channel="+string(android.UUID())), contact.URNs()[0].URN())
	assert.Equal(t, android, contact.URNs()[0].Channel())
	assert.Equal(t, urns.URN("tel:+18005555777?channel="+string(android.UUID())), contact.URNs()[1].URN())
	assert.Equal(t, android, contact.URNs()[1].Channel())
	assert.Equal(t, urns.URN("twitter:joey"), contact.URNs()[2].URN())
	assert.Nil(t, contact.URNs()[2].Channel())

	// same only applies to URNs of other schemes if they don't have a channel already
	contact.UpdatePreferredChannel(twitter1)
	assert.Equal(t, urns.URN("twitter:joey?channel="+string(twitter1.UUID())), contact.URNs()[0].URN())

	contact.UpdatePreferredChannel(twitter2)
	assert.Equal(t, urns.URN("twitter:joey?channel="+string(twitter1.UUID())), contact.URNs()[0].URN())

	// if they are already associated with the channel, then they become the preferred URN
	contact.UpdatePreferredChannel(android)
	contact.UpdatePreferredChannel(twitter1)

	assert.Equal(t, urns.URN("twitter:joey?channel="+string(twitter1.UUID())), contact.URNs()[0].URN())
	assert.Equal(t, twitter1, contact.URNs()[0].Channel())

	contact.UpdatePreferredChannel(android2)

	for _, urn := range contact.URNs() {
		assert.NotEqual(t, android2, urn.Channel())
	}

}

func TestReevaluateQueryBasedGroups(t *testing.T) {
	source, err := static.LoadSource("testdata/smart_groups.assets.json")
	require.NoError(t, err)

	tests := []struct {
		Description   string          `json:"description"`
		ContactBefore json.RawMessage `json:"contact_before"`
		RedactURNs    bool            `json:"redact_urns"`
		ContactAfter  json.RawMessage `json:"contact_after"`
	}{}

	testFile, err := os.ReadFile("testdata/smart_groups.json")
	require.NoError(t, err)
	err = jsonx.Unmarshal(testFile, &tests)
	require.NoError(t, err)

	for _, tc := range tests {
		envBuilder := envs.NewBuilder().
			WithAllowedLanguages([]envs.Language{"eng", "spa"}).
			WithDefaultCountry("RW")

		if tc.RedactURNs {
			envBuilder.WithRedactionPolicy(envs.RedactionPolicyURNs)
		}
		env := envBuilder.Build()

		sa, err := engine.NewSessionAssets(env, source, nil)
		require.NoError(t, err)

		contact, err := flows.ReadContact(sa, tc.ContactBefore, assets.IgnoreMissing)
		require.NoError(t, err)

		trigger := triggers.NewBuilder(
			env,
			assets.NewFlowReference("76f0a02f-3b75-4b86-9064-e9195e1b3a02", "Empty Flow"),
			contact,
		).Manual().Build()

		eng := engine.NewBuilder().Build()
		session, _, _ := eng.NewSession(sa, trigger)
		afterJSON := jsonx.MustMarshal(session.Contact())

		test.AssertEqualJSON(t, tc.ContactAfter, afterJSON, "contact JSON mismatch in '%s'", tc.Description)
	}
}

func TestContactEqual(t *testing.T) {
	session, _ := test.NewSessionBuilder().MustBuild()

	contact1JSON := []byte(`{
		"uuid": "ba96bf7f-bc2a-4873-a7c7-254d1927c4e3",
		"id": 1234567,
		"created_on": "2000-01-01T00:00:00.000000000-00:00",
		"fields": {
			"gender": {"text": "Male"}
		},
		"language": "eng",
		"name": "Ben Haggerty",
		"timezone": "America/Guayaquil",
		"urns": ["tel:+12065551212"]
	}`)

	contact1, err := flows.ReadContact(session.Assets(), contact1JSON, assets.PanicOnMissing)
	require.NoError(t, err)

	contact2, err := flows.ReadContact(session.Assets(), contact1JSON, assets.PanicOnMissing)
	require.NoError(t, err)

	assert.True(t, contact1.Equal(contact2))
	assert.True(t, contact2.Equal(contact1))
	assert.True(t, contact1.Equal(contact1.Clone()))
	assert.Equal(t, flows.ContactStatusActive, contact1.Status())

	// marshal and unmarshal contact 1 again
	contact1JSON, err = jsonx.Marshal(contact1)
	require.NoError(t, err)
	contact1, err = flows.ReadContact(session.Assets(), contact1JSON, assets.PanicOnMissing)
	require.NoError(t, err)

	assert.True(t, contact1.Equal(contact2))

	contact2.SetLanguage(envs.NilLanguage)
	assert.False(t, contact1.Equal(contact2))
}

func TestContactQuery(t *testing.T) {
	session, _ := test.NewSessionBuilder().MustBuild()

	contactJSON := []byte(`{
		"uuid": "ba96bf7f-bc2a-4873-a7c7-254d1927c4e3",
		"id": 1234567,
		"name": "Ben Haggerty",
		"fields": {
			"gender": {"text": "Male"},
			"age": {"text": "39!", "number": 39}
		},
		"groups": [
			{"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Testers"},
			{"uuid": "4f1f98fc-27a7-4a69-bbdb-24744ba739a9", "name": "Males"}
		],
		"tickets": [
			{
				"uuid": "e5f5a9b0-1c08-4e56-8f5c-92e00bc3cf52",
				"ticketer": {
					"uuid": "19dc6346-9623-4fe4-be80-538d493ecdf5",
					"name": "Support Tickets"
				},
				"subject": "Old ticket",
				"body": "I have a problem",
				"assignee": null
			}
		],
		"language": "eng",
		"timezone": "America/Guayaquil",
		"urns": [
			"tel:+12065551212", 
			"tel:+12065551313", 
			"twitter:ewok"
		],
		"created_on": "2020-01-24T13:24:30Z",
		"last_seen_on": "2020-08-06T15:41:30Z"
	}`)

	contact, err := flows.ReadContact(session.Assets(), contactJSON, assets.PanicOnMissing)
	require.NoError(t, err)

	testCases := []struct {
		query     string
		redaction envs.RedactionPolicy
		result    bool
		err       string
	}{
		{`name = "Ben Haggerty"`, envs.RedactionPolicyNone, true, ""},
		{`name = "Joe X"`, envs.RedactionPolicyNone, false, ""},
		{`name != "Joe X"`, envs.RedactionPolicyNone, true, ""},
		{`name != "Joe X"`, envs.RedactionPolicyNone, true, ""},
		{`name ~ Joe`, envs.RedactionPolicyNone, false, ""},
		{`name = ""`, envs.RedactionPolicyNone, false, ""},
		{`name != ""`, envs.RedactionPolicyNone, true, ""},

		{`uuid = ba96bf7f-bc2a-4873-a7c7-254d1927c4e3`, envs.RedactionPolicyNone, true, ""},
		{`uuid = 3bf7edda-b926-4a78-9131-d336df77d44f`, envs.RedactionPolicyNone, false, ""},

		{`id = 1234567`, envs.RedactionPolicyNone, true, ""},
		{`id = 5678889`, envs.RedactionPolicyNone, false, ""},

		{`language = ENG`, envs.RedactionPolicyNone, true, ""},
		{`language = FRA`, envs.RedactionPolicyNone, false, ""},
		{`language = ""`, envs.RedactionPolicyNone, false, ""},
		{`language != ""`, envs.RedactionPolicyNone, true, ""},

		{`created_on = 24-01-2020`, envs.RedactionPolicyNone, true, ""},
		{`created_on = 25-01-2020`, envs.RedactionPolicyNone, false, ""},
		{`created_on != 24-01-2020`, envs.RedactionPolicyNone, false, ""},
		{`created_on != 25-01-2020`, envs.RedactionPolicyNone, true, ""},
		{`created_on > 22-01-2020`, envs.RedactionPolicyNone, true, ""},
		{`created_on > 26-01-2020`, envs.RedactionPolicyNone, false, ""},

		{`last_seen_on = 06-08-2020`, envs.RedactionPolicyNone, true, ""},
		{`last_seen_on = 07-08-2020`, envs.RedactionPolicyNone, false, ""},
		{`last_seen_on > 05-08-2020`, envs.RedactionPolicyNone, true, ""},
		{`last_seen_on > 08-08-2020`, envs.RedactionPolicyNone, false, ""},
		{`last_seen_on != ""`, envs.RedactionPolicyNone, true, ""},
		{`last_seen_on = ""`, envs.RedactionPolicyNone, false, ""},

		{`tel = +12065551212`, envs.RedactionPolicyNone, true, ""},
		{`tel = +12065551313`, envs.RedactionPolicyNone, true, ""},
		{`tel = +13065551212`, envs.RedactionPolicyNone, false, ""},
		{`tel ~ 555`, envs.RedactionPolicyNone, true, ""},
		{`tel ~ 666`, envs.RedactionPolicyNone, false, ""},
		{`tel = ""`, envs.RedactionPolicyNone, false, ""},
		{`tel != ""`, envs.RedactionPolicyNone, true, ""},

		{`tel = +12065551212`, envs.RedactionPolicyURNs, false, "cannot query on redacted URNs"},
		{`tel ~ 555`, envs.RedactionPolicyURNs, false, "cannot query on redacted URNs"},
		{`tel = ""`, envs.RedactionPolicyURNs, false, ""},
		{`tel != ""`, envs.RedactionPolicyURNs, true, ""},

		{`twitter = ewok`, envs.RedactionPolicyNone, true, ""},
		{`twitter = nicp`, envs.RedactionPolicyNone, false, ""},
		{`twitter ~ wok`, envs.RedactionPolicyNone, true, ""},
		{`twitter ~ EWO`, envs.RedactionPolicyNone, true, ""},
		{`twitter ~ ijk`, envs.RedactionPolicyNone, false, ""},
		{`twitter = ""`, envs.RedactionPolicyNone, false, ""},
		{`twitter != ""`, envs.RedactionPolicyNone, true, ""},

		{`viber = ewok`, envs.RedactionPolicyNone, false, ""},
		{`viber ~ wok`, envs.RedactionPolicyNone, false, ""},
		{`viber = ""`, envs.RedactionPolicyNone, true, ""},
		{`viber != ""`, envs.RedactionPolicyNone, false, ""},

		{`urn = +12065551212`, envs.RedactionPolicyNone, true, ""},
		{`urn = ewok`, envs.RedactionPolicyNone, true, ""},
		{`urn = +13065551212`, envs.RedactionPolicyNone, false, ""},
		{`urn != +13065551212`, envs.RedactionPolicyNone, true, ""},
		{`urn ~ 555`, envs.RedactionPolicyNone, true, ""},
		{`urn ~ 666`, envs.RedactionPolicyNone, false, ""},
		{`urn = ""`, envs.RedactionPolicyNone, false, ""},
		{`urn != ""`, envs.RedactionPolicyNone, true, ""},

		{`urn = +12065551212`, envs.RedactionPolicyURNs, false, "cannot query on redacted URNs"},
		{`urn ~ 555`, envs.RedactionPolicyURNs, false, "cannot query on redacted URNs"},
		{`urn = ""`, envs.RedactionPolicyURNs, false, ""},
		{`urn != ""`, envs.RedactionPolicyURNs, true, ""},

		{`group = testers`, envs.RedactionPolicyNone, true, ""},
		{`group != testers`, envs.RedactionPolicyNone, false, ""},
		{`group = customers`, envs.RedactionPolicyNone, false, ""},
		{`group != customers`, envs.RedactionPolicyNone, true, ""},

		{`tickets = 1`, envs.RedactionPolicyNone, true, ""},
		{`tickets = 0`, envs.RedactionPolicyNone, false, ""},
		{`tickets != 1`, envs.RedactionPolicyNone, false, ""},
		{`tickets != 0`, envs.RedactionPolicyNone, true, ""},
		{`tickets > 0`, envs.RedactionPolicyNone, true, ""},

		{`age = 39`, envs.RedactionPolicyNone, true, ""},
		{`age != 39`, envs.RedactionPolicyNone, false, ""},
		{`age = 60`, envs.RedactionPolicyNone, false, ""},
		{`age != 60`, envs.RedactionPolicyNone, true, ""},
	}

	doQuery := func(q string, redaction envs.RedactionPolicy) (bool, error) {
		var env envs.Environment
		if redaction == envs.RedactionPolicyURNs {
			env = envs.NewBuilder().WithRedactionPolicy(envs.RedactionPolicyURNs).Build()
		} else {
			env = session.Environment()
		}

		parsed, err := contactql.ParseQuery(env, q, session.Assets())
		if err != nil {
			return false, err
		}

		return contactql.EvaluateQuery(env, parsed, contact), nil
	}

	for _, tc := range testCases {
		result, err := doQuery(tc.query, tc.redaction)

		if tc.err != "" {
			assert.EqualError(t, err, tc.err, "error mismatch evaluating '%s'", tc.query)
		} else {
			assert.NoError(t, err, "unexpected error evaluating '%s'", tc.query)
			assert.Equal(t, tc.result, result, "unexpected result for '%s'", tc.query)
		}
	}
}
