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
)

var webhookStatusCategories = map[flows.WebhookStatus]string{
	flows.WebhookStatusSuccess:         "Success",
	flows.WebhookStatusResponseError:   "Failure",
	flows.WebhookStatusConnectionError: "Failure",
}

var registeredTypes = map[string](func() flows.Action){}

// RegisterType registers a new type of router
func RegisterType(name string, initFunc func() flows.Action) {
	registeredTypes[name] = initFunc
}

var uuidRegex = regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

type eventLog struct {
	events []flows.Event
}

func NewEventLog() flows.EventLog {
	return &eventLog{events: make([]flows.Event, 0)}
}

func (l *eventLog) Events() []flows.Event { return l.events }

func (l *eventLog) Add(event flows.Event) {
	l.events = append(l.events, event)
}

// BaseAction is our base action
type BaseAction struct {
	UUID_ flows.ActionUUID `json:"uuid" validate:"required,uuid4"`
}

func NewBaseAction(uuid flows.ActionUUID) BaseAction {
	return BaseAction{UUID_: uuid}
}

// UUID returns the UUID of the action
func (a *BaseAction) UUID() flows.ActionUUID { return a.UUID_ }

func (a *BaseAction) evaluateLocalizableTemplate(run flows.FlowRun, localizationKey string, defaultValue string) (string, error) {
	localizedTemplate := run.GetText(utils.UUID(a.UUID()), localizationKey, defaultValue)
	return run.EvaluateTemplateAsString(localizedTemplate, false)
}

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

