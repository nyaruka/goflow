package legacy

import (
	"encoding/json"
	"fmt"

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

// StringOrNumber represents something we need to read as a string, but might actually be number value in the JSON source
type StringOrNumber string

// UnmarshalJSON unmarshals this from the given JSON
func (s *StringOrNumber) UnmarshalJSON(data []byte) error {
	c := data[0]
	if c == '"' {
		// data is a quoted string
		*s = StringOrNumber(data[1 : len(data)-1])
	} else if (c >= '0' && c <= '9') || c == '-' {
		// data is JSON number
		*s = StringOrNumber(data)
	} else {
		return fmt.Errorf("expected string or number, not %s", string(c))
	}
	return nil
}
