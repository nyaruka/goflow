package engine

import (
	"github.com/nyaruka/goflow/flows"
)

type sprint struct {
	modifiers []flows.Modifier
	events    []flows.Event
}

// NewEmptySprint creates a new sprint
func NewEmptySprint() flows.Sprint {
	return &sprint{
		modifiers: make([]flows.Modifier, 0),
		events:    make([]flows.Event, 0),
	}
}

// NewSprint creates a new sprint with the passed in modifiers and events
func NewSprint(modifiers []flows.Modifier, events []flows.Event) flows.Sprint {
	return &sprint{
		modifiers: modifiers,
		events:    events,
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
