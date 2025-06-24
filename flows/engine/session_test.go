package engine_test

import (
	"context"
	"encoding/json"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluateTemplate(t *testing.T) {
	testFile, err := os.ReadFile("testdata/templates.json")
	require.NoError(t, err)

	test.MockUniverse()

	server := test.NewTestHTTPServer(49992)
	defer server.Close()

	sessionWithURNs, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	require.NoError(t, err)
	sessionWithoutURNs, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyURNs)
	require.NoError(t, err)

	tests := []struct {
		Template   string `json:"template"`
		RedactURNs bool   `json:"redact_urns,omitempty"`

		Output     string          `json:"output,omitempty"`
		OutputJSON json.RawMessage `json:"output_json,omitempty"`
		Error      string          `json:"error,omitempty"`
	}{}

	jsonx.MustUnmarshal(testFile, &tests)

	for i, tc := range tests {
		var run flows.Run
		if tc.RedactURNs {
			run = sessionWithoutURNs.Runs()[0]
		} else {
			run = sessionWithURNs.Runs()[0]
		}

		log := test.NewEventLog()
		eval, ok := run.EvaluateTemplate(tc.Template, log.Log)

		// clone test case and populate with actual values
		actual := tc
		if tc.OutputJSON != nil {
			actual.OutputJSON = []byte(eval)
		} else {
			actual.Output = eval
		}
		if !ok {
			actual.Error = log.Error().Error()
		}

		if !test.UpdateSnapshots {
			if tc.OutputJSON != nil {
				test.AssertEqualJSON(t, tc.OutputJSON, actual.OutputJSON, "output mismatch evaluating template: '%s'", tc.Template)
			} else {
				assert.Equal(t, tc.Output, actual.Output, "output mismatch evaluating template: '%s'", tc.Template)
			}
			assert.Equal(t, tc.Error, actual.Error, "error mismatch evaluating template: '%s'", tc.Template)
		} else {
			tests[i] = actual
		}
	}

	if test.UpdateSnapshots {
		actualJSON, err := jsonx.MarshalPretty(tests)
		require.NoError(t, err)

		err = os.WriteFile("testdata/templates.json", actualJSON, 0666)
		require.NoError(t, err)
	}
}

func BenchmarkEvaluateTemplate(b *testing.B) {
	testFile, err := os.ReadFile("testdata/templates.json")
	require.NoError(b, err)

	session, _, err := test.CreateTestSession("http://localhost", envs.RedactionPolicyNone)
	require.NoError(b, err)

	run := session.Runs()[0]

	tests := []struct {
		Template   string `json:"template"`
		RedactURNs bool   `json:"redact_urns,omitempty"`

		Output string `json:"output,omitempty"`
		Error  string `json:"error,omitempty"`
	}{}

	jsonx.MustUnmarshal(testFile, &tests)

	logEvent := func(e flows.Event) {}

	for n := 0; n < b.N; n++ {
		for _, tc := range tests {
			run.EvaluateTemplate(tc.Template, logEvent)
		}
	}
}

func TestReadWithMissingAssets(t *testing.T) {
	// create standard test session and marshal to JSON
	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	sessionJSON, err := jsonx.Marshal(session)
	require.NoError(t, err)

	// try to read it back but with no assets
	sessionAssets, err := engine.NewSessionAssets(session.Environment(), static.NewEmptySource(), nil)
	require.NoError(t, err)

	missingAssets := make([]assets.Reference, 0)
	missing := func(a assets.Reference, err error) { missingAssets = append(missingAssets, a) }

	eng := engine.NewBuilder().Build()
	_, err = eng.ReadSession(sessionAssets, sessionJSON, session.Environment(), session.Contact(), nil, missing)
	require.NoError(t, err)

	refs := make([]string, len(missingAssets))
	for i := range missingAssets {
		refs[i] = missingAssets[i].String()
	}

	// ordering isn't deterministic so sort A-Z
	sort.Strings(refs)

	assert.Equal(t, []string{
		"channel[uuid=57f1078f-88aa-46f4-a59a-948a5739c03d,name=My Android Phone]",
		"flow[uuid=50c3706e-fedb-42c0-8eab-dda3335714b7,name=Registration]",
		"flow[uuid=b7cf0d83-f1c9-411c-96fd-c511a4cfa86d,name=Collect Age]",
	}, refs)
}

