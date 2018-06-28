package runs

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type step struct {
	stepUUID  flows.StepUUID
	nodeUUID  flows.NodeUUID
	exitUUID  flows.ExitUUID
	arrivedOn time.Time
}

func (s *step) UUID() flows.StepUUID     { return s.stepUUID }
func (s *step) NodeUUID() flows.NodeUUID { return s.nodeUUID }
func (s *step) ExitUUID() flows.ExitUUID { return s.exitUUID }
func (s *step) ArrivedOn() time.Time     { return s.arrivedOn }

func (s *step) Leave(exit flows.ExitUUID) {
	s.exitUUID = exit
}

func (s *step) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "uuid":
		return types.NewXText(string(s.UUID()))
	case "node_uuid":
		return types.NewXText(string(s.NodeUUID()))
	case "arrived_on":
		return types.NewXDateTime(s.ArrivedOn())
	case "exit_uuid":
		return types.NewXText(string(s.ExitUUID()))
	default:
		return types.NewXResolveError(s, key)
	}
}

// Describe returns a representation of this type for error messages
func (s *step) Describe() string { return "step" }

func (s *step) Reduce(env utils.Environment) types.XPrimitive {
	return types.NewXText(string(s.UUID()))
}

func (s *step) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, s, "uuid", "node_uuid", "arrived_on", "exit_uuid").ToXJSON(env)
}

var _ flows.Step = (*step)(nil)

type Path []flows.Step

func (p Path) Length() int {
	return len(p)
}

func (p Path) Index(index int) types.XValue {
	return p[index]
}

// Describe returns a representation of this type for error messages
func (p Path) Describe() string { return "path" }

func (p Path) Reduce(env utils.Environment) types.XPrimitive {
	array := types.NewXArray()
	for _, step := range p {
		array.Append(step)
	}
	return array
}

func (p Path) ToXJSON(env utils.Environment) types.XText {
	return p.Reduce(env).ToXJSON(env)
}

var _ types.XValue = (Path)(nil)
var _ types.XIndexable = (Path)(nil)

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
	var err error

	err = json.Unmarshal(data, &se)
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
	return json.Marshal(&stepEnvelope{
		UUID:      s.stepUUID,
		NodeUUID:  s.nodeUUID,
		ExitUUID:  s.exitUUID,
		ArrivedOn: s.arrivedOn,
	})
}
