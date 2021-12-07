package engine

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

type segment struct {
	flow        flows.Flow
	exit        flows.Exit
	destination flows.Node
}

func (s *segment) Flow() flows.Flow        { return s.flow }
func (s *segment) Exit() flows.Exit        { return s.exit }
func (s *segment) Destination() flows.Node { return s.destination }

type segmentEnvelope struct {
	FlowUUID        assets.FlowUUID `json:"flow_uuid"`
	ExitUUID        flows.ExitUUID  `json:"exit_uuid"`
	DestinationUUID flows.NodeUUID  `json:"destination_uuid"`
}

func (s *segment) MarshalJSON() ([]byte, error) {
	return json.Marshal(&segmentEnvelope{
		FlowUUID:        s.flow.UUID(),
		ExitUUID:        s.exit.UUID(),
		DestinationUUID: s.destination.UUID(),
	})
}

var _ flows.Segment = (*segment)(nil)

type sprint struct {
	modifiers []flows.Modifier
	events    []flows.Event
	segments  []flows.Segment
}

// creates a new empty sprint
func newEmptySprint() *sprint {
	return &sprint{
		modifiers: make([]flows.Modifier, 0, 10),
		events:    make([]flows.Event, 0, 10),
		segments:  make([]flows.Segment, 0, 10),
	}
}

func (s *sprint) Modifiers() []flows.Modifier { return s.modifiers }
func (s *sprint) Events() []flows.Event       { return s.events }
func (s *sprint) Segments() []flows.Segment   { return s.segments }

func (s *sprint) logModifier(m flows.Modifier) {
	s.modifiers = append(s.modifiers, m)
}

func (s *sprint) logEvent(e flows.Event) {
	s.events = append(s.events, e)
}

func (s *sprint) logSegment(flow flows.Flow, exit flows.Exit, dest flows.Node) {
	s.segments = append(s.segments, &segment{flow: flow, exit: exit, destination: dest})
}

var _ flows.Sprint = (*sprint)(nil)
