package contactql

import (
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/goflow/utils"
	"github.com/shopspring/decimal"
)

const (
	ImplicitKey string = "*"
)

type Queryable interface {
	ResolveQueryKey(string) interface{}
}

func EvaluateQuery(env utils.Environment, query *ContactQuery, queryable Queryable) (bool, error) {
	return query.Evaluate(env, queryable)
}

func icontains(s string, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func implicitComparison(objectVals []string, queryVal string) bool {
	for _, objectVal := range objectVals {
		if icontains(objectVal, queryVal) {
			return true
		}
	}
	return false
}

func stringComparison(objectVal string, comparator string, queryVal string) (bool, error) {
	switch comparator {
	case "=":
		return strings.ToLower(objectVal) == strings.ToLower(queryVal), nil
	case "~":
		return icontains(objectVal, queryVal), nil
	}
	return false, fmt.Errorf("can't query text fields with %s", comparator)
}

func decimalComparison(objectVal decimal.Decimal, comparator string, queryVal decimal.Decimal) (bool, error) {
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
	return false, fmt.Errorf("can't query text fields with %s", comparator)
}

func dateComparison(objectVal time.Time, comparator string, queryVal time.Time) (bool, error) {
	utcDayStart, utcDayEnd := utils.DateToUTCRange(queryVal, queryVal.Location())

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
	return false, fmt.Errorf("can't query location fields with %s", comparator)
}
