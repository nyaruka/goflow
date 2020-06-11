package contactql_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static/types"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type TestQueryable map[string][]interface{}

func (t TestQueryable) QueryProperty(env envs.Environment, key string, propType contactql.PropertyType) []interface{} {
	return t[key]
}

func TestEvaluateQuery(t *testing.T) {
	env := envs.NewBuilder().Build()
	var testObj = TestQueryable{
		"uuid":     []interface{}{"c7d9bece-6bbd-4b3b-8a86-eb0cf1ac9d05"},
		"id":       []interface{}{"12345"},
		"name":     []interface{}{"Bob Smithwick"},
		"tel":      []interface{}{"+59313145145"},
		"twitter":  []interface{}{"bob_smith"},
		"whatsapp": []interface{}{},
		"gender":   []interface{}{"male"},
		"age":      []interface{}{decimal.NewFromFloat(36)},
		"dob":      []interface{}{time.Date(1981, 5, 28, 13, 30, 23, 0, time.UTC)},
		"state":    []interface{}{"Kigali"},
		"district": []interface{}{"Gasabo"},
		"ward":     []interface{}{"Ndera"},
		"empty":    []interface{}{""},
		"nope":     []interface{}{envs.NewBuilder().Build()},
	}

	tests := []struct {
		query  string
		result bool
	}{
		// UUID condition
		{`uuid = "C7D9BECE-6bbd-4b3b-8a86-eb0cf1ac9d05"`, true},
		{`uuid = "xyz"`, false},

		// ID condition
		{`id = 12345`, true},
		{`id = 76543`, false},

		// name condition
		{`name = "BOB smithwick"`, true},
		{`name = "Bob"`, false},
		{`name ~ "Bob"`, true},
		{`name ~ "Bobby"`, false},
		{`name ~ "Sm"`, true},
		{`name ~ "Smithwicke"`, true}, // only compare up to 8 chars
		{`name ~ "Smithx"`, false},

		// URN condition
		{`tel = +59313145145`, true},
		{`tel = +59313140000`, false},
		{`tel:+59313145145`, true},
		{`tel:+59313140000`, false},
		{`tel has 45145`, true},
		{`tel ~ 33333`, false},
		{`TWITTER IS bob_smith`, true},
		{`twitter:bob_smith`, true},
		{`twitter = jim_smith`, false},
		{`twitter:jim_smith`, false},
		{`twitter ~ smith`, true},
		{`whatsapp = 4533343`, false},

		// text field condition
		{`Gender = male`, true},
		{`Gender is MALE`, true},
		{`gender = "female"`, false},
		{`gender != "female"`, true},
		{`gender != "male"`, false},
		{`empty != "male"`, true}, // this is true because "" is not "male"
		{`gender != ""`, true},

		// number field condition
		{`age = 36`, true},
		{`age is 35`, false},
		{`age > 36`, false},
		{`age > 35`, true},
		{`age >= 36`, true},
		{`age < 36`, false},
		{`age < 37`, true},
		{`age <= 36`, true},

		// datetime field condition
		{`dob = 1981/05/28`, true},
		{`dob > 1981/05/28`, false},
		{`dob > 1981/05/27`, true},
		{`dob >= 1981/05/28`, true},
		{`dob >= 1981/05/29`, false},
		{`dob < 1981/05/28`, false},
		{`dob < 1981/05/29`, true},
		{`dob <= 1981/05/28`, true},
		{`dob <= 1981/05/27`, false},

		// location field condition
		{`state = kigali`, true},
		{`state = "kigali"`, true},
		{`state = "NYC"`, false},
		{`district = "GASABO"`, true},
		{`district = "Brooklyn"`, false},
		{`ward = ndera`, true},
		{`ward = solano`, false},
		{`ward != ndera`, false},
		{`ward != solano`, true},

		// existence
		{`age = ""`, false},
		{`age != ""`, true},
		{`xyz = ""`, true},
		{`xyz != ""`, false},
		{`age != "" AND xyz != ""`, false},
		{`age != "" OR xyz != ""`, true},

		// boolean combinations
		{`age = 36 AND gender = male`, true},
		{`(age = 36) AND (gender = male)`, true},
		{`age = 36 AND gender = female`, false},
		{`age = 36 OR gender = female`, true},
		{`age = 35 OR gender = female`, false},
		{`(age = 36 OR gender = female) AND age > 35`, true},
	}

	resolver := contactql.NewMockResolver(map[string]assets.Field{
		"age":      types.NewField(assets.FieldUUID("f1b5aea6-6586-41c7-9020-1a6326cc6565"), "age", "Age", assets.FieldTypeNumber),
		"dob":      types.NewField(assets.FieldUUID("3810a485-3fda-4011-a589-7320c0b8dbef"), "dob", "DOB", assets.FieldTypeDatetime),
		"gender":   types.NewField(assets.FieldUUID("d66a7823-eada-40e5-9a3a-57239d4690bf"), "gender", "Gender", assets.FieldTypeText),
		"state":    types.NewField(assets.FieldUUID("369be3e2-0186-4e5d-93c4-6264736588f8"), "state", "State", assets.FieldTypeState),
		"district": types.NewField(assets.FieldUUID("e52f34ad-a5a7-4855-9040-05a910a75f57"), "district", "District", assets.FieldTypeDistrict),
		"ward":     types.NewField(assets.FieldUUID("e9e738ce-617d-4c61-bfce-3d3b55cfe3dd"), "ward", "Ward", assets.FieldTypeWard),
		"empty":    types.NewField(assets.FieldUUID("023f733d-ce00-4a61-96e4-b411987028ea"), "empty", "Empty", assets.FieldTypeText),
		"xyz":      types.NewField(assets.FieldUUID("81e25783-a1d8-42b9-85e4-68c7ab2df39d"), "xyz", "XYZ", assets.FieldTypeText),
	}, map[string]assets.Group{})

	for _, test := range tests {
		parsed, err := contactql.ParseQuery(test.query, envs.RedactionPolicyNone, "", resolver)
		assert.NoError(t, err, "unexpected error parsing '%s'", test.query)

		actualResult, err := contactql.EvaluateQuery(env, parsed, testObj)
		assert.NoError(t, err, "unexpected error evaluating '%s'", test.query)
		assert.Equal(t, test.result, actualResult, "unexpected result for '%s'", test.query)
	}
}

