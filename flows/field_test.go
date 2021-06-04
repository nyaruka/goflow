package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
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
	flows.NewFieldValues(session.Assets(), map[string]*flows.Value{}, assets.PanicOnMissing)

	// can have a value but not in the right type for that field (age below)
	fieldVals := flows.NewFieldValues(session.Assets(), map[string]*flows.Value{
		"gender": flows.NewValue(types.NewXText("Male"), nil, nil, envs.LocationPath(""), envs.LocationPath(""), envs.LocationPath("")),
		"age":    flows.NewValue(types.NewXText("nan"), nil, nil, envs.LocationPath(""), envs.LocationPath(""), envs.LocationPath("")),
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
	}), flows.Context(env, fieldVals))
}

func TestFieldValueParse(t *testing.T) {
	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	fields := session.Assets().Fields()
	gender := fields.Get("gender")
	age := fields.Get("age")
	state := fields.Get("state")

	xt := types.NewXText
	xn := func(s string) *types.XNumber { xn := types.RequireXNumberFromString(s); return &xn }
	nilLocPath := envs.LocationPath("")

	tcs := []struct {
		field    *flows.Field
		value    string
		expected *flows.Value
	}{
		{gender, "", nil},
		{gender, "M", flows.NewValue(xt("M"), nil, nil, nilLocPath, nilLocPath, nilLocPath)},
		{gender, " M ", flows.NewValue(xt(" M "), nil, nil, nilLocPath, nilLocPath, nilLocPath)},
		{gender, " 12 ", flows.NewValue(xt(" 12 "), nil, xn("12"), nilLocPath, nilLocPath, nilLocPath)},
		{age, "", nil},
		{age, "12", flows.NewValue(xt("12"), nil, xn("12"), nilLocPath, nilLocPath, nilLocPath)},
		{state, "", nil},
		{state, "kigali city", flows.NewValue(xt("kigali city"), nil, nil, envs.LocationPath("Rwanda > Kigali City"), nilLocPath, nilLocPath)},
		{state, "x", flows.NewValue(xt("x"), nil, nil, nilLocPath, nilLocPath, nilLocPath)},
	}

	for _, tc := range tcs {
		actual := session.Contact().Fields().Parse(session.Runs()[0].Environment(), fields, tc.field, tc.value)

		assert.Equal(t, tc.expected, actual, "parse mismatch for field %s and value '%s'", tc.field.Key(), tc.value)
	}
}

func TestValues(t *testing.T) {
	num1 := types.RequireXNumberFromString("23")
	num2 := types.RequireXNumberFromString("23")
	num3 := types.RequireXNumberFromString("45")

	v1 := flows.NewValue(types.NewXText("Male"), nil, nil, envs.LocationPath(""), envs.LocationPath(""), envs.LocationPath(""))
	v2 := flows.NewValue(types.NewXText("Male"), nil, nil, envs.LocationPath(""), envs.LocationPath(""), envs.LocationPath(""))
	v3 := flows.NewValue(types.NewXText("23"), nil, &num1, envs.LocationPath(""), envs.LocationPath(""), envs.LocationPath(""))
	v4 := flows.NewValue(types.NewXText("23x"), nil, &num2, envs.LocationPath(""), envs.LocationPath(""), envs.LocationPath(""))
	v5 := flows.NewValue(types.NewXText("23x"), nil, &num3, envs.LocationPath(""), envs.LocationPath(""), envs.LocationPath(""))
	v6 := (*flows.Value)(nil)

	assert.True(t, v1.Equals(v1))
	assert.True(t, v1.Equals(v2))
	assert.False(t, v2.Equals(v3))
	assert.False(t, v3.Equals(v4))
	assert.False(t, v4.Equals(v5))
	assert.False(t, v4.Equals(v6))
	assert.False(t, v6.Equals(v4))
	assert.True(t, v6.Equals(v6))
}
