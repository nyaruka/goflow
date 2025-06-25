package issues

import (
	"fmt"
	"regexp"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/inspect"
	"github.com/nyaruka/goflow/flows/routers"
)

func init() {
	registerType(TypeInvalidRegex, InvalidRegexCheck)
}

// TypeInvalidRegex is our type for an invalid regex
const TypeInvalidRegex string = "invalid_regex"

// InvalidRegex is an invalid regex issue
type InvalidRegex struct {
	baseIssue

	Regex string `json:"regex"`
}

func newInvalidRegex(nodeUUID flows.NodeUUID, actionUUID flows.ActionUUID, language i18n.Language, regex string) *InvalidRegex {
	return &InvalidRegex{
		baseIssue: newBaseIssue(
			TypeInvalidRegex,
			nodeUUID,
			actionUUID,
			language,
			fmt.Sprintf("invalid regex: %s", regex),
		),
		Regex: regex,
	}
}

// InvalidRegexCheck checks for invalid regexes
func InvalidRegexCheck(sa flows.SessionAssets, flow flows.Flow, tpls []flows.ExtractedTemplate, refs []flows.ExtractedReference, report func(flows.Issue)) {
	checkTemplate := func(n flows.Node, a flows.Action, l i18n.Language, t string) {
		// only check if template doesn't contain expressions
		if !excellent.HasExpressions(t, flows.RunContextTopLevels) {
			_, err := regexp.Compile(t)
			if err != nil {
				var actionUUID flows.ActionUUID
				if a != nil {
					actionUUID = a.UUID()
				}
				report(newInvalidRegex(n.UUID(), actionUUID, l, t))
			}
		}
	}

	// look for switch router cases which are regex tests
	for _, node := range flow.Nodes() {
		if node.Router() != nil && node.Router().Type() == routers.TypeSwitch {
			router := node.Router().(*routers.Switch)
			for _, kase := range router.Cases() {
				if kase.Type == "has_pattern" && len(kase.Arguments) > 0 {
					checkTemplate(node, nil, "", kase.Arguments[0])

					inspect.Translations(flow.Localization(), kase.LocalizationUUID(), "arguments", func(l i18n.Language, t string) {
						checkTemplate(node, nil, l, t)
					})
				}
			}
		}
	}
}
