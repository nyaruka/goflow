package contactql_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
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
		{query: `uuid = "C7D9BECE-6bbd-4b3b-8a86-eb0cf1ac9d05"`, result: true},
		{query: `uuid = "xyz"`, result: false},

		// ID condition
		{query: `id = 12345`, result: true},
		{query: `id = 76543`, result: false},

		// name condition
		{query: `name = "BOB smithwick"`, result: true},
		{query: `name = "Bob"`, result: false},
		{query: `name ~ "Bob"`, result: true},
		{query: `name ~ "Bobby"`, result: false},
		{query: `name ~ "Sm"`, result: true},
		{query: `name ~ "Smithwicke"`, result: true}, // only compare up to 8 chars
		{query: `name ~ "Smithx"`, result: false},

		// URN condition
		{query: `tel = +59313145145`, result: true},
		{query: `tel = +59313140000`, result: false},
		{query: `tel:+59313145145`, result: true},
		{query: `tel:+59313140000`, result: false},
		{query: `tel has 45145`, result: true},
		{query: `tel ~ 33333`, result: false},
		{query: `TWITTER IS bob_smith`, result: true},
		{query: `twitter:bob_smith`, result: true},
		{query: `twitter = jim_smith`, result: false},
		{query: `twitter:jim_smith`, result: false},
		{query: `twitter ~ smith`, result: true},
		{query: `whatsapp = 4533343`, result: false},

		// text field condition
		{query: `Gender = male`, result: true},
		{query: `Gender is MALE`, result: true},
		{query: `gender = "female"`, result: false},
		{query: `gender != "female"`, result: true},
		{query: `gender != "male"`, result: false},
		{query: `empty != "male"`, result: true}, // this is true because "" is not "male"
		{query: `gender != ""`, result: true},

		// number field condition
		{query: `age = 36`, result: true},
		{query: `age = 35`, result: false},
		{query: `age is 35`, result: false},
		{query: `age != 36`, result: false},
		{query: `age != 35`, result: true},
		{query: `age > 36`, result: false},
		{query: `age > 35`, result: true},
		{query: `age >= 36`, result: true},
		{query: `age < 36`, result: false},
		{query: `age < 37`, result: true},
		{query: `age <= 36`, result: true},

		// datetime field condition
		{query: `dob = 1981/05/28`, result: true},
		{query: `dob = 1981/05/29`, result: false},
		{query: `dob != 1981/05/28`, result: false},
		{query: `dob != 1981/05/29`, result: true},
		{query: `dob > 1981/05/28`, result: false},
		{query: `dob > 1981/05/27`, result: true},
		{query: `dob >= 1981/05/28`, result: true},
		{query: `dob >= 1981/05/29`, result: false},
		{query: `dob < 1981/05/28`, result: false},
		{query: `dob < 1981/05/29`, result: true},
		{query: `dob <= 1981/05/28`, result: true},
		{query: `dob <= 1981/05/27`, result: false},

		// location field condition
		{query: `state = kigali`, result: true},
		{query: `state = "kigali"`, result: true},
		{query: `state = "NYC"`, result: false},
		{query: `district = "GASABO"`, result: true},
		{query: `district = "Brooklyn"`, result: false},
		{query: `ward = ndera`, result: true},
		{query: `ward = solano`, result: false},
		{query: `ward != ndera`, result: false},
		{query: `ward != solano`, result: true},

		// existence
		{query: `age = ""`, result: false},
		{query: `age != ""`, result: true},
		{query: `xyz = ""`, result: true},
		{query: `xyz != ""`, result: false},
		{query: `age != "" AND xyz != ""`, result: false},
		{query: `age != "" OR xyz != ""`, result: true},

		// boolean combinations
		{query: `age = 36 AND gender = male`, result: true},
		{query: `(age = 36) AND (gender = male)`, result: true},
		{query: `age = 36 AND gender = female`, result: false},
		{query: `age = 36 OR gender = female`, result: true},
		{query: `age = 35 OR gender = female`, result: false},
		{query: `(age = 36 OR gender = female) AND age > 35`, result: true},
	}

	resolver := contactql.NewMockResolver(map[string]assets.Field{
		"age":      static.NewField(assets.FieldUUID("f1b5aea6-6586-41c7-9020-1a6326cc6565"), "age", "Age", assets.FieldTypeNumber),
		"dob":      static.NewField(assets.FieldUUID("3810a485-3fda-4011-a589-7320c0b8dbef"), "dob", "DOB", assets.FieldTypeDatetime),
		"gender":   static.NewField(assets.FieldUUID("d66a7823-eada-40e5-9a3a-57239d4690bf"), "gender", "Gender", assets.FieldTypeText),
		"state":    static.NewField(assets.FieldUUID("369be3e2-0186-4e5d-93c4-6264736588f8"), "state", "State", assets.FieldTypeState),
		"district": static.NewField(assets.FieldUUID("e52f34ad-a5a7-4855-9040-05a910a75f57"), "district", "District", assets.FieldTypeDistrict),
		"ward":     static.NewField(assets.FieldUUID("e9e738ce-617d-4c61-bfce-3d3b55cfe3dd"), "ward", "Ward", assets.FieldTypeWard),
		"empty":    static.NewField(assets.FieldUUID("023f733d-ce00-4a61-96e4-b411987028ea"), "empty", "Empty", assets.FieldTypeText),
		"xyz":      static.NewField(assets.FieldUUID("81e25783-a1d8-42b9-85e4-68c7ab2df39d"), "xyz", "XYZ", assets.FieldTypeText),
	}, map[string]assets.Group{})

	for _, test := range tests {
		parsed, err := contactql.ParseQuery(env, test.query, resolver)
		assert.NoError(t, err, "unexpected error parsing '%s'", test.query)

		actualResult := contactql.EvaluateQuery(env, parsed, testObj)
		assert.Equal(t, test.result, actualResult, "unexpected result for '%s'", test.query)
	}
}
