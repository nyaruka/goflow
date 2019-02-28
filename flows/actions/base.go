package actions

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

var webhookStatusCategories = map[flows.WebhookStatus]string{
	flows.WebhookStatusSuccess:         "Success",
	flows.WebhookStatusResponseError:   "Failure",
	flows.WebhookStatusConnectionError: "Failure",
	flows.WebhookStatusSubscriberGone:  "Failure",
}

var registeredTypes = map[string](func() flows.Action){}

// RegisterType registers a new type of action
func RegisterType(name string, initFunc func() flows.Action) {
	registeredTypes[name] = initFunc
}

// RegisteredTypes gets the registered types of action
func RegisteredTypes() map[string](func() flows.Action) {
	return registeredTypes
}

var uuidRegex = regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

// BaseAction is our base action
type BaseAction struct {
	Type_ string           `json:"type" validate:"required"`
	UUID_ flows.ActionUUID `json:"uuid" validate:"required,uuid4"`
}

func NewBaseAction(typeName string, uuid flows.ActionUUID) BaseAction {
	return BaseAction{Type_: typeName, UUID_: uuid}
}

// Type returns the type of this action
func (a *BaseAction) Type() string { return a.Type_ }

// UUID returns the UUID of the action
func (a *BaseAction) UUID() flows.ActionUUID { return a.UUID_ }

// LocalizationUUID gets the UUID which identifies this object for localization
func (a *BaseAction) LocalizationUUID() utils.UUID { return utils.UUID(a.UUID_) }

// EnumerateTemplates enumerates all expressions on this object and its children
func (a *BaseAction) EnumerateTemplates(localization flows.Localization, callback func(string)) {}

// RewriteTemplates rewrites all templates on this object and its children
func (a *BaseAction) RewriteTemplates(localization flows.Localization, rewrite func(string) string) {}

// EnumerateDependencies enumerates all dependencies on this object and its children
func (a *BaseAction) EnumerateDependencies(localization flows.Localization, callback func(assets.Reference)) {}

// helper function for actions that have a set of group references that must be validated
func (a *BaseAction) validateGroups(assets flows.SessionAssets, references []*assets.GroupReference) error {
	for _, ref := range references {
		if ref.UUID != "" {
			if _, err := assets.Groups().Get(ref.UUID); err != nil {
				return err
			}
		}
	}
	return nil
}

// helper function for actions that have a set of label references that must be validated
func (a *BaseAction) validateLabels(assets flows.SessionAssets, references []*assets.LabelReference) error {
	for _, ref := range references {
		if ref.UUID != "" {
			if _, err := assets.Labels().Get(ref.UUID); err != nil {
				return err
			}
		}
	}
	return nil
}

// helper function for actions that have a flow reference that must be validated
func (a *BaseAction) validateFlow(assets flows.SessionAssets, reference *assets.FlowReference, context *flows.ValidationContext) error {
	// check the flow exists
	flow, err := assets.Flows().Get(reference.UUID)
	if err != nil {
		return err
	}

	// and that it's valid
	return flow.Validate(assets, context)
}

// helper function for actions that have a set of group references that must be resolved to actual groups
func (a *BaseAction) resolveGroups(run flows.FlowRun, references []*assets.GroupReference, staticOnly bool, logEvent flows.EventCallback) ([]*flows.Group, error) {
	groupSet := run.Session().Assets().Groups()
	groups := make([]*flows.Group, 0, len(references))

	for _, ref := range references {
		var group *flows.Group
		var err error

		if ref.UUID != "" {
			// group is a fixed group with a UUID
			group, err = groupSet.Get(ref.UUID)
			if err != nil {
				return nil, err
			}
		} else {
			// group is an expression that evaluates to an existing group's name
			evaluatedGroupName, err := run.EvaluateTemplate(ref.NameMatch)
			if err != nil {
				logEvent(events.NewErrorEvent(err))
			} else {
				// look up the set of all groups to see if such a group exists
				group = groupSet.FindByName(evaluatedGroupName)
				if group == nil {
					logEvent(events.NewErrorEventf("no such group with name '%s'", evaluatedGroupName))
				}
			}
		}

		if group != nil {
			if staticOnly && group.IsDynamic() {
				logEvent(events.NewErrorEventf("can't add or remove contacts from a dynamic group '%s'", group.Name()))
			} else {
				groups = append(groups, group)
			}
		}
	}

	return groups, nil
}

