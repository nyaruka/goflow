package test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"text/template"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/services/email/smtp"
	"github.com/nyaruka/goflow/services/webhooks"
	"github.com/nyaruka/goflow/test/services"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/smtpx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var includeTests string
var testFilePattern = regexp.MustCompile(`(\w+)\.(\w+)\.json`)

func init() {
	flag.StringVar(&includeTests, "include", "", "include only test names containing")
}

type runnerTest struct {
	testName   string
	assetsName string
	outputFile string
	assetsFile string
}

func (t runnerTest) String() string {
	return fmt.Sprintf("%s.%s", t.assetsName, t.testName)
}

func loadTestCases() ([]runnerTest, error) {
	directory := "testdata/runner/"
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("error reading test directory: %w", err)
	}

	tests := make([]runnerTest, 0)

	for _, file := range files {
		groups := testFilePattern.FindStringSubmatch(file.Name())
		if groups != nil {
			testName := groups[2]
			assetsName := groups[1]
			assetsFile := directory + assetsName + ".json"
			outputFile := directory + groups[0]

			if includeTests == "" || strings.Contains(assetsName+"."+testName, includeTests) {
				tests = append(tests, runnerTest{testName, assetsName, outputFile, assetsFile})
			}
		}
	}

	return tests, nil
}

func marshalEventLog(eventLog []flows.Event) []json.RawMessage {
	marshaled := make([]json.RawMessage, len(eventLog))

	for i := range eventLog {
		marshaled[i] = jsonx.MustMarshal(eventLog[i])
	}
	return marshaled
}

type Output struct {
	Session  json.RawMessage   `json:"session"`
	Events   []json.RawMessage `json:"events"`
	Segments json.RawMessage   `json:"segments"`
}

type FlowTest struct {
	Environment json.RawMessage        `json:"environment"`
	Contact     *flows.ContactEnvelope `json:"contact"`
	Trigger     json.RawMessage        `json:"trigger"`
	Call        *flows.CallEnvelope    `json:"call,omitempty"`
	Resumes     []json.RawMessage      `json:"resumes"`
	Outputs     []json.RawMessage      `json:"outputs"`
	HTTPMocks   *httpx.MockRequestor   `json:"http_mocks,omitempty"`
}

type runResult struct {
	session flows.Session
	outputs []*Output
}

func runFlow(assetsPath string, rawEnv []byte, rawContact *flows.ContactEnvelope, rawTrigger []byte, rawCall *flows.CallEnvelope, rawResumes []json.RawMessage) (runResult, error) {
	ctx := context.Background()

	// load the test specific assets
	sa, err := LoadSessionAssets(envs.NewBuilder().Build(), assetsPath)
	if err != nil {
		return runResult{}, err
	}

	env, err := envs.ReadEnvironment(rawEnv)
	if err != nil {
		return runResult{}, fmt.Errorf("error unmarshalling environment: %w", err)
	}

	contact, err := rawContact.Unmarshal(sa, assets.PanicOnMissing)
	if err != nil {
		return runResult{}, fmt.Errorf("error unmarshalling contact: %w", err)
	}

	trigger, err := triggers.Read(sa, rawTrigger, assets.PanicOnMissing)
	if err != nil {
		return runResult{}, fmt.Errorf("error unmarshalling trigger: %w", err)
	}

	eng := engine.NewBuilder().
		WithLLMPrompts(map[string]*template.Template{
			"categorize": template.Must(template.New("").Parse("Categorize the following text into one of the following categories and only return that category or <CANT> if you can't: {{ .arg1 }}")),
		}).
		WithEmailServiceFactory(func(flows.SessionAssets) (flows.EmailService, error) {
			return smtp.NewService("smtp://nyaruka:pass123@mail.temba.io?from=flows@temba.io", nil)
		}).
		WithWebhookServiceFactory(webhooks.NewServiceFactory(http.DefaultClient, nil, nil, map[string]string{"User-Agent": "goflow-testing"}, 100000)).
		WithClassificationServiceFactory(func(c *flows.Classifier) (flows.ClassificationService, error) {
			return services.NewClassification(c), nil
		}).
		WithLLMServiceFactory(func(l *flows.LLM) (flows.LLMService, error) {
			return services.NewLLM(), nil
		}).
		WithAirtimeServiceFactory(func(flows.SessionAssets) (flows.AirtimeService, error) {
			return services.NewAirtime("RWF"), nil
		}).
		Build()

	var call *flows.Call
	if rawCall != nil {
		call = rawCall.Unmarshal(sa, assets.PanicOnMissing)
	}

	session, sprint, err := eng.NewSession(ctx, sa, env, contact, trigger, call)
	if err != nil {
		return runResult{}, err
	}

	outputs := make([]*Output, 0)

	// try to resume the session for each of the provided resumes
	for i, rawResume := range rawResumes {
		sessionJSON, err := jsonx.MarshalPretty(session)
		if err != nil {
			return runResult{}, fmt.Errorf("error marshalling output: %w", err)
		}

		outputs = append(outputs, &Output{
			Session:  sessionJSON,
			Events:   marshalEventLog(sprint.Events()),
			Segments: jsonx.MustMarshal(sprint.Segments()),
		})

		session, err = eng.ReadSession(sa, sessionJSON, env, contact, call, assets.PanicOnMissing)
		if err != nil {
			return runResult{}, fmt.Errorf("error marshalling output: %w", err)
		}

		// if session isn't waiting for another resume, that's an error
		if session.Status() != flows.SessionStatusWaiting {
			return runResult{}, fmt.Errorf("did not stop at expected wait, have unused resumes: %d", len(rawResumes[i:]))
		}

		resume, err := resumes.Read(sa, rawResume, assets.PanicOnMissing)
		if err != nil {
			return runResult{}, err
		}

		sprint, err = session.Resume(ctx, resume)
		if err != nil {
			return runResult{}, err
		}
	}

	sessionJSON, err := jsonx.MarshalPretty(session)
	if err != nil {
		return runResult{}, fmt.Errorf("error marshalling output: %w", err)
	}

	outputs = append(outputs, &Output{
		Session:  sessionJSON,
		Events:   marshalEventLog(sprint.Events()),
		Segments: jsonx.MustMarshal(sprint.Segments()),
	})

	return runResult{session, outputs}, nil
}

