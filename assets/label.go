package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/utils"
	validator "gopkg.in/go-playground/validator.v9"
)

func init() {
	utils.RegisterStructValidator(LabelReferenceValidation, LabelReference{})
}

// LabelUUID is the UUID of a label
type LabelUUID uuids.UUID

// Label is an organizational tag that can be applied to a message.
//
//   {
//     "uuid": "14782905-81a6-4910-bc9f-93ad287b23c3",
//     "name": "Spam"
//   }
//
// @asset label
type Label interface {
	UUID() LabelUUID
	Name() string
}

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

//------------------------------------------------------------------------------------------
// Validation
//------------------------------------------------------------------------------------------

// LabelReferenceValidation validates that the given label reference is either a concrete
// reference or a name matcher
func LabelReferenceValidation(sl validator.StructLevel) {
	ref := sl.Current().Interface().(LabelReference)
	if neitherOrBoth(string(ref.UUID), ref.NameMatch) {
		sl.ReportError(ref.UUID, "uuid", "UUID", "mutually_exclusive", "name_match")
		sl.ReportError(ref.NameMatch, "name_match", "NameMatch", "mutually_exclusive", "uuid")
	}
}
