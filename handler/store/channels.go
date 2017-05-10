package store

import (
	"fmt"
	"strings"

	"database/sql"

	"github.com/jmoiron/sqlx"
)

type ChannelID struct {
	sql.NullInt64
}

type Channel struct {
	Org         OrgID     `db:"org_id"`
	ID          ChannelID `db:"id"`
	UUID        string    `db:"uuid"`
	ChannelType string    `db:"channel_type"`
	Config      string    `db:"config"`
}

const lookupChannelFromUUIDSQL = `
SELECT org_id, id, uuid, channel_type, config 
FROM channels_channel 
WHERE channel_type = $1 AND uuid = $2 AND is_active = true AND org_id IS NOT NULL
`

// ChannelForUUID attempts to look up the channel with the passed in UUID, returning it
func ChannelForUUID(db *sqlx.DB, channelType string, uuid string) (*Channel, error) {
	// lowercase our uuid
	uuid = strings.ToLower(uuid)

	// select just the fields we need
	channel := Channel{ChannelType: channelType, UUID: uuid}
	err := db.Get(&channel, lookupChannelFromUUIDSQL, channelType, uuid)

	// we didn't find a match
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("No channel found with type %s and UUID: %s", channelType, uuid)
	}

	// other error
	if err != nil {
		return nil, err
	}

	// and return it
	return &channel, err
}
