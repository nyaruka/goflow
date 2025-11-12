package flows_test

import (
	"context"
	"encoding/json"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/i18n"
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

	uuids.SetGenerator(uuids.NewSeededGenerator(1234, time.Now))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	tz, _ := time.LoadLocation("America/Bogota")

	contact, err := flows.NewContact(
		sa,
		flows.NewContactUUID(),
		flows.ContactID(12345),
		"Joe Bloggs",
		i18n.Language("eng"),
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

	contact.AddURN(urns.URN("tel:+12024561111?channel=294a14d4-c998-41e5-a314-5941b97b89d7"))
	contact.AddURN(urns.URN("twitter:joey"))
	contact.AddURN(urns.URN("whatsapp:235423721788"))

	assert.Equal(t, "Joe Bloggs", contact.Name())
	assert.Equal(t, flows.ContactID(12345), contact.ID())
	assert.Equal(t, tz, contact.Timezone())
	assert.Equal(t, i18n.Language("eng"), contact.Language())
	assert.Equal(t, android, contact.PreferredChannel())
	assert.Equal(t, i18n.Country("US"), contact.Country())
	assert.Equal(t, i18n.Locale("eng-US"), contact.Locale(env))

	contact.SetStatus(flows.ContactStatusStopped)
	assert.Equal(t, flows.ContactStatusStopped, contact.Status())

	contact.SetStatus(flows.ContactStatusBlocked)
	assert.Equal(t, flows.ContactStatusBlocked, contact.Status())

	contact.SetStatus(flows.ContactStatusActive)
	assert.Equal(t, flows.ContactStatusActive, contact.Status())

	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"ext":        nil,
		"facebook":   nil,
		"fcm":        nil,
		"freshchat":  nil,
		"instagram":  nil,
		"jiochat":    nil,
		"line":       nil,
		"mailto":     nil,
		"rocketchat": nil,
		"slack":      nil,
		"tel":        flows.NewURN("tel", "+12024561111", "", android).ToXValue(env),
		"telegram":   nil,
		"twitter":    flows.NewURN("twitter", "joey", "", nil).ToXValue(env),
		"twitterid":  nil,
		"viber":      nil,
		"vk":         nil,
		"webchat":    nil,
		"wechat":     nil,
		"whatsapp":   flows.NewURN("whatsapp", "235423721788", "", nil).ToXValue(env),
	}), flows.ContextFunc(env, contact.URNs().MapContext))

	assert.Equal(t, 0, contact.Tickets().Open().Count())

	weather := sa.Topics().Get("472a7a73-96cb-4736-b567-056d987cc5b4")
	ticket := flows.OpenTicket(weather, nil)
	contact.Tickets().Add(ticket)

	assert.Equal(t, 1, contact.Tickets().Open().Count())

	clone := contact.Clone()
	assert.Equal(t, "Joe Bloggs", clone.Name())
	assert.Equal(t, flows.ContactID(12345), clone.ID())
	assert.Equal(t, tz, clone.Timezone())
	assert.Equal(t, i18n.Language("eng"), clone.Language())
	assert.Equal(t, i18n.Country("US"), clone.Country())
	assert.Equal(t, android, clone.PreferredChannel())
	assert.Equal(t, 0, clone.Tickets().Open().Count()) // not included

	// country can be resolved from tel urns if there's no preferred channel
	clone.UpdatePreferredChannel(nil)
	assert.Equal(t, i18n.Country("US"), clone.Country())

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
		"ref":          types.NewXText("A6YWQL"),
		"status":       types.NewXText(string(contact.Status())),
		"tickets":      contact.Tickets().ToXValue(env),
		"timezone":     types.NewXText("America/Bogota"),
		"urn":          contact.URNs()[0].ToXValue(env),
		"urns":         contact.URNs().ToXValue(env),
		"uuid":         types.NewXText(string(contact.UUID())),
	}), flows.Context(env, contact))

	marshaled, err := jsonx.Marshal(contact)
	require.NoError(t, err)

	unmarshaled, err := flows.ReadContact(sa, marshaled, assets.PanicOnMissing)
	require.NoError(t, err)

	assert.Equal(t, contact.UUID(), unmarshaled.UUID())
}

