package assets

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
)

// FlowUUID is the UUID of a flow
type FlowUUID uuids.UUID

// Flow is graph of nodes with actions and routers.
//
//   {
//     "uuid": "14782905-81a6-4910-bc9f-93ad287b23c3",
//     "name": "Registration",
//     "definition": {
//       "nodes": []
//     }
//   }
//
// @asset flow
type Flow interface {
	UUID() FlowUUID
	Name() string
	Definition() json.RawMessage
}

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
