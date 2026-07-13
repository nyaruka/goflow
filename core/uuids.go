package core

import "github.com/nyaruka/gocommon/uuids"

// NodeUUID is a UUID of a flow node
type NodeUUID uuids.UUID

// NewNodeUUID generates a new UUID for a node
func NewNodeUUID() NodeUUID { return NodeUUID(uuids.NewV4()) }

// InputUUID is the UUID of an input
type InputUUID uuids.UUID
