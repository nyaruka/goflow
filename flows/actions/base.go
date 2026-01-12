package actions

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/stringsx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/modifiers"
	"github.com/nyaruka/goflow/utils"
)

// max number of bytes to be saved to extra on a result
const resultExtraMaxBytes = 10000

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
	UUID_ flows.ActionUUID `json:"uuid" validate:"required,uuid"`
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

func (a *baseAction) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	// nothing to declare
}

// LocalizationUUID gets the UUID which identifies this object for localization
func (a *baseAction) LocalizationUUID() uuids.UUID { return uuids.UUID(a.UUID_) }

// helper function for actions that send a message (text + attachments) that must be localized and evalulated
func (a *baseAction) evaluateMessage(run flows.Run, languages []i18n.Language, actionText string, actionAttachments []string, actionQuickReplies []string, log flows.EventLogger) (*flows.MsgContent, i18n.Language) {
	// localize and evaluate the message text
	localizedText, txtLang := run.GetTextArray(uuids.UUID(a.UUID()), "text", []string{actionText}, languages)
	evaluatedText, _ := run.EvaluateTemplate(localizedText[0], log)

	// localize and evaluate the message attachments
	translatedAttachments, attLang := run.GetTextArray(uuids.UUID(a.UUID()), "attachments", actionAttachments, languages)
	evaluatedAttachments := make([]utils.Attachment, 0, len(translatedAttachments))
	for _, a := range translatedAttachments {
		evaluatedAttachment, _ := run.EvaluateTemplate(a, log)
		evaluatedAttachment = strings.TrimSpace(evaluatedAttachment)
		if !utils.IsValidAttachment(evaluatedAttachment) {
			log(events.NewError("attachment evaluated to invalid value, skipping", ""))
			continue
		}
		if len(evaluatedAttachment) > flows.MaxAttachmentLength {
			log(events.NewError(fmt.Sprintf("evaluated attachment is longer than %d limit, skipping", flows.MaxAttachmentLength), ""))
			continue
		}
		evaluatedAttachments = append(evaluatedAttachments, utils.Attachment(evaluatedAttachment))
	}

	// localize and evaluate the quick replies
	translatedQuickReplies, qrsLang := run.GetTextArray(uuids.UUID(a.UUID()), "quick_replies", actionQuickReplies, languages)
	evaluatedQuickReplies := make([]flows.QuickReply, 0, len(translatedQuickReplies))
	for _, qr := range translatedQuickReplies {
		evaluatedQuickReply, _ := run.EvaluateTemplate(qr, log)
		if evaluatedQuickReply == "" {
			log(events.NewError("quick reply evaluated to empty string, skipping", ""))
			continue
		}
		evaluatedQuickReplies = append(evaluatedQuickReplies, flows.QuickReply{Text: stringsx.TruncateEllipsis(evaluatedQuickReply, flows.MaxQuickReplyLength)})
	}

	// although it's possible for the different parts of the message to have different languages, we want to resolve
	// a single language based on what the user actually provided for this message
	var lang i18n.Language
	if localizedText[0] != "" {
		lang = txtLang
	} else if len(translatedAttachments) > 0 {
		lang = attLang
	} else if len(translatedQuickReplies) > 0 {
		lang = qrsLang
	}

	return &flows.MsgContent{Text: evaluatedText, Attachments: evaluatedAttachments, QuickReplies: evaluatedQuickReplies}, lang
}

// helper to save a run result and log it as an event
func (a *baseAction) saveResult(run flows.Run, step flows.Step, name, value, category, categoryLocalized string, input string, extra []byte, log flows.EventLogger) {
	result := flows.NewResult(name, value, category, categoryLocalized, step.NodeUUID(), input, extra, dates.Now())
	prev, changed := run.SetResult(result)
	if changed {
		log(events.NewRunResultChanged(result, prev))
	}
}

// helper to save a run result based on a webhook call and log it as an event.. new webhook nodes don't use this
func (a *baseAction) saveLegacyWebhookResult(run flows.Run, step flows.Step, name string, call *flows.WebhookCall, status flows.CallStatus, log flows.EventLogger) {
	input := fmt.Sprintf("%s %s", call.Method, call.URL)
	value := strconv.Itoa(call.ResponseStatus)
	category := webhookStatusCategories[status]

	var extra []byte
	if len(call.ResponseJSON) > 0 && len(call.ResponseJSON) < resultExtraMaxBytes {
		extra = call.ResponseJSON
	}

	a.saveResult(run, step, name, value, category, "", input, extra, log)
}

// helper to apply a contact modifier
func (a *baseAction) applyModifier(ctx context.Context, run flows.Run, mod flows.Modifier, log flows.EventLogger) (bool, error) {
	s := run.Session()
	return modifiers.Apply(ctx, s.Engine(), s.MergedEnvironment(), s.Assets(), run.Contact(), mod, log)
}

// helper to log a failure
func (a *baseAction) fail(run flows.Run, err error, log flows.EventLogger) {
	run.Exit(flows.RunStatusFailed)
	log(events.NewFailure(err))
	log(events.NewRunEnded(run.UUID(), run.FlowReference(), flows.RunStatusFailed))
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
	Groups       []*assets.GroupReference  `json:"groups,omitempty" validate:"dive"`
	Contacts     []*flows.ContactReference `json:"contacts,omitempty" validate:"dive"`
	ContactQuery string                    `json:"contact_query,omitempty" engine:"evaluated"`
	URNs         []urns.URN                `json:"urns,omitempty"`
	LegacyVars   []string                  `json:"legacy_vars,omitempty" engine:"evaluated"`
}

