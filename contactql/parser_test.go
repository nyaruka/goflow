package contactql

import (
	"strconv"
	"testing"
	"time"

	"github.com/nyaruka/goflow/utils"
	"github.com/shopspring/decimal"
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

		// boolean operator precedence is AND before OR, even when AND is implicit
		{`will and felix or matt amber`, "OR(AND(*=will, *=felix), AND(*=matt, *=amber))"},

		// boolean combinations can themselves be combined
		{`(Age < 18 and Gender = "male") or (Age > 18 and Gender = "female")`, "OR(AND(age<18, gender=male), AND(age>18, gender=female))"},
	}

	for _, test := range tests {
		parsed, err := ParseQuery(test.text)
		if err != nil {
			t.Errorf("Error parsing query '%s'\n  Error: %s\n", test.text, err.Error())
			continue
		}
		if parsed.String() != test.parsed {
			t.Errorf("Error parsing query '%s'\n  Expected: %s\n  Got: %s\n", test.text, test.parsed, parsed.String())
		}
	}
}

type TestQueryable struct{}

func (t *TestQueryable) ResolveQueryKey(key string) interface{} {
	switch key {
	case "*":
		return []string{"Bob Smith", "+59313145145", "bob_smith"}
	case "name":
		return "Bob Smith"
	case "tel":
		return "+59313145145"
	case "twitter":
		return "bob_smith"
	case "gender":
		return "male"
	case "age":
		return decimal.NewFromFloat(36)
	case "dob":
		return time.Date(1981, 5, 28, 13, 30, 23, 0, time.UTC)
	}
	return nil
}

func TestEvaluateQuery(t *testing.T) {
	env := utils.NewDefaultEnvironment()
	testObj := &TestQueryable{}

	tests := []struct {
		text   string
		result bool
	}{
		// implicit key (i.e. name or URN) condition
		{`Bob Smith`, true},
		{`Smith`, true},
		{`Jim`, false},
		{`+59313145145`, true},
		{`+593131`, true},
		{`+5931317777777`, false},
		{`bob_smith`, true},

		// attribute condition
		{`name = "Bob Smith"`, true},
		{`name IS "Bob Smith"`, true},
		{`name = "Jim Smith"`, false},
		{`NAME HAS "Bob"`, true},
		{`Name has Bob`, true},
		{`name has Jim`, false},

		// URN condition
		{`tel = +59313145145`, true},
		{`tel has 45145`, true},
		{`tel ~ 33333`, false},
		{`TWITTER IS bob_smith`, true},
		{`twitter = jim_smith`, false},
		{`twitter ~ smith`, true},

		// text field condition
		{`Gender = male`, true},
		{`Gender is MALE`, true},
		{`gender = "female"`, false},

		// decimal field condition
		{`age = 36`, true},
		{`age is 35`, false},
		{`age > 36`, false},
		{`age > 35`, true},
		{`age >= 36`, true},
		{`age < 36`, false},
		{`age < 37`, true},
		{`age <= 36`, true},

		// date field condition
		{`dob = 1981/05/28`, true},
		{`dob > 1981/05/28`, false},
		{`dob > 1981/05/27`, true},
		{`dob >= 1981/05/28`, true},
		{`dob >= 1981/05/29`, false},
		{`dob < 1981/05/28`, false},
		{`dob < 1981/05/29`, true},
		{`dob <= 1981/05/28`, true},
		{`dob <= 1981/05/27`, false},

		// boolean combinations
		{`name = "Bob Smith" AND gender = male`, true},
		{`(name = "Bob Smith") AND (gender = male)`, true},
		{`name = "Bob Smith" AND gender = female`, false},
		{`name = "Bob Smith" OR gender = female`, true},
		{`(name = "Bob Smith" OR gender = female) AND age > 35`, true},
	}

	for _, test := range tests {
		parsed, err := ParseQuery(test.text)
		if err != nil {
			t.Errorf("Error parsing query '%s'\n  Error: %s\n", test.text, err.Error())
			continue
		}

		actualResult, err := EvaluateQuery(env, parsed, testObj)
		if err != nil {
			t.Errorf("Error evaluating query '%s'\n  Error: %s\n", test.text, err.Error())
			continue
		}
		if actualResult != test.result {
			t.Errorf("Error evaluating query '%s'\n  Expected: %s  Got: %s\n", test.text, strconv.FormatBool(test.result), strconv.FormatBool(actualResult))
		}
	}
}
