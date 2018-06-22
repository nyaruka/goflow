package flows_test

import (
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

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
	env := test.NewTestEnvironment(utils.DateFormatYearMonthDay, time.UTC, nil)

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
