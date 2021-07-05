package legacy

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition/legacy/expressions"
	"github.com/nyaruka/goflow/utils"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

//------------------------------------------------------------------------------------------
// Legacy flow objects
//------------------------------------------------------------------------------------------

// Flow is a flow in the legacy format
type Flow struct {
	BaseLanguage envs.Language `json:"base_language"`
	FlowType     string        `json:"flow_type"`
	RuleSets     []RuleSet     `json:"rule_sets" validate:"dive"`
	ActionSets   []ActionSet   `json:"action_sets" validate:"dive"`
	Entry        uuids.UUID    `json:"entry" validate:"omitempty,uuid4"`
	Metadata     *Metadata     `json:"metadata"`

	// some flows have these set here instead of in metadata
	UUID uuids.UUID `json:"uuid"`
	Name string     `json:"name"`
}

// Metadata is the metadata section of a legacy flow
type Metadata struct {
	UUID     uuids.UUID `json:"uuid"`
	Name     string     `json:"name"`
	Revision int        `json:"revision"`
	Expires  int        `json:"expires"`
	Notes    []Note     `json:"notes,omitempty"`
}

type Rule struct {
	UUID            uuids.UUID    `json:"uuid" validate:"required,uuid4"`
	Destination     uuids.UUID    `json:"destination" validate:"omitempty,uuid4"`
	DestinationType string        `json:"destination_type" validate:"eq=A|eq=R"`
	Test            TypedEnvelope `json:"test"`
	Category        Translations  `json:"category"`
}

type RuleSet struct {
	Y           int             `json:"y"`
	X           int             `json:"x"`
	UUID        uuids.UUID      `json:"uuid" validate:"required,uuid4"`
	Type        string          `json:"ruleset_type"`
	Label       string          `json:"label"`
	Operand     string          `json:"operand"`
	Rules       []Rule          `json:"rules"`
	Config      json.RawMessage `json:"config"`
	FinishedKey string          `json:"finished_key"`
}

type ActionSet struct {
	Y           int        `json:"y"`
	X           int        `json:"x"`
	Destination uuids.UUID `json:"destination" validate:"omitempty,uuid4"`
	ExitUUID    uuids.UUID `json:"exit_uuid" validate:"required,uuid4"`
	UUID        uuids.UUID `json:"uuid" validate:"required,uuid4"`
	Actions     []Action   `json:"actions"`
}

type LabelReference struct {
	UUID uuids.UUID
	Name string
}

// UnmarshalJSON unmarshals a legacy label reference from the given JSON
func (l *LabelReference) UnmarshalJSON(data []byte) error {
	// label reference may be a string
	if data[0] == '"' {
		var nameExpression string
		if err := jsonx.Unmarshal(data, &nameExpression); err != nil {
			return err
		}

		// if it starts with @ then it's an expression
		if strings.HasPrefix(nameExpression, "@") {
			nameExpression, _ = expressions.MigrateTemplate(nameExpression, nil)
		}

		l.Name = nameExpression
		return nil
	}

	// or a JSON object with UUID/Name properties
	var raw map[string]string
	if err := jsonx.Unmarshal(data, &raw); err != nil {
		return err
	}

	l.UUID = uuids.UUID(raw["uuid"])
	l.Name = raw["name"]
	return nil
}

type ContactReference struct {
	UUID uuids.UUID `json:"uuid"`
	Name string     `json:"name"`
}

type GroupReference struct {
	UUID uuids.UUID
	Name string
}

func (g *GroupReference) Migrate() *assets.GroupReference {
	if len(g.UUID) > 0 {
		return assets.NewGroupReference(assets.GroupUUID(g.UUID), g.Name)
	}
	return assets.NewVariableGroupReference(g.Name)
}

// UnmarshalJSON unmarshals a legacy group reference from the given JSON
func (g *GroupReference) UnmarshalJSON(data []byte) error {
	// group reference may be a string
	if data[0] == '"' {
		var nameExpression string
		if err := jsonx.Unmarshal(data, &nameExpression); err != nil {
			return err
		}

		// if it starts with @ then it's an expression
		if strings.HasPrefix(nameExpression, "@") {
			nameExpression, _ = expressions.MigrateTemplate(nameExpression, nil)
		}

		g.Name = nameExpression
		return nil
	}

	// or a JSON object with UUID/Name properties
	var raw map[string]string
	if err := jsonx.Unmarshal(data, &raw); err != nil {
		return err
	}

	g.UUID = uuids.UUID(raw["uuid"])
	g.Name = raw["name"]
	return nil
}

