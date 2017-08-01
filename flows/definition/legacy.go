package definition

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"

	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/routers"
	"github.com/nyaruka/goflow/flows/waits"
	"github.com/nyaruka/goflow/utils"
	"github.com/satori/go.uuid"
)

// LegacyFlow imports an old-world flow so it can be exported anew
type LegacyFlow struct {
	flow
	envelope legacyFlowEnvelope
}

type legacyFlowEnvelope struct {
	BaseLanguage utils.Language         `json:"base_language"`
	Metadata     legacyMetadataEnvelope `json:"metadata"`
	RuleSets     []legacyRuleSet        `json:"rule_sets"`
	ActionSets   []legacyActionSet      `json:"action_sets"`
	Entry        flows.NodeUUID         `json:"entry"`
}

type legacyMetadataEnvelope struct {
	UUID flows.FlowUUID `json:"uuid"`
	Name string         `json:"name"`
}

type legacyRule struct {
	UUID        flows.ExitUUID            `json:"uuid"`
	Destination flows.NodeUUID            `json:"destination"`
	Test        utils.TypedEnvelope       `json:"test"`
	Category    map[utils.Language]string `json:"category"`
	ExitType    string                    `json:"exit_type"`
}

type legacyRuleSet struct {
	Y       int             `json:"y"`
	X       int             `json:"x"`
	UUID    flows.NodeUUID  `json:"uuid"`
	Type    string          `json:"ruleset_type"`
	Label   string          `json:"label"`
	Operand string          `json:"operand"`
	Rules   []legacyRule    `json:"rules"`
	Config  json.RawMessage `json:"config"`
}

type legacyActionSet struct {
	Y           int            `json:"y"`
	X           int            `json:"x"`
	Destination flows.NodeUUID `json:"destination"`
	UUID        flows.NodeUUID `json:"uuid"`
	Actions     []legacyAction `json:"actions"`
}

type legacyLabelReference struct {
	UUID flows.LabelUUID `json:"uuid"`
	Name string          `json:"name"`
}

func (l *legacyLabelReference) Migrate() *flows.Label {
	return flows.NewLabel(l.UUID, l.Name)
}

type legacyContactReference struct {
	UUID flows.ContactUUID `json:"uuid"`
}

func (c *legacyContactReference) Migrate() *flows.ContactReference {
	return flows.NewContactReference(c.UUID, "")
}

type legacyGroupReference struct {
	UUID flows.GroupUUID `json:"uuid"`
	Name string          `json:"name"`
}

func (g *legacyGroupReference) Migrate() *flows.Group {
	return flows.NewGroup(g.UUID, g.Name)
}

type legacyVariable struct {
	ID string `json:"id"`
}

