package features

import (
	"sort"
	"strings"

	"github.com/nyaruka/goflow/flows"
)

type checkFunc func(flows.Flow) bool

var registeredTypes = map[flows.Feature]checkFunc{}

// registers a new type of issue
func registerType(name flows.Feature, check checkFunc) {
	registeredTypes[name] = check
}

// Check returns all features in the given flow
func Check(flow flows.Flow) []flows.Feature {
	features := make([]flows.Feature, 0, len(registeredTypes))

	for f, fn := range registeredTypes {
		if fn(flow) {
			features = append(features, f)
		}
	}

	sort.Slice(features, func(i, j int) bool { return strings.Compare(string(features[i]), string(features[j])) < 0 })
	return features
}

func hasActionTypes(flow flows.Flow, types ...string) bool {
	for _, n := range flow.Nodes() {
		for _, a := range n.Actions() {
			for _, t := range types {
				if a.Type() == t {
					return true
				}
			}
		}
	}
	return false
}
