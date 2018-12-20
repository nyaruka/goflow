package contactql

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		text   string
		parsed string
	}{
		{`will`, "*=will"},
		{`will felix`, "AND(*=will, *=felix)"},     // implicit AND
		{`will and felix`, "AND(*=will, *=felix)"}, // explicit AND
		{`will or felix or matt`, "OR(OR(*=will, *=felix), *=matt)"},
		{`Name=will`, "name=will"},
		{`Name ~ "felix"`, "name~felix"},
		{`name is ""`, `name=""`},  // is not set
		{`name != ""`, `name!=""`}, // is set
		{`name=will or Name ~ "felix"`, "OR(name=will, name~felix)"},
		{`Name is will or Name has felix`, "OR(name=will, name~felix)"}, // comparator aliases
		{`will or Name ~ "felix"`, "OR(*=will, name~felix)"},
		{`email ~ user@example.com`, "email~user@example.com"},

		// boolean operator precedence is AND before OR, even when AND is implicit
		{`will and felix or matt amber`, "OR(AND(*=will, *=felix), AND(*=matt, *=amber))"},

		// boolean combinations can themselves be combined
		{`(Age < 18 and Gender = "male") or (Age > 18 and Gender = "female")`, "OR(AND(age<18, gender=male), AND(age>18, gender=female))"},
	}

	for _, test := range tests {
		parsed, err := ParseQuery(test.text)
		assert.NoError(t, err)
		assert.Equal(t, test.parsed, parsed.String(), "error parsing query '%s'", test.text)
	}
}

type TestQueryable struct{}

func (t *TestQueryable) ResolveQueryKey(env utils.Environment, key string) []interface{} {
	switch key {
	case "tel":
		return []interface{}{"+59313145145"}
	case "twitter":
		return []interface{}{"bob_smith"}
	case "whatsapp":
		return []interface{}{}
	case "gender":
		return []interface{}{"male"}
	case "age":
		return []interface{}{decimal.NewFromFloat(36)}
	case "dob":
		return []interface{}{time.Date(1981, 5, 28, 13, 30, 23, 0, time.UTC)}
	case "state":
		return []interface{}{"Kigali"}
	case "district":
		return []interface{}{"Gasabo"}
	case "ward":
		return []interface{}{"Ndera"}
	}
	return nil
}

func TestEvaluateQuery(t *testing.T) {
	env := utils.NewDefaultEnvironment()
	testObj := &TestQueryable{}

	tests := []struct {
		query  string
		result bool
	}{
		// URN condition
		{`tel = +59313145145`, true},
		{`tel has 45145`, true},
		{`tel ~ 33333`, false},
		{`TWITTER IS bob_smith`, true},
		{`twitter = jim_smith`, false},
		{`twitter ~ smith`, true},
		{`whatsapp = 4533343`, false},

		// text field condition
		{`Gender = male`, true},
		{`Gender is MALE`, true},
		{`gender = "female"`, false},

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
		{`state = "NY"`, false},
		{`state ~ KIG`, true},
		{`state ~ NY`, false},
		{`district = "GASABO"`, true},
		{`district = "Brooklyn"`, false},
		{`district ~ SAB`, true},
		{`district ~ BRO`, false},
		{`ward = ndera`, true},
		{`ward = solano`, false},
		{`ward ~ era`, true},

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

	for _, test := range tests {
		parsed, err := ParseQuery(test.query)
		assert.NoError(t, err, "unexpected error parsing '%s'", test.query)

		actualResult, err := EvaluateQuery(env, parsed, testObj)
		assert.NoError(t, err, "unexpected error evaluating '%s'", test.query)
		assert.Equal(t, test.result, actualResult, "unexpected result for '%s'", test.query)
	}
}

func TestParsingErrors(t *testing.T) {
	_, err := ParseQuery("name = ")
	assert.EqualError(t, err, "mismatched input '<EOF>' expecting {TEXT, STRING}")
}

func TestEvaluationErrors(t *testing.T) {
	env := utils.NewDefaultEnvironment()
	testObj := &TestQueryable{}

	tests := []struct {
		query  string
		errMsg string
	}{
		{`Bob`, "dynamic group queries can't contain implicit conditions"},
		{`gender > Male`, "can't query text fields with >"},
		{`age ~ 32`, "can't query number fields with ~"},
		{`dob = 32`, "string '32' couldn't be parsed as a date"},
		{`dob = 32 AND name = Bob`, "string '32' couldn't be parsed as a date"},
		{`name = Bob OR dob = 32`, "string '32' couldn't be parsed as a date"},
		{`dob ~ 2018-12-31`, "can't query datetime fields with ~"},
	}

	for _, test := range tests {
		parsed, err := ParseQuery(test.query)
		assert.NoError(t, err, "unexpected error parsing '%s'", test.query)

		actualResult, err := EvaluateQuery(env, parsed, testObj)
		assert.EqualError(t, err, test.errMsg, "unexpected error evaluating '%s'", test.query)
		assert.False(t, actualResult, "unexpected non-false result for '%s'", test.query)
	}
}
