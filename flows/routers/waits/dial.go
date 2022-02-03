package waits

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeDial, readDialWait)
}

// TypeDial is the type of our dial wait
const TypeDial string = "dial"

// DialWait is a wait which waits for a phone number to be dialed
type DialWait struct {
	baseWait

	phone string
}

// NewDialWait creates a new Dial wait
func NewDialWait(phone string) *DialWait {
	return &DialWait{
		baseWait: newBaseWait(TypeDial, nil),
		phone:    phone,
	}
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

	log(events.NewDialWait(urn, w.expiresOn(run)))

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

	Phone string `json:"phone" validate:"required"`
}

func readDialWait(data json.RawMessage) (flows.Wait, error) {
	e := &dialWaitEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	w := &DialWait{phone: e.Phone}

	return w, w.unmarshal(&e.baseWaitEnvelope)
}

// MarshalJSON marshals this wait into JSON
func (w *DialWait) MarshalJSON() ([]byte, error) {
	e := &dialWaitEnvelope{Phone: w.phone}

	if err := w.marshal(&e.baseWaitEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
