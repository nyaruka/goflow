package routers

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// RouterFromEnvelope attempts to build a router given the passed in TypedEnvelope
func RouterFromEnvelope(envelope *utils.TypedEnvelope) (flows.Router, error) {
	switch envelope.Type {

	case TypeFirst:
		router := FirstRouter{}
		return &router, nil

	case TypeSwitch:
		router := SwitchRouter{}
		return &router, utils.UnmarshalAndValidate(envelope.Data, &router, "router")

	case TypeRandom:
		router := RandomRouter{}
		return &router, nil

	case TypeRandomOnce:
		router := RandomOnceRouter{}
		return &router, utils.UnmarshalAndValidate(envelope.Data, &router, "router")

	default:
		return nil, fmt.Errorf("Unknown router type: %s", envelope.Type)
	}
}
