package definition

import (
	"github.com/nyaruka/goflow/utils"
)

// holds information about a possible result in the flow
type resultInfo struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func resultInfosFromNames(names []string) []resultInfo {
	r := make([]resultInfo, len(names))
	for n, name := range names {
		r[n] = resultInfo{name, utils.Snakify(name)}
	}
	return r
}
