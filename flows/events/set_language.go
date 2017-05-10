package events

import "github.com/nyaruka/goflow/flows"

const SET_LANGUAGE string = "set_language"

type SetLanguageEvent struct {
	Language flows.Language `json:"language"     validate:"nonzero"`
	BaseEvent
}

func (e *SetLanguageEvent) Type() string { return SET_LANGUAGE }
