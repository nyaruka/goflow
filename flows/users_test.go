package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	atypes "github.com/nyaruka/goflow/assets/static/types"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	ua1 := atypes.NewUser("bob@nyaruka.com", "Bob")
	ua2 := atypes.NewUser("jim@nyaruka.com", "")

	ua := flows.NewUserAssets([]assets.User{ua1, ua2})

	u1 := ua.Get("bob@nyaruka.com")

	assert.Equal(t, "Bob", u1.Name())
	assert.Equal(t, ua1, u1.Asset())
	assert.Equal(t, assets.NewUserReference("bob@nyaruka.com", "Bob"), u1.Reference())

	env := envs.NewBuilder().Build()

	// check use in expressions
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("bob@nyaruka.com"),
		"email":       types.NewXText("bob@nyaruka.com"),
		"name":        types.NewXText("Bob"),
	}), flows.Context(env, u1))
}
