package triggers

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type baseTriggerEnvelope struct {
	TriggeredOn time.Time `json:"triggered_on" validate:"required"`
}

func ReadTrigger(session flows.Session, envelope *utils.TypedEnvelope) (flows.Trigger, error) {
	switch envelope.Type {

	case TypeRun:
		return ReadRunTrigger(session, envelope)

	default:
		return nil, fmt.Errorf("unknown trigger type: %s", envelope.Type)
	}
}
