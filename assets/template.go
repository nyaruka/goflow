package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/uuids"
)

// TemplateUUID is the UUID of a template
type TemplateUUID uuids.UUID

// Template is a message template, currently only used by WhatsApp channels
//
//	{
//	  "name": "greeting",
//	  "uuid": "14782905-81a6-4910-bc9f-93ad287b23c3",
//	  "translations": [
//	    {
//	       "locale": "eng-US",
//	       "channel": {
//	         "uuid": "cf26be4c-875f-4094-9e08-162c3c9dcb5b",
//	         "name": "Twilio Channel"
//	       },
//	       "components": [
//	         {
//	           "name": "body",
//	           "type": "body/text",
//	           "content": "Hello {{1}}",
//	           "variables": {"1": 0}
//	         }
//	       ],
//	       "variables": [{"type": "text"}]
//	    },
//	    {
//	       "locale": "fra",
//	       "channel": {
//	         "uuid": "cf26be4c-875f-4094-9e08-162c3c9dcb5b",
//	         "name": "Twilio Channel"
//	       },
//	       "components": [
//	         {
//	           "name": "body",
//	           "type": "body/text",
//	           "content": "Bonjour {{1}}",
//	           "variables": {"1": 0}
//	         }
//	       ],
//	       "variables": [{"type": "text"}]
//	    }
//	  ]
//	}
//
// @asset template
type Template interface {
	UUID() TemplateUUID
	Name() string
	Translations() []TemplateTranslation
}

type TemplateVariable interface {
	Type() string
}

type TemplateComponent interface {
	Name() string
	Type() string
	Content() string
	Display() string
	Variables() map[string]int
}

// TemplateTranslation represents a single translation for a specific template and channel
type TemplateTranslation interface {
	Locale() i18n.Locale
	Channel() *ChannelReference
	Components() []TemplateComponent
	Variables() []TemplateVariable
}

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
