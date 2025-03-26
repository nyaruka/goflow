package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
)

// LLMUUID is the UUID of an LLM
type LLMUUID uuids.UUID

// LLM is a large language model.
//
//	{
//	  "uuid": "00cc7310-4bb9-473f-851e-39b0880aad78",
//	  "name": "ChatGPT-4",
//	  "type": "openai"
//	}
//
// @asset llm
type LLM interface {
	UUID() LLMUUID
	Name() string
	Type() string
}

// LLMReference is used to reference an LLM
type LLMReference struct {
	UUID LLMUUID `json:"uuid" validate:"required,uuid"`
	Name string  `json:"name"`
}

// NewLLMReference creates a new LLM reference with the given UUID and name
func NewLLMReference(uuid LLMUUID, name string) *LLMReference {
	return &LLMReference{UUID: uuid, Name: name}
}

// Type returns the name of the asset type
func (r *LLMReference) Type() string {
	return "llm"
}

// GenericUUID returns the untyped UUID
func (r *LLMReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity returns the unique identity of the asset
func (r *LLMReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *LLMReference) Variable() bool {
	return false
}

func (r *LLMReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ UUIDReference = (*LLMReference)(nil)
