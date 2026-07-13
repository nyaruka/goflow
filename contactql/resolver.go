package contactql

import (
	"github.com/nyaruka/goflow/assets"
)

// Fixed attributes that can be searched
const (
	AttributeUUID       = "uuid"
	AttributeID         = "id" // deprecated in favor of ref
	AttributeRef        = "ref"
	AttributeName       = "name"
	AttributeStatus     = "status"
	AttributeLanguage   = "language"
	AttributeURN        = "urn"
	AttributeGroup      = "group"
	AttributeFlow       = "flow"
	AttributeHistory    = "history"
	AttributeTickets    = "tickets"
	AttributeCreatedOn  = "created_on"
	AttributeLastSeenOn = "last_seen_on"
)

var attributes = map[string]assets.FieldType{
	AttributeUUID:       assets.FieldTypeText,
	AttributeID:         assets.FieldTypeText,
	AttributeRef:        assets.FieldTypeText,
	AttributeName:       assets.FieldTypeText,
	AttributeStatus:     assets.FieldTypeText,
	AttributeLanguage:   assets.FieldTypeText,
	AttributeURN:        assets.FieldTypeText,
	AttributeGroup:      assets.FieldTypeText,
	AttributeFlow:       assets.FieldTypeText,
	AttributeHistory:    assets.FieldTypeText,
	AttributeTickets:    assets.FieldTypeNumber,
	AttributeCreatedOn:  assets.FieldTypeDatetime,
	AttributeLastSeenOn: assets.FieldTypeDatetime,
}

// AttributeType returns the value type of the given fixed attribute, or false if it's not a fixed attribute
func AttributeType(key string) (assets.FieldType, bool) {
	t, ok := attributes[key]
	return t, ok
}

// Resolver provides functions for resolving assets referenced in queries
type Resolver interface {
	ResolveField(key string) assets.Field
	ResolveGroup(name string) assets.Group
	ResolveFlow(name string) assets.Flow
}
