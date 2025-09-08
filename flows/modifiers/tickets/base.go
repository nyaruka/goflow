package tickets

// base of all modifier types
type baseModifier struct {
	typ string
}

// creates new base modifier
func newBaseModifier(typ string) baseModifier {
	return baseModifier{typ: typ}
}

// Type returns the type of this modifier
func (m *baseModifier) Type() string { return m.typ }
