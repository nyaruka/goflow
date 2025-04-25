package migrations_test

import (
	"os"
	"testing"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/flows/definition/migrations"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrimitives(t *testing.T) {
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

func TestGetObjectUUID(t *testing.T) {
	assert.Equal(t, uuids.UUID(""), migrations.GetObjectUUID(nil))
	assert.Equal(t, uuids.UUID(""), migrations.GetObjectUUID(map[string]any{}))
	assert.Equal(t, uuids.UUID(""), migrations.GetObjectUUID(map[string]any{"uuid": 234}))
	assert.Equal(t, uuids.UUID("foo"), migrations.GetObjectUUID(map[string]any{"uuid": "foo", "name": "bar"}))
}

func TestLocalizationPrimitives(t *testing.T) {
	readLocalization := func(j string) migrations.Localization {
		m, err := jsonx.DecodeGeneric([]byte(j))
		require.NoError(t, err)
		return migrations.Localization(m.(map[string]any))
	}

	l10n1 := readLocalization(`{
			"spa": {
				"8eebd020-1af5-431c-b943-aa670fc74da9": {
					"text": ["Hola"],
					"params": ["Rojo", "Verde"],
					"empty": [],
					"bad": {}
				}
			}
	}`)

	spa := l10n1.GetLanguageTranslation("spa")
	assert.NotNil(t, spa)
	assert.Nil(t, spa.GetTranslation("6f865930-e783-4fde-8e28-34b93b3a17c6", "text")) // no such item
	assert.Nil(t, spa.GetTranslation("8eebd020-1af5-431c-b943-aa670fc74da9", "foo"))  // no such property
	assert.Equal(t, []string{"Hola"}, spa.GetTranslation("8eebd020-1af5-431c-b943-aa670fc74da9", "text"))
	assert.Equal(t, []string{"Rojo", "Verde"}, spa.GetTranslation("8eebd020-1af5-431c-b943-aa670fc74da9", "params"))
	assert.Equal(t, []string{}, spa.GetTranslation("8eebd020-1af5-431c-b943-aa670fc74da9", "empty"))
	assert.Nil(t, spa.GetTranslation("8eebd020-1af5-431c-b943-aa670fc74da9", "bad")) // not strings

	spa.SetTranslation("8eebd020-1af5-431c-b943-aa670fc74da9", "text", []string{"Que tal"})
	spa.SetTranslation("8eebd020-1af5-431c-b943-aa670fc74da9", "bad", []string{"Mal"})

	test.AssertEqualJSON(t, []byte(`{
		"spa": {
			"8eebd020-1af5-431c-b943-aa670fc74da9": {
				"text": ["Que tal"],
				"params": ["Rojo", "Verde"],
				"empty": [],
				"bad": ["Mal"]
			}
		}
	}`), jsonx.MustMarshal(l10n1))

	// if item doesn't exist, should be created
	spa.SetTranslation("6f865930-e783-4fde-8e28-34b93b3a17c6", "text", []string{"Uno"})

	test.AssertEqualJSON(t, []byte(`{
		"spa": {
			"8eebd020-1af5-431c-b943-aa670fc74da9": {
				"text": ["Que tal"],
				"params": ["Rojo", "Verde"],
				"empty": [],
				"bad": ["Mal"]
			},
			"6f865930-e783-4fde-8e28-34b93b3a17c6": {
				"text": ["Uno"]
			}
		}
	}`), jsonx.MustMarshal(l10n1))

	spa.DeleteTranslation("8eebd020-1af5-431c-b943-aa670fc74da9", "text")
	spa.DeleteTranslation("8eebd020-1af5-431c-b943-aa670fc74da9", "empty")
	spa.DeleteTranslation("8eebd020-1af5-431c-b943-aa670fc74da9", "foo") // doesn't exist
	spa.DeleteTranslation("6f865930-e783-4fde-8e28-34b93b3a17c6", "text")

	test.AssertEqualJSON(t, []byte(`{
		"spa": {
			"8eebd020-1af5-431c-b943-aa670fc74da9": {
				"params": ["Rojo", "Verde"],
				"bad": ["Mal"]
			}
		}
	}`), jsonx.MustMarshal(l10n1))
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
		assert.Equal(t, []i18n.Language{"spa"}, f.Localization().Languages())
		assert.NotNil(t, f.Localization().GetLanguageTranslation("spa"))
		assert.Nil(t, f.Localization().GetLanguageTranslation("kin"))
	}

	// error trying to load something that is not a flow
	_, err = migrations.ReadFlow([]byte(`[]`))
	assert.EqualError(t, err, "flow definition isn't an object")

	// but tolerate other ways that a flow might be invalid since validation is version specific
	f, err = migrations.ReadFlow([]byte(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"type": "messaging",
		"nodes": null
	}`))
	assert.NoError(t, err)
	assert.Len(t, f.Nodes(), 0)

	f, err = migrations.ReadFlow([]byte(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"type": "messaging",
		"nodes": [
			{
				"uuid": "365293c7-633c-45bd-96b7-0b059766588d"
			}, 
			null
		]
	}`))
	assert.NoError(t, err)
	assert.Len(t, f.Nodes(), 1)

	f, err = migrations.ReadFlow([]byte(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"type": "messaging",
		"nodes": [
			{
				"uuid": "365293c7-633c-45bd-96b7-0b059766588d",
				"actions": [
					{
						"uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
						"type": "spam"
					},
					null
				]
			}, 
			null
		]
	}`))
	assert.NoError(t, err)
	assert.Len(t, f.Nodes(), 1)
	assert.Len(t, f.Nodes()[0].Actions(), 1)
	assert.Equal(t, "spam", f.Nodes()[0].Actions()[0].Type())
}

func readFlow(t *testing.T, path string) migrations.Flow {
	d, err := os.ReadFile(path)
	require.NoError(t, err)
	f, err := migrations.ReadFlow(d)
	require.NoError(t, err)
	return f
}