type VariableReference struct {
	ID string `json:"id"`
}

type FlowReference struct {
	UUID uuids.UUID `json:"uuid"`
	Name string     `json:"name"`
}

// RulesetConfig holds the config dictionary for a legacy ruleset
type RulesetConfig struct {
	Flow           *FlowReference  `json:"flow"`
	FieldDelimiter string          `json:"field_delimiter"`
	FieldIndex     int             `json:"field_index"`
	Webhook        string          `json:"webhook"`
	WebhookAction  string          `json:"webhook_action"`
	WebhookHeaders []WebhookHeader `json:"webhook_headers"`
	Resthook       string          `json:"resthook"`
}

type WebhookHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Action struct {
	Type string     `json:"type"`
	UUID uuids.UUID `json:"uuid"`
	Name string     `json:"name"`

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
	Language envs.Language `json:"lang"`

	// add label action
	Labels []LabelReference `json:"labels"`

	// start/trigger flow
	Flow FlowReference `json:"flow"`

	// channel
	Channel uuids.UUID `json:"channel"`

	// email
	Emails  []string `json:"emails"`
	Subject string   `json:"subject"`

	// IVR
	Recording json.RawMessage `json:"recording"`
	URL       string          `json:"url"`
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
	Test StringOrNumber `json:"test"`
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

var relativeDateTest = regexp.MustCompile(`@\(date\.today\s+(\+|\-)\s+(\-?\d+)\)`)

//------------------------------------------------------------------------------------------
// Migrated flow objects
//------------------------------------------------------------------------------------------

var flowTypeMapping = map[string]string{
	"":  "messaging", // some campaign event flows are missing this
	"F": "messaging",
	"M": "messaging",
	"V": "voice",
	"S": "messaging_offline",
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
	"ward":                 "has_ward",
}

// migrates the given legacy action to a new action
func migrateAction(baseLanguage envs.Language, a Action, localization migratedLocalization, baseMediaURL string) (migratedAction, error) {
	switch a.Type {
	case "add_label":
		labels := make([]*assets.LabelReference, len(a.Labels))
		for i, label := range a.Labels {
			if len(label.UUID) > 0 {
				labels[i] = assets.NewLabelReference(assets.LabelUUID(label.UUID), label.Name)
			} else {
				labels[i] = assets.NewVariableLabelReference(label.Name)
			}
		}

		return newAddInputLabelsAction(a.UUID, labels), nil

	case "email":
		var msg string
		err := jsonx.Unmarshal(a.Msg, &msg)
		if err != nil {
			return nil, err
		}

		migratedSubject, _ := expressions.MigrateTemplate(a.Subject, nil)
		migratedBody, _ := expressions.MigrateTemplate(msg, nil)
		migratedEmails := make([]string, len(a.Emails))
		for i, email := range a.Emails {
			migratedEmails[i], _ = expressions.MigrateTemplate(email, nil)
		}

		return newSendEmailAction(a.UUID, migratedEmails, migratedSubject, migratedBody), nil

	case "lang":
		return newSetContactLanguageAction(a.UUID, string(a.Language)), nil
	case "channel":
		return newSetContactChannelAction(a.UUID, assets.NewChannelReference(assets.ChannelUUID(a.Channel), a.Name)), nil
	case "flow":
		flowRef := assets.NewFlowReference(assets.FlowUUID(a.Flow.UUID), a.Flow.Name)

		return newEnterFlowAction(a.UUID, flowRef, true), nil
	case "trigger-flow":
		flowRef := assets.NewFlowReference(assets.FlowUUID(a.Flow.UUID), a.Flow.Name)

		contacts := make([]*flows.ContactReference, len(a.Contacts))
		for i, contact := range a.Contacts {
			contacts[i] = flows.NewContactReference(flows.ContactUUID(contact.UUID), contact.Name)
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
				migratedVar, _ := expressions.MigrateTemplate(variable.ID, nil)
				variables = append(variables, migratedVar)
			}
		}

		return newStartSessionAction(a.UUID, flowRef, []urns.URN{}, contacts, groups, variables, createContact), nil
	case "reply", "send":
		media := make(Translations)
		var quickReplies map[envs.Language][]string

		msg, err := ReadTranslations(a.Msg)
		if err != nil {
			return nil, err
		}

		if a.Media != nil {
			err := jsonx.Unmarshal(a.Media, &media)
			if err != nil {
				return nil, err
			}
		}
		if a.QuickReplies != nil {
			legacyQuickReplies := make([]Translations, 0)

			err := jsonx.Unmarshal(a.QuickReplies, &legacyQuickReplies)
			if err != nil {
				return nil, err
			}

			quickReplies = TransformTranslations(legacyQuickReplies)
		}

		for lang, attachment := range media {
			parts := strings.SplitN(attachment, ":", 2)
			var mediaType, mediaURL string
			if len(parts) == 2 {
				mediaType = parts[0]
				mediaURL = parts[1]
			} else {
				// no media type defaults to image
				mediaType = "image"
				mediaURL = parts[0]
			}

			// attachment is a real upload and not just an expression, need to make it absolute
			if !strings.Contains(mediaURL, "@") {
				media[lang] = fmt.Sprintf("%s:%s", mediaType, URLJoin(baseMediaURL, mediaURL))
			}
		}

		migratedText := localization.addTranslationMap(baseLanguage, msg, uuids.UUID(a.UUID), "text")
		migratedMedia := localization.addTranslationMap(baseLanguage, media, uuids.UUID(a.UUID), "attachments")
		migratedQuickReplies := localization.addTranslationMultiMap(baseLanguage, quickReplies, uuids.UUID(a.UUID), "quick_replies")

		attachments := []string{}
		if migratedMedia != "" {
			attachments = append(attachments, migratedMedia)
		}

		if a.Type == "reply" {
			return newSendMsgAction(a.UUID, migratedText, attachments, migratedQuickReplies, a.SendAll), nil
		}

		contacts := make([]*flows.ContactReference, len(a.Contacts))
		for i, contact := range a.Contacts {
			contacts[i] = flows.NewContactReference(flows.ContactUUID(contact.UUID), contact.Name)
		}
		groups := make([]*assets.GroupReference, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = group.Migrate()
		}
		variables := make([]string, 0, len(a.Variables))
		for _, variable := range a.Variables {
			migratedVar, _ := expressions.MigrateTemplate(variable.ID, nil)
			variables = append(variables, migratedVar)
		}

		return newSendBroadcastAction(a.UUID, migratedText, attachments, migratedQuickReplies, []urns.URN{}, contacts, groups, variables), nil

	case "add_group":
		groups := make([]*assets.GroupReference, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = group.Migrate()
		}

		return newAddContactGroupsAction(a.UUID, groups), nil
	case "del_group":
		groups := make([]*assets.GroupReference, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = group.Migrate()
		}

		allGroups := len(groups) == 0
		return newRemoveContactGroupsAction(a.UUID, groups, allGroups), nil
	case "save":
		migratedValue, _ := expressions.MigrateTemplate(a.Value, nil)

		// flows now have different action for name changing
		if a.Field == "name" || a.Field == "first_name" {
			// we can emulate setting only the first name with an expression
			if a.Field == "first_name" {
				migratedValue = strings.TrimSpace(migratedValue)
				migratedValue = fmt.Sprintf("%s @(word_slice(contact.name, 1, -1))", migratedValue)
			}

			return newSetContactNameAction(a.UUID, migratedValue), nil
		}

		// and another new action for adding a URN
		if urns.IsValidScheme(a.Field) {
			return newAddContactURNAction(a.UUID, a.Field, migratedValue), nil
		} else if a.Field == "tel_e164" {
			return newAddContactURNAction(a.UUID, "tel", migratedValue), nil
		}

		return newSetContactFieldAction(a.UUID, assets.NewFieldReference(a.Field, a.Label), migratedValue), nil
	case "say":
		msg, err := ReadTranslations(a.Msg)
		if err != nil {
			return nil, err
		}
		recording, err := ReadTranslations(a.Recording)
		if err != nil {
			return nil, err
		}

		// make audio URLs absolute
		for lang, audioURL := range recording {
			if audioURL != "" {
				recording[lang] = URLJoin(baseMediaURL, audioURL)
			}
		}

		migratedText := localization.addTranslationMap(baseLanguage, msg, uuids.UUID(a.UUID), "text")
		migratedAudioURL := localization.addTranslationMap(baseLanguage, recording, uuids.UUID(a.UUID), "audio_url")

		return newSayMsgAction(a.UUID, migratedText, migratedAudioURL), nil
	case "play":
		// note this URL is already assumed to be absolute
		migratedAudioURL, _ := expressions.MigrateTemplate(a.URL, nil)

		return newPlayAudioAction(a.UUID, migratedAudioURL), nil
	default:
		return nil, errors.Errorf("unable to migrate legacy action type: %s", a.Type)
	}
}

