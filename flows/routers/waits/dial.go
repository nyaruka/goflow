package waits

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/gocommon/dates"
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
	registerType(TypeDial, readDial)
}

// TypeDial is the type of our dial wait
const TypeDial string = "dial"

// Dial is a wait which waits for a phone number to be dialed
type Dial struct {
	baseWait

	phone     string
	dialLimit time.Duration
	callLimit time.Duration
}

// NewDial creates a new Dial wait
func NewDial(phone string, dialLimit, callLimit time.Duration) *Dial {
	return &Dial{
		baseWait:  newBaseWait(TypeDial, nil),
		phone:     phone,
		dialLimit: dialLimit,
		callLimit: callLimit,
	}
}

// DialLimit returns the time limit for dialing
func (w *Dial) DialLimit() time.Duration {
	return w.dialLimit
}

// CallLimit returns the time limit for an answered call
func (w *Dial) CallLimit() time.Duration {
	return w.callLimit
}

// AllowedFlowTypes returns the flow types which this wait is allowed to occur in
func (w *Dial) AllowedFlowTypes() []flows.FlowType {
	return []flows.FlowType{flows.FlowTypeVoice}
}

// Begin beings waiting at this wait
func (w *Dial) Begin(run flows.Run, log flows.EventLogger) bool {
	phone, _ := run.EvaluateTemplate(w.phone, log)
	country := run.Session().MergedEnvironment().DefaultCountry()

	urn, err := urns.ParsePhone(phone, country, false, false)
	if err != nil {
		log(events.NewError(err.Error()))
		return false
	}

	// we don't want to expire the flow whilst the contact is in the forwarded call and appearing "inactive" in the
	// flow so calculate an expiry guaranteed to be after the wait returns
	expiresOn := dates.Now().Add(w.dialLimit + w.callLimit + time.Second*30)

	log(events.NewDialWait(urn, int(w.dialLimit/time.Second), int(w.callLimit/time.Second), expiresOn))

	return true
}

// Accept returns whether this wait accepts the given resume
func (w *Dial) Accepts(resume flows.Resume) bool {
	return resume.Type() == resumes.TypeDial
}

var _ flows.Wait = (*Dial)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type dialEnvelope struct {
	baseEnvelope

	Phone            string `json:"phone" validate:"required"`
	DialLimitSeconds int    `json:"dial_limit_seconds,omitempty"`
	CallLimitSeconds int    `json:"call_limit_seconds,omitempty"`
}

func readDial(data json.RawMessage) (flows.Wait, error) {
	e := &dialEnvelope{
		DialLimitSeconds: int(defaultDialLimit / time.Second),
		CallLimitSeconds: int(defaultCallLimit / time.Second),
	}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	w := &Dial{
		phone:     e.Phone,
		dialLimit: time.Second * time.Duration(e.DialLimitSeconds),
		callLimit: time.Second * time.Duration(e.CallLimitSeconds),
	}

	return w, w.unmarshal(&e.baseEnvelope)
}

// MarshalJSON marshals this wait into JSON
func (w *Dial) MarshalJSON() ([]byte, error) {
	e := &dialEnvelope{
		Phone:            w.phone,
		DialLimitSeconds: int(w.dialLimit / time.Second),
		CallLimitSeconds: int(w.callLimit / time.Second),
	}

	if err := w.marshal(&e.baseEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
