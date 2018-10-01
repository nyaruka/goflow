package resumes

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
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
}

func (t *baseResume) Environment() utils.Environment { return t.environment }
func (t *baseResume) Contact() *flows.Contact        { return t.contact }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type baseResumeEnvelope struct {
	Environment json.RawMessage `json:"environment,omitempty"`
	Contact     json.RawMessage `json:"contact,omitempty"`
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
