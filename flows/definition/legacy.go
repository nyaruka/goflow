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

type legacyFlow struct {
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

type legacyLabel struct {
	Name string          `json:"name"`
	UUID flows.LabelUUID `json:"uuid"`
}

type legacyGroup struct {
	Name string          `json:"name"`
	UUID flows.GroupUUID `json:"uuid"`
}

type legacyAction struct {
	Type string           `json:"type"`
	UUID flows.ActionUUID `json:"uuid"`
	Name string           `json:"name"`

	// message  and email
	Msg json.RawMessage `json:"msg"`

	// groups
	Groups []legacyGroup `json:"groups"`

	// save actions
	Field string `json:"field"`
	Value string `json:"value"`
	Label string `json:"label"`

	// set language
	Language utils.Language `json:"lang"`

	// webhook
	Action  string `json:"action"`
	Webhook string `json:"webhook"`

	// add lable action
	Labels []legacyLabel `json:"labels"`

	// Trigger flow
	Flow legacyFlowReference `json:"flow"`

	// channel
	Channel flows.ChannelUUID `json:"channel"`

	//email
	Emails  []string `json:"emails"`
	Subject string   `json:"subject"`
}

type legacyFlowReference struct {
	UUID flows.FlowUUID `json:"uuid"`
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
	Test string `jons:"test"`
}

type localizations map[utils.Language]flows.Action

// ReadLegacyFlows reads in legacy formatted flows
func ReadLegacyFlows(data json.RawMessage) ([]legacyFlow, error) {
	var flows []legacyFlow
	err := json.Unmarshal(data, &flows)
	return flows, err
}

type translationMap map[utils.Language]string

func addTranslationMap(baseLanguage utils.Language, translations *flowTranslations, mapped translationMap, uuid flows.UUID, key string) {
	for language, translation := range mapped {
		items := itemTranslations{}
		expression, _ := excellent.TranslateTemplate(translation)
		items[key] = expression
		if language != baseLanguage {
			addTranslation(baseLanguage, translations, language, uuid, items)
		}
	}
}

func addTranslation(baseLanguage utils.Language, translations *flowTranslations, lang utils.Language, uuid flows.UUID, items itemTranslations) {
	langTranslations, ok := (*translations)[lang]
	if !ok {
		langTranslations = &languageTranslations{}
	}

	(*langTranslations)[uuid] = items
	(*translations)[lang] = langTranslations
}

var testTranslations = map[string]string{
	"contains_any": "has_any_word",
	"not_empty":    "has_text",
}

func translateTest(test string) string {
	translated, ok := testTranslations[test]
	if ok {
		return translated
	}
	return test
}

func createAction(baseLanguage utils.Language, a legacyAction, fieldMap map[string]flows.FieldUUID, translations *flowTranslations) (flows.Action, error) {

	if a.UUID == "" {
		a.UUID = flows.ActionUUID(uuid.NewV4().String())
	}

	switch a.Type {
	case "add_label":

		labels := make([]*flows.Label, len(a.Labels))
		for i := range labels {
			labels[i] = flows.NewLabel(a.Labels[i].UUID, a.Labels[i].Name)
		}

		return &actions.AddLabelAction{
			Labels: labels,
			BaseAction: actions.BaseAction{
				UUID: a.UUID,
			},
		}, nil

	case "email":
		var msg string
		err := json.Unmarshal(a.Msg, &msg)
		if err != nil {
			return nil, err
		}

		return &actions.EmailAction{
			Subject: a.Subject,
			Body:    msg,
			Emails:  a.Emails,
			BaseAction: actions.BaseAction{
				UUID: a.UUID,
			},
		}, nil

	case "lang":
		return &actions.SaveContactField{
			FieldUUID: "language",
			FieldName: "Language",
			Value:     string(a.Language),
			BaseAction: actions.BaseAction{
				UUID: a.UUID,
			},
		}, nil
	case "channel":
		return &actions.PreferredChannelAction{
			ChannelUUID: a.Channel,
			Name:        a.Name,
			BaseAction: actions.BaseAction{
				UUID: a.UUID,
			},
		}, nil
	case "flow":
		return &actions.StartFlowAction{
			FlowUUID: a.Flow.UUID,
			FlowName: a.Name,
			BaseAction: actions.BaseAction{
				UUID: a.UUID,
			},
		}, nil
	case "reply":
		msg := make(map[utils.Language]string)
		err := json.Unmarshal(a.Msg, &msg)
		if err != nil {
			return nil, err
		}

		addTranslationMap(baseLanguage, translations, msg, flows.UUID(a.UUID), "text")
		expression, _ := excellent.TranslateTemplate(msg[baseLanguage])
		return &actions.SendMsgAction{
			Text: expression,
			BaseAction: actions.BaseAction{
				UUID: a.UUID,
			},
		}, nil
	case "add_group":
		groups := make([]*flows.Group, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = flows.NewGroup(group.UUID, group.Name)
		}

		return &actions.AddToGroupAction{
			Groups: groups,
			BaseAction: actions.BaseAction{
				UUID: a.UUID,
			},
		}, nil
	case "del_group":
		groups := make([]*flows.Group, len(a.Groups))
		for i, group := range a.Groups {
			groups[i] = flows.NewGroup(group.UUID, group.Name)
		}

		return &actions.RemoveFromGroupAction{
			Groups: groups,
			BaseAction: actions.BaseAction{
				UUID: a.UUID,
			},
		}, nil
	case "save":
		fieldUUID, ok := fieldMap[a.Value]
		if !ok {
			fieldUUID = flows.FieldUUID(uuid.NewV4().String())
			fieldMap[a.Value] = fieldUUID
		}

		translated, _ := excellent.TranslateTemplate(a.Value)

		return &actions.SaveContactField{
			FieldName: a.Label,
			Value:     translated,
			FieldUUID: fieldUUID,
			BaseAction: actions.BaseAction{
				UUID: a.UUID,
			},
		}, nil
	case "set_language":
		return &actions.SaveContactField{
			FieldUUID: "language",
			FieldName: "Language",
			Value:     string(a.Value),
			BaseAction: actions.BaseAction{
				UUID: a.UUID,
			},
		}, nil
	case "api":
		translated, _ := excellent.TranslateTemplate(a.Webhook)
		return &actions.WebhookAction{
			Method: a.Action,
			URL:    translated,
			BaseAction: actions.BaseAction{
				UUID: a.UUID,
			},
		}, nil
	default:
		return nil, fmt.Errorf("couldn't create action for %s", a.Type)
	}
}

func createCase(baseLanguage utils.Language, exitMap map[string]flows.Exit, r legacyRule, translations *flowTranslations) (routers.Case, error) {
	category := r.Category[baseLanguage]
	testType := r.Test.Type
	var arguments []string
	var err error

	caseUUID := flows.UUID(uuid.NewV4().String())

	switch r.Test.Type {
	case "subflow":
		testType = "has_any_word"
		test := subflowTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		if test.ExitType == "completed" {
			arguments = []string{"C"}
		} else {
			arguments = []string{"E"}
		}
	case "webhook_status":
		test := webhookTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		if test.Status == "success" {
			testType = "has_webhook_status"
			arguments = []string{"S"}
		} else {
			return routers.Case{}, fmt.Errorf("No failure test")
		}
	case "eq":
		test := stringTest{}
		err = json.Unmarshal(r.Test.Data, &test)
		arguments = []string{test.Test}

	case "regex":
		fallthrough

	case "contains_any":
		test := localizedStringTest{}
		err = json.Unmarshal(r.Test.Data, &test)

		testType = translateTest(testType)

		// TODO: arguments should be an array
		addTranslationMap(baseLanguage, translations, test.Test, caseUUID, "arguments")
		arguments = []string{test.Test[baseLanguage]}
	default:
		testType = translateTest(testType)
	}

	return routers.Case{
		UUID:      caseUUID,
		Type:      testType,
		ExitUUID:  exitMap[category].UUID(),
		Arguments: arguments,
	}, err
}

type categoryName struct {
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
				destination:  r.Rules[i].Destination,
				translations: r.Rules[i].Category,
				order:        order,
			}
			order++
		}
	}

	// create exists for each category
	exits := make([]flows.Exit, len(categoryMap))
	exitMap := make(map[string]flows.Exit)
	for k, category := range categoryMap {
		uuid := flows.ExitUUID(uuid.NewV4().String())

		addTranslationMap(baseLanguage, translations, category.translations, flows.UUID(uuid), "label")

		exits[category.order] = &exit{
			name:        k,
			uuid:        uuid,
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
			defaultExit = exitMap[r.Rules[i].Category[baseLanguage]].UUID()
		}
	}

	return exits, cases, defaultExit
}

