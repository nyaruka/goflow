package waits

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func WaitFromEnvelope(envelope *utils.TypedEnvelope) (flows.Wait, error) {
	var wait flows.Wait

	switch envelope.Type {
	case TypeMsg:
		wait = &MsgWait{}
	case TypeFlow:
		wait = &FlowWait{}
	default:
		return nil, fmt.Errorf("Unknown wait type: %s", envelope.Type)
	}

	return wait, utils.UnmarshalAndValidate(envelope.Data, wait, fmt.Sprintf("wait[type=%s]", envelope.Type))
}
