package contactql

import (
	"strings"
	"time"

	"github.com/nyaruka/goflow/dates"
	"github.com/nyaruka/goflow/envs"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const (
	ImplicitKey string = "*"
)

// Queryable is the interface objects must implement queried
type Queryable interface {
	ResolveQueryKey(envs.Environment, string) []interface{}
}

// EvaluateQuery evaluates the given parsed query against a queryable object
func EvaluateQuery(env envs.Environment, query *ContactQuery, queryable Queryable) (bool, error) {
	return query.Evaluate(env, queryable)
}

func icontains(s string, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func textComparison(objectVal string, comparator string, queryVal string) (bool, error) {
	switch comparator {
	case "=":
		return strings.ToLower(objectVal) == strings.ToLower(queryVal), nil
	case "!=":
		return strings.ToLower(objectVal) != strings.ToLower(queryVal), nil
	case "~":
		return icontains(objectVal, queryVal), nil
	}
	return false, errors.Errorf("can't query text fields with %s", comparator)
}

func numberComparison(objectVal decimal.Decimal, comparator string, queryVal decimal.Decimal) (bool, error) {
	switch comparator {
	case "=":
		return objectVal.Equal(queryVal), nil
	case ">":
		return objectVal.GreaterThan(queryVal), nil
	case ">=":
		return objectVal.GreaterThanOrEqual(queryVal), nil
	case "<":
		return objectVal.LessThan(queryVal), nil
	case "<=":
		return objectVal.LessThanOrEqual(queryVal), nil
	}
	return false, errors.Errorf("can't query number fields with %s", comparator)
}

func dateComparison(objectVal time.Time, comparator string, queryVal time.Time) (bool, error) {
	utcDayStart, utcDayEnd := dates.DayToUTCRange(queryVal, queryVal.Location())

	switch comparator {
	case "=":
		return (objectVal.Equal(utcDayStart) || objectVal.After(utcDayStart)) && objectVal.Before(utcDayEnd), nil
	case ">":
		return objectVal.After(utcDayEnd) || objectVal.Equal(utcDayEnd), nil
	case ">=":
		return objectVal.After(utcDayStart) || objectVal.Equal(utcDayStart), nil
	case "<":
		return objectVal.Before(utcDayStart), nil
	case "<=":
		return objectVal.Before(utcDayEnd), nil
	}
	return false, errors.Errorf("can't query datetime fields with %s", comparator)
}
