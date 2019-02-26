package actions

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeCallResthook, func() flows.Action { return &CallResthookAction{} })
}

// TypeCallResthook is the type for the call resthook action
const TypeCallResthook string = "call_resthook"

// CallResthookAction can be used to call a resthook.
//
// A [event:webhook_called] event will be created for each subscriber of the resthook with the results
// of the HTTP call. If the action has `result_name` set, a result will
// be created with that name, and if the resthook returns valid JSON, that will be accessible
// through `extra` on the result.
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
	onlineAction

	Resthook   string `json:"resthook" validate:"required"`
	ResultName string `json:"result_name,omitempty"`
}

// NewCallResthookAction creates a new call resthook action
func NewCallResthookAction(uuid flows.ActionUUID, resthook string, resultName string) *CallResthookAction {
	return &CallResthookAction{
		BaseAction: NewBaseAction(TypeCallResthook, uuid),
		Resthook:   resthook,
		ResultName: resultName,
	}
}

// Validate validates our action is valid and has all the assets it needs
func (a *CallResthookAction) Validate(assets flows.SessionAssets, context *flows.ValidationContext) error {
	return nil
}

// Execute runs this action
func (a *CallResthookAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	// NOOP if resthook doesn't exist
	resthook := run.Session().Assets().Resthooks().FindBySlug(a.Resthook)
	if resthook == nil {
		return nil
	}

	// build our payload
	payload, err := run.EvaluateTemplate(flows.DefaultWebhookPayload)
	if err != nil {
		logEvent(events.NewErrorEvent(err))
	}

	// regardless of what subscriber calls we make, we need to record the payload that would be sent
	logEvent(events.NewResthookCalledEvent(a.Resthook, json.RawMessage(payload)))

	// make a call to each subscriber URL
	webhooks := make([]*flows.WebhookCall, 0, len(resthook.Subscribers()))

	for _, url := range resthook.Subscribers() {
		req, err := http.NewRequest("POST", url, strings.NewReader(payload))
		if err != nil {
			logEvent(events.NewErrorEvent(err))
			return nil
		}

		req.Header.Add("Content-Type", "application/json")

		webhook, err := flows.MakeWebhookCall(run.Session(), req, a.Resthook)
		if err != nil {
			logEvent(events.NewErrorEvent(err))
		} else {
			webhooks = append(webhooks, webhook)
			logEvent(events.NewWebhookCalledEvent(webhook))
		}
	}

	asResult := a.pickResultWebhook(webhooks)
	if asResult != nil && a.ResultName != "" {
		a.saveWebhookResult(run, step, a.ResultName, asResult, logEvent)
	}

	return nil
}

// picks one of the resthook calls to become the result generated by this action
func (a *CallResthookAction) pickResultWebhook(calls []*flows.WebhookCall) *flows.WebhookCall {
	var lastFailure, asResult *flows.WebhookCall
	for _, call := range calls {
		if call.Status() == flows.WebhookStatusSuccess {
			asResult = call
		} else {
			lastFailure = call
		}
	}
	if lastFailure != nil {
		asResult = lastFailure
	}
	return asResult
}