// migrates the given legacy rulset to a node with a router
func migrateRuleSet(lang envs.Language, r RuleSet, validDests map[uuids.UUID]bool, localization migratedLocalization) (migratedNode, UINodeType, NodeUIConfig, error) {
	var newActions []migratedAction
	var router migratedRouter
	var wait migratedWait
	var uiType UINodeType
	uiConfig := make(NodeUIConfig)

	cases, categories, defaultCategory, timeoutCategory, exits, err := migrateRules(lang, r, validDests, localization, uiConfig)
	if err != nil {
		return nil, "", nil, err
	}

	resultName := r.Label

	// load the config for this ruleset
	var config RulesetConfig
	if r.Config != nil {
		err := jsonx.Unmarshal(r.Config, &config)
		if err != nil {
			return nil, "", nil, err
		}
	}

	// sometimes old flows don't have this set
	if r.Type == "" {
		r.Type = "wait_message"
	}

	switch r.Type {
	case "subflow":
		flowRef := assets.NewFlowReference(assets.FlowUUID(config.Flow.UUID), config.Flow.Name)

		newActions = []migratedAction{
			newEnterFlowAction(uuids.New(), flowRef, false),
		}

		// subflow rulesets operate on the child flow status
		router = newSwitchRouter(nil, resultName, categories, "@child.status", cases, defaultCategory)
		uiType = UINodeTypeSplitBySubflow

	case "webhook":
		migratedURL, _ := expressions.MigrateTemplate(config.Webhook, &expressions.MigrateOptions{URLEncode: true})
		headers := make(map[string]string, len(config.WebhookHeaders))
		body := ""
		method := strings.ToUpper(config.WebhookAction)
		if method == "" {
			method = "POST"
		}

		if method == "POST" {
			headers["Content-Type"] = "application/json"
			body = legacyWebhookPayload
		}

		for _, header := range config.WebhookHeaders {
			// ignore empty headers sometimes left in flow definitions
			if header.Name != "" {
				headers[header.Name], _ = expressions.MigrateTemplate(header.Value, nil)
			}
		}

		newActions = []migratedAction{
			newCallWebhookAction(uuids.New(), method, migratedURL, headers, body, resultName),
		}

		// webhook rulesets operate on the webhook status, saved as category
		operand := fmt.Sprintf("@results.%s.category", utils.Snakify(resultName))
		router = newSwitchRouter(nil, "", categories, operand, cases, defaultCategory)
		uiType = UINodeTypeSplitByWebhook

	case "resthook":
		newActions = []migratedAction{
			newCallResthookAction(uuids.New(), config.Resthook, resultName),
		}

		// resthook rulesets operate on the webhook status, saved as category
		operand := fmt.Sprintf("@results.%s.category", utils.Snakify(resultName))
		router = newSwitchRouter(nil, "", categories, operand, cases, defaultCategory)
		uiType = UINodeTypeSplitByResthook

	case "form_field":
		operand, _ := expressions.MigrateTemplate(r.Operand, nil)
		operand = fmt.Sprintf("@(field(%s, %d, \"%s\"))", operand[1:], config.FieldIndex, config.FieldDelimiter)
		router = newSwitchRouter(nil, resultName, categories, operand, cases, defaultCategory)

		lastDot := strings.LastIndex(r.Operand, ".")
		if lastDot > -1 {
			fieldKey := r.Operand[lastDot+1:]

			uiConfig["operand"] = map[string]string{"id": fieldKey}
			uiConfig["delimiter"] = config.FieldDelimiter
			uiConfig["index"] = config.FieldIndex
		}

		uiType = UINodeTypeSplitByRunResultDelimited

	case "group":
		// in legacy flows these rulesets have their operand as @step.value but it's not used
		router = newSwitchRouter(nil, resultName, categories, "@contact.groups", cases, defaultCategory)
		uiType = UINodeTypeSplitByGroups

	case "wait_message", "wait_audio", "wait_video", "wait_photo", "wait_gps", "wait_recording", "wait_digit", "wait_digits":
		// look for timeout test on the legacy ruleset
		timeoutSeconds := 0
		for _, rule := range r.Rules {
			if rule.Test.Type == "timeout" {
				test := timeoutTest{}
				if err := jsonx.Unmarshal(rule.Test.Data, &test); err != nil {
					return nil, "", nil, err
				}
				timeoutSeconds = 60 * test.Minutes
				break
			}
		}

		var timeout migratedTimeout
		if timeoutSeconds > 0 && timeoutCategory != "" {
			timeout = newTimeout(timeoutSeconds, timeoutCategory)
		}

		hint, operand := migrateWaitingRuleset(r)
		wait = newMsgWait(timeout, hint)
		uiType = UINodeTypeWaitForResponse

		router = newSwitchRouter(wait, resultName, categories, operand, cases, defaultCategory)
	case "flow_field", "contact_field", "expression":
		// unlike other templates, operands for expression rulesets need to be wrapped in such a way that if
		// they error, they evaluate to the original expression
		var defaultToSelf bool
		switch r.Type {
		case "flow_field":
			uiType = UINodeTypeSplitByRunResult
			lastDot := strings.LastIndex(r.Operand, ".")
			if lastDot > -1 {
				fieldKey := r.Operand[lastDot+1:]

				uiConfig["operand"] = map[string]string{"id": fieldKey}
			}
		case "contact_field":
			uiType = UINodeTypeSplitByContactField

			lastDot := strings.LastIndex(r.Operand, ".")
			if lastDot > -1 {
				fieldKey := r.Operand[lastDot+1:]
				if fieldKey == "name" {
					uiConfig["operand"] = map[string]string{
						"type": "property",
						"id":   "name",
						"name": "Name",
					}
				} else if fieldKey == "groups" {
					uiType = UINodeTypeSplitByExpression

				} else if urns.IsValidScheme(fieldKey) {
					uiConfig["operand"] = map[string]string{
						"type": "scheme",
						"id":   fieldKey,
					}
				} else {
					uiConfig["operand"] = map[string]string{
						"type": "field",
						"id":   fieldKey,
					}
				}
			}

		case "expression":
			defaultToSelf = true
			uiType = UINodeTypeSplitByExpression
		}

		operand, _ := expressions.MigrateTemplate(r.Operand, &expressions.MigrateOptions{DefaultToSelf: defaultToSelf})
		if operand == "" {
			operand = "@input"
		}

		router = newSwitchRouter(wait, resultName, categories, operand, cases, defaultCategory)
	case "random":
		router = newRandomRouter(resultName, categories)
		uiType = UINodeTypeSplitByRandom

	case "airtime":
		countryConfigs := map[string]struct {
			CurrencyCode string          `json:"currency_code"`
			Amount       decimal.Decimal `json:"amount"`
		}{}
		if err := jsonx.Unmarshal(r.Config, &countryConfigs); err != nil {
			return nil, "", nil, err
		}
		currencyAmounts := make(map[string]decimal.Decimal, len(countryConfigs))
		for _, countryCfg := range countryConfigs {
			// check if we already have a configuration for this currency
			existingAmount, alreadyDefined := currencyAmounts[countryCfg.CurrencyCode]
			if alreadyDefined && existingAmount != countryCfg.Amount {
				return nil, "", nil, errors.Errorf("unable to migrate airtime ruleset with different amounts in same currency")
			}

			currencyAmounts[countryCfg.CurrencyCode] = countryCfg.Amount
		}

		newActions = []migratedAction{
			newTransferAirtimeAction(uuids.New(), currencyAmounts, resultName),
		}

		operand := fmt.Sprintf("@results.%s", utils.Snakify(resultName))
		router = newSwitchRouter(nil, "", categories, operand, cases, defaultCategory)
		uiType = UINodeTypeSplitByAirtime

	default:
		return nil, "", nil, errors.Errorf("unrecognized ruleset type: %s", r.Type)
	}

	return newNode(r.UUID, newActions, router, exits), uiType, uiConfig, nil
}

