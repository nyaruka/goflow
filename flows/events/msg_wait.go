package events

const MSG_WAIT string = "msg_wait"

type MsgWaitEvent struct {
	Timeout int `json:"timeout"`
	BaseEvent
}

func (e *MsgWaitEvent) Type() string { return MSG_WAIT }
