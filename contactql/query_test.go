package contactql_test

import (
	"testing"

	"github.com/nyaruka/goflow/contactql"

	"github.com/stretchr/testify/assert"
)

func TestEscapeValue(t *testing.T) {
	assert.Equal(t, `""`, contactql.EscapeValue(``))
	assert.Equal(t, `"bobby tables"`, contactql.EscapeValue(`bobby tables`))
	assert.Equal(t, `"\"\" OR (id = 1)"`, contactql.EscapeValue(`"" OR (id = 1)`))
	assert.Equal(t, `"\\\"foo"`, contactql.EscapeValue(`\"foo`))
}
