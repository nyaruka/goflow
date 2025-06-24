package resumes

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeWaitExpiration, readWaitExpirationResume)
}

// TypeWaitExpiration is the type for resuming a session when a wait has expired
const TypeWaitExpiration string = "wait_expiration"

// WaitExpirationResume is used when a session is resumed because the waiting run has expired
//
//	{
//	  "type": "wait_expiration",
//	  "resumed_on": "2000-01-01T00:00:00.000000000-00:00",
//	  "event": {
//	    "type": "wait_expired",
//	    "created_on": "2006-01-02T15:04:05Z"
//	  }
//	}
//
// @resume wait_expiration
type WaitExpirationResume struct {
	baseResume

	event *events.WaitExpiredEvent
}

// NewWaitExpiration creates a new run expired resume with the passed in values
func NewWaitExpiration(event *events.WaitExpiredEvent) *WaitExpirationResume {
	return &WaitExpirationResume{
		baseResume: newBaseResume(TypeWaitExpiration),
		event:      event,
	}
}

// Event returns the event this resume is based on
func (r *WaitExpirationResume) Event() flows.Event { return r.event }

// Apply applies our state changes
func (r *WaitExpirationResume) Apply(run flows.Run, logEvent flows.EventCallback) {
	run.Exit(flows.RunStatusExpired)

	r.baseResume.Apply(run, logEvent)
}

var _ flows.Resume = (*WaitExpirationResume)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type waitExpirationResumeEnvelope struct {
	baseResumeEnvelope

	Event *events.WaitExpiredEvent `json:"event"` // TODO make required
}

func readWaitExpirationResume(sessionAssets flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Resume, error) {
	e := &waitExpirationResumeEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &WaitExpirationResume{event: e.Event}

	if err := r.unmarshal(sessionAssets, &e.baseResumeEnvelope, missing); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *WaitExpirationResume) MarshalJSON() ([]byte, error) {
	e := &waitExpirationResumeEnvelope{Event: r.event}

	if err := r.marshal(&e.baseResumeEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
