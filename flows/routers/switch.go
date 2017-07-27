package routers

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// TypeSwitch is the constant for our switch router
const TypeSwitch string = "switch"

// Case represents a single case and test in our switch
type Case struct {
	UUID      flows.UUID     `json:"uuid"                 validate:"required"`
	Type      string         `json:"type"                 validate:"required"`
	Arguments []string       `json:"arguments,omitempty"`
	ExitUUID  flows.ExitUUID `json:"exit_uuid"            validate:"required"`
}

// SwitchRouter is a router which allows specifying 0-n cases which should each be tested in order, following
// whichever case returns true, or if none do, then taking the default exit
type SwitchRouter struct {
	Default flows.ExitUUID `json:"default_exit_uuid"   validate:"omitempty,uuid4"`
	Operand string         `json:"operand"             validate:"required"`
	Cases   []Case         `json:"cases"`
	BaseRouter
}

// Type returns the type of this router
func (r *SwitchRouter) Type() string { return TypeSwitch }

// Validate validates the arguments for this router
func (r *SwitchRouter) Validate(exits []flows.Exit) error {
	err := utils.ValidateAll(r)
	for _, c := range r.Cases {
		err = utils.ValidateAll(c)
		if err != nil {
			return err
		}

		// find the matching exit
		found := false
		for _, e := range exits {
			if e.UUID() == c.ExitUUID {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Exit '%s' missing from node", c.ExitUUID)
		}
	}

	return err
}

// PickRoute evaluates each of the tests on our cases in order, returning the exit for the first case which
// evaluates to a true. If no cases evaluate to true, then the default exit (if specified) is returned
func (r *SwitchRouter) PickRoute(run flows.FlowRun, exits []flows.Exit, step flows.Step) (flows.Route, error) {
	env := run.Environment()

	// first evaluate our operand
	operand, err := excellent.EvaluateTemplate(env, run.Context(), r.Operand)
	if err != nil {
		run.AddError(step, err)
	}

	// each of our cases
	for _, c := range r.Cases {
		test := strings.ToLower(c.Type)

		// try to look up our function
		xtest := excellent.XTESTS[test]
		if xtest == nil {
			run.AddError(step, fmt.Errorf("Unknown test '%s', taking no exit", c.Type))
			return flows.NoRoute, nil
		}

		// build our argument list
		args := make([]interface{}, len(c.Arguments)+1)
		args[0] = operand

		localizedArgs := run.GetTextArray(c.UUID, "arguments", c.Arguments)
		for i := range c.Arguments {
			test := localizedArgs[i]
			args[i+1], err = excellent.EvaluateTemplate(env, run.Context(), test)
			if err != nil {
				run.AddError(step, err)
			}
		}

		// call our function
		rawResult := xtest(env, args...)
		err, isErr := rawResult.(error)
		if isErr {
			return flows.NoRoute, err
		}

		// ok, not an error, must be an XTestResult
		result, isResult := rawResult.(excellent.XTestResult)
		if !isResult {
			return flows.NoRoute, fmt.Errorf("Unexpected result type from test %v: %#v", xtest, result)
		}

		// looks truthy, lets return this exit
		if result.Matched() {
			asStr, err := utils.ToString(env, result.Match())
			if err != nil {
				return flows.NoRoute, err
			}

			return flows.NewRoute(c.ExitUUID, asStr), nil
		}
	}

	// we have a default exit, use that
	if r.Default != "" {
		// evaluate our operand as a string
		value, err := utils.ToString(env, operand)
		if err != nil {
			run.AddError(step, err)
		}

		return flows.NewRoute(r.Default, value), nil
	}

	// no matches, no defaults, no route
	return flows.NoRoute, nil
}