func TestQueryBasedGroupReevaluationOnTrigger(t *testing.T) {
	assetsJSON, err := os.ReadFile("testdata/smart_groups.json")
	require.NoError(t, err)

	sa, err := test.CreateSessionAssets(assetsJSON, "")
	require.NoError(t, err)

	// contact is in wrong groups
	contact, err := flows.ReadContact(sa, []byte(`{
		"uuid": "6d116680-eab9-460a-9c6e-1f05d3c5b5d6",
		"created_on": "2018-06-20T11:40:30.123456789-00:00",
        "groups": [
            {"uuid": "047de1c9-9189-4f4c-aa04-bff0a4c2efb6", "name": "Males"}
        ],
        "fields": {
            "gender": {
                "text": "Female"
			}
		}
	}`), assets.PanicOnMissing)
	require.NoError(t, err)

	env := envs.NewBuilder().Build()
	trigger := triggers.NewBuilder(assets.NewFlowReference("1b462ce8-983a-4393-b133-e15a0efdb70c", "")).Manual().Build()
	eng := engine.NewBuilder().Build()

	session, sprint, err := eng.NewSession(context.Background(), sa, env, contact, trigger, nil)
	require.NoError(t, err)

	assert.Equal(t, 1, len(sprint.Events()))
	assert.Equal(t, "contact_groups_changed", sprint.Events()[0].Type())
	assert.Equal(t, 1, session.Contact().Groups().Count())
	assert.Equal(t, "Females", session.Contact().Groups().All()[0].Name())
}

func TestRunResuming(t *testing.T) {
	assetsJSON, err := os.ReadFile("testdata/subflows.json")
	require.NoError(t, err)

	_, session, sprint := test.NewSessionBuilder().WithAssetsJSON(assetsJSON).WithFlow("72162f46-dce3-4798-9f19-384a2447efc5").MustBuild()

	// each run should be marked as completed
	assert.Equal(t, 3, len(session.Runs()))
	assert.Equal(t, flows.RunStatusCompleted, session.Runs()[0].Status())
	assert.Equal(t, flows.RunStatusCompleted, session.Runs()[1].Status())
	assert.Equal(t, flows.RunStatusCompleted, session.Runs()[2].Status())
	assert.Len(t, sprint.Flows(), 3)

	// change the UUID of the third flow so the enter_flow in the second flow will error
	assetsWithoutChildFlow := test.JSONReplace(assetsJSON, []string{"flows", "[2]", "uuid"}, []byte(`"653a3fa3-ff59-4a89-93c3-a8b9486ec479"`))

	_, session, _ = test.NewSessionBuilder().WithAssetsJSON(assetsWithoutChildFlow).WithFlow("72162f46-dce3-4798-9f19-384a2447efc5").MustBuild()

	// each run should be marked as failed
	assert.Equal(t, 2, len(session.Runs()))
	assert.Equal(t, flows.RunStatusFailed, session.Runs()[0].Status())
	assert.Equal(t, flows.RunStatusFailed, session.Runs()[1].Status())
}

