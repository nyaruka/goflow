package migrations_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/definition/migrations"
	"github.com/nyaruka/goflow/flows/inspect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	flow, err := migrations.ReadFlow([]byte(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"name": "Test Flow",
		"spec_version": "13.2.0",
		"language": "und",
		"type": "messaging",
		"nodes": [
			{
				"uuid": "365293c7-633c-45bd-96b7-0b059766588d",
				"actions": [
					{
						"uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
						"type": "send_msg",
						"text": "Hello"
					}
				],
				"exits": [
					{
						"uuid": "b6f4caf3-ec99-44d5-a40c-8600ac0e2eac"
					}
				]
			}
		]
	}`))
	require.NoError(t, err)

	migrations.RewriteTemplates(flow, migrations.GetTemplateCatalog(definition.CurrentSpecVersion), func(s string) string { return strings.ToUpper(s) })
}
