package es_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static/types"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/contactql/es"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils/jsonx"

	"github.com/olivere/elastic"
	"github.com/stretchr/testify/assert"
)

func newMockResolver() contactql.Resolver {
	return contactql.NewMockResolver(
		map[string]assets.Field{
			"age":      types.NewField("6b6a43fa-a26d-4017-bede-328bcdd5c93b", "age", "Age", assets.FieldTypeNumber),
			"color":    types.NewField("ecc7b13b-c698-4f46-8a90-24a8fab6fe34", "color", "Color", assets.FieldTypeText),
			"dob":      types.NewField("cbd3fc0e-9b74-4207-a8c7-248082bb4572", "dob", "DOB", assets.FieldTypeDatetime),
			"state":    types.NewField("67663ad1-3abc-42dd-a162-09df2dea66ec", "state", "State", assets.FieldTypeState),
			"district": types.NewField("54c72635-d747-4e45-883c-099d57dd998e", "district", "District", assets.FieldTypeDistrict),
			"ward":     types.NewField("fde8f740-c337-421b-8abb-83b954897c80", "ward", "Ward", assets.FieldTypeWard),
		},
		map[string]assets.Group{
			"u-reporters": types.NewGroup("8de30b78-d9ef-4db2-b2e8-4f7b6aef64cf", "U-Reporters", ""),
			"testers":     types.NewGroup("cf51cf8d-94da-447a-b27e-a42a900c37a6", "Testers", ""),
		},
	)
}

func TestElasticSort(t *testing.T) {
	resolver := newMockResolver()

	tcs := []struct {
		Label   string
		Sort    string
		Elastic string
		Error   error
	}{
		{"empty", "", `{"id":{"order":"desc"}}`, nil},
		{"descending created_on", "-created_on", `{"created_on":{"order":"desc"}}`, nil},
		{"ascending name", "name", `{"name.keyword":{"order":"asc"}}`, nil},
		{"descending language", "-language", `{"language":{"order":"desc"}}`, nil},
		{"descending numeric", "-AGE", `{"fields.number":{"nested":{"filter":{"term":{"fields.field":"6b6a43fa-a26d-4017-bede-328bcdd5c93b"}},"path":"fields"},"order":"desc"}}`, nil},
		{"ascending text", "color", `{"fields.text":{"nested":{"filter":{"term":{"fields.field":"ecc7b13b-c698-4f46-8a90-24a8fab6fe34"}},"path":"fields"},"order":"asc"}}`, nil},
		{"descending date", "-dob", `{"fields.datetime":{"nested":{"filter":{"term":{"fields.field":"cbd3fc0e-9b74-4207-a8c7-248082bb4572"}},"path":"fields"},"order":"desc"}}`, nil},
		{"descending state", "-state", `{"fields.state_keyword":{"nested":{"filter":{"term":{"fields.field":"67663ad1-3abc-42dd-a162-09df2dea66ec"}},"path":"fields"},"order":"desc"}}`, nil},
		{"ascending district", "district", `{"fields.district_keyword":{"nested":{"filter":{"term":{"fields.field":"54c72635-d747-4e45-883c-099d57dd998e"}},"path":"fields"},"order":"asc"}}`, nil},
		{"ascending ward", "ward", `{"fields.ward_keyword":{"nested":{"filter":{"term":{"fields.field":"fde8f740-c337-421b-8abb-83b954897c80"}},"path":"fields"},"order":"asc"}}`, nil},

		{"unknown field", "foo", "", fmt.Errorf("unable to find field with name: foo")},
	}

	for _, tc := range tcs {
		sort, err := es.ToElasticFieldSort(resolver, tc.Sort)

		if err != nil {
			assert.Equal(t, tc.Error.Error(), err.Error())
			continue
		}

		src, _ := sort.Source()
		encoded, _ := jsonx.Marshal(src)
		assert.Equal(t, tc.Elastic, string(encoded))
	}
}

func TestElasticQuery(t *testing.T) {
	resolver := newMockResolver()

	type TestCase struct {
		Description string          `json:"description"`
		Query       string          `json:"query"`
		Elastic     json.RawMessage `json:"elastic"`
		Error       string          `json:"error"`
		IsAnon      bool            `json:"is_anon"`
	}
	tcs := make([]TestCase, 0, 20)
	tcJSON, err := ioutil.ReadFile("testdata/to_query.json")
	assert.NoError(t, err)

	err = json.Unmarshal(tcJSON, &tcs)
	assert.NoError(t, err)

	ny, _ := time.LoadLocation("America/New_York")

	for _, tc := range tcs {
		testName := fmt.Sprintf("test '%s' for query '%s'", tc.Description, tc.Query)

		redactionPolicy := envs.RedactionPolicyNone
		if tc.IsAnon {
			redactionPolicy = envs.RedactionPolicyURNs
		}
		env := envs.NewBuilder().WithTimezone(ny).WithRedactionPolicy(redactionPolicy).Build()

		qlQuery, err := contactql.ParseQuery(tc.Query, env.RedactionPolicy(), env.DefaultCountry(), resolver)

		var query elastic.Query
		if err == nil {
			query, err = es.ToElasticQuery(env, resolver, qlQuery)
		}

		if tc.Error != "" {
			assert.Error(t, err, "expected error in %s", testName)
			if err != nil {
				assert.Contains(t, err.Error(), tc.Error)
			}
			continue
		}

		assert.NoError(t, err, "unexpected error in %s", testName)
		if err != nil {
			continue
		}

		assert.NotNil(t, query, tc.Description)
		if query == nil {
			continue
		}

		source, err := query.Source()
		assert.NoError(t, err, "error reqesting source for elastic query in ", testName)
		if err != nil {
			continue
		}

		asJSON, err := jsonx.Marshal(source)
		assert.NoError(t, err, "error marshaling elastic query in ", testName)
		if err != nil {
			continue
		}

		compacted := &bytes.Buffer{}
		json.Compact(compacted, tc.Elastic)

		assert.Equal(t, compacted.String(), string(asJSON), "elastic query mismatch in ", testName)
	}
}
