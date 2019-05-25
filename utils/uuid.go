package utils

import (
	"fmt"
	"math/rand"
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

// fixedUUID4Generator returns the same single fixed v4 UUID (useful for testing)
type fixedUUID4Generator struct {
	value UUID
}

// Next returns the next UUID
func (g *fixedUUID4Generator) Next() UUID {
	return g.value
}

// NewFixedUUID4Generator creates a new fixed UUID4 generator
func NewFixedUUID4Generator(value UUID) UUIDGenerator {
	return &fixedUUID4Generator{value}
}

// generates a seedable random v4 UUID using math/rand
type seededUUID4Generator struct {
	rnd *rand.Rand
}

// NewSeededUUID4Generator creates a new seeded UUID4 generator from the given seed
func NewSeededUUID4Generator(seed int64) UUIDGenerator {
	return &seededUUID4Generator{rnd: NewSeededRand(seed)}
}

// Next returns the next random UUID
func (g *seededUUID4Generator) Next() UUID {
	u := uuid.UUID{}
	if _, err := g.rnd.Read(u[:]); err != nil {
		panic(err)
	}
	u.SetVersion(uuid.V4)
	u.SetVariant(uuid.VariantRFC4122)
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
