package problems

import "github.com/nyaruka/goflow/flows"

type reportFunc func(flows.Flow, func(flows.Problem))

var registeredTypes = map[string]reportFunc{}

// registers a new type of problem
func registerType(name string, report reportFunc) {
	registeredTypes[name] = report
}

// base of all problem types
type baseProblem struct {
	Type_       string           `json:"type"`
	NodeUUID_   flows.NodeUUID   `json:"node_uuid"`
	ActionUUID_ flows.ActionUUID `json:"action_uuid"`
}

// Type returns the type of this problem
func (p *baseProblem) Type() string { return p.Type_ }

// NodeUUID returns the UUID of the node where problem is found
func (p *baseProblem) NodeUUID() flows.NodeUUID { return p.NodeUUID_ }

// ActionUUID returns the UUID of the action where problem is found
func (p *baseProblem) ActionUUID() flows.ActionUUID { return p.ActionUUID_ }

func Check(flow flows.Flow) []flows.Problem {
	problems := make([]flows.Problem, 0)
	report := func(p flows.Problem) {
		problems = append(problems, p)
	}

	for _, fn := range registeredTypes {
		fn(flow, report)
	}

	return problems
}
