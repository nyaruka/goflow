package store

import (
	"fmt"
	"strings"

	"database/sql"

	"github.com/jmoiron/sqlx"
)

const (
	SchemeTel     = "tel"
	SchemeTwitter = "twitter"
	SchemeEmail   = "email"
)

type ContactURNID struct {
	sql.NullInt64
}
type URN string

type ContactURN struct {
	Org      OrgID        `db:"org_id"`
	ID       ContactURNID `db:"id"`
	URN      URN          `db:"urn"`
	Scheme   string       `db:"scheme"`
	Path     string       `db:"path"`
	Priority int          `db:"priority"`
	Channel  ChannelID    `db:"channel_id"`
	Contact  ContactID    `db:"contact_id"`
}

const insertURN = `
INSERT INTO contacts_contacturn(org_id, urn, path, scheme, priority, channel_id, contact_id)
VALUES(:org_id, :urn, :path, :scheme, :priority, :channel_id, :contact_id)
RETURNING id
`

const updateURN = `
UPDATE contacts_contacturn
SET channel_id = :channel_id, contact_id = :contact_id
WHERE id = :id
`

const selectOrgURN = `
SELECT org_id, id, urn, scheme, path, priority, channel_id, contact_id 
FROM contacts_contacturn
WHERE org_id = $1 AND urn = $2
ORDER BY priority desc LIMIT 1
`

// NewContactURN returns a new ContactURN object for the passed in org, contact and string urn, this is not saved to the DB yet
func NewContactURN(org OrgID, channel ChannelID, contact ContactID, urn URN) *ContactURN {
	offset := strings.Index(string(urn), ":")
	scheme := string(urn)[:offset]
	path := string(urn)[offset+1:]

	return &ContactURN{Org: org, Channel: channel, Contact: contact, URN: urn, Scheme: scheme, Path: path}
}

// ContactURNForURN returns the ContactURN for the passed in org and URN, creating and associating
//  it with the passed in contact if necessary
func ContactURNForURN(db *sqlx.DB, org OrgID, channel ChannelID, contact ContactID, urn URN) (*ContactURN, error) {
	contactURN := NewContactURN(org, channel, contact, urn)
	err := db.Get(contactURN, selectOrgURN, org, urn)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// we didn't find it, let's insert it
	if err == sql.ErrNoRows {
		err = InsertContactURN(db, contactURN)
		if err != nil {
			return nil, err
		}
	}

	// make sure our contact URN is up to date
	if contactURN.Channel != channel || contactURN.Contact != contact {
		contactURN.Channel = channel
		contactURN.Contact = contact

		err = UpdateContactURN(db, contactURN)
	}

	return contactURN, nil
}

// InsertContactURN inserts the passed in urn, the id field will be populated with the result on success
func InsertContactURN(db *sqlx.DB, urn *ContactURN) error {
	_, err := db.NamedExec(insertURN, urn)
	return err
}

// UpdateContactURN updates the Channel and Contact on an existing URN
func UpdateContactURN(db *sqlx.DB, urn *ContactURN) error {
	rows, err := db.NamedQuery(updateURN, urn)
	if err != nil {
		return err
	}
	if rows.Next() {
		rows.Scan(&urn.ID)
	}
	return err
}

// NewURNFromParts builds a new URN from the passed in string and path
func NewURNFromParts(scheme string, path string) URN {
	// do some simple normalization when appropriate
	if scheme == SchemeTel {
		// TODO: do real normalization here
		path = strings.ToLower(path)
	} else if scheme == SchemeEmail {
		path = strings.ToLower(path)
	} else if scheme == SchemeTwitter {
		if strings.HasSuffix(path, "@") {
			path = path[1:]
		}
	}

	return URN(fmt.Sprintf("%s:%s", scheme, path))
}
