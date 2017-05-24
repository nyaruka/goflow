package actions

import (
	"github.com/nyaruka/goflow/flows"
)

// BaseAction is our base action type
type BaseAction struct {
	UUID flows.ActionUUID `json:"uuid"    validate:"required,uuid4"`
}
