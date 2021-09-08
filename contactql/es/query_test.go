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

func newMockResolver() contactql.Resolver {
	return contactql.NewMockResolver(
		map[string]assets.Field{
			"age":      static.NewField("6b6a43fa-a26d-4017-bede-328bcdd5c93b", "age", "Age", assets.FieldTypeNumber),
			"color":    static.NewField("ecc7b13b-c698-4f46-8a90-24a8fab6fe34", "color", "Color", assets.FieldTypeText),
			"dob":      static.NewField("cbd3fc0e-9b74-4207-a8c7-248082bb4572", "dob", "DOB", assets.FieldTypeDatetime),
			"state":    static.NewField("67663ad1-3abc-42dd-a162-09df2dea66ec", "state", "State", assets.FieldTypeState),
			"district": static.NewField("54c72635-d747-4e45-883c-099d57dd998e", "district", "District", assets.FieldTypeDistrict),
			"ward":     static.NewField("fde8f740-c337-421b-8abb-83b954897c80", "ward", "Ward", assets.FieldTypeWard),
		},
		map[string]assets.Group{
			"u-reporters": static.NewGroup("8de30b78-d9ef-4db2-b2e8-4f7b6aef64cf", "U-Reporters", ""),
			"testers":     static.NewGroup("cf51cf8d-94da-447a-b27e-a42a900c37a6", "Testers", ""),
		},
	)
}

func TestElasticQuery(t *testing.T) {
	resolver := newMockResolver()

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

	for _, tc := range tcs {
		testName := fmt.Sprintf("test '%s' for query '%s'", tc.Description, tc.Query)

		redactionPolicy := envs.RedactionPolicyNone
		if tc.RedactURNs {
			redactionPolicy = envs.RedactionPolicyURNs
		}
		env := envs.NewBuilder().WithTimezone(ny).WithRedactionPolicy(redactionPolicy).Build()

		parsed, err := contactql.ParseQuery(env, tc.Query, resolver)
		require.NoError(t, err)

		query := es.ToElasticQuery(env, parsed)
		assert.NotNil(t, query, tc.Description)

		source, err := query.Source()
		require.NoError(t, err, "error requesting source for elastic query in ", testName)

		asJSON, err := jsonx.Marshal(source)
		require.NoError(t, err)

		test.AssertEqualJSON(t, tc.Elastic, asJSON, "elastic mismatch in %s", testName)
	}
}
