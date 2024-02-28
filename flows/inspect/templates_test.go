package inspect_test

import (
	"testing"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/inspect"

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
	l.SetItemTranslation(i18n.Language("spa"), uuids.UUID("f50df34b-18f8-489b-b8e8-ccb14d720641"), "foo", []string{"Hola"})

	thing := &testFlowThing{UUID: uuids.UUID("f50df34b-18f8-489b-b8e8-ccb14d720641"), Foo: "Hello", Bar: "World"}

	templates := make(map[i18n.Language][]string)
	inspect.Templates(thing, l, func(l i18n.Language, t string) {
		templates[l] = append(templates[l], t)
	})

	assert.Equal(t, map[i18n.Language][]string{"": {"Hello", "World"}, "spa": {"Hola"}}, templates)

	// can also extract from slice of things
	templates = make(map[i18n.Language][]string)
	inspect.Templates([]*testFlowThing{thing}, l, func(l i18n.Language, t string) {
		templates[l] = append(templates[l], t)
	})

	assert.Equal(t, map[i18n.Language][]string{"": {"Hello", "World"}, "spa": {"Hola"}}, templates)

	// or a slice of actions
	actions := []flows.Action{
		actions.NewSetContactName(flows.ActionUUID("d5ecd045-a15f-467c-925a-54bcdc726b9f"), "Bob"),
		actions.NewSetContactLanguage(flows.ActionUUID("d5ecd045-a15f-467c-925a-54bcdc726b9f"), "Gibberish"),
	}

	templates = make(map[i18n.Language][]string)
	inspect.Templates(actions, nil, func(l i18n.Language, t string) {
		templates[l] = append(templates[l], t)
	})

	assert.Equal(t, map[i18n.Language][]string{"": {"Bob", "Gibberish"}}, templates)
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
