package waits

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/utils"
)

const (
	defaultDialLimit = time.Second * 60
	defaultCallLimit = time.Hour * 2
)

func init() {
	registerType(TypeDial, readDialWait)
}

// TypeDial is the type of our dial wait
const TypeDial string = "dial"

// DialWait is a wait which waits for a phone number to be dialed
type DialWait struct {
	baseWait

	phone     string
	dialLimit time.Duration
	callLimit time.Duration
}

// NewDialWait creates a new Dial wait
func NewDialWait(phone string, dialLimit, callLimit time.Duration) *DialWait {
	return &DialWait{
		baseWait:  newBaseWait(TypeDial, nil),
		phone:     phone,
		dialLimit: dialLimit,
		callLimit: callLimit,
	}
}

// DialLimit returns the time limit for dialing
func (w *DialWait) DialLimit() time.Duration {
	return w.dialLimit
}

// CallLimit returns the time limit for an answered call
func (w *DialWait) CallLimit() time.Duration {
	return w.callLimit
}

// AllowedFlowTypes returns the flow types which this wait is allowed to occur in
func (w *DialWait) AllowedFlowTypes() []flows.FlowType {
	return []flows.FlowType{flows.FlowTypeVoice}
}

// Begin beings waiting at this wait
func (w *DialWait) Begin(run flows.Run, log flows.EventCallback) bool {
	phone, err := run.EvaluateTemplate(w.phone)
	if err != nil {
		log(events.NewError(err))
	}

	urn, err := urns.NewTelURNForCountry(phone, string(run.Environment().DefaultCountry()))
	if err != nil {
		log(events.NewError(err))
		return false
	}

	log(events.NewDialWait(urn, int(w.dialLimit/time.Second), int(w.callLimit/time.Second), w.expiresOn(run)))

	return true
}

// Accept returns whether this wait accepts the given resume
func (w *DialWait) Accepts(resume flows.Resume) bool {
	return resume.Type() == resumes.TypeDial
}

var _ flows.Wait = (*DialWait)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type dialWaitEnvelope struct {
	baseWaitEnvelope

	Phone            string `json:"phone" validate:"required"`
	DialLimitSeconds int    `json:"dial_limit_seconds,omitempty"`
	CallLimitSeconds int    `json:"call_limit_seconds,omitempty"`
}

func readDialWait(data json.RawMessage) (flows.Wait, error) {
	e := &dialWaitEnvelope{
		DialLimitSeconds: int(defaultDialLimit / time.Second),
		CallLimitSeconds: int(defaultCallLimit / time.Second),
	}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	w := &DialWait{
		phone:     e.Phone,
		dialLimit: time.Second * time.Duration(e.DialLimitSeconds),
		callLimit: time.Second * time.Duration(e.CallLimitSeconds),
	}

	return w, w.unmarshal(&e.baseWaitEnvelope)
}

// MarshalJSON marshals this wait into JSON
func (w *DialWait) MarshalJSON() ([]byte, error) {
	e := &dialWaitEnvelope{
		Phone:            w.phone,
		DialLimitSeconds: int(w.dialLimit / time.Second),
		CallLimitSeconds: int(w.callLimit / time.Second),
	}

	if err := w.marshal(&e.baseWaitEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
