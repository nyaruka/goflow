package engine_test

import (
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunSummary(t *testing.T) {
	uuids.SetGenerator(uuids.NewSeededGenerator(123456, time.Now))
	dates.SetNowFunc(dates.NewSequentialNow(time.Date(2018, 7, 6, 12, 30, 0, 123456789, time.UTC), time.Second))
	defer uuids.SetGenerator(uuids.DefaultGenerator)
	defer dates.SetNowFunc(time.Now)

	server := test.NewHTTPServer(49999, test.MockWebhooksHandler)
	defer server.Close()

	session, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	require.NoError(t, err)

	run := session.Runs()[0]
	summary := run.Snapshot()

	assert.Equal(t, run.Flow(), summary.Flow())
	assert.Equal(t, run.Contact().UUID(), summary.Contact().UUID())
	assert.Equal(t, run.Contact().Name(), summary.Contact().Name())
	assert.Equal(t, run.Contact().Routes(), summary.Contact().Routes())
	assert.Equal(t, run.Contact().Status(), summary.Contact().Status())
	assert.Equal(t, run.Contact().Language(), summary.Contact().Language())
	assert.Equal(t, run.Contact().Fields(), summary.Contact().Fields())
	assert.Equal(t, run.Contact().Groups(), summary.Contact().Groups())
	assert.Equal(t, run.Results(), summary.Results())
	assert.Equal(t, run.Status(), summary.Status())
	assert.Equal(t, run.Results(), summary.Results())

	assert.Equal(t, "Ryan Lewis@Registration", engine.FormatRunSummary(session.Environment(), summary))

	// test marshaling and unmarshaling
	marshaled, err := jsonx.Marshal(summary)
	require.NoError(t, err)

	summary, err = engine.ReadRunSummary(session.Assets(), marshaled, assets.PanicOnMissing)
	require.NoError(t, err)

	assert.Equal(t, run.Flow().Name(), summary.Flow().Name())
	assert.Equal(t, run.Status(), summary.Status())
	assert.Equal(t, "Ryan Lewis@Registration", engine.FormatRunSummary(session.Environment(), summary))

	// try reading with missing assets
	emptyAssets, err := engine.NewSessionAssets(session.Environment(), static.NewEmptySource(), nil)
	assert.NoError(t, err)

	summary, err = engine.ReadRunSummary(emptyAssets, marshaled, assets.IgnoreMissing)
	require.NoError(t, err)

	assert.Nil(t, summary.Flow())
	assert.Equal(t, run.Status(), summary.Status())
	assert.Equal(t, "Ryan Lewis@<missing>", engine.FormatRunSummary(session.Environment(), summary))

	// try removing the contact (they're optional) and re-reading
	marshaled = test.JSONDelete(marshaled, []string{"contact"})

	summary, err = engine.ReadRunSummary(session.Assets(), marshaled, assets.PanicOnMissing)
	require.NoError(t, err)

	assert.Equal(t, run.Flow().Name(), summary.Flow().Name())
	assert.Equal(t, run.Status(), summary.Status())
	assert.Equal(t, "<nocontact>@Registration", engine.FormatRunSummary(session.Environment(), summary))
}
