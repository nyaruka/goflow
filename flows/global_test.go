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

func TestGlobals(t *testing.T) {
	ga1 := atypes.NewGlobal("org_name", "Org Name", "U-Report")
	ga2 := atypes.NewGlobal("access_token", "Access Token", "674372272")

	ga := flows.NewGlobalAssets([]assets.Global{ga1, ga2})

	g1 := ga.Get("org_name")

	assert.Equal(t, "Org Name", g1.Name())
	assert.Equal(t, ga1, g1.Asset())
	assert.Equal(t, assets.NewGlobalReference("org_name", "Org Name"), g1.Reference())

	env := envs.NewBuilder().Build()

	// check use in expressions
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__":  types.NewXText("Org Name: U-Report\nAccess Token: 674372272"),
		"access_token": types.NewXText("674372272"),
		"org_name":     types.NewXText("U-Report"),
	}), flows.Context(env, ga))
}
