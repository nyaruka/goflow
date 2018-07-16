package transferto

import (
	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/routers/tests"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	functions.RegisterXFunction("has_airtime_status", functions.TwoArgFunction(HasAirtimeStatus))
}

// HasAirtimeStatus returns whether the last airtime transfer has the given status. If there
// are no airtime transfer events, it returns false.
//
//   @(has_airtime_status(run, "success")) -> false
//
// @test has_airtime_status(run, status)
func HasAirtimeStatus(env utils.Environment, arg1 types.XValue, arg2 types.XValue) types.XValue {
	// first parameter needs to be a flow run
	run, isRun := arg1.(flows.FlowRun)
	if !isRun {
		return types.NewXErrorf("must be called with a run as first argument")
	}

	status, xerr := types.ToXText(env, arg2)
	if xerr != nil {
		return xerr
	}

	// look to see if the last transfer event has the given status
	runEvents := run.Events()
	for e := len(runEvents) - 1; e >= 0; e-- {
		event := runEvents[e]

		asTransfer, isTransfer := event.(*AirtimeTransferedEvent)
		if isTransfer {
			if status.Native() == asTransfer.Status {
				return tests.NewTrueResult(types.NewXText(asTransfer.Status))
			}
			break
		}
	}

	return tests.XFalseResult
}
