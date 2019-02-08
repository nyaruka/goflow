package flows_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContact(t *testing.T) {
	env := utils.NewDefaultEnvironment()

	utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(1234))
	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)

	contact := flows.NewContact(
		flows.ContactUUID(utils.NewUUID()), flows.ContactID(12345), "Joe Bloggs", utils.Language("eng"),
		nil, time.Now(), flows.URNList{}, flows.NewGroupList([]*flows.Group{}), make(flows.FieldValues),
	)

	assert.Equal(t, flows.URNList{}, contact.URNs())
	assert.Nil(t, contact.PreferredChannel())

	contact.SetTimezone(env.Timezone())
	contact.SetCreatedOn(time.Date(2017, 12, 15, 10, 0, 0, 0, time.UTC))
	contact.AddURN(flows.NewContactURN(urns.URN("tel:+16364646466?channel=294a14d4-c998-41e5-a314-5941b97b89d7"), nil))
	contact.AddURN(flows.NewContactURN(urns.URN("twitter:joey"), nil))

	assert.Equal(t, "Joe Bloggs", contact.Name())
	assert.Equal(t, flows.ContactID(12345), contact.ID())
	assert.Equal(t, env.Timezone(), contact.Timezone())
	assert.Equal(t, utils.Language("eng"), contact.Language())
	assert.Nil(t, contact.PreferredChannel())
	assert.True(t, contact.HasURN("tel:+16364646466"))
	assert.False(t, contact.HasURN("tel:+16300000000"))

	clone := contact.Clone()
	assert.Equal(t, "Joe Bloggs", clone.Name())
	assert.Equal(t, flows.ContactID(12345), clone.ID())
	assert.Equal(t, env.Timezone(), clone.Timezone())
	assert.Equal(t, utils.Language("eng"), clone.Language())
	assert.Nil(t, contact.PreferredChannel())

	// can also clone a null contact!
	mrNil := (*flows.Contact)(nil)
	assert.Nil(t, mrNil.Clone())

	assert.Equal(t, types.NewXText(string(contact.UUID())), contact.Resolve(env, "uuid"))
	assert.Equal(t, types.NewXNumberFromInt(12345), contact.Resolve(env, "id"))
	assert.Equal(t, types.NewXText("Joe Bloggs"), contact.Resolve(env, "name"))
	assert.Equal(t, types.NewXText("Joe"), contact.Resolve(env, "first_name"))
	assert.Equal(t, types.NewXDateTime(contact.CreatedOn()), contact.Resolve(env, "created_on"))
	assert.Equal(t, contact.URNs(), contact.Resolve(env, "urns"))
	assert.Equal(t, contact.URNs()[0], contact.Resolve(env, "urn"))
	assert.Equal(t, contact.Fields(), contact.Resolve(env, "fields"))
	assert.Equal(t, contact.Groups(), contact.Resolve(env, "groups"))
	assert.Nil(t, contact.Resolve(env, "channel"))
	assert.Equal(t, types.NewXResolveError(contact, "xxx"), contact.Resolve(env, "xxx"))
	assert.Equal(t, types.NewXText("Joe Bloggs"), contact.Reduce(env))
	assert.Equal(t, "contact", contact.Describe())
	assert.Equal(t, types.NewXText(`{"channel":null,"created_on":"2017-12-15T10:00:00.000000Z","fields":{},"groups":[],"language":"eng","name":"Joe Bloggs","timezone":"UTC","urns":[{"display":"(636) 464-6466","path":"+16364646466","scheme":"tel"},{"display":"joey","path":"joey","scheme":"twitter"}],"uuid":"c00e5d67-c275-4389-aded-7d8b151cbd5b"}`), contact.ToXJSON(env))
}

func TestContactFormat(t *testing.T) {
	env := utils.NewEnvironment(utils.DateFormatYearMonthDay, utils.TimeFormatHourMinute, time.UTC, utils.NilLanguage, nil, utils.NilCountry, utils.DefaultNumberFormat, utils.RedactionPolicyNone, 640)

	// name takes precedence if set
	contact := flows.NewEmptyContact("Joe", utils.NilLanguage, nil)
	contact.AddURN(flows.NewContactURN(urns.URN("twitter:joey"), nil))
	assert.Equal(t, "Joe", contact.Format(env))

	// if not we fallback to URN
	contact = flows.NewContact(
		flows.ContactUUID(utils.NewUUID()), flows.ContactID(1234), "", utils.NilLanguage, nil, time.Now(),
		flows.URNList{}, flows.NewGroupList([]*flows.Group{}), make(flows.FieldValues),
	)
	contact.AddURN(flows.NewContactURN(urns.URN("twitter:joey"), nil))
	assert.Equal(t, "joey", contact.Format(env))

	anonEnv := utils.NewEnvironment(utils.DateFormatYearMonthDay, utils.TimeFormatHourMinute, time.UTC, utils.NilLanguage, nil, utils.NilCountry, utils.DefaultNumberFormat, utils.RedactionPolicyURNs, 640)

	// unless URNs are redacted
	assert.Equal(t, "1234", contact.Format(anonEnv))

	// if we don't have name or URNs, then empty string
	contact = flows.NewEmptyContact("", utils.NilLanguage, nil)
	assert.Equal(t, "", contact.Format(env))
}

