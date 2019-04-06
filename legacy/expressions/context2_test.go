package expressions_test

import (
	"testing"

	"github.com/nyaruka/goflow/legacy/expressions"

	"github.com/stretchr/testify/assert"
)

var testCases = []struct {
	old string
	new string
}{
	{old: `contact`, new: `contact`},
	{old: `CONTACT`, new: `contact`},
	{old: `contact.uuid`, new: `contact.uuid`},
	{old: `contact.id`, new: `contact.id`},
	{old: `contact.name`, new: `contact.name`},
	{old: `contact.NAME`, new: `contact.name`},
	{old: `contact.first_name`, new: `contact.first_name`},
	{old: `contact.gender`, new: `fields.gender`},
	{old: `contact.groups`, new: `join(contact.groups, ",")`},
	{old: `contact.language`, new: `contact.language`},
	{old: `contact.created_on`, new: `contact.created_on`},

	// contact URN variables
	{old: `contact.tel`, new: `format_urn(urns.tel)`},
	{old: `contact.tel.display`, new: `format_urn(urns.tel)`},
	{old: `contact.tel.scheme`, new: `urn_parts(urns.tel).scheme`},
	{old: `contact.tel.path`, new: `urn_parts(urns.tel).path`},
	{old: `contact.tel.urn`, new: `urns.tel`},
	{old: `contact.tel_e164`, new: `urn_parts(urns.tel).path`},
	{old: `contact.twitterid`, new: `format_urn(urns.twitterid)`},
	{old: `contact.mailto`, new: `format_urn(urns.mailto)`},

	// run variables
	{old: `flow`, new: `results`},
	{old: `flow.favorite_color`, new: `results.favorite_color`},
	{old: `flow.favorite_color.category`, new: `results.favorite_color.category_localized`},
	{old: `flow.favorite_color.text`, new: `results.favorite_color.input`},
	{old: `flow.favorite_color.time`, new: `results.favorite_color.created_on`},
	{old: `flow.favorite_color.value`, new: `results.favorite_color.value`},
	{old: `flow.2factor`, new: `results["2factor"]`},
	{old: `flow.2factor.value`, new: `results["2factor"].value`},
	{old: `flow.1`, new: `results["1"]`},
	{old: `flow.1337`, new: `results["1337"]`},
	{old: `flow.1337.category`, new: `results["1337"].category_localized`},
	{old: `flow.contact`, new: `contact`},
	{old: `flow.contact.name`, new: `contact.name`},

	{old: `child.age`, new: `child.results.age`},
	{old: `child.age.category`, new: `child.results.age.category_localized`},
	{old: `child.age.text`, new: `child.results.age.input`},
	{old: `child.age.time`, new: `child.results.age.created_on`},
	{old: `child.age.value`, new: `child.results.age.value`},
	{old: `child.contact`, new: `contact`},
	{old: `child.contact.age`, new: `fields.age`},

	{old: `parent.role`, new: `parent.results.role`},
	{old: `parent.role.category`, new: `parent.results.role.category_localized`},
	{old: `parent.role.text`, new: `parent.results.role.input`},
	{old: `parent.role.time`, new: `parent.results.role.created_on`},
	{old: `parent.role.value`, new: `parent.results.role.value`},
	{old: `parent.contact`, new: `parent.contact`},
	{old: `parent.contact.name`, new: `parent.contact.name`},
	{old: `parent.contact.groups`, new: `join(parent.contact.groups, ",")`},
	{old: `parent.contact.gender`, new: `parent.fields.gender`},
	{old: `parent.contact.tel`, new: `format_urn(parent.urns.tel)`},
	{old: `parent.contact.tel.display`, new: `format_urn(parent.urns.tel)`},
	{old: `parent.contact.tel.scheme`, new: `urn_parts(parent.urns.tel).scheme`},
	{old: `parent.contact.tel.path`, new: `urn_parts(parent.urns.tel).path`},
	{old: `parent.contact.tel.urn`, new: `parent.urns.tel`},
	{old: `parent.contact.tel_e164`, new: `urn_parts(parent.urns.tel).path`},

	{old: `step`, new: `format_input(input)`},
	{old: `step.value`, new: `format_input(input)`},
	{old: `step.text`, new: `input.text`},
	{old: `step.attachments`, new: `foreach(foreach(input.attachments, attachment_parts), extract, "url")`},
	{old: `step.attachments.0`, new: `attachment_parts(input.attachments[0]).url`},
	{old: `step.attachments.10`, new: `attachment_parts(input.attachments[10]).url`},
	{old: `step.time`, new: `input.created_on`},
	{old: `step.contact`, new: `contact`},
	{old: `step.contact.name`, new: `contact.name`},
	{old: `step.contact.age`, new: `fields.age`},

	{old: `channel`, new: `contact.channel.address`},
	{old: `channel.tel`, new: `contact.channel.address`},
	{old: `channel.tel_e164`, new: `contact.channel.address`},
	{old: `channel.name`, new: `contact.channel.name`},

	{old: `date`, new: `now()`},
	{old: `date.now`, new: `now()`},
	{old: `date.today`, new: `format_date(today())`},
	{old: `date.tomorrow`, new: `format_date(datetime_add(now(), 1, "D"))`},
	{old: `date.yesterday`, new: `format_date(datetime_add(now(), -1, "D"))`},

	{old: `extra`, new: `legacy_extra`},
	{old: `extra.address.state`, new: `legacy_extra.address.state`},
	{old: `extra.results.1`, new: `legacy_extra.results.1`},
	{old: `extra.flow.role`, new: `parent.results.role`},
}

func TestMigrateContextReference(t *testing.T) {
	for _, tc := range testCases {
		actual, migrated := expressions.MigrateContextReference(tc.old)
		assert.True(t, migrated, "expected true for %s", tc.old)
		if migrated {
			assert.Equal(t, tc.new, actual, "migrated context reference mismatch for %s", tc.old)
		}
	}
}
