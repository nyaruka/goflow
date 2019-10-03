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
	l.AddItemTranslation(envs.Language("eng"), uuids.UUID("f50df34b-18f8-489b-b8e8-ccb14d720641"), "foo", []string{"Hola"})

	thing := &testFlowThing{UUID: uuids.UUID("f50df34b-18f8-489b-b8e8-ccb14d720641"), Foo: "Hello", Bar: "World"}

	templates := make([]string, 0)
	inspect.Templates(thing, l, func(t string) {
		templates = append(templates, t)
	})

	assert.Equal(t, []string{"Hello", "Hola", "World"}, templates)

	// can also extract from slice of things
	templates = make([]string, 0)
	inspect.Templates([]*testFlowThing{thing}, l, func(t string) {
		templates = append(templates, t)
	})

	assert.Equal(t, []string{"Hello", "Hola", "World"}, templates)

	// or a slice of actions
	actions := []flows.Action{
		actions.NewSetContactName(flows.ActionUUID("d5ecd045-a15f-467c-925a-54bcdc726b9f"), "Bob"),
		actions.NewSetContactLanguage(flows.ActionUUID("d5ecd045-a15f-467c-925a-54bcdc726b9f"), "Gibberish"),
	}

	templates = make([]string, 0)
	inspect.Templates(actions, nil, func(t string) {
		templates = append(templates, t)
	})

	assert.Equal(t, []string{"Bob", "Gibberish"}, templates)
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
		"$.nodes[*].actions[@.type=\"call_webhook\"].body",
		"$.nodes[*].actions[@.type=\"call_webhook\"].headers[*]",
		"$.nodes[*].actions[@.type=\"call_webhook\"].url",
		"$.nodes[*].actions[@.type=\"classify_text\"].input",
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

func TestExtractFieldReferences(t *testing.T) {
	testCases := []struct {
		template string
		refs     []*assets.FieldReference
	}{
		{``, []*assets.FieldReference{}},
		{`Hi @contact`, []*assets.FieldReference{}},
		{`You are @fields.age`, []*assets.FieldReference{assets.NewFieldReference("age", "")}},
		{`You are @contact.fields.age`, []*assets.FieldReference{assets.NewFieldReference("age", "")}},
		{`You are @CONTACT.FIELDS.AGE`, []*assets.FieldReference{assets.NewFieldReference("age", "")}},
		{`You are @parent.fields.age`, []*assets.FieldReference{assets.NewFieldReference("age", "")}},
		{`You are @parent.contact.fields.age`, []*assets.FieldReference{assets.NewFieldReference("age", "")}},
		{`You are @child.contact.fields.age today`, []*assets.FieldReference{assets.NewFieldReference("age", "")}},
		{`You are @(ABS(contact . fields . age) + 1)`, []*assets.FieldReference{assets.NewFieldReference("age", "")}},
		{`You are @CONTACT.fields.age on @(contact.fields["Birthday"])`, []*assets.FieldReference{
			assets.NewFieldReference("age", ""),
			assets.NewFieldReference("birthday", ""),
		}},
	}

	for _, tc := range testCases {
		actual := inspect.ExtractFieldReferences(tc.template)

		assert.Equal(t, tc.refs, actual, "field refs mismatch for template '%s'", tc.template)
	}
}
