package routers

import (
	"fmt"
	"strings"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/inspect"
	"github.com/nyaruka/goflow/flows/routers/cases"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

func init() {
	registerType(TypeSwitch, func() flows.Router { return &SwitchRouter{} })
}

// TypeSwitch is the constant for our switch router
const TypeSwitch string = "switch"

// Case represents a single case and test in our switch
type Case struct {
	UUID         uuids.UUID         `json:"uuid"                   validate:"required"`
	Type         string             `json:"type"                   validate:"required"`
	Arguments    []string           `json:"arguments,omitempty"    engine:"localized,evaluated"`
	CategoryUUID flows.CategoryUUID `json:"category_uuid"          validate:"required"`
}

// NewCase creates a new case
func NewCase(uuid uuids.UUID, type_ string, arguments []string, categoryUUID flows.CategoryUUID) *Case {
	return &Case{
		UUID:         uuid,
		Type:         type_,
		Arguments:    arguments,
		CategoryUUID: categoryUUID,
	}
}

// LocalizationUUID gets the UUID which identifies this object for localization
func (c *Case) LocalizationUUID() uuids.UUID { return uuids.UUID(c.UUID) }

// Dependencies enumerates the dependencies on this case
func (c *Case) Dependencies(localization flows.Localization, include func(i18n.Language, assets.Reference)) {
	groupRef := func(args []string) assets.Reference {
		// if we have two args, the second is name
		name := ""
		if len(args) == 2 {
			name = args[1]
		}
		return assets.NewGroupReference(assets.GroupUUID(args[0]), name)
	}

	// currently only the HAS_GROUP router test can produce a dependency
	if c.Type == "has_group" && len(c.Arguments) > 0 {
		include(i18n.NilLanguage, groupRef(c.Arguments))

		// the group UUID might be different in different translations
		for _, lang := range localization.Languages() {
			arguments := localization.GetItemTranslation(lang, c.UUID, "arguments")
			if len(arguments) > 0 {
				include(lang, groupRef(arguments))
			}
		}
	}
}

// SwitchRouter is a router which allows specifying 0-n cases which should each be tested in order, following
// whichever case returns true, or if none do, then taking the default category
type SwitchRouter struct {
	baseRouter

	Operand             string             `json:"operand"  validate:"required"  engine:"evaluated"`
	Cases               []*Case            `json:"cases"`
	DefaultCategoryUUID flows.CategoryUUID `json:"default_category_uuid" validate:"omitempty,uuid4"`
}

// NewSwitch creates a new switch router
func NewSwitch(wait flows.Wait, resultName string, categories []flows.Category, operand string, cases []*Case, defaultCategoryUUID flows.CategoryUUID) *SwitchRouter {
	return &SwitchRouter{
		baseRouter:          newBaseRouter(TypeSwitch, wait, resultName, categories),
		DefaultCategoryUUID: defaultCategoryUUID,
		Operand:             operand,
		Cases:               cases,
	}
}

// Validate validates the arguments for this router
func (r *SwitchRouter) Validate(flow flows.Flow, exits []flows.Exit) error {
	// check the default category is valid
	if r.DefaultCategoryUUID != "" && !r.isValidCategory(r.DefaultCategoryUUID) {
		return errors.Errorf("default category %s is not a valid category", r.DefaultCategoryUUID)
	}

	for _, c := range r.Cases {
		// check each case points to a valid category
		if !r.isValidCategory(c.CategoryUUID) {
			return errors.Errorf("case category %s is not a valid category", c.CategoryUUID)
		}

		// and each case test is valid
		if _, exists := cases.XTESTS[c.Type]; !exists {
			return errors.Errorf("case test %s is not a registered test function", c.Type)
		}
	}

	return r.validate(flow, exits)
}

// Route determines which exit to take from a node
func (r *SwitchRouter) Route(run flows.Run, step flows.Step, logEvent flows.EventCallback) (flows.ExitUUID, string, error) {
	env := run.Session().MergedEnvironment()

	// first evaluate our operand
	operand, err := run.EvaluateTemplateValue(r.Operand)
	if err != nil {
		run.LogError(step, err)
	}

	var operandAsStr string

	if operand != nil {
		asText, _ := types.ToXText(env, operand)
		operandAsStr = asText.Native()
	}

	// find first matching case
	match, categoryUUID, extra, err := r.matchCase(run, step, operand)
	if err != nil {
		return "", "", err
	}

	// none of our cases matched, so try to use the default
	if categoryUUID == "" && r.DefaultCategoryUUID != "" {
		// evaluate our operand as a string
		value, xerr := types.ToXText(env, operand)
		if xerr != nil {
			run.LogError(step, xerr)
		}

		match = value.Native()
		categoryUUID = r.DefaultCategoryUUID
	}

	exit, err := r.routeToCategory(run, step, categoryUUID, match, operandAsStr, extra, logEvent)
	return exit, operandAsStr, err
}

func (r *SwitchRouter) matchCase(run flows.Run, step flows.Step, operand types.XValue) (string, flows.CategoryUUID, *types.XObject, error) {
	for _, c := range r.Cases {
		test := strings.ToLower(c.Type)

		// try to look up our function
		xtest := cases.XTESTS[test]
		if xtest == nil {
			return "", "", nil, errors.Errorf("unknown case test '%s'", c.Type)
		}

		// build our argument list which starts with the operand
		args := []types.XValue{operand}

		localizedArgs, _ := run.GetTextArray(c.UUID, "arguments", c.Arguments, nil)

		// this shouldn't happen but if the number of localized args doesn't match the base arguments, ignore them
		if len(localizedArgs) != len(c.Arguments) {
			localizedArgs = c.Arguments
		}

		for _, localizedArg := range localizedArgs {
			arg, err := run.EvaluateTemplateValue(localizedArg)
			if err != nil {
				run.LogError(step, err)
			}
			args = append(args, arg)
		}

		// call our function
		result := xtest.Call(run.Session().MergedEnvironment(), args)

		// tests have to return either errors or test results
		switch typed := result.(type) {
		case types.XError:
			// test functions can return an error
			run.LogError(step, errors.Errorf("error calling test %s: %s", xtest.Describe(), typed.Error()))
		case *types.XObject:
			matched := typed.Truthy()
			if !matched {
				continue
			}

			match, _ := typed.Get("match")
			extra, _ := typed.Get("extra")

			extraAsObject, isObject := extra.(*types.XObject)
			if extra != nil && !isObject {
				run.LogError(step, errors.Errorf("test %s returned non-object extra", strings.ToUpper(test)))
			}

			resultAsStr, xerr := types.ToXText(run.Session().MergedEnvironment(), match)
			if xerr != nil {
				return "", "", nil, xerr
			}

			return resultAsStr.Native(), c.CategoryUUID, extraAsObject, nil
		default:
			panic(fmt.Sprintf("unexpected result type from test %v: %#v", xtest, result))
		}
	}
	return "", "", nil, nil
}

// EnumerateDependencies enumerates all dependencies on this object and its children
func (r *SwitchRouter) EnumerateDependencies(localization flows.Localization, include func(i18n.Language, assets.Reference)) {
	inspect.Dependencies(r.Cases, localization, include)
}

// EnumerateLocalizables enumerates all the localizable text on this object
func (r *SwitchRouter) EnumerateLocalizables(include func(uuids.UUID, string, []string, func([]string))) {
	inspect.LocalizableText(r.Cases, include)

	r.baseRouter.EnumerateLocalizables(include)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type switchRouterEnvelope struct {
	baseRouterEnvelope

	Operand             string             `json:"operand"  validate:"required"`
	Cases               []*Case            `json:"cases"`
	DefaultCategoryUUID flows.CategoryUUID `json:"default_category_uuid" validate:"omitempty,uuid4"`
}

func (r *SwitchRouter) UnmarshalJSON(data []byte) error {
	e := &switchRouterEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return err
	}

	r.Operand = e.Operand
	r.Cases = e.Cases
	r.DefaultCategoryUUID = e.DefaultCategoryUUID

	if err := r.unmarshal(&e.baseRouterEnvelope); err != nil {
		return err
	}

	return nil
}

// MarshalJSON marshals this router into JSON
func (r *SwitchRouter) MarshalJSON() ([]byte, error) {
	e := &switchRouterEnvelope{
		Operand:             r.Operand,
		Cases:               r.Cases,
		DefaultCategoryUUID: r.DefaultCategoryUUID,
	}

	if err := r.marshal(&e.baseRouterEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
