package routers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

var registeredTypes = map[string](func() flows.Router){}

// RegisterType registers a new type of router
func RegisterType(name string, initFunc func() flows.Router) {
	registeredTypes[name] = initFunc
}

// RegisteredTypes gets the registered types of router
func RegisteredTypes() map[string](func() flows.Router) {
	return registeredTypes
}

// BaseRouter is the base class for all our router classes
type BaseRouter struct {
	Type_       string      `json:"type"                  validate:"required"`
	ResultName_ string      `json:"result_name,omitempty"`
	Categories_ []*Category `json:"categories,omitempty"  validate:"required,min=1"`
}

func newBaseRouter(typeName string, resultName string, categories []*Category) BaseRouter {
	return BaseRouter{Type_: typeName, ResultName_: resultName, Categories_: categories}
}

// Type returns the type of this router
func (r *BaseRouter) Type() string { return r.Type_ }

// ResultName returns the name which the result of this router should be saved as (if any)
func (r *BaseRouter) ResultName() string { return r.ResultName_ }

// EnumerateTemplates enumerates all expressions on this object and its children
func (r *BaseRouter) EnumerateTemplates(localization flows.Localization, include func(string)) {}

// RewriteTemplates rewrites all templates on this object and its children
func (r *BaseRouter) RewriteTemplates(localization flows.Localization, rewrite func(string) string) {}

// EnumerateDependencies enumerates all dependencies on this object
func (r *BaseRouter) EnumerateDependencies(localization flows.Localization, include func(assets.Reference)) {
}

// EnumerateResults enumerates all potential results on this object
func (r *BaseRouter) EnumerateResults(include func(*flows.ResultSpec)) {
	if r.ResultName_ != "" {
		categoryNames := make([]string, len(r.Categories_))
		for c := range r.Categories_ {
			categoryNames[c] = r.Categories_[c].Name()
		}

		include(flows.NewResultSpec(r.ResultName_, categoryNames))
	}
}

func (r *BaseRouter) validate(exits []flows.Exit) error {
	// check each category points to a valid exit
	for _, c := range r.Categories_ {
		if c.ExitUUID() != "" && !r.isValidExit(c.ExitUUID(), exits) {
			return errors.Errorf("category exit %s is not a valid exit", c.ExitUUID())
		}
	}
	return nil
}

func (r *BaseRouter) isValidCategory(uuid flows.CategoryUUID) bool {
	for _, c := range r.Categories_ {
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

func (r *BaseRouter) routeToCategory(run flows.FlowRun, step flows.Step, categoryUUID flows.CategoryUUID, match string, input *string, extra map[string]string, logEvent flows.EventCallback) (flows.ExitUUID, error) {
	if categoryUUID == "" {
		return "", errors.New("switch router failed to pick an exit")
	}

	// find the actual category
	var category *Category
	for _, c := range r.Categories_ {
		if c.UUID() == categoryUUID {
			category = c
			break
		}
	}

	if category == nil {
		return "", errors.Errorf("category %s is not a valid category", categoryUUID)
	}

	// save result if we have a result name
	if r.ResultName_ != "" {
		// localize the category name
		localizedCategory := run.GetText(utils.UUID(category.UUID()), "name", "")

		var extraJSON json.RawMessage
		if extra != nil {
			extraJSON, _ = json.Marshal(extra)
		}
		result := flows.NewResult(r.ResultName_, match, category.Name(), localizedCategory, step.NodeUUID(), input, extraJSON, utils.Now())
		run.SaveResult(result)
		logEvent(events.NewRunResultChangedEvent(result))
	}

	return category.ExitUUID(), nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

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

	router := f()
	return router, utils.UnmarshalAndValidate(data, router)
}
