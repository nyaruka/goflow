package assets

import (
	"github.com/nyaruka/goflow/utils"
)

// GroupUUID is the UUID of a group
type GroupUUID utils.UUID

// Group is a set of contacts that
type Group interface {
	UUID() GroupUUID
	Name() string
	Query() string
}

// LabelUUID is the UUID of a label
type LabelUUID utils.UUID

// Label is something that can be applied a message
type Label interface {
	UUID() LabelUUID
	Name() string
}

// AssetSource is a source of assets
type AssetSource interface {
	Groups() ([]Group, error)
	Labels() ([]Label, error)
}
