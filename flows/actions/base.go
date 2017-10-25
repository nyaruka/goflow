package actions

import (
	"fmt"
	"regexp"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

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

// helper function for actions that have a set of group references that must be validated
func (a *BaseAction) validateGroups(assets flows.SessionAssets, references []*flows.GroupReference) error {
	for _, ref := range references {
		if ref.UUID != "" {
			if _, err := assets.GetGroup(ref.UUID); err != nil {
				return err
			}
		}
	}
	return nil
}

// helper function for actions that have a set of label references that must be validated
func (a *BaseAction) validateLabels(assets flows.SessionAssets, references []*flows.LabelReference) error {
	for _, ref := range references {
		if ref.UUID != "" {
			if _, err := assets.GetLabel(ref.UUID); err != nil {
				return err
			}
		}
	}
	return nil
}

// helper function for actions that have a set of group references that must be resolved to actual groups
func (a *BaseAction) resolveGroups(run flows.FlowRun, step flows.Step, references []*flows.GroupReference, log flows.EventLog) ([]*flows.Group, error) {
	groupSet, err := run.Session().Assets().GetGroupSet()
	if err != nil {
		return nil, err
	}

	groups := make([]*flows.Group, 0, len(references))

	for _, ref := range references {
		var group *flows.Group

		if ref.UUID != "" {
			// group is a fixed group with a UUID
			group = groupSet.FindByUUID(ref.UUID)
			if group == nil {
				return nil, fmt.Errorf("no such group with UUID '%s'", ref.UUID)
			}
		} else {
			// group is an expression that evaluates to an existing group's name
			evaluatedGroupName, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), ref.NameMatch)
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
func (a *BaseAction) resolveLabels(run flows.FlowRun, step flows.Step, references []*flows.LabelReference, log flows.EventLog) ([]*flows.Label, error) {
	labelSet, err := run.Session().Assets().GetLabelSet()
	if err != nil {
		return nil, err
	}

	labels := make([]*flows.Label, 0, len(references))

	for _, ref := range references {
		var label *flows.Label

		if ref.UUID != "" {
			// label is a fixed label with a UUID
			label = labelSet.FindByUUID(ref.UUID)
			if label == nil {
				return nil, fmt.Errorf("no such label with UUID '%s'", ref.UUID)
			}
		} else {
			// label is an expression that evaluates to an existing label's name
			evaluatedLabelName, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), ref.NameMatch)
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

// MsgAction is our mixin for an action that generates message events
type MsgAction struct {
	Text        string   `json:"text"`
	Attachments []string `json:"attachments"`
}

func (a *MsgAction) evaluateMessage(b *BaseAction, run flows.FlowRun, step flows.Step, log flows.EventLog) (string, []string) {
	// localize and evaluate the message text
	localizedText := run.GetText(flows.UUID(b.UUID()), "text", a.Text)
	evaluatedText, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), localizedText)
	if err != nil {
		log.Add(events.NewErrorEvent(err))
	}

	// localize and evaluate the message attachments
	translatedAttachments := run.GetTextArray(flows.UUID(b.UUID()), "attachments", a.Attachments)
	evaluatedAttachments := make([]string, 0, len(a.Attachments))
	for n := range translatedAttachments {
		evaluatedAttachment, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), translatedAttachments[n])
		if err != nil {
			log.Add(events.NewErrorEvent(err))
		} else if evaluatedAttachment == "" {
			log.Add(events.NewErrorEvent(fmt.Errorf("attachment text evaluated to empty string, skipping")))
			continue
		}
		evaluatedAttachments = append(evaluatedAttachments, evaluatedAttachment)
	}

	return evaluatedText, evaluatedAttachments
}

// ContactsAndGroupsAction is our mixin for an action that operates on other contacts and groups
type ContactsAndGroupsAction struct {
	URNs       []urns.URN                `json:"urns,omitempty"`
	Contacts   []*flows.ContactReference `json:"contacts,omitempty" validate:"dive"`
	Groups     []*flows.GroupReference   `json:"groups,omitempty" validate:"dive"`
	LegacyVars []string                  `json:"legacy_vars,omitempty"`
}

func (a *ContactsAndGroupsAction) resolveContactsAndGroups(b *BaseAction, run flows.FlowRun, step flows.Step, log flows.EventLog) ([]urns.URN, []*flows.ContactReference, []*flows.GroupReference, error) {
	// copy URNs
	urnList := make([]urns.URN, 0, len(a.URNs))
	for _, urn := range a.URNs {
		urnList = append(urnList, urn)
	}

	// copy contact references
	contactRefs := make([]*flows.ContactReference, 0, len(a.Contacts))
	for _, contactRef := range a.Contacts {
		contactRefs = append(contactRefs, contactRef)
	}

	// resolve group references
	groups, err := b.resolveGroups(run, step, a.Groups, log)
	if err != nil {
		return nil, nil, nil, err
	}
	groupRefs := make([]*flows.GroupReference, 0, len(groups))
	for _, group := range groups {
		groupRefs = append(groupRefs, group.Reference())
	}

	// get the list of all groups
	groupSet, err := run.Session().Assets().GetGroupSet()
	if err != nil {
		return nil, nil, nil, err
	}

	// evaluate the legacy variables
	for _, legacyVar := range a.LegacyVars {
		evaluatedLegacyVar, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), legacyVar)
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
			urnList = append(urnList, urns.NewURNFromParts(urns.TelScheme, evaluatedLegacyVar, ""))
		}
	}

	return urnList, contactRefs, groupRefs, nil
}
