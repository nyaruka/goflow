package expressions_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/legacy/expressions"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testServerPort = 49997

type testTemplate struct {
	old string
	new string

	extraAs expressions.ExtraVarsMapping
}

func TestMigrateTemplate(t *testing.T) {
	var tests = []testTemplate{

		// contact variables
		{old: `@contact`, new: `@contact`},
		{old: `@CONTACT`, new: `@contact`},
		{old: `@contact.uuid`, new: `@contact.uuid`},
		{old: `@contact.id`, new: `@contact.id`},
		{old: `@contact.name`, new: `@contact.name`},
		{old: `@contact.NAME`, new: `@contact.name`},
		{old: `@contact.first_name`, new: `@contact.first_name`},
		{old: `@contact.gender`, new: `@contact.fields.gender`},
		{old: `@contact.groups`, new: `@(join(contact.groups, ","))`},

		// contact URN variables
		{old: `@contact.tel`, new: `@(format_urn(contact.urns.tel))`},
		{old: `@contact.tel.display`, new: `@(format_urn(contact.urns.tel))`},
		{old: `@contact.tel.scheme`, new: `@contact.urns.tel.0.scheme`},
		{old: `@contact.tel.path`, new: `@contact.urns.tel.0.path`},
		{old: `@contact.tel.urn`, new: `@contact.urns.tel.0`},
		{old: `@contact.tel_e164`, new: `@contact.urns.tel.0.path`},
		{old: `@contact.twitterid`, new: `@(format_urn(contact.urns.twitterid))`},
		{old: `@contact.mailto`, new: `@(format_urn(contact.urns.mailto))`},

		// run variables
		{old: `@flow.favorite_color`, new: `@run.results.favorite_color`},
		{old: `@flow.favorite_color.category`, new: `@run.results.favorite_color.category_localized`},
		{old: `@flow.favorite_color.text`, new: `@run.results.favorite_color.input`},
		{old: `@flow.favorite_color.time`, new: `@run.results.favorite_color.created_on`},
		{old: `@flow.favorite_color.value`, new: `@run.results.favorite_color.value`},
		{old: `@flow.contact`, new: `@contact`},
		{old: `@flow.contact.name`, new: `@contact.name`},

		{old: `@child.age`, new: `@child.results.age`},
		{old: `@child.contact`, new: `@child.contact`},
		{old: `@child.contact.age`, new: `@child.contact.fields.age`},

		{old: `@parent.role`, new: `@parent.results.role`},
		{old: `@parent.role.category`, new: `@parent.results.role.category_localized`},
		{old: `@parent.contact`, new: `@parent.contact`},
		{old: `@parent.contact.name`, new: `@parent.contact.name`},
		{old: `@parent.contact.gender`, new: `@parent.contact.fields.gender`},

		// input
		{old: `@step`, new: `@run.input`},
		{old: `@step.value`, new: `@run.input`},
		{old: `@step.text`, new: `@run.input.text`},
		{old: `@step.attachments`, new: `@run.input.attachments`},
		{old: `@step.time`, new: `@run.input.created_on`},
		{old: `@step.contact`, new: `@contact`},
		{old: `@step.contact.name`, new: `@contact.name`},
		{old: `@step.contact.age`, new: `@contact.fields.age`},

		// dates
		{old: `@date`, new: `@(now())`},
		{old: `@date.now`, new: `@(now())`},
		{old: `@date.today`, new: `@(today())`},
		{old: `@date.tomorrow`, new: `@(datetime_add(today(), 1, "D"))`},
		{old: `@date.yesterday`, new: `@(datetime_add(today(), -1, "D"))`},

		// extra
		{old: `@extra.results.0.state`, new: `@run.webhook.json.results.0.state`, extraAs: expressions.ExtraAsWebhookJSON},
		{old: `@extra.address.state`, new: `@trigger.params.address.state`, extraAs: expressions.ExtraAsTriggerParams},
		{old: `@extra.address.state`, new: `@(if(is_error(run.webhook.json.address.state), trigger.params.address.state, run.webhook.json.address.state))`, extraAs: expressions.ExtraAsFunction},
		{old: `@extra.flow.role`, new: `@parent.results.role`},

		// variables in parens
		{old: `@(contact.tel)`, new: `@(format_urn(contact.urns.tel))`},
		{old: `@(contact.gender)`, new: `@(contact.fields.gender)`},
		{old: `@(flow.favorite_color)`, new: `@(run.results.favorite_color)`},

		// booleans
		{old: `@(TRUE)`, new: `@(true)`},
		{old: `@(False)`, new: `@(false)`},
		{old: `@(TRUE())`, new: `@(true)`},
		{old: `@(FALSE())`, new: `@(false)`},

		// arithmetic
		{old: `@(1 + 2)`, new: `@(1 + 2)`},
		{old: `@(1 - 2)`, new: `@(1 - 2)`},
		{old: `@(-2)`, new: `@(-2)`},
		{old: `@(2 ^ 4)`, new: `@(2 ^ 4)`},
		{old: `@(2 * 4)`, new: `@(2 * 4)`},
		{old: `@(2 / 4)`, new: `@(2 / 4)`},

		// comparisons
		{old: `@(1 = 4)`, new: `@(1 = 4)`},
		{old: `@(1 <> 4)`, new: `@(1 != 4)`},
		{old: `@(1 < 4)`, new: `@(1 < 4)`},
		{old: `@(1 <= 4)`, new: `@(1 <= 4)`},
		{old: `@(1 > 4)`, new: `@(1 > 4)`},
		{old: `@(1 >= 4)`, new: `@(1 >= 4)`},

		// string concatenation
		{old: `@("you" & " are " & contact.gender)`, new: `@("you" & " are " & contact.fields.gender)`},

		// functions
		{old: `@(REMOVE_FIRST_WORD(flow.favorite_color))`, new: `@(remove_first_word(run.results.favorite_color))`},
		{old: `@(WORD_SLICE(flow.favorite_color, 2))`, new: `@(word_slice(run.results.favorite_color, 2 - 1))`},
		{old: `@(WORD_SLICE(flow.favorite_color, 2, 4))`, new: `@(word_slice(run.results.favorite_color, 2 - 1, 4 - 1))`},
		{old: `@(FIELD(flow.favorite_color, 2, ","))`, new: `@(field(run.results.favorite_color, 2 - 1, ","))`},
		{old: `@(FIELD(flow.favorite_color, child.age, ","))`, new: `@(field(run.results.favorite_color, child.results.age - 1, ","))`},
		{old: `@(FIRST_WORD(WORD_SLICE("bee cat dog elf", 2, 4)))`, new: `@(word(word_slice("bee cat dog elf", 2 - 1, 4 - 1), 0))`},

		{old: `@(FIRST_WORD(flow.favorite_color))`, new: `@(word(run.results.favorite_color, 0))`},
		{old: `@(WORD(flow.favorite_color, child.age - 22))`, new: `@(word(run.results.favorite_color, legacy_add(child.results.age, -22) - 1))`},
		{old: `@(WORD(flow.favorite_color, 1))`, new: `@(word(run.results.favorite_color, 1 - 1))`},
		{old: `@(ABS(-5))`, new: `@(abs(-5))`},
		{old: `@(AVERAGE(1, 2, 3, 4, 5))`, new: `@(mean(1, 2, 3, 4, 5))`},
		{old: `@(AND(contact.age > 30, flow.amount < 5))`, new: `@(and(contact.fields.age > 30, run.results.amount < 5))`},
		{old: `@(DATEVALUE("2012-02-03"))`, new: `@(datetime("2012-02-03"))`},
		{old: `@(EDATE("2012-02-03", 1))`, new: `@(datetime_add("2012-02-03", 1, "M"))`},
		{old: `@(DATEDIF(contact.join_date, date.now, "M"))`, new: `@(datetime_diff(contact.fields.join_date, now(), "M"))`},
		{old: `@(DAYS("2016-02-28", "2015-02-28"))`, new: `@(datetime_diff("2016-02-28", "2015-02-28", "D"))`},
		{old: `@(DAY(contact.join_date))`, new: `@(format_datetime(contact.fields.join_date, "D"))`},
		{old: `@(HOUR(NOW()))`, new: `@(format_datetime(now(), "h"))`},
		{old: `@(MINUTE(NOW()))`, new: `@(format_datetime(now(), "m"))`},
		{old: `@(MONTH(NOW()))`, new: `@(format_datetime(now(), "M"))`},
		{old: `@(NOW())`, new: `@(now())`},
		{old: `@(SECOND(NOW()))`, new: `@(format_datetime(now(), "s"))`},

		// date addition should get converted to datetime_add
		{old: `@(date.now + 5)`, new: `@(datetime_add(now(), 5, "D"))`},
		{old: `@(now() + 5)`, new: `@(datetime_add(now(), 5, "D"))`},
		{old: `@(date + 5)`, new: `@(datetime_add(now(), 5, "D"))`},
		{old: `@(date.now + 5 + contact.age)`, new: `@(legacy_add(datetime_add(now(), 5, "D"), contact.fields.age))`},

		// legacy_add permutations
		{old: `@(contact.age + 5)`, new: `@(legacy_add(contact.fields.age, 5))`},
		{old: `@(contact.join_date + 5 + contact.age)`, new: `@(legacy_add(legacy_add(contact.fields.join_date, 5), contact.fields.age))`},
		{old: `@(contact.age + 100 - 5)`, new: `@(legacy_add(legacy_add(contact.fields.age, 100), -5))`},
		{old: `@(date.yesterday - 3 + 10)`, new: `@(legacy_add(legacy_add(datetime_add(today(), -1, "D"), -3), 10))`},

		{old: `@(3 + date.now)`, new: `@(datetime_add(now(), 3, "D"))`},
		{old: `@(date.tomorrow - 3)`, new: `@(legacy_add(datetime_add(today(), 1, "D"), -3))`},
		{old: `@(date.now + TIME(2, 30, 0))`, new: `@(datetime_add(now(), 9000, "s"))`},
		{old: `@(TIME(0, 1, 5) + contact.join_date)`, new: `@(datetime_add(contact.fields.join_date, 65, "s"))`},
		{old: `@(contact.join_date - TIME(0,0,12))`, new: `@(datetime_add(contact.fields.join_date, -12, "s"))`},

		// TODO: beware different org format, need a like function
		{old: `@(DATE(2012, 12, 25))`, new: `@(datetime("2012-12-25"))`},
		{old: `@(5 * contact.age)`, new: `@(5 * contact.fields.age)`},

		{old: `@((5 + contact.age) / 2)`, new: `@((legacy_add(5, contact.fields.age)) / 2)`},
		{old: `@(WEEKDAY(TODAY()))`, new: `@(weekday(today()))`},
		{old: `@(YEAR(date.now))`, new: `@(format_datetime(now(), "YYYY"))`},

		// booleans and conditionals
		{old: `@(AND(contact.gender = "F", contact.age >= 18))`, new: `@(and(contact.fields.gender = "F", contact.fields.age >= 18))`},
		{old: `@(IF(contact.gender = "M", "Sir", "Madam"))`, new: `@(if(contact.fields.gender = "M", "Sir", "Madam"))`},
		{old: `@(OR(contact.gender = "M", contact.gender = "F", contact.gender = "NB"))`, new: `@(or(contact.fields.gender = "M", contact.fields.gender = "F", contact.fields.gender = "NB"))`},

		// math functions
		{old: `@(ABS(-1))`, new: `@(abs(-1))`},
		{old: `@(MAX(child.age, 10))`, new: `@(max(child.results.age, 10))`},
		{old: `@(MIN(child.age, 10))`, new: `@(min(child.results.age, 10))`},
		{old: `@(POWER(2, 3))`, new: `@(2 ^ 3)`},
		{old: `@(RAND())`, new: `@(rand())`},
		{old: `@(RANDBETWEEN(1, 10))`, new: `@(rand_between(1, 10))`},
		{old: `@(ROUND(9.4378))`, new: `@(round(9.4378))`},
		{old: `@(ROUND(9.4378, 3))`, new: `@(round(9.4378, 3))`},
		{old: `@(ROUNDUP(9.4378))`, new: `@(round_up(9.4378))`},
		{old: `@(ROUNDUP(9.4378, 3))`, new: `@(round_up(9.4378, 3))`},
		{old: `@(ROUNDDOWN(9.4378))`, new: `@(round_down(9.4378))`},
		{old: `@(ROUNDDOWN(9.4378, 3))`, new: `@(round_down(9.4378, 3))`},
		{old: `@(SUM(contact.age, child.age))`, new: `@(contact.fields.age + child.results.age)`},
		{old: `@(CHAR(10))`, new: `@(char(10))`},
		{old: `@(CLEAN(contact.gender))`, new: `@(clean(contact.fields.gender))`},
		{old: `@(CODE("A"))`, new: `@(code("A"))`},
		{old: `@(CONCATENATE(contact.first_name, " ", contact.language))`, new: `@(contact.first_name & " " & contact.language)`},

		{old: `@(FIXED(contact.age))`, new: `@(format_number(contact.fields.age))`},
		{old: `@(FIXED(contact.age, 2))`, new: `@(format_number(contact.fields.age, 2))`},
		{old: `@(FIXED(contact.age, 2, false))`, new: `@(format_number(contact.fields.age, 2, false))`},
		{old: `@(INT(contact.age))`, new: `@(round_down(contact.fields.age))`},
		{old: `@(LEFT(contact.name, 4))`, new: `@(left(contact.name, 4))`},
		{old: `@(RIGHT(contact.name, 4))`, new: `@(right(contact.name, 4))`},
		{old: `@(LEN(contact.first_name))`, new: `@(length(contact.first_name))`},
		{old: `@(LOWER(contact.first_name))`, new: `@(lower(contact.first_name))`},
		{old: `@(mod(103, 4))`, new: `@(mod(103, 4))`},

		{old: `@(PROPER(contact))`, new: `@(title(contact))`},
		{old: `@(REPT("*", 10))`, new: `@(repeat("*", 10))`},
		{old: `@((DATEDIF(DATEVALUE("1970-01-01"), date.now, "D") * 24 * 60 * 60) + ((((HOUR(date.now)+7) * 60) + MINUTE(date.now)) * 60))`, new: `@(legacy_add((datetime_diff(datetime("1970-01-01"), now(), "D") * 24 * 60 * 60), ((legacy_add(((legacy_add(format_datetime(now(), "h"), 7)) * 60), format_datetime(now(), "m"))) * 60)))`},

		// non-expressions
		{old: `bob@nyaruka.com`, new: `bob@nyaruka.com`},
		{old: `@twitter_handle`, new: `@twitter_handle`},

		{old: `@`, new: `@`},
		{old: `Hi @@@flow.favorite_color @@flow.favorite_color @flow.favorite_color @nyaruka @ @`, new: `Hi @@@run.results.favorite_color @@flow.favorite_color @run.results.favorite_color @nyaruka @ @`},
	}

	for i := range tests {
		tests = append(tests, testTemplate{old: "Embedded " + tests[i].old + " text", new: "Embedded " + tests[i].new + " text", extraAs: tests[i].extraAs})
		tests = append(tests, testTemplate{old: "Replace " + tests[i].old + " two " + tests[i].old + " times", new: "Replace " + tests[i].new + " two " + tests[i].new + " times", extraAs: tests[i].extraAs})
	}

	server, err := test.NewTestHTTPServer(testServerPort)
	require.NoError(t, err)

	defer server.Close()

	session, err := test.CreateTestSession(testServerPort, nil)
	require.NoError(t, err)

	for _, test := range tests {

		for i := 0; i < 1; i++ {
			migratedTemplate, err := expressions.MigrateTemplate(test.old, test.extraAs)

			defer func() {
				if r := recover(); r != nil {
					t.Errorf("panic migrating template '%s': %#v", test.old, r)
				}
			}()

			assert.NoError(t, err, "error migrating template '%s'", test.old)
			assert.Equal(t, test.new, migratedTemplate, "migrating template '%s' failed", test.old)

			if migratedTemplate == test.new {
				// check that the migrated template can be evaluated
				_, err = session.Runs()[0].EvaluateTemplate(migratedTemplate)
				require.NoError(t, err, "unable to evaluate migrated template '%s'", migratedTemplate)
			}
		}
	}
}

