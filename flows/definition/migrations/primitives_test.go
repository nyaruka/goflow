package migrations_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows/definition/migrations"

	"github.com/stretchr/testify/assert"
)

func TestMigrationPrimitives(t *testing.T) {
	f := migrations.Flow(map[string]interface{}{}) // nodes not set
	assert.Equal(t, []migrations.Node{}, f.Nodes())

	f = migrations.Flow(map[string]interface{}{"nodes": nil}) // nodes is nil
	assert.Equal(t, []migrations.Node{}, f.Nodes())

	f = migrations.Flow(map[string]interface{}{"nodes": []interface{}{}}) // nodes is empty
	assert.Equal(t, []migrations.Node{}, f.Nodes())

	f = migrations.Flow(map[string]interface{}{"nodes": []interface{}{
		map[string]interface{}{},
	}})
	assert.Equal(t, []migrations.Node{migrations.Node(map[string]interface{}{})}, f.Nodes())

	n := migrations.Node(map[string]interface{}{}) // actions and router are not set
	assert.Equal(t, []migrations.Action{}, n.Actions())
	assert.Nil(t, n.Router())

	n = migrations.Node(map[string]interface{}{"actions": nil, "router": nil}) // actions and router are nil
	assert.Equal(t, []migrations.Action{}, n.Actions())
	assert.Nil(t, n.Router())

	n = migrations.Node(map[string]interface{}{
		"actions": []interface{}{},
		"router":  map[string]interface{}{},
	})
	assert.Equal(t, []migrations.Action{}, n.Actions())
	assert.Equal(t, migrations.Router(map[string]interface{}{}), n.Router())

	a := migrations.Action(map[string]interface{}{}) // type not set
	assert.Equal(t, "", a.Type())

	a = migrations.Action(map[string]interface{}{"type": "foo"}) // type set
	assert.Equal(t, "foo", a.Type())
}
