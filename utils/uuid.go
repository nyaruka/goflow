package utils

import (
	"github.com/satori/go.uuid"
)

// UUIDGenerator is something that can generate a UUID
type UUIDGenerator interface {
	Next() string
}

// RandomUUID4Generator generates a random v4 UUID
type RandomUUID4Generator struct{}

// Next returns the next random UUID
func (g RandomUUID4Generator) Next() string {
	return uuid.NewV4().String()
}

var currentUUIDGenerator UUIDGenerator = RandomUUID4Generator{}

// UUID returns a new UUID
func UUID() string {
	return currentUUIDGenerator.Next()
}

// SetUUIDGenerator sets the generator used by UUID4()
func SetUUIDGenerator(generator UUIDGenerator) {
	currentUUIDGenerator = generator
}
