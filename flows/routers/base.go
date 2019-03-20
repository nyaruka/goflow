package routers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
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
	Type_       string      `json:"type" validate:"required"`
	ResultName_ string      `json:"result_name,omitempty"`
	Categories_ []*Category `json:"categories,omitempty"`
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

// EnumerateResultNames enumerates all result names on this object
func (r *BaseRouter) EnumerateResultNames(include func(string)) {
	include(r.ResultName())
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
