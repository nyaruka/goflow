package events

func init() {
	registerType(TypeTypingStarted, func() Event { return &TypingStarted{} })
}

// TypeTypingStarted is the type of our typing started event
const TypeTypingStarted string = "typing_started"

// TypingDirection is the direction of a typing indicator
type TypingDirection string

// possible values for typing directions
const (
	TypingDirectionIncoming TypingDirection = "incoming" // the contact is typing
	TypingDirectionOutgoing TypingDirection = "outgoing" // a user is typing
)

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

	Direction TypingDirection `json:"direction" validate:"required,eq=incoming|eq=outgoing"`
}

// NewTypingStarted returns a new typing started event
func NewTypingStarted(direction TypingDirection) *TypingStarted {
	return &TypingStarted{
		BaseEvent: NewBaseEvent(TypeTypingStarted),
		Direction: direction,
	}
}

var _ Event = (*TypingStarted)(nil)
