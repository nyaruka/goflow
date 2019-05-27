package utils

import (
	"fmt"
	"regexp"

	"github.com/gofrs/uuid"
)

var UUID4Regex = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}`)
var UUID4OnlyRegex = regexp.MustCompile(`^` + UUID4Regex.String() + `$`)

// IsUUIDv4 returns whether the given string contains only a valid v4 UUID
func IsUUIDv4(s string) bool {
	return UUID4OnlyRegex.MatchString(s)
}

// UUID is a 36 character UUID
type UUID string

// UUIDGenerator is something that can generate a UUID
type UUIDGenerator interface {
	Next() UUID
}

// defaultUUID4Generator generates a random v4 UUID using crypto/rand
type defaultUUID4Generator struct{}

// Next returns the next random UUID
func (g defaultUUID4Generator) Next() UUID {
	u, err := uuid.NewV4()
	if err != nil {
		// if we can't generate a UUID.. we're done
		panic(fmt.Sprintf("unable to generate UUID: %s", err))
	}
	return UUID(u.String())
}

// DefaultUUIDGenerator is the default generator for calls to NewUUID
var DefaultUUIDGenerator UUIDGenerator = defaultUUID4Generator{}
var currentUUIDGenerator = DefaultUUIDGenerator

// NewUUID returns a new NewUUID
func NewUUID() UUID {
	return currentUUIDGenerator.Next()
}

// SetUUIDGenerator sets the generator used by UUID4()
func SetUUIDGenerator(generator UUIDGenerator) {
	currentUUIDGenerator = generator
}
