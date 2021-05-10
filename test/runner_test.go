package test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/services/airtime/dtone"
	"github.com/nyaruka/goflow/services/email/smtp"
	"github.com/nyaruka/goflow/services/webhooks"
	"github.com/nyaruka/goflow/utils/smtpx"

	"github.com/pkg/errors"
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
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, errors.Wrap(err, "error reading test directory")
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

func marshalEventLog(eventLog []flows.Event) ([]json.RawMessage, error) {
	marshaled := make([]json.RawMessage, len(eventLog))
	var err error

	for i := range eventLog {
		marshaled[i], err = jsonx.Marshal(eventLog[i])
		if err != nil {
			return nil, errors.Wrap(err, "error marshaling event")
		}
	}
	return marshaled, nil
}

type Output struct {
	Session json.RawMessage   `json:"session"`
	Events  []json.RawMessage `json:"events"`
}

type FlowTest struct {
	Trigger   json.RawMessage      `json:"trigger"`
	Resumes   []json.RawMessage    `json:"resumes"`
	Outputs   []json.RawMessage    `json:"outputs"`
	HTTPMocks *httpx.MockRequestor `json:"http_mocks,omitempty"`
}

type runResult struct {
	session flows.Session
	outputs []*Output
}

func runFlow(assetsPath string, rawTrigger json.RawMessage, rawResumes []json.RawMessage) (runResult, error) {
	// load the test specific assets
	sa, err := LoadSessionAssets(envs.NewBuilder().Build(), assetsPath)
	if err != nil {
		return runResult{}, err
	}

	trigger, err := triggers.ReadTrigger(sa, rawTrigger, assets.PanicOnMissing)
	if err != nil {
		return runResult{}, errors.Wrapf(err, "error unmarshalling trigger")
	}

	eng := engine.NewBuilder().
		WithEmailServiceFactory(func(flows.Session) (flows.EmailService, error) {
			return smtp.NewService("smtp://nyaruka:pass123@mail.temba.io?from=flows@temba.io", nil)
		}).
		WithWebhookServiceFactory(webhooks.NewServiceFactory(http.DefaultClient, nil, nil, map[string]string{"User-Agent": "goflow-testing"}, 100000)).
		WithClassificationServiceFactory(func(s flows.Session, c *flows.Classifier) (flows.ClassificationService, error) {
			return newClassificationService(c), nil
		}).
		WithAirtimeServiceFactory(func(flows.Session) (flows.AirtimeService, error) {
			return dtone.NewService(http.DefaultClient, nil, "nyaruka", "123456789"), nil
		}).
		WithTicketServiceFactory(func(s flows.Session, t *flows.Ticketer) (flows.TicketService, error) {
			return NewTicketService(t), nil
		}).
		Build()

	session, sprint, err := eng.NewSession(sa, trigger)
	if err != nil {
		return runResult{}, err
	}

	outputs := make([]*Output, 0)

	// try to resume the session for each of the provided resumes
	for i, rawResume := range rawResumes {
		sessionJSON, err := jsonx.MarshalPretty(session)
		if err != nil {
			return runResult{}, errors.Wrap(err, "error marshalling output")
		}
		marshalledEvents, err := marshalEventLog(sprint.Events())
		if err != nil {
			return runResult{}, err
		}

		outputs = append(outputs, &Output{sessionJSON, marshalledEvents})

		session, err = eng.ReadSession(sa, sessionJSON, assets.PanicOnMissing)
		if err != nil {
			return runResult{}, errors.Wrap(err, "error marshalling output")
		}

		// if we aren't at a wait, that's an error
		if session.Wait() == nil {
			return runResult{}, errors.Errorf("did not stop at expected wait, have unused resumes: %d", len(rawResumes[i:]))
		}

		resume, err := resumes.ReadResume(sa, rawResume, assets.PanicOnMissing)
		if err != nil {
			return runResult{}, err
		}

		sprint, err = session.Resume(resume)
		if err != nil {
			return runResult{}, err
		}
	}

	sessionJSON, err := jsonx.MarshalPretty(session)
	if err != nil {
		return runResult{}, errors.Wrap(err, "error marshalling output")
	}

	marshalledEvents, err := marshalEventLog(sprint.Events())
	if err != nil {
		return runResult{}, err
	}

	outputs = append(outputs, &Output{sessionJSON, marshalledEvents})

	return runResult{session, outputs}, nil
}