func TestFlows(t *testing.T) {
	testCases, err := loadTestCases()
	require.NoError(t, err)
	require.True(t, len(testCases) > 0)

	defer httpx.SetRequestor(httpx.DefaultRequestor)
	defer smtpx.SetSender(smtpx.DefaultSender)

	for _, tc := range testCases {
		MockUniverse()

		var httpMocksCopy *httpx.MockRequestor

		smtpx.SetSender(smtpx.NewMockSender(nil, nil, nil, nil, nil, nil))

		testJSON, err := os.ReadFile(tc.outputFile)
		require.NoError(t, err, "error reading output file %s", tc.outputFile)

		flowTest := &FlowTest{}
		err = utils.UnmarshalAndValidate([]byte(testJSON), flowTest)
		require.NoError(t, err, "error unmarshalling output file %s", tc.outputFile)

		if flowTest.HTTPMocks != nil {
			httpx.SetRequestor(flowTest.HTTPMocks)
			httpMocksCopy = flowTest.HTTPMocks.Clone()
		} else {
			httpx.SetRequestor(httpx.DefaultRequestor)
			httpMocksCopy = nil
		}

		// run our flow
		runResult, err := runFlow(tc.assetsFile, flowTest.Environment, flowTest.Contact, flowTest.Trigger, flowTest.Call, flowTest.Resumes)
		if err != nil {
			t.Errorf("error running flow for flow '%s' and output '%s': %s", tc.assetsFile, tc.outputFile, err)
			continue
		}

		// check all http mocks were used
		if flowTest.HTTPMocks != nil {
			require.False(t, flowTest.HTTPMocks.HasUnused(), "unused HTTP mocks for flow '%s' and output '%s'", tc.assetsFile, tc.outputFile)
		}

		if UpdateSnapshots {
			// we are writing new outputs, we write new files but don't test anything
			rawOutputs := make([]json.RawMessage, len(runResult.outputs))
			for i := range runResult.outputs {
				rawOutputs[i], err = jsonx.Marshal(runResult.outputs[i])
				require.NoError(t, err)
			}
			flowTest := &FlowTest{Contact: flowTest.Contact, Trigger: flowTest.Trigger, Call: flowTest.Call, Resumes: flowTest.Resumes, Outputs: rawOutputs, HTTPMocks: httpMocksCopy}
			testJSON, err := jsonx.MarshalPretty(flowTest)
			require.NoError(t, err, "Error marshalling test definition: %s", err)

			testJSON, _ = NormalizeJSON(testJSON)

			// write our output
			err = os.WriteFile(tc.outputFile, testJSON, 0644)
			require.NoError(t, err, "Error writing test file to %s: %s", tc.outputFile, err)
		} else {
			// start by checking we have the expected number of outputs
			if !assert.Equal(t, len(flowTest.Outputs), len(runResult.outputs), "wrong number of outputs in %s", tc) {
				continue
			}

			// then check each output
			for i, actual := range runResult.outputs {
				// unmarshal our expected outputsinto session+events
				expected := &Output{}
				err := jsonx.Unmarshal(flowTest.Outputs[i], expected)
				require.NoError(t, err, "error unmarshalling output")

				// first the session
				if !AssertEqualJSON(t, expected.Session, actual.Session, "session is different in output[%d] in %s", i, tc) {
					break
				}

				// and then each event
				for j := range actual.Events {
					if !AssertEqualJSON(t, expected.Events[j], actual.Events[j], "event[%d] is different in output[%d] in %s", j, i, tc) {
						break
					}
				}

				// and finally the path segments
				if !AssertEqualJSON(t, expected.Segments, actual.Segments, "segments are different in output[%d] in %s", i, tc) {
					break
				}
			}
		}
	}
}
