package events

func init() {
	registerType(TypeTypingStarted, func() Event { return &TypingStarted{} })
}

// TypeTypingStarted is the type of our typing started event
const TypeTypingStarted string = "typing_started"

// TypingStarted events are created when the contact (direction of incoming) or a user (direction of outgoing)
// starts typing.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "typing_started",
//	  "created_on": "2019-01-02T15:04:05Z",
//	  "direction": "incoming"
//	}
//
// @event typing_started
type TypingStarted struct {
	BaseEvent

	Direction Direction `json:"direction" validate:"required,direction"`
}

// NewTypingStarted returns a new typing started event
func NewTypingStarted(direction Direction) *TypingStarted {
	return &TypingStarted{
		BaseEvent: NewBaseEvent(TypeTypingStarted),
		Direction: direction,
	}
}

var _ Event = (*TypingStarted)(nil)