func migrateWaitingRuleset(r RuleSet) (migratedHint, string) {
	switch r.Type {
	case "wait_audio":
		return newAudioHint(), "@input"
	case "wait_video":
		return newVideoHint(), "@input"
	case "wait_photo":
		return newImageHint(), "@input"
	case "wait_gps":
		return newLocationHint(), "@input"
	case "wait_recording":
		return newAudioHint(), "@input"
	case "wait_digit":
		return newFixedDigitsHint(1), "@input.text"
	case "wait_digits":
		return newTerminatedDigitsHint(r.FinishedKey), "@input.text"
	}
	return nil, "@input"
}

type categoryAndExit struct {
	category migratedCategory
	exit     migratedExit
}

// migrates a set of legacy rules to sets of categories, cases and exits
func migrateRules(baseLanguage envs.Language, r RuleSet, validDests map[uuids.UUID]bool, localization migratedLocalization, uiConfig NodeUIConfig) ([]migratedCase, []migratedCategory, uuids.UUID, uuids.UUID, []migratedExit, error) {
	cases := make([]migratedCase, 0, len(r.Rules))
	categories := make([]migratedCategory, 0, len(r.Rules))
	exits := make([]migratedExit, 0, len(r.Rules))

	var defaultCategoryUUID, timeoutCategoryUUID uuids.UUID

	convertedByRuleUUID := make(map[uuids.UUID]*categoryAndExit, len(r.Rules))
	convertedByCategoryName := make(map[string]*categoryAndExit, len(r.Rules))

	// create categories and exits from the rules
	for _, rule := range r.Rules {
		baseName := rule.Category.Base(baseLanguage)
		var converted *categoryAndExit

		// check if we have previously created category/exits for this category name
		if rule.Test.Type != "true" {
			converted = convertedByCategoryName[baseName]
		}
		if converted == nil {
			// only set exit destination if it's valid
			var destinationUUID uuids.UUID
			if validDests[rule.Destination] {
				destinationUUID = rule.Destination
			}

			// rule UUIDs in legacy flows determine path data, so their UUIDs become the exit UUIDs
			exit := newExit(rule.UUID, destinationUUID)
			exits = append(exits, exit)

			category := newCategory(uuids.New(), baseName, exit.UUID())
			categories = append(categories, category)

			converted = &categoryAndExit{category, exit}
		}

		convertedByRuleUUID[rule.UUID] = converted
		convertedByCategoryName[baseName] = converted

		localization.addTranslationMap(baseLanguage, rule.Category, uuids.UUID(converted.category.UUID()), "name")
	}

	// and then a case for each rule
	for _, rule := range r.Rules {
		converted := convertedByRuleUUID[rule.UUID]

		if rule.Test.Type == "true" {
			// implicit Other rules don't become cases, but instead become the router default
			defaultCategoryUUID = uuids.UUID(converted.category.UUID())
			continue

		} else if rule.Test.Type == "timeout" {
			// timeout rules become category setting on the wait
			timeoutCategoryUUID = uuids.UUID(converted.category.UUID())
			continue

		} else if rule.Test.Type == "webhook_status" || rule.Test.Type == "airtime_status" {
			// default case for airtime or webhook rulesetss is the last migrated rule (failure)
			defaultCategoryUUID = uuids.UUID(converted.category.UUID())
		}

		kase, caseUI, err := migrateRule(baseLanguage, rule, converted.category, localization)
		if err != nil {
			return nil, nil, "", "", nil, err
		}

		if kase != nil {
			cases = append(cases, kase)

			if caseUI != nil {
				uiConfig.AddCaseConfig(kase.UUID(), caseUI)
			}
		}
	}

	return cases, categories, defaultCategoryUUID, timeoutCategoryUUID, exits, nil
}