func (a *otherContactsAction) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	for _, group := range a.Groups {
		dependency(group)
	}
	for _, contact := range a.Contacts {
		dependency(contact)
	}
}

func (a *otherContactsAction) resolveRecipients(run flows.Run, log flows.EventLogger) ([]*assets.GroupReference, []*flows.ContactReference, string, []urns.URN, error) {
	groupSet := run.Session().Assets().Groups()

	// copy URNs
	urnList := make([]urns.URN, 0, len(a.URNs))
	urnList = append(urnList, a.URNs...)

	// copy contact references
	contactRefs := make([]*flows.ContactReference, 0, len(a.Contacts))
	contactRefs = append(contactRefs, a.Contacts...)

	// resolve group references
	groups := resolveGroups(run, a.Groups, log)
	groupRefs := make([]*assets.GroupReference, 0, len(groups))
	for _, group := range groups {
		groupRefs = append(groupRefs, group.Reference())
	}

	// evaluate the legacy variables
	for _, legacyVar := range a.LegacyVars {
		evaluatedLegacyVar, _ := run.EvaluateTemplate(legacyVar, log)

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
				urnList = append(urnList, urn.Normalize())
			} else {
				// if that fails, try to parse as phone number
				parsedTel := utils.ParsePhoneNumber(evaluatedLegacyVar, run.Session().MergedEnvironment().DefaultCountry())
				if parsedTel != "" {
					urn, _ := urns.New(urns.Phone, parsedTel)
					urnList = append(urnList, urn)
				} else {
					log(events.NewError(fmt.Sprintf("'%s' couldn't be resolved to a contact, group or URN", evaluatedLegacyVar), ""))
				}
			}
		}
	}

	// evaluate contact query
	contactQuery, _ := run.EvaluateTemplateText(a.ContactQuery, flows.ContactQueryEscaping, true, log)
	contactQuery = strings.TrimSpace(contactQuery)

	return groupRefs, contactRefs, contactQuery, urnList, nil
}

// utility struct for actions which create a message
type createMsgAction struct {
	Text         string   `json:"text"                    validate:"required,max=10000"     engine:"localized,evaluated"`
	Attachments  []string `json:"attachments,omitempty"   validate:"max=10,dive,attachment" engine:"localized,evaluated"`
	QuickReplies []string `json:"quick_replies,omitempty" validate:"max=10,dive,max=1000"   engine:"localized,evaluated"`
}

// helper function for actions that have a set of group references that must be resolved to actual groups
func resolveGroups(run flows.Run, references []*assets.GroupReference, log flows.EventLogger) []*flows.Group {
	groupAssets := run.Session().Assets().Groups()
	groups := make([]*flows.Group, 0, len(references))

	for _, ref := range references {
		var group *flows.Group

		if ref.Variable() {
			// is an expression that evaluates to an existing group's name
			evaluatedName, ok := run.EvaluateTemplate(ref.NameMatch, log)
			if ok {

				// look up the set of all groups to see if such a group exists
				group = groupAssets.FindByName(evaluatedName)
				if group == nil {
					log(events.NewError(fmt.Sprintf("no such group with name '%s'", evaluatedName), ""))
				}
			}
		} else {
			// group is a fixed group with a UUID
			group = groupAssets.Get(ref.UUID)
			if group == nil {
				log(events.NewDependencyError(ref))
			}
		}

		if group != nil {
			groups = append(groups, group)
		}
	}

	return groups
}

// helper function for actions that have a set of label references that must be resolved to actual labels
func resolveLabels(run flows.Run, references []*assets.LabelReference, log flows.EventLogger) []*flows.Label {
	labelAssets := run.Session().Assets().Labels()
	labels := make([]*flows.Label, 0, len(references))

	for _, ref := range references {
		var label *flows.Label

		if ref.Variable() {
			// is an expression that evaluates to an existing label's name
			evaluatedName, ok := run.EvaluateTemplate(ref.NameMatch, log)
			if ok {
				// look up the set of all labels to see if such a label exists
				label = labelAssets.FindByName(evaluatedName)
				if label == nil {
					log(events.NewError(fmt.Sprintf("no such label with name '%s'", evaluatedName), ""))
				}
			}
		} else {
			// label is a fixed label with a UUID
			label = labelAssets.Get(ref.UUID)
			if label == nil {
				log(events.NewDependencyError(ref))
			}
		}

		if label != nil {
			labels = append(labels, label)
		}
	}

	return labels
}

// helper function to resolve a user reference to a user
func resolveUser(run flows.Run, ref *assets.UserReference, log flows.EventLogger) *flows.User {
	userAssets := run.Session().Assets().Users()
	var user *flows.User

	if ref.Variable() {
		// is an expression that evaluates to an existing user's email
		evaluatedEmail, ok := run.EvaluateTemplate(ref.EmailMatch, log)
		if ok {
			// look up to see if such a user exists
			user = userAssets.FindByEmail(evaluatedEmail)
			if user == nil {
				log(events.NewError(fmt.Sprintf("no such user with email '%s'", evaluatedEmail), ""))
			}
		}
	} else {
		// user is a fixed user with this UUID
		user = userAssets.Get(ref.UUID)
		if user == nil {
			log(events.NewDependencyError(ref))
		}
	}

	return user
}

func currentLocale(run flows.Run, lang i18n.Language) i18n.Locale {
	return i18n.NewLocale(lang, run.Session().MergedEnvironment().DefaultCountry())
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// Read reads an action from the given JSON
func Read(data []byte) (flows.Action, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, fmt.Errorf("unknown type: '%s'", typeName)
	}

	action := f()
	return action, utils.UnmarshalAndValidate(data, action)
}
