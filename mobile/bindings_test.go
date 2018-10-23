package mobile_test

import (
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/mobile"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMobileBindings(t *testing.T) {
	// error if we try to create assets from invalid JSON
	_, err := mobile.NewAssetsSource("{")
	assert.Error(t, err)

	// can load a standard assets file
	assetsJSON, err := ioutil.ReadFile("../test/testdata/flows/two_questions_offline.json")
	require.NoError(t, err)

	source, err := mobile.NewAssetsSource(string(assetsJSON))
	require.NoError(t, err)

	// and create a new session assets
	sessionAssets, err := mobile.NewSessionAssets(source)
	require.NoError(t, err)

	environment, err := mobile.NewEnvironment("DD-MM-YYYY", "tt:mm", "Africa/Kigali", "eng", []string{"eng", "fra"})
	require.NoError(t, err)

	contact := mobile.NewEmptyContact()

	trigger := mobile.NewManualTrigger(environment, contact, "7c3db26f-e12a-48af-9673-e2feefdf8516", "Two Questions")

	session := mobile.NewSession(sessionAssets, "mobile-test")

	err = session.Start(trigger)
	require.NoError(t, err)
}
