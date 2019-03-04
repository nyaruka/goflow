package definition

import (
	"github.com/nyaruka/goflow/utils"
)

type resultInfo struct {
	Names []string `json:"names"`
}

// holds information about what results a flow can generate, as a map of result
// keys to slices of result names, e.g.
//
//  { "age": {"names": ["Age"]}, "response_1": {"names": ["Response 1", "Response-1"]} }
//
type resultsInfo map[string]*resultInfo

func newResultsInfo(names []string) resultsInfo {
	r := make(resultsInfo)
	for _, name := range names {
		key := utils.Snakify(name)
		info, exists := r[key]
		if !exists {
			info = &resultInfo{}
			r[key] = info
		}

		info.Names = append(info.Names, name)
	}
	return r
}
