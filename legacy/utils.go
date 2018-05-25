package legacy

import (
	"encoding/json"

	"github.com/nyaruka/goflow/utils"
)

// Translations is an inline translation map used for localization
type Translations map[utils.Language]string

// Base looks up the translation in the given base language, or "base"
func (t Translations) Base(baseLanguage utils.Language) string {
	val, exists := t[baseLanguage]
	if exists {
		return val
	}
	return t["base"]
}

// UnmarshalJSON unmarshals legacy translations from the given JSON
func (t *Translations) UnmarshalJSON(data []byte) error {
	// sometimes legacy flows have a single string instead of a map
	if data[0] == '"' {
		var asString string
		if err := json.Unmarshal(data, &asString); err != nil {
			return err
		}
		*t = Translations{"base": asString}
		return nil
	}

	asMap := make(map[utils.Language]string)
	if err := json.Unmarshal(data, &asMap); err != nil {
		return err
	}

	*t = asMap
	return nil
}

// DecimalString represents a decimal value which may be provided as a string or floating point value
type DecimalString string

// UnmarshalJSON unmarshals a decimal string from the given JSON
func (s *DecimalString) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		// data is a quoted string
		*s = DecimalString(data[1 : len(data)-1])
	} else {
		// data is JSON float
		*s = DecimalString(data)
	}
	return nil
}
