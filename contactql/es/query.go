package es

import (
	"fmt"
	"strings"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/elastic"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils/obfuscate"
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
func ToElasticQuery(env envs.Environment, mapper AssetMapper, query *contactql.ContactQuery) elastic.Query {
	if query.Resolver() == nil {
		panic("can only convert queries parsed with a resolver")
	}

	return nodeToElastic(env, query.Resolver(), mapper, query.Root())
}

func nodeToElastic(env envs.Environment, resolver contactql.Resolver, mapper AssetMapper, node contactql.QueryNode) elastic.Query {
	switch n := node.(type) {
	case *contactql.BoolCombination:
		return boolCombination(env, resolver, mapper, n)
	case *contactql.Condition:
		return condition(env, resolver, mapper, n)
	default:
		panic(fmt.Sprintf("unsupported node type: %T", n))
	}
}

func boolCombination(env envs.Environment, resolver contactql.Resolver, mapper AssetMapper, combination *contactql.BoolCombination) elastic.Query {
	queries := make([]elastic.Query, len(combination.Children()))
	for i, child := range combination.Children() {
		queries[i] = nodeToElastic(env, resolver, mapper, child)
	}

	if combination.Operator() == contactql.BoolOperatorAnd {
		return elastic.All(queries...)
	}

	return elastic.Any(queries...)
}

func condition(env envs.Environment, resolver contactql.Resolver, mapper AssetMapper, c *contactql.Condition) elastic.Query {
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

func fieldCondition(env envs.Environment, resolver contactql.Resolver, c *contactql.Condition) elastic.Query {
	field := resolver.ResolveField(c.PropertyKey())
	fieldType := field.Type()
	fieldQuery := elastic.Term("fields.field", field.UUID())

	// special cases for set/unset
	if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && c.Value() == "" {
		query := elastic.Nested("fields", elastic.All(fieldQuery, elastic.Exists("fields."+string(fieldType))))

		// if we are looking for unset, inverse our query
		if c.Operator() == contactql.OpEqual {
			query = elastic.Not(query)
		}
		return query
	}

	if fieldType == assets.FieldTypeText {
		value := strings.ToLower(c.Value())

		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.Nested("fields", elastic.All(fieldQuery, elastic.Term("fields.text", value)))
		case contactql.OpNotEqual:
			query := elastic.All(fieldQuery, elastic.Term("fields.text", value), elastic.Exists("fields.text"))
			return elastic.Not(elastic.Nested("fields", query))
		default:
			panic(fmt.Sprintf("unsupported text field operator: %s", c.Operator()))
		}

	} else if fieldType == assets.FieldTypeNumber {
		value, _ := c.ValueAsNumber()
		var query elastic.Query

		switch c.Operator() {
		case contactql.OpEqual:
			query = elastic.Match("fields.number", value)
		case contactql.OpNotEqual:
			return elastic.Not(
				elastic.Nested("fields",
					elastic.All(fieldQuery, elastic.Match("fields.number", value)),
				),
			)
		case contactql.OpGreaterThan:
			query = elastic.GreaterThan("fields.number", value)
		case contactql.OpGreaterThanOrEqual:
			query = elastic.GreaterThanOrEqual("fields.number", value)
		case contactql.OpLessThan:
			query = elastic.LessThan("fields.number", value)
		case contactql.OpLessThanOrEqual:
			query = elastic.LessThanOrEqual("fields.number", value)
		default:
			panic(fmt.Sprintf("unsupported number field operator: %s", c.Operator()))
		}

		return elastic.Nested("fields", elastic.All(fieldQuery, query))

	} else if fieldType == assets.FieldTypeDatetime {
		value, _ := c.ValueAsDate(env)
		start, end := dates.DayToUTCRange(value, value.Location())
		var query elastic.Query

		switch c.Operator() {
		case contactql.OpEqual:
			query = elastic.Between("fields.datetime", start, end)
		case contactql.OpNotEqual:
			return elastic.Not(
				elastic.Nested("fields",
					elastic.All(fieldQuery, elastic.Between("fields.datetime", start, end)),
				),
			)
		case contactql.OpGreaterThan:
			query = elastic.GreaterThanOrEqual("fields.datetime", end)
		case contactql.OpGreaterThanOrEqual:
			query = elastic.GreaterThanOrEqual("fields.datetime", start)
		case contactql.OpLessThan:
			query = elastic.LessThan("fields.datetime", start)
		case contactql.OpLessThanOrEqual:
			query = elastic.LessThan("fields.datetime", end)
		default:
			panic(fmt.Sprintf("unsupported datetime field operator: %s", c.Operator()))
		}

		return elastic.Nested("fields", elastic.All(fieldQuery, query))

	} else if fieldType == assets.FieldTypeState || fieldType == assets.FieldTypeDistrict || fieldType == assets.FieldTypeWard {
		value := strings.ToLower(c.Value())
		name := fmt.Sprintf("fields.%s_keyword", string(fieldType))

		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.Nested("fields", elastic.All(fieldQuery, elastic.Term(name, value)))
		case contactql.OpNotEqual:
			return elastic.Not(
				elastic.Nested("fields",
					elastic.All(elastic.Term(name, value), elastic.Exists(name)),
				),
			)
		default:
			panic(fmt.Sprintf("unsupported location field operator: %s", c.Operator()))
		}
	}

	panic(fmt.Sprintf("unsupported field type: %s", fieldType))
}

