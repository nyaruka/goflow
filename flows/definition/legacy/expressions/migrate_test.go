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

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows/definition/legacy/expressions"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testTemplate struct {
	old           string
	new           string
	defaultToSelf bool
	dontEval      bool // do the migration test but don't try to evaluate the result
}

var tests = []testTemplate{

	// contact variables
	{old: `@contact`, new: `@contact`},
	{old: `@CONTACT`, new: `@contact`},
	{old: `@contact.uuid`, new: `@contact.uuid`},
	{old: `@contact.id`, new: `@contact.id`},
	{old: `@contact.name`, new: `@contact.name`},
	{old: `@contact.NAME`, new: `@contact.name`},
	{old: `@contact.first_name`, new: `@contact.first_name`},
	{old: `@contact.gender`, new: `@fields.gender`},
	{old: `@contact.groups`, new: `@(join(contact.groups, ","))`},
	{old: `@contact.language`, new: `@contact.language`},
	{old: `@contact.created_on`, new: `@contact.created_on`},

	// contact URN variables
	{old: `@contact.tel`, new: `@(format_urn(urns.tel))`},
	{old: `@contact.tel.display`, new: `@(format_urn(urns.tel))`},
	{old: `@contact.tel.scheme`, new: `@(urn_parts(urns.tel).scheme)`},
	{old: `@contact.tel.path`, new: `@(urn_parts(urns.tel).path)`},
	{old: `@contact.tel.urn`, new: `@urns.tel`},
	{old: `@contact.tel_e164`, new: `@(default(urn_parts(urns.tel).path, ""))`},
	{old: `@contact.twitterid`, new: `@(default(urn_parts(urns.twitterid).path, ""))`},
	{old: `@contact.mailto`, new: `@(default(urn_parts(urns.mailto).path, ""))`},

	// run variables
	{old: `@flow`, new: `@results`},
	{old: `@flow.favorite_color`, new: `@results.favorite_color`},
	{old: `@flow.favorite_color.category`, new: `@results.favorite_color.category_localized`},
	{old: `@flow.favorite_color.text`, new: `@results.favorite_color.input`},
	{old: `@flow.favorite_color.time`, new: `@results.favorite_color.created_on`},
	{old: `@flow.favorite_color.value`, new: `@results.favorite_color.value`},
	{old: `@flow.2factor`, new: `@(results["2factor"])`},
	{old: `@flow.2factor.value`, new: `@(results["2factor"].value)`},
	{old: `@flow.1`, new: `@(results["1"])`, dontEval: true},
	{old: `@(flow.1337)`, new: `@(results["1337"])`, dontEval: true},
	{old: `@(flow.1337.value)`, new: `@(results["1337"].value)`, dontEval: true},
	{old: `@(flow.1337.category)`, new: `@(results["1337"].category_localized)`, dontEval: true},
	{old: `@flow.contact`, new: `@contact`},
	{old: `@flow.contact.name`, new: `@contact.name`},
	{old: `@flow.contact.age`, new: `@fields.age`},

	{old: `@child`, new: `@child.results`},
	{old: `@child.age`, new: `@child.results.age`},
	{old: `@child.age.value`, new: `@child.results.age.value`},
	{old: `@child.age.category`, new: `@child.results.age.category_localized`},
	{old: `@child.age.text`, new: `@child.results.age.input`},
	{old: `@child.age.time`, new: `@child.results.age.created_on`},
	{old: `@child.contact`, new: `@child.contact`},
	{old: `@child.contact.name`, new: `@child.contact.name`},
	{old: `@child.contact.age`, new: `@child.fields.age`},

	{old: `@parent`, new: `@parent.results`},
	{old: `@parent.role`, new: `@parent.results.role`},
	{old: `@parent.role.value`, new: `@parent.results.role.value`},
	{old: `@parent.role.category`, new: `@parent.results.role.category_localized`},
	{old: `@parent.role.text`, new: `@parent.results.role.input`},
	{old: `@parent.role.time`, new: `@parent.results.role.created_on`},
	{old: `@parent.contact`, new: `@parent.contact`},
	{old: `@parent.contact.name`, new: `@parent.contact.name`},
	{old: `@parent.contact.groups`, new: `@(join(parent.contact.groups, ","))`},
	{old: `@parent.contact.gender`, new: `@parent.fields.gender`},
	{old: `@parent.contact.tel`, new: `@(format_urn(parent.urns.tel))`},
	{old: `@parent.contact.tel.display`, new: `@(format_urn(parent.urns.tel))`},
	{old: `@parent.contact.tel.scheme`, new: `@(urn_parts(parent.urns.tel).scheme)`},
	{old: `@parent.contact.tel.path`, new: `@(urn_parts(parent.urns.tel).path)`},
	{old: `@parent.contact.tel.urn`, new: `@parent.urns.tel`},
	{old: `@parent.contact.tel_e164`, new: `@(default(urn_parts(parent.urns.tel).path, ""))`},
	{old: `@parent.contact.twitterid`, new: `@(default(urn_parts(parent.urns.twitterid).path, ""))`},

	// input
	{old: `@step`, new: `@input`},
	{old: `@step.value`, new: `@input`},
	{old: `@step.text`, new: `@input.text`},
	{old: `@step.attachments`, new: `@(foreach(foreach(input.attachments, attachment_parts), extract, "url"))`},
	{old: `@step.attachments.0`, new: `@(attachment_parts(input.attachments[0]).url)`},
	{old: `@step.attachments.10`, new: `@(attachment_parts(input.attachments[10]).url)`, dontEval: true}, // out of range
	{old: `@step.time`, new: `@input.created_on`},
	{old: `@step.contact`, new: `@contact`},
	{old: `@step.contact.name`, new: `@contact.name`},
	{old: `@step.contact.age`, new: `@fields.age`},

	// dates
	{old: `@date`, new: `@(now())`},
	{old: `@date.now`, new: `@(now())`},
	{old: `@date.today`, new: `@(format_date(today()))`},
	{old: `@date.tomorrow`, new: `@(format_date(datetime_add(now(), 1, "D")))`},
	{old: `@date.yesterday`, new: `@(format_date(datetime_add(now(), -1, "D")))`},

	// channel
	{old: `@channel`, new: `@contact.channel.address`},
	{old: `@channel.address`, new: `@contact.channel.address`},
	{old: `@channel.tel`, new: `@contact.channel.address`},
	{old: `@channel.tel_e164`, new: `@contact.channel.address`},
	{old: `@channel.name`, new: `@contact.channel.name`},

	// extra
	{old: `@extra`, new: `@legacy_extra`},
	{old: `@extra.address.state`, new: `@legacy_extra.address.state`},
	{old: `@extra.results.1`, new: `@(legacy_extra.results["1"])`},
	{old: `@extra.flow.role`, new: `@parent.results.role`},

	// variables in parens
	{old: `@(contact.tel)`, new: `@(format_urn(urns.tel))`},
	{old: `@(contact.gender)`, new: `@fields.gender`},
	{old: `@(flow.favorite_color)`, new: `@results.favorite_color`},

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

	// strings
	{old: `@("")`, new: ``},
	{old: `@(" ")`, new: `@(" ")`},
	{old: `@(" "" ")`, new: `@(" \" ")`},
	{old: `@("you" & " are " & contact.gender)`, new: `@("you" & " are " & fields.gender)`},

	// number+number addition/subtraction should stay as addition/subtraction
	{old: `@(5 + 4)`, new: `@(5 + 4)`},
	{old: `@(5 - 4)`, new: `@(5 - 4)`},
	{old: `@(ABS(5) + MOD(7, 2))`, new: `@(abs(5) + mod(7, 2))`},

	// datetime+number addition should get converted to datetime_add
	{old: `@(date.now + 5)`, new: `@(datetime_add(now(), 5, "D"))`},
	{old: `@(now() + 5)`, new: `@(datetime_add(now(), 5, "D"))`},
	{old: `@(date + 5)`, new: `@(datetime_add(now(), 5, "D"))`},
	{old: `@(date.now + 5 + contact.age)`, new: `@(legacy_add(datetime_add(now(), 5, "D"), fields.age))`},

	// datetime+time addition should get converted to datetime_add
	{old: `@(date.now + TIME(2, 30, 0))`, new: `@(datetime_add(now(), format_time(time_from_parts(2, 30, 0), "tt") * 60 + format_time(time_from_parts(2, 30, 0), "m"), "m"))`},
	{old: `@(date.now - TIME(2, 30, 0))`, new: `@(datetime_add(now(), -(format_time(time_from_parts(2, 30, 0), "tt") * 60 + format_time(time_from_parts(2, 30, 0), "m")), "m"))`},

	// date+number addition should get converted to format_date(datetime_add(...))
	{old: `@(date.today + 5)`, new: `@(format_date(datetime_add(format_date(today()), 5, "D")))`},
	{old: `@(date.yesterday - 5)`, new: `@(format_date(datetime_add(format_date(datetime_add(now(), -1, "D")), -5, "D")))`},
	{old: `@(date.tomorrow - 3 + 10)`, new: `@(format_date(datetime_add(format_date(datetime_add(format_date(datetime_add(now(), 1, "D")), -3, "D")), 10, "D")))`},

	// date+time addition should get converted to replace_time
	{old: `@(today() + TIME(15, 30, 0))`, new: `@(replace_time(today(), time_from_parts(15, 30, 0)))`},
	{old: `@(TODAY()+TIMEVALUE("10:30"))`, new: `@(replace_time(today(), time("10:30")))`},
	{old: `@(DATEVALUE(date.today) + TIMEVALUE(CONCATENATE(flow.time_input, ":00")))`, new: `@(replace_time(date(format_date(today())), time(results.time_input & ":00")))`, dontEval: true},
	{old: `@(contact.join_date + TIME(2, 30, 0))`, new: `@(replace_time(fields.join_date, time_from_parts(2, 30, 0)))`},

	// legacy_add permutations
	{old: `@(contact.age + 5)`, new: `@(legacy_add(fields.age, 5))`},
	{old: `@(contact.join_date + 5 + contact.age)`, new: `@(legacy_add(legacy_add(fields.join_date, 5), fields.age))`},
	{old: `@(contact.age + 100 - 5)`, new: `@(legacy_add(legacy_add(fields.age, 100), -5))`},
	{old: `@((5 + contact.age) / 2)`, new: `@((legacy_add(5, fields.age)) / 2)`},
	{old: `@((DATEDIF(DATEVALUE("1970-01-01"), date.now, "D") * 24 * 60 * 60) + ((((HOUR(date.now)+7) * 60) + MINUTE(date.now)) * 60))`, new: `@(legacy_add((datetime_diff(date("1970-01-01"), now(), "D") * 24 * 60 * 60), ((legacy_add(((legacy_add(format_datetime(now(), "tt"), 7)) * 60), format_datetime(now(), "m"))) * 60)))`},

	// expressions that should default to themselves on error
	{old: `@("hello")`, new: `@(if(is_error("hello"), "@(\"hello\")", "hello"))`, defaultToSelf: true},
	{old: `@extra.exists`, new: `@(if(is_error(legacy_extra.exists), "@extra.exists", legacy_extra.exists))`, defaultToSelf: true},

	// non-expressions
	{old: `bob@nyaruka.com`, new: `bob@nyaruka.com`},
	{old: `@twitter_handle`, new: `@twitter_handle`},

	// misc edge cases
	{old: `@`, new: `@`},
	{old: `@contact.name...?`, new: `@contact.name...?`},
	{old: `Hi @@@flow.favorite_color @@flow.favorite_color @flow.favorite_color @nyaruka @ @`, new: `Hi @@@results.favorite_color @@flow.favorite_color @results.favorite_color @nyaruka @ @`},
}

