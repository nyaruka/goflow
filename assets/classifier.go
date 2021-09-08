package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
)

// ClassifierUUID is the UUID of an NLU classifier
type ClassifierUUID uuids.UUID

// Classifier is an NLU classifier.
//
//   {
//     "uuid": "37657cf7-5eab-4286-9cb0-bbf270587bad",
//     "name": "Booking",
//     "type": "wit",
//     "intents": ["book_flight", "book_hotel"]
//   }
//
// @asset classifier
type Classifier interface {
	UUID() ClassifierUUID
	Name() string
	Type() string
	Intents() []string
}

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
