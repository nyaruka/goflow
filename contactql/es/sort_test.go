package es_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/contactql/es"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestElasticSort(t *testing.T) {
	resolver := newMockResolver()

	type testCase struct {
		Description string          `json:"description"`
		SortBy      string          `json:"sort_by"`
		Elastic     json.RawMessage `json:"elastic,omitempty"`
		Error       string          `json:"error,omitempty"`
	}
	tcs := make([]testCase, 0, 20)
	tcJSON, err := os.ReadFile("testdata/to_sort.json")
	require.NoError(t, err)

	err = json.Unmarshal(tcJSON, &tcs)
	require.NoError(t, err)

	for _, tc := range tcs {
		sort, err := es.ToElasticFieldSort(tc.SortBy, resolver)

		if tc.Error != "" {
			assert.EqualError(t, err, tc.Error)
		} else {
			src, _ := sort.Source()
			encoded := jsonx.MustMarshal(src)
			test.AssertEqualJSON(t, []byte(tc.Elastic), encoded, "field sort mismatch for %s", tc.Description)
		}
	}
}
