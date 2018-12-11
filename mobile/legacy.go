package mobile

import (
	"fmt"

	"github.com/nyaruka/goflow/legacy"
	"github.com/nyaruka/goflow/utils"
)

// MigrateLegacyFlow migrates a legacy flow definitin
func MigrateLegacyFlow(definition string) (string, error) {
	legacyFlow, err := legacy.ReadLegacyFlow([]byte(definition))
	if err != nil {
		return "", fmt.Errorf("unable to read legacy flow: %s", err)
	}

	flow, err := legacyFlow.Migrate(false, false)
	if err != nil {
		return "", fmt.Errorf("unable to migrate legacy flow: %s", err)
	}

	marshaled, err := utils.JSONMarshal(flow)
	if err != nil {
		return "", fmt.Errorf("unable to marshal migrated flow: %s", err)
	}

	return string(marshaled), nil
}
