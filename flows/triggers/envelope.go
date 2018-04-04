package triggers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type baseTriggerEnvelope struct {
	Environment json.RawMessage      `json:"environment,omitempty"`
	Flow        *flows.FlowReference `json:"flow" validate:"required"`
	Contact     json.RawMessage      `json:"contact,omitempty"`
	Params      json.RawMessage      `json:"params,omitempty"`
	TriggeredOn time.Time            `json:"triggered_on" validate:"required"`
}

func ReadTrigger(session flows.Session, envelope *utils.TypedEnvelope) (flows.Trigger, error) {
	switch envelope.Type {

	case TypeManual:
		return ReadManualTrigger(session, envelope)
	case TypeFlowAction:
		return ReadFlowActionTrigger(session, envelope)

	default:
		return nil, fmt.Errorf("unknown trigger type: %s", envelope.Type)
	}
}

func unmarshalBaseTrigger(session flows.Session, base *baseTrigger, envelope *baseTriggerEnvelope) error {
	var err error

	base.triggeredOn = envelope.TriggeredOn

	if base.flow, err = session.Assets().GetFlow(envelope.Flow.UUID); err != nil {
		return err
	}

	if envelope.Environment != nil {
		if base.environment, err = utils.ReadEnvironment(envelope.Environment); err != nil {
			return err
		}
	}
	if envelope.Contact != nil {
		if base.contact, err = flows.ReadContact(session, envelope.Contact); err != nil {
			return err
		}
	}
	if envelope.Params != nil {
		base.params = types.JSONFragment(envelope.Params)
	} else {
		base.params = types.EmptyJSONFragment
	}

	return nil
}

func marshalBaseTrigger(t *baseTrigger, envelope *baseTriggerEnvelope) error {
	var err error
	envelope.Flow = t.flow.Reference()
	envelope.TriggeredOn = t.triggeredOn

	if t.environment != nil {
		envelope.Environment, err = json.Marshal(t.environment)
		if err != nil {
			return err
		}
	}
	if t.contact != nil {
		envelope.Contact, err = json.Marshal(t.contact)
		if err != nil {
			return err
		}
	}
	if t.params != nil {
		envelope.Params = json.RawMessage(t.params)
	}
	return nil
}
