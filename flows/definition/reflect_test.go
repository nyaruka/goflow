package definition_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/engine"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractTemplates(t *testing.T) {
	source, err := static.LoadSource("../../test/testdata/flows/two_questions.json")
	require.NoError(t, err)

	sessionAssets, err := engine.NewSessionAssets(source)
	require.NoError(t, err)

	flow, err := sessionAssets.Flows().Get("615b8a0f-588c-4d20-a05f-363b0b4ce6f4")
	require.NoError(t, err)

	assert.Equal(t, []string{
		`Hi @contact.name! What is your favorite color? (red/blue) Your number is @(format_urn(contact.urn))`,
		`Quelle est votres couleur preferee? (rouge/blue)`,
		`@input`,
		`red`,
		`rouge`,
		`blue`,
		`bleu`,
		`fra`,
		`@(TITLE(results.favorite_color.category_localized)) it is! What is your favorite soda? (pepsi/coke)`,
		`@(TITLE(results.favorite_color.category_localized))! Bien sur! Quelle est votes soda preferee? (pepsi/coke)`,
		`@input`,
		`pepsi`,
		`coke coca cola`,
		`http://localhost/?cmd=success`,
		`{ "contact": @(json(contact.uuid)), "soda": @(json(results.soda.value)) }`,
		`Great, you are done and like @results.soda! Webhook status was @results.webhook.value`,
		`Parfait, vous avez finis et tu aimes @results.soda.category`,
	}, definition.ExtractTemplates(flow))
}
