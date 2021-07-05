package runs

import (
	"time"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
)

type step struct {
	stepUUID  flows.StepUUID
	nodeUUID  flows.NodeUUID
	exitUUID  flows.ExitUUID
	arrivedOn time.Time
}

// NewStep creates a new step
func NewStep(node flows.Node, arrivedOn time.Time) flows.Step {
	return &step{
		stepUUID:  flows.StepUUID(uuids.New()),
		nodeUUID:  node.UUID(),
		arrivedOn: arrivedOn,
	}
}

func (s *step) UUID() flows.StepUUID     { return s.stepUUID }
func (s *step) NodeUUID() flows.NodeUUID { return s.nodeUUID }
func (s *step) ExitUUID() flows.ExitUUID { return s.exitUUID }
func (s *step) ArrivedOn() time.Time     { return s.arrivedOn }

func (s *step) Leave(exit flows.ExitUUID) {
	s.exitUUID = exit
}

// Context returns the properties available in expressions
func (s *step) Context(env envs.Environment) map[string]types.XValue {
	return map[string]types.XValue{
		"uuid":       types.NewXText(string(s.UUID())),
		"node_uuid":  types.NewXText(string(s.NodeUUID())),
		"arrived_on": types.NewXDateTime(s.ArrivedOn()),
		"exit_uuid":  types.NewXText(string(s.ExitUUID())),
	}
}

var _ flows.Step = (*step)(nil)

// Path is the steps taken in a run
type Path []flows.Step

// ToXValue returns a representation of this object for use in expressions
func (p Path) ToXValue(env envs.Environment) types.XValue {
	array := make([]types.XValue, len(p))
	for i, step := range p {
		array[i] = flows.Context(env, step)
	}
	return types.NewXArray(array...)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type stepEnvelope struct {
	UUID      flows.StepUUID `json:"uuid" validate:"required,uuid4"`
	NodeUUID  flows.NodeUUID `json:"node_uuid" validate:"required,uuid4"`
	ExitUUID  flows.ExitUUID `json:"exit_uuid,omitempty" validate:"omitempty,uuid4"`
	ArrivedOn time.Time      `json:"arrived_on"`
}

// UnmarshalJSON unmarshals a run step from the given JSON
func (s *step) UnmarshalJSON(data []byte) error {
	var se stepEnvelope

	err := jsonx.Unmarshal(data, &se)
	if err != nil {
		return err
	}

	s.stepUUID = se.UUID
	s.nodeUUID = se.NodeUUID
	s.exitUUID = se.ExitUUID
	s.arrivedOn = se.ArrivedOn
	return err
}

// MarshalJSON marshals this run step into JSON
func (s *step) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&stepEnvelope{
		UUID:      s.stepUUID,
		NodeUUID:  s.nodeUUID,
		ExitUUID:  s.exitUUID,
		ArrivedOn: s.arrivedOn,
	})
}
