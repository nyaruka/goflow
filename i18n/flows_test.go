package i18n_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/i18n"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/dates"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractFromFlows(t *testing.T) {
	defer dates.SetNowSource(dates.DefaultNowSource)
	dates.SetNowSource(dates.NewFixedNowSource(time.Date(2020, 3, 25, 13, 57, 30, 123456789, time.UTC)))

	tests := []struct {
		assets      string
		flowUUIDs   []assets.FlowUUID
		lang        envs.Language
		excludeArgs bool
		po          string
	}{
		{
			"../test/testdata/runner/two_questions.json",
			[]assets.FlowUUID{assets.FlowUUID(`615b8a0f-588c-4d20-a05f-363b0b4ce6f4`)},
			envs.NilLanguage,
			false,
			"two_questions.en.po",
		},
		{
			"../test/testdata/runner/two_questions.json",
			[]assets.FlowUUID{assets.FlowUUID(`615b8a0f-588c-4d20-a05f-363b0b4ce6f4`)},
			envs.Language(`fra`),
			false,
			"two_questions.fr.po",
		},
		{
			"../test/testdata/runner/two_questions.json",
			[]assets.FlowUUID{assets.FlowUUID(`615b8a0f-588c-4d20-a05f-363b0b4ce6f4`)},
			envs.Language(`fra`),
			true,
			"two_questions.noargs.fr.po",
		},
		{
			"testdata/translation_mismatches.json",
			[]assets.FlowUUID{assets.FlowUUID(`19cad1f2-9110-4271-98d4-1b968bf19410`)},
			envs.Language(`spa`),
			false,
			"translation_mismatches.noargs.es.po",
		},
	}

	for _, tc := range tests {
		env := envs.NewBuilder().Build()
		sa, err := test.LoadSessionAssets(env, tc.assets)
		require.NoError(t, err)

		sources := make([]flows.Flow, 0)
		for _, flowUUID := range tc.flowUUIDs {
			flow, err := sa.Flows().Get(flowUUID)
			require.NoError(t, err)
			sources = append(sources, flow)
		}

		po, err := i18n.ExtractFromFlows("Testing", tc.lang, tc.excludeArgs, sources...)
		assert.NoError(t, err)

		b := &strings.Builder{}
		po.Write(b)
		poAsStr := b.String()

		if !test.UpdateSnapshots {
			expected, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s", tc.po))
			require.NoError(t, err)

			assert.Equal(t, string(expected), poAsStr)
		} else {
			ioutil.WriteFile(fmt.Sprintf("testdata/%s", tc.po), []byte(poAsStr), 0666)
		}
	}
}
