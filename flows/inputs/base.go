package inputs

import (
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
)

type baseInput struct {
	uuid      flows.InputUUID
	channel   flows.Channel
	createdOn time.Time
}

func (i *baseInput) UUID() flows.InputUUID  { return i.uuid }
func (i *baseInput) Channel() flows.Channel { return i.channel }
func (i *baseInput) CreatedOn() time.Time   { return i.createdOn }

// Resolve resolves the given key when this input is referenced in an expression
func (i *baseInput) Resolve(key string) types.XValue {
	switch key {
	case "uuid":
		return types.NewXString(string(i.uuid))
	case "created_on":
		return types.NewXDate(i.createdOn)
	case "channel":
		return i.channel
	}

	return types.NewXResolveError(i, key)
}
