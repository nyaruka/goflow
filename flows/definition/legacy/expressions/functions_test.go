package expressions_test

import (
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows/definition/legacy/expressions"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrateFunctionCall(t *testing.T) {
	var tests = []struct {
		old string
		new string
		val string
	}{
		{old: `@(ABS(-1))`, new: `@(abs(-1))`, val: `1`},
		{old: `@(Abs(5))`, new: `@(abs(5))`, val: `5`},

		{old: `@(AND(contact.age > 30, flow.amount < 5))`, new: `@(and(fields.age > 30, results.amount < 5))`, val: `false`},
		{old: `@(AND(contact.gender = "F", contact.age >= 18))`, new: `@(and(fields.gender = "F", fields.age >= 18))`, val: `false`},

		{old: `@(AVERAGE(1, 2, 3, 4, 5))`, new: `@(mean(1, 2, 3, 4, 5))`, val: `3`},

		{old: `@(CHAR(10))`, new: `@(char(10))`, val: "\n"},

		{old: `@(CLEAN(contact.gender))`, new: `@(clean(fields.gender))`, val: `Male`},

		{old: `@(CODE("A"))`, new: `@(code("A"))`, val: `65`},

		{old: `@(CONCATENATE(contact.name, " ", contact.language))`, new: `@(contact.name & " " & contact.language)`, val: `Ryan Lewis eng`},

		{old: `@(DATE(2012, 12, 25))`, new: `@(date_from_parts(2012, 12, 25))`, val: `2012-12-25`},

		{old: `@(DATEDIF(contact.join_date, date.now, "M"))`, new: `@(datetime_diff(fields.join_date, now(), "M"))`, val: `4`},

		{old: `@(DATEVALUE("2012-02-03"))`, new: `@(date("2012-02-03"))`, val: `2012-02-03`},

		{old: `@(DAY(contact.join_date))`, new: `@(format_date(fields.join_date, "D"))`, val: `1`},

		{old: `@(DAYS("2016-02-28", "2015-02-28"))`, new: `@(datetime_diff("2015-02-28", "2016-02-28", "D"))`, val: `365`},
		{old: `@(DAYS("2016-02-28", "2016-02-29"))`, new: `@(datetime_diff("2016-02-29", "2016-02-28", "D"))`, val: `-1`},

		{old: `@(EDATE("2012-02-03", 1))`, new: `@(datetime_add("2012-02-03", 1, "M"))`, val: `2012-03-03T00:00:00.000000-05:00`},

		{old: `@(EPOCH(NOW()))`, new: `@(epoch(now()))`},

		{old: `@(FALSE())`, new: `@(false)`, val: `false`},

		{old: `@(FIELD(flow.favorite_color, 2, ","))`, new: `@(field(results.favorite_color, 1, ","))`, val: ``},
		{old: `@(FIELD(flow.favorite_color, child.age, ","))`, new: `@(field(results.favorite_color, child.results.age - 1, ","))`, val: ``},

		{old: `@(FIRST_WORD(flow.favorite_color))`, new: `@(word(results.favorite_color, 0))`, val: `red`},
		{old: `@(FIRST_WORD(WORD_SLICE("bee cat dog elf", 2, 4)))`, new: `@(word(word_slice("bee cat dog elf", 1, 3), 0))`, val: `cat`},

		{old: `@(FIXED(contact.age, 3, false))`, new: `@(format_number(fields.age, 3, false))`, val: `23.000`},
		{old: `@(FIXED(contact.age, 3))`, new: `@(format_number(fields.age, 3))`, val: `23.000`},
		{old: `@(FIXED(contact.age))`, new: `@(format_number(fields.age, 2))`, val: `23.00`},

		{old: `@(HOUR(NOW()))`, new: `@(format_datetime(now(), "tt"))`},

		{old: `@(IF(contact.gender = "M", "Sir", "Madam"))`, new: `@(if(fields.gender = "M", "Sir", "Madam"))`, val: `Madam`},

		{old: `@(INT(contact.age))`, new: `@(round_down(fields.age))`, val: `23`},

		{old: `@(LEFT(contact.name, 4))`, new: `@(text_slice(contact.name, 0, 4))`, val: `Ryan`},

		{old: `@(LEN(contact.first_name))`, new: `@(text_length(contact.first_name))`, val: `4`},

		{old: `@(LOWER(contact.first_name))`, new: `@(lower(contact.first_name))`, val: `ryan`},

		{old: `@(MAX(child.age, 10))`, new: `@(max(child.results.age, 10))`, val: `23`},

		{old: `@(MIN(child.age, 10))`, new: `@(min(child.results.age, 10))`, val: `10`},

		{old: `@(MINUTE(NOW()))`, new: `@(format_datetime(now(), "m"))`},

		{old: `@(MOD(103, 4))`, new: `@(mod(103, 4))`, val: `3`},

		{old: `@(MONTH(NOW()))`, new: `@(format_date(now(), "M"))`},

		{old: `@(NOW())`, new: `@(now())`},

		{old: `@(OR(contact.gender = "M", contact.gender = "F", contact.gender = "NB"))`, new: `@(or(fields.gender = "M", fields.gender = "F", fields.gender = "NB"))`},

		{old: `@(POWER(2, 3))`, new: `@(2 ^ 3)`, val: `8`},

		{old: `@(PROPER(contact))`, new: `@(title(contact))`, val: `Ryan Lewis`},

		{old: `@(RAND())`, new: `@(rand())`},

		{old: `@(RANDBETWEEN(1, 10))`, new: `@(rand_between(1, 10))`},

		{old: `@(REGEX_GROUP(flow.favorite_color, "\w(\w+)", 1))`, new: `@(regex_match(results.favorite_color, "\w(\w+)", 1))`},
		{old: `@(REGEX_GROUP(flow.favorite_color, "\w\w+"))`, new: `@(regex_match(results.favorite_color, "\w\w+"))`},

		{old: `@(REMOVE_FIRST_WORD(flow.favorite_color))`, new: `@(remove_first_word(results.favorite_color))`},

		{old: `@(REPT("*", 10))`, new: `@(repeat("*", 10))`},

		{old: `@(RIGHT(contact.name, 4))`, new: `@(text_slice(contact.name, -4))`, val: "ewis"},

		{old: `@(ROUND(9.4378, 3))`, new: `@(round(9.4378, 3))`},
		{old: `@(ROUND(9.4378))`, new: `@(round(9.4378))`},

		{old: `@(ROUNDDOWN(9.4378, 3))`, new: `@(round_down(9.4378, 3))`},
		{old: `@(ROUNDDOWN(9.4378))`, new: `@(round_down(9.4378))`},

		{old: `@(ROUNDUP(9.4378, 3))`, new: `@(round_up(9.4378, 3))`},
		{old: `@(ROUNDUP(9.4378))`, new: `@(round_up(9.4378))`},

		{old: `@(SECOND(NOW()))`, new: `@(format_datetime(now(), "s"))`},

		{old: `@(SUM(contact.age, child.age))`, new: `@(fields.age + child.results.age)`},

		{old: `@(TRUE())`, new: `@(true)`},

		{old: `@(WEEKDAY(TODAY()))`, new: `@(weekday(today()) + 1)`},

		{old: `@(WORD_COUNT(flow.favorite_color, FALSE))`, new: `@(word_count(results.favorite_color, NULL))`},
		{old: `@(WORD_COUNT(flow.favorite_color, TRUE))`, new: `@(word_count(results.favorite_color, " \t"))`},
		{old: `@(WORD_COUNT(flow.favorite_color))`, new: `@(word_count(results.favorite_color))`},

		{old: `@(WORD_SLICE(flow.favorite_color, 2, 4, FALSE))`, new: `@(word_slice(results.favorite_color, 1, 3, NULL))`},
		{old: `@(WORD_SLICE(flow.favorite_color, 2, 4, TRUE))`, new: `@(word_slice(results.favorite_color, 1, 3, " \t"))`},
		{old: `@(WORD_SLICE(flow.favorite_color, 2, 4))`, new: `@(word_slice(results.favorite_color, 1, 3))`},
		{old: `@(WORD_SLICE(flow.favorite_color, 2))`, new: `@(word_slice(results.favorite_color, 1))`},

		{old: `@(WORD(flow.favorite_color, 1, FALSE))`, new: `@(word(results.favorite_color, 0, NULL))`},
		{old: `@(WORD(flow.favorite_color, 1, TRUE))`, new: `@(word(results.favorite_color, 0, " \t"))`},
		{old: `@(WORD(flow.favorite_color, 1))`, new: `@(word(results.favorite_color, 0))`},
		{old: `@(WORD(flow.favorite_color, child.age - 22))`, new: `@(word(results.favorite_color, legacy_add(child.results.age, -22) - 1))`},

		{old: `@(YEAR(date.now))`, new: `@(format_date(now(), "YYYY"))`},
		{old: `@(YEAR(NOW()))`, new: `@(format_date(now(), "YYYY"))`},
	}

	server := test.NewTestHTTPServer(49991)
	defer server.Close()

	session, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	require.NoError(t, err)

	defer dates.SetNowSource(dates.DefaultNowSource)
	dates.SetNowSource(dates.NewFixedNowSource(time.Date(2018, 4, 11, 13, 24, 30, 123456000, time.UTC)))

	for _, tc := range tests {
		migratedTemplate, err := expressions.MigrateTemplate(tc.old, nil)

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic migrating function call '%s': %#v", tc.old, r)
			}
		}()

		assert.NoError(t, err, "error migrating function call '%s'", tc.old)
		assert.Equal(t, tc.new, migratedTemplate, "migrating function call '%s' failed", tc.old)

		if migratedTemplate == tc.new {
			// check that the migrated template can be evaluated
			val, err := session.Runs()[0].EvaluateTemplate(migratedTemplate)
			require.NoError(t, err, "unable to evaluate migrated function call '%s'", migratedTemplate)

			if tc.val != "" {
				assert.Equal(t, tc.val, val, "unexpected evaluated value for migrated function call '%s'", tc.old)
			}
		}
	}
}
