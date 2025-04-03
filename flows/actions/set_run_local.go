package actions

import (
	"context"
	"fmt"
	"strconv"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeSetRunLocal, func() flows.Action { return &SetRunLocalAction{} })
}

// TypeSetRunLocal is the type for the set run result action
const TypeSetRunLocal string = "set_run_local"

// SetRunLocalAction can be used to save a local variable. The local will be available in the context
// for the run as @locals.[name]. The value field can be a template and will be evaluated.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "set_run_local",
//	  "name": "my_var",
//	  "value": "1",
//	  "operation": "increment"
//	}
//
// @action set_run_local
type SetRunLocalAction struct {
	baseAction
	universalAction

	Name      string `json:"name"                         validate:"required,local_name"`
	Value     string `json:"value"     engine:"evaluated"`
	Operation string `json:"operation"                    validate:"required,eq=set|eq=increment"`
}

// NewSetRunLocal creates a new set run local action
func NewSetRunLocal(uuid flows.ActionUUID, name, value string) *SetRunLocalAction {
	return &SetRunLocalAction{
		baseAction: newBaseAction(TypeSetRunLocal, uuid),
		Name:       name,
		Value:      value,
	}
}

// Execute runs this action
func (a *SetRunLocalAction) Execute(ctx context.Context, run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	value, ok := run.EvaluateTemplate(a.Value, logEvent)
	if !ok {
		return nil
	}

	if a.Operation == "increment" {
		existing, _ := strconv.Atoi(run.Locals().Get(a.Name))
		increment, err := strconv.Atoi(value)
		if err != nil {
			logEvent(events.NewError("unable to convert value to an integer"))
		} else {
			run.Locals().Set(a.Name, fmt.Sprint(existing+increment))
		}
	} else {
		run.Locals().Set(a.Name, value)
	}

	return nil
}
