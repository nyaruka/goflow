package runs_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/jsonx"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunSummary(t *testing.T) {
	uuids.SetGenerator(uuids.NewSeededGenerator(123456))
	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2018, 7, 6, 12, 30, 0, 123456789, time.UTC)))
	defer uuids.SetGenerator(uuids.DefaultGenerator)
	defer dates.SetNowSource(dates.DefaultNowSource)

	server := test.NewTestHTTPServer(49999)
	defer server.Close()

	session, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	require.NoError(t, err)

	run := session.Runs()[0]
	summary := run.Snapshot()

	assert.Equal(t, run.Flow(), summary.Flow())
	assert.Equal(t, run.Contact(), summary.Contact())
	assert.Equal(t, run.Results(), summary.Results())
	assert.Equal(t, run.Status(), summary.Status())
	assert.Equal(t, run.Results(), summary.Results())

	assert.Equal(t, "Ryan Lewis@Registration", runs.FormatRunSummary(session.Environment(), summary))

	// test marshaling and unmarshaling
	marshaled, err := jsonx.Marshal(summary)
	require.NoError(t, err)

	summary, err = runs.ReadRunSummary(session.Assets(), marshaled, assets.PanicOnMissing)
	require.NoError(t, err)

	assert.Equal(t, run.Flow().Name(), summary.Flow().Name())
	assert.Equal(t, run.Status(), summary.Status())
	assert.Equal(t, "Ryan Lewis@Registration", runs.FormatRunSummary(session.Environment(), summary))

	// try reading with missing assets
	emptyAssets, err := engine.NewSessionAssets(session.Environment(), static.NewEmptySource(), nil)

	summary, err = runs.ReadRunSummary(emptyAssets, marshaled, assets.IgnoreMissing)
	require.NoError(t, err)

	assert.Nil(t, summary.Flow())
	assert.Equal(t, run.Status(), summary.Status())
	assert.Equal(t, "Ryan Lewis@<missing>", runs.FormatRunSummary(session.Environment(), summary))

	// try removing the contact (they're optional) and re-reading
	marshaled = test.JSONDelete(marshaled, []string{"contact"})

	summary, err = runs.ReadRunSummary(session.Assets(), marshaled, assets.PanicOnMissing)
	require.NoError(t, err)

	assert.Equal(t, run.Flow().Name(), summary.Flow().Name())
	assert.Equal(t, run.Status(), summary.Status())
	assert.Equal(t, "<nocontact>@Registration", runs.FormatRunSummary(session.Environment(), summary))
}
