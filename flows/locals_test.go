package flows_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
)

func TestLocals(t *testing.T) {
	l1 := flows.NewLocals()
	assert.True(t, l1.IsZero())

	l1.Set("foo", types.NewXText("bar"))
	l1.Set("int", types.NewXNumberFromInt(42))
	l1.Set("obj", types.NewXObject(map[string]types.XValue{"sub": types.NewXText("baz")}))

	assert.Equal(t, types.NewXText("bar"), l1.Get("foo"))
	assert.Equal(t, types.NewXNumberFromInt(42), l1.Get("int"))
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{"sub": types.NewXText("baz")}), l1.Get("obj"))
	assert.False(t, l1.IsZero())

	marshaled, err := json.Marshal(l1)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"foo":"bar","int":42,"obj":{"sub":"baz"}}`, string(marshaled))

	var l2 flows.Locals
	err = json.Unmarshal(marshaled, &l2)
	assert.NoError(t, err)

	assert.Equal(t, types.NewXText("bar"), l2.Get("foo"))
	assert.Equal(t, types.NewXNumberFromInt(42), l2.Get("int"))
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{"sub": types.NewXText("baz")}), l2.Get("obj"))
}
