package inputs

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type baseInputEnvelope struct {
	UUID      flows.InputUUID         `json:"uuid"`
	Channel   *flows.ChannelReference `json:"channel,omitempty" validate:"omitempty,dive"`
	CreatedOn time.Time               `json:"created_on" validate:"required"`
}

func ReadInput(session flows.Session, envelope *utils.TypedEnvelope) (flows.Input, error) {
	switch envelope.Type {

	case TypeMsg:
		return ReadMsgInput(session, envelope.Data)

	default:
		return nil, fmt.Errorf("Unknown input type: %s", envelope.Type)
	}
}
