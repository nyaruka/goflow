package definition

import (
	"github.com/nyaruka/goflow/utils"
)

// holds information about a possible result in the flow
type resultInfo struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// creates a set of result infos with unique keys
func resultInfosFromNames(names []string) []resultInfo {
	keysSeen := make(map[string]bool)

	r := make([]resultInfo, 0, len(names))

	for _, name := range names {
		key := utils.Snakify(name)

		if _, seen := keysSeen[key]; !seen {
			r = append(r, resultInfo{name, key})
			keysSeen[key] = true
		}
	}
	return r
}
