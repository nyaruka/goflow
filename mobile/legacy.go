package mobile

import (
	"github.com/nyaruka/goflow/flows/definition/legacy"
)

// IsLegacyDefinition peeks at the given flow definition to determine if it is in legacy format
func IsLegacyDefinition(definition string) bool {
	return legacy.IsLegacyDefinition([]byte(definition))
}

// MigrateLegacyDefinition migrates a legacy definition
func MigrateLegacyDefinition(definition string) (string, error) {
	migrated, err := legacy.MigrateLegacyDefinition([]byte(definition), "")
	if err != nil {
		return "", err
	}
	return string(migrated), nil
}
