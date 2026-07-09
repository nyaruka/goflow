package events

func init() {
	registerType(TypeTypingStopped, func() Event { return &TypingStopped{} })
}

// TypeTypingStopped is the type of our typing stopped event
const TypeTypingStopped string = "typing_stopped"

// TypingStopped events are created when the contact (direction of incoming) or a user (direction of outgoing)
// stops typing.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "typing_stopped",
//	  "created_on": "2019-01-02T15:04:05Z",
//	  "direction": "incoming"
//	}
//
// @event typing_stopped
type TypingStopped struct {
	BaseEvent

	Direction Direction `json:"direction" validate:"required,direction"`
}

// NewTypingStopped returns a new typing stopped event
func NewTypingStopped(direction Direction) *TypingStopped {
	return &TypingStopped{
		BaseEvent: NewBaseEvent(TypeTypingStopped),
		Direction: direction,
	}
}

var _ Event = (*TypingStopped)(nil)