func TestContactURNs(t *testing.T) {
	source, err := static.NewSource([]byte(`{}`))
	require.NoError(t, err)

	env := envs.NewBuilder().Build()

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	contact := flows.NewEmptyContact(sa, "", i18n.NilLanguage, nil)

	assert.Len(t, contact.URNs(), 0)
	assert.True(t, contact.AddURN("tel:+12024561111"))  // didn't have URN so returns true
	assert.False(t, contact.AddURN("tel:+12024561111")) // did have

	assert.Equal(t, flows.URNList{flows.NewURN("tel", "+12024561111", "", nil)}, contact.URNs())
	assert.True(t, contact.AddURN("tel:+12024562222"))
	assert.Equal(t, flows.URNList{flows.NewURN("tel", "+12024561111", "", nil), flows.NewURN("tel", "+12024562222", "", nil)}, contact.URNs())
	assert.False(t, contact.SetURNs([]urns.URN{"tel:+12024561111", "tel:+12024562222"})) // no change
	assert.Equal(t, flows.URNList{flows.NewURN("tel", "+12024561111", "", nil), flows.NewURN("tel", "+12024562222", "", nil)}, contact.URNs())
	assert.True(t, contact.SetURNs([]urns.URN{"tel:+12024562222", "tel:+12024561111"})) // order changed
	assert.Equal(t, flows.URNList{flows.NewURN("tel", "+12024562222", "", nil), flows.NewURN("tel", "+12024561111", "", nil)}, contact.URNs())
	assert.True(t, contact.SetURNs([]urns.URN{"tel:+12024562222", "tel:+12024561111", "tel:+12024563333"}))
	assert.Equal(t, flows.URNList{flows.NewURN("tel", "+12024562222", "", nil), flows.NewURN("tel", "+12024561111", "", nil), flows.NewURN("tel", "+12024563333", "", nil)}, contact.URNs())
	assert.True(t, contact.RemoveURN("tel:+12024561111"))
	assert.False(t, contact.RemoveURN("tel:+12024566666"))
	assert.Equal(t, flows.URNList{flows.NewURN("tel", "+12024562222", "", nil), flows.NewURN("tel", "+12024563333", "", nil)}, contact.URNs())
}

func TestReadContact(t *testing.T) {
	source, err := static.NewSource([]byte(`{}`))
	require.NoError(t, err)

	env := envs.NewBuilder().Build()

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	// read minimal contact
	contact, err := flows.ReadContact(sa, []byte(`{"uuid": "a20f7948-e497-4a4a-be3c-b17f79f7ab7d", "status": "active", "created_on": "2020-07-22T13:50:30.123456789Z"}`), assets.PanicOnMissing)
	assert.NoError(t, err)
	assert.Equal(t, flows.ContactUUID("a20f7948-e497-4a4a-be3c-b17f79f7ab7d"), contact.UUID())
	assert.Equal(t, flows.ContactStatusActive, contact.Status())

	// read invalid contact
	_, err = flows.ReadContact(sa, []byte(`{"uuid": "a20f7948-e497-4a4a-be3c-b17f79f7ab7d", "status": "drunk", "created_on": "2020-07-22T13:50:30.123456789Z"}`), assets.PanicOnMissing)
	assert.EqualError(t, err, "unable to read contact: field 'status' is not a valid contact status")
}

