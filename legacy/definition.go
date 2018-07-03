package legacy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/routers"
	"github.com/nyaruka/goflow/flows/waits"
	"github.com/nyaruka/goflow/legacy/expressions"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

var legacyWebhookBody = `{
	"contact": {"uuid": "@contact.uuid", "name": @(json(contact.name)), "urn": @(json(if(default(run.input.urn, default(contact.urns.0, null)), text(default(run.input.urn, default(contact.urns.0, null))), null)))},
	"flow": @(json(run.flow)),
	"path": @(json(run.path)),
	"results": @(json(run.results)),
	"run": {"uuid": "@run.uuid", "created_on": "@run.created_on"},
	"input": @(json(run.input)),
	"channel": @(json(if(run.input, run.input.channel, null)))
}`

// Flow is a flow in the legacy format
type Flow struct {
	BaseLanguage utils.Language `json:"base_language"`
	Metadata     Metadata       `json:"metadata"`
	RuleSets     []RuleSet      `json:"rule_sets" validate:"dive"`
	ActionSets   []ActionSet    `json:"action_sets" validate:"dive"`
	Entry        flows.NodeUUID `json:"entry" validate:"required,uuid4"`
}

// Note is a legacy sticky note
type Note struct {
	X     decimal.Decimal `json:"x"`
	Y     decimal.Decimal `json:"y"`
	Title string          `json:"title"`
	Body  string          `json:"body"`
}

// Sticky is a migrated note
type Sticky map[string]interface{}

// Migrate migrates this note to a new sticky note
func (n *Note) Migrate() Sticky {
	return Sticky{
		"position": map[string]interface{}{"left": n.X.IntPart(), "top": n.Y.IntPart()},
		"title":    n.Title,
		"body":     n.Body,
		"color":    "yellow",
	}
}

// Metadata is the metadata section of a legacy flow
type Metadata struct {
	UUID     flows.FlowUUID `json:"uuid" validate:"required,uuid4"`
	Name     string         `json:"name"`
	Revision int            `json:"revision"`
	Expires  int            `json:"expires"`
	Notes    []Note         `json:"notes,omitempty"`
}

type Rule struct {
	UUID            flows.ExitUUID      `json:"uuid" validate:"required,uuid4"`
	Destination     flows.NodeUUID      `json:"destination" validate:"omitempty,uuid4"`
	DestinationType string              `json:"destination_type" validate:"eq=A|eq=R"`
	Test            utils.TypedEnvelope `json:"test"`
	Category        Translations        `json:"category"`
}

type RuleSet struct {
	Y       int             `json:"y"`
	X       int             `json:"x"`
	UUID    flows.NodeUUID  `json:"uuid" validate:"required,uuid4"`
	Type    string          `json:"ruleset_type"`
	Label   string          `json:"label"`
	Operand string          `json:"operand"`
	Rules   []Rule          `json:"rules"`
	Config  json.RawMessage `json:"config"`
}

type ActionSet struct {
	Y           int            `json:"y"`
	X           int            `json:"x"`
	Destination flows.NodeUUID `json:"destination" validate:"omitempty,uuid4"`
	ExitUUID    flows.ExitUUID `json:"exit_uuid" validate:"required,uuid4"`
	UUID        flows.NodeUUID `json:"uuid" validate:"required,uuid4"`
	Actions     []Action       `json:"actions"`
}

type LabelReference struct {
	UUID flows.LabelUUID
	Name string
}

func (l *LabelReference) Migrate() *flows.LabelReference {
	if len(l.UUID) > 0 {
		return flows.NewLabelReference(l.UUID, l.Name)
	}
	return flows.NewVariableLabelReference(l.Name)
}

