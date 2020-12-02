package contactql

import (
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

// Queryable is the interface objects must implement queried
type Queryable interface {
	QueryProperty(envs.Environment, string, PropertyType) []interface{}
}

// EvaluateQuery evaluates the given parsed query against a queryable object
func EvaluateQuery(env envs.Environment, query *ContactQuery, queryable Queryable) (bool, error) {
	return query.Evaluate(env, queryable)
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
	}

	panic(fmt.Sprintf("can't query text fields with %s", op))
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
	}

	panic(fmt.Sprintf("can't query number fields with %s", op))
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
	}

	panic(fmt.Sprintf("can't query date fields with %s", op))
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
