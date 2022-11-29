package issues

import (
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
)

func init() {
	registerType(TypeLegacyVars, LegacyVarsCheck)
}

// TypeLegacyVars is our type for this issue
const TypeLegacyVars string = "legacy_vars"

// LegacyVars is a use of legacy vars issue
type LegacyVars struct {
	baseIssue

	Vars []string `json:"vars"`
}

func newLegacyVars(nodeUUID flows.NodeUUID, actionUUID flows.ActionUUID, language envs.Language, vars []string) *LegacyVars {
	return &LegacyVars{
		baseIssue: newBaseIssue(
			TypeLegacyVars,
			nodeUUID,
			actionUUID,
			language,
			"use of expressions instead of contact query",
		),
		Vars: vars,
	}
}

// LegacyVarsCheck checks for this issue
func LegacyVarsCheck(sa flows.SessionAssets, flow flows.Flow, tpls []flows.ExtractedTemplate, refs []flows.ExtractedReference, report func(flows.Issue)) {
	// look for start_session actions using legacy vars
	for _, node := range flow.Nodes() {
		for _, a := range node.Actions() {
			if a.Type() == actions.TypeStartSession {
				action := a.(*actions.StartSessionAction)
				if len(action.LegacyVars) > 0 {
					report(newLegacyVars(node.UUID(), a.UUID(), envs.NilLanguage, action.LegacyVars))
				}
			}
		}
	}
}
