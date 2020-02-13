package problems

import "github.com/nyaruka/goflow/flows"

type reportFunc func(flows.SessionAssets, flows.Flow, []flows.ExtractedReference, func(flows.Problem))

var registeredTypes = map[string]reportFunc{}

// registers a new type of problem
func registerType(name string, report reportFunc) {
	registeredTypes[name] = report
}

// base of all problem types
type baseProblem struct {
	Type_       string           `json:"type"`
	NodeUUID_   flows.NodeUUID   `json:"node_uuid"`
	ActionUUID_ flows.ActionUUID `json:"action_uuid,omitempty"`
}

// creates a new base problem
func newBaseProblem(typeName string, nodeUUID flows.NodeUUID, actionUUID flows.ActionUUID) baseProblem {
	return baseProblem{Type_: typeName, NodeUUID_: nodeUUID, ActionUUID_: actionUUID}
}

// Type returns the type of this problem
func (p *baseProblem) Type() string { return p.Type_ }

// NodeUUID returns the UUID of the node where problem is found
func (p *baseProblem) NodeUUID() flows.NodeUUID { return p.NodeUUID_ }

// ActionUUID returns the UUID of the action where problem is found
func (p *baseProblem) ActionUUID() flows.ActionUUID { return p.ActionUUID_ }

// Check returns all problems in the given flow
func Check(sa flows.SessionAssets, flow flows.Flow, refs []flows.ExtractedReference) []flows.Problem {
	problems := make([]flows.Problem, 0)
	report := func(p flows.Problem) {
		problems = append(problems, p)
	}

	for _, fn := range registeredTypes {
		fn(sa, flow, refs, report)
	}

	return problems
}
