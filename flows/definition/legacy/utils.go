package legacy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

// Translations is an inline translation map used for localization
type Translations map[envs.Language]string

// ReadTranslations reads a translations map
func ReadTranslations(data json.RawMessage) (Translations, error) {
	if data == nil {
		return nil, nil
	}
	t := make(Translations)
	return t, jsonx.Unmarshal(data, &t)
}

// Base looks up the translation in the given base language, or "base"
func (t Translations) Base(baseLanguage envs.Language) string {
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
		if err := jsonx.Unmarshal(data, &asString); err != nil {
			return err
		}
		*t = Translations{"base": asString}
		return nil
	}

	asMap := make(map[envs.Language]string)
	if err := jsonx.Unmarshal(data, &asMap); err != nil {
		return err
	}

	*t = asMap
	return nil
}

// TransformTranslations transforms a list of single item translations into a map of multi-item translations, e.g.
//
// [{"eng": "yes", "fra": "oui"}, {"eng": "no", "fra": "non"}] becomes {"eng": ["yes", "no"], "fra": ["oui", "non"]}
//
func TransformTranslations(items []Translations) map[envs.Language][]string {
	// re-organize into a map of arrays
	transformed := make(map[envs.Language][]string)

	for i := range items {
		for language, translation := range items[i] {
			perLanguage, found := transformed[language]
			if !found {
				perLanguage = make([]string, len(items))
				transformed[language] = perLanguage
			}
			perLanguage[i] = translation
		}
	}
	return transformed
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
		return errors.Errorf("expected string or number, not %s", string(c))
	}
	return nil
}

// TypedEnvelope represents a json blob with a type property
type TypedEnvelope struct {
	Type string `json:"type" validate:"required"`
	Data []byte `json:"-"`
}

type typeOnly struct {
	Type string `json:"type" validate:"required"`
}

// UnmarshalJSON unmarshals a typed envelope from the given JSON
func (e *TypedEnvelope) UnmarshalJSON(b []byte) error {
	t := &typeOnly{}
	if err := utils.UnmarshalAndValidate(b, t); err != nil {
		return err
	}
	e.Type = t.Type
	e.Data = make([]byte, len(b))
	copy(e.Data, b)
	return nil
}

// URLJoin joins two URL parts with /
func URLJoin(base, relative string) string {

	// noop if we are already absolute
	if strings.HasPrefix(relative, "http://") || strings.HasPrefix(relative, "https://") {
		return relative
	}

	base = strings.TrimSuffix(base, "/")
	relative = strings.TrimPrefix(relative, "/")

	return fmt.Sprintf("%s/%s", base, relative)
}
