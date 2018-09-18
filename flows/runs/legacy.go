package runs

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
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

	run           flows.FlowRun
	lastEventTime time.Time
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
	prevLastEventTime := e.lastEventTime

	for _, event := range e.run.Events() {
		if !event.CreatedOn().After(prevLastEventTime) {
			continue
		}

		switch typed := event.(type) {
		case *events.WebhookCalledEvent:
			e.addPossibleJSONResponse(typed.Response)
		case *events.ResthookCalledEvent:
			for _, call := range typed.Calls {
				e.addPossibleJSONResponse(call.Response)
			}
		case *events.RunResultChangedEvent:
			for k, v := range typed.Extra {
				e.legacyExtraMap[legacyExtraKey(k)] = v
			}
		}

		e.lastEventTime = event.CreatedOn()
	}
}

// tries to parse the given response as JSON and if successful adds it to this @extra
func (e *legacyExtra) addPossibleJSONResponse(response string) {
	parts := strings.SplitN(response, "\r\n\r\n", 2)
	if len(parts) != 2 {
		return
	}
	values, err := utils.JSONDecodeToMap(json.RawMessage(parts[1]))
	if err != nil {
		return
	}
	for k, v := range values {
		e.legacyExtraMap[legacyExtraKey(k)] = v
	}
}

var _ types.XValue = (*legacyExtra)(nil)
var _ types.XResolvable = (*legacyExtra)(nil)
