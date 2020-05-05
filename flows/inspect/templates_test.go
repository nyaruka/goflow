package inspect_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/inspect"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/stretchr/testify/assert"
)

type testFlowThing struct {
	UUID uuids.UUID `json:"uuid"`                             // not a template
	Foo  string     `json:"foo" engine:"evaluated,localized"` // a localizable template
	Bar  string     `json:"bar" engine:"evaluated"`           // a template
}

func (t *testFlowThing) LocalizationUUID() uuids.UUID {
	return t.UUID
}

func TestTemplates(t *testing.T) {
	l := definition.NewLocalization()
	l.SetItemTranslation(envs.Language("spa"), uuids.UUID("f50df34b-18f8-489b-b8e8-ccb14d720641"), "foo", []string{"Hola"})

	thing := &testFlowThing{UUID: uuids.UUID("f50df34b-18f8-489b-b8e8-ccb14d720641"), Foo: "Hello", Bar: "World"}

	templates := make(map[envs.Language][]string)
	inspect.Templates(thing, l, func(l envs.Language, t string) {
		templates[l] = append(templates[l], t)
	})

	assert.Equal(t, map[envs.Language][]string{"": []string{"Hello", "World"}, "spa": []string{"Hola"}}, templates)

	// can also extract from slice of things
	templates = make(map[envs.Language][]string)
	inspect.Templates([]*testFlowThing{thing}, l, func(l envs.Language, t string) {
		templates[l] = append(templates[l], t)
	})

	assert.Equal(t, map[envs.Language][]string{"": []string{"Hello", "World"}, "spa": []string{"Hola"}}, templates)

	// or a slice of actions
	actions := []flows.Action{
		actions.NewSetContactName(flows.ActionUUID("d5ecd045-a15f-467c-925a-54bcdc726b9f"), "Bob"),
		actions.NewSetContactLanguage(flows.ActionUUID("d5ecd045-a15f-467c-925a-54bcdc726b9f"), "Gibberish"),
	}

	templates = make(map[envs.Language][]string)
	inspect.Templates(actions, nil, func(l envs.Language, t string) {
		templates[l] = append(templates[l], t)
	})

	assert.Equal(t, map[envs.Language][]string{"": []string{"Bob", "Gibberish"}}, templates)
}

