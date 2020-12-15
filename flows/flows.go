package flows

// FlowType represents the different types of flows
type FlowType string

const (
	// FlowTypeMessaging is a flow that is run over a messaging channel
	FlowTypeMessaging FlowType = "messaging"

	// FlowTypeMessagingPassive is a non-interactive messaging flow (i.e. never waits for input)
	FlowTypeMessagingPassive FlowType = "messaging_passive"

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

