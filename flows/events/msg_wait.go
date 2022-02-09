package events

import (
	"encoding/json"
	"time"

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
//     "created_on": "2022-01-03T13:27:30Z",
//     "timeout_seconds": 300,
//     "expires_on": "2022-02-02T13:27:30Z",
//     "hint": {
//        "type": "image"
//     }
//   }
//
// @event msg_wait
type MsgWaitEvent struct {
	BaseEvent

	// when this wait times out and we can proceed assuming router has a timeout category. This value is relative
	// because we want it to start counting when the last message is actually sent, which the engine can't know.
	TimeoutSeconds *int `json:"timeout_seconds,omitempty"`

	// When this wait expires and the whole run can be expired
	ExpiresOn *time.Time `json:"expires_on,omitempty"`

	Hint flows.Hint `json:"hint,omitempty"`
}

// NewMsgWait returns a new msg wait with the passed in timeout
func NewMsgWait(timeoutSeconds *int, expiresOn *time.Time, hint flows.Hint) *MsgWaitEvent {
	return &MsgWaitEvent{
		BaseEvent:      NewBaseEvent(TypeMsgWait),
		TimeoutSeconds: timeoutSeconds,
		ExpiresOn:      expiresOn,
		Hint:           hint,
	}
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type msgWaitEnvelope struct {
	BaseEvent

	TimeoutSeconds *int            `json:"timeout_seconds,omitempty"`
	ExpiresOn      *time.Time      `json:"expires_on,omitempty"`
	Hint           json.RawMessage `json:"hint,omitempty"`
}

// UnmarshalJSON unmarshals this event from the given JSON
func (e *MsgWaitEvent) UnmarshalJSON(data []byte) error {
	v := &msgWaitEnvelope{}
	if err := utils.UnmarshalAndValidate(data, v); err != nil {
		return err
	}

	e.BaseEvent = v.BaseEvent
	e.TimeoutSeconds = v.TimeoutSeconds
	e.ExpiresOn = v.ExpiresOn

	var err error
	if v.Hint != nil {
		if e.Hint, err = hints.ReadHint(v.Hint); err != nil {
			return errors.Wrap(err, "unable to read hint")
		}
	}

	return nil
}
