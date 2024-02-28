package migrations_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/definition/migrations"
	"github.com/nyaruka/goflow/flows/inspect"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
)

func TestCurrentTemplateCatalog(t *testing.T) {
	s := &migrations.TemplateCatalog{
		Actions: make(map[string][]string),
		Routers: map[string][]string{
			"random": {".operand", ".cases[*].arguments[*]"},
			"switch": {".operand", ".cases[*].arguments[*]"},
		},
	}

	for typeName, fn := range actions.RegisteredTypes() {
		actionType := reflect.TypeOf(fn())

		s.Actions[typeName] = inspect.TemplatePaths(actionType)
	}

	assert.Equal(t, migrations.GetTemplateCatalog(definition.CurrentSpecVersion), s)
}

func TestRewriteTemplates(t *testing.T) {
	flow := readFlow(t, "testdata/templates1.json")
	expected := readFlow(t, "testdata/templates1.upper.json")

	migrations.RewriteTemplates(flow, migrations.GetTemplateCatalog(definition.CurrentSpecVersion), func(s string) string { return strings.ToUpper(s) })

	test.AssertEqualJSON(t, jsonx.MustMarshal(expected), jsonx.MustMarshal(flow), "template rewrite mismatch")
}
