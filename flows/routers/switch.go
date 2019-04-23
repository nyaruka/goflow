package routers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/routers/cases"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

func init() {
	RegisterType(TypeSwitch, readSwitchRouter)
}

// TypeSwitch is the constant for our switch router
const TypeSwitch string = "switch"

// Case represents a single case and test in our switch
type Case struct {
	UUID         utils.UUID         `json:"uuid"                   validate:"required"`
	Type         string             `json:"type"                   validate:"required"`
	Arguments    []string           `json:"arguments,omitempty"`
	CategoryUUID flows.CategoryUUID `json:"category_uuid"          validate:"required"`
}

// NewCase creates a new case
func NewCase(uuid utils.UUID, type_ string, arguments []string, categoryUUID flows.CategoryUUID) *Case {
	return &Case{
		UUID:         uuid,
		Type:         type_,
		Arguments:    arguments,
		CategoryUUID: categoryUUID,
	}
}

// LocalizationUUID gets the UUID which identifies this object for localization
func (c *Case) LocalizationUUID() utils.UUID { return utils.UUID(c.UUID) }

// Inspect inspects this object and any children
func (c *Case) Inspect(inspect func(flows.Inspectable)) {
	inspect(c)
}

// EnumerateTemplates enumerates all expressions on this object and its children
func (c *Case) EnumerateTemplates(localization flows.Localization, include func(string)) {
	for _, arg := range c.Arguments {
		include(arg)
	}

	flows.EnumerateTemplateTranslations(localization, c, "arguments", include)
}

// RewriteTemplates rewrites all templates on this object and its children
func (c *Case) RewriteTemplates(localization flows.Localization, rewrite func(string) string) {
	for a := range c.Arguments {
		c.Arguments[a] = rewrite(c.Arguments[a])
	}

	flows.RewriteTemplateTranslations(localization, c, "arguments", rewrite)
}

// EnumerateDependencies enumerates all dependencies on this object and its children
func (c *Case) EnumerateDependencies(localization flows.Localization, include func(assets.Reference)) {
	// currently only the HAS_GROUP router test can produce a dependency
	if c.Type == "has_group" && len(c.Arguments) > 0 {
		include(assets.NewGroupReference(assets.GroupUUID(c.Arguments[0]), ""))

		// the group UUID might be different in different translations
		for _, lang := range localization.Languages() {
			arguments := localization.GetTranslations(lang).GetTextArray(c.UUID, "arguments")
			if len(arguments) > 0 {
				include(assets.NewGroupReference(assets.GroupUUID(arguments[0]), ""))
			}
		}
	}
}

// EnumerateResults enumerates all potential results on this object
func (c *Case) EnumerateResults(include func(*flows.ResultSpec)) {}

// SwitchRouter is a router which allows specifying 0-n cases which should each be tested in order, following
// whichever case returns true, or if none do, then taking the default category
type SwitchRouter struct {
	BaseRouter

	operand  string
	cases    []*Case
	default_ flows.CategoryUUID
}

// NewSwitchRouter creates a new switch router
func NewSwitchRouter(wait flows.Wait, resultName string, categories []*Category, operand string, cases []*Case, defaultCategory flows.CategoryUUID) *SwitchRouter {
	return &SwitchRouter{
		BaseRouter: newBaseRouter(TypeSwitch, wait, resultName, categories),
		default_:   defaultCategory,
		operand:    operand,
		cases:      cases,
	}
}

func (r *SwitchRouter) Cases() []*Case { return r.cases }

// Validate validates the arguments for this router
func (r *SwitchRouter) Validate(exits []flows.Exit) error {
	// check the default category is valid
	if r.default_ != "" && !r.isValidCategory(r.default_) {
		return errors.Errorf("default category %s is not a valid category", r.default_)
	}

	for _, c := range r.cases {
		// check each case points to a valid category
		if !r.isValidCategory(c.CategoryUUID) {
			return errors.Errorf("case category %s is not a valid category", c.CategoryUUID)
		}

		// and each case test is valid
		if _, exists := cases.XTESTS[c.Type]; !exists {
			return errors.Errorf("case test %s is not a registered test function", c.Type)
		}
	}

	return r.validate(exits)
}

