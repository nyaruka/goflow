package flows_test

import (
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
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

	contact := flows.NewContact("Joe Bloggs", utils.Language("eng"), nil)
	contact.SetTimezone(env.Timezone())
	contact.SetID(12345)
	contact.SetCreatedOn(time.Date(2017, 12, 15, 10, 0, 0, 0, time.UTC))
	contact.AddURN(urns.URN("tel:+16364646466"))
	contact.AddURN(urns.URN("twitter:joey"))

	assert.Equal(t, "Joe Bloggs", contact.Name())
	assert.Equal(t, 12345, contact.ID())
	assert.Equal(t, env.Timezone(), contact.Timezone())
	assert.Equal(t, utils.Language("eng"), contact.Language())

	clone := contact.Clone()
	assert.Equal(t, "Joe Bloggs", clone.Name())
	assert.Equal(t, 12345, clone.ID())
	assert.Equal(t, env.Timezone(), clone.Timezone())
	assert.Equal(t, utils.Language("eng"), clone.Language())

	// can also clone a null contact!
	mrNil := (*flows.Contact)(nil)
	assert.Nil(t, mrNil.Clone())

	assert.Equal(t, types.NewXText(string(contact.UUID())), contact.Resolve(env, "uuid"))
	assert.Equal(t, types.NewXNumberFromInt(12345), contact.Resolve(env, "id"))
	assert.Equal(t, types.NewXText("Joe Bloggs"), contact.Resolve(env, "name"))
	assert.Equal(t, types.NewXText("Joe"), contact.Resolve(env, "first_name"))
	assert.Equal(t, types.NewXDateTime(contact.CreatedOn()), contact.Resolve(env, "created_on"))
	assert.Equal(t, contact.URNs(), contact.Resolve(env, "urns"))
	assert.Equal(t, types.NewXText("(636) 464-6466"), contact.Resolve(env, "urn"))
	assert.Equal(t, contact.Fields(), contact.Resolve(env, "fields"))
	assert.Equal(t, contact.Groups(), contact.Resolve(env, "groups"))
	assert.Nil(t, contact.Resolve(env, "channel"))
	assert.Equal(t, types.NewXResolveError(contact, "xxx"), contact.Resolve(env, "xxx"))
	assert.Equal(t, types.NewXText("Joe Bloggs"), contact.Reduce(env))
	assert.Equal(t, "contact", contact.Describe())
	assert.Equal(t, types.NewXText(`{"channel":null,"created_on":"2017-12-15T10:00:00.000000Z","fields":{},"groups":[],"language":"eng","name":"Joe Bloggs","timezone":"UTC","urns":[{"display":"","path":"+16364646466","scheme":"tel"},{"display":"","path":"joey","scheme":"twitter"}],"uuid":"c00e5d67-c275-4389-aded-7d8b151cbd5b"}`), contact.ToXJSON(env))
}

func TestContactFormat(t *testing.T) {
	env := utils.NewEnvironment(utils.DateFormatYearMonthDay, utils.TimeFormatHourMinute, time.UTC, nil, utils.RedactionPolicyNone)

	// name takes precedence if set
	contact := flows.NewContact("Joe", utils.NilLanguage, nil)
	contact.AddURN(urns.URN("twitter:joey"))
	assert.Equal(t, "Joe", contact.Format(env))

	// if not we fallback to URN
	contact = flows.NewContact("", utils.NilLanguage, nil)
	contact.AddURN(urns.URN("twitter:joey"))
	contact.SetID(12345)
	assert.Equal(t, "joey", contact.Format(env))

	anonEnv := utils.NewEnvironment(utils.DateFormatYearMonthDay, utils.TimeFormatHourMinute, time.UTC, nil, utils.RedactionPolicyURNs)

	// unless URNs are redacted
	assert.Equal(t, "12345", contact.Format(anonEnv))

	// if we don't have name or URNs, then empty string
	contact = flows.NewContact("", utils.NilLanguage, nil)
	assert.Equal(t, "", contact.Format(env))
}

