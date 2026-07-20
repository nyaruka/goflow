package core_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldValues(t *testing.T) {
	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	env := session.Environment()
	fields := session.Assets().Fields()
	gender := fields.Get("gender")
	age := fields.Get("age")

	// can have no values for any fields
	core.NewFieldValues(session.Assets().Fields(), map[string]*core.Value{}, assets.PanicOnMissing)

	// can have a value but not in the right type for that field (age below)
	fieldVals := core.NewFieldValues(session.Assets().Fields(), map[string]*core.Value{
		"gender": core.NewValue(types.NewXText("Male"), nil, nil, envs.LocationPath(""), envs.LocationPath(""), envs.LocationPath("")),
		"age":    core.NewValue(types.NewXText("nan"), nil, nil, envs.LocationPath(""), envs.LocationPath(""), envs.LocationPath("")),
	}, assets.PanicOnMissing)

	assert.Equal(t, types.NewXText("Male"), fieldVals.Get(gender).Text)
	assert.Equal(t, types.NewXText("nan"), fieldVals.Get(age).Text)

	genderVal := fieldVals["gender"]
	ageVal := fieldVals["age"]

	test.AssertXEqual(t, types.NewXText("Male"), genderVal.ToXValue(env))
	assert.Nil(t, ageVal.ToXValue(env)) // doesn't have a value in the right type

	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__":      types.NewXText("Gender: Male"),
		"activation_token": nil,
		"age":              nil,
		"gender":           types.NewXText("Male"),
		"join_date":        nil,
		"state":            nil,
		"not_set":          nil,
		"language":         nil,
	}), core.Context(env, fieldVals))
}

func TestFieldValueParse(t *testing.T) {
	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	fields := session.Assets().Fields()
	gender := fields.Get("gender")
	age := fields.Get("age")
	state := fields.Get("state")

	xt := types.NewXText
	xn := func(s string) *types.XNumber { xn := types.RequireXNumberFromString(s); return xn }
	nilLocPath := envs.LocationPath("")

	tcs := []struct {
		field    *core.Field
		value    string
		expected *core.Value
	}{
		{gender, "", nil},
		{gender, "M", core.NewValue(xt("M"), nil, nil, nilLocPath, nilLocPath, nilLocPath)},
		{gender, " M ", core.NewValue(xt(" M "), nil, nil, nilLocPath, nilLocPath, nilLocPath)},
		{gender, " 12 ", core.NewValue(xt(" 12 "), nil, xn("12"), nilLocPath, nilLocPath, nilLocPath)},
		{age, "", nil},
		{age, "12", core.NewValue(xt("12"), nil, xn("12"), nilLocPath, nilLocPath, nilLocPath)},
		{state, "", nil},
		{state, "kigali city", core.NewValue(xt("kigali city"), nil, nil, envs.LocationPath("Rwanda > Kigali City"), nilLocPath, nilLocPath)},
		{state, "x", core.NewValue(xt("x"), nil, nil, nilLocPath, nilLocPath, nilLocPath)},
	}

	for _, tc := range tcs {
		actual := session.Contact().Fields().Parse(session.MergedEnvironment(), fields, tc.field, tc.value)

		assert.Equal(t, tc.expected, actual, "parse mismatch for field %s and value '%s'", tc.field.Key(), tc.value)
	}
}

func TestValues(t *testing.T) {
	num1 := types.RequireXNumberFromString("23")
	num2 := types.RequireXNumberFromString("23")
	num3 := types.RequireXNumberFromString("45")

	v1 := core.NewValue(types.NewXText("Male"), nil, nil, envs.LocationPath(""), envs.LocationPath(""), envs.LocationPath(""))
	v2 := core.NewValue(types.NewXText("Male"), nil, nil, envs.LocationPath(""), envs.LocationPath(""), envs.LocationPath(""))
	v3 := core.NewValue(types.NewXText("23"), nil, num1, envs.LocationPath(""), envs.LocationPath(""), envs.LocationPath(""))
	v4 := core.NewValue(types.NewXText("23x"), nil, num2, envs.LocationPath(""), envs.LocationPath(""), envs.LocationPath(""))
	v5 := core.NewValue(types.NewXText("23x"), nil, num3, envs.LocationPath(""), envs.LocationPath(""), envs.LocationPath(""))
	v6 := (*core.Value)(nil)

	assert.True(t, v1.Equals(v1))
	assert.True(t, v1.Equals(v2))
	assert.False(t, v2.Equals(v3))
	assert.False(t, v3.Equals(v4))
	assert.False(t, v4.Equals(v5))
	assert.False(t, v4.Equals(v6))
	assert.False(t, v6.Equals(v4))
	assert.True(t, v6.Equals(v6))
}

func TestFieldAssets(t *testing.T) {
	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	// field assets are used as a resolver for query parsing
	fields := session.Assets().Fields()
	age := fields.ResolveField("age")
	assert.Equal(t, assets.FieldUUID(`f1b5aea6-6586-41c7-9020-1a6326cc6565`), age.UUID())
	assert.Equal(t, "age", age.Key())
	assert.Equal(t, "Age", age.Name())
	assert.Equal(t, assets.FieldTypeNumber, age.Type())

	// but groups don't support query conditions on groups or flows so those aren't included
	assert.Nil(t, fields.ResolveGroup(`b7cf0d83-f1c9-411c-96fd-c511a4cfa86d`))
	assert.Nil(t, fields.ResolveFlow(`50c3706e-fedb-42c0-8eab-dda3335714b7`))
}

// a location hierarchy can be constructed programmatically with an incomplete parent chain, which parsing a field
// value against must tolerate rather than panic on
type orphanHierarchy struct{ *envs.LocationHierarchy }

func TestFieldValueParseWithIncompleteLocationHierarchy(t *testing.T) {
	env := envs.NewBuilder().Build()

	tcs := []struct {
		level    envs.LocationLevel
		state    envs.LocationPath
		district envs.LocationPath
		ward     envs.LocationPath
	}{
		{core.LocationLevelState, "Rwanda > Kigali", "", ""},
		{core.LocationLevelDistrict, "", "Rwanda > Kigali", ""},
		{core.LocationLevelWard, "", "", "Rwanda > Kigali"},
	}

	for _, tc := range tcs {
		// a hierarchy rooted at this level, so its locations have no parents at all
		root := envs.NewLocation(tc.level, "Rwanda > Kigali")
		locations := core.NewLocationAssets([]assets.LocationHierarchy{
			&orphanHierarchy{envs.NewLocationHierarchy(env, root, 4)},
		})

		env := envs.NewBuilder().WithLocationResolver(locations).Build()
		fields := core.NewFieldAssets([]assets.Field{
			static.NewField("f4a0d0c1-4d0e-4e9c-9b2b-1a3e6f1b2c3d", "place", "Place", assets.FieldTypeWard),
		})
		field := fields.All()[0]

		actual := core.FieldValues{}.Parse(env, fields, field, "Rwanda > Kigali")

		assert.Equal(t, core.NewValue(types.NewXText("Rwanda > Kigali"), nil, nil, tc.state, tc.district, tc.ward), actual,
			"parse mismatch at location level %d", tc.level)
	}
}
