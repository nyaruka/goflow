package contactql_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInspect(t *testing.T) {
	tests := []struct {
		query      string
		inspection *contactql.Inspection
	}{
		{
			query: "bob",
			inspection: &contactql.Inspection{
				Attributes:   []string{"name"},
				Schemes:      []string{},
				Fields:       []*assets.FieldReference{},
				Groups:       []*assets.GroupReference{},
				AllowAsGroup: true,
			},
		},
		{
			query: "age > 18 AND name != \"\" OR twitter = bobby OR tel ~1234",
			inspection: &contactql.Inspection{
				Attributes: []string{"name"},
				Schemes:    []string{"tel", "twitter"},
				Fields: []*assets.FieldReference{
					assets.NewFieldReference("age", ""),
				},
				Groups:       []*assets.GroupReference{},
				AllowAsGroup: true,
			},
		},
		{
			query: "id = 123",
			inspection: &contactql.Inspection{
				Attributes:   []string{"id"},
				Schemes:      []string{},
				Fields:       []*assets.FieldReference{},
				Groups:       []*assets.GroupReference{},
				AllowAsGroup: false,
			},
		},
		{
			query: "group = U-reporters",
			inspection: &contactql.Inspection{
				Attributes: []string{"group"},
				Schemes:    []string{},
				Fields:     []*assets.FieldReference{},
				Groups: []*assets.GroupReference{
					assets.NewVariableGroupReference("U-reporters"),
				},
				AllowAsGroup: false,
			},
		},
	}

	env := envs.NewBuilder().Build()

	for _, tc := range tests {
		query, err := contactql.ParseQuery(env, tc.query)
		require.NoError(t, err, "error parsing %s", tc.query)

		assert.Equal(t, tc.inspection, contactql.Inspect(query), "inspect mismatch for query %s", tc.query)
	}

}
