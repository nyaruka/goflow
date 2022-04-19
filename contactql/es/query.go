package es

import (
	"fmt"
	"strings"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"

	"github.com/olivere/elastic/v7"
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
		return boolCombinationToElastic(env, resolver, mapper, n)
	case *contactql.Condition:
		return conditionToElastic(env, resolver, mapper, n)
	default:
		panic(fmt.Sprintf("unsupported node type: %T", n))
	}
}

func boolCombinationToElastic(env envs.Environment, resolver contactql.Resolver, mapper AssetMapper, combination *contactql.BoolCombination) elastic.Query {
	queries := make([]elastic.Query, len(combination.Children()))
	for i, child := range combination.Children() {
		queries[i] = nodeToElastic(env, resolver, mapper, child)
	}

	if combination.Operator() == contactql.BoolOperatorAnd {
		return elastic.NewBoolQuery().Must(queries...)
	}

	return elastic.NewBoolQuery().Should(queries...)
}

func conditionToElastic(env envs.Environment, resolver contactql.Resolver, mapper AssetMapper, c *contactql.Condition) elastic.Query {
	switch c.PropertyType() {
	case contactql.PropertyTypeField:
		return fieldConditionToElastic(env, resolver, c)
	case contactql.PropertyTypeAttribute:
		return attributeConditionToElastic(env, resolver, mapper, c)
	case contactql.PropertyTypeScheme:
		return schemeConditionToElastic(env, c)
	default:
		panic(fmt.Sprintf("unsupported property type: %s", c.PropertyType()))
	}
}

func fieldConditionToElastic(env envs.Environment, resolver contactql.Resolver, c *contactql.Condition) elastic.Query {
	var query elastic.Query

	field := resolver.ResolveField(c.PropertyKey())
	fieldType := field.Type()
	fieldQuery := elastic.NewTermQuery("fields.field", field.UUID())

	// special cases for set/unset
	if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && c.Value() == "" {
		query = elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(
			fieldQuery,
			elastic.NewExistsQuery("fields."+string(fieldType)),
		))

		// if we are looking for unset, inverse our query
		if c.Operator() == contactql.OpEqual {
			query = not(query)
		}
		return query
	}

	if fieldType == assets.FieldTypeText {
		value := strings.ToLower(c.Value())

		switch c.Operator() {
		case contactql.OpEqual:
			query = elastic.NewTermQuery("fields.text", value)
			return elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(fieldQuery, query))
		case contactql.OpNotEqual:
			query = elastic.NewBoolQuery().Must(
				fieldQuery,
				elastic.NewTermQuery("fields.text", value),
				elastic.NewExistsQuery("fields.text"),
			)
			return not(elastic.NewNestedQuery("fields", query))
		default:
			panic(fmt.Sprintf("unsupported text field operator: %s", c.Operator()))
		}

	} else if fieldType == assets.FieldTypeNumber {
		value, _ := c.ValueAsNumber()

		switch c.Operator() {
		case contactql.OpEqual:
			query = elastic.NewMatchQuery("fields.number", value)
		case contactql.OpNotEqual:
			return not(
				elastic.NewNestedQuery("fields",
					elastic.NewBoolQuery().Must(
						fieldQuery,
						elastic.NewMatchQuery("fields.number", value),
					),
				),
			)
		case contactql.OpGreaterThan:
			query = elastic.NewRangeQuery("fields.number").Gt(value)
		case contactql.OpGreaterThanOrEqual:
			query = elastic.NewRangeQuery("fields.number").Gte(value)
		case contactql.OpLessThan:
			query = elastic.NewRangeQuery("fields.number").Lt(value)
		case contactql.OpLessThanOrEqual:
			query = elastic.NewRangeQuery("fields.number").Lte(value)
		default:
			panic(fmt.Sprintf("unsupported number field operator: %s", c.Operator()))
		}

		return elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(fieldQuery, query))

	} else if fieldType == assets.FieldTypeDatetime {
		value, _ := c.ValueAsDate(env)
		start, end := dates.DayToUTCRange(value, value.Location())

		switch c.Operator() {
		case contactql.OpEqual:
			query = elastic.NewRangeQuery("fields.datetime").Gte(start).Lt(end)
		case contactql.OpNotEqual:
			return not(
				elastic.NewNestedQuery("fields",
					elastic.NewBoolQuery().Must(
						fieldQuery,
						elastic.NewRangeQuery("fields.datetime").Gte(start).Lt(end),
					),
				),
			)
		case contactql.OpGreaterThan:
			query = elastic.NewRangeQuery("fields.datetime").Gte(end)
		case contactql.OpGreaterThanOrEqual:
			query = elastic.NewRangeQuery("fields.datetime").Gte(start)
		case contactql.OpLessThan:
			query = elastic.NewRangeQuery("fields.datetime").Lt(start)
		case contactql.OpLessThanOrEqual:
			query = elastic.NewRangeQuery("fields.datetime").Lt(end)
		default:
			panic(fmt.Sprintf("unsupported datetime field operator: %s", c.Operator()))
		}

		return elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(fieldQuery, query))

	} else if fieldType == assets.FieldTypeState || fieldType == assets.FieldTypeDistrict || fieldType == assets.FieldTypeWard {
		value := strings.ToLower(c.Value())
		name := fmt.Sprintf("fields.%s_keyword", string(fieldType))

		switch c.Operator() {
		case contactql.OpEqual:
			query = elastic.NewTermQuery(name, value)
			return elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(fieldQuery, query))
		case contactql.OpNotEqual:
			return not(
				elastic.NewNestedQuery("fields",
					elastic.NewBoolQuery().Must(
						elastic.NewTermQuery(name, value),
						elastic.NewExistsQuery(name),
					),
				),
			)
		default:
			panic(fmt.Sprintf("unsupported location field operator: %s", c.Operator()))
		}
	}

	panic(fmt.Sprintf("unsupported field type: %s", fieldType))
}

