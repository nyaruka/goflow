package es_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/contactql/es"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockMapper struct {
	flows  map[assets.FlowUUID]int64
	groups map[assets.GroupUUID]int64
}

func (m *MockMapper) Flow(f assets.Flow) int64 {
	return m.flows[f.UUID()]
}

func (m *MockMapper) Group(g assets.Group) int64 {
	return m.groups[g.UUID()]
}

func newMockMapper(flows map[assets.FlowUUID]int64, groups map[assets.GroupUUID]int64) *MockMapper {
	return &MockMapper{flows, groups}
}

func newMockResolver() contactql.Resolver {
	return contactql.NewMockResolver(
		[]assets.Field{
			static.NewField("6b6a43fa-a26d-4017-bede-328bcdd5c93b", "age", "Age", assets.FieldTypeNumber),
			static.NewField("ecc7b13b-c698-4f46-8a90-24a8fab6fe34", "color", "Color", assets.FieldTypeText),
			static.NewField("cbd3fc0e-9b74-4207-a8c7-248082bb4572", "dob", "DOB", assets.FieldTypeDatetime),
			static.NewField("67663ad1-3abc-42dd-a162-09df2dea66ec", "state", "State", assets.FieldTypeState),
			static.NewField("54c72635-d747-4e45-883c-099d57dd998e", "district", "District", assets.FieldTypeDistrict),
			static.NewField("fde8f740-c337-421b-8abb-83b954897c80", "ward", "Ward", assets.FieldTypeWard),
		},
		[]assets.Flow{
			static.NewFlow("c261165a-f5b0-40ba-b916-76fb49667a4f", "Registration", []byte(`{}`)),
		},
		[]assets.Group{
			static.NewGroup("8de30b78-d9ef-4db2-b2e8-4f7b6aef64cf", "U-Reporters", ""),
			static.NewGroup("cf51cf8d-94da-447a-b27e-a42a900c37a6", "Testers", ""),
		},
	)
}

func TestElasticQuery(t *testing.T) {
	resolver := newMockResolver()
	mapper := newMockMapper(
		map[assets.FlowUUID]int64{
			"c261165a-f5b0-40ba-b916-76fb49667a4f": 234, // Registration
		},
		map[assets.GroupUUID]int64{
			"8de30b78-d9ef-4db2-b2e8-4f7b6aef64cf": 345, // U-Reporters
			"cf51cf8d-94da-447a-b27e-a42a900c37a6": 456, // Testers
		},
	)

	type testCase struct {
		Description string          `json:"description"`
		Query       string          `json:"query"`
		Elastic     json.RawMessage `json:"elastic"`
		RedactURNs  bool            `json:"redact_urns"`
	}
	tcs := make([]testCase, 0, 20)
	tcJSON, err := os.ReadFile("testdata/to_query.json")
	require.NoError(t, err)
	jsonx.MustUnmarshal(tcJSON, &tcs)

	ny, _ := time.LoadLocation("America/New_York")

	for i, tc := range tcs {
		testName := fmt.Sprintf("test '%s' for query '%s'", tc.Description, tc.Query)

		redactionPolicy := envs.RedactionPolicyNone
		if tc.RedactURNs {
			redactionPolicy = envs.RedactionPolicyURNs
		}
		env := envs.NewBuilder().WithTimezone(ny).WithRedactionPolicy(redactionPolicy).Build()

		parsed, err := contactql.ParseQuery(env, tc.Query, resolver)
		require.NoError(t, err)

		conv := es.NewConverter(env, mapper, false)
		query := conv.Query(parsed)
		assert.NotNil(t, query, tc.Description)

		asJSON, err := jsonx.Marshal(query)
		require.NoError(t, err)

		// clone test case and populate with actual values
		actual := tc
		actual.Elastic = asJSON

		if !test.UpdateSnapshots {
			test.AssertEqualJSON(t, tc.Elastic, actual.Elastic, "elastic mismatch in %s", testName)
		} else {
			tcs[i] = actual
		}
	}

	if test.UpdateSnapshots {
		actualJSON, err := jsonx.MarshalPretty(tcs)
		require.NoError(t, err)

		err = os.WriteFile("testdata/to_query.json", actualJSON, 0666)
		require.NoError(t, err)
	}
}

func TestElasticQueryUUIDAsDocID(t *testing.T) {
	resolver := newMockResolver()
	mapper := newMockMapper(
		map[assets.FlowUUID]int64{},
		map[assets.GroupUUID]int64{},
	)
	ny, _ := time.LoadLocation("America/New_York")
	env := envs.NewBuilder().WithTimezone(ny).Build()

	conv := es.NewConverter(env, mapper, true)

	// uuid = X should query _id
	parsed, err := contactql.ParseQuery(env, `uuid = "f81d4fae-7dec-11d0-a765-00a0c91e6bf6"`, resolver)
	require.NoError(t, err)
	asJSON := jsonx.MustMarshal(conv.Query(parsed))
	test.AssertEqualJSON(t, []byte(`{"ids":{"values":["f81d4fae-7dec-11d0-a765-00a0c91e6bf6"]}}`), asJSON, "uuid query mismatch")

	// uuid != X should query _id with NOT
	parsed, err = contactql.ParseQuery(env, `uuid != "f81d4fae-7dec-11d0-a765-00a0c91e6bf6"`, resolver)
	require.NoError(t, err)
	asJSON = jsonx.MustMarshal(conv.Query(parsed))
	test.AssertEqualJSON(t, []byte(`{"bool":{"must_not":{"ids":{"values":["f81d4fae-7dec-11d0-a765-00a0c91e6bf6"]}}}}`), asJSON, "uuid != query mismatch")

	// id = X should query id field
	parsed, err = contactql.ParseQuery(env, `id = 123`, resolver)
	require.NoError(t, err)
	asJSON = jsonx.MustMarshal(conv.Query(parsed))
	test.AssertEqualJSON(t, []byte(`{"term":{"id":{"value":"123"}}}`), asJSON, "id query mismatch")

	// id != X should query id field with NOT
	parsed, err = contactql.ParseQuery(env, `id != 123`, resolver)
	require.NoError(t, err)
	asJSON = jsonx.MustMarshal(conv.Query(parsed))
	test.AssertEqualJSON(t, []byte(`{"bool":{"must_not":{"term":{"id":{"value":"123"}}}}}`), asJSON, "id != query mismatch")
}
