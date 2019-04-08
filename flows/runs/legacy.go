package runs

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

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
	dict types.XDict

	run            flows.FlowRun
	lastResultTime time.Time
}

// creates a new legacy extra which will be lazily initialized on first call to .update()
func newLegacyExtra(run flows.FlowRun) *legacyExtra {
	return &legacyExtra{run: run}
}

func (e *legacyExtra) ToXValue(env utils.Environment) types.XValue {
	return e.dict
}

func (e *legacyExtra) initialize() {
	e.dict = types.NewEmptyXDict()

	// if trigger params is a JSON object, we include it in @extra
	triggerParams := e.run.Session().Trigger().Params()
	asDict, isDict := triggerParams.(types.XDict)
	if isDict {
		e.addValues(asDict)
	}

	// if trigger has results (i.e. a flow_action type trigger with a parent run) use them too
	asExtraContrib, isExtraContrib := e.run.Session().Trigger().(flows.LegacyExtraContributor)
	if isExtraContrib {
		e.addResults(asExtraContrib.LegacyExtra(), time.Time{})
	}
}

// updates @legacy_extra by looking for new results since we last updated
func (e *legacyExtra) update() {
	// lazy initialize if necessary
	if e.dict == nil {
		e.initialize()
	}

	prevLastResultTime := e.lastResultTime

	e.addResults(e.run.Results(), prevLastResultTime)
}

// adds any results with extra to this blob of all extras
func (e *legacyExtra) addResults(results flows.Results, after time.Time) {
	// get all results with extra created since the last update
	newExtras := make([]*flows.Result, 0)
	for _, result := range results {
		if result.Extra != nil && result.CreatedOn.After(after) {
			newExtras = append(newExtras, result)
		}
		e.lastResultTime = result.CreatedOn
	}
	// sort by created time
	sort.SliceStable(newExtras, func(i, j int) bool { return newExtras[i].CreatedOn.Before(newExtras[j].CreatedOn) })

	// add each extra blob to our master extra
	for _, result := range newExtras {
		e.dict.Put(utils.Snakify(result.Name), types.NewXText(string(result.Extra)))

		values := types.JSONToXValue(result.Extra)
		e.addValues(values)
	}
}

func (e *legacyExtra) addValues(values types.XValue) {
	switch typed := values.(type) {
	case types.XDict:
		for _, key := range typed.Keys() {
			e.dict.Put(legacyExtraKey(key), typed.Get(key))
		}
	case types.XArray:
		e.addValues(arrayToDict(typed))
	}
}

func arrayToDict(array types.XArray) types.XDict {
	m := make(map[string]types.XValue, array.Length())
	for i := 0; i < array.Length(); i++ {
		m[strconv.Itoa(i)] = array.Index(i)
	}
	return types.NewXDict(m)
}
