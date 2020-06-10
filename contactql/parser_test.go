package contactql_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static/types"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"

	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		text   string
		parsed string
		err    string
		redact envs.RedactionPolicy
	}{
		// implicit conditions
		{`"will"`, `name ~ "will"`, "", envs.RedactionPolicyNone},
		{`Will`, `name ~ "Will"`, "", envs.RedactionPolicyNone},
		{`wil`, `name ~ "wil"`, "", envs.RedactionPolicyNone},
		{`wi`, `name ~ "wi"`, "", envs.RedactionPolicyNone},
		{`w`, `name = "w"`, "", envs.RedactionPolicyNone}, // don't have at least 1 token of >= 2 chars
		{`w me`, `name = "w" AND name ~ "me"`, "", envs.RedactionPolicyNone},
		{`w m`, `name = "w" AND name = "m"`, "", envs.RedactionPolicyNone},
		{`tel:+0123456566`, `tel = "+0123456566"`, "", envs.RedactionPolicyNone}, // whole query is a URN
		{`twitter:bobby`, `twitter = "bobby"`, "", envs.RedactionPolicyNone},
		{`(202) 456-1111`, `tel = "+12024561111"`, "", envs.RedactionPolicyNone}, // whole query looks like a phone number
		{`+12024561111`, `tel = "+12024561111"`, "", envs.RedactionPolicyNone},
		{` 202.456.1111 `, `tel = "+12024561111"`, "", envs.RedactionPolicyNone},
		{`"+12024561111"`, `tel ~ "+12024561111"`, "", envs.RedactionPolicyNone},
		{`566`, `name ~ 566`, "", envs.RedactionPolicyNone}, // too short to be a phone number

		// implicit conditions with URN redaction
		{`will`, `name ~ "will"`, "", envs.RedactionPolicyURNs},
		{`tel:+0123456566`, `name ~ "tel:+0123456566"`, "", envs.RedactionPolicyURNs},
		{`twitter:bobby`, `name ~ "twitter:bobby"`, "", envs.RedactionPolicyURNs},
		{`0123456566`, `id = 123456566`, "", envs.RedactionPolicyURNs},
		{`+0123456566`, `id = 123456566`, "", envs.RedactionPolicyURNs},
		{`0123-456-566`, `name ~ "0123-456-566"`, "", envs.RedactionPolicyURNs},

		// explicit conditions on name
		{`Name=will`, `name = "will"`, "", envs.RedactionPolicyNone},
		{`Name ~ "felix"`, `name ~ "felix"`, "", envs.RedactionPolicyNone},
		{`Name HAS "Felix"`, `name ~ "Felix"`, "", envs.RedactionPolicyNone},
		{`name is ""`, `name = ""`, "", envs.RedactionPolicyNone},            // is not set
		{`name != ""`, `name != ""`, "", envs.RedactionPolicyNone},           // is set
		{`name != "felix"`, `name != "felix"`, "", envs.RedactionPolicyNone}, // is not equal to value
		{`Name ~ ""`, ``, "value must contain a word of at least 2 characters long for a contains condition on name", envs.RedactionPolicyNone},

		// explicit attribute conditions
		{`language = spa`, `language = "spa"`, "", envs.RedactionPolicyNone},
		{`Group IS U-Reporters`, `group = "U-Reporters"`, "", envs.RedactionPolicyNone},
		{`CREATED_ON>27-01-2020`, `created_on > "27-01-2020"`, "", envs.RedactionPolicyNone},

		// explicit conditions on URN
		{`tel=""`, `tel = ""`, "", envs.RedactionPolicyNone},
		{`tel!=""`, `tel != ""`, "", envs.RedactionPolicyNone},
		{`tel IS 233`, `tel = 233`, "", envs.RedactionPolicyNone},
		{`tel HAS 233`, `tel ~ 233`, "", envs.RedactionPolicyNone},
		{`tel ~ 23`, ``, "value must be least 3 characters long for a contains condition on a URN", envs.RedactionPolicyNone},
		{`mailto = user@example.com`, `mailto = "user@example.com"`, "", envs.RedactionPolicyNone},
		{`MAILTO ~ user@example.com`, `mailto ~ "user@example.com"`, "", envs.RedactionPolicyNone},
		{`URN=ewok`, `urn = "ewok"`, "", envs.RedactionPolicyNone},

		// explicit conditions on URN with URN redaction
		{`tel=""`, `tel = ""`, "", envs.RedactionPolicyURNs},
		{`tel!=""`, `tel != ""`, "", envs.RedactionPolicyURNs},
		{`mailto=""`, `mailto = ""`, "", envs.RedactionPolicyURNs},
		{`mailto!=""`, `mailto != ""`, "", envs.RedactionPolicyURNs},
		{`urn=""`, `urn = ""`, "", envs.RedactionPolicyURNs},
		{`urn!=""`, `urn != ""`, "", envs.RedactionPolicyURNs},
		{`tel = 233`, ``, "cannot query on redacted URNs", envs.RedactionPolicyURNs},
		{`tel ~ 233`, ``, "cannot query on redacted URNs", envs.RedactionPolicyURNs},
		{`mailto = user@example.com`, ``, "cannot query on redacted URNs", envs.RedactionPolicyURNs},
		{`MAILTO ~ user@example.com`, ``, "cannot query on redacted URNs", envs.RedactionPolicyURNs},
		{`URN=ewok`, ``, "cannot query on redacted URNs", envs.RedactionPolicyURNs},

		// field conditions
		{`Age IS 18`, `age = 18`, "", envs.RedactionPolicyNone},
		{`AGE != ""`, `age != ""`, "", envs.RedactionPolicyNone},
		{`age ~ 34`, ``, "contains conditions can only be used with name or URN values", envs.RedactionPolicyNone},
		{`gender ~ M`, ``, "contains conditions can only be used with name or URN values", envs.RedactionPolicyNone},

		// lt/lte/gt/gte comparisons
		{`Age > "18"`, `age > 18`, "", envs.RedactionPolicyNone},
		{`Age >= 18`, `age >= 18`, "", envs.RedactionPolicyNone},
		{`Age < 18`, `age < 18`, "", envs.RedactionPolicyNone},
		{`Age <= 18`, `age <= 18`, "", envs.RedactionPolicyNone},
		{`DOB > "27-01-2020"`, `dob > "27-01-2020"`, "", envs.RedactionPolicyNone},
		{`DOB >= 27-01-2020`, `dob >= "27-01-2020"`, "", envs.RedactionPolicyNone},
		{`DOB < 27/01/2020`, `dob < "27/01/2020"`, "", envs.RedactionPolicyNone},
		{`DOB <= 27.01.2020`, `dob <= "27.01.2020"`, "", envs.RedactionPolicyNone},
		{`name > Will`, ``, "comparisons with > can only be used with date and number fields", envs.RedactionPolicyNone},
		{`tel < 23425`, ``, "comparisons with < can only be used with date and number fields", envs.RedactionPolicyNone},

		// implicit combinations
		{`will felix`, `name ~ "will" AND name ~ "felix"`, "", envs.RedactionPolicyNone},
		{`will +123456566`, `name ~ "will" AND tel ~ "+123456566"`, "", envs.RedactionPolicyNone},

		// explicit combinations...
		{`will and felix`, `name ~ "will" AND name ~ "felix"`, "", envs.RedactionPolicyNone}, // explicit AND
		{`will or felix or matt`, `(name ~ "will" OR name ~ "felix") OR name ~ "matt"`, "", envs.RedactionPolicyNone},
		{`name=will or Name ~ "felix"`, `name = "will" OR name ~ "felix"`, "", envs.RedactionPolicyNone},
		{`Name is will or Name has felix`, `name = "will" OR name ~ "felix"`, "", envs.RedactionPolicyNone}, // comparator aliases
		{`will or Name ~ "felix"`, `name ~ "will" OR name ~ "felix"`, "", envs.RedactionPolicyNone},

		// boolean operator precedence is AND before OR, even when AND is implicit
		{`will and felix or matt amber`, `(name ~ "will" AND name ~ "felix") OR (name ~ "matt" AND name ~ "amber")`, "", envs.RedactionPolicyNone},

		// boolean combinations can themselves be combined
		{
			`(Age < 18 and Gender = "male") or (Age > 18 and Gender = "female")`,
			`(age < 18 AND gender = "male") OR (age > 18 AND gender = "female")`,
			"",
			envs.RedactionPolicyNone,
		},

		{`xyz != ""`, ``, "can't resolve 'xyz' to attribute, scheme or field", envs.RedactionPolicyNone},
		{`group != "Gamers"`, ``, "'Gamers' is not a valid group name", envs.RedactionPolicyNone},
		{`language = "xxxx"`, ``, "'xxxx' is not a valid language code", envs.RedactionPolicyNone},

		{`name = "O\"Leary"`, `name = "O\"Leary"`, "", envs.RedactionPolicyNone}, // string unquoting

		// = supported for everything
		{`uuid = f81d1eb5-215d-4ae8-90fa-38b3f2d6e328`, `uuid = "f81d1eb5-215d-4ae8-90fa-38b3f2d6e328"`, "", envs.RedactionPolicyNone},
		{`id = 02352`, `id = 02352`, "", envs.RedactionPolicyNone},
		{`name = felix`, `name = "felix"`, "", envs.RedactionPolicyNone},
		{`language = eng`, `language = "eng"`, "", envs.RedactionPolicyNone},
		{`group = u-reporters`, `group = "U-Reporters"`, "", envs.RedactionPolicyNone},
		{`created_on = 20-02-2020`, `created_on = "20-02-2020"`, "", envs.RedactionPolicyNone},
		{`tel = 02352`, `tel = 02352`, "", envs.RedactionPolicyNone},
		{`urn = 02352`, `urn = 02352`, "", envs.RedactionPolicyNone},
		{`age = 18`, `age = 18`, "", envs.RedactionPolicyNone},
		{`gender = male`, `gender = "male"`, "", envs.RedactionPolicyNone},
		{`dob = 20-02-2020`, `dob = "20-02-2020"`, "", envs.RedactionPolicyNone},
		{`state = Pichincha`, `state = "Pichincha"`, "", envs.RedactionPolicyNone},

		// != supported for everything
		{`uuid != f81d1eb5-215d-4ae8-90fa-38b3f2d6e328`, `uuid != "f81d1eb5-215d-4ae8-90fa-38b3f2d6e328"`, "", envs.RedactionPolicyNone},
		{`id != 02352`, `id != 02352`, "", envs.RedactionPolicyNone},
		{`name != felix`, `name != "felix"`, "", envs.RedactionPolicyNone},
		{`language != eng`, `language != "eng"`, "", envs.RedactionPolicyNone},
		{`group != u-reporters`, `group != "U-Reporters"`, "", envs.RedactionPolicyNone},
		{`created_on != 20-02-2020`, `created_on != "20-02-2020"`, "", envs.RedactionPolicyNone},
		{`tel != 02352`, `tel != 02352`, "", envs.RedactionPolicyNone},
		{`urn != 02352`, `urn != 02352`, "", envs.RedactionPolicyNone},
		{`age != 18`, `age != 18`, "", envs.RedactionPolicyNone},
		{`gender != male`, `gender != "male"`, "", envs.RedactionPolicyNone},
		{`dob != 20-02-2020`, `dob != "20-02-2020"`, "", envs.RedactionPolicyNone},
		{`state != Pichincha`, `state != "Pichincha"`, "", envs.RedactionPolicyNone},

		// = "" supported for name, language, fields and urns
		{`uuid = ""`, ``, "can't check whether 'uuid' is set or not set", envs.RedactionPolicyNone},
		{`id = ""`, ``, "can't check whether 'id' is set or not set", envs.RedactionPolicyNone},
		{`name = ""`, `name = ""`, "", envs.RedactionPolicyNone},
		{`language = ""`, `language = ""`, "", envs.RedactionPolicyNone},
		{`group = ""`, ``, "can't check whether 'group' is set or not set", envs.RedactionPolicyNone},
		{`created_on = ""`, ``, "can't check whether 'created_on' is set or not set", envs.RedactionPolicyNone},
		{`tel = ""`, `tel = ""`, "", envs.RedactionPolicyNone},
		{`urn = ""`, `urn = ""`, "", envs.RedactionPolicyNone},
		{`age = ""`, `age = ""`, "", envs.RedactionPolicyNone},
		{`gender = ""`, `gender = ""`, "", envs.RedactionPolicyNone},
		{`dob = ""`, `dob = ""`, "", envs.RedactionPolicyNone},
		{`state = ""`, `state = ""`, "", envs.RedactionPolicyNone},

		// ~ only supported for name and URNs
		{`uuid ~ 02352`, ``, "contains conditions can only be used with name or URN values", envs.RedactionPolicyNone},
		{`id ~ 02352`, ``, "contains conditions can only be used with name or URN values", envs.RedactionPolicyNone},
		{`name ~ felix`, `name ~ "felix"`, "", envs.RedactionPolicyNone},
		{`language ~ eng`, ``, "contains conditions can only be used with name or URN values", envs.RedactionPolicyNone},
		{`group ~ porters`, ``, "contains conditions can only be used with name or URN values", envs.RedactionPolicyNone},
		{`created_on ~ 2018`, ``, "contains conditions can only be used with name or URN values", envs.RedactionPolicyNone},
		{`tel ~ 02352`, `tel ~ 02352`, "", envs.RedactionPolicyNone},
		{`urn ~ 02352`, `urn ~ 02352`, "", envs.RedactionPolicyNone},
		{`age ~ 18`, ``, "contains conditions can only be used with name or URN values", envs.RedactionPolicyNone},
		{`gender ~ mal`, ``, "contains conditions can only be used with name or URN values", envs.RedactionPolicyNone},
		{`dob ~ 20-02-2020`, ``, "contains conditions can only be used with name or URN values", envs.RedactionPolicyNone},
		{`state ~ Pichincha`, ``, "contains conditions can only be used with name or URN values", envs.RedactionPolicyNone},

		// > >= < <= only supported for numeric or date fields
		{`uuid > 02352`, ``, "comparisons with > can only be used with date and number fields", envs.RedactionPolicyNone},
		{`id > 02352`, ``, "comparisons with > can only be used with date and number fields", envs.RedactionPolicyNone},
		{`name > felix`, ``, "comparisons with > can only be used with date and number fields", envs.RedactionPolicyNone},
		{`language > eng`, ``, "comparisons with > can only be used with date and number fields", envs.RedactionPolicyNone},
		{`group > reporters`, ``, "comparisons with > can only be used with date and number fields", envs.RedactionPolicyNone},
		{`created_on > 20-02-2020`, `created_on > "20-02-2020"`, "", envs.RedactionPolicyNone},
		{`tel > 02352`, ``, "comparisons with > can only be used with date and number fields", envs.RedactionPolicyNone},
		{`urn > 02352`, ``, "comparisons with > can only be used with date and number fields", envs.RedactionPolicyNone},
		{`age > 18`, `age > 18`, "", envs.RedactionPolicyNone},
		{`gender > male`, ``, "comparisons with > can only be used with date and number fields", envs.RedactionPolicyNone},
		{`dob > 20-02-2020`, `dob > "20-02-2020"`, "", envs.RedactionPolicyNone},
		{`state > Pichincha`, ``, "comparisons with > can only be used with date and number fields", envs.RedactionPolicyNone},
	}

	resolver := contactql.NewMockResolver(map[string]assets.Field{
		"age":    types.NewField(assets.FieldUUID("f1b5aea6-6586-41c7-9020-1a6326cc6565"), "age", "Age", assets.FieldTypeNumber),
		"gender": types.NewField(assets.FieldUUID("d66a7823-eada-40e5-9a3a-57239d4690bf"), "gender", "Gender", assets.FieldTypeText),
		"state":  types.NewField(assets.FieldUUID("165def68-3216-4ebf-96bc-f6f1ee5bd966"), "state", "State", assets.FieldTypeState),
		"dob":    types.NewField(assets.FieldUUID("85baf5e1-b57a-46dc-a726-a84e8c4229c7"), "dob", "DOB", assets.FieldTypeDatetime),
	}, map[string]assets.Group{
		"u-reporters": types.NewGroup(assets.GroupUUID(""), "U-Reporters", ""),
	})

	for _, tc := range tests {
		parsed, err := contactql.ParseQuery(tc.text, tc.redact, "US", resolver)
		if tc.err != "" {
			assert.EqualError(t, err, tc.err, "error mismatch for '%s'", tc.text)
			assert.Nil(t, parsed)
		} else {
			assert.NoError(t, err, "unexpected error for '%s'", tc.text)
			assert.Equal(t, tc.parsed, parsed.String(), "parse mismatch for '%s'", tc.text)
		}
	}
}

func TestParsingErrors(t *testing.T) {
	_, err := contactql.ParseQuery("name = ", envs.RedactionPolicyNone, "US", nil)
	assert.EqualError(t, err, "mismatched input '<EOF>' expecting {TEXT, STRING}")
}
