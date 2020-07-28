package contactql_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static/types"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInspect(t *testing.T) {
	tests := []struct {
		Query      string
		Inspection *contactql.Inspection
	}{
		{
			Query: "bob",
			Inspection: &contactql.Inspection{
				Attributes:   []string{"name"},
				Schemes:      []string{},
				Fields:       []*assets.FieldReference{},
				Groups:       []*assets.GroupReference{},
				AllowAsGroup: true,
			},
		},
		{
			Query: "age > 18 AND name != \"\" OR twitter = bobby OR tel ~1234",
			Inspection: &contactql.Inspection{
				Attributes: []string{"name"},
				Schemes:    []string{"tel", "twitter"},
				Fields: []*assets.FieldReference{
					assets.NewFieldReference("age", "Age"),
				},
				Groups:       []*assets.GroupReference{},
				AllowAsGroup: true,
			},
		},
		{
			Query: "id = 123",
			Inspection: &contactql.Inspection{
				Attributes:   []string{"id"},
				Schemes:      []string{},
				Fields:       []*assets.FieldReference{},
				Groups:       []*assets.GroupReference{},
				AllowAsGroup: false,
			},
		},
		{
			Query: "group = U-reporters",
			Inspection: &contactql.Inspection{
				Attributes: []string{"group"},
				Schemes:    []string{},
				Fields:     []*assets.FieldReference{},
				Groups: []*assets.GroupReference{
					assets.NewGroupReference(assets.GroupUUID("4eeca453-f474-4767-bdd0-434b180223db"), "U-Reporters"),
				},
				AllowAsGroup: false,
			},
		},
	}

	env := envs.NewBuilder().Build()
	resolver := contactql.NewMockResolver(map[string]assets.Field{
		"age":    types.NewField(assets.FieldUUID("f1b5aea6-6586-41c7-9020-1a6326cc6565"), "age", "Age", assets.FieldTypeNumber),
		"dob":    types.NewField(assets.FieldUUID("3810a485-3fda-4011-a589-7320c0b8dbef"), "dob", "DOB", assets.FieldTypeDatetime),
		"gender": types.NewField(assets.FieldUUID("d66a7823-eada-40e5-9a3a-57239d4690bf"), "gender", "Gender", assets.FieldTypeText),
	}, map[string]assets.Group{
		"u-reporters": types.NewGroup(assets.GroupUUID("4eeca453-f474-4767-bdd0-434b180223db"), "U-Reporters", ""),
	})

	for _, tc := range tests {
		query, err := contactql.ParseQuery(env, tc.Query, resolver)
		require.NoError(t, err, "error parsing %s", tc.Query)

		assert.Equal(t, tc.Inspection, contactql.Inspect(query), "inspect mismatch for query %s", tc.Query)
	}

}
