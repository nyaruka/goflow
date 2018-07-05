package actions

import (
	"net/http"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// TypeCallResthook is the type for the call resthook action
const TypeCallResthook string = "call_resthook"

// CallResthookAction can be used to call a resthook.
//
// A `resthook_subscriber_called` event will be created based on the results of the HTTP call
// to each subscriber of the resthook.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "call_resthook",
//     "resthook": "new-registration"
//   }
//
// @action call_resthook
type CallResthookAction struct {
	BaseAction
	Resthook string `json:"resthook" validate:"required"`
}

// Type returns the type of this action
func (a *CallResthookAction) Type() string { return TypeCallResthook }

// Validate validates our action is valid and has all the assets it needs
func (a *CallResthookAction) Validate(assets flows.SessionAssets) error {
	_, err := assets.GetResthook(a.Resthook)
	return err
}

// Execute runs this action
func (a *CallResthookAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	// lookup our resthook asset
	resthook, err := run.Session().Assets().GetResthook(a.Resthook)
	if err != nil {
		return err
	}

	// make a request for each URL
	for _, url := range resthook.Subscribers() {
		req, err := http.NewRequest("POST", url, strings.NewReader(flows.DefaultWebhookPayload))
		if err != nil {
			log.Add(events.NewErrorEvent(err))
			return nil
		}

		webhook, err := flows.MakeWebhookCall(run.Session(), req)

		if err != nil {
			log.Add(events.NewErrorEvent(err))
		} else {
			log.Add(events.NewResthookSubscriberCalledEvent(a.Resthook, webhook.URL(), webhook.Status(), webhook.StatusCode(), webhook.Request(), webhook.Response()))
		}
	}
	return nil
}
