package runs

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/goflow/flows"
)

type step struct {
	stepUUID  flows.StepUUID
	nodeUUID  flows.NodeUUID
	exitUUID  flows.ExitUUID
	arrivedOn time.Time
	leftOn    *time.Time
}

func (s *step) UUID() flows.StepUUID     { return s.stepUUID }
func (s *step) NodeUUID() flows.NodeUUID { return s.nodeUUID }
func (s *step) ExitUUID() flows.ExitUUID { return s.exitUUID }
func (s *step) ArrivedOn() time.Time     { return s.arrivedOn }
func (s *step) LeftOn() *time.Time       { return s.leftOn }

func (s *step) Leave(exit flows.ExitUUID) {
	now := time.Now().UTC()
	s.exitUUID = exit
	s.leftOn = &now
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type stepEnvelope struct {
	UUID      flows.StepUUID `json:"uuid" validate:"required,uuid4"`
	NodeUUID  flows.NodeUUID `json:"node_uuid" validate:"required,uuid4"`
	ExitUUID  flows.ExitUUID `json:"exit_uuid,omitempty" validate:"omitempty,uuid4"`
	ArrivedOn time.Time      `json:"arrived_on"`
	LeftOn    *time.Time     `json:"left_on,omitempty"`
}

// UnmarshalJSON unmarshals a run step from the given JSON
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
	return err
}

// MarshalJSON marshals this run step into JSON
func (s *step) MarshalJSON() ([]byte, error) {
	return json.Marshal(&stepEnvelope{
		UUID:      s.stepUUID,
		NodeUUID:  s.nodeUUID,
		ExitUUID:  s.exitUUID,
		ArrivedOn: s.arrivedOn,
		LeftOn:    s.leftOn,
	})
}
