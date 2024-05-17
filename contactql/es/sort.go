package es

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql"
	esq "github.com/nyaruka/goflow/utils/elastic"
	"github.com/pkg/errors"
)

// ToElasticFieldSort returns the elastic FieldSort for the passed in sort by string
func ToElasticFieldSort(sortBy string, resolver contactql.Resolver) (map[string]any, error) {
	// default to most recent first by id
	if sortBy == "" {
		return esq.Sort("id", false), nil
	}

	// figure out if we are ascending or descending (default is ascending, can be changed with leading -)
	property := sortBy
	ascending := true
	if strings.HasPrefix(sortBy, "-") {
		ascending = false
		property = sortBy[1:]
	}

	property = strings.ToLower(property)

	// name needs to be sorted by keyword field
	if property == contactql.AttributeName {
		return esq.Sort("name.keyword", ascending), nil
	}

	// other attributes are straight sorts
	if property == contactql.AttributeID || property == contactql.AttributeCreatedOn || property == contactql.AttributeLastSeenOn || property == contactql.AttributeLanguage {
		return esq.Sort(property, ascending), nil
	}

	// we are sorting by a custom field
	field := resolver.ResolveField(property)
	if field == nil {
		return nil, errors.Errorf("no such field with key: %s", property)
	}

	var key string
	switch field.Type() {
	case assets.FieldTypeState, assets.FieldTypeDistrict, assets.FieldTypeWard:
		key = fmt.Sprintf("fields.%s_keyword", field.Type())
	default:
		key = fmt.Sprintf("fields.%s", field.Type())
	}

	return esq.SortNested(key, esq.Term("fields.field", field.UUID()), "fields", ascending), nil
}