func TestEvaluationErrors(t *testing.T) {
	env := envs.NewBuilder().Build()
	var testObj = TestQueryable{
		"name":   []interface{}{"Bob Smithwick"},
		"gender": []interface{}{"male"},
		"age":    []interface{}{decimal.NewFromFloat(36)},
		"dob":    []interface{}{time.Date(1981, 5, 28, 13, 30, 23, 0, time.UTC)},
	}

	tests := []struct {
		query  string
		errMsg string
	}{
		{`age = 3X`, "can't convert '3X' to a number"},
		{`dob = 32`, "can't convert '32' to a date"},
		{`dob = 32 AND name = Bob`, "can't convert '32' to a date"},
		{`name = Bob OR dob = 32`, "can't convert '32' to a date"},
	}

	resolver := contactql.NewMockResolver(map[string]assets.Field{
		"age":    types.NewField(assets.FieldUUID("f1b5aea6-6586-41c7-9020-1a6326cc6565"), "age", "Age", assets.FieldTypeNumber),
		"dob":    types.NewField(assets.FieldUUID("3810a485-3fda-4011-a589-7320c0b8dbef"), "dob", "DOB", assets.FieldTypeDatetime),
		"gender": types.NewField(assets.FieldUUID("d66a7823-eada-40e5-9a3a-57239d4690bf"), "gender", "Gender", assets.FieldTypeText),
	}, map[string]assets.Group{})

	for _, test := range tests {
		parsed, err := contactql.ParseQuery(test.query, envs.RedactionPolicyNone, "", resolver)
		assert.NoError(t, err, "unexpected error parsing '%s'", test.query)

		actualResult, err := contactql.EvaluateQuery(env, parsed, testObj)
		assert.EqualError(t, err, test.errMsg, "unexpected error evaluating '%s'", test.query)
		assert.False(t, actualResult, "unexpected non-false result for '%s'", test.query)
	}
}
