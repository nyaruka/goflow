package engine

import (
	"github.com/nyaruka/goflow/flows"
)

type sprint struct {
	modifiers []flows.Modifier
	events    []flows.Event
}

// NewSprint creates a new sprint
func NewSprint() flows.Sprint {
	return &sprint{
		modifiers: make([]flows.Modifier, 0),
		events:    make([]flows.Event, 0),
	}
}

func (s *sprint) Modifiers() []flows.Modifier { return s.modifiers }
func (s *sprint) Events() []flows.Event       { return s.events }

func (s *sprint) LogModifier(m flows.Modifier) {
	s.modifiers = append(s.modifiers, m)
}

func (s *sprint) LogEvent(e flows.Event) {
	s.events = append(s.events, e)
}

var _ flows.Sprint = (*sprint)(nil)
