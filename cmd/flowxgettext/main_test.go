package main_test

import (
	"strings"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/i18n"
	main "github.com/nyaruka/goflow/cmd/flowxgettext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlowXGetText(t *testing.T) {
	defer dates.SetNowFunc(time.Now)
	dates.SetNowFunc(dates.NewFixedNow(time.Date(2020, 3, 25, 13, 57, 30, 123456789, time.UTC)))

	out := &strings.Builder{}

	err := main.FlowXGetText(i18n.Language("fra"), false, []string{"../../test/testdata/runner/two_questions.json"}, out)
	require.NoError(t, err)

	assert.Contains(t, out.String(), `
#: Two+Questions/2ab9b033-77a8-4e56-a558-b568c00c9492/name:0
msgid "Pepsi"
msgstr "Pepsi"
`)
}
