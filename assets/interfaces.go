package assets

import (
	"encoding/json"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils/uuids"
)

// ChannelUUID is the UUID of a channel
type ChannelUUID uuids.UUID

// ChannelRole is a role that a channel can perform
type ChannelRole string

// different roles that channels can perform
const (
	ChannelRoleSend    ChannelRole = "send"
	ChannelRoleReceive ChannelRole = "receive"
	ChannelRoleCall    ChannelRole = "call"
	ChannelRoleAnswer  ChannelRole = "answer"
	ChannelRoleUSSD    ChannelRole = "ussd"
)

// Channel is something that can send/receive messages.
//
//   {
//     "uuid": "14782905-81a6-4910-bc9f-93ad287b23c3",
//     "name": "My Android",
//     "address": "+593979011111",
//     "schemes": ["tel"],
//     "roles": ["send", "receive"],
//     "country": "EC"
//   }
//
// @asset channel
type Channel interface {
	UUID() ChannelUUID
	Name() string
	Address() string
	Schemes() []string
	Roles() []ChannelRole
	Parent() *ChannelReference
	Country() envs.Country
	MatchPrefixes() []string
	AllowInternational() bool
}

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

// FieldUUID is the UUID of a field
type FieldUUID uuids.UUID

// FieldType is the data type of values for each field
type FieldType string

// field value types
const (
	FieldTypeText     FieldType = "text"
	FieldTypeNumber   FieldType = "number"
	FieldTypeDatetime FieldType = "datetime"
	FieldTypeWard     FieldType = "ward"
	FieldTypeDistrict FieldType = "district"
	FieldTypeState    FieldType = "state"
)

// Field is a custom contact property.
//
//   {
//     "uuid": "d66a7823-eada-40e5-9a3a-57239d4690bf",
//     "key": "gender",
//     "name": "Gender",
//     "type": "text"
//   }
//
// @asset field
type Field interface {
	UUID() FieldUUID
	Key() string
	Name() string
	Type() FieldType
}

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

// Global is a named constant.
//
//   {
//     "key": "organization_name",
//     "name": "Organization Name",
//     "value": "U-Report"
//   }
//
// @asset global
type Global interface {
	Key() string
	Name() string
	Value() string
}

// GroupUUID is the UUID of a group
type GroupUUID uuids.UUID

// Group is a set of contacts which can be static or dynamic (i.e. based on a query).
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

// LocationHierarchy is a searchable hierarchy of locations.
//
//   {
//     "name": "Rwanda",
//     "aliases": ["Ruanda"],
//     "children": [
//       {
//         "name": "Kigali City",
//         "aliases": ["Kigali", "Kigari"],
//         "children": [
//           {
//             "name": "Gasabo",
//             "children": [
//               {
//                 "id": "575743222",
//                 "name": "Gisozi"
//               },
//               {
//                 "id": "457378732",
//                 "name": "Ndera"
//               }
//             ]
//           },
//           {
//             "name": "Nyarugenge",
//             "children": []
//           }
//         ]
//       },
//       {
//         "name": "Eastern Province"
//       }
//     ]
//   }
//
// @asset location
type LocationHierarchy interface {
	FindByPath(path envs.LocationPath) *envs.Location
	FindByName(name string, level envs.LocationLevel, parent *envs.Location) []*envs.Location
}

// Resthook is a set of URLs which are subscribed to the named event.
//
//   {
//     "slug": "new-registration",
//     "subscribers": [
//       "http://example.com/record.php?@contact.uuid"
//     ]
//   }
//
// @asset resthook
type Resthook interface {
	Slug() string
	Subscribers() []string
}

// TemplateUUID is the UUID of a template
type TemplateUUID uuids.UUID

// Template is a message template, currently only used by WhatsApp channels
//
//   {
//     "name": "revive-issue",
//     "uuid": "14782905-81a6-4910-bc9f-93ad287b23c3",
//     "translations": [
//       {
//          "language": "eng",
//          "content": "Hi {{1}}, are you still experiencing your issue?",
//          "channel": {
//            "uuid": "cf26be4c-875f-4094-9e08-162c3c9dcb5b",
//            "name": "Twilio Channel"
//          }
//       },
//       {
//          "language": "fra",
//          "content": "Bonjour {{1}}",
//          "channel": {
//            "uuid": "cf26be4c-875f-4094-9e08-162c3c9dcb5b",
//            "name": "Twilio Channel"
//          }
//       }
//     ]
//   }
//
// @asset template
type Template interface {
	UUID() TemplateUUID
	Name() string
	Translations() []TemplateTranslation
}

// TemplateTranslation represents a single translation for a specific template and channel
type TemplateTranslation interface {
	Content() string
	Language() envs.Language
	Country() envs.Country
	VariableCount() int
	Channel() ChannelReference
}

// TicketerUUID is the UUID of a ticketer
type TicketerUUID uuids.UUID

// Ticketer is a system which can open or close tickets
//
//   {
//     "uuid": "37657cf7-5eab-4286-9cb0-bbf270587bad",
//     "name": "Support Tickets",
//     "type": "mailgun"
//   }
//
// @asset ticketer
type Ticketer interface {
	UUID() TicketerUUID
	Name() string
	Type() string
}

// Source is a source of assets
type Source interface {
	Channels() ([]Channel, error)
	Classifiers() ([]Classifier, error)
	Fields() ([]Field, error)
	Flow(FlowUUID) (Flow, error)
	Globals() ([]Global, error)
	Groups() ([]Group, error)
	Labels() ([]Label, error)
	Locations() ([]LocationHierarchy, error)
	Resthooks() ([]Resthook, error)
	Templates() ([]Template, error)
	Ticketers() ([]Ticketer, error)
}
