package routers

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/routers/tests"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

func init() {
	RegisterType(TypeSwitch, func() flows.Router { return &SwitchRouter{} })
}

// TypeSwitch is the constant for our switch router
const TypeSwitch string = "switch"

// Case represents a single case and test in our switch
type Case struct {
	UUID         utils.UUID         `json:"uuid"                   validate:"required"`
	Type         string             `json:"type"                   validate:"required"`
	Arguments    []string           `json:"arguments,omitempty"`
	OmitOperand  bool               `json:"omit_operand,omitempty"`
	CategoryUUID flows.CategoryUUID `json:"category_uuid"          validate:"required"`
}

// NewCase creates a new case
func NewCase(uuid utils.UUID, type_ string, arguments []string, omitOperand bool, categoryUUID flows.CategoryUUID) *Case {
	return &Case{
		UUID:         uuid,
		Type:         type_,
		Arguments:    arguments,
		OmitOperand:  omitOperand,
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
	Operand string             `json:"operand"             validate:"required"`
	Cases   []*Case            `json:"cases"`
	Default flows.CategoryUUID `json:"default_category_uuid"   validate:"omitempty,uuid4"`
}

// NewSwitchRouter creates a new switch router
func NewSwitchRouter(resultName string, categories []*Category, operand string, cases []*Case, defaultCategory flows.CategoryUUID) *SwitchRouter {
	return &SwitchRouter{
		BaseRouter: newBaseRouter(TypeSwitch, resultName, categories),
		Default:    defaultCategory,
		Operand:    operand,
		Cases:      cases,
	}
}

// Validate validates the arguments for this router
func (r *SwitchRouter) Validate(exits []flows.Exit) error {
	// check the default category is valid
	if r.Default != "" && !r.isValidCategory(r.Default) {
		return errors.Errorf("default category %s is not a valid category", r.Default)
	}

	// check each case points to a valid category
	for _, c := range r.Cases {
		if !r.isValidCategory(c.CategoryUUID) {
			return errors.Errorf("case category %s is not a valid category", c.CategoryUUID)
		}
	}

	return r.validate(exits)
}

// PickExit determines which exit to take from a node
func (r *SwitchRouter) PickExit(run flows.FlowRun, step flows.Step, logEvent flows.EventCallback) (flows.ExitUUID, error) {
	env := run.Environment()

	// first evaluate our operand
	operand, err := run.EvaluateTemplateValue(r.Operand)
	if err != nil {
		run.LogError(step, err)
	}

	var input *string

	if operand != nil {
		asText, _ := types.ToXText(env, operand)
		asString := asText.Native()
		input = &asString
	}

	// find first matching case
	match, categoryUUID, extra, err := r.matchCase(run, step, operand)
	if err != nil {
		return "", err
	}

	// none of our cases matched, so try to use the default
	if categoryUUID == "" && r.Default != "" {
		// evaluate our operand as a string
		value, xerr := types.ToXText(env, operand)
		if xerr != nil {
			run.LogError(step, xerr)
		}

		match = value.Native()
		categoryUUID = r.Default
	}

	return r.routeToCategory(run, step, categoryUUID, match, input, extra, logEvent)
}

func (r *SwitchRouter) matchCase(run flows.FlowRun, step flows.Step, operand types.XValue) (string, flows.CategoryUUID, types.XDict, error) {
	for _, c := range r.Cases {
		test := strings.ToLower(c.Type)

		// try to look up our function
		xtest := tests.XTESTS[test]
		if xtest == nil {
			return "", "", nil, errors.Errorf("unknown case test '%s'", c.Type)
		}

		// build our argument list
		args := make([]types.XValue, 0, 1)
		if !c.OmitOperand {
			args = append(args, operand)
		}

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
		switch typedResult := result.(type) {
		case types.XError:
			// test functions can return an error
			run.LogError(step, errors.Errorf("error calling test %s: %s", strings.ToUpper(test), typedResult.Error()))
		case tests.XTestResult:
			if typedResult.Matched() {
				resultAsStr, xerr := types.ToXText(run.Environment(), typedResult.Match())
				if xerr != nil {
					return "", "", nil, xerr
				}

				return resultAsStr.Native(), c.CategoryUUID, typedResult.Extra(), nil
			}
		default:
			panic(fmt.Sprintf("unexpected result type from test %v: %#v", xtest, result))
		}
	}
	return "", "", nil, nil
}

// Inspect inspects this object and any children
func (r *SwitchRouter) Inspect(inspect func(flows.Inspectable)) {
	inspect(r)

	for _, cs := range r.Cases {
		cs.Inspect(inspect)
	}
}

// EnumerateTemplates enumerates all expressions on this object and its children
func (r *SwitchRouter) EnumerateTemplates(localization flows.Localization, include func(string)) {
	include(r.Operand)
}

// RewriteTemplates rewrites all templates on this object and its children
func (r *SwitchRouter) RewriteTemplates(localization flows.Localization, rewrite func(string) string) {
	r.Operand = rewrite(r.Operand)
}

// EnumerateDependencies enumerates all dependencies on this object and its children
func (r *SwitchRouter) EnumerateDependencies(localization flows.Localization, include func(assets.Reference)) {
}
