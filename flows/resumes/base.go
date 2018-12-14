package resumes

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

type readFunc func(session flows.Session, data json.RawMessage) (flows.Resume, error)

var registeredTypes = map[string]readFunc{}

// RegisterType registers a new type of trigger
func RegisterType(name string, f readFunc) {
	registeredTypes[name] = f
}

type baseResume struct {
	type_       string
	environment utils.Environment
	contact     *flows.Contact
	resumedOn   time.Time
}

func newBaseResume(typeName string, env utils.Environment, contact *flows.Contact) baseResume {
	return baseResume{type_: typeName, environment: env, contact: contact, resumedOn: utils.Now()}
}

// Type returns the type of this resume
func (r *baseResume) Type() string { return r.type_ }

func (r *baseResume) Environment() utils.Environment { return r.environment }
func (r *baseResume) Contact() *flows.Contact        { return r.contact }
func (r *baseResume) ResumedOn() time.Time           { return r.resumedOn }

// Apply applies our state changes and saves any events to the run
func (r *baseResume) Apply(run flows.FlowRun, logEvent func(flows.Event)) error {
	if r.environment != nil {
		if !run.Session().Environment().Equal(r.environment) {
			logEvent(events.NewEnvironmentRefreshedEvent(r.environment))
		}

		run.Session().SetEnvironment(r.environment)
	}
	if r.contact != nil {
		if !run.Session().Contact().Equal(r.contact) {
			logEvent(events.NewContactRefreshedEvent(r.contact))
		}

		run.Session().SetContact(r.contact)

		triggers.EnsureDynamicGroups(run.Session())
	}

	if run.Status() == flows.RunStatusWaiting {
		run.SetStatus(flows.RunStatusActive)
	}

	return nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type baseResumeEnvelope struct {
	Type        string          `json:"type" validate:"required"`
	Environment json.RawMessage `json:"environment,omitempty"`
	Contact     json.RawMessage `json:"contact,omitempty"`
	ResumedOn   time.Time       `json:"resumed_on" validate:"required"`
}

// ReadResume reads a resume from the given JSON
func ReadResume(session flows.Session, data json.RawMessage) (flows.Resume, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, errors.Errorf("unknown type: '%s'", typeName)
	}
	return f(session, data)
}

func (r *baseResume) unmarshal(session flows.Session, e *baseResumeEnvelope) error {
	var err error

	r.type_ = e.Type
	r.resumedOn = e.ResumedOn

	if e.Environment != nil {
		if r.environment, err = utils.ReadEnvironment(e.Environment); err != nil {
			return errors.Wrap(err, "unable to read environment")
		}
	}
	if e.Contact != nil {
		if r.contact, err = flows.ReadContact(session.Assets(), e.Contact, true); err != nil {
			return errors.Wrap(err, "unable to read contact")
		}
	}
	return nil
}

func (r *baseResume) marshal(e *baseResumeEnvelope) error {
	var err error
	e.Type = r.type_
	e.ResumedOn = r.resumedOn

	if r.environment != nil {
		e.Environment, err = json.Marshal(r.environment)
		if err != nil {
			return err
		}
	}
	if r.contact != nil {
		e.Contact, err = json.Marshal(r.contact)
		if err != nil {
			return err
		}
	}
	return nil
}
