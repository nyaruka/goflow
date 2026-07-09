package flows

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"maps"

	"github.com/nyaruka/goflow/core"

	"github.com/nyaruka/gocommon/stringsx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// Results is our wrapper around a map of snakified result names to result objects
type Results map[string]*core.Result

// NewResults creates a new empty set of results
func NewResults() Results {
	return make(Results)
}

// Clone returns a clone of this results set
func (r Results) Clone() Results {
	clone := make(Results, len(r))
	maps.Copy(clone, r)
	return clone
}

// Save saves a new result in our map using the snakified name as the key. Returns the old result if it existed.
func (r Results) Save(result *core.Result) (*core.Result, bool) {
	key := utils.Snakify(result.Name)
	old := r[key]
	r[key] = result

	if old == nil || (old.Value != result.Value || old.Category != result.Category) {
		return old, true
	}
	return nil, false
}

// Get returns the result with the given key
func (r Results) Get(key string) *core.Result {
	return r[key]
}

// Context returns the properties available in expressions
func (r Results) Context(env envs.Environment) map[string]types.XValue {
	entries := make(map[string]types.XValue, len(r)+1)
	entries["__default__"] = types.NewXText(r.format())

	for k, v := range r {
		entries[k] = core.Context(env, v)
	}
	return entries
}

func (r Results) format() string {
	lines := make([]string, 0, len(r))
	for _, v := range r {
		lines = append(lines, fmt.Sprintf("%s: %s", v.Name, v.Value))
	}

	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

func (r *Results) UnmarshalJSON(data []byte) error {
	var m map[string]*core.Result
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	*r = make(Results, len(m))

	// we enforce result names being at most 64 chars but old sessions may have longer names
	for _, v := range m {
		v.Name = strings.TrimSpace(stringsx.Truncate(v.Name, 64))
		(*r)[utils.Snakify(v.Name)] = v
	}

	return nil
}
