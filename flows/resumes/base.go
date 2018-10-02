package resumes

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

type readFunc func(session flows.Session, data json.RawMessage) (flows.Resume, error)

var registeredTypes = map[string]readFunc{}

// RegisterType registers a new type of trigger
func RegisterType(name string, f readFunc) {
	registeredTypes[name] = f
}

type baseResume struct {
	environment utils.Environment
	contact     *flows.Contact
	resumedOn   time.Time
}

func newBaseResume(env utils.Environment, contact *flows.Contact) baseResume {
	return baseResume{environment: env, contact: contact, resumedOn: utils.Now()}
}

func (r *baseResume) Environment() utils.Environment { return r.environment }
func (r *baseResume) Contact() *flows.Contact        { return r.contact }
func (r *baseResume) ResumedOn() time.Time           { return r.resumedOn }

// Apply applies our state changes and saves any events to the run
func (r *baseResume) Apply(run flows.FlowRun, step flows.Step) error {
	if r.environment != nil {
		run.Session().SetEnvironment(r.environment)

		// TODO diffing

		run.AddEvent(step, events.NewEnvironmentChangedEvent(r.Environment()))
	}
	if r.contact != nil {
		run.Session().SetContact(r.contact.Clone())

		// TODO diffing

		run.AddEvent(step, events.NewContactChangedEvent(r.contact))
	}

	return nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type baseResumeEnvelope struct {
	Environment json.RawMessage `json:"environment,omitempty"`
	Contact     json.RawMessage `json:"contact,omitempty"`
	ResumedOn   time.Time       `json:"resumed_on" validate:"required"`
}

// ReadResume reads a resume from the given typed envelope
func ReadResume(session flows.Session, envelope *utils.TypedEnvelope) (flows.Resume, error) {
	f := registeredTypes[envelope.Type]
	if f == nil {
		return nil, fmt.Errorf("unknown type: %s", envelope.Type)
	}
	return f(session, envelope.Data)
}

func unmarshalBaseResume(session flows.Session, base *baseResume, envelope *baseResumeEnvelope) error {
	var err error

	base.resumedOn = envelope.ResumedOn

	if envelope.Environment != nil {
		if base.environment, err = utils.ReadEnvironment(envelope.Environment); err != nil {
			return fmt.Errorf("unable to read environment: %s", err)
		}
	}
	if envelope.Contact != nil {
		if base.contact, err = flows.ReadContact(session.Assets(), envelope.Contact, true); err != nil {
			return fmt.Errorf("unable to read contact: %s", err)
		}
	}
	return nil
}
