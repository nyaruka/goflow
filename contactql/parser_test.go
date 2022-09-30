package contactql_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"
	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {
	resolver := contactql.NewMockResolver(
		[]assets.Field{
			static.NewField("f1b5aea6-6586-41c7-9020-1a6326cc6565", "age", "Age", assets.FieldTypeNumber),
			static.NewField("d66a7823-eada-40e5-9a3a-57239d4690bf", "gender", "Gender", assets.FieldTypeText),
			static.NewField("165def68-3216-4ebf-96bc-f6f1ee5bd966", "state", "State", assets.FieldTypeState),
			static.NewField("85baf5e1-b57a-46dc-a726-a84e8c4229c7", "dob", "DOB", assets.FieldTypeDatetime),
		},
		[]assets.Flow{
			static.NewFlow("f87fd7cd-e501-4394-9cff-62309af85138", "Registration", []byte(`{}`)),
		},
		[]assets.Group{
			static.NewGroup("a9b5b0a0-1098-4bc2-8384-eea09ae43e6b", "U-Reporters", ""),
		},
	)

	tests := []struct {
		text       string
		parsed     string
		err        string
		redactURNs bool
		resolver   contactql.Resolver
	}{
		// implicit conditions
		{text: `"will"`, parsed: `name ~ "will"`, resolver: resolver},
		{text: `Will`, parsed: `name ~ "Will"`, resolver: resolver},
		{text: `wil`, parsed: `name ~ "wil"`, resolver: resolver},
		{text: `wi`, parsed: `name ~ "wi"`, resolver: resolver},
		{text: `w`, parsed: `name = "w"`, resolver: resolver}, // don't have at least 1 token of >= 2 chars
		{text: `w me`, parsed: `name = "w" AND name ~ "me"`, resolver: resolver},
		{text: `w m`, parsed: `name = "w" AND name = "m"`, resolver: resolver},
		{text: `tel:+0123456566`, parsed: `tel = "+0123456566"`, resolver: resolver}, // whole query is a URN
		{text: `twitter:bobby`, parsed: `twitter = "bobby"`, resolver: resolver},
		{text: `(202) 456-1111`, parsed: `tel = "+12024561111"`, resolver: resolver}, // whole query looks like a phone number
		{text: `+12024561111`, parsed: `tel = "+12024561111"`, resolver: resolver},
		{text: ` 202.456.1111 `, parsed: `tel = "+12024561111"`, resolver: resolver},
		{text: `"+12024561111"`, parsed: `tel ~ "+12024561111"`, resolver: resolver},
		{text: `566`, parsed: `name ~ 566`, resolver: resolver}, // too short to be a phone number

		// implicit conditions with URN redaction
		{text: `will`, parsed: `name ~ "will"`, redactURNs: true, resolver: resolver},
		{text: `tel:+0123456566`, parsed: `name ~ "tel:+0123456566"`, redactURNs: true, resolver: resolver},
		{text: `twitter:bobby`, parsed: `name ~ "twitter:bobby"`, redactURNs: true, resolver: resolver},
		{text: `0123456566`, parsed: `id = 123456566`, redactURNs: true, resolver: resolver},
		{text: `+0123456566`, parsed: `id = 123456566`, redactURNs: true, resolver: resolver},
		{text: `0123-456-566`, parsed: `name ~ "0123-456-566"`, redactURNs: true, resolver: resolver},

		// explicit conditions on name
		{text: `Name=will`, parsed: `name = "will"`, resolver: resolver},
		{text: `Name=O'Shea`, parsed: `name = "O'Shea"`, resolver: resolver},
		{text: `Name ~ "felix"`, parsed: `name ~ "felix"`, resolver: resolver},
		{text: `Name HAS "Felix"`, parsed: `name ~ "Felix"`, resolver: resolver},
		{text: `name is ""`, parsed: `name = ""`, resolver: resolver},            // is not set
		{text: `name != ""`, parsed: `name != ""`, resolver: resolver},           // is set
		{text: `name != "felix"`, parsed: `name != "felix"`, resolver: resolver}, // is not equal to value
		{text: `Name ~ ""`, err: "contains operator on name requires token of minimum length 2", resolver: resolver},

		// explicit attribute conditions
		{text: `language = spa`, parsed: `language = "spa"`, resolver: resolver},
		{text: `Group IS U-Reporters`, parsed: `group = "U-Reporters"`, resolver: resolver},
		{text: `CREATED_ON>27-01-2020`, parsed: `created_on > "27-01-2020"`, resolver: resolver},

		// explicit conditions on URN
		{text: `tel=""`, parsed: `tel = ""`, resolver: resolver},
		{text: `tel!=""`, parsed: `tel != ""`, resolver: resolver},
		{text: `tel IS 233`, parsed: `tel = 233`, resolver: resolver},
		{text: `tel HAS 233`, parsed: `tel ~ 233`, resolver: resolver},
		{text: `tel ~ 23`, err: "contains operator on URN requires value of minimum length 3", resolver: resolver},
		{text: `mailto = user@example.com`, parsed: `mailto = "user@example.com"`, resolver: resolver},
		{text: `MAILTO ~ user@example.com`, parsed: `mailto ~ "user@example.com"`, resolver: resolver},
		{text: `URN=ewok`, parsed: `urn = "ewok"`, resolver: resolver},

		// explicit conditions on URN with URN redaction
		{text: `tel=""`, parsed: `tel = ""`, redactURNs: true, resolver: resolver},
		{text: `tel!=""`, parsed: `tel != ""`, redactURNs: true, resolver: resolver},
		{text: `mailto=""`, parsed: `mailto = ""`, redactURNs: true, resolver: resolver},
		{text: `mailto!=""`, parsed: `mailto != ""`, redactURNs: true, resolver: resolver},
		{text: `urn=""`, parsed: `urn = ""`, redactURNs: true, resolver: resolver},
		{text: `urn!=""`, parsed: `urn != ""`, redactURNs: true, resolver: resolver},
		{text: `tel = 233`, err: "cannot query on redacted URNs", redactURNs: true, resolver: resolver},
		{text: `tel ~ 233`, err: "cannot query on redacted URNs", redactURNs: true, resolver: resolver},
		{text: `mailto = user@example.com`, err: "cannot query on redacted URNs", redactURNs: true, resolver: resolver},
		{text: `MAILTO ~ user@example.com`, err: "cannot query on redacted URNs", redactURNs: true, resolver: resolver},
		{text: `URN=ewok`, err: "cannot query on redacted URNs", redactURNs: true, resolver: resolver},

		// field conditions
		{text: `Age IS 18`, parsed: `age = 18`, resolver: resolver},
		{text: `AGE != ""`, parsed: `age != ""`, resolver: resolver},
		{text: `age ~ 34`, err: "contains conditions can only be used with name or URN values", resolver: resolver},
		{text: `gender ~ M`, err: "contains conditions can only be used with name or URN values", resolver: resolver},

		// lt/lte/gt/gte comparisons
		{text: `Age > "18"`, parsed: `age > 18`, resolver: resolver},
		{text: `Age >= 18`, parsed: `age >= 18`, resolver: resolver},
		{text: `Age < 18`, parsed: `age < 18`, resolver: resolver},
		{text: `Age <= 18`, parsed: `age <= 18`, resolver: resolver},
		{text: `DOB > "27-01-2020"`, parsed: `dob > "27-01-2020"`, resolver: resolver},
		{text: `DOB >= 27-01-2020`, parsed: `dob >= "27-01-2020"`, resolver: resolver},
		{text: `DOB < 27/01/2020`, parsed: `dob < "27/01/2020"`, resolver: resolver},
		{text: `DOB <= 27.01.2020`, parsed: `dob <= "27.01.2020"`, resolver: resolver},
		{text: `name > Will`, err: "comparisons with > can only be used with date and number fields", resolver: resolver},
		{text: `tel < 23425`, err: "comparisons with < can only be used with date and number fields", resolver: resolver},

		// implicit combinations
		{text: `will felix`, parsed: `name ~ "will" AND name ~ "felix"`, resolver: resolver},
		{text: `will +123456566`, parsed: `name ~ "will" AND tel ~ "+123456566"`, resolver: resolver},

		// explicit combinations...
		{text: `will and felix`, parsed: `name ~ "will" AND name ~ "felix"`, resolver: resolver}, // explicit AND
		{text: `will AND felix AND matt`, parsed: `name ~ "will" AND name ~ "felix" AND name ~ "matt"`, resolver: resolver},
		{text: `will or felix or matt`, parsed: `name ~ "will" OR name ~ "felix" OR name ~ "matt"`, resolver: resolver},
		{text: `name = will AND age > 18 AND tickets = 0`, parsed: `name = "will" AND age > 18 AND tickets = 0`, resolver: resolver},
		{text: `name = will OR age > 18 AND tickets = 0`, parsed: `name = "will" OR (age > 18 AND tickets = 0)`, resolver: resolver},
		{text: `name = will AND age > 18 OR tickets = 0`, parsed: `(name = "will" AND age > 18) OR tickets = 0`, resolver: resolver},
		{text: `(name = will AND age > 18) AND tickets = 0`, parsed: `name = "will" AND age > 18 AND tickets = 0`, resolver: resolver},
		{text: `(name = will AND age > 18) AND (tickets = 0 AND language = eng)`, parsed: `name = "will" AND age > 18 AND tickets = 0 AND language = "eng"`, resolver: resolver},
		{text: `name=will or Name ~ "felix"`, parsed: `name = "will" OR name ~ "felix"`, resolver: resolver},
		{text: `Name is will or Name has felix`, parsed: `name = "will" OR name ~ "felix"`, resolver: resolver}, // operator aliases
		{text: `will or Name ~ "felix"`, parsed: `name ~ "will" OR name ~ "felix"`, resolver: resolver},

		// boolean operator precedence is AND before OR, even when AND is implicit
		{text: `will and felix or matt and amber`, parsed: `(name ~ "will" AND name ~ "felix") OR (name ~ "matt" AND name ~ "amber")`, resolver: resolver},
		{text: `will and felix or matt amber`, parsed: `(name ~ "will" AND name ~ "felix") OR (name ~ "matt" AND name ~ "amber")`, resolver: resolver},

		// boolean combinations can themselves be combined
		{
			text:     `(Age < 18 and Gender = "male") or (Age > 18 and Gender = "female")`,
			parsed:   `(age < 18 AND gender = "male") OR (age > 18 AND gender = "female")`,
			resolver: resolver,
		},
		{
			text:     `age > 10 and age < 20 or age > 30 and age < 40 or age > 50 and age < 60`,
			parsed:   `(age > 10 AND age < 20) OR (age > 30 AND age < 40) OR (age > 50 AND age < 60)`,
			resolver: resolver,
		},

		{text: `xyz != ""`, err: "can't resolve 'xyz' to attribute, scheme or field", resolver: resolver},
		{text: `group != "Gamers"`, err: "'Gamers' is not a valid group name", resolver: resolver},
		{text: `flow = "Catch All"`, err: "'Catch All' is not a valid flow name", resolver: resolver},
		{text: `status = "xxxx"`, err: "'xxxx' is not a valid contact status", resolver: resolver},
		{text: `language = "xxxx"`, err: "'xxxx' is not a valid language code", resolver: resolver},

		{text: `name = "O\"Leary"`, parsed: `name = "O\"Leary"`, resolver: resolver}, // string unquoting

		// = supported for everything
		{text: `uuid = f81d1eb5-215d-4ae8-90fa-38b3f2d6e328`, parsed: `uuid = "f81d1eb5-215d-4ae8-90fa-38b3f2d6e328"`, resolver: resolver},
		{text: `id = 02352`, parsed: `id = 02352`, resolver: resolver},
		{text: `name = felix`, parsed: `name = "felix"`, resolver: resolver},
		{text: `status = ACTIVE`, parsed: `status = "ACTIVE"`, resolver: resolver},
		{text: `language = eng`, parsed: `language = "eng"`, resolver: resolver},
		{text: `tickets = 0`, parsed: `tickets = 0`, resolver: resolver},
		{text: `group = u-reporters`, parsed: `group = "u-reporters"`, resolver: resolver},
		{text: `flow = registration`, parsed: `flow = "registration"`, resolver: resolver},
		{text: `created_on = 20-02-2020`, parsed: `created_on = "20-02-2020"`, resolver: resolver},
		{text: `tel = 02352`, parsed: `tel = 02352`, resolver: resolver},
		{text: `urn = 02352`, parsed: `urn = 02352`, resolver: resolver},
		{text: `age = 18`, parsed: `age = 18`, resolver: resolver},
		{text: `gender = male`, parsed: `gender = "male"`, resolver: resolver},
		{text: `dob = 20-02-2020`, parsed: `dob = "20-02-2020"`, resolver: resolver},
		{text: `state = Pichincha`, parsed: `state = "Pichincha"`, resolver: resolver},

		// != supported for everything
		{text: `uuid != f81d1eb5-215d-4ae8-90fa-38b3f2d6e328`, parsed: `uuid != "f81d1eb5-215d-4ae8-90fa-38b3f2d6e328"`, resolver: resolver},
		{text: `id != 02352`, parsed: `id != 02352`, resolver: resolver},
		{text: `name != felix`, parsed: `name != "felix"`, resolver: resolver},
		{text: `status != blocked`, parsed: `status != "blocked"`, resolver: resolver},
		{text: `language != eng`, parsed: `language != "eng"`, resolver: resolver},
		{text: `group != u-reporters`, parsed: `group != "u-reporters"`, resolver: resolver},
		{text: `flow != registration`, parsed: `flow != "registration"`, resolver: resolver},
		{text: `tickets != 0`, parsed: `tickets != 0`, resolver: resolver},
		{text: `created_on != 20-02-2020`, parsed: `created_on != "20-02-2020"`, resolver: resolver},
		{text: `tel != 02352`, parsed: `tel != 02352`, resolver: resolver},
		{text: `urn != 02352`, parsed: `urn != 02352`, resolver: resolver},
		{text: `age != 18`, parsed: `age != 18`, resolver: resolver},
		{text: `gender != male`, parsed: `gender != "male"`, resolver: resolver},
		{text: `dob != 20-02-2020`, parsed: `dob != "20-02-2020"`, resolver: resolver},
		{text: `state != Pichincha`, parsed: `state != "Pichincha"`, resolver: resolver},

		// = "" supported for name, language, flow, groups, fields, urns
		{text: `uuid = ""`, err: "can't check whether 'uuid' is set or not set", resolver: resolver},
		{text: `id = ""`, err: "can't check whether 'id' is set or not set", resolver: resolver},
		{text: `name = ""`, parsed: `name = ""`, resolver: resolver},
		{text: `status = ""`, err: "can't check whether 'status' is set or not set", resolver: resolver},
		{text: `language = ""`, parsed: `language = ""`, resolver: resolver},
		{text: `group = ""`, parsed: `group = ""`, resolver: resolver},
		{text: `flow = ""`, parsed: `flow = ""`, resolver: resolver},
		{text: `tickets = ""`, err: "can't check whether 'tickets' is set or not set", resolver: resolver},
		{text: `created_on = ""`, err: "can't check whether 'created_on' is set or not set", resolver: resolver},
		{text: `tel = ""`, parsed: `tel = ""`, resolver: resolver},
		{text: `urn = ""`, parsed: `urn = ""`, resolver: resolver},
		{text: `age = ""`, parsed: `age = ""`, resolver: resolver},
		{text: `gender = ""`, parsed: `gender = ""`, resolver: resolver},
		{text: `dob = ""`, parsed: `dob = ""`, resolver: resolver},
		{text: `state = ""`, parsed: `state = ""`, resolver: resolver},

		// ~ only supported for name and URNs
		{text: `uuid ~ 02352`, err: "contains conditions can only be used with name or URN values", resolver: resolver},
		{text: `id ~ 02352`, err: "contains conditions can only be used with name or URN values", resolver: resolver},
		{text: `name ~ felix`, parsed: `name ~ "felix"`, resolver: resolver},
		{text: `status ~ sto`, err: "contains conditions can only be used with name or URN values", resolver: resolver},
		{text: `language ~ eng`, err: "contains conditions can only be used with name or URN values", resolver: resolver},
		{text: `group ~ porters`, err: "contains conditions can only be used with name or URN values", resolver: resolver},
		{text: `flow ~ reg`, err: "contains conditions can only be used with name or URN values", resolver: resolver},
		{text: `tickets ~ 12`, err: "contains conditions can only be used with name or URN values", resolver: resolver},
		{text: `created_on ~ 2018`, err: "contains conditions can only be used with name or URN values", resolver: resolver},
		{text: `tel ~ 02352`, parsed: `tel ~ 02352`, resolver: resolver},
		{text: `urn ~ 02352`, parsed: `urn ~ 02352`, resolver: resolver},
		{text: `age ~ 18`, err: "contains conditions can only be used with name or URN values", resolver: resolver},
		{text: `gender ~ mal`, err: "contains conditions can only be used with name or URN values", resolver: resolver},
		{text: `dob ~ 20-02-2020`, err: "contains conditions can only be used with name or URN values", resolver: resolver},
		{text: `state ~ Pichincha`, err: "contains conditions can only be used with name or URN values", resolver: resolver},

		// > >= < <= only supported for numeric or date fields
		{text: `uuid > 02352`, err: "comparisons with > can only be used with date and number fields", resolver: resolver},
		{text: `id > 02352`, err: "comparisons with > can only be used with date and number fields", resolver: resolver},
		{text: `name > felix`, err: "comparisons with > can only be used with date and number fields", resolver: resolver},
		{text: `status > blo`, err: "comparisons with > can only be used with date and number fields", resolver: resolver},
		{text: `language > eng`, err: "comparisons with > can only be used with date and number fields", resolver: resolver},
		{text: `group > reporters`, err: "comparisons with > can only be used with date and number fields", resolver: resolver},
		{text: `flow > registration`, err: "comparisons with > can only be used with date and number fields", resolver: resolver},
		{text: `tickets > 0`, parsed: `tickets > 0`, resolver: resolver},
		{text: `created_on > 20-02-2020`, parsed: `created_on > "20-02-2020"`, resolver: resolver},
		{text: `tel > 02352`, err: "comparisons with > can only be used with date and number fields", resolver: resolver},
		{text: `urn > 02352`, err: "comparisons with > can only be used with date and number fields", resolver: resolver},
		{text: `age > 18`, parsed: `age > 18`, resolver: resolver},
		{text: `gender > male`, err: "comparisons with > can only be used with date and number fields", resolver: resolver},
		{text: `dob > 20-02-2020`, parsed: `dob > "20-02-2020"`, resolver: resolver},
		{text: `state > Pichincha`, err: "comparisons with > can only be used with date and number fields", resolver: resolver},

		// however if we don't provide a resolver, we don't know the field type, so allowed for all
		{text: `age > 18`, parsed: `age > 18`},
		{text: `gender > male`, parsed: `gender > "male"`},
		{text: `dob > 20-02-2020`, parsed: `dob > "20-02-2020"`},
		{text: `state > Pichincha`, parsed: `state > "Pichincha"`},
	}

	for _, tc := range tests {
		redact := envs.RedactionPolicyNone
		if tc.redactURNs {
			redact = envs.RedactionPolicyURNs
		}

		env := envs.NewBuilder().WithDateFormat(envs.DateFormatDayMonthYear).WithDefaultCountry("US").WithRedactionPolicy(redact).Build()

		parsed, err := contactql.ParseQuery(env, tc.text, tc.resolver)

		if tc.err != "" {
			assert.EqualError(t, err, tc.err, "error mismatch for '%s'", tc.text)
		} else {
			assert.NoError(t, err, "unexpected error for '%s'", tc.text)
			assert.Equal(t, tc.parsed, parsed.String(), "parse mismatch for '%s'", tc.text)
		}
	}
}

