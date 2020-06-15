package contactql

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// Inspection holds the result of inspecting a query
type Inspection struct {
	Attributes   []string                 `json:"attributes"`
	Schemes      []string                 `json:"schemes"`
	Fields       []*assets.FieldReference `json:"fields"`
	Groups       []*assets.GroupReference `json:"groups"`
	AllowAsGroup bool                     `json:"allow_as_group"`
}

// Inspect extracts information about a query
func Inspect(query *ContactQuery) *Inspection {
	attributes := make(map[string]bool)
	schemes := make(map[string]bool)
	fieldRefs := make([]*assets.FieldReference, 0)
	groupRefs := make([]*assets.GroupReference, 0)
	refsSeen := make(map[string]bool)

	walk(query.Root(), func(c *Condition) {
		switch c.propType {
		case PropertyTypeAttribute:
			attributes[c.propKey] = true
		case PropertyTypeScheme:
			schemes[c.propKey] = true
		}

		switch typed := c.reference.(type) {
		case *assets.FieldReference:
			if !refsSeen[typed.String()] {
				fieldRefs = append(fieldRefs, typed)
				refsSeen[typed.String()] = true
			}
		case *assets.GroupReference:
			if !refsSeen[typed.String()] {
				groupRefs = append(groupRefs, typed)
				refsSeen[typed.String()] = true
			}
		}
	})

	allowAsGroup := !(attributes[AttributeID] || attributes[AttributeGroup])

	return &Inspection{
		Attributes:   utils.StringSetKeys(attributes),
		Schemes:      utils.StringSetKeys(schemes),
		Fields:       fieldRefs,
		Groups:       groupRefs,
		AllowAsGroup: allowAsGroup,
	}
}

func walk(node QueryNode, conditionCallback func(*Condition)) {
	switch n := node.(type) {
	case *BoolCombination:
		for _, n := range n.Children() {
			walk(n, conditionCallback)
		}

	case *Condition:
		conditionCallback(n)
	}
}
