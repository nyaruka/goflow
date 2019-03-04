package definition

import (
	"github.com/nyaruka/goflow/utils"
)

// holds information about what results a flow can generate, as a map of result
// keys to slices of result names, e.g.
//
//  { "age": ["Age"], "response_1": ["Response 1", "Response-1"] }
//
type resultsInfo map[string][]string

func newResultsInfo(names []string) resultsInfo {
	r := make(resultsInfo)
	for _, name := range names {
		key := utils.Snakify(name)
		r[key] = append(r[key], name)
	}
	return r
}
