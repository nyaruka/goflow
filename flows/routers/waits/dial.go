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
	registerType(TypeDial, readDialWait, readActivatedDialWait)
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
func (w *DialWait) Begin(run flows.FlowRun, log flows.EventCallback) flows.ActivatedWait {
	phone, err := run.EvaluateTemplate(w.phone)
	if err != nil {
		log(events.NewError(err))
	}

	urn, err := urns.NewTelURNForCountry(phone, string(run.Environment().DefaultCountry()))
	if err != nil {
		log(events.NewError(err))
		return nil
	}

	log(events.NewDialWait(urn))

	return NewActivatedDialWait(urn)
}

// End ends this wait or returns an error
func (w *DialWait) End(resume flows.Resume) error {
	if resume.Type() == resumes.TypeDial {
		return nil
	}
	return w.resumeTypeError(resume)
}

var _ flows.Wait = (*DialWait)(nil)

type ActivatedDialWait struct {
	baseActivatedWait

	urn urns.URN
}

func NewActivatedDialWait(urn urns.URN) *ActivatedDialWait {
	return &ActivatedDialWait{
		baseActivatedWait: baseActivatedWait{type_: TypeDial},
		urn:               urn,
	}
}

func (w *ActivatedDialWait) URN() urns.URN {
	return w.urn
}

var _ flows.ActivatedWait = (*ActivatedDialWait)(nil)

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

type activatedDialWaitEnvelope struct {
	baseActivatedWaitEnvelope

	URN urns.URN `json:"urn" validate:"required,urn"`
}

func readActivatedDialWait(data json.RawMessage) (flows.ActivatedWait, error) {
	e := &activatedDialWaitEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	w := &ActivatedDialWait{urn: e.URN}

	return w, w.unmarshal(&e.baseActivatedWaitEnvelope)
}

// MarshalJSON marshals this wait into JSON
func (w *ActivatedDialWait) MarshalJSON() ([]byte, error) {
	e := &activatedDialWaitEnvelope{URN: w.urn}

	if err := w.marshal(&e.baseActivatedWaitEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