type legacyVariables map[string]interface{}

func (v legacyVariables) Resolve(env utils.Environment, key string) types.XValue {
	key = strings.ToLower(key)
	for k, val := range v {
		if strings.ToLower(k) == key {
			return toXType(val)
		}
	}
	return nil
}

func (v legacyVariables) Describe() string {
	return "legacy vars"
}

func (v legacyVariables) Reduce(env utils.Environment) types.XPrimitive {
	return toXType(v["*"]).(types.XPrimitive)
}

func (v legacyVariables) ToXJSON(env utils.Environment) types.XText { return types.NewXText("LEGACY") }

func toXType(val interface{}) types.XValue {
	if utils.IsNil(val) {
		return nil
	}

	switch typed := val.(type) {
	case string:
		return types.NewXText(typed)
	case json.Number:
		return types.RequireXNumberFromString(string(typed))
	case map[string]interface{}:
		return legacyVariables(typed)
	}
	panic(fmt.Sprintf("unsupported type: %s", reflect.TypeOf(val)))
}

func (v legacyVariables) Migrate() legacyVariables {
	migrated := make(map[string]interface{})

	for key, val := range v {
		key = strings.ToLower(key)
		switch key {
		case "flow":
			migrated["run"] = map[string]interface{}{"results": val}
		case "contact":
			asMap, isMap := val.(map[string]interface{})
			if isMap {
				migrated["contact"] = migrateContact(asMap)
			} else {
				migrated["contact"] = val
			}
		default:
			migrated[key] = val
		}
	}
	return migrated
}