func TestMigrateTemplate(t *testing.T) {
	for _, tc := range tests {
		tests = append(tests, testTemplate{
			old:           "Embedded " + tc.old + " text",
			new:           "Embedded " + tc.new + " text",
			defaultToSelf: tc.defaultToSelf,
			dontEval:      tc.dontEval,
		})
		tests = append(tests, testTemplate{
			old:           "Replace " + tc.old + " two " + tc.old + " times",
			new:           "Replace " + tc.new + " two " + tc.new + " times",
			defaultToSelf: tc.defaultToSelf,
			dontEval:      tc.dontEval,
		})
	}

	server := test.NewTestHTTPServer(49997)
	defer server.Close()

	session, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	require.NoError(t, err)

	for _, tc := range tests {
		options := &expressions.MigrateOptions{DefaultToSelf: tc.defaultToSelf}
		migratedTemplate, err := expressions.MigrateTemplate(tc.old, options)

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic migrating template '%s': %#v", tc.old, r)
			}
		}()

		assert.NoError(t, err, "error migrating template '%s'", tc.old)
		assert.Equal(t, tc.new, migratedTemplate, "migrating template '%s' failed", tc.old)

		if migratedTemplate == tc.new && !tc.dontEval {
			// check that the migrated template can be evaluated
			_, err := session.Runs()[0].EvaluateTemplate(migratedTemplate)
			require.NoError(t, err, "unable to evaluate migrated template '%s'", migratedTemplate)
		}
	}
}

