package runs

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

var invalidLegacyExtraKeyChars = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// keys in legacy @extra have non-word chars replaced with underscores and are limited to 255 chars
func legacyExtraKey(key string) string {
	key = invalidLegacyExtraKeyChars.ReplaceAllString(strings.ToLower(key), "_")
	return key[0:utils.MinInt(len(key), 255)]
}

type legacyExtra struct {
	dict *types.XDict
}

// creates a new legacy extra which will be lazily initialized on first call to .update()
func newLegacyExtra(run flows.FlowRun) *legacyExtra {
	e := &legacyExtra{dict: types.NewEmptyXDict()}

	// if trigger params is a JSON object, we include it in @extra
	triggerParams := run.Session().Trigger().Params()
	asDict, isDict := triggerParams.(*types.XDict)
	if isDict {
		e.addValues(asDict)
	}

	// if trigger has results (i.e. a flow_action type trigger with a parent run) use them too
	asExtraContrib, isExtraContrib := run.Session().Trigger().(flows.LegacyExtraContributor)
	if isExtraContrib {
		e.addResults(asExtraContrib.LegacyExtra())
	}

	// add any existing results from this run
	e.addResults(run.Results())
	return e
}

func (e *legacyExtra) ToXValue(env utils.Environment) types.XValue {
	return e.dict
}

func (e *legacyExtra) addResults(results flows.Results) {
	// sort by created time
	sortedResults := make([]*flows.Result, 0)
	for _, result := range results {
		sortedResults = append(sortedResults, result)

	}
	sort.SliceStable(sortedResults, func(i, j int) bool { return sortedResults[i].CreatedOn.Before(sortedResults[j].CreatedOn) })

	// add each result in order
	for _, result := range sortedResults {
		e.addResult(result)
	}
}

// adds any extra from the given result
func (e *legacyExtra) addResult(result *flows.Result) {
	if result.Extra == nil {
		return
	}

	e.dict.Put(utils.Snakify(result.Name), types.NewXText(string(result.Extra)))

	values := types.JSONToXValue(result.Extra)
	e.addValues(values)
}

func (e *legacyExtra) addValues(values types.XValue) {
	switch typed := values.(type) {
	case *types.XDict:
		for _, key := range typed.Keys() {
			value, _ := typed.Get(key)
			e.dict.Put(legacyExtraKey(key), value)
		}
	case *types.XArray:
		e.addValues(arrayToDict(typed))
	}
}

func arrayToDict(array *types.XArray) *types.XDict {
	m := make(map[string]types.XValue, array.Length())
	for i := 0; i < array.Length(); i++ {
		m[strconv.Itoa(i)] = array.Get(i)
	}
	return types.NewXDict(m)
}