func TestParsingErrors(t *testing.T) {
	tests := []struct {
		query    string
		errMsg   string
		errCode  string
		errExtra map[string]string
	}{
		{
			query:    `$`,
			errMsg:   "mismatched input '$' expecting {'(', TEXT, STRING}",
			errCode:  "unexpected_token",
			errExtra: map[string]string{"token": "$"},
		},
		{
			query:    `name = `,
			errMsg:   "mismatched input '<EOF>' expecting {TEXT, STRING}",
			errCode:  "unexpected_token",
			errExtra: map[string]string{"token": "<EOF>"},
		},
		{
			query:    `name = "x`,
			errMsg:   "extraneous input '\"' expecting {TEXT, STRING}",
			errCode:  "",
			errExtra: nil,
		},
		{
			query:    `age = XZ`,
			errMsg:   "can't convert 'XZ' to a number",
			errCode:  "invalid_number",
			errExtra: map[string]string{"value": "XZ"},
		},
		{
			query:    `dob = AB`,
			errMsg:   "can't convert 'AB' to a date",
			errCode:  "invalid_date",
			errExtra: map[string]string{"value": "AB"},
		},
		{
			query:    `created_on = AB`,
			errMsg:   "can't convert 'AB' to a date",
			errCode:  "invalid_date",
			errExtra: map[string]string{"value": "AB"},
		},
		{
			query:    `group = "Cool Kids"`,
			errMsg:   "'Cool Kids' is not a valid group name",
			errCode:  "invalid_group",
			errExtra: map[string]string{"value": "Cool Kids"},
		},
		{
			query:    `language = "zzzzzz"`,
			errMsg:   "'zzzzzz' is not a valid language code",
			errCode:  "invalid_language",
			errExtra: map[string]string{"value": "zzzzzz"},
		},
		{
			query:    `name ~ "x"`,
			errMsg:   "contains operator on name requires token of minimum length 2",
			errCode:  "invalid_partial_name",
			errExtra: map[string]string{"min_token_length": "2"},
		},
		{
			query:    `urn ~ "23"`,
			errMsg:   "contains operator on URN requires value of minimum length 3",
			errCode:  "invalid_partial_urn",
			errExtra: map[string]string{"min_value_length": "3"},
		},
		{
			query:    `uuid ~ 234`,
			errMsg:   "contains conditions can only be used with name or URN values",
			errCode:  "unsupported_contains",
			errExtra: map[string]string{"property": "uuid"},
		},
		{
			query:    `uuid > 123`,
			errMsg:   "comparisons with > can only be used with date and number fields",
			errCode:  "unsupported_comparison",
			errExtra: map[string]string{"property": "uuid", "operator": ">"},
		},
		{
			query:    `uuid = ""`,
			errMsg:   "can't check whether 'uuid' is set or not set",
			errCode:  "unsupported_setcheck",
			errExtra: map[string]string{"property": "uuid", "operator": "="},
		},
		{
			query:    `uuid != ""`,
			errMsg:   "can't check whether 'uuid' is set or not set",
			errCode:  "unsupported_setcheck",
			errExtra: map[string]string{"property": "uuid", "operator": "!="},
		},
		{
			query:    `beers = 12`,
			errMsg:   "can't resolve 'beers' to attribute, scheme or field",
			errCode:  "unknown_property",
			errExtra: map[string]string{"property": "beers"},
		},
	}

	env := envs.NewBuilder().WithDefaultCountry("US").Build()
	resolver := contactql.NewMockResolver(
		[]assets.Field{
			static.NewField("f1b5aea6-6586-41c7-9020-1a6326cc6565", "age", "Age", assets.FieldTypeNumber),
			static.NewField("3810a485-3fda-4011-a589-7320c0b8dbef", "dob", "DOB", assets.FieldTypeDatetime),
			static.NewField("d66a7823-eada-40e5-9a3a-57239d4690bf", "gender", "Gender", assets.FieldTypeText),
		},
		[]assets.Flow{},
		[]assets.Group{},
	)

	for _, tc := range tests {
		_, err := contactql.ParseQuery(env, tc.query, resolver)

		assert.EqualError(t, err, tc.errMsg, "error mismatch for '%s'", tc.query)

		qerr := err.(*contactql.QueryError)
		assert.Equal(t, tc.errCode, qerr.Code())
		assert.Equal(t, tc.errExtra, qerr.Extra())
	}
}