// Route determines which exit to take from a node
func (r *SwitchRouter) Route(run flows.FlowRun, step flows.Step, logEvent flows.EventCallback) (flows.ExitUUID, error) {
	env := run.Environment()

	// first evaluate our operand
	operand, err := run.EvaluateTemplateValue(r.operand)
	if err != nil {
		run.LogError(step, err)
	}

	var input string

	if operand != nil {
		asText, _ := types.ToXText(env, operand)
		input = asText.Native()
	}

	// find first matching case
	match, categoryUUID, extra, err := r.matchCase(run, step, operand)
	if err != nil {
		return "", err
	}

	// none of our cases matched, so try to use the default
	if categoryUUID == "" && r.default_ != "" {
		// evaluate our operand as a string
		value, xerr := types.ToXText(env, operand)
		if xerr != nil {
			run.LogError(step, xerr)
		}

		match = value.Native()
		categoryUUID = r.default_
	}

	return r.routeToCategory(run, step, categoryUUID, match, input, extra, logEvent)
}

func (r *SwitchRouter) matchCase(run flows.FlowRun, step flows.Step, operand types.XValue) (string, flows.CategoryUUID, *types.XObject, error) {
	for _, c := range r.cases {
		test := strings.ToLower(c.Type)

		// try to look up our function
		xtest := cases.XTESTS[test]
		if xtest == nil {
			return "", "", nil, errors.Errorf("unknown case test '%s'", c.Type)
		}

		// build our argument list which starts with the operand
		args := []types.XValue{operand}

		localizedArgs := run.GetTextArray(c.UUID, "arguments", c.Arguments)
		for i := range c.Arguments {
			test := localizedArgs[i]
			arg, err := run.EvaluateTemplateValue(test)
			if err != nil {
				run.LogError(step, err)
			}
			args = append(args, arg)
		}

		// call our function
		result := xtest(run.Environment(), args...)

		// tests have to return either errors or test results
		switch typed := result.(type) {
		case types.XError:
			// test functions can return an error
			run.LogError(step, errors.Errorf("error calling test %s: %s", strings.ToUpper(test), typed.Error()))
		case *types.XObject:
			match, _ := typed.Get("match")
			extra, _ := typed.Get("extra")

			extraAsObject, isObject := extra.(*types.XObject)
			if extra != nil && !isObject {
				run.LogError(step, errors.Errorf("test %s returned non-object extra", strings.ToUpper(test)))
			}

			resultAsStr, xerr := types.ToXText(run.Environment(), match)
			if xerr != nil {
				return "", "", nil, xerr
			}

			return resultAsStr.Native(), c.CategoryUUID, extraAsObject, nil
		case nil:
			continue
		default:
			panic(fmt.Sprintf("unexpected result type from test %v: %#v", xtest, result))
		}
	}
	return "", "", nil, nil
}

// Inspect inspects this object and any children
func (r *SwitchRouter) Inspect(inspect func(flows.Inspectable)) {
	inspect(r)

	for _, cs := range r.cases {
		cs.Inspect(inspect)
	}
}

// EnumerateTemplates enumerates all expressions on this object and its children
func (r *SwitchRouter) EnumerateTemplates(localization flows.Localization, include func(string)) {
	include(r.operand)
}

// RewriteTemplates rewrites all templates on this object and its children
func (r *SwitchRouter) RewriteTemplates(localization flows.Localization, rewrite func(string) string) {
	r.operand = rewrite(r.operand)
}

// EnumerateDependencies enumerates all dependencies on this object and its children
func (r *SwitchRouter) EnumerateDependencies(localization flows.Localization, include func(assets.Reference)) {
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type switchRouterEnvelope struct {
	baseRouterEnvelope

	Operand string             `json:"operand"               validate:"required"`
	Cases   []*Case            `json:"cases"`
	Default flows.CategoryUUID `json:"default_category_uuid" validate:"omitempty,uuid4"`
}

func readSwitchRouter(data json.RawMessage) (flows.Router, error) {
	e := &switchRouterEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &SwitchRouter{
		operand:  e.Operand,
		cases:    e.Cases,
		default_: e.Default,
	}

	if err := r.unmarshal(&e.baseRouterEnvelope); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *SwitchRouter) MarshalJSON() ([]byte, error) {
	e := &switchRouterEnvelope{
		Operand: r.operand,
		Cases:   r.cases,
		Default: r.default_,
	}

	if err := r.marshal(&e.baseRouterEnvelope); err != nil {
		return nil, err
	}

	return json.Marshal(e)
}
