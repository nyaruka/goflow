package expressions

import (
	"regexp"
	"strings"

	"github.com/nyaruka/gocommon/urns"
)

type mapping struct {
	pattern *regexp.Regexp
	replace string
	isDate  bool // is a date (not an ISO datetime) that may have to be formatted
}

var mappings []mapping

func init() {
	schemes := make([]string, 0, len(urns.ValidSchemes))
	for s := range urns.ValidSchemes {
		schemes = append(schemes, s)
	}

	schemesRe := strings.Join(schemes, `|`)
	var re = regexp.MustCompile

	mappings = []mapping{
		{re(`^(?:(?:flow|step)\.)?((?:parent|child)\.)?contact$`), `${1}contact`, false},
		{re(`^(?:(?:flow|step)\.)?((?:parent|child)\.)?contact\.uuid$`), `${1}contact.uuid`, false},
		{re(`^(?:(?:flow|step)\.)?((?:parent|child)\.)?contact\.id$`), `${1}contact.id`, false},
		{re(`^(?:(?:flow|step)\.)?((?:parent|child)\.)?contact\.name$`), `${1}contact.name`, false},
		{re(`^(?:(?:flow|step)\.)?((?:parent|child)\.)?contact\.first_name$`), `${1}contact.first_name`, false},
		{re(`^(?:(?:flow|step)\.)?((?:parent|child)\.)?contact\.created_on$`), `${1}contact.created_on`, false},
		{re(`^(?:(?:flow|step)\.)?((?:parent|child)\.)?contact\.language$`), `${1}contact.language`, false},
		{re(`^(?:(?:flow|step)\.)?((?:parent|child)\.)?contact\.groups$`), `join(${1}contact.groups, ",")`, false},
		{re(`^(?:(?:flow|step)\.)?((?:parent|child)\.)?contact\.tel_e164$`), `default(urn_parts(${1}urns.tel).path, "")`, false},
		{re(`^(?:(?:flow|step)\.)?((?:parent|child)\.)?contact\.tel$`), `format_urn(${1}urns.tel)`, false},
		{re(`^(?:(?:flow|step)\.)?((?:parent|child)\.)?contact\.(` + schemesRe + `)$`), `default(urn_parts(${1}urns.$2).path, "")`, false},
		{re(`^(?:(?:flow|step)\.)?((?:parent|child)\.)?contact\.(` + schemesRe + `)\.display$`), `format_urn(${1}urns.$2)`, false},
		{re(`^(?:(?:flow|step)\.)?((?:parent|child)\.)?contact\.(` + schemesRe + `)\.path$`), `urn_parts(${1}urns.$2).path`, false},
		{re(`^(?:(?:flow|step)\.)?((?:parent|child)\.)?contact\.(` + schemesRe + `)\.scheme$`), `urn_parts(${1}urns.$2).scheme`, false},
		{re(`^(?:(?:flow|step)\.)?((?:parent|child)\.)?contact\.(` + schemesRe + `)\.urn$`), `${1}urns.$2`, false},
		{re(`^(?:(?:flow|step)\.)?((?:parent|child)\.)?contact\.(\w+)$`), `${1}fields.$2`, false},

		{re(`^flow$`), `results`, false},
		{re(`^flow\.(\w+)$`), `results.$1`, false},
		{re(`^flow\.(\w+)\.value$`), `results.$1.value`, false},
		{re(`^flow\.(\w+)\.category$`), `results.$1.category_localized`, false},
		{re(`^flow\.(\w+)\.text$`), `results.$1.input`, false},
		{re(`^flow\.(\w+)\.time$`), `results.$1.created_on`, false},

		{re(`^child$`), `child.results`, false},
		{re(`^child\.(\w+)$`), `child.results.$1`, false},
		{re(`^child\.(\w+)\.value$`), `child.results.$1.value`, false},
		{re(`^child\.(\w+)\.category$`), `child.results.$1.category_localized`, false},
		{re(`^child\.(\w+)\.text$`), `child.results.$1.input`, false},
		{re(`^child\.(\w+)\.time$`), `child.results.$1.created_on`, false},

		{re(`^(?:parent|extra\.flow)$`), `parent.results`, false},
		{re(`^(?:parent|extra\.flow)\.(\w+)$`), `parent.results.$1`, false},
		{re(`^(?:parent|extra\.flow)\.(\w+)\.value$`), `parent.results.$1.value`, false},
		{re(`^(?:parent|extra\.flow)\.(\w+)\.category$`), `parent.results.$1.category_localized`, false},
		{re(`^(?:parent|extra\.flow)\.(\w+)\.text$`), `parent.results.$1.input`, false},
		{re(`^(?:parent|extra\.flow)\.(\w+)\.time$`), `parent.results.$1.created_on`, false},

		{re(`^step(\.value)?$`), `input`, false},
		{re(`^step\.text$`), `input.text`, false},
		{re(`^step\.time$`), `input.created_on`, false},
		{re(`^step\.attachments$`), `foreach(foreach(input.attachments, attachment_parts), extract, "url")`, false},
		{re(`^step\.attachments\.(\d+)$`), `attachment_parts(input.attachments[$1]).url`, false},

		{re(`^channel$`), `contact.channel.address`, false},
		{re(`^channel\.(address|tel|tel_e164)$`), `contact.channel.address`, false},
		{re(`^channel\.name$`), `contact.channel.name`, false},

		{re(`^date(\.now)?$`), `now()`, false},
		{re(`^date\.today$`), `today()`, true},
		{re(`^date\.tomorrow$`), `datetime_add(now(), 1, "D")`, true},
		{re(`^date\.yesterday$`), `datetime_add(now(), -1, "D")`, true},

		{re(`^extra$`), `legacy_extra`, false},
		{re(`^extra\.([\w\.]+)$`), `legacy_extra.${1}`, false},
	}
}

// MigrateContextReference migrates a context reference in a legacy expression
func MigrateContextReference(path string, rawDates bool) string {
	path = strings.ToLower(path)

	for _, mapping := range mappings {
		if mapping.pattern.MatchString(path) {
			migrated := mapping.pattern.ReplaceAllString(path, mapping.replace)
			if mapping.isDate && !rawDates {
				migrated = wrap(migrated, "format_date")
			}
			return fixLookups(migrated)
		}
	}

	return path
}

var numericLookupRegex = regexp.MustCompile(`\.\d+\w*`)

// fixes property lookups
//  .1 => ["1"]
//  .1foo  => ["1foo"]
func fixLookups(path string) string {
	return numericLookupRegex.ReplaceAllStringFunc(path, func(lookup string) string {
		return `["` + lookup[1:] + `"]`
	})
}
