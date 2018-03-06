package utils

import (
	"github.com/satori/go.uuid"
)

// UUID is a 36 character UUID
type UUID string

// UUIDGenerator is something that can generate a UUID
type UUIDGenerator interface {
	Next() UUID
}

// RandomUUID4Generator generates a random v4 UUID
type RandomUUID4Generator struct{}

// Next returns the next random UUID
func (g RandomUUID4Generator) Next() UUID {
	return UUID(uuid.NewV4().String())
}

var currentUUIDGenerator UUIDGenerator = RandomUUID4Generator{}

// NewUUID returns a new NewUUID
func NewUUID() UUID {
	return currentUUIDGenerator.Next()
}

// SetUUIDGenerator sets the generator used by UUID4()
func SetUUIDGenerator(generator UUIDGenerator) {
	currentUUIDGenerator = generator
}
