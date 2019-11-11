package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/require"
)

func TestGlobals(t *testing.T) {
	source, err := static.NewSource([]byte(`{
		"globals": [
			{"key": "org_name", "name": "Org Name", "value": "U-Report"},
			{"key": "access_token", "name": "Access Token", "value": "674372272"}
		]
	}`))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(source, nil)
	require.NoError(t, err)

	env := envs.NewBuilder().Build()

	// check use in expressions
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__":  types.NewXText("org_name: U-Report\naccess_token: 674372272"),
		"access_token": types.NewXText("674372272"),
		"org_name":     types.NewXText("U-Report"),
	}), flows.Context(env, sa.Globals()))
}
