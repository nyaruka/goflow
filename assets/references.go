package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
)

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

//------------------------------------------------------------------------------------------
// Callbacks for missing assets
//------------------------------------------------------------------------------------------

// MissingCallback is callback to be invoked when an asset is missing
type MissingCallback func(Reference, error)

// PanicOnMissing panics if an asset is reported missing
var PanicOnMissing MissingCallback = func(a Reference, err error) { panic(fmt.Sprintf("missing asset: %s, due to: %s", a, err)) }

// IgnoreMissing does nothing if an asset is reported missing
var IgnoreMissing MissingCallback = func(Reference, error) {}

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
