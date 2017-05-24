package waits

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func WaitFromEnvelope(envelope *utils.TypedEnvelope) (flows.Wait, error) {
	switch envelope.Type {

	case TypeMsg:
		wait := MsgWait{}
		err := json.Unmarshal(envelope.Data, &wait)
		return &wait, envelope.TraceError(err)

	case TypeFlow:
		wait := FlowWait{}
		err := json.Unmarshal(envelope.Data, &wait)
		return &wait, envelope.TraceError(err)

	default:
		return nil, fmt.Errorf("Unknown wait type: %s", envelope.Type)
	}
}
