package legacy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/extensions/transferto"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/routers"
	"github.com/nyaruka/goflow/flows/waits"
	"github.com/nyaruka/goflow/legacy/expressions"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

// Flow is a flow in the legacy format
type Flow struct {
	BaseLanguage utils.Language `json:"base_language"`
	FlowType     string         `json:"flow_type"`
	Metadata     Metadata       `json:"metadata"`
	RuleSets     []RuleSet      `json:"rule_sets" validate:"dive"`
	ActionSets   []ActionSet    `json:"action_sets" validate:"dive"`
	Entry        flows.NodeUUID `json:"entry" validate:"omitempty,uuid4"`
}

// Metadata is the metadata section of a legacy flow
type Metadata struct {
	UUID     assets.FlowUUID `json:"uuid" validate:"required,uuid4"`
	Name     string          `json:"name"`
	Revision int             `json:"revision"`
	Expires  int             `json:"expires"`
	Notes    []Note          `json:"notes,omitempty"`
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
	UUID assets.LabelUUID
	Name string
}

func (l *LabelReference) Migrate() *assets.LabelReference {
	if len(l.UUID) > 0 {
		return assets.NewLabelReference(l.UUID, l.Name)
	}
	return assets.NewVariableLabelReference(l.Name)
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
			nameExpression, _ = expressions.MigrateTemplate(nameExpression, expressions.ExtraAsFunction, false)
		}

		l.Name = nameExpression
		return nil
	}

	// or a JSON object with UUID/Name properties
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	l.UUID = assets.LabelUUID(raw["uuid"].(string))
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
	UUID assets.GroupUUID
	Name string
}

func (g *GroupReference) Migrate() *assets.GroupReference {
	if len(g.UUID) > 0 {
		return assets.NewGroupReference(g.UUID, g.Name)
	}
	return assets.NewVariableGroupReference(g.Name)
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
			nameExpression, _ = expressions.MigrateTemplate(nameExpression, expressions.ExtraAsFunction, false)
		}

		g.Name = nameExpression
		return nil
	}

	// or a JSON object with UUID/Name properties
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	g.UUID = assets.GroupUUID(raw["uuid"].(string))
	g.Name = raw["name"].(string)
	return nil
}

type VariableReference struct {
	ID string `json:"id"`
}

type FlowReference struct {
	UUID assets.FlowUUID `json:"uuid"`
	Name string          `json:"name"`
}

func (f *FlowReference) Migrate() *assets.FlowReference {
	return assets.NewFlowReference(f.UUID, f.Name)
}

