package legacy

import (
	"testing"
)

type testTemplate struct {
	old string
	new string

	extraAs ExtraVarsMapping
}

func TestMigrateTemplate(t *testing.T) {
	var tests = []testTemplate{

		// contact variables
		{old: "@contact", new: "@contact"},
		{old: "@contact.uuid", new: "@contact.uuid"},
		{old: "@contact.name", new: "@contact.name"},
		{old: "@contact.first_name", new: "@contact.first_name"},
		{old: "@contact.blerg", new: "@contact.fields.blerg"},

		// contact URN variables
		{old: "@contact.tel", new: "@(format_urn(contact.urns.tel))"},
		{old: "@contact.tel.display", new: "@(format_urn(contact.urns.tel))"},
		{old: "@contact.tel.scheme", new: "@contact.urns.tel.0.scheme"},
		{old: "@contact.tel.path", new: "@contact.urns.tel.0.path"},
		{old: "@contact.tel.urn", new: "@contact.urns.tel.0"},
		{old: "@contact.tel_e164", new: "@contact.urns.tel.0.path"},
		{old: "@contact.telegram", new: "@(format_urn(contact.urns.telegram))"},
		{old: "@contact.twitter", new: "@(format_urn(contact.urns.twitter))"},
		{old: "@contact.facebook", new: "@(format_urn(contact.urns.facebook))"},
		{old: "@contact.mailto", new: "@(format_urn(contact.urns.mailto))"},

		// run variables
		{old: "@flow.blerg", new: "@run.results.blerg"},
		{old: "@flow.blerg.category", new: "@run.results.blerg.category_localized"},
		{old: "@child.blerg", new: "@child.results.blerg"},
		{old: "@child.contact", new: "@child.contact"},
		{old: "@child.contact.age", new: "@child.contact.fields.age"},
		{old: "@parent.blerg", new: "@parent.results.blerg"},
		{old: "@parent.blerg.category", new: "@parent.results.blerg.category_localized"},
		{old: "@parent.contact", new: "@parent.contact"},
		{old: "@parent.contact.name", new: "@parent.contact.name"},

		// input
		{old: "@step", new: "@run.input"},
		{old: "@step.value", new: "@run.input"},
		{old: "@step.text", new: "@run.input.text"},
		{old: "@step.attachments", new: "@run.input.attachments"},
		{old: "@step.time", new: "@run.input.created_on"},
		{old: "@step.contact", new: "@contact"},
		{old: "@step.contact.name", new: "@contact.name"},
		{old: "@step.contact.age", new: "@contact.fields.age"},

		// dates
		{old: "@date", new: "@(now())"},
		{old: "@date.now", new: "@(now())"},
		{old: "@date.today", new: "@(today())"},
		{old: "@date.tomorrow", new: "@(tomorrow())"},
		{old: "@date.yesterday", new: "@(yesterday())"},

		// variables in parens
		{old: "@(contact.tel)", new: "@(format_urn(contact.urns.tel))"},
		{old: "@(contact.blerg)", new: "@(contact.fields.blerg)"},
		{old: "@(flow.blerg)", new: "@(run.results.blerg)"},

		// arithmetic
		{old: "@(1 + 2)", new: "@(1 + 2)"},
		{old: "@(1 - 2)", new: "@(1 - 2)"},
		{old: "@(-2)", new: "@(-2)"},
		{old: "@(2 ^ 4)", new: "@(2 ^ 4)"},
		{old: "@(2 * 4)", new: "@(2 * 4)"},
		{old: "@(2 / 4)", new: "@(2 / 4)"},

		// comparisons
		{old: "@(1 < 4)", new: "@(1 < 4)"},
		{old: "@(1 <= 4)", new: "@(1 <= 4)"},
		{old: "@(1 > 4)", new: "@(1 > 4)"},
		{old: "@(1 >= 4)", new: "@(1 >= 4)"},

		// functions
		{old: "@(REMOVE_FIRST_WORD(flow.blerg))", new: "@(remove_first_word(run.results.blerg))"},
		{old: "@(WORD_SLICE(flow.blerg, 2))", new: "@(word_slice(run.results.blerg, 2))"},
		{old: "@(WORD_SLICE(flow.blerg, 2, 4))", new: "@(word_slice(run.results.blerg, 2, 4))"},
		{old: "@(WORD_SLICE(flow.blerg, 2, 4, true))", new: "@(word_slice(run.results.blerg, 2, 4, true))"},
		{old: "@(FIELD(flow.blerg, 2, \",\"))", new: "@(field(run.results.blerg, 2 - 1, \",\"))"},
		{old: "@(FIELD(flow.blerg, flow.index, \",\"))", new: "@(field(run.results.blerg, run.results.index - 1, \",\"))"},
		{old: "@(FIRST_WORD(WORD_SLICE(contact.blerg, 2, 4)))", new: "@(split(word_slice(contact.fields.blerg, 2, 4), \" \")[0])"},
		{old: "@(FIRST_WORD(WORD_SLICE(contact.blerg, 2, 4)))", new: "@(split(word_slice(contact.fields.blerg, 2, 4), \" \")[0])"},
		{old: "@(\"this\" & contact.that)", new: "@(\"this\" & contact.fields.that)"},
		{old: "@(FIRST_WORD(flow.blerg))", new: "@(split(run.results.blerg, \" \")[0])"},
		{old: "@(WORD(flow.blerg, flow.index))", new: "@(word(run.results.blerg, run.results.index - 1))"},
		{old: "@(WORD(flow.blerg, 1))", new: "@(word(run.results.blerg, 1 - 1))"},
		{old: "@(WORD_SLICE(flow.blerg, 2))", new: "@(word_slice(run.results.blerg, 2))"},
		{old: "@(ABS(-5))", new: "@(abs(-5))"},
		{old: "@(AVERAGE(1, 2, 3, 4, 5))", new: "@(mean(1, 2, 3, 4, 5))"},
		{old: "@(AND(contact.age > 30, flow.amount < 5))", new: "@(and(contact.fields.age > 30, run.results.amount < 5))"},
		{old: "@(DATEVALUE(\"2012-02-03\"))", new: "@(date(\"2012-02-03\"))"},
		{old: "@(EDATE(\"2012-02-03\", 1))", new: "@(date_add(\"2012-02-03\", \"m\", 1))"},
		{old: "@(DATEDIF(contact.join_date, date.now, \"M\"))", new: "@(date_diff(contact.fields.join_date, now(), \"M\"))"},
		{old: "@(DAYS(\"02-28-2016\", \"02-28-2015\"))", new: "@(date_diff(\"02-28-2016\", \"02-28-2015\", \"D\"))"},
		{old: "@(DAY(contact.joined_date))", new: "@(format_date(contact.fields.joined_date, \"d\"))"},
		{old: "@(HOUR(NOW()))", new: "@(format_date(now(), \"h\"))"},
		{old: "@(MINUTE(NOW()))", new: "@(format_date(now(), \"m\"))"},
		{old: "@(MONTH(NOW()))", new: "@(format_date(now(), \"M\"))"},
		{old: "@(NOW())", new: "@(now())"},
		{old: "@(SECOND(NOW()))", new: "@(format_date(now(), \"s\"))"},
		{old: "@(FIRST_WORD(WORD_SLICE(contact.blerg, 2, 4)))", new: "@(split(word_slice(contact.fields.blerg, 2, 4), \" \")[0])"},

		// date addition should get converted to date_add
		{old: "@(date.now + 5)", new: "@(date_add(now(), \"d\", 5))"},
		{old: "@(now() + 5)", new: "@(date_add(now(), \"d\", 5))"},
		{old: "@(date + 5)", new: "@(date_add(now(), \"d\", 5))"},
		{old: "@(date.now + 5 + contact.days)", new: "@(legacy_add(date_add(now(), \"d\", 5), contact.fields.days))"},

		// legacy_add permutations
		{old: "@(contact.blerg + 5)", new: "@(legacy_add(contact.fields.blerg, 5))"},
		{old: "@(contact.registered + 5 + contact.quit)", new: "@(legacy_add(legacy_add(contact.fields.registered, 5), contact.fields.quit))"},
		{old: "@(contact.blerg + contact.bloop + 5)", new: "@(legacy_add(legacy_add(contact.fields.blerg, contact.fields.bloop), 5))"},

		{old: "@(date.yesterday - 3 + 10)", new: "@(legacy_add(date_add(yesterday(), \"d\", -3), 10))"},
		{old: "@(contact.blerg - contact.bloop + 5)", new: "@(legacy_add(legacy_add(contact.fields.blerg, -contact.fields.bloop), 5))"},

		{old: "@(3 + date.now)", new: "@(date_add(now(), \"d\", 3))"},
		{old: "@(date.yesterday - 3)", new: "@(date_add(yesterday(), \"d\", -3))"},
		{old: "@(date.now + TIME(2, 30, 0))", new: "@(date_add(now(), \"s\", 9000))"},
		{old: "@(TIME(0, 1, 5) + contact.join_date)", new: "@(date_add(contact.fields.join_date, \"s\", 65))"},
		{old: "@(contact.join_date - TIME(0,0,12))", new: "@(date_add(contact.fields.join_date, \"s\", -12))"},

		// TODO: beware different org format, need a like function
		{old: "@(DATE(2012, 12, 25))", new: "@(date(\"2012-12-25\"))"},
		{old: "@(5 * contact.balance)", new: "@(5 * contact.fields.balance)"},

		{old: "@((5 + contact.balance) / 2)", new: "@((legacy_add(5, contact.fields.balance)) / 2)"},
		{old: "@(WEEKDAY(TODAY()))", new: "@(weekday(today()))"},
		{old: "@(YEAR(date.now))", new: "@(format_date(now(), \"YYYY\"))"},

		// booleans and conditionals
		{old: "@(AND(contact.gender = \"F\", contact.age >= 18))", new: "@(and(contact.fields.gender == \"F\", contact.fields.age >= 18))"},
		{old: "@(IF(contact.gender = \"M\", \"Sir\", \"Madam\"))", new: "@(if(contact.fields.gender == \"M\", \"Sir\", \"Madam\"))"},
		{old: "@(OR(contact.state = \"GA\", contact.state = \"WA\", contact.state = \"IN\"))", new: "@(or(contact.fields.state == \"GA\", contact.fields.state == \"WA\", contact.fields.state == \"IN\"))"},

		// math functions
		{old: "@(ABS(-1))", new: "@(abs(-1))"},
		{old: "@(MAX(flow.questions, 10))", new: "@(max(run.results.questions, 10))"},
		{old: "@(MIN(flow.questions, 10))", new: "@(min(run.results.questions, 10))"},
		{old: "@(POWER(2, 3))", new: "@(2 ^ 3)"},
		{old: "@(RAND())", new: "@(rand())"},
		{old: "@(RANDBETWEEN(1, 10))", new: "@(rand(1, 10))"},
		{old: "@(ROUND(9.4378, 3))", new: "@(round(9.4378, 3))"},
		{old: "@(ROUNDUP(9.4378, 3))", new: "@(round_up(9.4378, 3))"},
		{old: "@(ROUNDDOWN(9.4378, 3))", new: "@(round_down(9.4378, 3))"},
		{old: "@(SUM(contact.reports, contact.forms))", new: "@(contact.fields.reports + contact.fields.forms)"},
		{old: "@(CHAR(10))", new: "@(char(10))"},
		{old: "@(CLEAN(contact.blerg))", new: "@(clean(contact.fields.blerg))"},
		{old: "@(CODE(\"A\"))", new: "@(code(\"A\"))"},
		{old: "@(CONCATENATE(contact.first_name, \" \", contact.language))", new: "@(contact.first_name & \" \" & contact.language)"},

		{old: "@(FIXED(contact.balance))", new: "@(format_num(contact.fields.balance))"},
		{old: "@(FIXED(contact.balance, 2))", new: "@(format_num(contact.fields.balance, 2))"},
		{old: "@(FIXED(contact.balance, 2, false))", new: "@(format_num(contact.fields.balance, 2, false))"},
		{old: "@(INT(contact.balance))", new: "@(round_down(contact.fields.balance))"},
		{old: "@(LEFT(contact.account_number, 4))", new: "@(left(contact.fields.account_number, 4))"},
		{old: "@(RIGHT(contact.ssn, 4))", new: "@(right(contact.fields.ssn, 4))"},
		{old: "@(LEN(contact.first_name))", new: "@(length(contact.first_name))"},
		{old: "@(LOWER(contact.first_name))", new: "@(lower(contact.first_name))"},
		{old: "@(mod(103, 4))", new: "@(mod(103, 4))"},

		{old: "@(PROPER(contact))", new: "@(title(contact))"},
		{old: "@(REPT(\"*\", 10))", new: "@(repeat(\"*\", 10))"},
		// {old: "@((DATEDIF(DATEVALUE(\"01-01-1970\"), date.now, \"D\") * 24 * 60 * 60) + ((((HOUR(date.now)+7) * 60) + MINUTE(date.now)) * 60))", new: ""},

		{old: "@extra.blerg.foo", new: "@run.webhook.json.blerg.foo", extraAs: ExtraAsWebhookJSON},
		{old: "@extra.blerg.foo", new: "@trigger.params.blerg.foo", extraAs: ExtraAsTriggerParams},
		{old: "@extra.blerg.foo", new: "@(if(is_error(run.webhook.json.blerg.foo), trigger.params.blerg.foo, run.webhook.json.blerg.foo))", extraAs: ExtraAsFunction},

		// non-expressions
		{old: "bob@nyaruka.com", new: "bob@nyaruka.com"},
		{old: "@twitter_handle", new: "@twitter_handle"},
	}

	for i := range tests {
		tests = append(tests, testTemplate{old: "Embedded " + tests[i].old + " text", new: "Embedded " + tests[i].new + " text", extraAs: tests[i].extraAs})
		tests = append(tests, testTemplate{old: "Replace " + tests[i].old + " two " + tests[i].old + " times", new: "Replace " + tests[i].new + " two " + tests[i].new + " times", extraAs: tests[i].extraAs})
	}

	for _, test := range tests {

		for i := 0; i < 1; i++ {
			translation, err := MigrateTemplate(test.old, test.extraAs)

			if err != nil {
				t.Errorf("Parse Error: %v", err)
			}

			if translation != test.new {
				t.Errorf("%s failed, expected: '%s' but got '%s'", test.old, test.new, translation)
			}
		}
	}
}
