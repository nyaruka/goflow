package contactql_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInspect(t *testing.T) {
	env := envs.NewBuilder().Build()
	resolver := contactql.NewMockResolver(
		[]assets.Field{
			static.NewField(assets.FieldUUID("f1b5aea6-6586-41c7-9020-1a6326cc6565"), "age", "Age", assets.FieldTypeNumber),
			static.NewField(assets.FieldUUID("3810a485-3fda-4011-a589-7320c0b8dbef"), "dob", "DOB", assets.FieldTypeDatetime),
			static.NewField(assets.FieldUUID("d66a7823-eada-40e5-9a3a-57239d4690bf"), "gender", "Gender", assets.FieldTypeText),
		},
		[]assets.Flow{},
		[]assets.Group{
			static.NewGroup(assets.GroupUUID("4eeca453-f474-4767-bdd0-434b180223db"), "U-Reporters", ""),
		},
	)

	tests := []struct {
		query      string
		resolver   contactql.Resolver
		inspection *contactql.Inspection
	}{
		{
			query:    "bob",
			resolver: resolver,
			inspection: &contactql.Inspection{
				Attributes:   []string{"name"},
				Schemes:      []string{},
				Fields:       []*assets.FieldReference{},
				Groups:       []*assets.GroupReference{},
				AllowAsGroup: true,
			},
		},
		{
			query: "AGE > 18 AND name != \"\" OR twitter = bobby OR tel ~1234 AND tickets > 0",
			inspection: &contactql.Inspection{
				Attributes: []string{"name", "tickets"},
				Schemes:    []string{"tel", "twitter"},
				Fields: []*assets.FieldReference{
					assets.NewFieldReference("age", ""),
				},
				Groups:       []*assets.GroupReference{},
				AllowAsGroup: true,
			},
		},
		{
			query:    "AGE > 18 AND name != \"\" OR twitter = bobby OR tel ~1234 AND tickets > 0",
			resolver: resolver,
			inspection: &contactql.Inspection{
				Attributes: []string{"name", "tickets"},
				Schemes:    []string{"tel", "twitter"},
				Fields: []*assets.FieldReference{
					assets.NewFieldReference("age", "Age"),
				},
				Groups:       []*assets.GroupReference{},
				AllowAsGroup: true,
			},
		},
		{
			query:    "id = 123",
			resolver: resolver,
			inspection: &contactql.Inspection{
				Attributes:   []string{"id"},
				Schemes:      []string{},
				Fields:       []*assets.FieldReference{},
				Groups:       []*assets.GroupReference{},
				AllowAsGroup: false,
			},
		},
		{
			query: "group = U-reporters", // group condition parsed without resolver
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
		{
			query:    "group = U-reporters", // group condition parsed with resolver
			resolver: resolver,
			inspection: &contactql.Inspection{
				Attributes: []string{"group"},
				Schemes:    []string{},
				Fields:     []*assets.FieldReference{},
				Groups: []*assets.GroupReference{
					assets.NewGroupReference("4eeca453-f474-4767-bdd0-434b180223db", "U-Reporters"),
				},
				AllowAsGroup: false,
			},
		},
	}

	for _, tc := range tests {
		query, err := contactql.ParseQuery(env, tc.query, tc.resolver)
		require.NoError(t, err, "error parsing %s", tc.query)

		assert.Equal(t, tc.inspection, contactql.Inspect(query), "inspect mismatch for query %s", tc.query)
	}

}
