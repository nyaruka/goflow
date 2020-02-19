package issues

import (
	"sort"

	"github.com/nyaruka/goflow/envs"
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
	NodeUUID_    flows.NodeUUID   `json:"node_uuid"`
	ActionUUID_  flows.ActionUUID `json:"action_uuid,omitempty"`
	Language_    envs.Language    `json:"language,omitempty"`
	Description_ string           `json:"description"`
}

// creates a new base issue
func newBaseIssue(typeName string, nodeUUID flows.NodeUUID, actionUUID flows.ActionUUID, language envs.Language, description string) baseIssue {
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
func (p *baseIssue) NodeUUID() flows.NodeUUID { return p.NodeUUID_ }

// ActionUUID returns the UUID of the action where issue is found
func (p *baseIssue) ActionUUID() flows.ActionUUID { return p.ActionUUID_ }

// Language returns the translation language if the issue was found in a translation
func (p *baseIssue) Language() envs.Language { return p.Language_ }

// Description returns the description of the issue
func (p *baseIssue) Description() string { return p.Description_ }

// Check returns all issues in the given flow
func Check(sa flows.SessionAssets, flow flows.Flow, tpls []flows.ExtractedTemplate, refs []flows.ExtractedReference) []flows.Issue {
	issues := make([]flows.Issue, 0)
	report := func(i flows.Issue) {
		issues = append(issues, i)
	}

	for _, fn := range RegisteredTypes {
		fn(sa, flow, tpls, refs, report)
	}

	// sort issues by node order
	nodeOrder := make(map[flows.NodeUUID]int, len(flow.Nodes()))
	for i, node := range flow.Nodes() {
		nodeOrder[node.UUID()] = i
	}
	sort.SliceStable(issues, func(i, j int) bool {
		return nodeOrder[issues[i].NodeUUID()] < nodeOrder[issues[j].NodeUUID()]
	})

	return issues
}
