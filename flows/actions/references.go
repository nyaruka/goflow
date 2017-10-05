package actions

import (
	"github.com/nyaruka/goflow/flows"
)

// GroupReference is used to reference a contact group
type GroupReference struct {
	UUID flows.GroupUUID `json:"uuid,omitempty" validate:"omitempty,uuid4"`
	Name string          `json:"name"`
}

// NewGroupReference creates a new flow reference with the given UUID and name
func NewGroupReference(uuid flows.GroupUUID, name string) *GroupReference {
	return &GroupReference{UUID: uuid, Name: name}
}

// ContactReference is used to reference a contact
type ContactReference struct {
	UUID flows.ContactUUID `json:"uuid" validate:"required,uuid4"`
	Name string            `json:"name"`
}

// NewContactReference creates a new flow reference with the given UUID and name
func NewContactReference(uuid flows.ContactUUID, name string) *ContactReference {
	return &ContactReference{UUID: uuid, Name: name}
}

// LabelReference is used to reference a label
type LabelReference struct {
	UUID flows.LabelUUID `json:"uuid,omitempty" validate:"omitempty,uuid4"`
	Name string          `json:"name"`
}

// NewLabelReference creates a new flow reference with the given UUID and name
func NewLabelReference(uuid flows.LabelUUID, name string) *LabelReference {
	return &LabelReference{UUID: uuid, Name: name}
}

// FlowReference is used to reference a flow from another flow
type FlowReference struct {
	UUID flows.FlowUUID `json:"uuid" validate:"uuid4"`
	Name string         `json:"name"`
}

// NewFlowReference creates a new flow reference with the given UUID and name
func NewFlowReference(uuid flows.FlowUUID, name string) *FlowReference {
	return &FlowReference{UUID: uuid, Name: name}
}

// FieldReference is a reference to field
type FieldReference struct {
	Key   flows.FieldKey `json:"key" validate:"required"`
	Label string         `json:"label"`
}

// NewFieldReference creates a new field reference with the given key and label
func NewFieldReference(key flows.FieldKey, label string) *FieldReference {
	return &FieldReference{Key: key, Label: label}
}
