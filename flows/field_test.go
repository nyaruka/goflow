package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldValues(t *testing.T) {
	session, _, err := test.CreateTestSession("http://localhost", nil)
	require.NoError(t, err)

	env := session.Environment()
	fields := session.Assets().Fields()
	gender := fields.Get("gender")
	age := fields.Get("age")

	// can have no values for any fields
	fieldVals, err := flows.NewFieldValues(session.Assets(), map[string]*flows.Value{}, assets.PanicOnMissing)
	assert.NoError(t, err)

	assert.Equal(t, 0, fieldVals.Length())
	assert.Equal(t, "field values", fieldVals.Describe())
	assert.Nil(t, fieldVals.Resolve(env, "gender"))
	assert.Nil(t, fieldVals.Resolve(env, "age"))

	// can have a value but not in the right type for that field (age below)
	fieldVals, err = flows.NewFieldValues(session.Assets(), map[string]*flows.Value{
		"gender": flows.NewValue(types.NewXText("Male"), nil, nil, flows.LocationPath(""), flows.LocationPath(""), flows.LocationPath("")),
		"age":    flows.NewValue(types.NewXText("nan"), nil, nil, flows.LocationPath(""), flows.LocationPath(""), flows.LocationPath("")),
	}, assets.PanicOnMissing)
	assert.NoError(t, err)

	assert.Equal(t, types.NewXText("Male"), fieldVals.Get(gender).Text)
	assert.Equal(t, types.NewXText("nan"), fieldVals.Get(age).Text)

	genderVal := fieldVals["gender"]
	ageVal := fieldVals["age"]

	assert.Equal(t, 2, fieldVals.Length())
	assert.Equal(t, genderVal, fieldVals.Resolve(env, "gender"))
	assert.Equal(t, ageVal, fieldVals.Resolve(env, "age"))
	assert.Nil(t, fieldVals.Resolve(env, "join_date"))

	assert.Equal(t, types.NewXText("Male"), genderVal.Reduce(env))
	assert.Equal(t, types.NewXText("Male"), genderVal.Resolve(env, "text"))
	assert.Equal(t, "field value", genderVal.Describe())

	assert.Nil(t, ageVal.Reduce(env)) // doesn't have a value in the right type
	assert.Equal(t, types.NewXText("nan"), ageVal.Resolve(env, "text"))
}

func TestValues(t *testing.T) {
	num1 := types.RequireXNumberFromString("23")
	num2 := types.RequireXNumberFromString("23")
	num3 := types.RequireXNumberFromString("45")

	v1 := flows.NewValue(types.NewXText("Male"), nil, nil, flows.LocationPath(""), flows.LocationPath(""), flows.LocationPath(""))
	v2 := flows.NewValue(types.NewXText("Male"), nil, nil, flows.LocationPath(""), flows.LocationPath(""), flows.LocationPath(""))
	v3 := flows.NewValue(types.NewXText("23"), nil, &num1, flows.LocationPath(""), flows.LocationPath(""), flows.LocationPath(""))
	v4 := flows.NewValue(types.NewXText("23x"), nil, &num2, flows.LocationPath(""), flows.LocationPath(""), flows.LocationPath(""))
	v5 := flows.NewValue(types.NewXText("23x"), nil, &num3, flows.LocationPath(""), flows.LocationPath(""), flows.LocationPath(""))
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
