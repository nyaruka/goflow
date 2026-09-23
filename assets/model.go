package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
)

// ModelUUID is the UUID of a model
type ModelUUID uuids.UUID

// ModelRole is a role that a model can perform
type ModelRole string

// different roles that models can perform
const (
	ModelRoleEditing ModelRole = "editing"
	ModelRoleEngine  ModelRole = "engine"
)

// Model is an AI model, e.g. a large language model.
//
//	{
//	  "uuid": "00cc7310-4bb9-473f-851e-39b0880aad78",
//	  "name": "ChatGPT-4",
//	  "type": "openai",
//	  "roles": ["editing", "engine"]
//	}
//
// @asset model
type Model interface {
	UUID() ModelUUID
	Name() string
	Type() string
	Roles() []ModelRole
}

// ModelReference is used to reference a model
type ModelReference struct {
	UUID ModelUUID `json:"uuid" validate:"required,uuid"`
	Name string    `json:"name" validate:"max=64"`
}

// NewModelReference creates a new model reference with the given UUID and name
func NewModelReference(uuid ModelUUID, name string) *ModelReference {
	return &ModelReference{UUID: uuid, Name: name}
}

// Type returns the name of the asset type
func (r *ModelReference) Type() string {
	return "llm"
}

// GenericUUID returns the untyped UUID
func (r *ModelReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity returns the unique identity of the asset
func (r *ModelReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *ModelReference) Variable() bool {
	return false
}

func (r *ModelReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ UUIDReference = (*ModelReference)(nil)
