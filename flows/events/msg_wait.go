package events

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/routers/waits/hints"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

func init() {
	registerType(TypeMsgWait, func() flows.Event { return &MsgWaitEvent{} })
}

// TypeMsgWait is the type of our msg wait event
const TypeMsgWait string = "msg_wait"

// MsgWaitEvent events are created when a flow pauses waiting for a response from
// a contact. If a timeout is set, then the caller should resume the flow after
// the number of seconds in the timeout to resume it.
//
//   {
//     "type": "msg_wait",
//     "created_on": "2019-01-02T15:04:05Z",
//     "timeout_seconds": 300,
//     "hint": {
//        "type": "image"
//     }
//   }
//
// @event msg_wait
type MsgWaitEvent struct {
	baseEvent

	TimeoutSeconds *int       `json:"timeout_seconds,omitempty"`
	Hint           flows.Hint `json:"hint,omitempty"`
}

// NewMsgWait returns a new msg wait with the passed in timeout
func NewMsgWait(timeoutSeconds *int, hint flows.Hint) *MsgWaitEvent {
	return &MsgWaitEvent{
		baseEvent:      newBaseEvent(TypeMsgWait),
		TimeoutSeconds: timeoutSeconds,
		Hint:           hint,
	}
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type msgWaitEnvelope struct {
	baseEvent

	TimeoutSeconds *int            `json:"timeout_seconds,omitempty"`
	Hint           json.RawMessage `json:"hint,omitempty"`
}

// UnmarshalJSON unmarshals this event from the given JSON
func (e *MsgWaitEvent) UnmarshalJSON(data []byte) error {
	v := &msgWaitEnvelope{}
	if err := utils.UnmarshalAndValidate(data, v); err != nil {
		return err
	}

	e.baseEvent = v.baseEvent
	e.TimeoutSeconds = v.TimeoutSeconds

	var err error
	if v.Hint != nil {
		if e.Hint, err = hints.ReadHint(v.Hint); err != nil {
			return errors.Wrap(err, "unable to read hint")
		}
	}

	return nil
}
