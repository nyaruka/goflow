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

// BaseRouter is the base class for all our router classes
type BaseRouter struct {
	Type_ string `json:"type" validate:"required"`

	// ResultName_ is the name of the which the result of this router should be saved as (if any)
	ResultName_ string `json:"result_name,omitempty"`
}

func newBaseRouter(typeName string, resultName string) BaseRouter {
	return BaseRouter{Type_: typeName, ResultName_: resultName}
}

// Type returns the type of this router
func (r *BaseRouter) Type() string { return r.Type_ }

// ResultName returns the name which the result of this router should be saved as (if any)
func (r *BaseRouter) ResultName() string { return r.ResultName_ }

// EnumerateTemplates enumerates all expressions on this object and its children
func (r *BaseRouter) EnumerateTemplates(localization flows.Localization, callback func(string)) {}

// RewriteTemplates rewrites all templates on this object and its children
func (r *BaseRouter) RewriteTemplates(localization flows.Localization, rewrite func(string) string) {}

// EnumerateDependencies enumerates all dependencies on this object
func (r *BaseRouter) EnumerateDependencies(localization flows.Localization, callback func(assets.Reference)) {
}

// EnumerateResultNames enumerates all result names on this object
func (r *BaseRouter) EnumerateResultNames(callback func(string)) {}

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
