package flows

// ChannelReference is used to reference a channel
type ChannelReference struct {
	UUID ChannelUUID `json:"uuid" validate:"required,uuid4"`
	Name string      `json:"name"`
}

// NewChannelReference creates a new channel reference with the given UUID and name
func NewChannelReference(uuid ChannelUUID, name string) *ChannelReference {
	return &ChannelReference{UUID: uuid, Name: name}
}

// ContactReference is used to reference a contact
type ContactReference struct {
	UUID ContactUUID `json:"uuid" validate:"required,uuid4"`
	Name string      `json:"name"`
}

// NewContactReference creates a new contact reference with the given UUID and name
func NewContactReference(uuid ContactUUID, name string) *ContactReference {
	return &ContactReference{UUID: uuid, Name: name}
}

// GroupReference is used to reference a contact group
type GroupReference struct {
	UUID GroupUUID `json:"uuid,omitempty" validate:"omitempty,uuid4"`
	Name string    `json:"name"`
}

// NewGroupReference creates a new flow reference with the given UUID and name
func NewGroupReference(uuid GroupUUID, name string) *GroupReference {
	return &GroupReference{UUID: uuid, Name: name}
}

// LabelReference is used to reference a label
type LabelReference struct {
	UUID LabelUUID `json:"uuid,omitempty" validate:"omitempty,uuid4"`
	Name string    `json:"name"`
}

// NewLabelReference creates a new flow reference with the given UUID and name
func NewLabelReference(uuid LabelUUID, name string) *LabelReference {
	return &LabelReference{UUID: uuid, Name: name}
}

// FlowReference is used to reference a flow from another flow
type FlowReference struct {
	UUID FlowUUID `json:"uuid" validate:"uuid4"`
	Name string   `json:"name"`
}

// NewFlowReference creates a new flow reference with the given UUID and name
func NewFlowReference(uuid FlowUUID, name string) *FlowReference {
	return &FlowReference{UUID: uuid, Name: name}
}

// FieldReference is a reference to field
type FieldReference struct {
	Key   FieldKey `json:"key" validate:"required"`
	Label string   `json:"label"`
}

// NewFieldReference creates a new field reference with the given key and label
func NewFieldReference(key FieldKey, label string) *FieldReference {
	return &FieldReference{Key: key, Label: label}
}