// UnmarshalJSON unmarshals a legacy label reference from the given JSON
func (l *LabelReference) UnmarshalJSON(data []byte) error {
	// label reference may be a string
	if data[0] == '"' {
		var nameExpression string
		if err := json.Unmarshal(data, &nameExpression); err != nil {
			return err
		}

		// if it starts with @ then it's an expression
		if strings.HasPrefix(nameExpression, "@") {
			nameExpression, _ = expressions.MigrateTemplate(nameExpression, expressions.ExtraAsFunction)
		}

		l.Name = nameExpression
		return nil
	}

	// or a JSON object with UUID/Name properties
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	l.UUID = flows.LabelUUID(raw["uuid"].(string))
	l.Name = raw["name"].(string)
	return nil
}

type ContactReference struct {
	UUID flows.ContactUUID `json:"uuid"`
	Name string            `json:"name"`
}

func (c *ContactReference) Migrate() *flows.ContactReference {
	return flows.NewContactReference(c.UUID, c.Name)
}

type GroupReference struct {
	UUID flows.GroupUUID
	Name string
}

func (g *GroupReference) Migrate() *flows.GroupReference {
	if len(g.UUID) > 0 {
		return flows.NewGroupReference(g.UUID, g.Name)
	}
	return flows.NewVariableGroupReference(g.Name)
}

// UnmarshalJSON unmarshals a legacy group reference from the given JSON
func (g *GroupReference) UnmarshalJSON(data []byte) error {
	// group reference may be a string
	if data[0] == '"' {
		var nameExpression string
		if err := json.Unmarshal(data, &nameExpression); err != nil {
			return err
		}

		// if it starts with @ then it's an expression
		if strings.HasPrefix(nameExpression, "@") {
			nameExpression, _ = expressions.MigrateTemplate(nameExpression, expressions.ExtraAsFunction)
		}

		g.Name = nameExpression
		return nil
	}

	// or a JSON object with UUID/Name properties
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	g.UUID = flows.GroupUUID(raw["uuid"].(string))
	g.Name = raw["name"].(string)
	return nil
}

type VariableReference struct {
	ID string `json:"id"`
}

type FlowReference struct {
	UUID flows.FlowUUID `json:"uuid"`
	Name string         `json:"name"`
}

func (f *FlowReference) Migrate() *flows.FlowReference {
	return flows.NewFlowReference(f.UUID, f.Name)
}

type WebhookConfig struct {
	Webhook string          `json:"webhook"`
	Action  string          `json:"webhook_action"`
	Headers []WebhookHeader `json:"webhook_headers"`
}

type WebhookHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Action struct {
	Type string           `json:"type"`
	UUID flows.ActionUUID `json:"uuid"`
	Name string           `json:"name"`

	// message and email
	Msg          json.RawMessage `json:"msg"`
	Media        json.RawMessage `json:"media"`
	QuickReplies json.RawMessage `json:"quick_replies"`
	SendAll      bool            `json:"send_all"`

	// variable contact actions
	Contacts  []ContactReference  `json:"contacts"`
	Groups    []GroupReference    `json:"groups"`
	Variables []VariableReference `json:"variables"`

	// save actions
	Field string `json:"field"`
	Value string `json:"value"`
	Label string `json:"label"`

	// set language
	Language utils.Language `json:"lang"`

	// webhook
	Action         string          `json:"action"`
	Webhook        string          `json:"webhook"`
	WebhookHeaders []WebhookHeader `json:"webhook_headers"`

	// add lable action
	Labels []LabelReference `json:"labels"`

	// Start/Trigger flow
	Flow FlowReference `json:"flow"`

	// channel
	Channel flows.ChannelUUID `json:"channel"`

	//email
	Emails  []string `json:"emails"`
	Subject string   `json:"subject"`
}

type subflowTest struct {
	ExitType string `json:"exit_type"`
}

type webhookTest struct {
	Status string `json:"status"`
}

type localizedStringTest struct {
	Test Translations `json:"test"`
}

type stringTest struct {
	Test string `json:"test"`
}

type numericTest struct {
	Test DecimalString `json:"test"`
}

type betweenTest struct {
	Min string `json:"min"`
	Max string `json:"max"`
}

type timeoutTest struct {
	Minutes int `json:"minutes"`
}

type groupTest struct {
	Test GroupReference `json:"test"`
}

type wardTest struct {
	State    string `json:"state"`
	District string `json:"district"`
}