func TestTemplatePaths(t *testing.T) {
	paths := make([]string, 0)
	for typeName, fn := range actions.RegisteredTypes() {
		actionType := reflect.TypeOf(fn())

		inspect.TemplatePaths(actionType, fmt.Sprintf("$.nodes[*].actions[@.type=\"%s\"]", typeName), func(path string) {
			paths = append(paths, path)
		})
	}

	sort.Strings(paths)

	assert.Equal(t, []string{
		"$.nodes[*].actions[@.type=\"add_contact_groups\"].groups[*].name_match",
		"$.nodes[*].actions[@.type=\"add_contact_urn\"].path",
		"$.nodes[*].actions[@.type=\"add_input_labels\"].labels[*].name_match",
		"$.nodes[*].actions[@.type=\"call_classifier\"].input",
		"$.nodes[*].actions[@.type=\"call_webhook\"].body",
		"$.nodes[*].actions[@.type=\"call_webhook\"].headers[*]",
		"$.nodes[*].actions[@.type=\"call_webhook\"].url",
		"$.nodes[*].actions[@.type=\"open_ticket\"].body",
		"$.nodes[*].actions[@.type=\"open_ticket\"].subject",
		"$.nodes[*].actions[@.type=\"play_audio\"].audio_url",
		"$.nodes[*].actions[@.type=\"remove_contact_groups\"].groups[*].name_match",
		"$.nodes[*].actions[@.type=\"say_msg\"].text",
		"$.nodes[*].actions[@.type=\"send_broadcast\"].attachments[*]",
		"$.nodes[*].actions[@.type=\"send_broadcast\"].contact_query",
		"$.nodes[*].actions[@.type=\"send_broadcast\"].groups[*].name_match",
		"$.nodes[*].actions[@.type=\"send_broadcast\"].legacy_vars[*]",
		"$.nodes[*].actions[@.type=\"send_broadcast\"].quick_replies[*]",
		"$.nodes[*].actions[@.type=\"send_broadcast\"].text",
		"$.nodes[*].actions[@.type=\"send_email\"].addresses[*]",
		"$.nodes[*].actions[@.type=\"send_email\"].body",
		"$.nodes[*].actions[@.type=\"send_email\"].subject",
		"$.nodes[*].actions[@.type=\"send_msg\"].attachments[*]",
		"$.nodes[*].actions[@.type=\"send_msg\"].quick_replies[*]",
		"$.nodes[*].actions[@.type=\"send_msg\"].templating.variables[*]",
		"$.nodes[*].actions[@.type=\"send_msg\"].text",
		"$.nodes[*].actions[@.type=\"set_contact_field\"].value",
		"$.nodes[*].actions[@.type=\"set_contact_language\"].language",
		"$.nodes[*].actions[@.type=\"set_contact_name\"].name",
		"$.nodes[*].actions[@.type=\"set_contact_timezone\"].timezone",
		"$.nodes[*].actions[@.type=\"set_run_result\"].value",
		"$.nodes[*].actions[@.type=\"start_session\"].contact_query",
		"$.nodes[*].actions[@.type=\"start_session\"].groups[*].name_match",
		"$.nodes[*].actions[@.type=\"start_session\"].legacy_vars[*]",
	}, paths)
}

func TestExtractFromTemplate(t *testing.T) {
	testCases := []struct {
		template   string
		assetRefs  []assets.Reference
		parentRefs []string
	}{
		{``, []assets.Reference{}, []string{}},
		{`Hi @contact`, []assets.Reference{}, []string{}},
		{
			`You are @fields.age`,
			[]assets.Reference{assets.NewFieldReference("age", "")},
			[]string{},
		},
		{
			`You are @contact.fields.age`,
			[]assets.Reference{assets.NewFieldReference("age", "")},
			[]string{},
		},
		{
			`You are @CONTACT.FIELDS.AGE`,
			[]assets.Reference{assets.NewFieldReference("age", "")},
			[]string{},
		},
		{
			`You are @parent.fields.age`,
			[]assets.Reference{assets.NewFieldReference("age", "")},
			[]string{},
		},
		{
			`You are @parent.contact.fields.age`,
			[]assets.Reference{assets.NewFieldReference("age", "")},
			[]string{},
		},
		{
			`You are @child.contact.fields.age today`,
			[]assets.Reference{assets.NewFieldReference("age", "")},
			[]string{},
		},
		{
			`You are @(ABS(contact . fields . age) + 1)`,
			[]assets.Reference{assets.NewFieldReference("age", "")},
			[]string{},
		},
		{
			`You are @FIELDS.AGE in @GLOBALS.ORG_NAME from @PARENT.RESULTS.STATE `,
			[]assets.Reference{
				assets.NewFieldReference("age", ""),
				assets.NewGlobalReference("org_name", ""),
			},
			[]string{"state"},
		},
		{
			`You are @(fields["age"]) in @(globals["org_name"]) from @(parent.results["state"])`,
			[]assets.Reference{
				assets.NewFieldReference("age", ""),
				assets.NewGlobalReference("org_name", ""),
			},
			[]string{"state"},
		},
	}

	for _, tc := range testCases {
		assetRefs, parentRefs := inspect.ExtractFromTemplate(tc.template)

		assert.Equal(t, tc.assetRefs, assetRefs, "asset refs mismatch for template '%s'", tc.template)
		assert.Equal(t, tc.parentRefs, parentRefs, "parent result refs mismatch for template '%s'", tc.template)
	}
}