func TestSimplify(t *testing.T) {
	env := envs.NewBuilder().WithDateFormat(envs.DateFormatDayMonthYear).WithDefaultCountry("US").Build()
	resolver := contactql.NewMockResolver(
		[]assets.Field{
			static.NewField("f1b5aea6-6586-41c7-9020-1a6326cc6565", "age", "Age", assets.FieldTypeNumber),
		},
		[]assets.Flow{},
		[]assets.Group{},
	)

	tests := []struct {
		text   string
		parsed string
	}{
		{
			text:   `age > 10 and age < 20 and age < 40 and age < 60`,
			parsed: `age > 10 AND age < 20 AND age < 40 AND age < 60`,
		},
		{
			text:   `age > 10 and age < 20 and age < 40 or age < 60`,
			parsed: `(age > 10 AND age < 20 AND age < 40) OR age < 60`,
		},
		{
			text:   `age > 10 or age < 20 and age < 40 and age < 60`,
			parsed: `age > 10 OR (age < 20 AND age < 40 AND age < 60)`,
		},
		{
			text:   `age > 10 and age < 20 or age > 30 and age < 40 or age > 50 and age < 60`,
			parsed: `(age > 10 AND age < 20) OR (age > 30 AND age < 40) OR (age > 50 AND age < 60)`,
		},
		{
			text:   `age > 10 and age < 20 or age > 30 and age < 40 or age > 50 and age < 60 or age > 70 and age < 80`,
			parsed: `(age > 10 AND age < 20) OR (age > 30 AND age < 40) OR (age > 50 AND age < 60) OR (age > 70 AND age < 80)`,
		},
		{
			text:   `Jim McJim or Bob McBob or Ann McAnn`,
			parsed: `(name ~ "Jim" AND name ~ "McJim") OR (name ~ "Bob" AND name ~ "McBob") OR (name ~ "Ann" AND name ~ "McAnn")`,
		},
	}

	for _, tc := range tests {
		parsed, err := contactql.ParseQuery(env, tc.text, resolver)
		assert.NoError(t, err)
		assert.Equal(t, tc.parsed, parsed.String(), "parsed mismatch for input '%s'", tc.text)
	}
}

