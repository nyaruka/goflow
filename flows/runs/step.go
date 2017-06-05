package runs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

type step struct {
	uuid      flows.StepUUID
	node      flows.NodeUUID
	exit      flows.ExitUUID
	arrivedOn time.Time
	leftOn    *time.Time
	events    []flows.Event
}

func (s *step) UUID() flows.StepUUID  { return s.uuid }
func (s *step) Node() flows.NodeUUID  { return s.node }
func (s *step) Exit() flows.ExitUUID  { return s.exit }
func (s *step) ArrivedOn() time.Time  { return s.arrivedOn }
func (s *step) LeftOn() *time.Time    { return s.leftOn }
func (s *step) Events() []flows.Event { return s.events }

func (s *step) Resolve(key string) interface{} {
	switch key {

	case "uuid":
		return s.UUID

	case "node":
		return s.Node

	case "exit":
		return s.Exit

	case "arrived_on":
		return s.ArrivedOn

	case "left_on":
		return s.LeftOn
	}

	return fmt.Errorf("No field '%s' on step", key)
}

func (s *step) Default() interface{} {
	return s
}

func (s *step) String() string {
	return string(s.node)
}

func (s *step) Leave(exit flows.ExitUUID) {
	now := time.Now().In(time.UTC)
	s.exit = exit
	s.leftOn = &now
}

func (s *step) addEvent(e flows.Event) {
	e.SetCreatedOn(time.Now().In(time.UTC))
	s.events = append(s.events, e)
}

func (s *step) addError(err error) {
	s.addEvent(&events.ErrorEvent{Text: err.Error()})
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type stepEnvelope struct {
	UUID      flows.StepUUID         `json:"uuid"`
	Node      flows.NodeUUID         `json:"node_uuid"`
	Exit      flows.ExitUUID         `json:"exit,omitempty"`
	ArrivedOn time.Time              `json:"arrived_on"`
	LeftOn    *time.Time             `json:"left_on,omitempty"`
	Events    []*utils.TypedEnvelope `json:"events,omitempty"`
}

func (s *step) UnmarshalJSON(data []byte) error {
	var se stepEnvelope
	var err error

	err = json.Unmarshal(data, &se)
	if err != nil {
		return err
	}

	s.uuid = se.UUID
	s.node = se.Node
	s.exit = se.Exit
	s.arrivedOn = se.ArrivedOn
	s.leftOn = se.LeftOn

	s.events = make([]flows.Event, len(se.Events))
	for i := range s.events {
		s.events[i], err = events.EventFromEnvelope(se.Events[i])
		if err != nil {
			return err
		}
	}

	return err
}

func (s *step) MarshalJSON() ([]byte, error) {
	var se stepEnvelope

	se.UUID = s.uuid
	se.Node = s.node
	se.Exit = s.exit
	se.ArrivedOn = s.arrivedOn
	se.LeftOn = s.leftOn

	se.Events = make([]*utils.TypedEnvelope, len(s.events))
	for i, event := range s.events {
		eventData, err := json.Marshal(event)
		if err != nil {
			return nil, err
		}
		se.Events[i] = &utils.TypedEnvelope{event.Type(), eventData}
	}

	return json.Marshal(se)
}
