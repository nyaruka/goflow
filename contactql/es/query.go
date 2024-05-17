package es

import (
	"fmt"
	"strings"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"
)

// AssetMapper is used to map engine assets to however ES identifies them
type AssetMapper interface {
	Flow(assets.Flow) int64
	Group(assets.Group) int64
}

// we store contact status in elastic as single char codes
var contactStatusCodes = map[string]string{
	"active":   "A",
	"blocked":  "B",
	"stopped":  "S",
	"archived": "V",
}

// ToElasticQuery converts a contactql query to an Elastic query
func ToElasticQuery(env envs.Environment, mapper AssetMapper, query *contactql.ContactQuery) map[string]any {
	if query.Resolver() == nil {
		panic("can only convert queries parsed with a resolver")
	}

	return nodeToElastic(env, query.Resolver(), mapper, query.Root())
}

func nodeToElastic(env envs.Environment, resolver contactql.Resolver, mapper AssetMapper, node contactql.QueryNode) map[string]any {
	switch n := node.(type) {
	case *contactql.BoolCombination:
		return boolCombination(env, resolver, mapper, n)
	case *contactql.Condition:
		return condition(env, resolver, mapper, n)
	default:
		panic(fmt.Sprintf("unsupported node type: %T", n))
	}
}

func boolCombination(env envs.Environment, resolver contactql.Resolver, mapper AssetMapper, combination *contactql.BoolCombination) map[string]any {
	queries := make([]map[string]any, len(combination.Children()))
	for i, child := range combination.Children() {
		queries[i] = nodeToElastic(env, resolver, mapper, child)
	}

	if combination.Operator() == contactql.BoolOperatorAnd {
		return All(queries...)
	}

	return Any(queries...)
}

func condition(env envs.Environment, resolver contactql.Resolver, mapper AssetMapper, c *contactql.Condition) map[string]any {
	switch c.PropertyType() {
	case contactql.PropertyTypeField:
		return fieldCondition(env, resolver, c)
	case contactql.PropertyTypeAttribute:
		return attributeCondition(env, resolver, mapper, c)
	case contactql.PropertyTypeURN:
		return schemeCondition(c)
	default:
		panic(fmt.Sprintf("unsupported property type: %s", c.PropertyType()))
	}
}

func fieldCondition(env envs.Environment, resolver contactql.Resolver, c *contactql.Condition) map[string]any {
	field := resolver.ResolveField(c.PropertyKey())
	fieldType := field.Type()
	fieldQuery := Term("fields.field", field.UUID())

	// special cases for set/unset
	if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && c.Value() == "" {
		query := Nested("fields", All(fieldQuery, Exists("fields."+string(fieldType))))

		// if we are looking for unset, inverse our query
		if c.Operator() == contactql.OpEqual {
			query = Not(query)
		}
		return query
	}

	if fieldType == assets.FieldTypeText {
		value := strings.ToLower(c.Value())

		switch c.Operator() {
		case contactql.OpEqual:
			return Nested("fields", All(fieldQuery, Term("fields.text", value)))
		case contactql.OpNotEqual:
			query := All(fieldQuery, Term("fields.text", value), Exists("fields.text"))
			return Not(Nested("fields", query))
		default:
			panic(fmt.Sprintf("unsupported text field operator: %s", c.Operator()))
		}

	} else if fieldType == assets.FieldTypeNumber {
		value, _ := c.ValueAsNumber()
		var query map[string]any

		switch c.Operator() {
		case contactql.OpEqual:
			query = Match("fields.number", value)
		case contactql.OpNotEqual:
			return Not(
				Nested("fields",
					All(fieldQuery, Match("fields.number", value)),
				),
			)
		case contactql.OpGreaterThan:
			query = GreaterThan("fields.number", value)
		case contactql.OpGreaterThanOrEqual:
			query = GreaterThanOrEqual("fields.number", value)
		case contactql.OpLessThan:
			query = LessThan("fields.number", value)
		case contactql.OpLessThanOrEqual:
			query = LessThanOrEqual("fields.number", value)
		default:
			panic(fmt.Sprintf("unsupported number field operator: %s", c.Operator()))
		}

		return Nested("fields", All(fieldQuery, query))

	} else if fieldType == assets.FieldTypeDatetime {
		value, _ := c.ValueAsDate(env)
		start, end := dates.DayToUTCRange(value, value.Location())
		var query map[string]any

		switch c.Operator() {
		case contactql.OpEqual:
			query = Between("fields.datetime", start, end)
		case contactql.OpNotEqual:
			return Not(
				Nested("fields",
					All(fieldQuery, Between("fields.datetime", start, end)),
				),
			)
		case contactql.OpGreaterThan:
			query = GreaterThanOrEqual("fields.datetime", end)
		case contactql.OpGreaterThanOrEqual:
			query = GreaterThanOrEqual("fields.datetime", start)
		case contactql.OpLessThan:
			query = LessThan("fields.datetime", start)
		case contactql.OpLessThanOrEqual:
			query = LessThan("fields.datetime", end)
		default:
			panic(fmt.Sprintf("unsupported datetime field operator: %s", c.Operator()))
		}

		return Nested("fields", All(fieldQuery, query))

	} else if fieldType == assets.FieldTypeState || fieldType == assets.FieldTypeDistrict || fieldType == assets.FieldTypeWard {
		value := strings.ToLower(c.Value())
		name := fmt.Sprintf("fields.%s_keyword", string(fieldType))

		switch c.Operator() {
		case contactql.OpEqual:
			return Nested("fields", All(fieldQuery, Term(name, value)))
		case contactql.OpNotEqual:
			return Not(
				Nested("fields",
					All(Term(name, value), Exists(name)),
				),
			)
		default:
			panic(fmt.Sprintf("unsupported location field operator: %s", c.Operator()))
		}
	}

	panic(fmt.Sprintf("unsupported field type: %s", fieldType))
}