// migrates the given legacy rule to a router case
func migrateRule(baseLanguage envs.Language, r Rule, category migratedCategory, localization migratedLocalization) (migratedCase, map[string]interface{}, error) {
	newType := testTypeMappings[r.Test.Type]
	var arguments []string
	var err error

	caseUUID := uuids.New()
	var caseUI map[string]interface{}

	switch r.Test.Type {

	// tests that take no arguments
	case "date", "has_email", "not_empty", "number", "phone", "state":
		arguments = []string{}

	// tests against a single numeric value
	case "eq", "gt", "gte", "lt", "lte":
		test := numericTest{}
		err = jsonx.Unmarshal(r.Test.Data, &test)
		if err != nil {
			return nil, nil, err
		}
		migratedTest, _ := expressions.MigrateTemplate(string(test.Test), nil)
		arguments = []string{migratedTest}

	case "between":
		test := betweenTest{}
		err = jsonx.Unmarshal(r.Test.Data, &test)
		if err != nil {
			return nil, nil, err
		}

		migratedMin, _ := expressions.MigrateTemplate(test.Min, nil)
		migratedMax, _ := expressions.MigrateTemplate(test.Max, nil)

		arguments = []string{migratedMin, migratedMax}

	// tests against a single localized string
	case "contains", "contains_any", "contains_phrase", "contains_only_phrase", "regex", "starts":
		test := localizedStringTest{}
		err = jsonx.Unmarshal(r.Test.Data, &test)
		if err != nil {
			return nil, nil, err
		}

		baseTest := test.Test.Base(baseLanguage)

		// all the tests are evaluated as templates.. except regex
		if r.Test.Type != "regex" {
			baseTest, _ = expressions.MigrateTemplate(baseTest, nil)

		}
		arguments = []string{baseTest}

		localization.addTranslationMap(baseLanguage, test.Test, caseUUID, "arguments")

	// tests against a single date value
	case "date_equal", "date_after", "date_before":
		test := stringTest{}
		err = jsonx.Unmarshal(r.Test.Data, &test)
		if err != nil {
			return nil, nil, err
		}
		migratedTest, _ := expressions.MigrateTemplate(test.Test, &expressions.MigrateOptions{RawDates: true})

		var delta int
		match := relativeDateTest.FindStringSubmatch(test.Test)
		if match != nil {
			delta, _ = strconv.Atoi(match[2])
			if match[1] == "-" {
				delta = -delta
			}
		}

		arguments = []string{migratedTest}

		caseUI = map[string]interface{}{
			"arguments": []string{strconv.Itoa(delta)},
		}

	// tests against a single group value
	case "in_group":
		test := groupTest{}
		err = jsonx.Unmarshal(r.Test.Data, &test)
		arguments = []string{string(test.Test.UUID), string(test.Test.Name)}

	case "subflow":
		newType = "has_only_text"
		test := subflowTest{}
		err = jsonx.Unmarshal(r.Test.Data, &test)
		arguments = []string{test.ExitType}

	case "webhook_status":
		newType = "has_only_text"
		test := webhookTest{}
		err = jsonx.Unmarshal(r.Test.Data, &test)
		if test.Status == "success" {
			arguments = []string{"Success"}
		} else {
			arguments = []string{"Failure"}
		}

	case "airtime_status":
		newType = "has_category"
		test := airtimeTest{}
		err = jsonx.Unmarshal(r.Test.Data, &test)
		if test.ExitStatus == "success" {
			arguments = []string{"Success"}
		} else {
			return nil, nil, nil // failure just becomes default category
		}

	case "district":
		test := stringTest{}
		err = jsonx.Unmarshal(r.Test.Data, &test)
		if err != nil {
			return nil, nil, err
		}

		migratedState, _ := expressions.MigrateTemplate(test.Test, nil)

		arguments = []string{migratedState}

	case "ward":
		test := wardTest{}
		err = jsonx.Unmarshal(r.Test.Data, &test)
		if err != nil {
			return nil, nil, err
		}

		migratedDistrict, _ := expressions.MigrateTemplate(test.District, nil)
		migratedState, _ := expressions.MigrateTemplate(test.State, nil)

		arguments = []string{migratedDistrict, migratedState}

	default:
		return nil, nil, errors.Errorf("migration of '%s' tests not supported", r.Test.Type)
	}

	return newCase(caseUUID, newType, arguments, category.UUID()), caseUI, err
}

