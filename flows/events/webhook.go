package events

import "github.com/nyaruka/goflow/utils"

const WEBHOOK string = "webhook"

type WebhookEvent struct {
	URL        string                      `json:"url"         validate:"required"`
	Status     utils.RequestResponseStatus `json:"status"      validate:"required"`
	StatusCode int                         `json:"status_code" validate:"required"`
	Request    string                      `json:"request"     validate:"required"`
	Response   string                      `json:"response"`

	BaseEvent
}

func (e *WebhookEvent) Type() string { return WEBHOOK }