// RulesetConfig holds the config dictionary for a legacy ruleset
type RulesetConfig struct {
	Flow           *assets.FlowReference `json:"flow"`
	FieldDelimiter string                `json:"field_delimiter"`
	FieldIndex     int                   `json:"field_index"`
	Webhook        string                `json:"webhook"`
	WebhookAction  string                `json:"webhook_action"`
	WebhookHeaders []WebhookHeader       `json:"webhook_headers"`
	Resthook       string                `json:"resthook"`
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
	Channel assets.ChannelUUID `json:"channel"`

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

type airtimeTest struct {
	ExitStatus string `json:"exit_status"`
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

var flowTypeMapping = map[string]flows.FlowType{
	"F": flows.FlowTypeMessaging,
	"M": flows.FlowTypeMessaging,
	"V": flows.FlowTypeVoice,
	"S": flows.FlowTypeMessagingOffline,
}

func addTranslationMap(baseLanguage utils.Language, localization flows.Localization, mapped Translations, uuid utils.UUID, property string) string {
	var inBaseLanguage string
	for language, item := range mapped {
		expression, _ := expressions.MigrateTemplate(item, expressions.ExtraAsFunction, false)
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
			expression, _ := expressions.MigrateTemplate(items[i], expressions.ExtraAsFunction, false)
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
	"webhook_status":       "has_webhook_status",
	"airtime_status":       "has_airtime_status",
}

// migrates the given legacy action to a new action
func migrateAction(baseLanguage utils.Language, a Action, localization flows.Localization) (flows.Action, error) {
	switch a.Type {
	case "add_label":
		labels := make([]*assets.LabelReference, len(a.Labels))
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

		migratedSubject, _ := expressions.MigrateTemplate(a.Subject, expressions.ExtraAsFunction, false)
		migratedBody, _ := expressions.MigrateTemplate(msg, expressions.ExtraAsFunction, false)
		migratedEmails := make([]string, len(a.Emails))
		for e, email := range a.Emails {
			migratedEmails[e], _ = expressions.MigrateTemplate(email, expressions.ExtraAsFunction, false)
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
			Channel:    assets.NewChannelReference(a.Channel, a.Name),
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
		groups := make([]*assets.GroupReference, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = group.Migrate()
		}
		var createContact bool
		variables := make([]string, 0, len(a.Variables))
		for _, variable := range a.Variables {
			if variable.ID == "@new_contact" {
				createContact = true
			} else {
				migratedVar, _ := expressions.MigrateTemplate(variable.ID, expressions.ExtraAsFunction, false)
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
		groups := make([]*assets.GroupReference, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = group.Migrate()
		}
		variables := make([]string, 0, len(a.Variables))
		for _, variable := range a.Variables {
			migratedVar, _ := expressions.MigrateTemplate(variable.ID, expressions.ExtraAsFunction, false)
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
		groups := make([]*assets.GroupReference, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = group.Migrate()
		}

		return &actions.AddContactGroupsAction{
			Groups:     groups,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil
	case "del_group":
		groups := make([]*assets.GroupReference, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = group.Migrate()
		}

		return &actions.RemoveContactGroupsAction{
			Groups:     groups,
			AllGroups:  len(groups) == 0,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil
	case "save":
		migratedValue, _ := expressions.MigrateTemplate(a.Value, expressions.ExtraAsFunction, false)

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
			Field:      assets.NewFieldReference(a.Field, a.Label),
			Value:      migratedValue,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil
	case "api":
		migratedURL, _ := expressions.MigrateTemplate(a.Webhook, expressions.ExtraAsFunction, false)

		headers := make(map[string]string, len(a.WebhookHeaders))
		body := ""
		method := strings.ToUpper(a.Action)
		if method == "" {
			method = "POST"
		}

		if method == "POST" {
			headers["Content-Type"] = "application/json"
			body = flows.DefaultWebhookPayload
		}

		for _, header := range a.WebhookHeaders {
			headers[header.Name] = header.Value
		}

		return &actions.CallWebhookAction{
			BaseAction: actions.NewBaseAction(a.UUID),
			Method:     method,
			URL:        migratedURL,
			Body:       body,
			Headers:    headers,
		}, nil
	default:
		return nil, fmt.Errorf("unable to migrate legacy action type: %s", a.Type)
	}
}

// migrates the given legacy rulset to a node with a router
func migrateRuleSet(lang utils.Language, r RuleSet, localization flows.Localization, collapseExits bool) (flows.Node, UINodeType, error) {
	var newActions []flows.Action
	var router flows.Router
	var wait flows.Wait
	var uiType UINodeType

	cases, exits, defaultExit, err := migrateRules(lang, r, localization, collapseExits)
	if err != nil {
		return nil, "", err
	}

	resultName := r.Label

	// load the config for this ruleset
	var config RulesetConfig
	if r.Config != nil {
		err := json.Unmarshal(r.Config, &config)
		if err != nil {
			return nil, "", err
		}
	}

	switch r.Type {
	case "subflow":
		newActions = []flows.Action{
			&actions.StartFlowAction{
				BaseAction: actions.NewBaseAction(flows.ActionUUID(utils.NewUUID())),
				Flow:       config.Flow,
			},
		}

		// subflow rulesets operate on the child flow status
		router = routers.NewSwitchRouter(defaultExit, "@child.status", cases, resultName)
		uiType = UINodeTypeSplitBySubflow

	case "webhook":
		migratedURL, _ := expressions.MigrateTemplate(config.Webhook, expressions.ExtraAsFunction, false)
		headers := make(map[string]string, len(config.WebhookHeaders))
		body := ""
		method := strings.ToUpper(config.WebhookAction)
		if method == "" {
			method = "POST"
		}

		if method == "POST" {
			headers["Content-Type"] = "application/json"
			body = flows.DefaultWebhookPayload
		}

		for _, header := range config.WebhookHeaders {
			headers[header.Name], _ = expressions.MigrateTemplate(header.Value, expressions.ExtraAsFunction, false)
		}

		newActions = []flows.Action{
			&actions.CallWebhookAction{
				BaseAction: actions.NewBaseAction(flows.ActionUUID(utils.NewUUID())),
				URL:        migratedURL,
				Method:     method,
				Headers:    headers,
				Body:       body,
			},
		}

		// webhook rulesets operate on the webhook call
		router = routers.NewSwitchRouter(defaultExit, "@run.webhook", cases, resultName)
		uiType = UINodeTypeSplitByWebhook

	case "resthook":
		newActions = []flows.Action{
			&actions.CallResthookAction{
				BaseAction: actions.NewBaseAction(flows.ActionUUID(utils.NewUUID())),
				Resthook:   config.Resthook,
			},
		}

		// webhook rulesets operate on the webhook call
		router = routers.NewSwitchRouter(defaultExit, "@run.webhook", cases, resultName)

	case "form_field":
		operand, _ := expressions.MigrateTemplate(r.Operand, expressions.ExtraAsFunction, false)
		operand = fmt.Sprintf("@(field(%s, %d, \"%s\"))", operand[1:], config.FieldIndex, config.FieldDelimiter)
		router = routers.NewSwitchRouter(defaultExit, operand, cases, resultName)
		uiType = UINodeTypeSplitByRunResultDelimited

	case "group":
		// in legacy flows these rulesets have their operand as @step.value but it's not used
		router = routers.NewSwitchRouter(defaultExit, "@contact", cases, resultName)
		uiType = UINodeTypeSplitByGroups

	case "wait_message":
		// look for timeout test on the legacy ruleset
		var timeout *int
		for _, rule := range r.Rules {
			if rule.Test.Type == "timeout" {
				test := timeoutTest{}
				if err := json.Unmarshal(rule.Test.Data, &test); err != nil {
					return nil, "", err
				}
				t := 60 * test.Minutes
				timeout = &t
				break
			}
		}

		wait = waits.NewMsgWait(timeout)
		uiType = UINodeTypeWaitForResponse

		fallthrough
	case "flow_field", "contact_field", "expression":
		// unlike other templates, operands for expression rulesets need to be wrapped in such a way that if
		// they error, they evaluate to the original expression
		var defaultToSelf bool
		switch r.Type {
		case "flow_field":
			uiType = UINodeTypeSplitByRunResult
		case "contact_field":
			uiType = UINodeTypeSplitByContactField
		case "expression":
			defaultToSelf = true
			uiType = UINodeTypeSplitByExpression
		}

		operand, _ := expressions.MigrateTemplate(r.Operand, expressions.ExtraAsFunction, defaultToSelf)
		if operand == "" {
			operand = "@run.input"
		}

		router = routers.NewSwitchRouter(defaultExit, operand, cases, resultName)
	case "random":
		router = routers.NewRandomRouter(resultName)
		uiType = UINodeTypeSplitByRandom

	case "airtime":
		countryConfigs := map[string]struct {
			CurrencyCode string          `json:"currency_code"`
			Amount       decimal.Decimal `json:"amount"`
		}{}
		if err := json.Unmarshal(r.Config, &countryConfigs); err != nil {
			return nil, "", err
		}
		currencyAmounts := make(map[string]decimal.Decimal, len(countryConfigs))
		for _, countryCfg := range countryConfigs {
			// check if we already have a configuration for this currency
			existingAmount, alreadyDefined := currencyAmounts[countryCfg.CurrencyCode]
			if alreadyDefined && existingAmount != countryCfg.Amount {
				return nil, "", fmt.Errorf("unable to migrate airtime ruleset with different amounts in same currency")
			}

			currencyAmounts[countryCfg.CurrencyCode] = countryCfg.Amount
		}

		newActions = []flows.Action{
			&transferto.TransferAirtimeAction{
				BaseAction: actions.NewBaseAction(flows.ActionUUID(utils.NewUUID())),
				Amounts:    currencyAmounts,
			},
		}

		router = routers.NewSwitchRouter(defaultExit, "@run", cases, resultName)

	default:
		return nil, "", fmt.Errorf("unrecognized ruleset type: %s", r.Type)
	}

	return definition.NewNode(r.UUID, newActions, router, exits, wait), uiType, nil
}

// migrates a set of legacy rules to sets of cases and exits
func migrateRules(baseLanguage utils.Language, r RuleSet, localization flows.Localization, collapseExits bool) ([]routers.Case, []flows.Exit, flows.ExitUUID, error) {
	cases := make([]routers.Case, 0, len(r.Rules))
	exits := make([]flows.Exit, 0, len(r.Rules))
	var defaultExitUUID flows.ExitUUID

	ruleUUIDsToExits := make(map[flows.ExitUUID]flows.Exit, len(r.Rules))
	categoriesToExits := make(map[string]flows.Exit, len(r.Rules))

	// creating exits from the rules
	for _, rule := range r.Rules {
		baseName := rule.Category.Base(baseLanguage)
		var exit flows.Exit

		// if we're collapsing exits, then we can use the exit previously created for this category
		if collapseExits && rule.Test.Type != "true" {
			exit = categoriesToExits[baseName]
		}
		if exit == nil {
			exit = definition.NewExit(rule.UUID, rule.Destination, baseName)
			exits = append(exits, exit)
		}

		ruleUUIDsToExits[rule.UUID] = exit
		categoriesToExits[baseName] = exit

		addTranslationMap(baseLanguage, localization, rule.Category, utils.UUID(exit.UUID()), "name")
	}

	// and then a case for each rule
	for _, rule := range r.Rules {
		// implicit Other rules don't become cases
		if rule.Test.Type == "true" {
			defaultExitUUID = rule.UUID
			continue
		} else if rule.Test.Type == "webhook_status" {
			// default case for a webhook ruleset is the last migrated rule (failure)
			defaultExitUUID = rule.UUID
		}

		exit := ruleUUIDsToExits[rule.UUID]

		kase, err := migrateRule(baseLanguage, rule, exit, localization)
		if err != nil {
			return nil, nil, "", err
		}

		cases = append(cases, kase)
	}

	// for webhook rulesets we need to add an additional case/error pair for connection errors
	if r.Type == "webhook" || r.Type == "resthook" {
		connectionErrorCategory := "Connection Error"
		connectionErrorExit := definition.NewExit(flows.ExitUUID(utils.NewUUID()), exits[1].DestinationNodeUUID(), connectionErrorCategory)

		cases = append(cases, routers.Case{
			UUID:        utils.UUID(utils.NewUUID()),
			Type:        "has_webhook_status",
			Arguments:   []string{"connection_error"},
			OmitOperand: false,
			ExitUUID:    connectionErrorExit.UUID(),
		})
		exits = append(exits, connectionErrorExit)
	}

	return cases, exits, defaultExitUUID, nil
}

// migrates the given legacy rule to a router case
func migrateRule(baseLanguage utils.Language, r Rule, exit flows.Exit, localization flows.Localization) (routers.Case, error) {
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
		migratedTest, err := expressions.MigrateTemplate(string(test.Test), expressions.ExtraAsFunction, false)
		if err != nil {
			return routers.Case{}, err
		}
		arguments = []string{migratedTest}

	case "between":
		test := betweenTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		migratedMin, err := expressions.MigrateTemplate(test.Min, expressions.ExtraAsFunction, false)
		if err != nil {
			return routers.Case{}, err
		}
		migratedMax, err := expressions.MigrateTemplate(test.Max, expressions.ExtraAsFunction, false)
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
		test := webhookTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		if test.Status == "success" {
			arguments = []string{"success"}
		} else {
			arguments = []string{"response_error"}
		}

	case "airtime_status":
		test := airtimeTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		arguments = []string{test.ExitStatus}

	case "timeout":
		omitOperand = true
		arguments = []string{"@run"}

	case "district":
		test := stringTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		migratedState, err := expressions.MigrateTemplate(test.Test, expressions.ExtraAsFunction, false)
		if err != nil {
			return routers.Case{}, err
		}
		arguments = []string{migratedState}

	case "ward":
		test := wardTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		migratedDistrict, err := expressions.MigrateTemplate(test.District, expressions.ExtraAsFunction, false)
		if err != nil {
			return routers.Case{}, err
		}
		migratedState, err := expressions.MigrateTemplate(test.State, expressions.ExtraAsFunction, false)
		if err != nil {
			return routers.Case{}, err
		}
		arguments = []string{migratedDistrict, migratedState}

	default:
		return routers.Case{}, fmt.Errorf("migration of '%s' tests no supported", r.Test.Type)
	}

	return routers.Case{
		UUID:        caseUUID,
		Type:        newType,
		Arguments:   arguments,
		OmitOperand: omitOperand,
		ExitUUID:    exit.UUID(),
	}, err
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

// ReadLegacyFlow reads a single legacy formatted flow
func ReadLegacyFlow(data json.RawMessage) (*Flow, error) {
	flow := &Flow{}
	if err := utils.UnmarshalAndValidate(data, flow); err != nil {
		return nil, err
	}
	return flow, nil
}

// Migrate migrates this legacy flow to the new format
func (f *Flow) Migrate(collapseExits bool, includeUI bool) (flows.Flow, error) {
	localization := definition.NewLocalization()
	numNodes := len(f.ActionSets) + len(f.RuleSets)
	nodes := make([]flows.Node, numNodes)
	nodeUITypes := make(map[flows.NodeUUID]UINodeType, numNodes)

	for i := range f.ActionSets {
		node, err := migateActionSet(f.BaseLanguage, f.ActionSets[i], localization)
		if err != nil {
			return nil, fmt.Errorf("error migrating action_set[uuid=%s]: %s", f.ActionSets[i].UUID, err)
		}
		nodes[i] = node
	}

	for i := range f.RuleSets {
		node, uiType, err := migrateRuleSet(f.BaseLanguage, f.RuleSets[i], localization, collapseExits)
		if err != nil {
			return nil, fmt.Errorf("error migrating rule_set[uuid=%s]: %s", f.RuleSets[i].UUID, err)
		}
		nodes[len(f.ActionSets)+i] = node
		nodeUITypes[node.UUID()] = uiType
	}

	// make sure our entry node is first
	if f.Entry != "" {
		for i := range nodes {
			if nodes[i].UUID() == f.Entry {
				firstNode := nodes[0]
				nodes[0] = nodes[i]
				nodes[i] = firstNode
			}
		}
	}

	var ui UI

	if includeUI {
		ui = NewUI()

		for _, actionSet := range f.ActionSets {
			ui.AddNode(actionSet.UUID, actionSet.X, actionSet.Y, "")
		}
		for _, ruleSet := range f.RuleSets {
			ui.AddNode(ruleSet.UUID, ruleSet.X, ruleSet.Y, nodeUITypes[ruleSet.UUID])
		}
		for _, note := range f.Metadata.Notes {
			ui.AddSticky(note.Migrate())
		}
	}

	return definition.NewFlow(
		f.Metadata.UUID,
		f.Metadata.Name,
		f.BaseLanguage,
		flowTypeMapping[f.FlowType],
		f.Metadata.Revision,
		f.Metadata.Expires,
		localization,
		nodes,
		ui,
	)
}