func addTranslationMap(baseLanguage utils.Language, localization flows.Localization, mapped Translations, uuid utils.UUID, property string) string {
	var inBaseLanguage string
	for language, item := range mapped {
		expression, _ := expressions.MigrateTemplate(item, expressions.ExtraAsFunction)
		if language != baseLanguage && language != "base" {
			localization.AddItemTranslation(language, uuid, property, []string{expression})
		} else {
			inBaseLanguage = expression
		}
	}

	return inBaseLanguage
}

func addTranslationMultiMap(baseLanguage utils.Language, localization flows.Localization, mapped map[utils.Language][]string, uuid utils.UUID, property string) []string {
	var inBaseLanguage []string
	for language, items := range mapped {
		templates := make([]string, len(items))
		for i := range items {
			expression, _ := expressions.MigrateTemplate(items[i], expressions.ExtraAsFunction)
			templates[i] = expression
		}
		if language != baseLanguage {
			localization.AddItemTranslation(language, uuid, property, templates)
		} else {
			inBaseLanguage = templates
		}
	}
	return inBaseLanguage
}

// TransformTranslations transforms a list of single item translations into a map of multi-item translations, e.g.
//
// [{"eng": "yes", "fra": "oui"}, {"eng": "no", "fra": "non"}] becomes {"eng": ["yes", "no"], "fra": ["oui", "non"]}
//
func TransformTranslations(items []Translations) map[utils.Language][]string {
	// re-organize into a map of arrays
	transformed := make(map[utils.Language][]string)

	for i := range items {
		for language, translation := range items[i] {
			perLanguage, found := transformed[language]
			if !found {
				perLanguage = make([]string, len(items))
				transformed[language] = perLanguage
			}
			perLanguage[i] = translation
		}
	}
	return transformed
}

var testTypeMappings = map[string]string{
	"between":              "has_number_between",
	"contains":             "has_all_words",
	"contains_any":         "has_any_word",
	"contains_only_phrase": "has_only_phrase",
	"contains_phrase":      "has_phrase",
	"date":                 "has_date",
	"date_after":           "has_date_gt",
	"date_before":          "has_date_lt",
	"date_equal":           "has_date_eq",
	"district":             "has_district",
	"has_email":            "has_email",
	"eq":                   "has_number_eq",
	"gt":                   "has_number_gt",
	"gte":                  "has_number_gte",
	"in_group":             "has_group",
	"lt":                   "has_number_lt",
	"lte":                  "has_number_lte",
	"not_empty":            "has_text",
	"number":               "has_number",
	"phone":                "has_phone",
	"regex":                "has_pattern",
	"starts":               "has_beginning",
	"state":                "has_state",
	"timeout":              "has_wait_timed_out",
	"ward":                 "has_ward",
}

