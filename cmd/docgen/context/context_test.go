package context_test

import (
	"testing"

	"github.com/nyaruka/goflow/cmd/docgen/context"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	groupType := context.NewStaticType("group", []*context.Property{
		context.NewProperty("uuid", "the UUID of the group", "text"),
		context.NewProperty("name", "the name of the group", "text"),
	})

	fieldsType := context.NewDynamicType("fields", "field-keys", context.NewProperty("{key}", "the value of {key}", "any"))

	contactType := context.NewStaticType("contact", []*context.Property{
		context.NewProperty("name", "the full name of the contact", "text"),
		context.NewProperty("fields", "the custom field values of the contact", "fields"),
		context.NewArrayProperty("groups", "the groups that the contact belongs to", "group"),
	})

	resultType := context.NewStaticType("result", []*context.Property{
		context.NewProperty("__default__", "the value of the result", "text"),
		context.NewProperty("value", "the value of the result", "text"),
		context.NewProperty("category", "the category of the result", "text"),
	})

	resultsType := context.NewDynamicType("results", "result-keys", context.NewProperty("{key}", "the result of {key}", "result"))

	ctx := context.NewContext()
	ctx.AddType(groupType)
	ctx.AddType(fieldsType)
	ctx.AddType(contactType)
	ctx.AddType(resultType)
	ctx.AddType(resultsType)
	ctx.SetRoot([]*context.Property{
		context.NewProperty("contact", "the run contact", "contact"),
		context.NewProperty("results", "the run results", "results"),
	})

	assert.Nil(t, ctx.Validate())
}