// helper function for actions that have a set of group references that must be resolved to actual groups
func (a *BaseAction) resolveGroups(run flows.FlowRun, step flows.Step, references []*assets.GroupReference, log flows.EventLog) ([]*flows.Group, error) {
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
			evaluatedGroupName, err := run.EvaluateTemplateAsString(ref.NameMatch, false)
			if err != nil {
				log.Add(events.NewErrorEvent(err))
			} else {
				// look up the set of all groups to see if such a group exists
				group = groupSet.FindByName(evaluatedGroupName)
				if group == nil {
					log.Add(events.NewErrorEvent(fmt.Errorf("no such group with name '%s'", evaluatedGroupName)))
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
func (a *BaseAction) resolveLabels(run flows.FlowRun, step flows.Step, references []*assets.LabelReference, log flows.EventLog) ([]*flows.Label, error) {
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
			evaluatedLabelName, err := run.EvaluateTemplateAsString(ref.NameMatch, false)
			if err != nil {
				log.Add(events.NewErrorEvent(err))
			} else {
				// look up the set of all labels to see if such a label exists
				label = labelSet.FindByName(evaluatedLabelName)
				if label == nil {
					log.Add(events.NewErrorEvent(fmt.Errorf("no such label with name '%s'", evaluatedLabelName)))
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
func (a *BaseAction) evaluateMessage(run flows.FlowRun, languages utils.LanguageList, actionText string, actionAttachments []string, actionQuickReplies []string, log flows.EventLog) (string, []flows.Attachment, []string) {
	// localize and evaluate the message text
	localizedText := run.GetTranslatedTextArray(utils.UUID(a.UUID()), "text", []string{actionText}, languages)[0]
	evaluatedText, err := run.EvaluateTemplateAsString(localizedText, false)
	if err != nil {
		log.Add(events.NewErrorEvent(err))
	}

	// localize and evaluate the message attachments
	translatedAttachments := run.GetTranslatedTextArray(utils.UUID(a.UUID()), "attachments", actionAttachments, languages)
	evaluatedAttachments := make([]flows.Attachment, 0, len(translatedAttachments))
	for n := range translatedAttachments {
		evaluatedAttachment, err := run.EvaluateTemplateAsString(translatedAttachments[n], true)
		if err != nil {
			log.Add(events.NewErrorEvent(err))
		} else if evaluatedAttachment == "" {
			log.Add(events.NewErrorEvent(fmt.Errorf("attachment text evaluated to empty string, skipping")))
			continue
		}
		evaluatedAttachments = append(evaluatedAttachments, flows.Attachment(evaluatedAttachment))
	}

	// localize and evaluate the quick replies
	translatedQuickReplies := run.GetTranslatedTextArray(utils.UUID(a.UUID()), "quick_replies", actionQuickReplies, languages)
	evaluatedQuickReplies := make([]string, 0, len(translatedQuickReplies))
	for n := range translatedQuickReplies {
		evaluatedQuickReply, err := run.EvaluateTemplateAsString(translatedQuickReplies[n], false)
		if err != nil {
			log.Add(events.NewErrorEvent(err))
		} else if evaluatedQuickReply == "" {
			log.Add(events.NewErrorEvent(fmt.Errorf("quick reply text evaluated to empty string, skipping")))
			continue
		}
		evaluatedQuickReplies = append(evaluatedQuickReplies, evaluatedQuickReply)
	}

	return evaluatedText, evaluatedAttachments, evaluatedQuickReplies
}

func (a *BaseAction) resolveContactsAndGroups(run flows.FlowRun, step flows.Step, actionURNs []urns.URN, actionContacts []*flows.ContactReference, actionGroups []*assets.GroupReference, actionLegacyVars []string, log flows.EventLog) ([]urns.URN, []*flows.ContactReference, []*assets.GroupReference, error) {
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
	groups, err := a.resolveGroups(run, step, actionGroups, log)
	if err != nil {
		return nil, nil, nil, err
	}
	groupRefs := make([]*assets.GroupReference, 0, len(groups))
	for _, group := range groups {
		groupRefs = append(groupRefs, group.Reference())
	}

	// evaluate the legacy variables
	for _, legacyVar := range actionLegacyVars {
		evaluatedLegacyVar, err := run.EvaluateTemplateAsString(legacyVar, false)
		if err != nil {
			log.Add(events.NewErrorEvent(err))
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
				return nil, nil, nil, err
			}
			urnList = append(urnList, urn)
		}
	}

	return urnList, contactRefs, groupRefs, nil
}

// saves a run result and returns the corresponding event
func (a *BaseAction) saveResult(run flows.FlowRun, step flows.Step, name, value, category, categoryLocalized string, input *string, extra json.RawMessage, log flows.EventLog) {
	run.Results().Save(name, value, category, categoryLocalized, step.NodeUUID(), input, extra, utils.Now())
	log.Add(events.NewRunResultChangedEvent(name, value, category, categoryLocalized, input, extra))
}

// saves a run result based on a webhook call and returns the corresponding event
func (a *BaseAction) saveWebhookResult(run flows.FlowRun, step flows.Step, name string, webhook *flows.WebhookCall, log flows.EventLog) {
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

	a.saveResult(run, step, name, value, category, "", &input, extra, log)
}

func (a *BaseAction) reevaluateDynamicGroups(run flows.FlowRun, log flows.EventLog) {
	added, removed, errors := run.Contact().ReevaluateDynamicGroups(run.Session())

	// add error event for each group we couldn't re-evaluate
	for _, err := range errors {
		log.Add(events.NewErrorEvent(err))
	}

	// add groups changed event for the groups we were added/removed to/from
	if len(added) > 0 || len(removed) > 0 {
		log.Add(events.NewContactGroupsChangedEvent(added, removed))
	}
}

func (a *BaseAction) logError(err error, log flows.EventLog) {
	log.Add(events.NewErrorEvent(err))
}

func (a *BaseAction) logFatalError(err error, log flows.EventLog) {
	log.Add(events.NewFatalErrorEvent(err))
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

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadAction reads an action from the given typed envelope
func ReadAction(envelope *utils.TypedEnvelope) (flows.Action, error) {
	f := registeredTypes[envelope.Type]
	if f == nil {
		return nil, fmt.Errorf("unknown type: %s", envelope.Type)
	}

	action := f()
	return action, utils.UnmarshalAndValidate(envelope.Data, action)
}
