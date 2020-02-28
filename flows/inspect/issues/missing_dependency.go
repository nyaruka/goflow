package issues

import (
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/inspect"
)

func init() {
	registerType(TypeMissingDependency, MissingDependencyCheck)
}

// TypeMissingDependency is our type for a missing dependency issue
const TypeMissingDependency string = "missing_dependency"

// MissingDependency is a missing asset dependency
type MissingDependency struct {
	baseIssue

	Dependency assets.TypedReference `json:"dependency"`
}

func newMissingDependency(nodeUUID flows.NodeUUID, actionUUID flows.ActionUUID, language envs.Language, ref assets.Reference) *MissingDependency {
	return &MissingDependency{
		baseIssue: newBaseIssue(
			TypeMissingDependency,
			nodeUUID,
			actionUUID,
			language,
			fmt.Sprintf("missing %s dependency '%s'", ref.Type(), ref.Identity()),
		),
		Dependency: assets.NewTypedReference(ref),
	}
}

// MissingDependencyCheck checks for missing dependencies
func MissingDependencyCheck(sa flows.SessionAssets, flow flows.Flow, tpls []flows.ExtractedTemplate, refs []flows.ExtractedReference, report func(flows.Issue)) {
	// skip check if we don't have assets
	if sa == nil {
		return
	}

	for _, ref := range refs {
		if !inspect.CheckReference(sa, ref.Reference) {
			var actionUUID flows.ActionUUID
			if ref.Action != nil {
				actionUUID = ref.Action.UUID()
			}
			report(newMissingDependency(ref.Node.UUID(), actionUUID, ref.Language, ref.Reference))
		}
	}
}