// migrates the given legacy actionset to a node with a set of migrated actions and a single exit
func migrateActionSet(lang envs.Language, a ActionSet, validDests map[uuids.UUID]bool, localization migratedLocalization, baseMediaURL string) (migratedNode, error) {
	actions := make([]migratedAction, len(a.Actions))

	// migrate each action
	for i := range a.Actions {
		action, err := migrateAction(lang, a.Actions[i], localization, baseMediaURL)
		if err != nil {
			return nil, errors.Wrapf(err, "error migrating action[type=%s]", a.Actions[i].Type)
		}
		actions[i] = action
	}

	// only set exit destination if it's valid
	var destinationUUID uuids.UUID
	if validDests[a.Destination] {
		destinationUUID = a.Destination
	}

	exit := newExit(a.ExitUUID, destinationUUID)

	return newNode(a.UUID, actions, nil, []migratedExit{exit}), nil
}

func readLegacyFlow(data json.RawMessage) (*Flow, error) {
	f := &Flow{}
	if err := utils.UnmarshalAndValidate(data, f); err != nil {
		return nil, err
	}

	if f.Metadata == nil {
		f.Metadata = &Metadata{}
	}

	return f, nil
}

func migrateNodes(f *Flow, baseMediaURL string) ([]migratedNode, map[uuids.UUID]*NodeUI, migratedLocalization, error) {
	localization := make(migratedLocalization)
	numNodes := len(f.ActionSets) + len(f.RuleSets)
	nodes := make([]migratedNode, numNodes)
	nodeUIs := make(map[uuids.UUID]*NodeUI, numNodes)

	// get set of all node UUIDs, i.e. the valid destinations for any exit
	validDestinations := make(map[uuids.UUID]bool, numNodes)
	for _, as := range f.ActionSets {
		validDestinations[as.UUID] = true
	}
	for _, rs := range f.RuleSets {
		validDestinations[rs.UUID] = true
	}

	for i, actionSet := range f.ActionSets {
		node, err := migrateActionSet(f.BaseLanguage, actionSet, validDestinations, localization, baseMediaURL)
		if err != nil {
			return nil, nil, nil, errors.Wrapf(err, "error migrating action_set[uuid=%s]", actionSet.UUID)
		}
		nodes[i] = node
		nodeUIs[node.UUID()] = NewNodeUI(UINodeTypeActionSet, actionSet.X, actionSet.Y, nil)
	}

	for i, ruleSet := range f.RuleSets {
		node, uiType, uiNodeConfig, err := migrateRuleSet(f.BaseLanguage, ruleSet, validDestinations, localization)
		if err != nil {
			return nil, nil, nil, errors.Wrapf(err, "error migrating rule_set[uuid=%s]", ruleSet.UUID)
		}
		nodes[len(f.ActionSets)+i] = node
		nodeUIs[node.UUID()] = NewNodeUI(uiType, ruleSet.X, ruleSet.Y, uiNodeConfig)
	}

	// make sure our entry node is first
	entryNodes := []migratedNode{}
	otherNodes := []migratedNode{}
	for _, node := range nodes {
		if node.UUID() == f.Entry {
			entryNodes = []migratedNode{node}
		} else {
			otherNodes = append(otherNodes, node)
		}
	}

	// and sort remaining nodes by their top position (Y)
	sort.SliceStable(otherNodes, func(i, j int) bool {
		u1 := nodeUIs[otherNodes[i].UUID()]
		u2 := nodeUIs[otherNodes[j].UUID()]

		if u1 != nil && u2 != nil {
			return u1.Position.Top < u2.Position.Top
		}
		return false
	})

	nodes = append(entryNodes, otherNodes...)

	return nodes, nodeUIs, localization, nil
}

