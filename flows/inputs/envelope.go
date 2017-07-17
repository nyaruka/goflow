package inputs

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type baseInputEnvelope struct {
	ChannelUUID flows.ChannelUUID `json:"channel_uuid"` // TODO validate:"required,uuid4"`
	CreatedOn   time.Time         `json:"created_on"   validate:"required"`
}

func InputFromEnvelope(env flows.FlowEnvironment, envelope *utils.TypedEnvelope) (flows.Input, error) {
	switch envelope.Type {

	case TypeMsg:
		input, err := ReadMsgInput(env, envelope)
		return input, utils.ValidateAllUnlessErr(err, input)

	default:
		return nil, fmt.Errorf("Unknown input type: %s", envelope.Type)
	}
}
