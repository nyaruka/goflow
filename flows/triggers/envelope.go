package triggers

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type baseTriggerEnvelope struct {
	FlowUUID    flows.FlowUUID `json:"flow_uuid" validate:"uuid4"`
	TriggeredOn time.Time      `json:"triggered_on" validate:"required"`
}

func ReadTrigger(session flows.Session, envelope *utils.TypedEnvelope) (flows.Trigger, error) {
	switch envelope.Type {

	case TypeUser:
		return ReadUserTrigger(session, envelope)
	case TypeRun:
		return ReadRunTrigger(session, envelope)

	default:
		return nil, fmt.Errorf("unknown trigger type: %s", envelope.Type)
	}
}

func readBaseTrigger(session flows.Session, base *baseTrigger, envelope *baseTriggerEnvelope) error {
	var err error

	base.triggeredOn = envelope.TriggeredOn

	if base.flow, err = session.Assets().GetFlow(envelope.FlowUUID); err != nil {
		return err
	}

	return nil
}
