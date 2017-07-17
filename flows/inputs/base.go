package inputs

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
)

type BaseInput struct {
	channelUUID flows.ChannelUUID
	createdOn   time.Time
}

func (i *BaseInput) ChannelUUID() flows.ChannelUUID               { return i.channelUUID }
func (i *BaseInput) SetChannelUUID(channelUUID flows.ChannelUUID) { i.channelUUID = channelUUID }

func (i *BaseInput) CreatedOn() time.Time        { return i.createdOn }
func (i *BaseInput) SetCreatedOn(time time.Time) { i.createdOn = time }

// Resolve resolves the passed in key to a value, returning an error if the key is unknown
func (i *BaseInput) Resolve(key string) interface{} {
	switch key {
	case "time":
		return i.createdOn
	case "channel_uuid":
		return i.channelUUID
	}
	return fmt.Errorf("No such field '%s' on input", key)
}
