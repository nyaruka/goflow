package main_test

import (
	"io/ioutil"
	"strings"
	"testing"
	"time"

	main "github.com/nyaruka/goflow/cmd/flowxgettext"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils/dates"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlowXGetText(t *testing.T) {
	defer dates.SetNowSource(dates.DefaultNowSource)
	dates.SetNowSource(dates.NewFixedNowSource(time.Date(2020, 3, 25, 13, 57, 30, 123456789, time.UTC)))

	out := &strings.Builder{}

	err := main.FlowXGetText(envs.Language("fra"), false, []string{"../../test/testdata/runner/two_questions.json"}, out)
	require.NoError(t, err)

	expected, err := ioutil.ReadFile("testdata/two_questions.fr.po")
	require.NoError(t, err)

	assert.Equal(t, string(expected), out.String())
}
