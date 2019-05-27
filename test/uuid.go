package test

import (
	"math/rand"

	"github.com/nyaruka/goflow/utils"

	"github.com/gofrs/uuid"
)

// generates a seedable random v4 UUID using math/rand
type seededUUIDGenerator struct {
	rnd *rand.Rand
}

// NewSeededUUIDGenerator creates a new seeded UUID4 generator from the given seed
func NewSeededUUIDGenerator(seed int64) utils.UUIDGenerator {
	return &seededUUIDGenerator{rnd: utils.NewSeededRand(seed)}
}

// Next returns the next random UUID
func (g *seededUUIDGenerator) Next() utils.UUID {
	u := uuid.UUID{}
	if _, err := g.rnd.Read(u[:]); err != nil {
		panic(err)
	}
	u.SetVersion(uuid.V4)
	u.SetVariant(uuid.VariantRFC4122)
	return utils.UUID(u.String())
}