func TestContactSetPreferredChannel(t *testing.T) {
	roles := []flows.ChannelRole{flows.ChannelRoleSend}

	android := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Android", "+250961111111", []string{"tel"}, roles)
	twitter := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Twitter", "nyaruka", []string{"twitter", "twitterid"}, roles)
	//nexmo := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Nexmo", "+250961111111", []string{"tel"}, roles)

	contact := flows.NewContact("Joe", utils.NilLanguage, nil)
	contact.AddURN(urns.URN("twitter:joey"))
	contact.AddURN(urns.URN("tel:+12345678999"))
	contact.AddURN(urns.URN("tel:+18005555777"))

	contact.UpdatePreferredChannel(android)

	// tel channels should be re-assigned to that channel, and moved to front of list
	assert.Equal(t, urns.URN("tel:+12345678999"), contact.URNs()[0].URN)
	assert.Equal(t, android, contact.URNs()[0].Channel())
	assert.Equal(t, urns.URN("tel:+18005555777"), contact.URNs()[1].URN)
	assert.Equal(t, android, contact.URNs()[1].Channel())
	assert.Equal(t, urns.URN("twitter:joey"), contact.URNs()[2].URN)
	assert.Nil(t, contact.URNs()[2].Channel())

	contact.UpdatePreferredChannel(twitter)

	// same doesn't apply to URNs of other schemes
	assert.Equal(t, urns.URN("twitter:joey"), contact.URNs()[2].URN)
	assert.Nil(t, contact.URNs()[2].Channel())

	// unless they are already associated with that channel
	contact.URNs()[2].SetChannel(twitter)
	contact.UpdatePreferredChannel(twitter)

	assert.Equal(t, urns.URN("twitter:joey"), contact.URNs()[0].URN)
	assert.Equal(t, twitter, contact.URNs()[0].Channel())
}

func TestReevaluateDynamicGroups(t *testing.T) {
	session, err := test.CreateTestSession("http://localhost", nil)
	require.NoError(t, err)

	env := session.Runs()[0].Environment()

	fieldSet := flows.NewFieldSet([]*flows.Field{
		flows.NewField("gender", "Gender", flows.FieldValueTypeText),
		flows.NewField("age", "Age", flows.FieldValueTypeNumber),
	})

	males := flows.NewGroup(flows.GroupUUID(utils.NewUUID()), "Males", `gender="M"`)
	old := flows.NewGroup(flows.GroupUUID(utils.NewUUID()), "Old", `age>30`)
	english := flows.NewGroup(flows.GroupUUID(utils.NewUUID()), "English", `language=eng`)
	spanish := flows.NewGroup(flows.GroupUUID(utils.NewUUID()), "Español", `language=spa`)
	lastYear := flows.NewGroup(flows.GroupUUID(utils.NewUUID()), "Old", `created_on <= 2017-12-31`)
	tel1800 := flows.NewGroup(flows.GroupUUID(utils.NewUUID()), "Tel with 1800", `tel ~ 1800`)
	twitterCrazies := flows.NewGroup(flows.GroupUUID(utils.NewUUID()), "Twitter Crazies", `twitter ~ crazy`)
	groups := []*flows.Group{males, old, english, spanish, lastYear, tel1800, twitterCrazies}

	contact := flows.NewContact("Joe", "eng", nil)
	contact.AddURN(urns.URN("tel:+12345678999"))

	assert.Equal(t, []*flows.Group{english}, evaluateGroups(t, env, contact, groups))

	contact.SetLanguage(utils.Language("spa"))
	contact.AddURN(urns.URN("twitter:crazy_joe"))
	contact.AddURN(urns.URN("tel:+18005555777"))
	contact.SetFieldValue(env, fieldSet, "gender", "M")
	contact.SetFieldValue(env, fieldSet, "age", "37")
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
