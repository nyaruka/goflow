package migrations_test

import (
	"os"
	"testing"

	"github.com/nyaruka/goflow/flows/definition/migrations"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrationPrimitives(t *testing.T) {
	f := migrations.Flow(map[string]any{}) // nodes not set
	assert.Equal(t, []migrations.Node{}, f.Nodes())

	f = migrations.Flow(map[string]any{"nodes": nil}) // nodes is nil
	assert.Equal(t, []migrations.Node{}, f.Nodes())

	f = migrations.Flow(map[string]any{"nodes": []any{}}) // nodes is empty
	assert.Equal(t, []migrations.Node{}, f.Nodes())

	f = migrations.Flow(map[string]any{"nodes": []any{
		map[string]any{},
	}})
	assert.Equal(t, []migrations.Node{migrations.Node(map[string]any{})}, f.Nodes())

	n := migrations.Node(map[string]any{}) // actions and router are not set
	assert.Equal(t, []migrations.Action{}, n.Actions())
	assert.Nil(t, n.Router())

	n = migrations.Node(map[string]any{"actions": nil, "router": nil}) // actions and router are nil
	assert.Equal(t, []migrations.Action{}, n.Actions())
	assert.Nil(t, n.Router())

	n = migrations.Node(map[string]any{
		"actions": []any{},
		"router":  map[string]any{},
	})
	assert.Equal(t, []migrations.Action{}, n.Actions())
	assert.Equal(t, migrations.Router(map[string]any{}), n.Router())

	a := migrations.Action(map[string]any{}) // type not set
	assert.Equal(t, "", a.Type())

	a = migrations.Action(map[string]any{"type": "foo"}) // type set
	assert.Equal(t, "foo", a.Type())
}

func TestReadFlow(t *testing.T) {
	f, err := migrations.ReadFlow([]byte(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"name": "Test Flow",
		"spec_version": "13.2.0",
		"language": "eng",
		"type": "messaging",
		"localization": {
			"spa": {
				"8eebd020-1af5-431c-b943-aa670fc74da9": {
					"text": ["Hola"]
				}
			}
		},
		"nodes": [
			{
				"uuid": "365293c7-633c-45bd-96b7-0b059766588d",
				"actions": [
					{
						"uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
						"type": "send_msg",
						"text": "Hello",
						"attachments": ["image/jpeg:foo.jpg", "audio/mp3:foo.mp3"]
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
	assert.NoError(t, err)
	if assert.Len(t, f.Nodes(), 1) {
		assert.Len(t, f.Nodes()[0].Actions(), 1)
		assert.Nil(t, f.Nodes()[0].Router())
	}
	if assert.NotNil(t, f.Localization()) {
		assert.Equal(t, []string{"spa"}, f.Localization().Languages())
		assert.NotNil(t, f.Localization().GetLanguageTranslation("spa"))
		assert.Nil(t, f.Localization().GetLanguageTranslation("kin"))
	}

	_, err = migrations.ReadFlow([]byte(`[]`))
	assert.EqualError(t, err, "flow definition isn't an object")
}

func readFlow(t *testing.T, path string) migrations.Flow {
	d, err := os.ReadFile(path)
	require.NoError(t, err)
	f, err := migrations.ReadFlow(d)
	require.NoError(t, err)
	return f
}