func TestContactSetPreferredChannel(t *testing.T) {
	roles := []assets.ChannelRole{assets.ChannelRoleSend}

	android := test.NewTelChannel("Android", "+250961111111", roles, nil, "RW", nil)
	twitter := test.NewChannel("Twitter", "nyaruka", []string{"twitter", "twitterid"}, roles, nil)

	contact := flows.NewEmptyContact("Joe", utils.NilLanguage, nil)
	contact.AddURN(flows.NewContactURN(urns.URN("twitter:joey"), nil))
	contact.AddURN(flows.NewContactURN(urns.URN("tel:+12345678999"), nil))
	contact.AddURN(flows.NewContactURN(urns.URN("tel:+18005555777"), nil))

	contact.UpdatePreferredChannel(android)

	// tel channels should be re-assigned to that channel, and moved to front of list
	assert.Equal(t, urns.URN("tel:+12345678999?channel="+string(android.UUID())), contact.URNs()[0].URN())
	assert.Equal(t, android, contact.URNs()[0].Channel())
	assert.Equal(t, urns.URN("tel:+18005555777?channel="+string(android.UUID())), contact.URNs()[1].URN())
	assert.Equal(t, android, contact.URNs()[1].Channel())
	assert.Equal(t, urns.URN("twitter:joey"), contact.URNs()[2].URN())
	assert.Nil(t, contact.URNs()[2].Channel())

	contact.UpdatePreferredChannel(twitter)

	// same doesn't apply to URNs of other schemes
	assert.Equal(t, urns.URN("twitter:joey"), contact.URNs()[2].URN())
	assert.Nil(t, contact.URNs()[2].Channel())

	// unless they are already associated with that channel
	contact.URNs()[2].SetChannel(twitter)
	contact.UpdatePreferredChannel(twitter)

	assert.Equal(t, urns.URN("twitter:joey?channel="+string(twitter.UUID())), contact.URNs()[0].URN())
	assert.Equal(t, twitter, contact.URNs()[0].Channel())
}

func TestReevaluateDynamicGroups(t *testing.T) {
	session, _, err := test.CreateTestSession("http://localhost", nil)
	require.NoError(t, err)

	env := session.Runs()[0].Environment()

	gender := test.NewField("gender", "Gender", assets.FieldTypeText)
	age := test.NewField("age", "Age", assets.FieldTypeNumber)

	fieldSet := flows.NewFieldAssets([]assets.Field{gender.Asset(), age.Asset()})

	males := test.NewGroup("Males", `gender="M"`)
	old := test.NewGroup("Old", `age>30`)
	english := test.NewGroup("English", `language=eng`)
	spanish := test.NewGroup("Español", `language=spa`)
	lastYear := test.NewGroup("Old", `created_on <= 2017-12-31`)
	tel1800 := test.NewGroup("Tel with 1800", `tel ~ 1800`)
	twitterCrazies := test.NewGroup("Twitter Crazies", `twitter ~ crazy`)
	groups := []*flows.Group{males, old, english, spanish, lastYear, tel1800, twitterCrazies}

	contact := flows.NewEmptyContact("Joe", "eng", nil)
	contact.AddURN(flows.NewContactURN(urns.URN("tel:+12345678999"), nil))

	assert.Equal(t, []*flows.Group{english}, evaluateGroups(t, env, contact, groups))

	contact.SetLanguage(utils.Language("spa"))
	contact.AddURN(flows.NewContactURN(urns.URN("twitter:crazy_joe"), nil))
	contact.AddURN(flows.NewContactURN(urns.URN("tel:+18005555777"), nil))

	genderValue := contact.Fields().Parse(env, fieldSet, gender, "M")
	contact.Fields().Set(gender, genderValue)

	ageValue := contact.Fields().Parse(env, fieldSet, age, "37")
	contact.Fields().Set(age, ageValue)

	contact.SetCreatedOn(time.Date(2017, 12, 15, 10, 0, 0, 0, time.UTC))

	assert.Equal(t, []*flows.Group{males, old, spanish, lastYear, tel1800, twitterCrazies}, evaluateGroups(t, env, contact, groups))
}

func evaluateGroups(t *testing.T, env utils.Environment, contact *flows.Contact, groups []*flows.Group) []*flows.Group {
	matching := make([]*flows.Group, 0)
	for _, group := range groups {
		isMember, err := group.CheckDynamicMembership(env, contact)
		assert.NoError(t, err)
		if isMember {
			matching = append(matching, group)
		}
	}
	return matching
}

func TestContactEqual(t *testing.T) {
	session, _, err := test.CreateTestSession("http://localhost", nil)
	require.NoError(t, err)

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

	// marshal and unmarshal contact 1 again
	contact1JSON, err = json.Marshal(contact1)
	require.NoError(t, err)
	contact1, err = flows.ReadContact(session.Assets(), contact1JSON, assets.PanicOnMissing)
	require.NoError(t, err)

	assert.True(t, contact1.Equal(contact2))

	contact2.SetLanguage(utils.NilLanguage)
	assert.False(t, contact1.Equal(contact2))
}
