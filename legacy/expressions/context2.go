package expressions

import (
	"regexp"
	"strings"

	"github.com/nyaruka/gocommon/urns"
)

type mapping struct {
	pattern *regexp.Regexp
	replace string
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
		{re(`^(?:(?:flow|child|step)\.)?(parent\.)?contact$`), `${1}contact`},
		{re(`^(?:(?:flow|child|step)\.)?(parent\.)?contact\.uuid$`), `${1}contact.uuid`},
		{re(`^(?:(?:flow|child|step)\.)?(parent\.)?contact\.id$`), `${1}contact.id`},
		{re(`^(?:(?:flow|child|step)\.)?(parent\.)?contact\.name$`), `${1}contact.name`},
		{re(`^(?:(?:flow|child|step)\.)?(parent\.)?contact\.first_name$`), `${1}contact.first_name`},
		{re(`^(?:(?:flow|child|step)\.)?(parent\.)?contact\.created_on$`), `${1}contact.created_on`},
		{re(`^(?:(?:flow|child|step)\.)?(parent\.)?contact\.language$`), `${1}contact.language`},
		{re(`^(?:(?:flow|child|step)\.)?(parent\.)?contact\.groups$`), `join(${1}contact.groups, ",")`},
		{re(`^(?:(?:flow|child|step)\.)?(parent\.)?contact\.tel_e164$`), `urn_parts(${1}urns.tel).path`},
		{re(`^(?:(?:flow|child|step)\.)?(parent\.)?contact\.(` + schemesRe + `)(\.display)?$`), `format_urn(${1}urns.$2)`},
		{re(`^(?:(?:flow|child|step)\.)?(parent\.)?contact\.(` + schemesRe + `)\.path$`), `urn_parts(${1}urns.$2).path`},
		{re(`^(?:(?:flow|child|step)\.)?(parent\.)?contact\.(` + schemesRe + `)\.scheme$`), `urn_parts(${1}urns.$2).scheme`},
		{re(`^(?:(?:flow|child|step)\.)?(parent\.)?contact\.(` + schemesRe + `)\.urn$`), `${1}urns.$2`},
		{re(`^(?:(?:flow|child|step)\.)?(parent\.)?contact\.(\w+)$`), `${1}fields.$2`},

		{re(`^flow$`), `results`},
		{re(`^flow\.(\d\w*)$`), `results["$1"]`},
		{re(`^flow\.(\d\w*)\.value$`), `results["$1"].value`},
		{re(`^flow\.(\d\w*)\.category$`), `results["$1"].category_localized`},
		{re(`^flow\.(\d\w*)\.text$`), `results["$1"].input`},
		{re(`^flow\.(\d\w*)\.time$`), `results["$1"].created_on`},
		{re(`^flow\.(\w+)$`), `results.$1`},
		{re(`^flow\.(\w+)\.value$`), `results.$1.value`},
		{re(`^flow\.(\w+)\.category$`), `results.$1.category_localized`},
		{re(`^flow\.(\w+)\.text$`), `results.$1.input`},
		{re(`^flow\.(\w+)\.time$`), `results.$1.created_on`},

		{re(`^child\.(\d\w*)\.value$`), `child.results["$1"].value`},
		{re(`^child\.(\d\w*)\.category$`), `child.results["$1"].category_localized`},
		{re(`^child\.(\d\w*)\.text$`), `child.results["$1"].input`},
		{re(`^child\.(\d\w*)\.time$`), `child.results["$1"].created_on`},
		{re(`^child\.(\w+)$`), `child.results.$1`},
		{re(`^child\.(\w+)\.value$`), `child.results.$1.value`},
		{re(`^child\.(\w+)\.category$`), `child.results.$1.category_localized`},
		{re(`^child\.(\w+)\.text$`), `child.results.$1.input`},
		{re(`^child\.(\w+)\.time$`), `child.results.$1.created_on`},

		{re(`^(?:parent|extra\.flow)\.(\d\w*)\.value$`), `parent.results["$1"].value`},
		{re(`^(?:parent|extra\.flow)\.(\d\w*)\.category$`), `parent.results["$1"].category_localized`},
		{re(`^(?:parent|extra\.flow)\.(\d\w*)\.text$`), `parent.results["$1"].input`},
		{re(`^(?:parent|extra\.flow)\.(\d\w*)\.time$`), `parent.results["$1"].created_on`},
		{re(`^(?:parent|extra\.flow)\.(\w+)$`), `parent.results.$1`},
		{re(`^(?:parent|extra\.flow)\.(\w+)\.value$`), `parent.results.$1.value`},
		{re(`^(?:parent|extra\.flow)\.(\w+)\.category$`), `parent.results.$1.category_localized`},
		{re(`^(?:parent|extra\.flow)\.(\w+)\.text$`), `parent.results.$1.input`},
		{re(`^(?:parent|extra\.flow)\.(\w+)\.time$`), `parent.results.$1.created_on`},

		{re(`^step(\.value)?$`), `format_input(input)`},
		{re(`^step\.text$`), `input.text`},
		{re(`^step\.time$`), `input.created_on`},
		{re(`^step\.attachments$`), `foreach(foreach(input.attachments, attachment_parts), extract, "url")`},
		{re(`^step\.attachments\.(\d+)$`), `attachment_parts(input.attachments[$1]).url`},

		{re(`^channel$`), `contact.channel.address`},
		{re(`^channel\.(tel|tel_e164)$`), `contact.channel.address`},
		{re(`^channel\.name$`), `contact.channel.name`},

		{re(`^date(\.now)?$`), `now()`},
		{re(`^date\.today$`), `format_date(today())`},
		{re(`^date\.tomorrow$`), `format_date(datetime_add(now(), 1, "D"))`},
		{re(`^date\.yesterday$`), `format_date(datetime_add(now(), -1, "D"))`},

		{re(`^extra$`), `legacy_extra`},
		{re(`^extra\.([\w\.]+)$`), `legacy_extra.${1}`},
	}
}

// MigrateContextReference migrates a context reference in a legacy expression
func MigrateContextReference(path string) (string, bool) {
	path = strings.ToLower(path)

	for _, mapping := range mappings {
		if mapping.pattern.MatchString(path) {
			//fmt.Printf("context ref '%s' matched '%s'\n", path, mapping.pattern)

			return mapping.pattern.ReplaceAllString(path, mapping.replace), true
		}
	}

	return path, false
}
