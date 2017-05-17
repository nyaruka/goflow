package store

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type MessageID int64

type Msg struct {
	Org        OrgID     `db:"org_id"`
	ID         MessageID `db:"id"`
	Direction  string    `db:"direction"`
	Text       string    `db:"text"`
	Priority   int       `db:"priority"`
	Status     string    `db:"status"`
	Visibility string    `db:"visibility"`
	ExternalID string    `db:"external_id"`

	MessageCount int `db:"msg_count"`
	ErrorCount   int `db:"error_count"`

	Channel    ChannelID    `db:"channel_id"`
	Contact    ContactID    `db:"contact_id"`
	ContactURN ContactURNID `db:"contact_urn_id"`

	NextAttempt time.Time `db:"next_attempt"`
	CreatedOn   time.Time `db:"created_on"`
	ModifiedOn  time.Time `db:"modified_on"`
	QueuedOn    time.Time `db:"queued_on"`
	SentOn      time.Time `db:"sent_on"`
}

const insertMsgSQL = `
INSERT INTO msgs_msg(org_id, direction, has_template_error, text, msg_count, error_count, priority, status, visibility, external_id, channel_id, contact_id, contact_urn_id, created_on, modified_on, next_attempt)
VALUES(:org_id, :direction, FALSE, :text, :msg_count, :error_count, :priority, :status, :visibility, :external_id, :channel_id, :contact_id, :contact_urn_id, :created_on, :modified_on, :next_attempt)
RETURNING id
`

// InsertMsg inserts the passed in msg, the id field will be populated with the result on success
func InsertMsg(db *sqlx.DB, msg *Msg) error {
	rows, err := db.NamedQuery(insertMsgSQL, msg)
	if err != nil {
		return err
	}
	if rows.Next() {
		rows.Scan(&msg.ID)
	}
	return err
}

// NewMsg creates a new contact with the passed in parameters
func NewMsg(org OrgID, channel ChannelID, contact ContactID, urn ContactURNID, direction string, text string, status string, externalID string) *Msg {
	msg := Msg{Org: org, Channel: channel, Contact: contact, ContactURN: urn, Direction: direction, Text: text, Status: status, ExternalID: externalID}
	now := time.Now()
	msg.CreatedOn = now
	msg.ModifiedOn = now
	msg.Visibility = "V"
	msg.MessageCount = 1
	return &msg
}
