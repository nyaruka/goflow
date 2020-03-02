package issues

import (
	"strings"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/tools"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeLegacyExtra, LegacyExtraCheck)
}

// TypeLegacyExtra is our type for a use of @legacy_extra
const TypeLegacyExtra string = "legacy_extra"

// LegacyExtra is a legacy extra issue
type LegacyExtra struct {
	baseIssue
}

func newLegacyExtra(nodeUUID flows.NodeUUID, actionUUID flows.ActionUUID, language envs.Language) *LegacyExtra {
	return &LegacyExtra{
		baseIssue: newBaseIssue(
			TypeLegacyExtra,
			nodeUUID,
			actionUUID,
			language,
			"use of @legacy_extra in an expression",
		),
	}
}

// LegacyExtraCheck checks for legacy extra usage
func LegacyExtraCheck(sa flows.SessionAssets, flow flows.Flow, tpls []flows.ExtractedTemplate, refs []flows.ExtractedReference, report func(flows.Issue)) {
	for _, t := range tpls {
		usesLegacyExtra := false

		tools.FindContextRefsInTemplate(t.Template, flows.RunContextTopLevels, func(path []string) {
			if strings.ToLower(path[0]) == "legacy_extra" {
				usesLegacyExtra = true
			}
		})

		if usesLegacyExtra {
			var actionUUID flows.ActionUUID
			if t.Action != nil {
				actionUUID = t.Action.UUID()
			}
			report(newLegacyExtra(t.Node.UUID(), actionUUID, t.Language))
		}
	}
}
