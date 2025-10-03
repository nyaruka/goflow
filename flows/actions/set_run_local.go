package actions

import (
	"context"
	"fmt"
	"strconv"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeSetRunLocal, func() flows.Action { return &SetRunLocal{} })
}

// TypeSetRunLocal is the type for the set run result action
const TypeSetRunLocal string = "set_run_local"

type LocalOperation string

const (
	LocalOperationSet       LocalOperation = "set"
	LocalOperationIncrement LocalOperation = "increment"
	LocalOperationClear     LocalOperation = "clear"
)

// SetRunLocal can be used to save a local variable. The local will be available in the context
// for the run as @locals.[local]. The value field can be a template and will be evaluated.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "set_run_local",
//	  "local": "my_var",
//	  "value": "1",
//	  "operation": "increment"
//	}
//
// @action set_run_local
type SetRunLocal struct {
	baseAction
	universalAction

	Local     string         `json:"local"                              validate:"required,local_ref"`
	Value     string         `json:"value,omitempty" engine:"evaluated" validate:"max=1000"`
	Operation LocalOperation `json:"operation"                          validate:"required,eq=set|eq=increment|eq=clear"`
}

// NewSetRunLocal creates a new set run local action
func NewSetRunLocal(uuid flows.ActionUUID, local, value string) *SetRunLocal {
	return &SetRunLocal{
		baseAction: newBaseAction(TypeSetRunLocal, uuid),
		Local:      local,
		Value:      value,
	}
}

// Execute runs this action
func (a *SetRunLocal) Execute(ctx context.Context, run flows.Run, step flows.Step, logEvent flows.EventCallback) error {
	value, ok := run.EvaluateTemplate(a.Value, logEvent)
	if !ok {
		return nil
	}

	if a.Operation == LocalOperationSet {
		run.Locals().Set(a.Local, value)
	} else if a.Operation == LocalOperationIncrement {
		existing, _ := strconv.Atoi(run.Locals().Get(a.Local))
		increment, err := strconv.Atoi(value)
		if err != nil {
			logEvent(events.NewError("increment value is not an integer"))
		} else {
			run.Locals().Set(a.Local, fmt.Sprint(existing+increment))
		}
	} else {
		run.Locals().Clear(a.Local)
	}

	return nil
}

func (a *SetRunLocal) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	local(a.Local)
}
