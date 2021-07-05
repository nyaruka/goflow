package actions

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

// max number of bytes to be saved to extra on a result
const resultExtraMaxBytes = 10000

// max length of a message attachment (type:url)
const maxAttachmentLength = 2048

// common category names
const (
	CategorySuccess = "Success"
	CategorySkipped = "Skipped"
	CategoryFailure = "Failure"
)

var webhookCategories = []string{CategorySuccess, CategoryFailure}
var webhookStatusCategories = map[flows.CallStatus]string{
	flows.CallStatusSuccess:         CategorySuccess,
	flows.CallStatusResponseError:   CategoryFailure,
	flows.CallStatusConnectionError: CategoryFailure,
	flows.CallStatusSubscriberGone:  CategoryFailure,
}

var registeredTypes = map[string](func() flows.Action){}

// registers a new type of action
func registerType(name string, initFunc func() flows.Action) {
	registeredTypes[name] = initFunc
}

// RegisteredTypes gets the registered types of action
func RegisteredTypes() map[string](func() flows.Action) {
	return registeredTypes
}

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// the base of all action types
type baseAction struct {
	Type_ string           `json:"type" validate:"required"`
	UUID_ flows.ActionUUID `json:"uuid" validate:"required,uuid4"`
}

// creates a new base action
func newBaseAction(typeName string, uuid flows.ActionUUID) baseAction {
	return baseAction{Type_: typeName, UUID_: uuid}
}

// Type returns the type of this action
func (a *baseAction) Type() string { return a.Type_ }

// UUID returns the UUID of the action
func (a *baseAction) UUID() flows.ActionUUID { return a.UUID_ }

// Validate validates our action is valid
func (a *baseAction) Validate() error { return nil }

// LocalizationUUID gets the UUID which identifies this object for localization
func (a *baseAction) LocalizationUUID() uuids.UUID { return uuids.UUID(a.UUID_) }

// helper function for actions that send a message (text + attachments) that must be localized and evalulated
func (a *baseAction) evaluateMessage(run flows.FlowRun, languages []envs.Language, actionText string, actionAttachments []string, actionQuickReplies []string, logEvent flows.EventCallback) (string, []utils.Attachment, []string) {
	// localize and evaluate the message text
	localizedText := run.GetTranslatedTextArray(uuids.UUID(a.UUID()), "text", []string{actionText}, languages)[0]
	evaluatedText, err := run.EvaluateTemplate(localizedText)
	if err != nil {
		logEvent(events.NewError(err))
	}

	// localize and evaluate the message attachments
	translatedAttachments := run.GetTranslatedTextArray(uuids.UUID(a.UUID()), "attachments", actionAttachments, languages)
	evaluatedAttachments := make([]utils.Attachment, 0, len(translatedAttachments))
	for _, a := range translatedAttachments {
		evaluatedAttachment, err := run.EvaluateTemplate(a)
		if err != nil {
			logEvent(events.NewError(err))
		}
		if evaluatedAttachment == "" {
			logEvent(events.NewErrorf("attachment text evaluated to empty string, skipping"))
			continue
		}
		if len(evaluatedAttachment) > maxAttachmentLength {
			logEvent(events.NewErrorf("evaluated attachment is longer than %d limit, skipping", maxAttachmentLength))
			continue
		}
		evaluatedAttachments = append(evaluatedAttachments, utils.Attachment(evaluatedAttachment))
	}

	// localize and evaluate the quick replies
	translatedQuickReplies := run.GetTranslatedTextArray(uuids.UUID(a.UUID()), "quick_replies", actionQuickReplies, languages)
	evaluatedQuickReplies := make([]string, 0, len(translatedQuickReplies))
	for _, qr := range translatedQuickReplies {
		evaluatedQuickReply, err := run.EvaluateTemplate(qr)
		if err != nil {
			logEvent(events.NewError(err))
		}
		if evaluatedQuickReply == "" {
			logEvent(events.NewErrorf("quick reply text evaluated to empty string, skipping"))
			continue
		}
		evaluatedQuickReplies = append(evaluatedQuickReplies, evaluatedQuickReply)
	}

	return evaluatedText, evaluatedAttachments, evaluatedQuickReplies
}

// helper to save a run result and log it as an event
func (a *baseAction) saveResult(run flows.FlowRun, step flows.Step, name, value, category, categoryLocalized string, input string, extra json.RawMessage, logEvent flows.EventCallback) {
	result := flows.NewResult(name, value, category, categoryLocalized, step.NodeUUID(), input, extra, dates.Now())
	run.SaveResult(result)
	logEvent(events.NewRunResultChanged(result))
}

