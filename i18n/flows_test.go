package i18n_test

import (
	"encoding/json"
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
		assets       string
		flowUUIDs    []assets.FlowUUID
		lang         envs.Language
		excludeProps []string
		po           string
	}{
		{
			"../test/testdata/runner/two_questions.json",
			[]assets.FlowUUID{assets.FlowUUID(`615b8a0f-588c-4d20-a05f-363b0b4ce6f4`)},
			envs.NilLanguage,
			nil,
			"two_questions.po",
		},
		{
			"../test/testdata/runner/two_questions.json",
			[]assets.FlowUUID{assets.FlowUUID(`615b8a0f-588c-4d20-a05f-363b0b4ce6f4`)},
			envs.Language("eng"),
			nil,
			"two_questions.en.po",
		},
		{
			"../test/testdata/runner/two_questions.json",
			[]assets.FlowUUID{assets.FlowUUID(`615b8a0f-588c-4d20-a05f-363b0b4ce6f4`)},
			envs.Language(`fra`),
			nil,
			"two_questions.fr.po",
		},
		{
			"../test/testdata/runner/two_questions.json",
			[]assets.FlowUUID{assets.FlowUUID(`615b8a0f-588c-4d20-a05f-363b0b4ce6f4`)},
			envs.Language(`fra`),
			[]string{"arguments"},
			"two_questions.noargs.fr.po",
		},
		{
			"testdata/translation_mismatches.json",
			[]assets.FlowUUID{assets.FlowUUID(`19cad1f2-9110-4271-98d4-1b968bf19410`)},
			envs.Language(`spa`),
			nil,
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

		po, err := i18n.ExtractFromFlows("Testing", tc.lang, tc.excludeProps, sources...)
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

func TestExtractFromFlowsWithDiffLanguages(t *testing.T) {
	sa, err := test.LoadSessionAssets(envs.NewBuilder().Build(), "testdata/different_languages.json")
	require.NoError(t, err)

	engFlow, _ := sa.Flows().Get("76f0a02f-3b75-4b86-9064-e9195e1b3a02")
	spaFlow, _ := sa.Flows().Get("e9e1d54f-f213-44ca-883a-eb96d15151aa")

	_, err = i18n.ExtractFromFlows("", "fra", nil, engFlow, spaFlow)
	assert.EqualError(t, err, "can't extract from flows with differing base languages")
}

func TestImportIntoFlows(t *testing.T) {
	sa, err := test.LoadSessionAssets(envs.NewBuilder().Build(), "testdata/translation_mismatches.json")
	require.NoError(t, err)

	flow, err := sa.Flows().Get("19cad1f2-9110-4271-98d4-1b968bf19410")
	require.NoError(t, err)

	po := i18n.NewPO(nil)

	// all instances of "Red" should be translated as "Rojo" (ignores the one which is already "Rojo")
	po.AddEntry(&i18n.POEntry{
		MsgID:  "Red",
		MsgStr: "Rojo",
	})

	// all instances of "Blue" should be translated as "Azul oscura"
	po.AddEntry(&i18n.POEntry{
		MsgID:  "Blue",
		MsgStr: "Azul oscura",
	})

	// except the quick reply instance of "Blue" should be translated as "Azul clara"
	po.AddEntry(&i18n.POEntry{
		MsgContext: "e42deebf-90fa-4636-81cb-d247a3d3ba75/quick_replies:1",
		MsgID:      "Blue",
		MsgStr:     "Azul clara",
	})

	updates := i18n.CalculateFlowUpdates(po, envs.Language("spa"), flow)
	assert.Equal(t, 3, len(updates))
	assert.Equal(t, `Translated/d1ce3c92-7025-4607-a910-444361a6b9b3/name:0 "Roja" -> "Rojo"`, updates[0].String())
	assert.Equal(t, `Translated/e42deebf-90fa-4636-81cb-d247a3d3ba75/quick_replies:1 "Azul" -> "Azul clara"`, updates[1].String())
	assert.Equal(t, `Translated/43f7e69e-727d-4cfe-81b8-564e7833052b/name:0 "Azul" -> "Azul oscura"`, updates[2].String())

	err = i18n.ImportIntoFlows(po, envs.Language("spa"), flow)
	require.NoError(t, err)

	localJSON, _ := json.Marshal(flow.Localization())
	test.AssertEqualJSON(t, []byte(`{
		"spa": {
			"e42deebf-90fa-4636-81cb-d247a3d3ba75": {
				"text": [
					"Cual pastilla?"
				],
				"quick_replies": [
					"Rojo",
					"Azul clara"
				]
			},
			"d1ce3c92-7025-4607-a910-444361a6b9b3": {
				"name": [
					"Rojo"
				]
			},
			"43f7e69e-727d-4cfe-81b8-564e7833052b": {
				"name": [
					"Azul oscura"
				]
			},
			"3a044264-81d1-4ba7-882a-e98740c8e724": {
				"name": [
					"Otro"
				]
			},
			"61bc5ed3-e216-4457-8ce5-ad658e697f29": {
				"arguments": [
					"rojo"
				]
			},
			"5f5fa09f-bf88-4719-ba64-cab9cf2f67b5": {
				"arguments": [
					"azul"
				]
			}
		}
	}`), localJSON, "post-import localization mismatch")
}

func TestImportIntoFlowsWithDiffLanguages(t *testing.T) {
	sa, err := test.LoadSessionAssets(envs.NewBuilder().Build(), "testdata/different_languages.json")
	require.NoError(t, err)

	engFlow, _ := sa.Flows().Get("76f0a02f-3b75-4b86-9064-e9195e1b3a02")
	spaFlow, _ := sa.Flows().Get("e9e1d54f-f213-44ca-883a-eb96d15151aa")

	err = i18n.ImportIntoFlows(nil, "fra", engFlow, spaFlow)
	assert.EqualError(t, err, "can't import into flows with differing base languages")
}