func BenchmarkMigrateTemplate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, tc := range tests {
			options := &expressions.MigrateOptions{DefaultToSelf: tc.defaultToSelf}
			expressions.MigrateTemplate(tc.old, options)
		}
	}
}

type legacyVariables map[string]interface{}

func (v legacyVariables) Context(env envs.Environment) *types.XObject {
	entries := make(map[string]types.XValue, len(v))

	for k, val := range v {
		entries[strings.ToLower(k)] = toXType(env, val)
	}
	return types.NewXObject(entries)
}

func toXType(env envs.Environment, val interface{}) types.XValue {
	if utils.IsNil(val) {
		return nil
	}

	switch typed := val.(type) {
	case string:
		return types.NewXText(typed)
	case json.Number:
		return types.RequireXNumberFromString(string(typed))
	case map[string]interface{}:
		return legacyVariables(typed).Context(env)
	}
	panic(fmt.Sprintf("unsupported type: %s", reflect.TypeOf(val)))
}

func (v legacyVariables) Migrate() legacyVariables {
	migrated := map[string]interface{}{
		"fields":  map[string]interface{}{},
		"results": map[string]interface{}{},
	}

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
		migratedTemplate, err := expressions.MigrateTemplate(tc.Template, nil)

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

			env := envs.NewBuilder().WithDateFormat(envs.DateFormatDayMonthYear).WithTimezone(tz).Build()
			if tc.Context.Now != nil {
				dates.SetNowSource(dates.NewFixedNowSource(*tc.Context.Now))
				defer dates.SetNowSource(dates.DefaultNowSource)
			}

			migratedVars := tc.Context.Variables.Migrate().Context(env)
			migratedVarsJSON := jsonx.MustMarshal(migratedVars)

			_, err = excellent.EvaluateTemplate(env, migratedVars, migratedTemplate, nil)

			if len(tc.Errors) > 0 {
				assert.Error(t, err, "expecting error evaluating template '%s' (migrated from '%s') with context %s", migratedTemplate, tc.Template, migratedVarsJSON)
			}
		}
	}
}

func TestMigrateStringLiteral(t *testing.T) {
	assert.Equal(t, `""`, expressions.MigrateStringLiteral(`""`))
	assert.Equal(t, `"abc"`, expressions.MigrateStringLiteral(`"abc"`))
	assert.Equal(t, `"\"hello\""`, expressions.MigrateStringLiteral(`"""hello"""`))
	assert.Equal(t, `"line1\nline2\ttabbed"`, expressions.MigrateStringLiteral(`"line1\nline2\ttabbed"`))
	assert.Equal(t, `"\D\w+[\.*]"`, expressions.MigrateStringLiteral(`"\D\w+[\.*]"`))
}
