package store

import "database/sql"

type OrgID struct {
	sql.NullInt64
}
