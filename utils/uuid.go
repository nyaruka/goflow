package utils

import (
	"math/rand"

	"github.com/satori/go.uuid"
)

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
	return UUID(uuid.NewV4().String())
}

// SeededUUID4Generator generates a seedable random v4 UUID using math/rand
type SeededUUID4Generator struct {
	rnd *rand.Rand
}

// NewSeededUUID4Generator creates a new SeededUUID4Generator from the given seed
func NewSeededUUID4Generator(seed int64) *SeededUUID4Generator {
	return &SeededUUID4Generator{rnd: rand.New(rand.NewSource(seed))}
}

// Next returns the next random UUID
func (g *SeededUUID4Generator) Next() UUID {
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