func attributeCondition(env envs.Environment, resolver contactql.Resolver, mapper AssetMapper, c *contactql.Condition) elastic.Query {
	key := c.PropertyKey()
	value := strings.ToLower(c.Value())

	// special case for set/unset for name and language
	if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && value == "" &&
		(key == contactql.AttributeName || key == contactql.AttributeLanguage) {

		query := elastic.All(elastic.Exists(key), elastic.Not(elastic.Term(fmt.Sprintf("%s.keyword", key), "")))

		if c.Operator() == contactql.OpEqual {
			query = elastic.Not(query)
		}

		return query
	}

	switch c.PropertyKey() {
	case contactql.AttributeUUID:
		return textAttributeQuery(c, "uuid", strings.ToLower)
	case contactql.AttributeID:
		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.Ids(value)
		case contactql.OpNotEqual:
			return elastic.Not(elastic.Ids(value))
		default:
			panic(fmt.Sprintf("unsupported ID attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeRef:
		value, _ := obfuscate.DecodeID(c.Value(), env.ObfuscationKey()) // if can't be decoded value will be zero which is fine and just means no match

		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.Ids(fmt.Sprint(value))
		case contactql.OpNotEqual:
			return elastic.Not(elastic.Ids(fmt.Sprint(value)))
		default:
			panic(fmt.Sprintf("unsupported ref attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeName:
		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.Term("name.keyword", c.Value())
		case contactql.OpNotEqual:
			return elastic.Not(elastic.Term("name.keyword", c.Value()))
		case contactql.OpContains:
			return elastic.Match("name", value)
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
			return elastic.Between("created_on", start, end)
		case contactql.OpNotEqual:
			return elastic.Not(elastic.Between("created_on", start, end))
		case contactql.OpGreaterThan:
			return elastic.GreaterThanOrEqual("created_on", end)
		case contactql.OpGreaterThanOrEqual:
			return elastic.GreaterThanOrEqual("created_on", start)
		case contactql.OpLessThan:
			return elastic.LessThan("created_on", start)
		case contactql.OpLessThanOrEqual:
			return elastic.LessThan("created_on", end)
		default:
			panic(fmt.Sprintf("unsupported created_on attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeLastSeenOn:
		// special case for set/unset
		if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && value == "" {
			query := elastic.Exists("last_seen_on")
			if c.Operator() == contactql.OpEqual {
				query = elastic.Not(query)
			}
			return query
		}

		value, _ := c.ValueAsDate(env)
		start, end := dates.DayToUTCRange(value, value.Location())

		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.Between("last_seen_on", start, end)
		case contactql.OpNotEqual:
			return elastic.Not(elastic.Between("last_seen_on", start, end))
		case contactql.OpGreaterThan:
			return elastic.GreaterThanOrEqual("last_seen_on", end)
		case contactql.OpGreaterThanOrEqual:
			return elastic.GreaterThanOrEqual("last_seen_on", start)
		case contactql.OpLessThan:
			return elastic.LessThan("last_seen_on", start)
		case contactql.OpLessThanOrEqual:
			return elastic.LessThan("last_seen_on", end)
		default:
			panic(fmt.Sprintf("unsupported last_seen_on attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeURN:
		value := strings.ToLower(c.Value())

		// special case for set/unset
		if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && value == "" {
			query := elastic.Nested("urns", elastic.Exists("urns.path"))
			if c.Operator() == contactql.OpEqual {
				query = elastic.Not(query)
			}
			return query
		}

		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.Nested("urns", elastic.Term("urns.path.keyword", value))
		case contactql.OpNotEqual:
			return elastic.Not(elastic.Nested("urns", elastic.Term("urns.path.keyword", value)))
		case contactql.OpContains:
			return elastic.Nested("urns", elastic.MatchPhrase("urns.path", value))
		default:
			panic(fmt.Sprintf("unsupported URN attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeGroup:
		// special case for set/unset
		if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && value == "" {
			query := elastic.Exists("group_ids")
			if c.Operator() == contactql.OpEqual {
				query = elastic.Not(query)
			}
			return query
		}

		group := c.ValueAsGroup(resolver)

		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.Term("group_ids", mapper.Group(group))
		case contactql.OpNotEqual:
			return elastic.Not(elastic.Term("group_ids", mapper.Group(group)))
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
			query := elastic.Exists(fieldName)
			if c.Operator() == contactql.OpEqual {
				query = elastic.Not(query)
			}
			return query
		}

		flow := c.ValueAsFlow(resolver)

		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.Term(fieldName, mapper.Flow(flow))
		case contactql.OpNotEqual:
			return elastic.Not(elastic.Term(fieldName, mapper.Flow(flow)))
		default:
			panic(fmt.Sprintf("unsupported flow attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeTickets:
		return numericalAttributeQuery(c, "tickets")
	default:
		panic(fmt.Sprintf("unsupported contact attribute: %s", key))
	}
}

func schemeCondition(c *contactql.Condition) elastic.Query {
	key := c.PropertyKey()
	value := strings.ToLower(c.Value())

	// special case for set/unset
	if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && value == "" {
		query := elastic.Nested("urns", elastic.All(elastic.Term("urns.scheme", key), elastic.Exists("urns.path")))
		if c.Operator() == contactql.OpEqual {
			query = elastic.Not(query)
		}
		return query
	}

	switch c.Operator() {
	case contactql.OpEqual:
		return elastic.Nested("urns", elastic.All(elastic.Term("urns.path.keyword", value), elastic.Term("urns.scheme", key)))
	case contactql.OpNotEqual:
		return elastic.Not(elastic.Nested("urns", elastic.All(elastic.Term("urns.path.keyword", value), elastic.Term("urns.scheme", key))))
	case contactql.OpContains:
		return elastic.Nested("urns", elastic.All(elastic.MatchPhrase("urns.path", value), elastic.Term("urns.scheme", key)))
	default:
		panic(fmt.Sprintf("unsupported scheme operator: %s", c.Operator()))
	}
}

func textAttributeQuery(c *contactql.Condition, name string, tx func(string) string) elastic.Query {
	value := tx(c.Value())

	switch c.Operator() {
	case contactql.OpEqual:
		return elastic.Term(name, value)
	case contactql.OpNotEqual:
		return elastic.Not(elastic.Term(name, value))
	default:
		panic(fmt.Sprintf("unsupported %s attribute operator: %s", name, c.Operator()))
	}
}

func numericalAttributeQuery(c *contactql.Condition, name string) elastic.Query {
	value, _ := c.ValueAsNumber()

	switch c.Operator() {
	case contactql.OpEqual:
		return elastic.Match(name, value)
	case contactql.OpNotEqual:
		return elastic.Not(elastic.Match(name, value))
	case contactql.OpGreaterThan:
		return elastic.GreaterThan(name, value)
	case contactql.OpGreaterThanOrEqual:
		return elastic.GreaterThanOrEqual(name, value)
	case contactql.OpLessThan:
		return elastic.LessThan(name, value)
	case contactql.OpLessThanOrEqual:
		return elastic.LessThanOrEqual(name, value)
	default:
		panic(fmt.Sprintf("unsupported %s attribute operator: %s", name, c.Operator()))
	}
}
