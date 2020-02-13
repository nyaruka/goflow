package problems

import "github.com/nyaruka/goflow/flows"

func init() {
	registerType(TypeMissingDependency, MissingDependencyCheck)
}

// TypeMissingDependency is our type for a missing dependency problem
const TypeMissingDependency string = "missing_dependency"

func MissingDependencyCheck(flow flows.Flow, report func(flows.Problem)) {
	// TODO
}