// migrates the given legacy action to a new action
func migrateAction(baseLanguage utils.Language, a Action, localization flows.Localization) (flows.Action, error) {
	switch a.Type {
	case "add_label":
		labels := make([]*flows.LabelReference, len(a.Labels))
		for i, label := range a.Labels {
			labels[i] = label.Migrate()
		}

		return &actions.AddInputLabelsAction{
			Labels:     labels,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil

	case "email":
		var msg string
		err := json.Unmarshal(a.Msg, &msg)
		if err != nil {
			return nil, err
		}

		migratedSubject, _ := expressions.MigrateTemplate(a.Subject, expressions.ExtraAsFunction)
		migratedBody, _ := expressions.MigrateTemplate(msg, expressions.ExtraAsFunction)
		migratedEmails := make([]string, len(a.Emails))
		for e, email := range a.Emails {
			migratedEmails[e], _ = expressions.MigrateTemplate(email, expressions.ExtraAsFunction)
		}

		return &actions.SendEmailAction{
			Subject:    migratedSubject,
			Body:       migratedBody,
			Addresses:  migratedEmails,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil

	case "lang":
		return &actions.SetContactLanguageAction{
			Language:   string(a.Language),
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil
	case "channel":
		return &actions.SetContactChannelAction{
			Channel:    flows.NewChannelReference(a.Channel, a.Name),
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil
	case "flow":
		return &actions.StartFlowAction{
			BaseAction: actions.NewBaseAction(a.UUID),
			Flow:       a.Flow.Migrate(),
		}, nil
	case "trigger-flow":
		contacts := make([]*flows.ContactReference, len(a.Contacts))
		for i, contact := range a.Contacts {
			contacts[i] = contact.Migrate()
		}
		groups := make([]*flows.GroupReference, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = group.Migrate()
		}
		var createContact bool
		variables := make([]string, 0, len(a.Variables))
		for _, variable := range a.Variables {
			if variable.ID == "@new_contact" {
				createContact = true
			} else {
				migratedVar, _ := expressions.MigrateTemplate(variable.ID, expressions.ExtraAsFunction)
				variables = append(variables, migratedVar)
			}
		}

		return &actions.StartSessionAction{
			BaseAction:    actions.NewBaseAction(a.UUID),
			Flow:          a.Flow.Migrate(),
			URNs:          []urns.URN{},
			Contacts:      contacts,
			Groups:        groups,
			LegacyVars:    variables,
			CreateContact: createContact,
		}, nil
	case "reply", "send":
		msg := make(Translations)
		media := make(Translations)
		var quickReplies map[utils.Language][]string

		err := json.Unmarshal(a.Msg, &msg)
		if err != nil {
			return nil, err
		}

		if a.Media != nil {
			err := json.Unmarshal(a.Media, &media)
			if err != nil {
				return nil, err
			}
		}
		if a.QuickReplies != nil {
			legacyQuickReplies := make([]Translations, 0)

			err := json.Unmarshal(a.QuickReplies, &legacyQuickReplies)
			if err != nil {
				return nil, err
			}

			quickReplies = TransformTranslations(legacyQuickReplies)
		}

		migratedText := addTranslationMap(baseLanguage, localization, msg, utils.UUID(a.UUID), "text")
		migratedMedia := addTranslationMap(baseLanguage, localization, media, utils.UUID(a.UUID), "attachments")
		migratedQuickReplies := addTranslationMultiMap(baseLanguage, localization, quickReplies, utils.UUID(a.UUID), "quick_replies")

		attachments := []string{}
		if migratedMedia != "" {
			attachments = append(attachments, migratedMedia)
		}

		if a.Type == "reply" {
			return &actions.SendMsgAction{
				BaseAction:   actions.NewBaseAction(a.UUID),
				Text:         migratedText,
				Attachments:  attachments,
				QuickReplies: migratedQuickReplies,
				AllURNs:      a.SendAll,
			}, nil
		}

		contacts := make([]*flows.ContactReference, len(a.Contacts))
		for i, contact := range a.Contacts {
			contacts[i] = contact.Migrate()
		}
		groups := make([]*flows.GroupReference, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = group.Migrate()
		}
		variables := make([]string, 0, len(a.Variables))
		for _, variable := range a.Variables {
			migratedVar, _ := expressions.MigrateTemplate(variable.ID, expressions.ExtraAsFunction)
			variables = append(variables, migratedVar)
		}

		return &actions.SendBroadcastAction{
			BaseAction:  actions.NewBaseAction(a.UUID),
			Text:        migratedText,
			Attachments: attachments,
			URNs:        []urns.URN{},
			Contacts:    contacts,
			Groups:      groups,
			LegacyVars:  variables,
		}, nil

	case "add_group":
		groups := make([]*flows.GroupReference, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = group.Migrate()
		}

		return &actions.AddContactGroupsAction{
			Groups:     groups,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil
	case "del_group":
		groups := make([]*flows.GroupReference, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = group.Migrate()
		}

		return &actions.RemoveContactGroupsAction{
			Groups:     groups,
			AllGroups:  len(groups) == 0,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil
	case "save":
		migratedValue, _ := expressions.MigrateTemplate(a.Value, expressions.ExtraAsFunction)

		// flows now have different action for name changing
		if a.Field == "name" || a.Field == "first_name" {
			// we can emulate setting only the first name with an expression
			if a.Field == "first_name" {
				migratedValue = strings.TrimSpace(migratedValue)
				migratedValue = fmt.Sprintf("%s @(word_slice(contact.name, 1, -1))", migratedValue)
			}

			return &actions.SetContactNameAction{
				Name:       migratedValue,
				BaseAction: actions.NewBaseAction(a.UUID),
			}, nil
		}

		// and another new action for adding a URN
		if urns.IsValidScheme(a.Field) {
			return &actions.AddContactURNAction{
				Scheme:     a.Field,
				Path:       migratedValue,
				BaseAction: actions.NewBaseAction(a.UUID),
			}, nil
		} else if a.Field == "tel_e164" {
			return &actions.AddContactURNAction{
				Scheme:     "tel",
				Path:       migratedValue,
				BaseAction: actions.NewBaseAction(a.UUID),
			}, nil
		}

		return &actions.SetContactFieldAction{
			Field:      flows.NewFieldReference(a.Field, a.Label),
			Value:      migratedValue,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil
	case "api":
		migratedURL, _ := expressions.MigrateTemplate(a.Webhook, expressions.ExtraAsFunction)

		headers := make(map[string]string, len(a.WebhookHeaders))
		body := ""

		if strings.ToUpper(a.Action) == "POST" {
			headers["Content-Type"] = "application/json"
			body = legacyWebhookBody
		}

		for _, header := range a.WebhookHeaders {
			headers[header.Name] = header.Value
		}

		return &actions.CallWebhookAction{
			BaseAction: actions.NewBaseAction(a.UUID),
			Method:     a.Action,
			URL:        migratedURL,
			Body:       body,
			Headers:    headers,
		}, nil
	default:
		return nil, fmt.Errorf("unable to migrate legacy action type: %s", a.Type)
	}
}

// migrates the given legacy rule to a router case
func migrateRule(baseLanguage utils.Language, exitMap map[string]flows.Exit, r Rule, localization flows.Localization) (routers.Case, error) {
	category := r.Category.Base(baseLanguage)

	newType, _ := testTypeMappings[r.Test.Type]
	var omitOperand bool
	var arguments []string
	var err error

	caseUUID := utils.UUID(utils.NewUUID())

	switch r.Test.Type {

	// tests that take no arguments
	case "date", "has_email", "not_empty", "number", "phone", "state":
		arguments = []string{}

	// tests against a single numeric value
	case "eq", "gt", "gte", "lt", "lte":
		test := numericTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		migratedTest, err := expressions.MigrateTemplate(string(test.Test), expressions.ExtraAsFunction)
		if err != nil {
			return routers.Case{}, err
		}
		arguments = []string{migratedTest}

	case "between":
		test := betweenTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		migratedMin, err := expressions.MigrateTemplate(test.Min, expressions.ExtraAsFunction)
		if err != nil {
			return routers.Case{}, err
		}
		migratedMax, err := expressions.MigrateTemplate(test.Max, expressions.ExtraAsFunction)
		if err != nil {
			return routers.Case{}, err
		}
		arguments = []string{migratedMin, migratedMax}

	// tests against a single localized string
	case "contains", "contains_any", "contains_phrase", "contains_only_phrase", "regex", "starts":
		test := localizedStringTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		arguments = []string{test.Test.Base(baseLanguage)}

		addTranslationMap(baseLanguage, localization, test.Test, caseUUID, "arguments")

	// tests against a single date value
	case "date_equal", "date_after", "date_before":
		test := stringTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		arguments = []string{test.Test}

	// tests against a single group value
	case "in_group":
		test := groupTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		arguments = []string{string(test.Test.UUID)}

	case "subflow":
		newType = "is_text_eq"
		test := subflowTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		arguments = []string{test.ExitType}

	case "webhook_status":
		newType = "has_webhook_status"
		test := webhookTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		if test.Status == "success" {
			arguments = []string{"success"}
		} else {
			arguments = []string{"response_error"}
		}

	case "timeout":
		omitOperand = true
		arguments = []string{"@run"}

	case "district":
		test := stringTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		migratedState, err := expressions.MigrateTemplate(test.Test, expressions.ExtraAsFunction)
		if err != nil {
			return routers.Case{}, err
		}
		arguments = []string{migratedState}

	case "ward":
		test := wardTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		migratedDistrict, err := expressions.MigrateTemplate(test.District, expressions.ExtraAsFunction)
		if err != nil {
			return routers.Case{}, err
		}
		migratedState, err := expressions.MigrateTemplate(test.State, expressions.ExtraAsFunction)
		if err != nil {
			return routers.Case{}, err
		}
		arguments = []string{migratedDistrict, migratedState}

	default:
		return routers.Case{}, fmt.Errorf("migration of '%s' tests no supported", r.Test.Type)
	}

	// TODO
	// airtime_status
	// ward / district / state
	// interrupted_status

	return routers.Case{
		UUID:        caseUUID,
		Type:        newType,
		Arguments:   arguments,
		OmitOperand: omitOperand,
		ExitUUID:    exitMap[category].UUID(),
	}, err
}

// temporary struct for migrating categories to cases and exits
type categoryName struct {
	uuid         flows.ExitUUID
	destination  flows.NodeUUID
	translations Translations
	order        int
}

func parseRules(baseLanguage utils.Language, r RuleSet, localization flows.Localization) ([]flows.Exit, []routers.Case, flows.ExitUUID, error) {

	// find our discrete categories and Other category (which uses the true rule)
	categoryMap := make(map[string]categoryName)
	var otherCategory *categoryName
	var otherCategoryBaseName string

	order := 0
	for _, rule := range r.Rules {
		categoryBaseName := rule.Category.Base(baseLanguage)

		if rule.Test.Type == "true" {
			otherCategoryBaseName = categoryBaseName
			otherCategory = &categoryName{
				uuid:         rule.UUID,
				destination:  rule.Destination,
				translations: rule.Category,
				order:        -1,
			}
		} else {
			_, ok := categoryMap[categoryBaseName]
			if !ok {
				categoryMap[categoryBaseName] = categoryName{
					uuid:         rule.UUID,
					destination:  rule.Destination,
					translations: rule.Category,
					order:        order,
				}
				order++
			}
		}
	}

	// create exits for each category
	exits := make([]flows.Exit, len(categoryMap))
	exitMap := make(map[string]flows.Exit)
	for k, category := range categoryMap {
		addTranslationMap(baseLanguage, localization, category.translations, utils.UUID(category.uuid), "name")

		exits[category.order] = definition.NewExit(category.uuid, category.destination, k)
		exitMap[k] = exits[category.order]
	}

	var defaultExitUUID flows.ExitUUID

	if otherCategory != nil {
		addTranslationMap(baseLanguage, localization, otherCategory.translations, utils.UUID(otherCategory.uuid), "name")
		defaultExit := definition.NewExit(otherCategory.uuid, otherCategory.destination, otherCategoryBaseName)
		exits = append(exits, defaultExit)

		defaultExitUUID = defaultExit.UUID()
	}

	// create any cases to map to our new exits
	var cases []routers.Case
	for i := range r.Rules {
		// skip Other rules
		if r.Rules[i].Test.Type == "true" {
			continue
		}

		c, err := migrateRule(baseLanguage, exitMap, r.Rules[i], localization)
		if err != nil {
			return nil, nil, "", err
		}

		cases = append(cases, c)

		if r.Rules[i].Test.Type == "webhook_status" {
			// webhook failures don't have a case, instead they become the default
			defaultExitUUID = exitMap[r.Rules[i].Category.Base(baseLanguage)].UUID()
		}
	}

	// for webhook rulesets we need to map 2 rules (success/failure) to 3 cases and exits (success/response_error/connection_error)
	if r.Type == "webhook" {
		connectionErrorCategory := "Connection Error"
		connectionErrorExitUUID := flows.ExitUUID(utils.NewUUID())
		connectionErrorExit := definition.NewExit(connectionErrorExitUUID, exits[1].DestinationNodeUUID(), connectionErrorCategory)

		exits = append(exits, connectionErrorExit)
		cases = append(cases, routers.Case{
			UUID:        utils.UUID(utils.NewUUID()),
			Type:        "has_webhook_status",
			Arguments:   []string{"connection_error"},
			OmitOperand: false,
			ExitUUID:    connectionErrorExitUUID,
		})
	}

	return exits, cases, defaultExitUUID, nil
}

type fieldConfig struct {
	FieldDelimiter string `json:"field_delimiter"`
	FieldIndex     int    `json:"field_index"`
}

// migrates the given legacy rulset to a node with a router
func migrateRuleSet(lang utils.Language, r RuleSet, localization flows.Localization) (flows.Node, error) {
	var newActions []flows.Action
	var router flows.Router
	var wait flows.Wait

	exits, cases, defaultExit, err := parseRules(lang, r, localization)
	if err != nil {
		return nil, err
	}

	resultName := r.Label

	switch r.Type {
	case "subflow":
		config := make(map[string]map[string]string)
		err := json.Unmarshal(r.Config, &config)
		if err != nil {
			return nil, err
		}

		flowUUID := flows.FlowUUID(config["flow"]["uuid"])
		flowName := config["flow"]["name"]

		newActions = []flows.Action{
			&actions.StartFlowAction{
				BaseAction: actions.NewBaseAction(flows.ActionUUID(utils.NewUUID())),
				Flow:       flows.NewFlowReference(flowUUID, flowName),
			},
		}

		// subflow rulesets operate on the child flow status
		router = routers.NewSwitchRouter(defaultExit, "@child.status", cases, resultName)

	case "webhook":
		var config WebhookConfig
		err := json.Unmarshal(r.Config, &config)
		if err != nil {
			return nil, err
		}

		migratedURL, _ := expressions.MigrateTemplate(config.Webhook, expressions.ExtraAsFunction)
		headers := make(map[string]string, len(config.Headers))
		body := ""

		if strings.ToUpper(config.Action) == "POST" {
			headers["Content-Type"] = "application/json"
			body = legacyWebhookBody
		}

		for _, header := range config.Headers {
			headers[header.Name], _ = expressions.MigrateTemplate(header.Value, expressions.ExtraAsFunction)
		}

		newActions = []flows.Action{
			&actions.CallWebhookAction{
				BaseAction: actions.NewBaseAction(flows.ActionUUID(utils.NewUUID())),
				URL:        migratedURL,
				Method:     config.Action,
				Headers:    headers,
				Body:       body,
			},
		}

		// webhook rulesets operate on the webhook call
		router = routers.NewSwitchRouter(defaultExit, "@run.webhook", cases, resultName)

	case "form_field":
		var config fieldConfig
		json.Unmarshal(r.Config, &config)

		operand, _ := expressions.MigrateTemplate(r.Operand, expressions.ExtraAsFunction)
		operand = fmt.Sprintf("@(field(%s, %d, \"%s\"))", operand[1:], config.FieldIndex, config.FieldDelimiter)
		router = routers.NewSwitchRouter(defaultExit, operand, cases, resultName)

	case "group":
		// in legacy flows these rulesets have their operand as @step.value but it's not used
		router = routers.NewSwitchRouter(defaultExit, "@contact", cases, resultName)

	case "wait_message":
		// look for timeout test on the legacy ruleset
		var timeout *int
		for _, rule := range r.Rules {
			if rule.Test.Type == "timeout" {
				test := timeoutTest{}
				if err := json.Unmarshal(rule.Test.Data, &test); err != nil {
					return nil, err
				}
				t := 60 * test.Minutes
				timeout = &t
				break
			}
		}

		wait = waits.NewMsgWait(timeout)

		fallthrough
	case "flow_field":
		fallthrough
	case "contact_field":
		fallthrough
	case "expression":
		operand, _ := expressions.MigrateTemplate(r.Operand, expressions.ExtraAsFunction)
		if operand == "" {
			operand = "@run.input"
		}

		router = routers.NewSwitchRouter(defaultExit, operand, cases, resultName)
	case "random":
		router = routers.NewRandomRouter(resultName)
	default:
		return nil, fmt.Errorf("unrecognized ruleset type: %s", r.Type)
	}

	return definition.NewNode(r.UUID, newActions, router, exits, wait), nil
}

// migrates the given legacy actionset to a node with a set of migrated actions and a single exit
func migateActionSet(lang utils.Language, a ActionSet, localization flows.Localization) (flows.Node, error) {
	actions := make([]flows.Action, len(a.Actions))

	// migrate each action
	for i := range a.Actions {
		action, err := migrateAction(lang, a.Actions[i], localization)
		if err != nil {
			return nil, fmt.Errorf("error migrating action[type=%s]: %s", a.Actions[i].Type, err)
		}
		actions[i] = action
	}

	return definition.NewNode(a.UUID, actions, nil, []flows.Exit{definition.NewExit(a.ExitUUID, a.Destination, "")}, nil), nil
}

// ReadLegacyFlows reads in legacy formatted flows
func ReadLegacyFlows(data []json.RawMessage) ([]*Flow, error) {
	var err error
	flows := make([]*Flow, len(data))
	for f := range data {
		flows[f], err = ReadLegacyFlow(data[f])
		if err != nil {
			return nil, err
		}
	}

	return flows, nil
}

// ReadLegacyFlow reads a single legacy formatted flow
func ReadLegacyFlow(data json.RawMessage) (*Flow, error) {
	flow := &Flow{}
	if err := utils.UnmarshalAndValidate(data, flow, ""); err != nil {
		return nil, err
	}
	return flow, nil
}

// Migrate migrates this legacy flow to the new format
func (f *Flow) Migrate(includeUI bool) (flows.Flow, error) {
	localization := definition.NewLocalization()
	nodes := make([]flows.Node, len(f.ActionSets)+len(f.RuleSets))

	for i := range f.ActionSets {
		node, err := migateActionSet(f.BaseLanguage, f.ActionSets[i], localization)
		if err != nil {
			return nil, fmt.Errorf("error migrating action_set[uuid=%s]: %s", f.ActionSets[i].UUID, err)
		}
		nodes[i] = node
	}

	for i := range f.RuleSets {
		node, err := migrateRuleSet(f.BaseLanguage, f.RuleSets[i], localization)
		if err != nil {
			return nil, fmt.Errorf("error migrating rule_set[uuid=%s]: %s", f.RuleSets[i].UUID, err)
		}
		nodes[len(f.ActionSets)+i] = node
	}

	// make sure our entry node is first
	for i := range nodes {
		if nodes[i].UUID() == f.Entry {
			firstNode := nodes[0]
			nodes[0] = nodes[i]
			nodes[i] = firstNode
		}
	}

	ui := make(map[string]interface{})

	if includeUI {
		// convert our UI metadata
		nodesUI := make(map[flows.NodeUUID]interface{})

		for i := range f.ActionSets {
			actionset := f.ActionSets[i]
			nmd := make(map[string]interface{})
			nmd["position"] = map[string]int{
				"left": actionset.X,
				"top":  actionset.Y,
			}
			nodesUI[actionset.UUID] = nmd
		}

		for i := range f.RuleSets {
			ruleset := f.RuleSets[i]
			nmd := make(map[string]interface{})
			nmd["position"] = map[string]int{
				"left": ruleset.X,
				"top":  ruleset.Y,
			}
			nodesUI[ruleset.UUID] = nmd
		}

		stickies := make(map[utils.UUID]Sticky, len(f.Metadata.Notes))
		for _, note := range f.Metadata.Notes {
			stickies[utils.NewUUID()] = note.Migrate()
		}

		ui["nodes"] = nodesUI
		ui["stickies"] = stickies
	}

	return definition.NewFlow(
		f.Metadata.UUID,
		f.Metadata.Name,
		f.Metadata.Revision,
		f.BaseLanguage,
		f.Metadata.Expires,
		localization,
		nodes,
		ui,
	)
}
