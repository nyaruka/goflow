package completion_test

import (
	"testing"

	"github.com/nyaruka/goflow/cmd/docgen/completion"
	"github.com/stretchr/testify/assert"
)

func TestCompletion(t *testing.T) {
	groupType := completion.NewStaticType("group", []*completion.Property{
		completion.NewProperty("uuid", "the UUID of the group", "text"),
		completion.NewProperty("name", "the name of the group", "text"),
	})

	fieldsType := completion.NewDynamicType("fields", "fields", completion.NewProperty("{key}", "the value of {key}", "any"))

	contactType := completion.NewStaticType("contact", []*completion.Property{
		completion.NewProperty("name", "the full name of the contact", "text"),
		completion.NewProperty("fields", "the custom field values of the contact", "fields"),
		completion.NewArrayProperty("groups", "the groups that the contact belongs to", "group"),
	})

	c := completion.NewCompletion(
		[]completion.Type{groupType, fieldsType, contactType},
		[]*completion.Property{
			completion.NewProperty("contact", "the run contact", "contact"),
		},
	)

	// all type refs should be valid
	assert.Nil(t, c.Validate())

	nodes := c.EnumerateNodes(completion.NewContext(map[string][]string{
		"fields": {"age", "gender"},
	}))

	assert.Equal(t, []completion.Node{
		{Path: "contact", Help: "the run contact"},
		{Path: "contact.name", Help: "the full name of the contact"},
		{Path: "contact.fields", Help: "the custom field values of the contact"},
		{Path: "contact.fields.age", Help: "the value of age"},
		{Path: "contact.fields.gender", Help: "the value of gender"},
		{Path: "contact.groups", Help: "the groups that the contact belongs to"},
		{Path: "contact.groups[0]", Help: "first of the groups that the contact belongs to"},
		{Path: "contact.groups[0].uuid", Help: "the UUID of the group"},
		{Path: "contact.groups[0].name", Help: "the name of the group"},
	}, nodes)
}
