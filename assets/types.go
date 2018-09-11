package assets

import (
	"github.com/nyaruka/goflow/utils"
)

// LabelUUID is the UUID of a label
type LabelUUID utils.UUID

type Label interface {
	UUID() LabelUUID
	Name() string
}

type AssetSource interface {
	Labels() ([]Label, error)
}