func TestQueryBuilding(t *testing.T) {
	tests := []struct {
		node  contactql.QueryNode
		query string
	}{
		{
			node:  contactql.NewCondition("age", contactql.PropertyTypeField, ">", "10"),
			query: "age > 10",
		},
		{
			node: contactql.NewBoolCombination(contactql.BoolOperatorAnd,
				contactql.NewCondition("age", contactql.PropertyTypeField, ">", "10"),
				contactql.NewCondition("age", contactql.PropertyTypeField, "<", "20"),
			),
			query: "age > 10 AND age < 20",
		},
		{
			node: contactql.NewBoolCombination(contactql.BoolOperatorOr,
				contactql.NewCondition("name", contactql.PropertyTypeField, "=", "bob"),
				contactql.NewBoolCombination(contactql.BoolOperatorAnd,
					contactql.NewCondition("age", contactql.PropertyTypeField, ">", "10"),
					contactql.NewCondition("age", contactql.PropertyTypeField, "<", "20"),
				),
			),
			query: `name = "bob" OR (age > 10 AND age < 20)`,
		},
		{
			node: contactql.NewBoolCombination(contactql.BoolOperatorAnd,
				contactql.NewCondition("age", contactql.PropertyTypeField, ">", "10"),
			),
			query: "age > 10",
		},
		{
			node:  contactql.NewBoolCombination(contactql.BoolOperatorAnd),
			query: "",
		},
		{
			node: contactql.NewBoolCombination(contactql.BoolOperatorAnd,
				contactql.NewCondition("name", contactql.PropertyTypeField, "=", "bob"),
				contactql.NewBoolCombination(contactql.BoolOperatorAnd,
					contactql.NewBoolCombination(contactql.BoolOperatorAnd),
				),
			),
			query: `name = "bob"`,
		},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.query, contactql.Stringify(tc.node.Simplify()))
	}
}
