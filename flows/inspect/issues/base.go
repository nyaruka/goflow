package issues

import (
	"maps"
	"slices"
	"sort"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/flows"
)

type reportFunc func(flows.SessionAssets, flows.Flow, []flows.ExtractedTemplate, []flows.ExtractedReference, func(flows.Issue))

var RegisteredTypes = map[string]reportFunc{}

// registers a new type of issue
func registerType(name string, report reportFunc) {
	RegisteredTypes[name] = report
}

// base of all issue types
type baseIssue struct {
	Type_        string           `json:"type"`
	NodeUUID_    core.NodeUUID    `json:"node_uuid"`
	ActionUUID_  flows.ActionUUID `json:"action_uuid,omitempty"`
	Language_    i18n.Language    `json:"language,omitempty"`
	Description_ string           `json:"description"`
}

// creates a new base issue
func newBaseIssue(typeName string, nodeUUID core.NodeUUID, actionUUID flows.ActionUUID, language i18n.Language, description string) baseIssue {
	return baseIssue{
		Type_:        typeName,
		NodeUUID_:    nodeUUID,
		ActionUUID_:  actionUUID,
		Language_:    language,
		Description_: description,
	}
}

// Type returns the type of this issue
func (p *baseIssue) Type() string { return p.Type_ }

// NodeUUID returns the UUID of the node where issue is found
func (p *baseIssue) NodeUUID() core.NodeUUID { return p.NodeUUID_ }

// ActionUUID returns the UUID of the action where issue is found
func (p *baseIssue) ActionUUID() flows.ActionUUID { return p.ActionUUID_ }

// Language returns the translation language if the issue was found in a translation
func (p *baseIssue) Language() i18n.Language { return p.Language_ }

// Description returns the description of the issue
func (p *baseIssue) Description() string { return p.Description_ }

// Check returns all issues in the given flow
func Check(sa flows.SessionAssets, flow flows.Flow, tpls []flows.ExtractedTemplate, refs []flows.ExtractedReference) []flows.Issue {
	issues := make([]flows.Issue, 0)
	report := func(i flows.Issue) {
		issues = append(issues, i)
	}

	// run checks in a fixed order so that issues on the same action are reported deterministically
	for _, typeName := range slices.Sorted(maps.Keys(RegisteredTypes)) {
		RegisteredTypes[typeName](sa, flow, tpls, refs, report)
	}

	// sort issues by node order, then by action order within the node (router issues come last)
	nodeOrder := make(map[core.NodeUUID]int, len(flow.Nodes()))
	actionOrder := make(map[flows.ActionUUID]int)
	for i, node := range flow.Nodes() {
		nodeOrder[node.UUID()] = i
		for j, action := range node.Actions() {
			actionOrder[action.UUID()] = j
		}
	}
	actionIndex := func(i flows.Issue) int {
		if idx, ok := actionOrder[i.ActionUUID()]; ok {
			return idx
		}
		return len(actionOrder)
	}
	sort.SliceStable(issues, func(i, j int) bool {
		ni, nj := nodeOrder[issues[i].NodeUUID()], nodeOrder[issues[j].NodeUUID()]
		if ni != nj {
			return ni < nj
		}
		return actionIndex(issues[i]) < actionIndex(issues[j])
	})

	return issues
}
