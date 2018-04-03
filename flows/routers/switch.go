package routers

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/routers/tests"
	"github.com/nyaruka/goflow/utils"
)

// TypeSwitch is the constant for our switch router
const TypeSwitch string = "switch"

// Case represents a single case and test in our switch
type Case struct {
	UUID        utils.UUID     `json:"uuid"                 validate:"required"`
	Type        string         `json:"type"                 validate:"required"`
	Arguments   []string       `json:"arguments,omitempty"`
	OmitOperand bool           `json:"omit_operand,omitempty"`
	ExitUUID    flows.ExitUUID `json:"exit_uuid"            validate:"required"`
}

// SwitchRouter is a router which allows specifying 0-n cases which should each be tested in order, following
// whichever case returns true, or if none do, then taking the default exit
type SwitchRouter struct {
	BaseRouter
	Default flows.ExitUUID `json:"default_exit_uuid"   validate:"omitempty,uuid4"`
	Operand string         `json:"operand"             validate:"required"`
	Cases   []Case         `json:"cases"`
}

func NewSwitchRouter(defaultExit flows.ExitUUID, operand string, cases []Case, resultName string) *SwitchRouter {
	return &SwitchRouter{
		BaseRouter: BaseRouter{ResultName_: resultName},
		Default:    defaultExit,
		Operand:    operand,
		Cases:      cases,
	}
}

// Type returns the type of this router
func (r *SwitchRouter) Type() string { return TypeSwitch }

// Validate validates the arguments for this router
func (r *SwitchRouter) Validate(exits []flows.Exit) error {
	// helper to look for the given exit UUID
	hasExit := func(exitUUID flows.ExitUUID) bool {
		found := false
		for _, e := range exits {
			if e.UUID() == exitUUID {
				found = true
				break
			}
		}
		return found
	}

	if r.Default != "" && !hasExit(r.Default) {
		return fmt.Errorf("default exit %s is not a valid exit", r.Default)
	}

	for _, c := range r.Cases {
		if !hasExit(c.ExitUUID) {
			return fmt.Errorf("case exit %s is not a valid exit", c.ExitUUID)
		}
	}

	return nil
}

// PickRoute evaluates each of the tests on our cases in order, returning the exit for the first case which
// evaluates to a true. If no cases evaluate to true, then the default exit (if specified) is returned
func (r *SwitchRouter) PickRoute(run flows.FlowRun, exits []flows.Exit, step flows.Step) (interface{}, flows.Route, error) {
	env := run.Environment()

	// first evaluate our operand
	operand, err := run.EvaluateTemplate(r.Operand)
	if err != nil {
		run.AddError(step, nil, err)
	}

	// each of our cases
	for _, c := range r.Cases {
		test := strings.ToLower(c.Type)

		// try to look up our function
		xtest := tests.XTESTS[test]
		if xtest == nil {
			return operand, flows.NoRoute, fmt.Errorf("Unknown test '%s', taking no exit", c.Type)
		}

		// build our argument list
		args := make([]interface{}, 0, 1)
		if !c.OmitOperand {
			args = append(args, operand)
		}

		localizedArgs := run.GetTextArray(c.UUID, "arguments", c.Arguments)
		for i := range c.Arguments {
			test := localizedArgs[i]
			arg, err := run.EvaluateTemplate(test)
			if err != nil {
				run.AddError(step, nil, err)
			}
			args = append(args, arg)
		}

		// call our function
		rawResult := xtest(env, args...)
		err, isErr := rawResult.(error)
		if isErr {
			return operand, flows.NoRoute, err
		}

		// ok, not an error, must be an XTestResult
		result, isResult := rawResult.(tests.XTestResult)
		if !isResult {
			return operand, flows.NoRoute, fmt.Errorf("Unexpected result type from test %v: %#v", xtest, result)
		}

		// looks truthy, lets return this exit
		if result.Matched() {
			asStr, err := utils.ToString(env, result.Match())
			if err != nil {
				return operand, flows.NoRoute, err
			}

			return operand, flows.NewRoute(c.ExitUUID, asStr), nil
		}
	}

	// we have a default exit, use that
	if r.Default != "" {
		// evaluate our operand as a string
		value, err := utils.ToString(env, operand)
		if err != nil {
			run.AddError(step, nil, err)
		}

		return operand, flows.NewRoute(r.Default, value), nil
	}

	// no matches, no defaults, no route
	return operand, flows.NoRoute, nil
}
