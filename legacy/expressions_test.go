package legacy

import (
	"testing"

	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testTemplate struct {
	old string
	new string

	extraAs ExtraVarsMapping
}

func TestMigrateTemplate(t *testing.T) {
	var tests = []testTemplate{

		// contact variables
		{old: `@contact`, new: `@contact`},
		{old: `@contact.uuid`, new: `@contact.uuid`},
		{old: `@contact.name`, new: `@contact.name`},
		{old: `@contact.first_name`, new: `@contact.first_name`},
		{old: `@contact.gender`, new: `@contact.fields.gender`},

		// contact URN variables
		{old: `@contact.tel`, new: `@(format_urn(contact.urns.tel.0))`},
		{old: `@contact.tel.display`, new: `@(format_urn(contact.urns.tel.0))`},
		{old: `@contact.tel.scheme`, new: `@contact.urns.tel.0.scheme`},
		{old: `@contact.tel.path`, new: `@contact.urns.tel.0.path`},
		{old: `@contact.tel.urn`, new: `@contact.urns.tel.0`},
		{old: `@contact.tel_e164`, new: `@contact.urns.tel.0.path`},
		{old: `@contact.twitterid`, new: `@(format_urn(contact.urns.twitterid.0))`},
		{old: `@contact.mailto`, new: `@(format_urn(contact.urns.mailto.0))`},

		// run variables
		{old: `@flow.favorite_color`, new: `@run.results.favorite_color`},
		{old: `@flow.favorite_color.category`, new: `@run.results.favorite_color.category_localized`},
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
		//{old: `@date.tomorrow`, new: `@(tomorrow())`}, // TODO
		//{old: `@date.yesterday`, new: `@(yesterday())`}, // TODO

		// variables in parens
		{old: `@(contact.tel)`, new: `@(format_urn(contact.urns.tel.0))`},
		{old: `@(contact.gender)`, new: `@(contact.fields.gender)`},
		{old: `@(flow.favorite_color)`, new: `@(run.results.favorite_color)`},

		// arithmetic
		{old: `@(1 + 2)`, new: `@(1 + 2)`},
		{old: `@(1 - 2)`, new: `@(1 - 2)`},
		{old: `@(-2)`, new: `@(-2)`},
		{old: `@(2 ^ 4)`, new: `@(2 ^ 4)`},
		{old: `@(2 * 4)`, new: `@(2 * 4)`},
		{old: `@(2 / 4)`, new: `@(2 / 4)`},

		// comparisons
		{old: `@(1 < 4)`, new: `@(1 < 4)`},
		{old: `@(1 <= 4)`, new: `@(1 <= 4)`},
		{old: `@(1 > 4)`, new: `@(1 > 4)`},
		{old: `@(1 >= 4)`, new: `@(1 >= 4)`},

		// string concatenation
		{old: `@("you" & " are " & contact.gender)`, new: `@("you" & " are " & contact.fields.gender)`},

		// functions
		{old: "@(REMOVE_FIRST_WORD(flow.favorite_color))", new: "@(remove_first_word(run.results.favorite_color))"},
		//{old: "@(WORD_SLICE(flow.favorite_color, 2))", new: "@(word_slice(run.results.favorite_color, 2 - 1))"}, // TODO
		{old: "@(WORD_SLICE(flow.favorite_color, 2, 4))", new: "@(word_slice(run.results.favorite_color, 2 - 1, 4 - 1))"},
		{old: "@(FIELD(flow.favorite_color, 2, \",\"))", new: "@(field(run.results.favorite_color, 2 - 1, \",\"))"},
		{old: "@(FIELD(flow.favorite_color, child.age, \",\"))", new: "@(field(run.results.favorite_color, child.results.age - 1, \",\"))"},
		{old: "@(FIRST_WORD(WORD_SLICE(\"bee cat dog elf\", 2, 4)))", new: "@(word(word_slice(\"bee cat dog elf\", 2 - 1, 4 - 1), 0))"},

		{old: "@(FIRST_WORD(flow.favorite_color))", new: "@(word(run.results.favorite_color, 0))"},
		{old: "@(WORD(flow.favorite_color, child.age - 22))", new: "@(word(run.results.favorite_color, legacy_add(child.results.age, -22) - 1))"},
		{old: "@(WORD(flow.favorite_color, 1))", new: "@(word(run.results.favorite_color, 1 - 1))"},
		{old: "@(ABS(-5))", new: "@(abs(-5))"},
		{old: "@(AVERAGE(1, 2, 3, 4, 5))", new: "@(mean(1, 2, 3, 4, 5))"},
		{old: "@(AND(contact.age > 30, flow.amount < 5))", new: "@(and(contact.fields.age > 30, run.results.amount < 5))"},
		{old: "@(DATEVALUE(\"2012-02-03\"))", new: "@(date(\"2012-02-03\"))"},
		{old: "@(EDATE(\"2012-02-03\", 1))", new: "@(date_add(\"2012-02-03\", 1, \"M\"))"},
		{old: "@(DATEDIF(contact.join_date, date.now, \"M\"))", new: "@(date_diff(contact.fields.join_date, now(), \"M\"))"},
		{old: "@(DAYS(\"2016-02-28\", \"2015-02-28\"))", new: "@(date_diff(\"2016-02-28\", \"2015-02-28\", \"D\"))"},
		{old: "@(DAY(contact.join_date))", new: "@(format_date(contact.fields.join_date, \"D\"))"},
		{old: "@(HOUR(NOW()))", new: "@(format_date(now(), \"h\"))"},
		{old: "@(MINUTE(NOW()))", new: "@(format_date(now(), \"m\"))"},
		{old: "@(MONTH(NOW()))", new: "@(format_date(now(), \"M\"))"},
		{old: "@(NOW())", new: "@(now())"},
		{old: "@(SECOND(NOW()))", new: "@(format_date(now(), \"s\"))"},

		// date addition should get converted to date_add
		{old: "@(date.now + 5)", new: "@(date_add(now(), 5, \"D\"))"},
		{old: "@(now() + 5)", new: "@(date_add(now(), 5, \"D\"))"},
		{old: "@(date + 5)", new: "@(date_add(now(), 5, \"D\"))"},
		{old: "@(date.now + 5 + contact.age)", new: "@(legacy_add(date_add(now(), 5, \"D\"), contact.fields.age))"},

		// legacy_add permutations
		{old: "@(contact.age + 5)", new: "@(legacy_add(contact.fields.age, 5))"},
		{old: "@(contact.join_date + 5 + contact.age)", new: "@(legacy_add(legacy_add(contact.fields.join_date, 5), contact.fields.age))"},
		{old: "@(contact.age + 100 - 5)", new: "@(legacy_add(legacy_add(contact.fields.age, 100), -5))"},
		// {old: "@(date.yesterday - 3 + 10)", new: "@(legacy_add(date_add(yesterday(), -3, \"D\"), 10))"}, // TODO

		{old: "@(3 + date.now)", new: "@(date_add(now(), 3, \"D\"))"},
		//{old: "@(date.yesterday - 3)", new: "@(date_add(yesterday(), -3, \"D\"))"},  // TODO
		{old: "@(date.now + TIME(2, 30, 0))", new: "@(date_add(now(), 9000, \"s\"))"},
		{old: "@(TIME(0, 1, 5) + contact.join_date)", new: "@(date_add(contact.fields.join_date, 65, \"s\"))"},
		{old: "@(contact.join_date - TIME(0,0,12))", new: "@(date_add(contact.fields.join_date, -12, \"s\"))"},

		// TODO: beware different org format, need a like function
		{old: "@(DATE(2012, 12, 25))", new: "@(date(\"2012-12-25\"))"},
		{old: "@(5 * contact.age)", new: "@(5 * contact.fields.age)"},

		{old: "@((5 + contact.age) / 2)", new: "@((legacy_add(5, contact.fields.age)) / 2)"},
		{old: "@(WEEKDAY(TODAY()))", new: "@(weekday(today()))"},
		{old: "@(YEAR(date.now))", new: "@(format_date(now(), \"YYYY\"))"},

		// booleans and conditionals
		// TODO
		//{old: "@(AND(contact.gender = \"F\", contact.age >= 18))", new: "@(and(contact.fields.gender = \"F\", contact.fields.age >= 18))"},
		//{old: "@(IF(contact.gender = \"M\", \"Sir\", \"Madam\"))", new: "@(if(contact.fields.gender = \"M\", \"Sir\", \"Madam\"))"},
		//{old: "@(OR(contact.state = \"GA\", contact.state = \"WA\", contact.state = \"IN\"))", new: "@(or(contact.fields.state = \"GA\", contact.fields.state == \"WA\", contact.fields.state == \"IN\"))"},

		// math functions
		{old: "@(ABS(-1))", new: "@(abs(-1))"},
		{old: "@(MAX(child.age, 10))", new: "@(max(child.results.age, 10))"},
		{old: "@(MIN(child.age, 10))", new: "@(min(child.results.age, 10))"},
		{old: "@(POWER(2, 3))", new: "@(2 ^ 3)"},
		{old: "@(RAND())", new: "@(rand())"},
		{old: "@(RANDBETWEEN(1, 10))", new: "@(rand_between(1, 10))"},
		{old: "@(ROUND(9.4378, 3))", new: "@(round(9.4378, 3))"},
		//{old: "@(ROUNDUP(9.4378))", new: "@(round_up(9.4378, 3))"}, // TODO
		//{old: "@(ROUNDDOWN(9.4378, 3))", new: "@(round_down(9.4378, 3))"},  // TODO
		{old: "@(SUM(contact.age, child.age))", new: "@(contact.fields.age + child.results.age)"},
		{old: "@(CHAR(10))", new: "@(char(10))"},
		{old: "@(CLEAN(contact.gender))", new: "@(clean(contact.fields.gender))"},
		{old: "@(CODE(\"A\"))", new: "@(code(\"A\"))"},
		{old: "@(CONCATENATE(contact.first_name, \" \", contact.language))", new: "@(contact.first_name & \" \" & contact.language)"},

		//{old: "@(FIXED(contact.balance))", new: "@(format_num(contact.fields.balance))"}, // TODO
		//{old: "@(FIXED(contact.balance, 2))", new: "@(format_num(contact.fields.balance, 2))"}, // TODO
		{old: "@(FIXED(contact.age, 2, false))", new: "@(format_num(contact.fields.age, 2, false))"},
		{old: "@(INT(contact.age))", new: "@(round_down(contact.fields.age))"},
		{old: "@(LEFT(contact.name, 4))", new: "@(left(contact.name, 4))"},
		{old: "@(RIGHT(contact.name, 4))", new: "@(right(contact.name, 4))"},
		{old: "@(LEN(contact.first_name))", new: "@(length(contact.first_name))"},
		{old: "@(LOWER(contact.first_name))", new: "@(lower(contact.first_name))"},
		{old: "@(mod(103, 4))", new: "@(mod(103, 4))"},

		{old: "@(PROPER(contact))", new: "@(title(contact))"},
		{old: "@(REPT(\"*\", 10))", new: "@(repeat(\"*\", 10))"},
		// {old: "@((DATEDIF(DATEVALUE(\"01-01-1970\"), date.now, \"D\") * 24 * 60 * 60) + ((((HOUR(date.now)+7) * 60) + MINUTE(date.now)) * 60))", new: ""},

		{old: "@extra.results.0.state", new: "@run.webhook.json.results.0.state", extraAs: ExtraAsWebhookJSON},
		{old: "@extra.address.state", new: "@trigger.params.address.state", extraAs: ExtraAsTriggerParams},
		{old: "@extra.address.state", new: "@(if(is_error(run.webhook.json.address.state), trigger.params.address.state, run.webhook.json.address.state))", extraAs: ExtraAsFunction},

		// non-expressions
		{old: "bob@nyaruka.com", new: "bob@nyaruka.com"},
		{old: "@twitter_handle", new: "@twitter_handle"},
	}

	for i := range tests {
		tests = append(tests, testTemplate{old: "Embedded " + tests[i].old + " text", new: "Embedded " + tests[i].new + " text", extraAs: tests[i].extraAs})
		tests = append(tests, testTemplate{old: "Replace " + tests[i].old + " two " + tests[i].old + " times", new: "Replace " + tests[i].new + " two " + tests[i].new + " times", extraAs: tests[i].extraAs})
	}

	server, err := test.NewTestHTTPServer()
	require.NoError(t, err)

	defer server.Close()

	session, err := test.CreateTestSession(nil)
	require.NoError(t, err)

	for _, test := range tests {

		for i := 0; i < 1; i++ {
			migratedTemplate, err := MigrateTemplate(test.old, test.extraAs)

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
