package contactql

import (
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
	"github.com/shopspring/decimal"
)

// Queryable is the interface objects must implement queried
type Queryable interface {
	QueryProperty(envs.Environment, string, PropertyType) []interface{}
}

// EvaluateQuery evaluates the given query against the given queryable. That query must have been parsed
// with a resolver to ensure all fields and groups resolve. If not function panics.
func EvaluateQuery(env envs.Environment, query *ContactQuery, queryable Queryable) bool {
	if query.Resolver() == nil {
		panic("can only evaluate queries parsed with a resolver")
	}

	return evaluateNode(env, query.Resolver(), query.Root(), queryable)
}

func evaluateNode(env envs.Environment, resolver Resolver, node QueryNode, queryable Queryable) bool {
	switch n := node.(type) {
	case *BoolCombination:
		return evaluateBoolCombination(env, resolver, n, queryable)
	case *Condition:
		return evaluateCondition(env, resolver, n, queryable)
	default:
		panic(fmt.Sprintf("unsupported node type: %T", n))
	}
}

func evaluateBoolCombination(env envs.Environment, resolver Resolver, b *BoolCombination, queryable Queryable) bool {
	if b.op == BoolOperatorAnd {
		for _, child := range b.children {
			if !evaluateNode(env, resolver, child, queryable) {
				return false
			}
		}
		return true
	}

	for _, child := range b.children {
		if evaluateNode(env, resolver, child, queryable) {
			return true
		}
	}
	return false
}

func evaluateCondition(env envs.Environment, resolver Resolver, c *Condition, queryable Queryable) bool {
	// contacts can return multiple values per key, e.g. multiple phone numbers in a "tel = x" condition
	vals := queryable.QueryProperty(env, c.PropertyKey(), c.PropertyType())

	// is this an existence check?
	if c.value == "" {
		if c.operator == OpEqual {
			return len(vals) == 0 // x = "" is true if x doesn't exist
		} else if c.operator == OpNotEqual {
			return len(vals) > 0 // x != "" is false if x doesn't exist (i.e. true if x does exist)
		}
	}

	// evaluate condition against each resolved value
	anyTrue := false
	allTrue := true
	for _, val := range vals {
		if evaluateConditionWithValue(env, resolver, c, val) {
			anyTrue = true
		} else {
			allTrue = false
		}
	}

	// foo != x is only true if all values of foo are not x
	if c.operator == OpNotEqual {
		return allTrue
	}

	// foo = x is true if any value of foo is x
	return anyTrue
}

func evaluateConditionWithValue(env envs.Environment, resolver Resolver, c *Condition, val interface{}) bool {
	valueType := c.resolveValueType(resolver)

	switch valueType {
	case assets.FieldTypeNumber:
		asNumber, _ := c.ValueAsNumber()
		return numberComparison(val.(decimal.Decimal), c.operator, asNumber)
	case assets.FieldTypeDatetime:
		asDate, _ := c.ValueAsDate(env)
		return dateComparison(val.(time.Time), c.operator, asDate)
	default:
		isName := c.propKey == AttributeName // needs to be handled as special case
		return textComparison(val.(string), c.operator, c.value, isName)
	}
}

func textComparison(objectVal string, op Operator, queryVal string, isName bool) bool {
	objectVal = strings.TrimSpace(strings.ToLower(objectVal))
	queryVal = strings.TrimSpace(strings.ToLower(queryVal))

	switch op {
	case OpEqual:
		return objectVal == queryVal
	case OpNotEqual:
		return objectVal != queryVal
	case OpContains:
		// name is special case
		if isName {
			return tokenizedPrefixMatch(objectVal, queryVal, 8)
		}
		return strings.Contains(objectVal, queryVal)
	default:
		panic(fmt.Sprintf("can't query text fields with %s", op))
	}
}

func numberComparison(objectVal decimal.Decimal, op Operator, queryVal decimal.Decimal) bool {
	switch op {
	case OpEqual:
		return objectVal.Equal(queryVal)
	case OpNotEqual:
		return !objectVal.Equal(queryVal)
	case OpGreaterThan:
		return objectVal.GreaterThan(queryVal)
	case OpGreaterThanOrEqual:
		return objectVal.GreaterThanOrEqual(queryVal)
	case OpLessThan:
		return objectVal.LessThan(queryVal)
	case OpLessThanOrEqual:
		return objectVal.LessThanOrEqual(queryVal)
	default:
		panic(fmt.Sprintf("can't query number fields with %s", op))
	}
}

func dateComparison(objectVal time.Time, op Operator, queryVal time.Time) bool {
	utcDayStart, utcDayEnd := dates.DayToUTCRange(queryVal, queryVal.Location())

	switch op {
	case OpEqual:
		return (objectVal.Equal(utcDayStart) || objectVal.After(utcDayStart)) && objectVal.Before(utcDayEnd)
	case OpNotEqual:
		return !((objectVal.Equal(utcDayStart) || objectVal.After(utcDayStart)) && objectVal.Before(utcDayEnd))
	case OpGreaterThan:
		return objectVal.After(utcDayEnd) || objectVal.Equal(utcDayEnd)
	case OpGreaterThanOrEqual:
		return objectVal.After(utcDayStart) || objectVal.Equal(utcDayStart)
	case OpLessThan:
		return objectVal.Before(utcDayStart)
	case OpLessThanOrEqual:
		return objectVal.Before(utcDayEnd)
	default:
		panic(fmt.Sprintf("can't query date fields with %s", op))
	}
}

// performs a prefix match which should be equivalent to an edge_ngram filter in ES
func tokenizedPrefixMatch(objectVal string, queryVal string, length int) bool {
	objectTokens := tokenizeNameValue(objectVal)
	queryTokens := tokenizeNameValue(queryVal)

	for _, objectToken := range objectTokens {
		for _, queryToken := range queryTokens {
			objectTokenVal := utils.Truncate(objectToken, length)
			queryTokenVal := utils.Truncate(queryToken, length)

			if strings.HasPrefix(objectTokenVal, queryTokenVal) {
				return true
			}
		}
	}
	return false
}
