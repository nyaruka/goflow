package problems

import (
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeMissingDependency, MissingDependencyCheck)
}

// TypeMissingDependency is our type for a missing dependency problem
const TypeMissingDependency string = "missing_dependency"

// MissingDependency is a missing asset dependency
type MissingDependency struct {
	baseProblem

	Dependency assets.TypedReference `json:"dependency"`
}

func newMissingDependency(nodeUUID flows.NodeUUID, actionUUID flows.ActionUUID, ref assets.Reference) *MissingDependency {
	return &MissingDependency{
		baseProblem: newBaseProblem(
			TypeMissingDependency,
			nodeUUID,
			actionUUID,
			fmt.Sprintf("missing %s dependency '%s'", ref.Type(), ref.Identity()),
		),
		Dependency: assets.NewTypedReference(ref),
	}
}

// MissingDependencyCheck checks for missing dependencies
func MissingDependencyCheck(sa flows.SessionAssets, flow flows.Flow, refs []flows.ExtractedReference, report func(flows.Problem)) {
	for _, ref := range refs {
		if !ref.Check(sa) {
			var actionUUID flows.ActionUUID
			if ref.Action != nil {
				actionUUID = ref.Action.UUID()
			}
			report(newMissingDependency(ref.Node.UUID(), actionUUID, ref.Reference))
		}
	}
}