func TestReadContactWithMissingAssets(t *testing.T) {
	sessionAssets, err := engine.NewSessionAssets(envs.NewBuilder().Build(), static.NewEmptySource(), nil)
	require.NoError(t, err)

	missingAssets := make([]assets.Reference, 0)
	missing := func(a assets.Reference, err error) { missingAssets = append(missingAssets, a) }

	flows.ReadContact(sessionAssets, []byte(`{
		"uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f",
		"id": 1234567,
		"name": "Ryan Lewis",
		"status": "active",
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
				"uuid": "78d1fe0d-7e39-461e-81c3-a6a25f15ed69",
				"status": "open",
				"topic": {
					"uuid": "472a7a73-96cb-4736-b567-056d987cc5b4",
					"name": "Weather"
				},
				"assignee": {"uuid": "0c78ef47-7d56-44d8-8f57-96e0f30e8f44", "name": "Bob"}
			}
		]
	}`), missing)

	refs := make([]string, len(missingAssets))
	for i := range missingAssets {
		refs[i] = missingAssets[i].String()
	}

	// ordering isn't deterministic so sort A-Z
	sort.Strings(refs)

	assert.Equal(t, []string{
		"channel[uuid=57f1078f-88aa-46f4-a59a-948a5739c03d,name=]",
		"field[key=activation_token,name=]",
		"field[key=gender,name=]",
		"field[key=join_date,name=]",
		"group[uuid=4f1f98fc-27a7-4a69-bbdb-24744ba739a9,name=Males]",
		"group[uuid=b7cf0d83-f1c9-411c-96fd-c511a4cfa86d,name=Testers]",
		"topic[uuid=472a7a73-96cb-4736-b567-056d987cc5b4,name=Weather]",
		"user[uuid=0c78ef47-7d56-44d8-8f57-96e0f30e8f44,name=Bob]",
	}, refs)
}

