package definition

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"

	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/routers"
	"github.com/nyaruka/goflow/flows/waits"
	"github.com/nyaruka/goflow/utils"
	"github.com/satori/go.uuid"
)

// represents a decimal value which may be provided as a string or floating point value
type decimalString string

func (s *decimalString) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		// data is a quoted string
		*s = decimalString(data[1 : len(data)-1])
	} else {
		// data is JSON float
		*s = decimalString(data)
	}
	return nil
}

// LegacyFlow imports an old-world flow so it can be exported anew
type LegacyFlow struct {
	flow
	envelope legacyFlowEnvelope
}

type legacyFlowEnvelope struct {
	BaseLanguage utils.Language         `json:"base_language"`
	Metadata     legacyMetadataEnvelope `json:"metadata"`
	RuleSets     []legacyRuleSet        `json:"rule_sets" validate:"dive"`
	ActionSets   []legacyActionSet      `json:"action_sets" validate:"dive"`
	Entry        flows.NodeUUID         `json:"entry" validate:"required,uuid4"`
}

type legacyMetadataEnvelope struct {
	UUID    flows.FlowUUID `json:"uuid" validate:"required,uuid4"`
	Name    string         `json:"name"`
	Expires int            `json:"expires"`
}

type legacyRule struct {
	UUID            flows.ExitUUID            `json:"uuid" validate:"required,uuid4"`
	Destination     flows.NodeUUID            `json:"destination" validate:"omitempty,uuid4"`
	DestinationType string                    `json:"destination_type" validate:"eq=A|eq=R"`
	Test            utils.TypedEnvelope       `json:"test"`
	Category        map[utils.Language]string `json:"category"`
}

type legacyRuleSet struct {
	Y       int             `json:"y"`
	X       int             `json:"x"`
	UUID    flows.NodeUUID  `json:"uuid" validate:"required,uuid4"`
	Type    string          `json:"ruleset_type"`
	Label   string          `json:"label"`
	Operand string          `json:"operand"`
	Rules   []legacyRule    `json:"rules"`
	Config  json.RawMessage `json:"config"`
}

type legacyActionSet struct {
	Y           int            `json:"y"`
	X           int            `json:"x"`
	Destination flows.NodeUUID `json:"destination" validate:"omitempty,uuid4"`
	ExitUUID    flows.ExitUUID `json:"exit_uuid" validate:"required,uuid4"`
	UUID        flows.NodeUUID `json:"uuid" validate:"required,uuid4"`
	Actions     []legacyAction `json:"actions"`
}

type legacyLabelReference struct {
	UUID flows.LabelUUID
	Name string
}

func (l *legacyLabelReference) Migrate() *flows.LabelReference {
	if len(l.UUID) > 0 {
		return flows.NewLabelReference(l.UUID, l.Name)
	}
	return flows.NewVariableLabelReference(l.Name)
}

