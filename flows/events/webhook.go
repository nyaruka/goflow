package events

import "github.com/nyaruka/goflow/utils"

const WEBHOOK string = "webhook"

type WebhookEvent struct {
	URL        string                      `json:"url" validate:"nonzero"`
	Status     utils.RequestResponseStatus `json:"status" validate:"nonzero"`
	StatusCode int                         `json:"status_code" validate:"nonzero"`
	Request    string                      `json:"request" validate:"nonzero"`
	Response   string                      `json:"response"`

	BaseEvent
}

func (e *WebhookEvent) Type() string { return WEBHOOK }
