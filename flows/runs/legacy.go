package runs

import (
	"encoding/json"
	"regexp"
	"sort"
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

func (e legacyExtraMap) convertToXValue(val interface{}) types.XValue {
	switch typed := val.(type) {
	case map[string]interface{}:
		normalized := make(map[string]interface{}, len(typed))
		for k, v := range typed {
			normalized[legacyExtraKey(k)] = v
		}
		return legacyExtraMap(normalized)
	case []interface{}:
		xvals := make([]types.XValue, len(typed))
		for v := range typed {
			xvals[v] = e.convertToXValue(typed[v])
		}
		arr := types.NewXArray(xvals...)
		return arr
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

type legacyExtra struct {
	legacyExtraMap

	run            flows.FlowRun
	lastResultTime time.Time
}

func newLegacyExtra(run flows.FlowRun) *legacyExtra {
	values := make(map[string]interface{})

	// if trigger params is a JSON object, we include it in @extra
	triggerParams := run.Session().Trigger().Params()
	asJSON, isJSON := triggerParams.(types.XJSONObject)
	if isJSON {
		asMap, err := utils.JSONDecodeToMap(json.RawMessage(asJSON.XJSON))
		if err == nil {
			for k, v := range asMap {
				values[legacyExtraKey(k)] = v
			}
		}
	}

	return &legacyExtra{legacyExtraMap: values, run: run}
}

// updates @legacy_extra by looking for new events since we last updated
func (e *legacyExtra) update() {
	prevLastResultTime := e.lastResultTime

	// get all results with extra created since the last update
	newExtras := make([]*flows.Result, 0)
	for _, result := range e.run.Results() {
		if result.Extra != nil && result.CreatedOn.After(prevLastResultTime) {
			newExtras = append(newExtras, result)
		}
		e.lastResultTime = result.CreatedOn
	}
	// sort by created time
	sort.SliceStable(newExtras, func(i, j int) bool { return newExtras[i].CreatedOn.Before(newExtras[j].CreatedOn) })

	// add each extra blob to our master extra
	for _, result := range newExtras {
		e.legacyExtraMap["webhook"] = string(result.Extra)
		values, err := utils.JSONDecodeToMap(result.Extra)
		if err == nil {
			for k, v := range values {
				e.legacyExtraMap[legacyExtraKey(k)] = v
			}
		}
	}
}

var _ types.XValue = (*legacyExtra)(nil)
var _ types.XResolvable = (*legacyExtra)(nil)
