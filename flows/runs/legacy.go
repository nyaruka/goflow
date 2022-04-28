package runs

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

var invalidLegacyExtraKeyChars = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// keys in legacy @extra have non-word chars replaced with underscores and are limited to 255 chars
func legacyExtraKey(key string) string {
	key = invalidLegacyExtraKeyChars.ReplaceAllString(strings.ToLower(key), "_")
	return key[0:utils.Min(len(key), 255)]
}

type legacyExtra struct {
	values map[string]types.XValue
}

// creates a new legacy extra which will be lazily initialized on first call to .update()
func newLegacyExtra(run flows.Run) *legacyExtra {
	e := &legacyExtra{values: make(map[string]types.XValue)}

	// if trigger params is set, we include it in @extra
	triggerParams := run.Session().Trigger().Params()
	if triggerParams != nil {
		e.addValues(triggerParams)
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

func (e *legacyExtra) ToXValue(env envs.Environment) types.XValue {
	return types.NewXObject(e.values)
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

	e.values[utils.Snakify(result.Name)] = types.NewXText(string(result.Extra))

	values := types.JSONToXValue(result.Extra)
	e.addValues(values)
}

func (e *legacyExtra) addValues(values types.XValue) {
	switch typed := values.(type) {
	case *types.XObject:
		for _, key := range typed.Properties() {
			value, _ := typed.Get(key)
			e.values[legacyExtraKey(key)] = value
		}
	case *types.XArray:
		e.addValues(arrayToObject(typed))
	}
}

func arrayToObject(array *types.XArray) *types.XObject {
	properties := make(map[string]types.XValue, array.Count())
	for i := 0; i < array.Count(); i++ {
		properties[strconv.Itoa(i)] = array.Get(i)
	}
	return types.NewXObject(properties)
}

// finds the last webhook response that was saved as extra on a result
func lastWebhookSavedAsExtra(r *flowRun) types.XValue {
	for i := len(r.events) - 1; i >= 0; i-- {
		switch typed := r.events[i].(type) {
		case *events.WebhookCalledEvent:
			// look for a run result changed event on the same step
			resultEvent := r.findEvent(typed.StepUUID(), events.TypeRunResultChanged)

			if resultEvent != nil {
				asResultEvent := resultEvent.(*events.RunResultChangedEvent)
				if asResultEvent.Extra != nil {
					return types.JSONToXValue([]byte(asResultEvent.Extra))
				}
			}
		default:
			continue
		}
	}
	return nil
}