// Migrate migrates this legacy flow to the new format
func (f *Flow) Migrate(baseMediaURL string) ([]byte, error) {
	nodes, nodeUIs, localization, err := migrateNodes(f, baseMediaURL)
	if err != nil {
		return nil, err
	}

	// build UI section
	ui := NewUI()
	for _, actionSet := range f.ActionSets {
		ui.AddNode(actionSet.UUID, nodeUIs[actionSet.UUID])
	}
	for _, ruleSet := range f.RuleSets {
		ui.AddNode(ruleSet.UUID, nodeUIs[ruleSet.UUID])
	}
	for _, note := range f.Metadata.Notes {
		ui.AddSticky(note.Migrate())
	}

	uuid := f.Metadata.UUID
	name := f.Metadata.Name

	// some flows have these set on root-level instead.. or not set at all
	if uuid == "" {
		uuid = f.UUID
		if uuid == "" {
			uuid = uuids.New()
		}
	}
	if name == "" {
		name = f.Name
	}

	migrated := map[string]interface{}{
		"uuid":                 uuid,
		"name":                 name,
		"spec_version":         "13.0.0",
		"language":             f.BaseLanguage,
		"type":                 flowTypeMapping[f.FlowType],
		"revision":             f.Metadata.Revision,
		"expire_after_minutes": f.Metadata.Expires,
		"localization":         localization,
		"nodes":                nodes,
		"_ui":                  ui,
	}

	return jsonx.Marshal(migrated)
}

// IsPossibleDefinition peeks at the given flow definition to determine if it could be in legacy format
func IsPossibleDefinition(data json.RawMessage) bool {
	// any JSON blob with one of the following keys could be a legacy definition
	frag1, _, _, _ := jsonparser.Get(data, "action_sets")
	frag2, _, _, _ := jsonparser.Get(data, "rule_sets")
	frag3, _, _, _ := jsonparser.Get(data, "flow_type")
	return frag1 != nil || frag2 != nil || frag3 != nil
}

// MigrateDefinition migrates a legacy definition to 13.0.0
func MigrateDefinition(data json.RawMessage, baseMediaURL string) (json.RawMessage, error) {
	legacyFlow, err := readLegacyFlow(data)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read legacy flow")
	}

	return legacyFlow.Migrate(baseMediaURL)
}