// helper to save a run result based on a webhook call and log it as an event
func (a *baseAction) saveWebhookResult(run flows.FlowRun, step flows.Step, name string, call *flows.WebhookCall, status flows.CallStatus, logEvent flows.EventCallback) {
	input := fmt.Sprintf("%s %s", call.Request.Method, call.Request.URL.String())
	value := "0"
	category := webhookStatusCategories[status]
	var extra json.RawMessage

	if call.Response != nil {
		value = strconv.Itoa(call.Response.StatusCode)

		if len(call.ResponseJSON) > 0 && len(call.ResponseJSON) < resultExtraMaxBytes {
			extra = call.ResponseBody
		}
	}

	a.saveResult(run, step, name, value, category, "", input, extra, logEvent)
}

func (a *baseAction) updateWebhook(run flows.FlowRun, call *flows.WebhookCall) {
	parsed := types.JSONToXValue(call.ResponseBody)

	switch typed := parsed.(type) {
	case nil, types.XError:
		run.SetWebhook(types.XObjectEmpty)
	default:
		run.SetWebhook(typed)
	}
}

// helper to apply a contact modifier
func (a *baseAction) applyModifier(run flows.FlowRun, mod flows.Modifier, logModifier flows.ModifierCallback, logEvent flows.EventCallback) {
	mod.Apply(run.Environment(), run.Session().Assets(), run.Contact(), logEvent)
	logModifier(mod)
}

// helper to log a failure
func (a *baseAction) fail(run flows.FlowRun, err error, logEvent flows.EventCallback) {
	run.Exit(flows.RunStatusFailed)
	logEvent(events.NewFailure(err))
}

// utility struct which sets the allowed flow types to any
type universalAction struct{}

// AllowedFlowTypes returns the flow types which this action is allowed to occur in
func (a *universalAction) AllowedFlowTypes() []flows.FlowType {
	return []flows.FlowType{flows.FlowTypeMessaging, flows.FlowTypeMessagingBackground, flows.FlowTypeMessagingOffline, flows.FlowTypeVoice}
}

// utility struct which sets the allowed flow types to non-background
type interactiveAction struct{}

// AllowedFlowTypes returns the flow types which this action is allowed to occur in
func (a *interactiveAction) AllowedFlowTypes() []flows.FlowType {
	return []flows.FlowType{flows.FlowTypeMessaging, flows.FlowTypeMessagingOffline, flows.FlowTypeVoice}
}

// utility struct which sets the allowed flow types to any which run online
type onlineAction struct{}

// AllowedFlowTypes returns the flow types which this action is allowed to occur in
func (a *onlineAction) AllowedFlowTypes() []flows.FlowType {
	return []flows.FlowType{flows.FlowTypeMessaging, flows.FlowTypeMessagingBackground, flows.FlowTypeVoice}
}

// utility struct which sets the allowed flow types to just voice
type voiceAction struct{}

// AllowedFlowTypes returns the flow types which this action is allowed to occur in
func (a *voiceAction) AllowedFlowTypes() []flows.FlowType {
	return []flows.FlowType{flows.FlowTypeVoice}
}

// utility struct for actions which operate on other contacts
type otherContactsAction struct {
	URNs         []urns.URN                `json:"urns,omitempty"`
	Groups       []*assets.GroupReference  `json:"groups,omitempty" validate:"dive"`
	Contacts     []*flows.ContactReference `json:"contacts,omitempty" validate:"dive"`
	ContactQuery string                    `json:"contact_query,omitempty" engine:"evaluated"`
	LegacyVars   []string                  `json:"legacy_vars,omitempty" engine:"evaluated"`
}

