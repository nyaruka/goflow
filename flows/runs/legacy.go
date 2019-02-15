package runs

import (
	"encoding/json"
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

type legacyExtraMap map[string]interface{}

func (e legacyExtraMap) Describe() string { return "legacy extra" }

func (e legacyExtraMap) ToXJSON(env utils.Environment) types.XText {
	return types.MustMarshalToXText(e)
}

func (e legacyExtraMap) Reduce(env utils.Environment) types.XPrimitive { return e.ToXJSON(env) }

func (e legacyExtraMap) Resolve(env utils.Environment, key string) types.XValue {
	key = strings.ToLower(key)
	val, found := e[key]
	if !found {
		return types.NewXResolveError(e, key)
	}
	return e.convertToXValue(val)
}

func (e legacyExtraMap) Index(index int) types.XValue {
	return e.Resolve(nil, strconv.Itoa(index))
}

func (e legacyExtraMap) Length() int {
	return len(e)
}

func (e legacyExtraMap) convertToXValue(val interface{}) types.XValue {
	switch typed := val.(type) {
	case map[string]interface{}:
		normalized := make(map[string]interface{}, len(typed))
		for k, v := range typed {
			normalized[legacyExtraKey(k)] = v
		}
		return legacyExtraMap(normalized)
	case []interface{}:
		normalized := sliceToMap(typed)
		return legacyExtraMap(normalized)
	case json.Number:
		return types.RequireXNumberFromString(string(typed))
	case string:
		return types.NewXText(typed)
	case bool:
		return types.NewXBoolean(typed)
	}
	return nil
}

var _ types.XValue = (legacyExtraMap)(nil)
var _ types.XResolvable = (legacyExtraMap)(nil)
var _ types.XIndexable = (legacyExtraMap)(nil)

type legacyExtra struct {
	legacyExtraMap

	run            flows.FlowRun
	lastResultTime time.Time
}

// creates a new legacy extra which will be lazily initialized on first call to .update()
func newLegacyExtra(run flows.FlowRun) *legacyExtra {
	return &legacyExtra{run: run}
}

func (e *legacyExtra) initialize() {
	e.legacyExtraMap = make(map[string]interface{})

	// if trigger params is a JSON object, we include it in @extra
	triggerParams := e.run.Session().Trigger().Params()
	asJSON, isJSON := triggerParams.(types.XJSONObject)
	if isJSON {
		values, err := utils.JSONDecodeGeneric(json.RawMessage(asJSON.XJSON))
		if err == nil {
			e.addValues(values)
		}
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
	if e.legacyExtraMap == nil {
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
		e.legacyExtraMap[utils.Snakify(result.Name)] = string(result.Extra)
		values, err := utils.JSONDecodeGeneric(result.Extra)

		// ignore unparseable extra
		if err != nil {
			return
		}

		e.addValues(values)
	}
}

func (e *legacyExtra) addValues(values interface{}) {
	switch typed := values.(type) {
	case map[string]interface{}:
		for k, v := range typed {
			e.legacyExtraMap[legacyExtraKey(k)] = v
		}
	case []interface{}:
		e.addValues(sliceToMap(typed))
	}
}

func sliceToMap(s []interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(s))
	for i, v := range s {
		m[strconv.Itoa(i)] = v
	}
	return m
}

var _ types.XValue = (*legacyExtra)(nil)
var _ types.XResolvable = (*legacyExtra)(nil)
