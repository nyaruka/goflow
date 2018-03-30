package flows

import (
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/nyaruka/goflow/utils"
)

func init() {
	utils.Validator.RegisterStructValidation(GroupReferenceValidation, GroupReference{})
	utils.Validator.RegisterStructValidation(LabelReferenceValidation, LabelReference{})
}

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

// GroupReference is used to reference a group
type GroupReference struct {
	UUID      GroupUUID `json:"uuid,omitempty" validate:"omitempty,uuid4"`
	Name      string    `json:"name,omitempty"`
	NameMatch string    `json:"name_match,omitempty"`
}

// NewGroupReference creates a new group reference with the given UUID and name
func NewGroupReference(uuid GroupUUID, name string) *GroupReference {
	return &GroupReference{UUID: uuid, Name: name}
}

// NewVariableGroupReference creates a new group reference from the given templatized name match
func NewVariableGroupReference(nameMatch string) *GroupReference {
	return &GroupReference{NameMatch: nameMatch}
}

// FieldReference is a reference to field
type FieldReference struct {
	Key  FieldKey `json:"key" validate:"required"`
	Name string   `json:"name"`
}

// NewFieldReference creates a new field reference with the given key and label
func NewFieldReference(key FieldKey, label string) *FieldReference {
	return &FieldReference{Key: key, Name: label}
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

// LabelReference is used to reference a label
type LabelReference struct {
	UUID      LabelUUID `json:"uuid,omitempty" validate:"omitempty,uuid4"`
	Name      string    `json:"name,omitempty"`
	NameMatch string    `json:"name_match,omitempty"`
}

// NewLabelReference creates a new label reference with the given UUID and name
func NewLabelReference(uuid LabelUUID, name string) *LabelReference {
	return &LabelReference{UUID: uuid, Name: name}
}

// NewVariableLabelReference creates a new label reference from the given templatized name match
func NewVariableLabelReference(nameMatch string) *LabelReference {
	return &LabelReference{NameMatch: nameMatch}
}

//------------------------------------------------------------------------------------------
// Validation
//------------------------------------------------------------------------------------------

// GroupReferenceValidation validates that the given group reference is either a concrete
// reference or a name matcher
func GroupReferenceValidation(sl validator.StructLevel) {
	ref := sl.Current().Interface().(GroupReference)
	if neitherOrBoth(string(ref.UUID), ref.NameMatch) {
		sl.ReportError(ref.UUID, "UUID", "uuid", "mutually_exclusive", "name_match")
		sl.ReportError(ref.NameMatch, "NameMatch", "name_match", "mutually_exclusive", "uuid")
	}
}

// LabelReferenceValidation validates that the given label reference is either a concrete
// reference or a name matcher
func LabelReferenceValidation(sl validator.StructLevel) {
	ref := sl.Current().Interface().(LabelReference)
	if neitherOrBoth(string(ref.UUID), ref.NameMatch) {
		sl.ReportError(ref.UUID, "UUID", "uuid", "mutually_exclusive", "name_match")
		sl.ReportError(ref.NameMatch, "NameMatch", "name_match", "mutually_exclusive", "uuid")
	}
}

// utility method which returns true if both string values or neither string values is defined
func neitherOrBoth(s1 string, s2 string) bool {
	return (len(s1) > 0) == (len(s2) > 0)
}
