package uuids

import (
	"fmt"
	"math/rand"
	"regexp"

	"github.com/nyaruka/gocommon/random"

	"github.com/gofrs/uuid"
)

// V4Regex matches a string containing a valid v4 UUID
var V4Regex = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}`)

// V4OnlyRegex matches a string containing only a valid v4 UUID
var V4OnlyRegex = regexp.MustCompile(`^` + V4Regex.String() + `$`)

// New returns a new v4 UUID
func New() UUID {
	return currentGenerator.Next()
}

// IsV4 returns whether the given string contains only a valid v4 UUID
func IsV4(s string) bool {
	return V4OnlyRegex.MatchString(s)
}

// UUID is a 36 character UUID
type UUID string

// Generator is something that can generate a UUID
type Generator interface {
	Next() UUID
}

// defaultGenerator generates a random v4 UUID using a 3rd party library
type defaultGenerator struct{}

// Next returns the next random UUID
func (g defaultGenerator) Next() UUID {
	u, err := uuid.NewV4()
	if err != nil {
		// if we can't generate a UUID.. we're done
		panic(fmt.Sprintf("unable to generate UUID: %s", err))
	}
	return UUID(u.String())
}

// DefaultGenerator is the default generator for calls to NewUUID
var DefaultGenerator Generator = defaultGenerator{}
var currentGenerator = DefaultGenerator

// SetGenerator sets the generator used by UUID4()
func SetGenerator(generator Generator) {
	currentGenerator = generator
}

// generates a seedable random v4 UUID using math/rand
type seededGenerator struct {
	rnd *rand.Rand
}

// NewSeededGenerator creates a new seeded UUID4 generator from the given seed
func NewSeededGenerator(seed int64) Generator {
	return &seededGenerator{rnd: random.NewSeededGenerator(seed)}
}

// Next returns the next random UUID
func (g *seededGenerator) Next() UUID {
	u := uuid.UUID{}
	if _, err := g.rnd.Read(u[:]); err != nil {
		panic(err)
	}
	u.SetVersion(uuid.V4)
	u.SetVariant(uuid.VariantRFC4122)
	return UUID(u.String())
}