func attributeCondition(env envs.Environment, resolver contactql.Resolver, mapper AssetMapper, c *contactql.Condition) map[string]any {
	key := c.PropertyKey()
	value := strings.ToLower(c.Value())

	// special case for set/unset for name and language
	if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && value == "" &&
		(key == contactql.AttributeName || key == contactql.AttributeLanguage) {

		query := All(Exists(key), Not(Term(fmt.Sprintf("%s.keyword", key), "")))

		if c.Operator() == contactql.OpEqual {
			query = Not(query)
		}

		return query
	}

	switch c.PropertyKey() {
	case contactql.AttributeUUID:
		return textAttributeQuery(c, "uuid", strings.ToLower)
	case contactql.AttributeID:
		switch c.Operator() {
		case contactql.OpEqual:
			return Ids(value)
		case contactql.OpNotEqual:
			return Not(Ids(value))
		default:
			panic(fmt.Sprintf("unsupported ID attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeName:
		switch c.Operator() {
		case contactql.OpEqual:
			return Term("name.keyword", c.Value())
		case contactql.OpNotEqual:
			return Not(Term("name.keyword", c.Value()))
		case contactql.OpContains:
			return Match("name", value)
		default:
			panic(fmt.Sprintf("unsupported name attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeStatus:
		return textAttributeQuery(c, "status", func(v string) string {
			return contactStatusCodes[strings.ToLower(v)]
		})
	case contactql.AttributeLanguage:
		return textAttributeQuery(c, "language", strings.ToLower)
	case contactql.AttributeCreatedOn:
		value, _ := c.ValueAsDate(env)
		start, end := dates.DayToUTCRange(value, value.Location())

		switch c.Operator() {
		case contactql.OpEqual:
			return Between("created_on", start, end)
		case contactql.OpNotEqual:
			return Not(Between("created_on", start, end))
		case contactql.OpGreaterThan:
			return GreaterThanOrEqual("created_on", end)
		case contactql.OpGreaterThanOrEqual:
			return GreaterThanOrEqual("created_on", start)
		case contactql.OpLessThan:
			return LessThan("created_on", start)
		case contactql.OpLessThanOrEqual:
			return LessThan("created_on", end)
		default:
			panic(fmt.Sprintf("unsupported created_on attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeLastSeenOn:
		// special case for set/unset
		if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && value == "" {
			query := Exists("last_seen_on")
			if c.Operator() == contactql.OpEqual {
				query = Not(query)
			}
			return query
		}

		value, _ := c.ValueAsDate(env)
		start, end := dates.DayToUTCRange(value, value.Location())

		switch c.Operator() {
		case contactql.OpEqual:
			return Between("last_seen_on", start, end)
		case contactql.OpNotEqual:
			return Not(Between("last_seen_on", start, end))
		case contactql.OpGreaterThan:
			return GreaterThanOrEqual("last_seen_on", end)
		case contactql.OpGreaterThanOrEqual:
			return GreaterThanOrEqual("last_seen_on", start)
		case contactql.OpLessThan:
			return LessThan("last_seen_on", start)
		case contactql.OpLessThanOrEqual:
			return LessThan("last_seen_on", end)
		default:
			panic(fmt.Sprintf("unsupported last_seen_on attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeURN:
		value := strings.ToLower(c.Value())

		// special case for set/unset
		if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && value == "" {
			query := Nested("urns", Exists("urns.path"))
			if c.Operator() == contactql.OpEqual {
				query = Not(query)
			}
			return query
		}

		switch c.Operator() {
		case contactql.OpEqual:
			return Nested("urns", Term("urns.path.keyword", value))
		case contactql.OpNotEqual:
			return Not(Nested("urns", Term("urns.path.keyword", value)))
		case contactql.OpContains:
			return Nested("urns", MatchPhrase("urns.path", value))
		default:
			panic(fmt.Sprintf("unsupported URN attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeGroup:
		// special case for set/unset
		if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && value == "" {
			query := Exists("group_ids")
			if c.Operator() == contactql.OpEqual {
				query = Not(query)
			}
			return query
		}

		group := c.ValueAsGroup(resolver)

		switch c.Operator() {
		case contactql.OpEqual:
			return Term("group_ids", mapper.Group(group))
		case contactql.OpNotEqual:
			return Not(Term("group_ids", mapper.Group(group)))
		default:
			panic(fmt.Sprintf("unsupported group attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeFlow, contactql.AttributeHistory:
		fieldName := "flow_id"
		if c.PropertyKey() == contactql.AttributeHistory {
			fieldName = "flow_history_ids"
		}

		// special case for set/unset
		if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && value == "" {
			query := Exists(fieldName)
			if c.Operator() == contactql.OpEqual {
				query = Not(query)
			}
			return query
		}

		flow := c.ValueAsFlow(resolver)

		switch c.Operator() {
		case contactql.OpEqual:
			return Term(fieldName, mapper.Flow(flow))
		case contactql.OpNotEqual:
			return Not(Term(fieldName, mapper.Flow(flow)))
		default:
			panic(fmt.Sprintf("unsupported flow attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeTickets:
		return numericalAttributeQuery(c, "tickets")
	default:
		panic(fmt.Sprintf("unsupported contact attribute: %s", key))
	}
}

func schemeCondition(c *contactql.Condition) map[string]any {
	key := c.PropertyKey()
	value := strings.ToLower(c.Value())

	// special case for set/unset
	if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && value == "" {
		query := Nested("urns", All(Term("urns.scheme", key), Exists("urns.path")))
		if c.Operator() == contactql.OpEqual {
			query = Not(query)
		}
		return query
	}

	switch c.Operator() {
	case contactql.OpEqual:
		return Nested("urns", All(Term("urns.path.keyword", value), Term("urns.scheme", key)))
	case contactql.OpNotEqual:
		return Not(Nested("urns", All(Term("urns.path.keyword", value), Term("urns.scheme", key))))
	case contactql.OpContains:
		return Nested("urns", All(MatchPhrase("urns.path", value), Term("urns.scheme", key)))
	default:
		panic(fmt.Sprintf("unsupported scheme operator: %s", c.Operator()))
	}
}

func textAttributeQuery(c *contactql.Condition, name string, tx func(string) string) map[string]any {
	value := tx(c.Value())

	switch c.Operator() {
	case contactql.OpEqual:
		return Term(name, value)
	case contactql.OpNotEqual:
		return Not(Term(name, value))
	default:
		panic(fmt.Sprintf("unsupported %s attribute operator: %s", name, c.Operator()))
	}
}

func numericalAttributeQuery(c *contactql.Condition, name string) map[string]any {
	value, _ := c.ValueAsNumber()

	switch c.Operator() {
	case contactql.OpEqual:
		return Match(name, value)
	case contactql.OpNotEqual:
		return Not(Match(name, value))
	case contactql.OpGreaterThan:
		return GreaterThan(name, value)
	case contactql.OpGreaterThanOrEqual:
		return GreaterThanOrEqual(name, value)
	case contactql.OpLessThan:
		return LessThan(name, value)
	case contactql.OpLessThanOrEqual:
		return LessThanOrEqual(name, value)
	default:
		panic(fmt.Sprintf("unsupported %s attribute operator: %s", name, c.Operator()))
	}
}
