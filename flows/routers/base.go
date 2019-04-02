package routers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/waits"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

type readFunc func(json.RawMessage) (flows.Router, error)

var registeredTypes = map[string]readFunc{}

// RegisterType registers a new type of router
func RegisterType(name string, f readFunc) {
	registeredTypes[name] = f
}

// RegisteredTypes gets the registered types of router
func RegisteredTypes() []string {
	typeNames := make([]string, 0, len(registeredTypes))
	for typeName := range registeredTypes {
		typeNames = append(typeNames, typeName)
	}
	return typeNames
}

// BaseRouter is the base class for all our router classes
type BaseRouter struct {
	type_      string
	wait       flows.Wait
	resultName string
	categories []*Category
}

func newBaseRouter(typeName string, wait flows.Wait, resultName string, categories []*Category) BaseRouter {
	return BaseRouter{type_: typeName, wait: wait, resultName: resultName, categories: categories}
}

// Type returns the type of this router
func (r *BaseRouter) Type() string { return r.type_ }

// Wait returns the optional wait on this router
func (r *BaseRouter) Wait() flows.Wait { return r.wait }

// ResultName returns the name which the result of this router should be saved as (if any)
func (r *BaseRouter) ResultName() string { return r.resultName }

// EnumerateTemplates enumerates all expressions on this object and its children
func (r *BaseRouter) EnumerateTemplates(localization flows.Localization, include func(string)) {}

// RewriteTemplates rewrites all templates on this object and its children
func (r *BaseRouter) RewriteTemplates(localization flows.Localization, rewrite func(string) string) {}

// EnumerateDependencies enumerates all dependencies on this object
func (r *BaseRouter) EnumerateDependencies(localization flows.Localization, include func(assets.Reference)) {
}

// EnumerateResults enumerates all potential results on this object
func (r *BaseRouter) EnumerateResults(include func(*flows.ResultSpec)) {
	if r.resultName != "" {
		categoryNames := make([]string, len(r.categories))
		for c := range r.categories {
			categoryNames[c] = r.categories[c].Name()
		}

		include(flows.NewResultSpec(r.resultName, categoryNames))
	}
}

func (r *BaseRouter) validate(exits []flows.Exit) error {
	// check each category points to a valid exit
	for _, c := range r.categories {
		if c.ExitUUID() != "" && !r.isValidExit(c.ExitUUID(), exits) {
			return errors.Errorf("category exit %s is not a valid exit", c.ExitUUID())
		}
	}
	return nil
}

func (r *BaseRouter) isValidCategory(uuid flows.CategoryUUID) bool {
	for _, c := range r.categories {
		if c.UUID() == uuid {
			return true
		}
	}
	return false
}

func (r *BaseRouter) isValidExit(uuid flows.ExitUUID, exits []flows.Exit) bool {
	for _, e := range exits {
		if e.UUID() == uuid {
			return true
		}
	}
	return false
}

func (r *BaseRouter) routeToCategory(run flows.FlowRun, step flows.Step, categoryUUID flows.CategoryUUID, match string, input *string, extra types.XDict, logEvent flows.EventCallback) (flows.ExitUUID, error) {
	if categoryUUID == "" {
		return "", errors.New("switch router failed to pick an exit")
	}

	// find the actual category
	var category *Category
	for _, c := range r.categories {
		if c.UUID() == categoryUUID {
			category = c
			break
		}
	}

	if category == nil {
		return "", errors.Errorf("category %s is not a valid category", categoryUUID)
	}

	// save result if we have a result name
	if r.resultName != "" {
		// localize the category name
		localizedCategory := run.GetText(utils.UUID(category.UUID()), "name", "")

		var extraJSON json.RawMessage
		if extra != nil {
			extraJSON, _ = json.Marshal(extra)
		}
		result := flows.NewResult(r.resultName, match, category.Name(), localizedCategory, step.NodeUUID(), input, extraJSON, utils.Now())
		run.SaveResult(result)
		logEvent(events.NewRunResultChangedEvent(result))
	}

	return category.ExitUUID(), nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type baseRouterEnvelope struct {
	Type       string          `json:"type"                  validate:"required"`
	Wait       json.RawMessage `json:"wait,omitempty"`
	ResultName string          `json:"result_name,omitempty"`
	Categories []*Category     `json:"categories,omitempty"  validate:"required,min=1"`
}

// ReadRouter reads a router from the given JSON
func ReadRouter(data json.RawMessage) (flows.Router, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, errors.Errorf("unknown type: '%s'", typeName)
	}

	return f(data)
}

func (r *BaseRouter) unmarshal(e *baseRouterEnvelope) error {
	r.type_ = e.Type
	r.resultName = e.ResultName
	r.categories = e.Categories

	var err error

	if e.Wait != nil {
		r.wait, err = waits.ReadWait(e.Wait)
		if err != nil {
			return errors.Wrap(err, "unable to read wait")
		}
	}

	return nil
}

func (r *BaseRouter) marshal(e *baseRouterEnvelope) error {
	e.Type = r.type_
	e.ResultName = r.resultName
	e.Categories = r.categories

	var err error

	if r.wait != nil {
		e.Wait, err = json.Marshal(r.wait)
		if err != nil {
			return err
		}
	}
	return nil
}