func TestResumeAfterWaitWithMissingFlowAssets(t *testing.T) {
	assetsJSON, err := os.ReadFile("../../test/testdata/runner/subflow.json")
	require.NoError(t, err)

	_, session1, _ := test.NewSessionBuilder().WithAssetsJSON(assetsJSON).WithFlow("76f0a02f-3b75-4b86-9064-e9195e1b3a02").MustBuild()

	assert.Equal(t, flows.SessionStatusWaiting, session1.Status())
	assert.Equal(t, flows.RunStatusActive, session1.Runs()[0].Status())
	assert.Equal(t, flows.RunStatusWaiting, session1.Runs()[1].Status())

	// change the UUID of the child flow so it will effectively be missing
	assetsWithoutChildFlow := test.JSONReplace(assetsJSON, []string{"flows", "[1]", "uuid"}, []byte(`"653a3fa3-ff59-4a89-93c3-a8b9486ec479"`))

	sa, err := test.CreateSessionAssets(assetsWithoutChildFlow, "")
	require.NoError(t, err)

	session2, _, err := test.ResumeSession(session1, sa, "Hello")
	require.NoError(t, err)

	// should have a failed session (with no runs left was active/waiting)
	assert.Equal(t, flows.SessionStatusFailed, session2.Status())
	assert.Equal(t, flows.RunStatusFailed, session2.Runs()[0].Status())
	assert.Equal(t, flows.RunStatusFailed, session2.Runs()[1].Status())

	// change the UUID of the parent flow so it will effectively be missing
	assetsWithoutParentFlow := test.JSONReplace(assetsJSON, []string{"flows", "[0]", "uuid"}, []byte(`"653a3fa3-ff59-4a89-93c3-a8b9486ec479"`))

	sa, err = test.CreateSessionAssets(assetsWithoutParentFlow, "")
	require.NoError(t, err)

	session3, _, err := test.ResumeSession(session1, sa, "Hello")
	require.NoError(t, err)

	// should have an failed session
	assert.Equal(t, flows.SessionStatusFailed, session3.Status())
	assert.Equal(t, flows.RunStatusFailed, session3.Runs()[0].Status())
	assert.Equal(t, flows.RunStatusCompleted, session3.Runs()[1].Status())
}

func TestWaitTimeout(t *testing.T) {
	defer dates.SetNowFunc(time.Now)

	t1 := time.Date(2018, 4, 11, 13, 24, 30, 123456000, time.UTC)
	dates.SetNowFunc(dates.NewFixedNow(t1))

	_, session, sprint := test.NewSessionBuilder().WithAssetsPath("testdata/timeout_test.json").WithFlow("76f0a02f-3b75-4b86-9064-e9195e1b3a02").MustBuild()

	require.Equal(t, 1, len(session.Runs()[0].Path()))
	run := session.Runs()[0]

	require.Equal(t, 2, len(sprint.Events()))
	require.Equal(t, "msg_created", sprint.Events()[0].Type())
	require.Equal(t, "msg_wait", sprint.Events()[1].Type())

	// check our wait has a timeout
	waitEvent := run.Events()[1].(*events.MsgWaitEvent)
	require.Equal(t, 600, *waitEvent.TimeoutSeconds)

	_, err := session.Resume(context.Background(), resumes.NewWaitTimeout(events.NewWaitTimedOut()))
	require.NoError(t, err)

	require.Equal(t, flows.SessionStatusCompleted, session.Status())
	require.Equal(t, 2, len(run.Path()))
	require.Equal(t, 5, len(run.Events()))

	result := run.Results().Get("favorite_color")
	require.Equal(t, "Timeout", result.Category)
	require.Equal(t, "", result.Value)
	require.Equal(t, "", result.Input)
}

func TestCurrentContext(t *testing.T) {
	ctx := context.Background()

	_, session, _ := test.NewSessionBuilder().WithAssetsPath("../../test/testdata/runner/subflow_loop_with_wait.json").WithFlow("76f0a02f-3b75-4b86-9064-e9195e1b3a02").MustBuild()

	assert.Equal(t, string(flows.SessionStatusWaiting), string(session.Status()))

	context := session.CurrentContext()
	assert.NotNil(t, context)

	runContext, _ := context.Get("run")
	flowContext, _ := runContext.(*types.XObject).Get("flow")
	flowName, _ := flowContext.(*types.XObject).Get("name")
	assert.Equal(t, types.NewXText("Child flow"), flowName)

	// check we can marshal it
	_, err := jsonx.Marshal(context)
	assert.NoError(t, err)

	// end it
	session.Resume(ctx, resumes.NewWaitExpiration(events.NewWaitExpired()))
	assert.Equal(t, flows.SessionStatusCompleted, session.Status())

	// can still get context of completed session
	context = session.CurrentContext()
	assert.NotNil(t, context)

	runContext, _ = context.Get("run")
	flowContext, _ = runContext.(*types.XObject).Get("flow")
	flowName, _ = flowContext.(*types.XObject).Get("name")
	assert.Equal(t, types.NewXText("Parent Flow"), flowName)
}

