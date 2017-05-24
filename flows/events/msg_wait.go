package events

// TypeMsgWait is the type of our msg wait event
const TypeMsgWait string = "msg_wait"

// MsgWaitEvent events are created when a flow pauses waiting for a response from
// a contact. If a timeout is set, then the caller should resume the flow after
// the number of seconds in the timeout to resume it.
//
// ```
//   {
//    "step": "8eebd020-1af5-431c-b943-aa670fc74da9",
//    "created_on": "2006-01-02T15:04:05Z",
//    "type": "msg_wait",
//    "timeout": 300
//   }
// ```
//
// @event msg_wait
type MsgWaitEvent struct {
	Timeout int `json:"timeout"`
	BaseEvent
}

// NewMsgWait returns a new msg wait with the passed in timeout
func NewMsgWait(timeout int) *MsgWaitEvent {
	return &MsgWaitEvent{Timeout: timeout}
}

// Type returns the type of this event
func (e *MsgWaitEvent) Type() string { return TypeMsgWait }
