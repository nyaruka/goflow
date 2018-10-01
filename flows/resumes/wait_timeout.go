package resumes

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeWaitTimeout, ReadWaitTimeoutResume)
}

// TypeWaitTimeout is the type for resuming a session when a wait has timed out
const TypeWaitTimeout string = "wait_timeout"

// WaitTimeoutResume is used when a session is resumed because a wait has timed out
//
//   {
//     "type": "wait_timeout",
//     "contact": {
//       "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
//       "name": "Bob",
//       "language": "fra",
//       "fields": {"gender": {"text": "Male"}},
//       "groups": []
//     }
//   }
//
// @resume wait_timeout
type WaitTimeoutResume struct {
	baseResume
}

// NewWaitTimeoutResume creates a new timeout resume with the passed in values
func NewWaitTimeoutResume(env utils.Environment, contact *flows.Contact) *WaitTimeoutResume {
	return &WaitTimeoutResume{
		baseResume: baseResume{
			environment: env,
			contact:     contact,
		},
	}
}

// Type returns the type of this resume
func (t *WaitTimeoutResume) Type() string { return TypeWaitTimeout }

var _ flows.Resume = (*WaitTimeoutResume)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadWaitTimeoutResume reads a timeout resume
func ReadWaitTimeoutResume(session flows.Session, data json.RawMessage) (flows.Resume, error) {
	resume := &WaitTimeoutResume{}
	e := baseResumeEnvelope{}
	if err := utils.UnmarshalAndValidate(data, &e); err != nil {
		return nil, err
	}

	if err := unmarshalBaseResume(session, &resume.baseResume, &e); err != nil {
		return nil, err
	}

	return resume, nil
}