func (l *legacyLabelReference) UnmarshalJSON(data []byte) error {
	// label reference may be a string
	if data[0] == '"' {
		var nameExpression string
		if err := json.Unmarshal(data, &nameExpression); err != nil {
			return err
		}

		// if it starts with @ then it's an expression
		if strings.HasPrefix(nameExpression, "@") {
			nameExpression, _ = excellent.MigrateTemplate(nameExpression)
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

type legacyContactReference struct {
	UUID flows.ContactUUID `json:"uuid"`
}

func (c *legacyContactReference) Migrate() *flows.ContactReference {
	return flows.NewContactReference(c.UUID, "")
}

type legacyGroupReference struct {
	UUID flows.GroupUUID
	Name string
}

func (g *legacyGroupReference) Migrate() *flows.GroupReference {
	if len(g.UUID) > 0 {
		return flows.NewGroupReference(g.UUID, g.Name)
	}
	return flows.NewVariableGroupReference(g.Name)
}

func (g *legacyGroupReference) UnmarshalJSON(data []byte) error {
	// group reference may be a string
	if data[0] == '"' {
		var nameExpression string
		if err := json.Unmarshal(data, &nameExpression); err != nil {
			return err
		}

		// if it starts with @ then it's an expression
		if strings.HasPrefix(nameExpression, "@") {
			nameExpression, _ = excellent.MigrateTemplate(nameExpression)
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

type legacyVariable struct {
	ID string `json:"id"`
}

type legacyFlowReference struct {
	UUID flows.FlowUUID `json:"uuid"`
	Name string         `json:"name"`
}

func (f *legacyFlowReference) Migrate() *flows.FlowReference {
	return flows.NewFlowReference(f.UUID, f.Name)
}

type legacyWebhookConfig struct {
	Webhook string                `json:"webhook"`
	Action  string                `json:"webhook_action"`
	Headers []legacyWebhookHeader `json:"webhook_headers"`
}

type legacyWebhookHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type legacyAction struct {
	Type string           `json:"type"`
	UUID flows.ActionUUID `json:"uuid"`
	Name string           `json:"name"`

	// message and email
	Msg          json.RawMessage `json:"msg"`
	Media        json.RawMessage `json:"media"`
	QuickReplies json.RawMessage `json:"quick_replies"`
	SendAll      bool            `json:"send_all"`

	// variable contact actions
	Contacts  []legacyContactReference `json:"contacts"`
	Groups    []legacyGroupReference   `json:"groups"`
	Variables []legacyVariable         `json:"variables"`

	// save actions
	Field string `json:"field"`
	Value string `json:"value"`
	Label string `json:"label"`

	// set language
	Language utils.Language `json:"lang"`

	// webhook
	Action         string                `json:"action"`
	Webhook        string                `json:"webhook"`
	WebhookHeaders []legacyWebhookHeader `json:"webhook_headers"`

	// add lable action
	Labels []legacyLabelReference `json:"labels"`

	// Start/Trigger flow
	Flow legacyFlowReference `json:"flow"`

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
	Test map[utils.Language]string `json:"test"`
}

type stringTest struct {
	Test string `json:"test"`
}

type numericTest struct {
	Test decimalString `json:"test"`
}

type betweenTest struct {
	Min string `json:"min"`
	Max string `json:"max"`
}

type timeoutTest struct {
	Minutes int `json:"minutes"`
}

type groupTest struct {
	Test legacyGroupReference `json:"test"`
}

type wardTest struct {
	State    string `json:"state"`
	District string `json:"district"`
}

type localizations map[utils.Language]flows.Action

func addTranslationMap(baseLanguage utils.Language, translations *flowTranslations, mapped map[utils.Language]string, uuid flows.UUID, key string) string {
	var inBaseLanguage string
	for language, item := range mapped {
		expression, _ := excellent.MigrateTemplate(item)
		if language != baseLanguage {
			addTranslation(translations, language, uuid, key, []string{expression})
		} else {
			inBaseLanguage = expression
		}
	}

	return inBaseLanguage
}

func addTranslationMultiMap(baseLanguage utils.Language, translations *flowTranslations, mapped map[utils.Language][]string, uuid flows.UUID, key string) []string {
	var inBaseLanguage []string
	for language, items := range mapped {
		expressions := make([]string, len(items))
		for i := range items {
			expression, _ := excellent.MigrateTemplate(items[i])
			expressions[i] = expression
		}
		if language != baseLanguage {
			addTranslation(translations, language, uuid, key, expressions)
		} else {
			inBaseLanguage = expressions
		}
	}
	return inBaseLanguage
}

func addTranslation(translations *flowTranslations, lang utils.Language, itemUUID flows.UUID, propKey string, translation []string) {
	// ensure we have a translation set for this language
	langTranslations, found := (*translations)[lang]
	if !found {
		langTranslations = &languageTranslations{}
		(*translations)[lang] = langTranslations
	}

	// ensure we have a translation set for this item
	itemTrans, found := (*langTranslations)[itemUUID]
	if !found {
		itemTrans = itemTranslations{}
		(*langTranslations)[itemUUID] = itemTrans
	}

	itemTrans[propKey] = translation
}

// Transforms a list of single item translations into a map of multi-item translations, e.g.
//
// [{"eng": "yes", "fra": "oui"}, {"eng": "no", "fra": "non"}] becomes {"eng": ["yes", "no"], "fra": ["oui", "non"]}
//
func transformTranslations(items []map[utils.Language]string) map[utils.Language][]string {
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
	"email":                "has_email",
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
func migrateAction(baseLanguage utils.Language, a legacyAction, translations *flowTranslations) (flows.Action, error) {
	switch a.Type {
	case "add_label":
		labels := make([]*flows.LabelReference, len(a.Labels))
		for i, label := range a.Labels {
			labels[i] = label.Migrate()
		}

		return &actions.AddLabelAction{
			Labels:     labels,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil

	case "email":
		var msg string
		err := json.Unmarshal(a.Msg, &msg)
		if err != nil {
			return nil, err
		}

		migratedSubject, _ := excellent.MigrateTemplate(a.Subject)
		migratedBody, _ := excellent.MigrateTemplate(msg)
		migratedEmails := make([]string, len(a.Emails))
		for e, email := range a.Emails {
			migratedEmails[e], _ = excellent.MigrateTemplate(email)
		}

		return &actions.EmailAction{
			Subject:    migratedSubject,
			Body:       migratedBody,
			Addresses:  migratedEmails,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil

	case "lang":
		return &actions.UpdateContactAction{
			FieldName:  "language",
			Value:      string(a.Language),
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil
	case "channel":
		return &actions.PreferredChannelAction{
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
				migratedVar, _ := excellent.MigrateTemplate(variable.ID)
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
		msg := make(map[utils.Language]string)
		media := make(map[utils.Language]string)
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
			legacyQuickReplies := make([]map[utils.Language]string, 0)

			err := json.Unmarshal(a.QuickReplies, &legacyQuickReplies)
			if err != nil {
				return nil, err
			}

			quickReplies = transformTranslations(legacyQuickReplies)
		}

		migratedText := addTranslationMap(baseLanguage, translations, msg, flows.UUID(a.UUID), "text")
		migratedMedia := addTranslationMap(baseLanguage, translations, media, flows.UUID(a.UUID), "attachments")
		migratedQuickReplies := addTranslationMultiMap(baseLanguage, translations, quickReplies, flows.UUID(a.UUID), "quick_replies")

		attachments := []string{}
		if migratedMedia != "" {
			attachments = append(attachments, migratedMedia)
		}

		if a.Type == "reply" {
			return &actions.ReplyAction{
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
			migratedVar, _ := excellent.MigrateTemplate(variable.ID)
			variables = append(variables, migratedVar)
		}

		return &actions.SendMsgAction{
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

		return &actions.AddToGroupAction{
			Groups:     groups,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil
	case "del_group":
		groups := make([]*flows.GroupReference, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = group.Migrate()
		}

		return &actions.RemoveFromGroupAction{
			Groups:     groups,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil
	case "save":
		migratedValue, _ := excellent.MigrateTemplate(a.Value)

		// flows now have different action for name changing
		if a.Field == "name" || a.Field == "first_name" {
			// we can emulate setting only the first name with an expression
			if a.Field == "first_name" {
				migratedValue = fmt.Sprintf("%s @(word_slice(contact.name, 2, -1))", migratedValue)
			}

			return &actions.UpdateContactAction{
				FieldName:  "name",
				Value:      migratedValue,
				BaseAction: actions.NewBaseAction(a.UUID),
			}, nil
		}

		// and another new action for adding a URN
		if urns.IsValidScheme(a.Field) {
			return &actions.AddURNAction{
				Scheme:     a.Field,
				Path:       migratedValue,
				BaseAction: actions.NewBaseAction(a.UUID),
			}, nil
		}

		return &actions.SaveContactField{
			Field:      flows.NewFieldReference(flows.FieldKey(a.Field), a.Label),
			Value:      migratedValue,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil
	case "api":
		migratedURL, _ := excellent.MigrateTemplate(a.Webhook)

		headers := make(map[string]string, len(a.WebhookHeaders))
		for _, header := range a.WebhookHeaders {
			headers[header.Name] = header.Value
		}

		return &actions.WebhookAction{
			Method:     a.Action,
			URL:        migratedURL,
			Headers:    headers,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil
	default:
		return nil, fmt.Errorf("couldn't create action for %s", a.Type)
	}
}

// migrates the given legacy rule to a router case
func migrateRule(baseLanguage utils.Language, exitMap map[string]flows.Exit, r legacyRule, translations *flowTranslations) (routers.Case, error) {
	category := r.Category[baseLanguage]

	newType, _ := testTypeMappings[r.Test.Type]
	var omitOperand bool
	var arguments []string
	var err error

	caseUUID := flows.UUID(uuid.NewV4().String())

	switch r.Test.Type {

	// tests that take no arguments
	case "date", "email", "not_empty", "number", "phone", "state":
		arguments = []string{}

	// tests against a single numeric value
	case "eq", "gt", "gte", "lt", "lte":
		test := numericTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		migratedTest, err := excellent.MigrateTemplate(string(test.Test))
		if err != nil {
			return routers.Case{}, err
		}
		arguments = []string{migratedTest}

	case "between":
		test := betweenTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		migratedMin, err := excellent.MigrateTemplate(test.Min)
		if err != nil {
			return routers.Case{}, err
		}
		migratedMax, err := excellent.MigrateTemplate(test.Max)
		if err != nil {
			return routers.Case{}, err
		}
		arguments = []string{migratedMin, migratedMax}

	// tests against a single localized string
	case "contains", "contains_any", "contains_phrase", "contains_only_phrase", "regex", "starts":
		test := localizedStringTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		arguments = []string{test.Test[baseLanguage]}

		addTranslationMap(baseLanguage, translations, test.Test, caseUUID, "arguments")

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
		newType = "is_string_eq"
		test := subflowTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		arguments = []string{test.ExitType}

	case "webhook_status":
		newType = "is_string_eq"
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
		arguments = []string{test.Test}

	case "ward":
		test := wardTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		arguments = []string{test.District, test.State}

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

type categoryName struct {
	uuid         flows.ExitUUID
	destination  flows.NodeUUID
	translations map[utils.Language]string
	order        int
}

func parseRules(baseLanguage utils.Language, r legacyRuleSet, translations *flowTranslations) ([]flows.Exit, []routers.Case, flows.ExitUUID, error) {

	// find our discrete categories
	categoryMap := make(map[string]categoryName)
	order := 0
	for i := range r.Rules {
		category := r.Rules[i].Category[baseLanguage]
		_, ok := categoryMap[category]
		if !ok {
			categoryMap[category] = categoryName{
				uuid:         r.Rules[i].UUID,
				destination:  r.Rules[i].Destination,
				translations: r.Rules[i].Category,
				order:        order,
			}
			order++
		}
	}

	// create exits for each category
	exits := make([]flows.Exit, len(categoryMap))
	exitMap := make(map[string]flows.Exit)
	for k, category := range categoryMap {
		addTranslationMap(baseLanguage, translations, category.translations, flows.UUID(category.uuid), "name")

		exits[category.order] = NewExit(category.uuid, category.destination, k)
		exitMap[k] = exits[category.order]
	}

	var defaultExit flows.ExitUUID

	// create any cases to map to our new exits
	var cases []routers.Case
	for i := range r.Rules {
		if r.Rules[i].Test.Type == "true" {
			// take the first true rule as our default exit
			if defaultExit == "" {
				defaultExit = exitMap[r.Rules[i].Category[baseLanguage]].UUID()
			}
			continue
		}

		c, err := migrateRule(baseLanguage, exitMap, r.Rules[i], translations)
		if err != nil {
			return nil, nil, "", err
		}

		cases = append(cases, c)

		if r.Rules[i].Test.Type == "webhook_status" {
			// webhook failures don't have a case, instead they are the default
			defaultExit = exitMap[r.Rules[i].Category[baseLanguage]].UUID()
		}
	}

	// for webhook rulesets we need to map 2 rules (success/failure) to 3 cases and exits (success/response_error/connection_error)
	if r.Type == "webhook" {
		connectionErrorCategory := "Connection Error"
		connectionErrorExitUUID := flows.ExitUUID(uuid.NewV4().String())
		connectionErrorExit := NewExit(connectionErrorExitUUID, exits[1].(*exit).destination, connectionErrorCategory)

		exits = append(exits, connectionErrorExit)
		cases = append(cases, routers.Case{
			UUID:        flows.UUID(uuid.NewV4().String()),
			Type:        "is_string_eq",
			Arguments:   []string{"connection_error"},
			OmitOperand: false,
			ExitUUID:    connectionErrorExitUUID,
		})
	}

	return exits, cases, defaultExit, nil
}

type fieldConfig struct {
	FieldDelimiter string `json:"field_delimiter"`
	FieldIndex     int    `json:"field_index"`
}

// migrates the given legacy rulset to a node with a router
func migrateRuleSet(lang utils.Language, r legacyRuleSet, translations *flowTranslations) (*node, error) {
	node := &node{}
	node.uuid = r.UUID

	exits, cases, defaultExit, err := parseRules(lang, r, translations)
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

		node.actions = []flows.Action{
			&actions.StartFlowAction{
				BaseAction: actions.NewBaseAction(flows.ActionUUID(uuid.NewV4().String())),
				Flow:       flows.NewFlowReference(flowUUID, flowName),
			},
		}

		// subflow rulesets operate on the child flow status
		node.router = routers.NewSwitchRouter(defaultExit, "@child.status", cases, resultName)

	case "webhook":
		var config legacyWebhookConfig
		err := json.Unmarshal(r.Config, &config)
		if err != nil {
			return nil, err
		}

		migratedHeaders := make(map[string]string, len(config.Headers))
		for _, header := range config.Headers {
			migratedHeaders[header.Name] = header.Value
		}

		node.actions = []flows.Action{
			&actions.WebhookAction{
				BaseAction: actions.NewBaseAction(flows.ActionUUID(uuid.NewV4().String())),
				URL:        config.Webhook,
				Method:     config.Action,
				Headers:    migratedHeaders,
			},
		}

		// subflow rulesets operate on the child flow status
		node.router = routers.NewSwitchRouter(defaultExit, "@run.webhook.status", cases, resultName)

	case "form_field":
		var config fieldConfig
		json.Unmarshal(r.Config, &config)

		operand, _ := excellent.MigrateTemplate(r.Operand)
		operand = fmt.Sprintf("@(field(%s, %d, \"%s\"))", operand[1:], config.FieldIndex, config.FieldDelimiter)
		node.router = routers.NewSwitchRouter(defaultExit, operand, cases, resultName)

	case "group":
		// in legacy flows these rulesets have their operand as @step.value but it's not used
		node.router = routers.NewSwitchRouter(defaultExit, "@contact", cases, resultName)

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

		node.wait = waits.NewMsgWait(timeout)

		fallthrough
	case "flow_field":
		fallthrough
	case "contact_field":
		fallthrough
	case "expression":
		operand, _ := excellent.MigrateTemplate(r.Operand)
		if operand == "" {
			operand = "@run.input"
		}

		node.router = routers.NewSwitchRouter(defaultExit, operand, cases, resultName)
	case "random":
		node.router = routers.NewRandomRouter(resultName)
	default:
		fmt.Printf("Unable to migrate unrecognized ruleset type: '%s'\n", r.Type)
	}

	node.exits = exits

	return node, nil
}

// migrates the given legacy actionset to a node with a set of migrated actions and a single exit
func migateActionSet(lang utils.Language, a legacyActionSet, translations *flowTranslations) (*node, error) {
	node := &node{
		uuid:    a.UUID,
		actions: make([]flows.Action, len(a.Actions)),
		exits: []flows.Exit{
			NewExit(a.ExitUUID, a.Destination, ""),
		},
	}

	// migrate each action
	for i := range a.Actions {
		action, err := migrateAction(lang, a.Actions[i], translations)
		if err != nil {
			return nil, err
		}
		node.actions[i] = action
	}

	return node, nil
}

// ReadLegacyFlows reads in legacy formatted flows
func ReadLegacyFlows(data []json.RawMessage) ([]*LegacyFlow, error) {
	var err error
	flows := make([]*LegacyFlow, len(data))
	for f := range data {
		flows[f], err = ReadLegacyFlow(data[f])
		if err != nil {
			return nil, err
		}
	}

	return flows, nil
}

func ReadLegacyFlow(data json.RawMessage) (*LegacyFlow, error) {
	var envelope legacyFlowEnvelope
	var err error

	if err := utils.UnmarshalAndValidate(data, &envelope, ""); err != nil {
		return nil, err
	}

	f := &LegacyFlow{}
	f.uuid = envelope.Metadata.UUID
	f.name = envelope.Metadata.Name
	f.language = envelope.BaseLanguage
	f.expireAfterMinutes = envelope.Metadata.Expires

	translations := &flowTranslations{}

	f.nodes = make([]flows.Node, len(envelope.ActionSets)+len(envelope.RuleSets))
	for i := range envelope.ActionSets {
		node, err := migateActionSet(f.language, envelope.ActionSets[i], translations)
		if err != nil {
			return nil, err
		}
		f.nodes[i] = node
	}

	for i := range envelope.RuleSets {
		node, err := migrateRuleSet(f.language, envelope.RuleSets[i], translations)
		if err != nil {
			return nil, err
		}
		f.nodes[len(envelope.ActionSets)+i] = node
	}

	// make sure our entry node is first
	for i := range f.nodes {
		if f.nodes[i].UUID() == envelope.Entry {
			firstNode := f.nodes[0]
			f.nodes[0] = f.nodes[i]
			f.nodes[i] = firstNode
		}
	}

	f.translations = translations
	f.envelope = envelope

	return f, err
}

// MarshalJSON sends turns our legacy flow into bytes
func (f *LegacyFlow) MarshalJSON() ([]byte, error) {

	var fe = flowEnvelope{}
	fe.UUID = f.uuid
	fe.Name = f.name
	fe.Language = f.language
	fe.ExpireAfterMinutes = f.expireAfterMinutes

	if f.translations != nil {
		fe.Localization = *f.translations.(*flowTranslations)
	}

	fe.Nodes = make([]*node, len(f.nodes))
	for i := range f.nodes {
		fe.Nodes[i] = f.nodes[i].(*node)
	}

	// add in our ui metadata
	fe.Metadata = make(map[string]interface{})
	fe.Metadata["nodes"] = make(map[flows.NodeUUID]interface{})
	nodes := fe.Metadata["nodes"].(map[flows.NodeUUID]interface{})

	for i := range f.envelope.ActionSets {
		actionset := f.envelope.ActionSets[i]
		nmd := make(map[string]interface{})
		nmd["position"] = map[string]int{
			"x": actionset.X,
			"y": actionset.Y,
		}
		nodes[actionset.UUID] = nmd
	}

	for i := range f.envelope.RuleSets {
		ruleset := f.envelope.RuleSets[i]
		nmd := make(map[string]interface{})
		nmd["position"] = map[string]int{
			"x": ruleset.X,
			"y": ruleset.Y,
		}
		nodes[ruleset.UUID] = nmd
	}

	return json.Marshal(&fe)
}
