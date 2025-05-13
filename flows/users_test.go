package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	ua1 := static.NewUser("0c78ef47-7d56-44d8-8f57-96e0f30e8f44", "Bob McTickets", "bob@nyaruka.com")
	ua2 := static.NewUser("c8945bcc-5d4b-495f-b3ea-2662c6070fe3", "", "jim@nyaruka.com")

	ua := flows.NewUserAssets([]assets.User{ua1, ua2})

	u1 := ua.Get("bob@nyaruka.com")

	assert.Equal(t, "Bob McTickets", u1.Format())
	assert.Equal(t, "Bob McTickets", u1.Name())
	assert.Equal(t, ua1, u1.Asset())
	assert.Equal(t, assets.NewUserReference("bob@nyaruka.com", "Bob McTickets"), u1.Reference())

	// nil object returns nil reference
	assert.Nil(t, (*flows.User)(nil).Reference())

	env := envs.NewBuilder().Build()

	// check use in expressions
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("Bob McTickets"),
		"email":       types.NewXText("bob@nyaruka.com"),
		"name":        types.NewXText("Bob McTickets"),
		"first_name":  types.NewXText("Bob"),
	}), flows.Context(env, u1))

	u2 := ua.Get("jim@nyaruka.com")

	assert.Equal(t, "jim@nyaruka.com", u2.Format()) // fallsback on email

	// check use in expressions
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("jim@nyaruka.com"),
		"email":       types.NewXText("jim@nyaruka.com"),
		"name":        types.NewXText(""),
		"first_name":  types.NewXText(""),
	}), flows.Context(env, u2))
}
