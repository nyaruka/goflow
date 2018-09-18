package runs

import (
	"encoding/json"
	"regexp"
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
	values map[string]interface{}
}

// NewLegacyExtra creates a new resolveable which provides backwards compatibility with @extra in legacy flows
func NewLegacyExtra(run flows.FlowRun) types.XValue {
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

	// if we have a webhook with a JSON payload, we include that too
	if run.Webhook() != nil {
		asMap, err := utils.JSONDecodeToMap(json.RawMessage(run.Webhook().Body()))
		if err == nil {
			for k, v := range asMap {
				values[legacyExtraKey(k)] = v
			}
		}
	}

	return &legacyExtra{values: values}
}

func (e *legacyExtra) Describe() string { return "legacy extra" }

func (e *legacyExtra) ToXJSON(env utils.Environment) types.XText {
	return types.MustMarshalToXText(e.values)
}

func (e *legacyExtra) Reduce(env utils.Environment) types.XPrimitive { return e.ToXJSON(env) }

func (e *legacyExtra) Resolve(env utils.Environment, key string) types.XValue {
	key = strings.ToLower(key)
	val, found := e.values[key]
	if !found {
		return types.NewXResolveError(e, key)
	}
	return e.convertToXValue(val)
}

func (e *legacyExtra) convertToXValue(val interface{}) types.XValue {
	switch typed := val.(type) {
	case map[string]interface{}:
		normalized := make(map[string]interface{}, len(typed))
		for k, v := range typed {
			normalized[legacyExtraKey(k)] = v
		}
		return &legacyExtra{values: normalized}
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

var _ types.XValue = (*legacyExtra)(nil)
var _ types.XResolvable = (*legacyExtra)(nil)
