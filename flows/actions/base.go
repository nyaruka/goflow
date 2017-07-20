package actions

import (
	"github.com/nyaruka/goflow/flows"
)

// BaseAction is our base action type
type BaseAction struct {
	UUID_ flows.ActionUUID `json:"uuid"    validate:"required,uuid4"`
}

func NewBaseAction(uuid flows.ActionUUID) BaseAction {
	return BaseAction{UUID_: uuid}
}

func (a *BaseAction) UUID() flows.ActionUUID { return a.UUID_ }
