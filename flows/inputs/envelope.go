package inputs

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type baseInputEnvelope struct {
	ChannelUUID flows.ChannelUUID `json:"channel_uuid" validate:"uuid4"`
	CreatedOn   time.Time         `json:"created_on"   validate:"required"`
}

func ReadInput(session flows.Session, envelope *utils.TypedEnvelope) (flows.Input, error) {
	switch envelope.Type {

	case TypeMsg:
		return ReadMsgInput(session, envelope)

	default:
		return nil, fmt.Errorf("Unknown input type: %s", envelope.Type)
	}
}