func (a *otherContactsAction) resolveRecipients(run flows.FlowRun, logEvent flows.EventCallback) ([]*assets.GroupReference, []*flows.ContactReference, string, []urns.URN, error) {
	groupSet := run.Session().Assets().Groups()

	// copy URNs
	urnList := make([]urns.URN, 0, len(a.URNs))
	urnList = append(urnList, a.URNs...)

	// copy contact references
	contactRefs := make([]*flows.ContactReference, 0, len(a.Contacts))
	contactRefs = append(contactRefs, a.Contacts...)

	// resolve group references
	groups, err := resolveGroups(run, a.Groups, logEvent)
	if err != nil {
		return nil, nil, "", nil, err
	}
	groupRefs := make([]*assets.GroupReference, 0, len(groups))
	for _, group := range groups {
		groupRefs = append(groupRefs, group.Reference())
	}

	// evaluate the legacy variables
	for _, legacyVar := range a.LegacyVars {
		evaluatedLegacyVar, err := run.EvaluateTemplate(legacyVar)
		if err != nil {
			logEvent(events.NewError(err))
		}

		evaluatedLegacyVar = strings.TrimSpace(evaluatedLegacyVar)

		if uuidRegex.MatchString(evaluatedLegacyVar) {
			// if variable evaluates to a UUID, we assume it's a contact UUID
			contactRefs = append(contactRefs, flows.NewContactReference(flows.ContactUUID(evaluatedLegacyVar), ""))

		} else if groupByName := groupSet.FindByName(evaluatedLegacyVar); groupByName != nil {
			// next up we look for a group with a matching name
			groupRefs = append(groupRefs, groupByName.Reference())
		} else {
			// next up try it as a URN
			urn := urns.URN(evaluatedLegacyVar)
			if urn.Validate() == nil {
				urn = urn.Normalize(string(run.Environment().DefaultCountry()))
				urnList = append(urnList, urn)
			} else {
				// if that fails, assume this is a phone number, and let the caller worry about validation
				urn, err := urns.NewURNFromParts(urns.TelScheme, evaluatedLegacyVar, "", "")
				if err != nil {
					logEvent(events.NewError(err))
				} else {
					urn = urn.Normalize(string(run.Environment().DefaultCountry()))
					urnList = append(urnList, urn)
				}
			}
		}
	}

	// evaluate contact query
	contactQuery, _ := run.EvaluateTemplateText(a.ContactQuery, flows.ContactQueryEscaping, true)
	contactQuery = strings.TrimSpace(contactQuery)

	return groupRefs, contactRefs, contactQuery, urnList, nil
}

// utility struct for actions which create a message
type createMsgAction struct {
	Text         string   `json:"text" validate:"required" engine:"localized,evaluated"`
	Attachments  []string `json:"attachments,omitempty" engine:"localized,evaluated"`
	QuickReplies []string `json:"quick_replies,omitempty" engine:"localized,evaluated"`
}

// helper function for actions that have a set of group references that must be resolved to actual groups
func resolveGroups(run flows.FlowRun, references []*assets.GroupReference, logEvent flows.EventCallback) ([]*flows.Group, error) {
	groupSet := run.Session().Assets().Groups()
	groups := make([]*flows.Group, 0, len(references))

	for _, ref := range references {
		var group *flows.Group

		if ref.UUID != "" {
			// group is a fixed group with a UUID
			group = groupSet.Get(ref.UUID)
			if group == nil {
				logEvent(events.NewDependencyError(ref))
			}
		} else {
			// group is an expression that evaluates to an existing group's name
			evaluatedGroupName, err := run.EvaluateTemplate(ref.NameMatch)
			if err != nil {
				logEvent(events.NewError(err))
			} else {
				// look up the set of all groups to see if such a group exists
				group = groupSet.FindByName(evaluatedGroupName)
				if group == nil {
					logEvent(events.NewErrorf("no such group with name '%s'", evaluatedGroupName))
				}
			}
		}

		if group != nil {
			groups = append(groups, group)
		}
	}

	return groups, nil
}

// helper function for actions that have a set of label references that must be resolved to actual labels
func resolveLabels(run flows.FlowRun, references []*assets.LabelReference, logEvent flows.EventCallback) ([]*flows.Label, error) {
	labelSet := run.Session().Assets().Labels()
	labels := make([]*flows.Label, 0, len(references))

	for _, ref := range references {
		var label *flows.Label

		if ref.UUID != "" {
			// label is a fixed label with a UUID
			label = labelSet.Get(ref.UUID)
			if label == nil {
				logEvent(events.NewDependencyError(ref))
			}
		} else {
			// label is an expression that evaluates to an existing label's name
			evaluatedLabelName, err := run.EvaluateTemplate(ref.NameMatch)
			if err != nil {
				logEvent(events.NewError(err))
			} else {
				// look up the set of all labels to see if such a label exists
				label = labelSet.FindByName(evaluatedLabelName)
				if label == nil {
					logEvent(events.NewErrorf("no such label with name '%s'", evaluatedLabelName))
				}
			}
		}

		if label != nil {
			labels = append(labels, label)
		}
	}

	return labels, nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadAction reads an action from the given JSON
func ReadAction(data json.RawMessage) (flows.Action, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, errors.Errorf("unknown type: '%s'", typeName)
	}

	action := f()
	return action, utils.UnmarshalAndValidate(data, action)
}
