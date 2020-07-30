package es

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql"

	"github.com/olivere/elastic"
)

// ToElasticFieldSort returns the FieldSort for the passed in field
func ToElasticFieldSort(resolver contactql.Resolver, fieldName string) (*elastic.FieldSort, error) {
	// no field name? default to most recent first by id
	if fieldName == "" {
		return elastic.NewFieldSort("id").Desc(), nil
	}

	// figure out if we are ascending or descending (default is ascending, can be changed with leading -)
	ascending := true
	if strings.HasPrefix(fieldName, "-") {
		ascending = false
		fieldName = fieldName[1:]
	}

	fieldName = strings.ToLower(fieldName)

	// name needs to be sorted by keyword field
	if fieldName == contactql.AttributeName {
		return elastic.NewFieldSort("name.keyword").Order(ascending), nil
	}

	// other attributes are straight sorts
	if fieldName == contactql.AttributeID || fieldName == contactql.AttributeCreatedOn || fieldName == contactql.AttributeLanguage {
		return elastic.NewFieldSort(fieldName).Order(ascending), nil
	}

	// we are sorting by a custom field
	field := resolver.ResolveField(fieldName)
	if field == nil {
		return nil, contactql.NewQueryError("", "unable to find field with name: %s", fieldName)
	}

	var key string
	switch field.Type() {
	case assets.FieldTypeState, assets.FieldTypeDistrict, assets.FieldTypeWard:
		key = fmt.Sprintf("fields.%s_keyword", field.Type())
	default:
		key = fmt.Sprintf("fields.%s", field.Type())
	}

	sort := elastic.NewFieldSort(key)
	sort = sort.Nested(elastic.NewNestedSort("fields").Filter(elastic.NewTermQuery("fields.field", field.UUID())))
	sort = sort.Order(ascending)
	return sort, nil
}
