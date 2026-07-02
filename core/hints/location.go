package hints

func init() {
	registerType(TypeLocation, func() Hint { return &Location{} })
}

// TypeLocation is the type of our location hint
const TypeLocation string = "location"

// Location requests a message with a location attachment, i.e. a geo:<lat>,<long>
type Location struct {
	baseHint
}

// NewLocation creates a new location hint
func NewLocation() *Location {
	return &Location{
		baseHint: newBaseHint(TypeLocation),
	}
}
