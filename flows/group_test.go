package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestGroupList(t *testing.T) {
	customers := test.NewGroup("Customers", "")
	testers := test.NewGroup("Testers", "")
	males := test.NewGroup("Males", "gender = \"M\"")
	groups := flows.NewGroupList([]*flows.Group{customers, testers, males})

	env := utils.NewEnvironmentBuilder().Build()

	// check use in expressions
	assert.Equal(t, types.NewXArray(
		customers.ToXValue(env),
		testers.ToXValue(env),
		males.ToXValue(env),
	), groups.ToXValue(env))
}
