package routers

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func RouterFromEnvelope(envelope *utils.TypedEnvelope) (flows.Router, error) {
	switch envelope.Type {

	case FIRST:
		router := FirstRouter{}
		return &router, nil

	case SWITCH:
		router := SwitchRouter{}
		err := json.Unmarshal(envelope.Data, &router)
		return &router, envelope.TraceError(err)

	case RANDOM:
		router := RandomRouter{}
		return &router, nil

	case RANDOM_ONCE:
		router := RandomOnceRouter{}
		err := json.Unmarshal(envelope.Data, &router)
		return &router, envelope.TraceError(err)

	default:
		return nil, fmt.Errorf("Unknown router type: %s", envelope.Type)
	}
}
