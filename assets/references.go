package assets

import (
	"fmt"

	validator "gopkg.in/go-playground/validator.v9"

	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/jsonx"
	"github.com/nyaruka/goflow/utils/uuids"
)

func init() {
	utils.Validator.RegisterStructValidation(GroupReferenceValidation, GroupReference{})
	utils.Validator.RegisterStructValidation(LabelReferenceValidation, LabelReference{})
}

// Reference is interface for all reference types
type Reference interface {
	fmt.Stringer

	Type() string
	Identity() string
	Variable() bool
}

// UUIDReference is interface for all reference types that contain a UUID
type UUIDReference interface {
	Reference
	GenericUUID() uuids.UUID
}

// ChannelReference is used to reference a channel
type ChannelReference struct {
	UUID ChannelUUID `json:"uuid" validate:"required,uuid"`
	Name string      `json:"name"`
}

// NewChannelReference creates a new channel reference with the given UUID and name
func NewChannelReference(uuid ChannelUUID, name string) *ChannelReference {
	return &ChannelReference{UUID: uuid, Name: name}
}

// Type returns the name of the asset type
func (r *ChannelReference) Type() string {
	return "channel"
}

// GenericUUID returns the untyped UUID
func (r *ChannelReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity returns the unique identity of the asset
func (r *ChannelReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *ChannelReference) Variable() bool {
	return false
}

func (r *ChannelReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ UUIDReference = (*ChannelReference)(nil)

// ClassifierReference is used to reference a classifier
type ClassifierReference struct {
	UUID ClassifierUUID `json:"uuid" validate:"required,uuid"`
	Name string         `json:"name"`
}

// NewClassifierReference creates a new classifier reference with the given UUID and name
func NewClassifierReference(uuid ClassifierUUID, name string) *ClassifierReference {
	return &ClassifierReference{UUID: uuid, Name: name}
}

// Type returns the name of the asset type
func (r *ClassifierReference) Type() string {
	return "classifier"
}

// GenericUUID returns the untyped UUID
func (r *ClassifierReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity returns the unique identity of the asset
func (r *ClassifierReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *ClassifierReference) Variable() bool {
	return false
}

func (r *ClassifierReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ UUIDReference = (*ClassifierReference)(nil)

// GroupReference is used to reference a group
type GroupReference struct {
	UUID      GroupUUID `json:"uuid,omitempty" validate:"omitempty,uuid4"`
	Name      string    `json:"name,omitempty"`
	NameMatch string    `json:"name_match,omitempty" engine:"evaluated"`
}

// NewGroupReference creates a new group reference with the given UUID and name
func NewGroupReference(uuid GroupUUID, name string) *GroupReference {
	return &GroupReference{UUID: uuid, Name: name}
}

// NewVariableGroupReference creates a new group reference from the given templatized name match
func NewVariableGroupReference(nameMatch string) *GroupReference {
	return &GroupReference{NameMatch: nameMatch}
}

// Type returns the name of the asset type
func (r *GroupReference) Type() string {
	return "group"
}

// GenericUUID returns the untyped UUID
func (r *GroupReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity returns the unique identity of the asset
func (r *GroupReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *GroupReference) Variable() bool {
	return r.Identity() == ""
}

func (r *GroupReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ UUIDReference = (*GroupReference)(nil)

// FieldReference is a reference to a field
type FieldReference struct {
	Key  string `json:"key" validate:"required"`
	Name string `json:"name"`
}

// NewFieldReference creates a new field reference with the given key and name
func NewFieldReference(key string, name string) *FieldReference {
	return &FieldReference{Key: key, Name: name}
}

// Type returns the name of the asset type
func (r *FieldReference) Type() string {
	return "field"
}

// Identity returns the unique identity of the asset
func (r *FieldReference) Identity() string {
	return string(r.Key)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *FieldReference) Variable() bool {
	return false
}

func (r *FieldReference) String() string {
	return fmt.Sprintf("%s[key=%s,name=%s]", r.Type(), r.Identity(), r.Name)
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

// Type returns the name of the asset type
func (r *FlowReference) Type() string {
	return "flow"
}

// GenericUUID returns the untyped UUID
func (r *FlowReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity returns the unique identity of the asset
func (r *FlowReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *FlowReference) Variable() bool {
	return false
}

func (r *FlowReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ UUIDReference = (*FlowReference)(nil)

// GlobalReference is a reference to a global
type GlobalReference struct {
	Key  string `json:"key" validate:"required"`
	Name string `json:"name"`
}

// NewGlobalReference creates a new global reference with the given key and name
func NewGlobalReference(key string, name string) *GlobalReference {
	return &GlobalReference{Key: key, Name: name}
}

// Type returns the name of the asset type
func (r *GlobalReference) Type() string {
	return "global"
}

// Identity returns the unique identity of the asset
func (r *GlobalReference) Identity() string {
	return string(r.Key)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *GlobalReference) Variable() bool {
	return false
}

func (r *GlobalReference) String() string {
	return fmt.Sprintf("%s[key=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ Reference = (*GlobalReference)(nil)

// LabelReference is used to reference a label
type LabelReference struct {
	UUID      LabelUUID `json:"uuid,omitempty" validate:"omitempty,uuid4"`
	Name      string    `json:"name,omitempty"`
	NameMatch string    `json:"name_match,omitempty" engine:"evaluated"`
}

// NewLabelReference creates a new label reference with the given UUID and name
func NewLabelReference(uuid LabelUUID, name string) *LabelReference {
	return &LabelReference{UUID: uuid, Name: name}
}

// NewVariableLabelReference creates a new label reference from the given templatized name match
func NewVariableLabelReference(nameMatch string) *LabelReference {
	return &LabelReference{NameMatch: nameMatch}
}

// Type returns the name of the asset type
func (r *LabelReference) Type() string {
	return "label"
}

// GenericUUID returns the untyped UUID
func (r *LabelReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity returns the unique identity of the asset
func (r *LabelReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *LabelReference) Variable() bool {
	return r.Identity() == ""
}

func (r *LabelReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ UUIDReference = (*LabelReference)(nil)

// TemplateReference is used to reference a Template
type TemplateReference struct {
	UUID TemplateUUID `json:"uuid" validate:"required,uuid"`
	Name string       `json:"name"`
}

// NewTemplateReference creates a new template reference with the given UUID and name
func NewTemplateReference(uuid TemplateUUID, name string) *TemplateReference {
	return &TemplateReference{UUID: uuid, Name: name}
}

// GenericUUID returns the untyped UUID
func (r *TemplateReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity returns the unique identity of the asset
func (r *TemplateReference) Identity() string {
	return string(r.UUID)
}

// Type returns the name of the asset type
func (r *TemplateReference) Type() string {
	return "template"
}

func (r *TemplateReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *TemplateReference) Variable() bool {
	return false
}

var _ UUIDReference = (*TemplateReference)(nil)

// TicketerReference is used to reference a ticketer
type TicketerReference struct {
	UUID TicketerUUID `json:"uuid" validate:"required,uuid"`
	Name string       `json:"name"`
}

// NewTicketerReference creates a new classifier reference with the given UUID and name
func NewTicketerReference(uuid TicketerUUID, name string) *TicketerReference {
	return &TicketerReference{UUID: uuid, Name: name}
}

// Type returns the name of the asset type
func (r *TicketerReference) Type() string {
	return "ticketer"
}

// GenericUUID returns the untyped UUID
func (r *TicketerReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity returns the unique identity of the asset
func (r *TicketerReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *TicketerReference) Variable() bool {
	return false
}

func (r *TicketerReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ UUIDReference = (*TicketerReference)(nil)

//------------------------------------------------------------------------------------------
// Callbacks for missing assets
//------------------------------------------------------------------------------------------

// MissingCallback is callback to be invoked when an asset is missing
type MissingCallback func(Reference, error)

// PanicOnMissing panics if an asset is reported missing
var PanicOnMissing MissingCallback = func(a Reference, err error) { panic(fmt.Sprintf("missing asset: %s, due to: %s", a, err)) }

// IgnoreMissing does nothing if an asset is reported missing
var IgnoreMissing MissingCallback = func(Reference, error) {}

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

// TypedReference is a utility struct for when we need to serialize a reference with a type
type TypedReference struct {
	Reference Reference `json:"-"`
	Type      string    `json:"type"`
}

// NewTypedReference creates a new typed reference
func NewTypedReference(r Reference) TypedReference {
	return TypedReference{Reference: r, Type: r.Type()}
}

func (r TypedReference) MarshalJSON() ([]byte, error) {
	type typed TypedReference // need to alias type to avoid circular calls to this method
	return jsonx.MarshalMerged(r.Reference, typed(r))
}