func attributeConditionToElastic(env envs.Environment, resolver contactql.Resolver, mapper AssetMapper, c *contactql.Condition) elastic.Query {
	key := c.PropertyKey()
	value := strings.ToLower(c.Value())
	var query elastic.Query

	// special case for set/unset for name and language
	if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && value == "" &&
		(key == contactql.AttributeName || key == contactql.AttributeLanguage) {

		query = elastic.NewBoolQuery().Must(
			elastic.NewExistsQuery(key),
			not(elastic.NewTermQuery(fmt.Sprintf("%s.keyword", key), "")),
		)

		if c.Operator() == contactql.OpEqual {
			query = not(query)
		}

		return query
	}

	switch c.PropertyKey() {
	case contactql.AttributeUUID:
		return textAttributeQuery(c, "uuid", strings.ToLower)
	case contactql.AttributeID:
		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.NewIdsQuery().Ids(value)
		case contactql.OpNotEqual:
			return not(elastic.NewIdsQuery().Ids(value))
		default:
			panic(fmt.Sprintf("unsupported ID attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeName:
		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.NewTermQuery("name.keyword", c.Value())
		case contactql.OpNotEqual:
			return not(elastic.NewTermQuery("name.keyword", c.Value()))
		case contactql.OpContains:
			return elastic.NewMatchQuery("name", value)
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
			return elastic.NewRangeQuery("created_on").Gte(start).Lt(end)
		case contactql.OpNotEqual:
			return not(elastic.NewRangeQuery("created_on").Gte(start).Lt(end))
		case contactql.OpGreaterThan:
			return elastic.NewRangeQuery("created_on").Gte(end)
		case contactql.OpGreaterThanOrEqual:
			return elastic.NewRangeQuery("created_on").Gte(start)
		case contactql.OpLessThan:
			return elastic.NewRangeQuery("created_on").Lt(start)
		case contactql.OpLessThanOrEqual:
			return elastic.NewRangeQuery("created_on").Lt(end)
		default:
			panic(fmt.Sprintf("unsupported created_on attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeLastSeenOn:
		// special case for set/unset
		if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && value == "" {
			query = elastic.NewExistsQuery("last_seen_on")
			if c.Operator() == contactql.OpEqual {
				query = not(query)
			}
			return query
		}

		value, _ := c.ValueAsDate(env)
		start, end := dates.DayToUTCRange(value, value.Location())

		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.NewRangeQuery("last_seen_on").Gte(start).Lt(end)
		case contactql.OpNotEqual:
			return not(elastic.NewRangeQuery("last_seen_on").Gte(start).Lt(end))
		case contactql.OpGreaterThan:
			return elastic.NewRangeQuery("last_seen_on").Gte(end)
		case contactql.OpGreaterThanOrEqual:
			return elastic.NewRangeQuery("last_seen_on").Gte(start)
		case contactql.OpLessThan:
			return elastic.NewRangeQuery("last_seen_on").Lt(start)
		case contactql.OpLessThanOrEqual:
			return elastic.NewRangeQuery("last_seen_on").Lt(end)
		default:
			panic(fmt.Sprintf("unsupported last_seen_on attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeURN:
		value := strings.ToLower(c.Value())

		// special case for set/unset
		if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && value == "" {
			query = elastic.NewNestedQuery("urns", elastic.NewExistsQuery("urns.path"))
			if c.Operator() == contactql.OpEqual {
				query = not(query)
			}
			return query
		}

		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.NewNestedQuery("urns", elastic.NewTermQuery("urns.path.keyword", value))
		case contactql.OpNotEqual:
			return not(elastic.NewNestedQuery("urns", elastic.NewTermQuery("urns.path.keyword", value)))
		case contactql.OpContains:
			return elastic.NewNestedQuery("urns", elastic.NewMatchPhraseQuery("urns.path", value))
		default:
			panic(fmt.Sprintf("unsupported URN attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeGroup:
		// special case for set/unset
		if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && value == "" {
			query = elastic.NewExistsQuery("group_ids")
			if c.Operator() == contactql.OpEqual {
				query = not(query)
			}
			return query
		}

		group := c.ValueAsGroup(resolver)

		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.NewTermQuery("group_ids", mapper.Group(group))
		case contactql.OpNotEqual:
			return not(elastic.NewTermQuery("group_ids", mapper.Group(group)))
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
			query = elastic.NewExistsQuery(fieldName)
			if c.Operator() == contactql.OpEqual {
				query = not(query)
			}
			return query
		}

		flow := c.ValueAsFlow(resolver)

		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.NewTermQuery(fieldName, mapper.Flow(flow))
		case contactql.OpNotEqual:
			return not(elastic.NewTermQuery(fieldName, mapper.Flow(flow)))
		default:
			panic(fmt.Sprintf("unsupported flow attribute operator: %s", c.Operator()))
		}
	case contactql.AttributeTickets:
		return numericalAttributeQuery(c, "tickets")
	default:
		panic(fmt.Sprintf("unsupported contact attribute: %s", key))
	}
}

func schemeConditionToElastic(env envs.Environment, c *contactql.Condition) elastic.Query {
	key := c.PropertyKey()
	value := strings.ToLower(c.Value())

	// special case for set/unset
	if (c.Operator() == contactql.OpEqual || c.Operator() == contactql.OpNotEqual) && value == "" {
		var query elastic.Query
		query = elastic.NewNestedQuery("urns", elastic.NewBoolQuery().Must(
			elastic.NewTermQuery("urns.scheme", key),
			elastic.NewExistsQuery("urns.path"),
		))
		if c.Operator() == contactql.OpEqual {
			query = not(query)
		}
		return query
	}

	switch c.Operator() {
	case contactql.OpEqual:
		return elastic.NewNestedQuery("urns", elastic.NewBoolQuery().Must(
			elastic.NewTermQuery("urns.path.keyword", value),
			elastic.NewTermQuery("urns.scheme", key)),
		)
	case contactql.OpNotEqual:
		return not(elastic.NewNestedQuery("urns", elastic.NewBoolQuery().Must(
			elastic.NewTermQuery("urns.path.keyword", value),
			elastic.NewTermQuery("urns.scheme", key)),
		))
	case contactql.OpContains:
		return elastic.NewNestedQuery("urns", elastic.NewBoolQuery().Must(
			elastic.NewMatchPhraseQuery("urns.path", value),
			elastic.NewTermQuery("urns.scheme", key)),
		)
	default:
		panic(fmt.Sprintf("unsupported scheme operator: %s", c.Operator()))
	}
}

func textAttributeQuery(c *contactql.Condition, name string, tx func(string) string) elastic.Query {
	value := tx(c.Value())

	switch c.Operator() {
	case contactql.OpEqual:
		return elastic.NewTermQuery(name, value)
	case contactql.OpNotEqual:
		return not(elastic.NewTermQuery(name, value))
	default:
		panic(fmt.Sprintf("unsupported %s attribute operator: %s", name, c.Operator()))
	}
}

func numericalAttributeQuery(c *contactql.Condition, name string) elastic.Query {
	value, _ := c.ValueAsNumber()

	switch c.Operator() {
	case contactql.OpEqual:
		return elastic.NewMatchQuery(name, value)
	case contactql.OpNotEqual:
		return not(elastic.NewMatchQuery(name, value))
	case contactql.OpGreaterThan:
		return elastic.NewRangeQuery(name).Gt(value)
	case contactql.OpGreaterThanOrEqual:
		return elastic.NewRangeQuery(name).Gte(value)
	case contactql.OpLessThan:
		return elastic.NewRangeQuery(name).Lt(value)
	case contactql.OpLessThanOrEqual:
		return elastic.NewRangeQuery(name).Lte(value)
	default:
		panic(fmt.Sprintf("unsupported %s attribute operator: %s", name, c.Operator()))
	}
}

// convenience utility to create a not boolean query
func not(queries ...elastic.Query) *elastic.BoolQuery {
	return elastic.NewBoolQuery().MustNot(queries...)
}
