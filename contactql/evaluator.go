package contactql

import (
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/dates"

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

func textComparison(objectVal string, comparator Comparator, queryVal string, isName bool) bool {
	objectVal = strings.TrimSpace(strings.ToLower(objectVal))
	queryVal = strings.TrimSpace(strings.ToLower(queryVal))

	switch comparator {
	case ComparatorEqual:
		return objectVal == queryVal
	case ComparatorNotEqual:
		return objectVal != queryVal
	case ComparatorContains:
		// name is special case
		if isName {
			return tokenizedPrefixMatch(objectVal, queryVal, 8)
		}
		return strings.Contains(objectVal, queryVal)
	}

	panic(fmt.Sprintf("can't query text fields with %s", comparator))
}

func numberComparison(objectVal decimal.Decimal, comparator Comparator, queryVal decimal.Decimal) bool {
	switch comparator {
	case ComparatorEqual:
		return objectVal.Equal(queryVal)
	case ComparatorGreaterThan:
		return objectVal.GreaterThan(queryVal)
	case ComparatorGreaterThanOrEqual:
		return objectVal.GreaterThanOrEqual(queryVal)
	case ComparatorLessThan:
		return objectVal.LessThan(queryVal)
	case ComparatorLessThanOrEqual:
		return objectVal.LessThanOrEqual(queryVal)
	}

	panic(fmt.Sprintf("can't query number fields with %s", comparator))
}

func dateComparison(objectVal time.Time, comparator Comparator, queryVal time.Time) bool {
	utcDayStart, utcDayEnd := dates.DayToUTCRange(queryVal, queryVal.Location())

	switch comparator {
	case ComparatorEqual:
		return (objectVal.Equal(utcDayStart) || objectVal.After(utcDayStart)) && objectVal.Before(utcDayEnd)
	case ComparatorGreaterThan:
		return objectVal.After(utcDayEnd) || objectVal.Equal(utcDayEnd)
	case ComparatorGreaterThanOrEqual:
		return objectVal.After(utcDayStart) || objectVal.Equal(utcDayStart)
	case ComparatorLessThan:
		return objectVal.Before(utcDayStart)
	case ComparatorLessThanOrEqual:
		return objectVal.Before(utcDayEnd)
	}

	panic(fmt.Sprintf("can't query date fields with %s", comparator))
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
