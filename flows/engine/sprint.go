package engine

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

type segment struct {
	flow        flows.Flow
	node        flows.Node
	exit        flows.Exit
	operand     string
	destination flows.Node
	time        time.Time
}

func (s *segment) Flow() flows.Flow        { return s.flow }
func (s *segment) Node() flows.Node        { return s.node }
func (s *segment) Exit() flows.Exit        { return s.exit }
func (s *segment) Operand() string         { return s.operand }
func (s *segment) Destination() flows.Node { return s.destination }
func (s *segment) Time() time.Time         { return s.time }

type segmentEnvelope struct {
	FlowUUID        assets.FlowUUID `json:"flow_uuid"`
	NodeUUID        flows.NodeUUID  `json:"node_uuid"`
	ExitUUID        flows.ExitUUID  `json:"exit_uuid"`
	Operand         string          `json:"operand,omitempty"`
	DestinationUUID flows.NodeUUID  `json:"destination_uuid,omitempty"`
	Time            time.Time       `json:"time"`
}

func (s *segment) MarshalJSON() ([]byte, error) {
	return json.Marshal(&segmentEnvelope{
		FlowUUID:        s.flow.UUID(),
		NodeUUID:        s.node.UUID(),
		ExitUUID:        s.exit.UUID(),
		Operand:         s.operand,
		DestinationUUID: s.destination.UUID(),
		Time:            s.time,
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

// NewSprint creates a new sprint - engine doesn't use this but we do it when handling surveyor responses
func NewSprint(modifiers []flows.Modifier, events []flows.Event, segments []flows.Segment) flows.Sprint {
	return &sprint{modifiers: modifiers, events: events, segments: segments}
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

func (s *sprint) logSegment(flow flows.Flow, node flows.Node, exit flows.Exit, operand string, dest flows.Node) {
	s.segments = append(s.segments, &segment{
		flow:        flow,
		node:        node,
		exit:        exit,
		operand:     operand,
		destination: dest,
		time:        dates.Now(),
	})
}

var _ flows.Sprint = (*sprint)(nil)
