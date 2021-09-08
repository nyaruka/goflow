package static_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"

	"github.com/stretchr/testify/assert"
)

func TestLabel(t *testing.T) {
	label := static.NewLabel(assets.LabelUUID("f5263dca-469b-47c2-be4f-845d3a14eedf"), "Spam")
	assert.Equal(t, assets.LabelUUID("f5263dca-469b-47c2-be4f-845d3a14eedf"), label.UUID())
	assert.Equal(t, "Spam", label.Name())
}
