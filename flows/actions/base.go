package actions

import (
	"github.com/nyaruka/goflow/flows"
)

type BaseAction struct {
	Uuid flows.ActionUUID `json:"uuid"                     validate:"required"`
}
