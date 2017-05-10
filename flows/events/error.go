package events

const ERROR string = "error"

type ErrorEvent struct {
	Text string `json:"text"     validate:"nonzero"`
	BaseEvent
}

func (e *ErrorEvent) Type() string { return ERROR }
