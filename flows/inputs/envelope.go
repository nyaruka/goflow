package inputs

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func InputFromEnvelope(envelope *utils.TypedEnvelope) (flows.Input, error) {
	switch envelope.Type {

	case events.TypeMsgIn:
		event := events.MsgInEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		if err == nil {
			err = utils.ValidateAll(event)
		}
		return &event, envelope.TraceError(err)

	default:
		return nil, fmt.Errorf("Unknown input type: %s", envelope.Type)
	}
}
