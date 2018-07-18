package transferto

import (
	"encoding/json"
	"strings"
)

// CSVList is a list of strings which can be automatically unmarshalled from a CSV list
type CSVList []string

// UnmarshalJSON unmarshals this list from a CSV string
func (l *CSVList) UnmarshalJSON(data []byte) error {
	var asString string
	if err := json.Unmarshal(data, &asString); err != nil {
		return err
	}
	*l = strings.Split(asString, ",")
	return nil
}
