package definition_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssetsValidation(t *testing.T) {
	sa, err := test.LoadSessionAssets("testdata/assets_with_only_flow.json")
	require.NoError(t, err)

	missingAssets, err := sa.Validate("70def698-3ecb-41e1-a9b8-28a104828a13")
	assert.NoError(t, err)
	assert.Equal(t, []assets.Reference{
		assets.NewFieldReference("nick_name", ""),
		assets.NewGroupReference(assets.GroupUUID("7be2f40b-38a0-4b06-9e6d-522dca592cc8"), "Registered"),
	}, missingAssets)
}
