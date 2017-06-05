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

	case events.TypeMsgReceived:
		event := events.MsgReceivedEvent{}
		err := json.Unmarshal(envelope.Data, &event)
		return &event, utils.ValidateAllUnlessErr(err, &event)

	default:
		return nil, fmt.Errorf("Unknown input type: %s", envelope.Type)
	}
}
