package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldValues(t *testing.T) {
	session, err := test.CreateTestSession("http://localhost", nil)
	require.NoError(t, err)

	env := session.Environment()
	fields := session.Assets().Fields()
	gender, _ := fields.Get("gender")
	age, _ := fields.Get("age")

	// can have no values for any fields
	fieldVals, err := flows.NewFieldValues(session.Assets(), map[string]*flows.Value{}, true)
	assert.NoError(t, err)

	assert.Equal(t, 0, fieldVals.Length())
	assert.Equal(t, "field values", fieldVals.Describe())
	assert.Nil(t, fieldVals.Resolve(env, "gender"))
	assert.Nil(t, fieldVals.Resolve(env, "age"))

	// can have a value but not in the right type for that field (age below)
	fieldVals, err = flows.NewFieldValues(session.Assets(), map[string]*flows.Value{
		"gender": flows.NewValue(types.NewXText("Male"), nil, nil, flows.LocationPath(""), flows.LocationPath(""), flows.LocationPath("")),
		"age":    flows.NewValue(types.NewXText("nan"), nil, nil, flows.LocationPath(""), flows.LocationPath(""), flows.LocationPath("")),
	}, true)
	assert.NoError(t, err)

	genderVal := fieldVals.GetValue(gender)
	ageVal := fieldVals.GetValue(age)
	assert.NotNil(t, genderVal)
	assert.NotNil(t, ageVal)

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