func migrateContact(contact map[string]interface{}) map[string]interface{} {
	fields := make(map[string]interface{})
	migrated := map[string]interface{}{"fields": fields}
	for key, val := range contact {
		key = strings.ToLower(key)
		if key == "*" || key == "name" {
			migrated[key] = val
		} else {
			fields[key] = val
		}
	}
	return migrated
}

type legacyTestContext struct {
	Variables legacyVariables `json:"variables"`
	Timezone  string          `json:"timezone"`
	DateStyle string          `json:"date_style"`
	Now       *time.Time      `json:"now"`
}

type legacyTest struct {
	Template  string            `json:"template"`
	Context   legacyTestContext `json:"context"`
	URLEncode bool              `json:"url_encode"`
	Output    string            `json:"output"`
	Errors    []string          `json:"errors"`
}

// TestLegacyTests runs the tests from https://github.com/rapidpro/expressions,  migrating each template first
func TestLegacyTests(t *testing.T) {
	legacyTestData, err := ioutil.ReadFile("testdata/legacy_tests.json")
	require.NoError(t, err)

	var tests []legacyTest
	d := json.NewDecoder(bytes.NewReader(legacyTestData))
	d.UseNumber()
	err = d.Decode(&tests)
	require.NoError(t, err)

	for _, tc := range tests {
		migratedTemplate, err := expressions.MigrateTemplate(tc.Template, expressions.ExtraAsFunction)

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic migrating template '%s': %#v", tc.Template, r)
			}
		}()

		if err != nil {
			assert.Equal(t, tc.Output, migratedTemplate, "migrated template should match input on error")
		} else {
			// evaluate the migrated template
			tz, err := time.LoadLocation(tc.Context.Timezone)
			require.NoError(t, err)

			env := test.NewTestEnvironment(utils.DateFormatDayMonthYear, tz, tc.Context.Now)

			migratedVars := tc.Context.Variables.Migrate()
			migratedVarsJSON, _ := json.Marshal(migratedVars)

			_, err = excellent.EvaluateTemplateAsString(env, migratedVars, migratedTemplate, tc.URLEncode, runs.RunContextTopLevels)

			if len(tc.Errors) > 0 {
				assert.Error(t, err, "expecting error evaluating template '%s' (migrated from '%s') with context %s", migratedTemplate, tc.Template, migratedVarsJSON)
			} else {
				// TODO enable checking of output
				//assert.Equal(t, tc.Output, output, "output mismatch for template '%s' (migrated from '%s') with context %s", migratedTemplate, tc.Template, migratedVarsJSON)
			}
		}
	}
}