func TestSessionHistory(t *testing.T) {
	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(`{
		"flows": [
			{
				"uuid": "5472a1c3-63e1-484f-8485-cc8ecb16a058",
				"name": "Empty",
				"spec_version": "13.1.0",
				"language": "eng",
				"type": "messaging",
				"nodes": []
			}
		]
	}`))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	flow := assets.NewFlowReference("5472a1c3-63e1-484f-8485-cc8ecb16a058", "Inception")
	contact := flows.NewEmptyContact(sa, "Bob", i18n.Language("eng"), nil)

	// trigger session manually which will have no history
	eng := engine.NewBuilder().Build()
	session1, _, err := eng.NewSession(context.Background(), sa, env, contact, triggers.NewBuilder(flow).Manual().Build(), nil)
	require.NoError(t, err)

	assert.Equal(t, flows.EmptyHistory, session1.History())

	// trigger another session from that session
	runSummary := session1.Runs()[0].Snapshot()
	runSummaryJSON := jsonx.MustMarshal(runSummary)
	history := flows.NewChildHistory(session1)

	session2, _, err := eng.NewSession(context.Background(), sa, env, contact, triggers.NewBuilder(flow).FlowAction(history, runSummaryJSON).Build(), nil)
	require.NoError(t, err)

	assert.Equal(t, &flows.SessionHistory{
		ParentUUID:          session1.UUID(),
		Ancestors:           1,
		AncestorsSinceInput: 1,
	}, session2.History())
}

func TestMaxResumesPerSession(t *testing.T) {
	ctx := context.Background()
	_, session, _ := test.NewSessionBuilder().WithAssetsPath("../../test/testdata/runner/two_questions.json").WithFlow("615b8a0f-588c-4d20-a05f-363b0b4ce6f4").MustBuild()
	require.Equal(t, flows.SessionStatusWaiting, session.Status())

	numResumes := 0
	for {
		msg := flows.NewMsgIn(flows.NewMsgUUID(), "tel:+593979123456", nil, "Teal", nil, "SMS1234")
		resume := resumes.NewMsg(events.NewMsgReceived(msg))
		numResumes++

		_, err := session.Resume(ctx, resume)
		require.NoError(t, err)

		if session.Status() == flows.SessionStatusFailed {
			break
		}
	}

	assert.Equal(t, 500, numResumes)
}

func TestFindStep(t *testing.T) {
	_, session, sprint := test.NewSessionBuilder().MustBuild()
	evts := sprint.Events()

	run, step := session.FindStep(evts[0].StepUUID())
	assert.Equal(t, "Registration", run.Flow().Name())
	assert.Equal(t, step.UUID(), evts[0].StepUUID())

	run, step = session.FindStep(flows.StepUUID("4f33917a-d562-4c20-88bd-f1a4c6827848"))
	assert.Nil(t, run)
	assert.Nil(t, step)
}

func TestEngineErrors(t *testing.T) {
	ctx := context.Background()

	// create a completed session and try to resume it
	_, session, _ := test.NewSessionBuilder().WithAssetsPath("../../test/testdata/runner/empty.json").WithFlow("76f0a02f-3b75-4b86-9064-e9195e1b3a02").MustBuild()
	require.Equal(t, flows.SessionStatusCompleted, session.Status())

	_, err := session.Resume(ctx, nil)
	assert.EqualError(t, err, "only waiting sessions can be resumed")
	assert.Equal(t, engine.ErrorResumeNonWaitingSession, err.(*engine.Error).Code())

	// create a session which is waiting for a message and try to resume it with a dial
	_, session, _ = test.NewSessionBuilder().MustBuild()
	require.Equal(t, flows.SessionStatusWaiting, session.Status())

	_, err = session.Resume(ctx, resumes.NewDial(events.NewDialEnded(flows.NewDial(flows.DialStatusAnswered, 10))))
	assert.EqualError(t, err, "resume of type dial not accepted by wait of type msg")
	assert.Equal(t, engine.ErrorResumeRejectedByWait, err.(*engine.Error).Code())
}