// helper function for actions that have a set of label references that must be resolved to actual labels
func (a *BaseAction) resolveLabels(run flows.FlowRun, references []*assets.LabelReference, logEvent flows.EventCallback) ([]*flows.Label, error) {
	labelSet := run.Session().Assets().Labels()
	labels := make([]*flows.Label, 0, len(references))

	for _, ref := range references {
		var label *flows.Label
		var err error

		if ref.UUID != "" {
			// label is a fixed label with a UUID
			label, err = labelSet.Get(ref.UUID)
			if err != nil {
				return nil, err
			}
		} else {
			// label is an expression that evaluates to an existing label's name
			evaluatedLabelName, err := run.EvaluateTemplate(ref.NameMatch)
			if err != nil {
				logEvent(events.NewErrorEvent(err))
			} else {
				// look up the set of all labels to see if such a label exists
				label = labelSet.FindByName(evaluatedLabelName)
				if label == nil {
					logEvent(events.NewErrorEventf("no such label with name '%s'", evaluatedLabelName))
				}
			}
		}

		if label != nil {
			labels = append(labels, label)
		}
	}

	return labels, nil
}

// helper function for actions that send a message (text + attachments) that must be localized and evalulated
func (a *BaseAction) evaluateMessage(run flows.FlowRun, languages []utils.Language, actionText string, actionAttachments []string, actionQuickReplies []string, logEvent flows.EventCallback) (string, []flows.Attachment, []string) {
	// localize and evaluate the message text
	localizedText := run.GetTranslatedTextArray(utils.UUID(a.UUID()), "text", []string{actionText}, languages)[0]
	evaluatedText, err := run.EvaluateTemplate(localizedText)
	if err != nil {
		logEvent(events.NewErrorEvent(err))
	}

	// localize and evaluate the message attachments
	translatedAttachments := run.GetTranslatedTextArray(utils.UUID(a.UUID()), "attachments", actionAttachments, languages)
	evaluatedAttachments := make([]flows.Attachment, 0, len(translatedAttachments))
	for n := range translatedAttachments {
		evaluatedAttachment, err := run.EvaluateTemplate(translatedAttachments[n])
		if err != nil {
			logEvent(events.NewErrorEvent(err))
		} else if evaluatedAttachment == "" {
			logEvent(events.NewErrorEventf("attachment text evaluated to empty string, skipping"))
			continue
		}
		evaluatedAttachments = append(evaluatedAttachments, flows.Attachment(evaluatedAttachment))
	}

	// localize and evaluate the quick replies
	translatedQuickReplies := run.GetTranslatedTextArray(utils.UUID(a.UUID()), "quick_replies", actionQuickReplies, languages)
	evaluatedQuickReplies := make([]string, 0, len(translatedQuickReplies))
	for n := range translatedQuickReplies {
		evaluatedQuickReply, err := run.EvaluateTemplate(translatedQuickReplies[n])
		if err != nil {
			logEvent(events.NewErrorEvent(err))
		} else if evaluatedQuickReply == "" {
			logEvent(events.NewErrorEventf("quick reply text evaluated to empty string, skipping"))
			continue
		}
		evaluatedQuickReplies = append(evaluatedQuickReplies, evaluatedQuickReply)
	}

	return evaluatedText, evaluatedAttachments, evaluatedQuickReplies
}

