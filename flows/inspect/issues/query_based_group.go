package issues

import (
	"fmt"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
)

func init() {
	registerType(TypeQueryBasedGroup, QueryBasedGroupCheck)
}

// TypeQueryBasedGroup is our type for a query based group issue
const TypeQueryBasedGroup string = "query_based_group"

// QueryBasedGroup is an action which tries to add or remove contacts to or from a query based group
type QueryBasedGroup struct {
	baseIssue

	Group *assets.GroupReference `json:"group"`
}

func newQueryBasedGroup(nodeUUID core.NodeUUID, actionUUID flows.ActionUUID, language i18n.Language, description string, ref *assets.GroupReference) *QueryBasedGroup {
	return &QueryBasedGroup{
		baseIssue: newBaseIssue(TypeQueryBasedGroup, nodeUUID, actionUUID, language, description),
		Group:     ref,
	}
}

// QueryBasedGroupCheck checks for actions which add or remove contacts to or from query based groups
func QueryBasedGroupCheck(sa flows.SessionAssets, flow flows.Flow, tpls []flows.ExtractedTemplate, refs []flows.ExtractedReference, report func(flows.Issue)) {
	// skip check if we don't have assets
	if sa == nil {
		return
	}

	for _, ref := range refs {
		groupRef, isGroup := ref.Reference.(*assets.GroupReference)
		if !isGroup || ref.Action == nil || groupRef.Variable() {
			continue
		}

		var format string
		switch ref.Action.Type() {
		case actions.TypeAddContactGroups:
			format = "can't add contacts to query based group '%s'"
		case actions.TypeRemoveContactGroups:
			format = "can't remove contacts from query based group '%s'"
		default:
			continue
		}

		// missing groups are reported as missing dependencies
		group := sa.Groups().Get(groupRef.UUID)
		if group != nil && group.UsesQuery() {
			report(newQueryBasedGroup(ref.Node.UUID(), ref.Action.UUID(), ref.Language, fmt.Sprintf(format, group.Name()), groupRef))
		}
	}
}
