package assets

import (
	"fmt"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/nyaruka/goflow/utils"
)

func init() {
	utils.Validator.RegisterStructValidation(GroupReferenceValidation, GroupReference{})
	utils.Validator.RegisterStructValidation(LabelReferenceValidation, LabelReference{})
}

// Reference is interface for all reference types
type Reference interface {
	fmt.Stringer
}

type MissingCallback func(Reference)

var PanicOnMissing MissingCallback = func(a Reference) { panic(fmt.Sprintf("unable to find asset %s", a.String())) }

var IgnoreOnMissing MissingCallback = func(Reference) {}

// ChannelReference is used to reference a channel
type ChannelReference struct {
	UUID ChannelUUID `json:"uuid" validate:"required,uuid"`
	Name string      `json:"name"`
}

func (r *ChannelReference) String() string {
	return fmt.Sprintf("channel[uuid=%s,name=%s]", r.UUID, r.Name)
}

var _ Reference = (*ChannelReference)(nil)

// NewChannelReference creates a new channel reference with the given UUID and name
func NewChannelReference(uuid ChannelUUID, name string) *ChannelReference {
	return &ChannelReference{UUID: uuid, Name: name}
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

func (r *GroupReference) String() string {
	return fmt.Sprintf("group[uuid=%s,name=%s]", r.UUID, r.Name)
}

var _ Reference = (*GroupReference)(nil)

// FieldReference is a reference to field
type FieldReference struct {
	Key  string `json:"key" validate:"required"`
	Name string `json:"name"`
}

// NewFieldReference creates a new field reference with the given key and label
func NewFieldReference(key string, name string) *FieldReference {
	return &FieldReference{Key: key, Name: name}
}

func (r *FieldReference) String() string {
	return fmt.Sprintf("field[key=%s,name=%s]", r.Key, r.Name)
}

var _ Reference = (*FieldReference)(nil)

// FlowReference is used to reference a flow from another flow
type FlowReference struct {
	UUID FlowUUID `json:"uuid" validate:"required,uuid4"`
	Name string   `json:"name"`
}

// NewFlowReference creates a new flow reference with the given UUID and name
func NewFlowReference(uuid FlowUUID, name string) *FlowReference {
	return &FlowReference{UUID: uuid, Name: name}
}

func (r *FlowReference) String() string {
	return fmt.Sprintf("flow[uuid=%s,name=%s]", r.UUID, r.Name)
}

var _ Reference = (*FlowReference)(nil)

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

func (r *LabelReference) String() string {
	return fmt.Sprintf("label[uuid=%s,name=%s]", r.UUID, r.Name)
}

var _ Reference = (*LabelReference)(nil)

//------------------------------------------------------------------------------------------
// Validation
//------------------------------------------------------------------------------------------

// GroupReferenceValidation validates that the given group reference is either a concrete
// reference or a name matcher
func GroupReferenceValidation(sl validator.StructLevel) {
	ref := sl.Current().Interface().(GroupReference)
	if neitherOrBoth(string(ref.UUID), ref.NameMatch) {
		sl.ReportError(ref.UUID, "uuid", "UUID", "mutually_exclusive", "name_match")
		sl.ReportError(ref.NameMatch, "name_match", "NameMatch", "mutually_exclusive", "uuid")
	}
}

// LabelReferenceValidation validates that the given label reference is either a concrete
// reference or a name matcher
func LabelReferenceValidation(sl validator.StructLevel) {
	ref := sl.Current().Interface().(LabelReference)
	if neitherOrBoth(string(ref.UUID), ref.NameMatch) {
		sl.ReportError(ref.UUID, "uuid", "UUID", "mutually_exclusive", "name_match")
		sl.ReportError(ref.NameMatch, "name_match", "NameMatch", "mutually_exclusive", "uuid")
	}
}

// utility method which returns true if both string values or neither string values is defined
func neitherOrBoth(s1 string, s2 string) bool {
	return (len(s1) > 0) == (len(s2) > 0)
}