func TestFlows(t *testing.T) {
	testCases, err := loadTestCases()
	require.NoError(t, err)
	require.True(t, len(testCases) > 0)

	defer uuids.SetGenerator(uuids.DefaultGenerator)
	defer dates.SetNowSource(dates.DefaultNowSource)
	defer httpx.SetRequestor(httpx.DefaultRequestor)
	defer smtpx.SetSender(smtpx.DefaultSender)

	for _, tc := range testCases {
		var httpMocksCopy *httpx.MockRequestor
		fmt.Printf("running %s\n", tc)

		uuids.SetGenerator(uuids.NewSeededGenerator(123456))
		dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2018, 7, 6, 12, 30, 0, 123456789, time.UTC)))
		smtpx.SetSender(smtpx.NewMockSender(nil, nil, nil, nil, nil, nil))

		testJSON, err := ioutil.ReadFile(tc.outputFile)
		require.NoError(t, err, "error reading output file %s", tc.outputFile)

		flowTest := &FlowTest{}
		err = jsonx.Unmarshal(json.RawMessage(testJSON), &flowTest)
		require.NoError(t, err, "error unmarshalling output file %s", tc.outputFile)

		if flowTest.HTTPMocks != nil {
			httpx.SetRequestor(flowTest.HTTPMocks)
			httpMocksCopy = flowTest.HTTPMocks.Clone()
		} else {
			httpx.SetRequestor(httpx.DefaultRequestor)
			httpMocksCopy = nil
		}

		// run our flow
		runResult, err := runFlow(tc.assetsFile, flowTest.Trigger, flowTest.Resumes)
		if err != nil {
			t.Errorf("error running flow for flow '%s' and output '%s': %s", tc.assetsFile, tc.outputFile, err)
			continue
		}

		if UpdateSnapshots {
			// we are writing new outputs, we write new files but don't test anything
			rawOutputs := make([]json.RawMessage, len(runResult.outputs))
			for i := range runResult.outputs {
				rawOutputs[i], err = jsonx.Marshal(runResult.outputs[i])
				require.NoError(t, err)
			}
			flowTest := &FlowTest{Trigger: flowTest.Trigger, Resumes: flowTest.Resumes, Outputs: rawOutputs, HTTPMocks: httpMocksCopy}
			testJSON, err := jsonx.MarshalPretty(flowTest)
			require.NoError(t, err, "Error marshalling test definition: %s", err)

			testJSON, _ = NormalizeJSON(testJSON)

			// write our output
			err = ioutil.WriteFile(tc.outputFile, testJSON, 0644)
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
				if !AssertEqualJSON(t, expected.Session, actual.Session, fmt.Sprintf("session is different in output[%d] in %s", i, tc)) {
					break
				}

				// and then each event
				for j := range actual.Events {
					if !AssertEqualJSON(t, expected.Events[j], actual.Events[j], fmt.Sprintf("event[%d] is different in output[%d] in %s", j, i, tc)) {
						break
					}
				}
			}
		}
	}
}

func BenchmarkFlows(b *testing.B) {
	testCases, _ := loadTestCases()

	for n := 0; n < b.N; n++ {
		for _, tc := range testCases {
			testJSON, err := ioutil.ReadFile(tc.outputFile)
			require.NoError(b, err, "error reading output file %s", tc.outputFile)

			flowTest := &FlowTest{}
			err = jsonx.Unmarshal(json.RawMessage(testJSON), &flowTest)
			require.NoError(b, err, "error unmarshalling output file %s", tc.outputFile)

			_, err = runFlow(tc.assetsFile, flowTest.Trigger, flowTest.Resumes)
			require.NoError(b, err, "error running flow %s", tc.testName)
		}
	}
}
