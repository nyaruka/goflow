package resumes

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeWaitExpiration, readWaitExpiration)
}

// TypeWaitExpiration is the type for resuming a session when a wait has expired
const TypeWaitExpiration string = "wait_expiration"

// WaitExpiration is used when a session is resumed because the waiting run has expired
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
type WaitExpiration struct {
	baseResume

	event *events.WaitExpired
}

// NewWaitExpiration creates a new run expired resume with the passed in values
func NewWaitExpiration(event *events.WaitExpired) *WaitExpiration {
	return &WaitExpiration{
		baseResume: newBaseResume(TypeWaitExpiration),
		event:      event,
	}
}

// Event returns the event this resume is based on
func (r *WaitExpiration) Event() flows.Event { return r.event }

// Apply applies our state changes
func (r *WaitExpiration) Apply(run flows.Run, logEvent flows.EventCallback) {
	run.Exit(flows.RunStatusExpired)

	r.baseResume.Apply(run, logEvent)
}

var _ flows.Resume = (*WaitExpiration)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type waitExpirationEnvelope struct {
	baseEnvelope

	Event *events.WaitExpired `json:"event" validate:"required"`
}

func readWaitExpiration(sessionAssets flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Resume, error) {
	e := &waitExpirationEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &WaitExpiration{event: e.Event}

	if err := r.unmarshal(sessionAssets, &e.baseEnvelope, missing); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *WaitExpiration) MarshalJSON() ([]byte, error) {
	e := &waitExpirationEnvelope{Event: r.event}

	if err := r.marshal(&e.baseEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
