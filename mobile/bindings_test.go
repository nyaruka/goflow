package mobile_test

import (
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/mobile"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMobileBindings(t *testing.T) {
	// if we try to create assets from invalid JSON
	_, err := mobile.NewAssetsSource("{")
	assert.Error(t, err)

	assetsJSON, err := ioutil.ReadFile("testdata/two_questions_offline.json")
	require.NoError(t, err)

	source, err := mobile.NewAssetsSource(string(assetsJSON))
	assert.NoError(t, err)

	sessionAssets, err := mobile.NewSessionAssets(source)
	assert.NoError(t, err)

	session := mobile.NewSession(sessionAssets, "mobile-test")

}
