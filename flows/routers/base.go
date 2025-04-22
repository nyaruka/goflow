package routers

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/routers/waits"
	"github.com/nyaruka/goflow/utils"
)

var registeredTypes = map[string](func() flows.Router){}

// registers a new type of router
func registerType(name string, initFunc func() flows.Router) {
	registeredTypes[name] = initFunc
}

// RegisteredTypes gets the registered types
func RegisteredTypes() map[string](func() flows.Router) {
	return registeredTypes
}

// baseRouter is the base class for all router types
type baseRouter struct {
	type_      string
	wait       flows.Wait
	resultName string
	categories []flows.Category
}

// creates a new base router
func newBaseRouter(typeName string, wait flows.Wait, resultName string, categories []flows.Category) baseRouter {
	return baseRouter{type_: typeName, wait: wait, resultName: resultName, categories: categories}
}

// Type returns the type of this router
func (r *baseRouter) Type() string { return r.type_ }

// Wait returns the optional wait on this router
func (r *baseRouter) Wait() flows.Wait { return r.wait }

// Categories returns the categories on this router
func (r *baseRouter) Categories() []flows.Category { return r.categories }

// AllowTimeout returns whether this router can be resumed at with a timeout
func (r *baseRouter) AllowTimeout() bool {
	return r.wait != nil && r.wait.Timeout() != nil
}

// ResultName returns the name which the result of this router should be saved as (if any)
func (r *baseRouter) ResultName() string { return r.resultName }

func (r *baseRouter) Inspect(result func(*flows.ResultInfo), dependency func(assets.Reference)) {
	if r.resultName != "" {
		categoryNames := make([]string, len(r.categories))
		for i := range r.categories {
			categoryNames[i] = r.categories[i].Name()
		}

		result(flows.NewResultInfo(r.resultName, categoryNames))
	}
}

// EnumerateTemplates enumerates all expressions on this object and its children
func (r *baseRouter) EnumerateTemplates(localization flows.Localization, include func(i18n.Language, string)) {
}

// EnumerateLocalizables enumerates all the localizable text on this object
func (r *baseRouter) EnumerateLocalizables(include func(uuids.UUID, string, []string, func([]string))) {
	for _, cat := range r.categories {
		w := func(v []string) {
			cat.(*Category).name = v[0]
		}
		include(cat.LocalizationUUID(), "name", []string{cat.Name()}, w)
	}
}

func (r *baseRouter) validate(flow flows.Flow, exits []flows.Exit) error {
	if len(r.categories) > flows.MaxCategoriesPerRouter {
		return fmt.Errorf("can't have more than %d categories (has %d)", flows.MaxCategoriesPerRouter, len(r.categories))
	}

	// check wait timeout category is valid
	if r.AllowTimeout() && !r.isValidCategory(r.wait.Timeout().CategoryUUID()) {
		return fmt.Errorf("timeout category %s is not a valid category", r.wait.Timeout().CategoryUUID())
	}

	// check each category points to a valid exit
	for _, c := range r.categories {
		if c.ExitUUID() != "" && !r.isValidExit(c.ExitUUID(), exits) {
			return fmt.Errorf("category exit %s is not a valid exit", c.ExitUUID())
		}
	}

	if r.wait != nil && !flow.Type().Allows(r.wait) {
		return fmt.Errorf("wait type '%s' is not allowed in a flow of type '%s'", r.wait.Type(), flow.Type())
	}

	return nil
}

func (r *baseRouter) isValidCategory(uuid flows.CategoryUUID) bool {
	for _, c := range r.categories {
		if c.UUID() == uuid {
			return true
		}
	}
	return false
}

func (r *baseRouter) isValidExit(uuid flows.ExitUUID, exits []flows.Exit) bool {
	for _, e := range exits {
		if e.UUID() == uuid {
			return true
		}
	}
	return false
}

// RouteTimeout routes in the case that this router's wait timed out
func (r *baseRouter) RouteTimeout(run flows.Run, step flows.Step, logEvent flows.EventCallback) (flows.ExitUUID, error) {
	if !r.AllowTimeout() {
		return "", errors.New("can't call route timeout on router with no timeout")
	}

	// find last timeout event to use as time of timeout
	var timedOutOn time.Time
	runEvents := run.Events()
	for i := len(runEvents) - 1; i >= 0; i-- {
		event := runEvents[i]

		_, isTimeout := event.(*events.WaitTimedOutEvent)
		if isTimeout {
			timedOutOn = event.CreatedOn()
		}
	}

	return r.routeToCategory(run, step, r.wait.Timeout().CategoryUUID(), dates.FormatISO(timedOutOn), "", nil, logEvent)
}

func (r *baseRouter) routeToCategory(run flows.Run, step flows.Step, categoryUUID flows.CategoryUUID, match string, operand string, extra *types.XObject, logEvent flows.EventCallback) (flows.ExitUUID, error) {
	// router failed to pick a category
	if categoryUUID == "" {
		return "", nil
	}

	// find the actual category
	var category flows.Category
	for _, c := range r.categories {
		if c.UUID() == categoryUUID {
			category = c
			break
		}
	}

	if category == nil {
		return "", fmt.Errorf("category %s is not a valid category", categoryUUID)
	}

	// save result if we have a result name
	if r.resultName != "" {
		// localize the category name
		localizedCategory, _ := run.GetText(uuids.UUID(category.UUID()), "name", "")

		var extraJSON json.RawMessage
		if extra != nil {
			extraJSON, _ = jsonx.Marshal(extra)
		}
		result := flows.NewResult(r.resultName, match, category.Name(), localizedCategory, step.NodeUUID(), operand, extraJSON, dates.Now())
		prev, changed := run.SetResult(result)
		if changed {
			logEvent(events.NewRunResultChanged(result, prev))
		}
	}

	return category.ExitUUID(), nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type baseRouterEnvelope struct {
	Type       string            `json:"type"                  validate:"required"`
	Wait       json.RawMessage   `json:"wait,omitempty"`
	ResultName string            `json:"result_name,omitempty" validate:"omitempty,result_name"`
	Categories []json.RawMessage `json:"categories,omitempty"  validate:"required,min=1"`
}

// ReadRouter reads a router from the given JSON
func ReadRouter(data json.RawMessage) (flows.Router, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, fmt.Errorf("unknown type: '%s'", typeName)
	}

	router := f()
	return router, utils.UnmarshalAndValidate(data, router)
}

func (r *baseRouter) unmarshal(e *baseRouterEnvelope) error {
	var err error

	r.type_ = e.Type
	r.resultName = e.ResultName
	r.categories = make([]flows.Category, len(e.Categories))

	for i, c := range e.Categories {
		r.categories[i], err = ReadCategory(c)
		if err != nil {
			return err
		}
	}

	if e.Wait != nil {
		r.wait, err = waits.ReadWait(e.Wait)
		if err != nil {
			return fmt.Errorf("unable to read wait: %w", err)
		}
	}

	return nil
}

func (r *baseRouter) marshal(e *baseRouterEnvelope) error {
	var err error

	e.Type = r.type_
	e.ResultName = r.resultName
	e.Categories = make([]json.RawMessage, len(r.categories))

	for i, c := range r.categories {
		e.Categories[i], err = jsonx.Marshal(c)
		if err != nil {
			return err
		}
	}

	if r.wait != nil {
		e.Wait, err = jsonx.Marshal(r.wait)
		if err != nil {
			return err
		}
	}
	return nil
}
