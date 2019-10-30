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
	"github.com/pkg/errors"
)

var telRegex = regexp.MustCompile(`^[+ \d\-\(\)]+$`)
var cleanSpecialCharsRegex = regexp.MustCompile(`[+ \-\(\)]+`)

var comparatorAliases = map[string]string{
	"has": "~",
	"is":  "=",
}

// Fixed attributes that can be searched
const (
	AttributeID        = "id"
	AttributeName      = "name"
	AttributeLanguage  = "language"
	AttributeCreatedOn = "created_on"
)

var attributes = map[string]assets.FieldType{
	AttributeID:        assets.FieldTypeNumber,
	AttributeName:      assets.FieldTypeText,
	AttributeLanguage:  assets.FieldTypeText,
	AttributeCreatedOn: assets.FieldTypeDatetime,
}

// FieldResolverFunc resolves a query property key to a possible contact field
type FieldResolverFunc func(string) assets.Field

type visitor struct {
	gen.BaseContactQLVisitor

	redaction     envs.RedactionPolicy
	fieldResolver FieldResolverFunc

	errors []error
}

// creates a new ContactQL visitor
func newVisitor(redaction envs.RedactionPolicy, fieldResolver FieldResolverFunc) *visitor {
	return &visitor{redaction: redaction, fieldResolver: fieldResolver}
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
	value := ctx.TEXT().GetText()

	if v.redaction == envs.RedactionPolicyURNs {
		num, err := strconv.Atoi(value)
		if err == nil {
			return newCondition(PropertyTypeAttribute, AttributeID, "=", strconv.Itoa(num), attributes[AttributeID])
		}
	} else if telRegex.MatchString(value) {
		value = cleanSpecialCharsRegex.ReplaceAllString(value, "")

		return newCondition(PropertyTypeScheme, urns.TelScheme, "~", value, assets.FieldTypeText)
	}

	return newCondition(PropertyTypeAttribute, AttributeName, "~", value, attributes[AttributeName])
}

// expression : TEXT COMPARATOR literal
func (v *visitor) VisitCondition(ctx *gen.ConditionContext) interface{} {
	propKey := strings.ToLower(ctx.TEXT().GetText())
	comparator := strings.ToLower(ctx.COMPARATOR().GetText())
	value := v.Visit(ctx.Literal()).(string)

	resolvedAlias, isAlias := comparatorAliases[comparator]
	if isAlias {
		comparator = resolvedAlias
	}

	var propType PropertyType

	// first try to match a fixed attribute
	valueType, isAttribute := attributes[propKey]
	if isAttribute {
		propType = PropertyTypeAttribute

	} else if urns.IsValidScheme(propKey) {
		// second try to match a URN scheme
		propType = PropertyTypeScheme
		valueType = assets.FieldTypeText

		if v.redaction == envs.RedactionPolicyURNs {
			v.errors = append(v.errors, errors.New("cannot query on redacted URNs"))
		}
	} else {
		field := v.fieldResolver(propKey)
		if field != nil {
			propType = PropertyTypeField
			valueType = field.Type()
		} else {
			v.errors = append(v.errors, errors.Errorf("can't resolve '%s' to attribute, scheme or field", propKey))
		}
	}

	return newCondition(propType, propKey, comparator, value, valueType)
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