type legacyFlowReference struct {
	UUID flows.FlowUUID `json:"uuid"`
	Name string         `json:"name"`
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
	Msg     json.RawMessage `json:"msg"`
	Media   json.RawMessage `json:"media"`
	SendAll bool            `json:"send_all"`

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

type betweenTest struct {
	Min string `json:"min"`
	Max string `json:"max"`
}

type groupTest struct {
	Test legacyGroupReference `json:"test"`
}

type localizations map[utils.Language]flows.Action

// ReadLegacyFlows reads in legacy formatted flows
func ReadLegacyFlows(data json.RawMessage) ([]LegacyFlow, error) {
	var flows []LegacyFlow
	err := json.Unmarshal(data, &flows)
	return flows, err
}

type translationMap map[utils.Language]string

func addTranslationMap(baseLanguage utils.Language, translations *flowTranslations, mapped translationMap, uuid flows.UUID, key string) {
	for language, translation := range mapped {
		expression, _ := excellent.MigrateTemplate(translation)
		if language != baseLanguage {
			addTranslation(translations, language, uuid, key, []string{expression})
		}
	}
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
	"webhook_status":       "has_webhook_status",
}

func createAction(baseLanguage utils.Language, a legacyAction, fieldMap map[string]flows.FieldUUID, translations *flowTranslations) (flows.Action, error) {

	if a.UUID == "" {
		a.UUID = flows.ActionUUID(uuid.NewV4().String())
	}

	switch a.Type {
	case "add_label":

		labels := make([]*flows.Label, len(a.Labels))
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
			Emails:     migratedEmails,
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
			ChannelUUID: a.Channel,
			ChannelName: a.Name,
			BaseAction:  actions.NewBaseAction(a.UUID),
		}, nil
	case "flow":
		return &actions.StartFlowAction{
			FlowUUID:   a.Flow.UUID,
			FlowName:   a.Flow.Name,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil
	case "reply", "send":
		msg := make(map[utils.Language]string)
		media := make(map[utils.Language]string)

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

		addTranslationMap(baseLanguage, translations, msg, flows.UUID(a.UUID), "text")
		addTranslationMap(baseLanguage, translations, media, flows.UUID(a.UUID), "attachments")

		migratedText, _ := excellent.MigrateTemplate(msg[baseLanguage])
		migratedMedia, _ := excellent.MigrateTemplate(media[baseLanguage])

		attachments := []string{}
		if migratedMedia != "" {
			attachments = append(attachments, migratedMedia)
		}

		if a.Type == "reply" {
			return &actions.ReplyAction{
				Text:        migratedText,
				Attachments: attachments,
				BaseAction:  actions.NewBaseAction(a.UUID),
				AllURNs:     a.SendAll,
			}, nil
		}

		contacts := make([]*flows.ContactReference, len(a.Contacts))
		for i, contact := range a.Contacts {
			contacts[i] = contact.Migrate()
		}
		groups := make([]*flows.Group, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = group.Migrate()
		}

		return &actions.SendMsgAction{
			Text:        migratedText,
			Attachments: attachments,
			BaseAction:  actions.NewBaseAction(a.UUID),
			URNs:        []flows.URN{},
			Contacts:    contacts,
			Groups:      groups,
		}, nil

	case "add_group":
		groups := make([]*flows.Group, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = group.Migrate()
		}

		return &actions.AddToGroupAction{
			Groups:     groups,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil
	case "del_group":
		groups := make([]*flows.Group, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = group.Migrate()
		}

		return &actions.RemoveFromGroupAction{
			Groups:     groups,
			BaseAction: actions.NewBaseAction(a.UUID),
		}, nil
	case "save":
		fieldUUID, ok := fieldMap[a.Value]
		if !ok {
			fieldUUID = flows.FieldUUID(uuid.NewV4().String())
			fieldMap[a.Value] = fieldUUID
		}

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

		return &actions.SaveContactField{
			FieldName:  a.Label,
			Value:      migratedValue,
			FieldUUID:  fieldUUID,
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

func createCase(baseLanguage utils.Language, exitMap map[string]flows.Exit, r legacyRule, translations *flowTranslations) (routers.Case, error) {
	category := r.Category[baseLanguage]

	newType, _ := testTypeMappings[r.Test.Type]
	var arguments []string
	var err error

	caseUUID := flows.UUID(uuid.NewV4().String())

	switch r.Test.Type {

	// tests that take no arguments
	case "date", "email", "not_empty", "number", "phone":
		arguments = []string{}

	// tests against a single numeric value
	case "eq", "gt", "gte", "lt", "lte":
		test := stringTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		migratedTest, err := excellent.MigrateTemplate(test.Test)
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
		newType = "has_any_word"
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

	default:
		return routers.Case{}, fmt.Errorf("Migration of '%s' tests no supported", r.Test.Type)
	}

	// TODO
	// airtime_status
	// ward / district / state
	// interrupted_status
	// timeout

	return routers.Case{
		UUID:      caseUUID,
		Type:      newType,
		ExitUUID:  exitMap[category].UUID(),
		Arguments: arguments,
	}, err
}

type categoryName struct {
	uuid         flows.ExitUUID
	destination  flows.NodeUUID
	translations translationMap
	order        int
}

func parseRules(baseLanguage utils.Language, r legacyRuleSet, translations *flowTranslations) ([]flows.Exit, []routers.Case, flows.ExitUUID) {

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

		exits[category.order] = &exit{
			name:        k,
			uuid:        category.uuid,
			destination: category.destination,
		}
		exitMap[k] = exits[category.order]
	}

	var defaultExit flows.ExitUUID

	// create any cases to map to our new exits
	var cases []routers.Case
	for i := range r.Rules {
		if r.Rules[i].Test.Type != "true" {

			c, err := createCase(baseLanguage, exitMap, r.Rules[i], translations)
			if err == nil {
				cases = append(cases, c)
			} else if r.Rules[i].Test.Type == "webhook_status" {
				// webhook failures don't have a case, instead they are the default
				defaultExit = exitMap[r.Rules[i].Category[baseLanguage]].UUID()
			}
		} else {
			// take the first true rule as our default exit
			if defaultExit == "" {
				defaultExit = exitMap[r.Rules[i].Category[baseLanguage]].UUID()
			}
		}
	}

	return exits, cases, defaultExit
}

type fieldConfig struct {
	FieldDelimiter string `json:"field_delimiter"`
	FieldIndex     int    `json:"field_index"`
}

func createRuleNode(lang utils.Language, r legacyRuleSet, translations *flowTranslations) (*node, error) {
	node := &node{}
	node.uuid = r.UUID

	exits, cases, defaultExit := parseRules(lang, r, translations)
	resultName := r.Label

	switch r.Type {
	case "subflow":
		// subflow rulesets operate on the child flow status
		node.router = routers.NewSwitchRouter(defaultExit, "@child.status", cases, resultName)

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
				FlowUUID:   flowUUID,
				FlowName:   flowName,
			},
		}

		node.wait = &waits.FlowWait{
			FlowUUID: flowUUID,
		}

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
		node.router = routers.NewSwitchRouter(defaultExit, "@run.webhook", cases, resultName)

	case "form_field":
		var config fieldConfig
		json.Unmarshal(r.Config, &config)

		operand, _ := excellent.MigrateTemplate(r.Operand)
		operand = fmt.Sprintf("@(field(%s, %d, \"%s\"))", operand[1:], config.FieldIndex, config.FieldDelimiter)
		node.router = routers.NewSwitchRouter(defaultExit, operand, cases, resultName)

	case "wait_message":
		// TODO: add in timeout
		node.wait = &waits.MsgWait{}

		fallthrough
	case "flow_field":
		fallthrough
	case "group":
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

func createActionNode(lang utils.Language, a legacyActionSet, fieldMap map[string]flows.FieldUUID, translations *flowTranslations) *node {
	node := &node{}

	node.uuid = a.UUID
	node.actions = make([]flows.Action, len(a.Actions))
	for i := range a.Actions {
		action, err := createAction(lang, a.Actions[i], fieldMap, translations)
		if err == nil {
			node.actions[i] = action
		} else {
			fmt.Println(err)
		}
	}

	node.exits = make([]flows.Exit, 1)
	node.exits[0] = &exit{
		destination: a.Destination,
		uuid:        flows.ExitUUID(uuid.NewV4().String()),
	}
	return node

}

// UnmarshalJSON imports our JSON into a LegacyFlow object
func (f *LegacyFlow) UnmarshalJSON(data []byte) error {
	var envelope legacyFlowEnvelope
	var err error

	err = json.Unmarshal(data, &envelope)
	if err != nil {
		return err
	}

	fieldMap := make(map[string]flows.FieldUUID)

	f.language = envelope.BaseLanguage
	f.name = envelope.Metadata.Name
	f.uuid = envelope.Metadata.UUID

	translations := &flowTranslations{}

	f.nodes = make([]flows.Node, len(envelope.ActionSets)+len(envelope.RuleSets))
	for i := range envelope.ActionSets {
		node := createActionNode(f.language, envelope.ActionSets[i], fieldMap, translations)
		f.nodes[i] = node
	}

	for i := range envelope.RuleSets {
		node, err := createRuleNode(f.language, envelope.RuleSets[i], translations)
		if err != nil {
			return err
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

	//for i := range f.nodes {
	//	fmt.Println("Node:", f.nodes[i])
	//}

	return err
}

// MarshalJSON sends turns our legacy flow into bytes
func (f *LegacyFlow) MarshalJSON() ([]byte, error) {

	var fe = flowEnvelope{}
	fe.Name = f.name
	fe.Language = f.language
	fe.UUID = f.uuid

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
