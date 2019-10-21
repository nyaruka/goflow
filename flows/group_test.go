package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
)

func TestGroupList(t *testing.T) {
	customers := test.NewGroup("Customers", "")
	testers := test.NewGroup("Testers", "")
	males := test.NewGroup("Males", `gender = "M"`)

	assert.Equal(t, "Customers", customers.Name())
	assert.Equal(t, `gender = "M"`, males.Query())

	groups := flows.NewGroupList([]*flows.Group{customers, testers, males})

	env := envs.NewBuilder().Build()

	// check use in expressions
	test.AssertXEqual(t, types.NewXArray(
		customers.ToXValue(env),
		testers.ToXValue(env),
		males.ToXValue(env),
	), groups.ToXValue(env))
}
