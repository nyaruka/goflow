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
	refs := make([]assets.Reference, 0)
	refsSeen := make(map[string]bool)

	addRef := func(ref assets.Reference) {
		if !refsSeen[ref.String()] {
			refs = append(refs, ref)
			refsSeen[ref.String()] = true
		}
	}

	walk(query.Root(), func(c *Condition) {
		switch c.propType {
		case PropertyTypeAttribute:
			attributes[c.propKey] = true

			if c.propKey == AttributeGroup {
				if query.resolver != nil {
					group := query.resolver.ResolveGroup(c.value)
					addRef(assets.NewGroupReference(group.UUID(), group.Name()))
				} else {
					addRef(assets.NewVariableGroupReference(c.value))
				}
			}
		case PropertyTypeScheme:
			schemes[c.propKey] = true
		case PropertyTypeField:
			if query.resolver != nil {
				field := query.resolver.ResolveField(c.propKey)
				addRef(assets.NewFieldReference(field.Key(), field.Name()))
			} else {
				addRef(assets.NewFieldReference(c.propKey, ""))
			}

		}
	})

	fieldRefs := make([]*assets.FieldReference, 0)
	groupRefs := make([]*assets.GroupReference, 0)
	for _, ref := range refs {
		switch typed := ref.(type) {
		case *assets.FieldReference:
			fieldRefs = append(fieldRefs, typed)
		case *assets.GroupReference:
			groupRefs = append(groupRefs, typed)
		}
	}

	// can't turn a query into a group if it uses id, status, group, flow or history
	allowAsGroup := !(attributes[AttributeID] || attributes[AttributeStatus] || attributes[AttributeGroup] || attributes[AttributeFlow] || attributes[AttributeHistory])

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