func TestContactFormat(t *testing.T) {
	env := envs.NewBuilder().Build()
	sa, _ := engine.NewSessionAssets(env, static.NewEmptySource(), nil)

	// name takes precedence if set
	contact := flows.NewEmptyContact(sa, "Joe", i18n.NilLanguage, nil)
	contact.AddURN(urns.URN("twitter:joey"))
	assert.Equal(t, "Joe", contact.Format(env))

	// if not we fallback to URN
	contact, _ = flows.NewContact(
		sa,
		flows.NewContactUUID(),
		flows.ContactID(1234),
		"",
		i18n.NilLanguage,
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
	contact.AddURN(urns.URN("twitter:joey"))
	assert.Equal(t, "joey", contact.Format(env))

	anonEnv := envs.NewBuilder().WithRedactionPolicy(envs.RedactionPolicyURNs).Build()

	// unless URNs are redacted
	assert.Equal(t, "1234", contact.Format(anonEnv))

	// if we don't have name or URNs, then empty string
	contact = flows.NewEmptyContact(sa, "", i18n.NilLanguage, nil)
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
	whatsapp1 := test.NewChannel("Whatsapp", "+250961111113", []string{"whatsapp"}, roles, nil)
	whatsapp2 := test.NewChannel("Whatsapp", "+250961111114", []string{"whatsapp"}, roles, nil)

	contact := flows.NewEmptyContact(sa, "Joe", i18n.NilLanguage, nil)
	contact.AddURN(urns.URN("twitter:joey"))
	contact.AddURN(urns.URN("tel:+12345678999"))
	contact.AddURN(urns.URN("tel:+18005555777"))
	contact.AddURN(urns.URN("whatsapp:18005555888"))

	contact.UpdatePreferredChannel(android)

	// tel channels should be re-assigned to that channel, and moved to front of list
	assert.Equal(t, urns.URN("tel:+12345678999?channel="+string(android.UUID())), contact.URNs()[0].Encode())
	assert.Equal(t, android, contact.URNs()[0].Channel)
	assert.Equal(t, urns.URN("tel:+18005555777?channel="+string(android.UUID())), contact.URNs()[1].Encode())
	assert.Equal(t, android, contact.URNs()[1].Channel)
	assert.Equal(t, urns.URN("twitter:joey"), contact.URNs()[2].Encode())
	assert.Nil(t, contact.URNs()[2].Channel)

	// same only applies to URNs of other schemes if they don't have a channel already
	contact.UpdatePreferredChannel(twitter1)
	assert.Equal(t, urns.URN("twitter:joey?channel="+string(twitter1.UUID())), contact.URNs()[0].Encode())

	contact.UpdatePreferredChannel(twitter2)
	assert.Equal(t, urns.URN("twitter:joey?channel="+string(twitter1.UUID())), contact.URNs()[0].Encode())

	contact.UpdatePreferredChannel(whatsapp1)
	assert.Equal(t, urns.URN("whatsapp:18005555888?channel="+string(whatsapp1.UUID())), contact.URNs()[0].Encode())

	contact.UpdatePreferredChannel(whatsapp2)
	assert.Equal(t, urns.URN("whatsapp:18005555888?channel="+string(whatsapp2.UUID())), contact.URNs()[0].Encode())

	// if they are already associated with the channel, then they become the preferred URN
	contact.UpdatePreferredChannel(android)
	contact.UpdatePreferredChannel(twitter1)

	assert.Equal(t, urns.URN("twitter:joey?channel="+string(twitter1.UUID())), contact.URNs()[0].Encode())
	assert.Equal(t, twitter1, contact.URNs()[0].Channel)

	contact.UpdatePreferredChannel(android2)

	for _, urn := range contact.URNs() {
		assert.NotEqual(t, android2, urn.Channel)
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
	jsonx.MustUnmarshal(testFile, &tests)

	for _, tc := range tests {
		envBuilder := envs.NewBuilder().
			WithAllowedLanguages("eng", "spa").
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
			assets.NewFlowReference("76f0a02f-3b75-4b86-9064-e9195e1b3a02", "Empty Flow"),
		).Manual().Build()

		eng := engine.NewBuilder().Build()
		session, _, _ := eng.NewSession(context.Background(), sa, env, contact, trigger, nil)
		afterJSON := jsonx.MustMarshal(session.Contact())

		test.AssertEqualJSON(t, tc.ContactAfter, afterJSON, "contact JSON mismatch in '%s'", tc.Description)
	}
}

func TestContactQuery(t *testing.T) {
	_, session, _ := test.NewSessionBuilder().MustBuild()

	contactJSON := []byte(`{
		"uuid": "ba96bf7f-bc2a-4873-a7c7-254d1927c4e3",
		"id": 1234567,
		"name": "Ben Haggerty",
		"status": "active",
		"fields": {
			"gender": {"text": "Male"},
			"age": {"text": "39!", "number": 39},
			"language": {"text": "en"}
		},
		"groups": [
			{"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Testers"},
			{"uuid": "4f1f98fc-27a7-4a69-bbdb-24744ba739a9", "name": "Males"}
		],
		"tickets": [
			{
				"uuid": "e5f5a9b0-1c08-4e56-8f5c-92e00bc3cf52",
				"status": "open",
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

		{`tickets = 1`, envs.RedactionPolicyNone, true, ""},
		{`tickets = 0`, envs.RedactionPolicyNone, false, ""},
		{`tickets != 1`, envs.RedactionPolicyNone, false, ""},
		{`tickets != 0`, envs.RedactionPolicyNone, true, ""},
		{`tickets > 0`, envs.RedactionPolicyNone, true, ""},

		{`age = 39`, envs.RedactionPolicyNone, true, ""},
		{`age != 39`, envs.RedactionPolicyNone, false, ""},
		{`age = 60`, envs.RedactionPolicyNone, false, ""},
		{`age != 60`, envs.RedactionPolicyNone, true, ""},

		// field with key that conflicts with attribute has to be prefixed
		{`fields.language = EN`, envs.RedactionPolicyNone, true, ""},
		{`fields.language = FR`, envs.RedactionPolicyNone, false, ""},

		// check querying on a field that isn't set for this contact
		{`activation_token = ""`, envs.RedactionPolicyNone, true, ""},
		{`activation_token != ""`, envs.RedactionPolicyNone, false, ""},
		{`activation_token = "xx"`, envs.RedactionPolicyNone, false, ""},
		{`activation_token != "xx"`, envs.RedactionPolicyNone, true, ""},
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
