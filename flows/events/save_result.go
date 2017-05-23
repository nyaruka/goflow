package events

import "github.com/nyaruka/goflow/flows"

const SAVE_RESULT string = "save_result"

func NewResultEvent(node flows.NodeUUID, name string, value string, category string) *SaveResultEvent {
	return &SaveResultEvent{Node: node, Name: name, Value: value, Category: category}
}

type SaveResultEvent struct {
	Node     flows.NodeUUID `json:"node"        validate:"required"`
	Name     string         `json:"name"        validate:"required"`
	Value    string         `json:"value"       validate:"required"`
	Category string         `json:"category"    validate:"required"`
	BaseEvent
}

func (e *SaveResultEvent) Type() string { return SAVE_RESULT }
