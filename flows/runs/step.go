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
	stepUUID  flows.StepUUID
	nodeUUID  flows.NodeUUID
	exitUUID  flows.ExitUUID
	arrivedOn time.Time
	leftOn    *time.Time
	events    []flows.Event
}

func (s *step) UUID() flows.StepUUID     { return s.stepUUID }
func (s *step) NodeUUID() flows.NodeUUID { return s.nodeUUID }
func (s *step) ExitUUID() flows.ExitUUID { return s.exitUUID }
func (s *step) ArrivedOn() time.Time     { return s.arrivedOn }
func (s *step) LeftOn() *time.Time       { return s.leftOn }
func (s *step) Events() []flows.Event    { return s.events }

func (s *step) Resolve(key string) interface{} {
	switch key {
	case "uuid":
		return s.UUID
	case "node_uuid":
		return s.NodeUUID
	case "exit_uuid":
		return s.ExitUUID
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

var _ utils.VariableResolver = (*step)(nil)

func (s *step) String() string {
	return string(s.nodeUUID)
}

func (s *step) Leave(exit flows.ExitUUID) {
	now := time.Now().UTC()
	s.exitUUID = exit
	s.leftOn = &now
}

func (s *step) addEvent(e flows.Event) {
	e.SetCreatedOn(time.Now().UTC())
	s.events = append(s.events, e)
}

func (s *step) addError(err error) {
	s.addEvent(&events.ErrorEvent{Text: err.Error()})
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type stepEnvelope struct {
	UUID      flows.StepUUID         `json:"uuid" validate:"required,uuid4"`
	NodeUUID  flows.NodeUUID         `json:"node_uuid" validate:"required,uuid4"`
	ExitUUID  flows.ExitUUID         `json:"exit_uuid,omitempty" validate:"omitempty,uuid4"`
	ArrivedOn time.Time              `json:"arrived_on"`
	LeftOn    *time.Time             `json:"left_on,omitempty"`
	Events    []*utils.TypedEnvelope `json:"events,omitempty" validate:"omitempty,dive"`
}

func (s *step) UnmarshalJSON(data []byte) error {
	var se stepEnvelope
	var err error

	err = json.Unmarshal(data, &se)
	if err != nil {
		return err
	}

	s.stepUUID = se.UUID
	s.nodeUUID = se.NodeUUID
	s.exitUUID = se.ExitUUID
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

	se.UUID = s.stepUUID
	se.NodeUUID = s.nodeUUID
	se.ExitUUID = s.exitUUID
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
