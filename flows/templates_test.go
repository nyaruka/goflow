package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"

	"github.com/stretchr/testify/assert"
)

func TestExtractFieldReferences(t *testing.T) {
	testCases := []struct {
		template string
		refs     []*assets.FieldReference
	}{
		{``, []*assets.FieldReference{}},
		{`Hi @contact`, []*assets.FieldReference{}},
		{`You are @contact.fields.age`, []*assets.FieldReference{assets.NewFieldReference("age", "")}},
		{`You are @CONTACT.FIELDS.AGE`, []*assets.FieldReference{assets.NewFieldReference("age", "")}},
		{`You are @parent.contact.fields.age`, []*assets.FieldReference{assets.NewFieldReference("age", "")}},
		{`You are @child.contact.fields.age today`, []*assets.FieldReference{assets.NewFieldReference("age", "")}},
		{`You are @(ABS(contact . fields . age) + 1)`, []*assets.FieldReference{assets.NewFieldReference("age", "")}},
		{`You are @CONTACT.fields.age on @(contact.fields["Birthday"])`, []*assets.FieldReference{
			assets.NewFieldReference("age", ""),
			assets.NewFieldReference("birthday", ""),
		}},
	}

	for _, tc := range testCases {
		actual := flows.ExtractFieldReferences(tc.template)

		assert.Equal(t, tc.refs, actual, "field refs mismatch for template '%s'", tc.template)
	}
}
