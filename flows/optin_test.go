package flows_test

import (
	"testing"
	"time"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptIns(t *testing.T) {
	defer uuids.SetGenerator(uuids.DefaultGenerator)
	uuids.SetGenerator(uuids.NewSeededGenerator(12345, time.Now))

	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(`{
		"optins": [
			{
				"uuid": "248be71d-78e9-4d71-a6c4-9981d369e5cb",
				"name": "Joke Of The Day"
			}
		]
	}`))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	jotd := sa.OptIns().Get("248be71d-78e9-4d71-a6c4-9981d369e5cb")
	assert.Equal(t, assets.OptInUUID("248be71d-78e9-4d71-a6c4-9981d369e5cb"), jotd.UUID())
	assert.Equal(t, "Joke Of The Day", jotd.Name())
	assert.Equal(t, assets.NewOptInReference("248be71d-78e9-4d71-a6c4-9981d369e5cb", "Joke Of The Day"), jotd.Reference())

	// check use in expressions
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("Joke Of The Day"),
		"uuid":        types.NewXText("248be71d-78e9-4d71-a6c4-9981d369e5cb"),
		"name":        types.NewXText("Joke Of The Day"),
	}), flows.Context(env, jotd))

	assert.Nil(t, (*flows.OptIn)(nil).Reference())
}
