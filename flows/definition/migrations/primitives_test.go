package migrations_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows/definition/migrations"

	"github.com/stretchr/testify/assert"
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
