package es

import (
	"fmt"
	"strings"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"

	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
)

// ToElasticQuery converts a contactql query to an Elastic query
func ToElasticQuery(env envs.Environment, resolver contactql.Resolver, query *contactql.ContactQuery) (elastic.Query, error) {
	return nodeToElastic(env, resolver, query.Root())
}

func nodeToElastic(env envs.Environment, resolver contactql.Resolver, node contactql.QueryNode) (elastic.Query, error) {
	switch n := node.(type) {
	case *contactql.BoolCombination:
		return boolCombinationToElastic(env, resolver, n)
	case *contactql.Condition:
		return conditionToElastic(env, resolver, n)
	default:
		panic(fmt.Sprintf("unsupported node type: %T", n))
	}
}

func boolCombinationToElastic(env envs.Environment, resolver contactql.Resolver, combination *contactql.BoolCombination) (elastic.Query, error) {
	var err error
	queries := make([]elastic.Query, len(combination.Children()))
	for i, child := range combination.Children() {
		queries[i], err = nodeToElastic(env, resolver, child)
		if err != nil {
			return nil, err
		}
	}

	if combination.Operator() == contactql.BoolOperatorAnd {
		return elastic.NewBoolQuery().Must(queries...), nil
	}

	return elastic.NewBoolQuery().Should(queries...), nil
}

func conditionToElastic(env envs.Environment, resolver contactql.Resolver, c *contactql.Condition) (elastic.Query, error) {
	switch c.PropertyType() {
	case contactql.PropertyTypeField:
		return fieldConditionToElastic(env, resolver, c)
	case contactql.PropertyTypeAttribute:
		return attributeConditionToElastic(env, resolver, c), nil
	case contactql.PropertyTypeScheme:
		return schemeConditionToElastic(env, c), nil
	default:
		panic(fmt.Sprintf("unsupported property type: %s", c.PropertyType()))
	}
}

func fieldConditionToElastic(env envs.Environment, resolver contactql.Resolver, c *contactql.Condition) (elastic.Query, error) {
	var query elastic.Query

	field := resolver.ResolveField(c.PropertyKey())
	if field == nil {
		return nil, errors.Errorf("no such field '%s'", c.PropertyKey())
	}

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
		return query, nil
	}

	if fieldType == assets.FieldTypeText {
		value := strings.ToLower(c.Value())

		switch c.Operator() {
		case contactql.OpEqual:
			query = elastic.NewTermQuery("fields.text", value)
			return elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(fieldQuery, query)), nil
		case contactql.OpNotEqual:
			query = elastic.NewBoolQuery().Must(
				fieldQuery,
				elastic.NewTermQuery("fields.text", value),
				elastic.NewExistsQuery("fields.text"),
			)
			return not(elastic.NewNestedQuery("fields", query)), nil
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
			), nil
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

		return elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(fieldQuery, query)), nil

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
			), nil
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

		return elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(fieldQuery, query)), nil

	} else if fieldType == assets.FieldTypeState || fieldType == assets.FieldTypeDistrict || fieldType == assets.FieldTypeWard {
		value := strings.ToLower(c.Value())
		name := fmt.Sprintf("fields.%s_keyword", string(fieldType))

		switch c.Operator() {
		case contactql.OpEqual:
			query = elastic.NewTermQuery(name, value)
			return elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(fieldQuery, query)), nil
		case contactql.OpNotEqual:
			return not(
				elastic.NewNestedQuery("fields",
					elastic.NewBoolQuery().Must(
						elastic.NewTermQuery(name, value),
						elastic.NewExistsQuery(name),
					),
				),
			), nil
		default:
			panic(fmt.Sprintf("unsupported location field operator: %s", c.Operator()))
		}
	}

	panic(fmt.Sprintf("unsupported field type: %s", fieldType))
}

func attributeConditionToElastic(env envs.Environment, resolver contactql.Resolver, c *contactql.Condition) elastic.Query {
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
		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.NewTermQuery("uuid", value)
		case contactql.OpNotEqual:
			return not(elastic.NewTermQuery("uuid", value))
		default:
			panic(fmt.Sprintf("unsupported UUID attribute operator: %s", c.Operator()))
		}
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
	case contactql.AttributeLanguage:
		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.NewTermQuery("language", value)
		case contactql.OpNotEqual:
			return not(elastic.NewTermQuery("language", value))
		default:
			panic(fmt.Sprintf("unsupported language attribute operator: %s", c.Operator()))
		}
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
		group := c.ValueAsGroup(resolver)

		switch c.Operator() {
		case contactql.OpEqual:
			return elastic.NewTermQuery("groups", group.UUID())
		case contactql.OpNotEqual:
			return not(elastic.NewTermQuery("groups", group.UUID()))
		default:
			panic(fmt.Sprintf("unsupported group attribute operator: %s", c.Operator()))
		}
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

// convenience utility to create a not boolean query
func not(queries ...elastic.Query) *elastic.BoolQuery {
	return elastic.NewBoolQuery().MustNot(queries...)
}
