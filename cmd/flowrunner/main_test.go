package main_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	main "github.com/nyaruka/goflow/cmd/flowrunner"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunFlow(t *testing.T) {
	// create an input than can be scanned for two answers
	in := strings.NewReader("I like red\npepsi\n")
	out := &strings.Builder{}

	_, err := main.RunFlow(test.NewEngine(), "testdata/two_questions.json", assets.FlowUUID("615b8a0f-588c-4d20-a05f-363b0b4ce6f4"), "", "eng", in, out)
	require.NoError(t, err)

	// remove input prompts and split output by line to get each event
	lines := strings.Split(strings.Replace(out.String(), "> ", "", -1), "\n")

	assert.Equal(t, []string{
		"Starting flow 'Two Questions'....",
		"---------------------------------------",
		"💬 message created \"Hi Ben Haggerty! What is your favorite color? (red/blue)\"",
		"⏳ waiting for message (600 sec timeout, type /timeout to simulate)...",
		"📥 message received \"I like red\"",
		"📈 run result 'Favorite Color' changed to 'red' with category 'Red'",
		"🌐 language changed to 'fra'",
		"💬 message created \"Red it is! What is your favorite soda? (pepsi/coke)\"",
		"⏳ waiting for message...",
		"📥 message received \"pepsi\"",
		"📈 run result 'Soda' changed to 'pepsi' with category 'Pepsi'",
		"💬 message created \"Great, you are done!\"",
		"",
	}, lines)

	// run again but don't specify the flow
	in = strings.NewReader("I like red\npepsi\n")
	out = &strings.Builder{}
	_, err = main.RunFlow(test.NewEngine(), "testdata/two_questions.json", "", "", "eng", in, out)
	require.NoError(t, err)

	assert.Contains(t, out.String(), "Starting flow 'Two Questions'")
}

func TestPrintEvent(t *testing.T) {
	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	sa := session.Assets()
	flow, _ := sa.Flows().Get("50c3706e-fedb-42c0-8eab-dda3335714b7")
	timeout := 3
	expiresOn := time.Date(2022, 2, 3, 13, 45, 30, 0, time.UTC)

	tests := []struct {
		event    flows.Event
		expected string
	}{
		{events.NewBroadcastCreated(map[envs.Language]*events.BroadcastTranslation{"eng": {Text: "hello"}}, "eng", nil, nil, nil), `🔉 broadcasted 'hello' to ...`},
		{events.NewContactFieldChanged(sa.Fields().Get("gender"), flows.NewValue(types.NewXText("M"), nil, nil, "", "", "")), `✏️ field 'gender' changed to 'M'`},
		{events.NewContactFieldChanged(sa.Fields().Get("gender"), nil), `✏️ field 'gender' cleared`},
		{events.NewContactGroupsChanged([]*flows.Group{sa.Groups().Get("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d")}, nil), `👪 added to 'Testers'`},
		{events.NewContactGroupsChanged(nil, []*flows.Group{sa.Groups().Get("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d")}), `👪 removed from 'Testers'`},
		{events.NewContactLanguageChanged("eng"), `🌐 language changed to 'eng'`},
		{events.NewContactNameChanged("Jim"), `📛 name changed to 'Jim'`},
		{events.NewContactRefreshed(session.Contact()), `👤 contact refreshed on resume`},
		{events.NewContactTimezoneChanged(session.Environment().Timezone()), `🕑 timezone changed to 'America/Guayaquil'`},
		{events.NewDialEnded(flows.NewDial(flows.DialStatusBusy, 3)), `☎️ dial ended with 'busy'`},
		{events.NewDialWait(urns.URN(`tel:+1234567890`), nil), `⏳ waiting for dial (type /dial <answered|no_answer|busy|failed>)...`},
		{events.NewEmailSent([]string{"code@example.com"}, "Hi", "What up?"), `✉️ email sent with subject 'Hi'`},
		{events.NewEnvironmentRefreshed(session.Environment()), `⚙️ environment refreshed on resume`},
		{events.NewErrorf("this didn't work"), `⚠️ this didn't work`},
		{events.NewFailure(errors.New("this really didn't work")), `🛑 this really didn't work`},
		{events.NewFlowEntered(flow.Reference(), "", false), `↪️ entered flow 'Registration'`},
		{events.NewInputLabelsAdded("2a786bbc-2314-4d57-a0c9-b66e1642e5e2", []*flows.Label{sa.Labels().FindByName("Spam")}), `🏷️ labeled with 'Spam'`},
		{events.NewMsgWait(nil, nil, nil), `⏳ waiting for message...`},
		{events.NewMsgWait(&timeout, &expiresOn, nil), `⏳ waiting for message (3 sec timeout, type /timeout to simulate)...`},
	}

	for _, tc := range tests {
		out := &strings.Builder{}
		main.PrintEvent(tc.event, out)
		assert.Equal(t, tc.expected, out.String(), "event print mismatch for event type '%s'", tc.event.Type())
	}
}
