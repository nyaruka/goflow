package inspect_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/inspect"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

type testFlowThing struct {
	UUID utils.UUID `json:"uuid"`                             // not a template
	Foo  string     `json:"foo" engine:"evaluated,localized"` // a localizable template
	Bar  string     `json:"bar" engine:"evaluated"`           // a template
}

func (t *testFlowThing) LocalizationUUID() utils.UUID {
	return t.UUID
}

func TestTemplates(t *testing.T) {
	l := definition.NewLocalization()
	l.AddItemTranslation(utils.Language("eng"), utils.UUID("f50df34b-18f8-489b-b8e8-ccb14d720641"), "foo", []string{"Hola"})

	thing := &testFlowThing{UUID: utils.UUID("f50df34b-18f8-489b-b8e8-ccb14d720641"), Foo: "Hello", Bar: "World"}

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
		actions.NewSetContactNameAction(flows.ActionUUID("d5ecd045-a15f-467c-925a-54bcdc726b9f"), "Bob"),
		actions.NewSetContactLanguageAction(flows.ActionUUID("d5ecd045-a15f-467c-925a-54bcdc726b9f"), "Gibberish"),
	}

	templates = make([]string, 0)
	inspect.Templates(actions, nil, func(t string) {
		templates = append(templates, t)
	})

	assert.Equal(t, []string{"Bob", "Gibberish"}, templates)
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
