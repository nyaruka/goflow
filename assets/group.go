package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/utils"
	validator "gopkg.in/go-playground/validator.v9"
)

func init() {
	utils.RegisterStructValidator(GroupReferenceValidation, GroupReference{})
}

// GroupUUID is the UUID of a group
type GroupUUID uuids.UUID

// Group is a set of contacts which can be added to and removed from manually, or based on a query.
//
//   {
//     "uuid": "14782905-81a6-4910-bc9f-93ad287b23c3",
//     "name": "Youth",
//     "query": "age <= 18"
//   }
//
// @asset group
type Group interface {
	UUID() GroupUUID
	Name() string
	Query() string
}

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
