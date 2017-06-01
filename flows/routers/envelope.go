package routers

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func RouterFromEnvelope(envelope *utils.TypedEnvelope) (flows.Router, error) {
	switch envelope.Type {

	case TypeFirst:
		router := FirstRouter{}
		return &router, nil

	case TypeSwitch:
		router := SwitchRouter{}
		err := json.Unmarshal(envelope.Data, &router)
		return &router, utils.ValidateAll(err, &router)

	case TypeRandom:
		router := RandomRouter{}
		return &router, nil

	case TypeRandomOnce:
		router := RandomOnceRouter{}
		err := json.Unmarshal(envelope.Data, &router)
		return &router, utils.ValidateAll(err, &router)

	default:
		return nil, fmt.Errorf("Unknown router type: %s", envelope.Type)
	}
}
