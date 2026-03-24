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

// Converter converts contactql queries and sorts to Elasticsearch queries and sorts
type Converter struct {
	env      envs.Environment
	assets   AssetMapper
	resolver contactql.Resolver // set from the parsed query during Query()
}

// NewConverter creates a new Converter
func NewConverter(env envs.Environment, assets AssetMapper) *Converter {
	return &Converter{env: env, assets: assets}
}

// we store contact status in elastic as single char codes
var contactStatusCodes = map[string]string{
	"active":   "A",
	"blocked":  "B",
	"stopped":  "S",
	"archived": "V",
}

// Query converts a contactql query to an Elastic query
func (c *Converter) Query(query *contactql.ContactQuery) elastic.Query {
	if query.Resolver() == nil {
		panic("can only convert queries parsed with a resolver")
	}

	c.resolver = query.Resolver()

	return c.nodeToElastic(query.Root())
}

func (c *Converter) nodeToElastic(node contactql.QueryNode) elastic.Query {
	switch n := node.(type) {
	case *contactql.BoolCombination:
		return c.boolCombination(n)
	case *contactql.Condition:
		return c.condition(n)
	default:
		panic(fmt.Sprintf("unsupported node type: %T", n))
	}
}

func (c *Converter) boolCombination(combination *contactql.BoolCombination) elastic.Query {
	queries := make([]elastic.Query, len(combination.Children()))
	for i, child := range combination.Children() {
		queries[i] = c.nodeToElastic(child)
	}

	if combination.Operator() == contactql.BoolOperatorAnd {
		return elastic.All(queries...)
	}

	return elastic.Any(queries...)
}

func (c *Converter) condition(cond *contactql.Condition) elastic.Query {
	switch cond.PropertyType() {
	case contactql.PropertyTypeField:
		return c.fieldCondition(cond)
	case contactql.PropertyTypeAttribute:
		return c.attributeCondition(cond)
	case contactql.PropertyTypeURN:
		return c.schemeCondition(cond)
	default:
		panic(fmt.Sprintf("unsupported property type: %s", cond.PropertyType()))
	}
}

func (c *Converter) fieldCondition(cond *contactql.Condition) elastic.Query {
	field := c.resolver.ResolveField(cond.PropertyKey())
	fieldType := field.Type()
	fieldQuery := elastic.Term("fields.field", field.UUID())

	// special cases for set/unset
	if (cond.Operator() == contactql.OpEqual || cond.Operator() == contactql.OpNotEqual) && cond.Value() == "" {
		query := elastic.Nested("fields", elastic.All(fieldQuery, elastic.Exists("fields."+string(fieldType))))

		// if we are looking for unset, inverse our query
		if cond.Operator() == contactql.OpEqual {
			query = elastic.Not(query)
		}
		return query
	}

	if fieldType == assets.FieldTypeText {
		value := strings.ToLower(cond.Value())

		switch cond.Operator() {
		case contactql.OpEqual:
			return elastic.Nested("fields", elastic.All(fieldQuery, elastic.Term("fields.text", value)))
		case contactql.OpNotEqual:
			query := elastic.All(fieldQuery, elastic.Term("fields.text", value), elastic.Exists("fields.text"))
			return elastic.Not(elastic.Nested("fields", query))
		default:
			panic(fmt.Sprintf("unsupported text field operator: %s", cond.Operator()))
		}

	} else if fieldType == assets.FieldTypeNumber {
		value, _ := cond.ValueAsNumber()
		var query elastic.Query

		switch cond.Operator() {
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
			panic(fmt.Sprintf("unsupported number field operator: %s", cond.Operator()))
		}

		return elastic.Nested("fields", elastic.All(fieldQuery, query))

	} else if fieldType == assets.FieldTypeDatetime {
		value, _ := cond.ValueAsDate(c.env)
		start, end := dates.DayToUTCRange(value, value.Location())
		var query elastic.Query

		switch cond.Operator() {
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
			panic(fmt.Sprintf("unsupported datetime field operator: %s", cond.Operator()))
		}

		return elastic.Nested("fields", elastic.All(fieldQuery, query))

	} else if fieldType == assets.FieldTypeState || fieldType == assets.FieldTypeDistrict || fieldType == assets.FieldTypeWard {
		value := strings.ToLower(cond.Value())
		name := fmt.Sprintf("fields.%s_keyword", string(fieldType))

		switch cond.Operator() {
		case contactql.OpEqual:
			return elastic.Nested("fields", elastic.All(fieldQuery, elastic.Term(name, value)))
		case contactql.OpNotEqual:
			return elastic.Not(
				elastic.Nested("fields",
					elastic.All(elastic.Term(name, value), elastic.Exists(name)),
				),
			)
		default:
			panic(fmt.Sprintf("unsupported location field operator: %s", cond.Operator()))
		}
	}

	panic(fmt.Sprintf("unsupported field type: %s", fieldType))
}

