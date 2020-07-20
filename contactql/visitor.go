package contactql

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql/gen"
	"github.com/nyaruka/goflow/envs"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// an implicit condition like +123-124-6546 or 1234 will be interpreted as a tel ~ condition
var implicitIsPhoneNumberRegex = regexp.MustCompile(`^\+?[\-\d]{4,}$`)

// used to strip formatting from phone number values
var cleanPhoneNumberRegex = regexp.MustCompile(`[^+\d]+`)

var comparatorAliases = map[string]Comparator{
	"has": ComparatorContains,
	"is":  ComparatorEqual,
}

// Fixed attributes that can be searched
const (
	AttributeUUID      = "uuid"
	AttributeID        = "id"
	AttributeName      = "name"
	AttributeLanguage  = "language"
	AttributeURN       = "urn"
	AttributeGroup     = "group"
	AttributeCreatedOn = "created_on"
)

var attributes = map[string]assets.FieldType{
	AttributeUUID:      assets.FieldTypeText,
	AttributeID:        assets.FieldTypeText,
	AttributeName:      assets.FieldTypeText,
	AttributeLanguage:  assets.FieldTypeText,
	AttributeURN:       assets.FieldTypeText,
	AttributeGroup:     assets.FieldTypeText,
	AttributeCreatedOn: assets.FieldTypeDatetime,
}

// Resolver provides functions for resolving fields and groups referenced in queries
type Resolver interface {
	ResolveField(key string) assets.Field
	ResolveGroup(name string) assets.Group
}

type visitor struct {
	gen.BaseContactQLVisitor

	redaction envs.RedactionPolicy
	resolver  Resolver

	errors []error
}

// creates a new ContactQL visitor
func newVisitor(redaction envs.RedactionPolicy, resolver Resolver) *visitor {
	return &visitor{redaction: redaction, resolver: resolver}
}

// Visit the top level parse tree
func (v *visitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

// parse: expression
func (v *visitor) VisitParse(ctx *gen.ParseContext) interface{} {
	return v.Visit(ctx.Expression())
}

// expression : TEXT
func (v *visitor) VisitImplicitCondition(ctx *gen.ImplicitConditionContext) interface{} {
	value := v.Visit(ctx.Literal()).(string)

	asURN, _ := urns.Parse(value)

	if v.redaction == envs.RedactionPolicyURNs {
		num, err := strconv.Atoi(value)
		if err == nil {
			return newCondition(PropertyTypeAttribute, AttributeID, ComparatorEqual, strconv.Itoa(num), attributes[AttributeID], nil)
		}
	} else if asURN != urns.NilURN {
		scheme, path, _, _ := asURN.ToParts()

		return newCondition(PropertyTypeScheme, scheme, ComparatorEqual, path, assets.FieldTypeText, nil)

	} else if implicitIsPhoneNumberRegex.MatchString(value) {
		value = cleanPhoneNumberRegex.ReplaceAllLiteralString(value, "")

		return newCondition(PropertyTypeScheme, urns.TelScheme, ComparatorContains, value, assets.FieldTypeText, nil)
	}

	// convert to contains condition only if we have the right tokens, otherwise make equals check
	comparator := ComparatorContains
	if len(tokenizeNameValue(value)) == 0 {
		comparator = ComparatorEqual
	}

	condition := newCondition(PropertyTypeAttribute, AttributeName, comparator, value, attributes[AttributeName], nil)

	if err := condition.Validate(v.resolver); err != nil {
		v.addError(err)
	}

	return condition
}

// expression : TEXT COMPARATOR literal
func (v *visitor) VisitCondition(ctx *gen.ConditionContext) interface{} {
	propKey := strings.ToLower(ctx.TEXT().GetText())
	comparatorText := strings.ToLower(ctx.COMPARATOR().GetText())
	value := v.Visit(ctx.Literal()).(string)

	comparator, isAlias := comparatorAliases[comparatorText]
	if !isAlias {
		comparator = Comparator(comparatorText)
	}

	var propType PropertyType
	var reference assets.Reference

	// first try to match a fixed attribute
	valueType, isAttribute := attributes[propKey]
	if isAttribute {
		propType = PropertyTypeAttribute

		if propKey == AttributeURN && v.redaction == envs.RedactionPolicyURNs && value != "" {
			v.addError(NewQueryErrorf("cannot query on redacted URNs"))
		}

	} else if urns.IsValidScheme(propKey) {
		// second try to match a URN scheme
		propType = PropertyTypeScheme
		valueType = assets.FieldTypeText

		if v.redaction == envs.RedactionPolicyURNs && value != "" {
			v.addError(NewQueryErrorf("cannot query on redacted URNs"))
		}
	} else {
		field := v.resolver.ResolveField(propKey)
		if field != nil {
			propType = PropertyTypeField
			valueType = field.Type()
			reference = assets.NewFieldReference(field.Key(), field.Name())
		} else {
			v.addError(NewQueryErrorf("can't resolve '%s' to attribute, scheme or field", propKey))
		}
	}

	condition := newCondition(propType, propKey, comparator, value, valueType, reference)

	if err := condition.Validate(v.resolver); err != nil {
		v.addError(err)
	}

	return condition
}

// expression : expression AND expression
func (v *visitor) VisitCombinationAnd(ctx *gen.CombinationAndContext) interface{} {
	child1 := v.Visit(ctx.Expression(0)).(QueryNode)
	child2 := v.Visit(ctx.Expression(1)).(QueryNode)
	return NewBoolCombination(BoolOperatorAnd, child1, child2)
}

// expression : expression expression
func (v *visitor) VisitCombinationImpicitAnd(ctx *gen.CombinationImpicitAndContext) interface{} {
	child1 := v.Visit(ctx.Expression(0)).(QueryNode)
	child2 := v.Visit(ctx.Expression(1)).(QueryNode)
	return NewBoolCombination(BoolOperatorAnd, child1, child2)
}

// expression : expression OR expression
func (v *visitor) VisitCombinationOr(ctx *gen.CombinationOrContext) interface{} {
	child1 := v.Visit(ctx.Expression(0)).(QueryNode)
	child2 := v.Visit(ctx.Expression(1)).(QueryNode)
	return NewBoolCombination(BoolOperatorOr, child1, child2)
}

// expression : LPAREN expression RPAREN
func (v *visitor) VisitExpressionGrouping(ctx *gen.ExpressionGroupingContext) interface{} {
	return v.Visit(ctx.Expression())
}

// literal : TEXT
func (v *visitor) VisitTextLiteral(ctx *gen.TextLiteralContext) interface{} {
	return ctx.GetText()
}

// literal : STRING
func (v *visitor) VisitStringLiteral(ctx *gen.StringLiteralContext) interface{} {
	value := ctx.GetText()

	// unquote, this takes care of escape sequences as well
	unquoted, err := strconv.Unquote(value)

	// if we had an error, just strip surrounding quotes
	if err != nil {
		unquoted = value[1 : len(value)-1]
	}

	return unquoted
}

func (v *visitor) addError(err error) {
	v.errors = append(v.errors, err)
}