func createRuleNode(lang utils.Language, r legacyRuleSet, translations *flowTranslations) (*node, error) {
	node := &node{}
	node.uuid = r.UUID

	exits, cases, defaultExit := parseRules(lang, r, translations)

	switch r.Type {
	case "subflow":
		// subflow rulesets operate on the child flow status
		node.router = &routers.SwitchRouter{
			Default: defaultExit,
			Operand: "@child.status",
			Cases:   cases,
		}

		config := make(map[string]map[string]string)
		err := json.Unmarshal(r.Config, &config)
		if err != nil {
			return nil, err
		}

		flowUUID := flows.FlowUUID(config["flow"]["uuid"])
		flowName := config["flow"]["name"]

		node.actions = []flows.Action{
			&actions.StartFlowAction{
				FlowUUID: flowUUID,
				FlowName: flowName,
			},
		}

		node.wait = &waits.FlowWait{
			FlowUUID: flowUUID,
		}

	case "webhook":
		config := make(map[string]string)
		err := json.Unmarshal(r.Config, &config)
		if err != nil {
			return nil, err
		}

		node.actions = []flows.Action{
			&actions.WebhookAction{
				URL:    config["webhook"],
				Method: config["webhook_action"],
			},
		}

		// subflow rulesets operate on the child flow status
		node.router = &routers.SwitchRouter{
			Default: defaultExit,
			Operand: "@webhook",
			Cases:   cases,
		}

	case "flow_field":
		fallthrough
	case "form_field":
		fallthrough
	case "wait_message":

		// TODO: add in timeout
		node.wait = &waits.MsgWait{}

		fallthrough
	case "group":
		fallthrough
	case "contact_field":
		fallthrough
	case "expression":
		operand, _ := excellent.TranslateTemplate(r.Operand)
		fmt.Println(operand)
		node.router = &routers.SwitchRouter{
			Default: defaultExit,
			Operand: operand,
			Cases:   cases,
			BaseRouter: routers.BaseRouter{
				Name_: r.Label,
			},
		}
	case "random":
		node.router = &routers.RandomRouter{}
	default:
		fmt.Printf("No router for %s\n", r.Type)
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

func (f *legacyFlow) UnmarshalJSON(data []byte) error {
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
		f.nodes[len(envelope.ActionSets)+i], err = createRuleNode(f.language, envelope.RuleSets[i], translations)
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

func (f *legacyFlow) MarshalJSON() ([]byte, error) {

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