func (c *Converter) attributeCondition(cond *contactql.Condition) elastic.Query {
	key := cond.PropertyKey()
	value := strings.ToLower(cond.Value())

	// special case for set/unset for name and language
	if (cond.Operator() == contactql.OpEqual || cond.Operator() == contactql.OpNotEqual) && value == "" &&
		(key == contactql.AttributeName || key == contactql.AttributeLanguage) {

		query := elastic.All(elastic.Exists(key), elastic.Not(elastic.Term(fmt.Sprintf("%s.keyword", key), "")))

		if cond.Operator() == contactql.OpEqual {
			query = elastic.Not(query)
		}

		return query
	}

	switch cond.PropertyKey() {
	case contactql.AttributeUUID:
		return textAttributeQuery(cond, "uuid", strings.ToLower)
	case contactql.AttributeID:
		switch cond.Operator() {
		case contactql.OpEqual:
			return elastic.Ids(value)
		case contactql.OpNotEqual:
			return elastic.Not(elastic.Ids(value))
		default:
			panic(fmt.Sprintf("unsupported ID attribute operator: %s", cond.Operator()))
		}
	case contactql.AttributeRef:
		value, _ := obfuscate.DecodeID(cond.Value(), c.env.ObfuscationKey()) // if can't be decoded value will be zero which is fine and just means no match

		switch cond.Operator() {
		case contactql.OpEqual:
			return elastic.Ids(fmt.Sprint(value))
		case contactql.OpNotEqual:
			return elastic.Not(elastic.Ids(fmt.Sprint(value)))
		default:
			panic(fmt.Sprintf("unsupported ref attribute operator: %s", cond.Operator()))
		}
	case contactql.AttributeName:
		switch cond.Operator() {
		case contactql.OpEqual:
			return elastic.Term("name.keyword", cond.Value())
		case contactql.OpNotEqual:
			return elastic.Not(elastic.Term("name.keyword", cond.Value()))
		case contactql.OpContains:
			return elastic.Match("name", value)
		default:
			panic(fmt.Sprintf("unsupported name attribute operator: %s", cond.Operator()))
		}
	case contactql.AttributeStatus:
		return textAttributeQuery(cond, "status", func(v string) string {
			return contactStatusCodes[strings.ToLower(v)]
		})
	case contactql.AttributeLanguage:
		return textAttributeQuery(cond, "language", strings.ToLower)
	case contactql.AttributeCreatedOn:
		value, _ := cond.ValueAsDate(c.env)
		start, end := dates.DayToUTCRange(value, value.Location())

		switch cond.Operator() {
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
			panic(fmt.Sprintf("unsupported created_on attribute operator: %s", cond.Operator()))
		}
	case contactql.AttributeLastSeenOn:
		// special case for set/unset
		if (cond.Operator() == contactql.OpEqual || cond.Operator() == contactql.OpNotEqual) && value == "" {
			query := elastic.Exists("last_seen_on")
			if cond.Operator() == contactql.OpEqual {
				query = elastic.Not(query)
			}
			return query
		}

		value, _ := cond.ValueAsDate(c.env)
		start, end := dates.DayToUTCRange(value, value.Location())

		switch cond.Operator() {
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
			panic(fmt.Sprintf("unsupported last_seen_on attribute operator: %s", cond.Operator()))
		}
	case contactql.AttributeURN:
		value := strings.ToLower(cond.Value())

		// special case for set/unset
		if (cond.Operator() == contactql.OpEqual || cond.Operator() == contactql.OpNotEqual) && value == "" {
			query := elastic.Nested("urns", elastic.Exists("urns.path"))
			if cond.Operator() == contactql.OpEqual {
				query = elastic.Not(query)
			}
			return query
		}

		switch cond.Operator() {
		case contactql.OpEqual:
			return elastic.Nested("urns", elastic.Term("urns.path.keyword", value))
		case contactql.OpNotEqual:
			return elastic.Not(elastic.Nested("urns", elastic.Term("urns.path.keyword", value)))
		case contactql.OpContains:
			return elastic.Nested("urns", elastic.MatchPhrase("urns.path", value))
		default:
			panic(fmt.Sprintf("unsupported URN attribute operator: %s", cond.Operator()))
		}
	case contactql.AttributeGroup:
		// special case for set/unset
		if (cond.Operator() == contactql.OpEqual || cond.Operator() == contactql.OpNotEqual) && value == "" {
			query := elastic.Exists("group_ids")
			if cond.Operator() == contactql.OpEqual {
				query = elastic.Not(query)
			}
			return query
		}

		group := cond.ValueAsGroup(c.resolver)

		switch cond.Operator() {
		case contactql.OpEqual:
			return elastic.Term("group_ids", c.assets.Group(group))
		case contactql.OpNotEqual:
			return elastic.Not(elastic.Term("group_ids", c.assets.Group(group)))
		default:
			panic(fmt.Sprintf("unsupported group attribute operator: %s", cond.Operator()))
		}
	case contactql.AttributeFlow, contactql.AttributeHistory:
		fieldName := "flow_id"
		if cond.PropertyKey() == contactql.AttributeHistory {
			fieldName = "flow_history_ids"
		}

		// special case for set/unset
		if (cond.Operator() == contactql.OpEqual || cond.Operator() == contactql.OpNotEqual) && value == "" {
			query := elastic.Exists(fieldName)
			if cond.Operator() == contactql.OpEqual {
				query = elastic.Not(query)
			}
			return query
		}

		flow := cond.ValueAsFlow(c.resolver)

		switch cond.Operator() {
		case contactql.OpEqual:
			return elastic.Term(fieldName, c.assets.Flow(flow))
		case contactql.OpNotEqual:
			return elastic.Not(elastic.Term(fieldName, c.assets.Flow(flow)))
		default:
			panic(fmt.Sprintf("unsupported flow attribute operator: %s", cond.Operator()))
		}
	case contactql.AttributeTickets:
		return numericalAttributeQuery(cond, "tickets")
	default:
		panic(fmt.Sprintf("unsupported contact attribute: %s", key))
	}
}

func (c *Converter) schemeCondition(cond *contactql.Condition) elastic.Query {
	key := cond.PropertyKey()
	value := strings.ToLower(cond.Value())

	// special case for set/unset
	if (cond.Operator() == contactql.OpEqual || cond.Operator() == contactql.OpNotEqual) && value == "" {
		query := elastic.Nested("urns", elastic.All(elastic.Term("urns.scheme", key), elastic.Exists("urns.path")))
		if cond.Operator() == contactql.OpEqual {
			query = elastic.Not(query)
		}
		return query
	}

	switch cond.Operator() {
	case contactql.OpEqual:
		return elastic.Nested("urns", elastic.All(elastic.Term("urns.path.keyword", value), elastic.Term("urns.scheme", key)))
	case contactql.OpNotEqual:
		return elastic.Not(elastic.Nested("urns", elastic.All(elastic.Term("urns.path.keyword", value), elastic.Term("urns.scheme", key))))
	case contactql.OpContains:
		return elastic.Nested("urns", elastic.All(elastic.MatchPhrase("urns.path", value), elastic.Term("urns.scheme", key)))
	default:
		panic(fmt.Sprintf("unsupported scheme operator: %s", cond.Operator()))
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