func (a *BaseAction) resolveContactsAndGroups(run flows.FlowRun, actionURNs []urns.URN, actionContacts []*flows.ContactReference, actionGroups []*assets.GroupReference, actionLegacyVars []string, logEvent flows.EventCallback) ([]urns.URN, []*flows.ContactReference, []*assets.GroupReference, error) {
	groupSet := run.Session().Assets().Groups()

	// copy URNs
	urnList := make([]urns.URN, 0, len(actionURNs))
	for _, urn := range actionURNs {
		urnList = append(urnList, urn)
	}

	// copy contact references
	contactRefs := make([]*flows.ContactReference, 0, len(actionContacts))
	for _, contactRef := range actionContacts {
		contactRefs = append(contactRefs, contactRef)
	}

	// resolve group references
	groups, err := a.resolveGroups(run, actionGroups, false, logEvent)
	if err != nil {
		return nil, nil, nil, err
	}
	groupRefs := make([]*assets.GroupReference, 0, len(groups))
	for _, group := range groups {
		groupRefs = append(groupRefs, group.Reference())
	}

	// evaluate the legacy variables
	for _, legacyVar := range actionLegacyVars {
		evaluatedLegacyVar, err := run.EvaluateTemplate(legacyVar)
		if err != nil {
			logEvent(events.NewErrorEvent(err))
		}

		if uuidRegex.MatchString(evaluatedLegacyVar) {
			// if variable evaluates to a UUID, we assume it's a contact UUID
			contactRefs = append(contactRefs, flows.NewContactReference(flows.ContactUUID(evaluatedLegacyVar), ""))

		} else if groupByName := groupSet.FindByName(evaluatedLegacyVar); groupByName != nil {
			// next up we look for a group with a matching name
			groupRefs = append(groupRefs, groupByName.Reference())
		} else {
			// if that fails, assume this is a phone number, and let the caller worry about validation
			urn, err := urns.NewURNFromParts(urns.TelScheme, evaluatedLegacyVar, "", "")
			if err != nil {
				logEvent(events.NewErrorEvent(err))
			} else {
				urnList = append(urnList, urn)
			}
		}
	}

	return urnList, contactRefs, groupRefs, nil
}

// helper to save a run result and log it as an event
func (a *BaseAction) saveResult(run flows.FlowRun, step flows.Step, name, value, category, categoryLocalized string, input *string, extra json.RawMessage, logEvent flows.EventCallback) {
	result := flows.NewResult(name, value, category, categoryLocalized, step.NodeUUID(), input, extra, utils.Now())
	run.SaveResult(result)
	logEvent(events.NewRunResultChangedEvent(result))
}

// helper to save a run result based on a webhook call and log it as an event
func (a *BaseAction) saveWebhookResult(run flows.FlowRun, step flows.Step, name string, webhook *flows.WebhookCall, logEvent flows.EventCallback) {
	input := fmt.Sprintf("%s %s", webhook.Method(), webhook.URL())
	value := strconv.Itoa(webhook.StatusCode())
	category := webhookStatusCategories[webhook.Status()]

	body := []byte(webhook.Body())
	var extra json.RawMessage

	// try to parse body as JSON
	if utils.IsValidJSON(body) {
		// if that was successful, the body is valid JSON and extra is the body
		extra = body
	} else {
		// if not, treat body as text and encode as a JSON string
		extra, _ = json.Marshal(string(body))
	}

	a.saveResult(run, step, name, value, category, "", &input, extra, logEvent)
}

// helper to apply a contact modifier
func (a *BaseAction) applyModifier(run flows.FlowRun, mod flows.Modifier, logModifier flows.ModifierCallback, logEvent flows.EventCallback) {
	mod.Apply(run.Session().Environment(), run.Session().Assets(), run.Contact(), logEvent)
	logModifier(mod)
}

// utility struct which sets the allowed flow types to any
type universalAction struct{}

// AllowedFlowTypes returns the flow types which this action is allowed to occur in
func (a *universalAction) AllowedFlowTypes() []flows.FlowType {
	return []flows.FlowType{flows.FlowTypeMessaging, flows.FlowTypeMessagingOffline, flows.FlowTypeVoice}
}

// utility struct which sets the allowed flow types to any which run online
type onlineAction struct{}

// AllowedFlowTypes returns the flow types which this action is allowed to occur in
func (a *onlineAction) AllowedFlowTypes() []flows.FlowType {
	return []flows.FlowType{flows.FlowTypeMessaging, flows.FlowTypeVoice}
}

// utility struct which sets the allowed flow types to just voice
type voiceAction struct{}

// AllowedFlowTypes returns the flow types which this action is allowed to occur in
func (a *voiceAction) AllowedFlowTypes() []flows.FlowType {
	return []flows.FlowType{flows.FlowTypeVoice}
}

// utility struct for actions which operate on other contacts
type otherContactsAction struct {
	URNs       []urns.URN                `json:"urns,omitempty"`
	Contacts   []*flows.ContactReference `json:"contacts,omitempty" validate:"dive"`
	Groups     []*assets.GroupReference  `json:"groups,omitempty" validate:"dive"`
	LegacyVars []string                  `json:"legacy_vars,omitempty"`
}

// utility struct for actions which create a message
type createMsgAction struct {
	Text         string   `json:"text"`
	Attachments  []string `json:"attachments,omitempty"`
	QuickReplies []string `json:"quick_replies,omitempty"`
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
