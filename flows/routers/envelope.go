package routers

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// RouterFromEnvelope attempts to build a router given the passed in TypedEnvelope
func RouterFromEnvelope(envelope *utils.TypedEnvelope) (flows.Router, error) {
	var router flows.Router

	switch envelope.Type {
	case TypeFirst:
		router = &FirstRouter{}
	case TypeSwitch:
		router = &SwitchRouter{}
	case TypeRandom:
		router = &RandomRouter{}
	case TypeRandomOnce:
		router = &RandomOnceRouter{}
	default:
		return nil, fmt.Errorf("Unknown router type: %s", envelope.Type)
	}

	return router, utils.UnmarshalAndValidate(envelope.Data, router, fmt.Sprintf("router[type=%s]", envelope.Type))
}
