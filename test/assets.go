package test

import (
	"github.com/satori/go.uuid"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/simple"
	"github.com/nyaruka/goflow/flows"
)

func NewGroup(name string, query string) *flows.Group {
	return flows.NewGroup(simple.NewGroup(assets.GroupUUID(uuid.NewV4().String()), name, query))
}
