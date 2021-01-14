package flows

import (
	"github.com/nyaruka/goflow/utils"

	validator "gopkg.in/go-playground/validator.v9"
)

func init() {
	utils.RegisterValidatorAlias("flow_type", "eq=messaging|eq=messaging_background|eq=messaging_offline|eq=voice", func(validator.FieldError) string {
		return "is not a valid flow type"
	})
}

// FlowType represents the different types of flows
type FlowType string

const (
	// FlowTypeMessaging is a flow that is run over a messaging channel
	FlowTypeMessaging FlowType = "messaging"

	// FlowTypeMessagingBackground is a non-interactive messaging flow (i.e. never waits for input)
	FlowTypeMessagingBackground FlowType = "messaging_background"

	// FlowTypeMessagingOffline is a flow which is run over an offline messaging client like Surveyor
	FlowTypeMessagingOffline FlowType = "messaging_offline"

	// FlowTypeVoice is a flow which is run over IVR
	FlowTypeVoice FlowType = "voice"
)

// Allows returns whether this flow type allows the given item
func (t FlowType) Allows(r FlowTypeRestricted) bool {
	for _, allowedType := range r.AllowedFlowTypes() {
		if t == allowedType {
			return true
		}
	}
	return false
}

// FlowTypeRestricted is a part of a flow which can be restricted to certain flow types
type FlowTypeRestricted interface {
	AllowedFlowTypes() []FlowType
}
