package contactql

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/nyaruka/gocommon/urns"
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

var attributeKeys = map[string]bool{
	"id":         true,
	"name":       true,
	"language":   true,
	"created_on": true,
}

type Visitor struct {
	gen.BaseContactQLVisitor

	redaction envs.RedactionPolicy

	errors []error
}

// NewVisitor creates a new ContactQL visitor
func NewVisitor(redaction envs.RedactionPolicy) *Visitor {
	return &Visitor{redaction: redaction}
}

// Visit the top level parse tree
func (v *Visitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

// parse: expression
func (v *Visitor) VisitParse(ctx *gen.ParseContext) interface{} {
	return v.Visit(ctx.Expression())
}

// expression : TEXT
func (v *Visitor) VisitImplicitCondition(ctx *gen.ImplicitConditionContext) interface{} {
	value := ctx.TEXT().GetText()

	if v.redaction == envs.RedactionPolicyURNs {
		num, err := strconv.Atoi(value)
		if err == nil {
			return newCondition(PropertyTypeAttribute, "id", "=", strconv.Itoa(num))
		}
	} else if telRegex.MatchString(value) {
		value = cleanSpecialCharsRegex.ReplaceAllString(value, "")

		return newCondition(PropertyTypeScheme, urns.TelScheme, "~", value)
	}

	return newCondition(PropertyTypeAttribute, "name", "~", value)
}

// expression : TEXT COMPARATOR literal
func (v *Visitor) VisitCondition(ctx *gen.ConditionContext) interface{} {
	propKey := strings.ToLower(ctx.TEXT().GetText())
	comparator := strings.ToLower(ctx.COMPARATOR().GetText())
	value := v.Visit(ctx.Literal()).(string)

	resolvedAlias, isAlias := comparatorAliases[comparator]
	if isAlias {
		comparator = resolvedAlias
	}

	var propType PropertyType

	if attributeKeys[propKey] {
		propType = PropertyTypeAttribute
	} else if urns.IsValidScheme(propKey) {
		propType = PropertyTypeScheme

		if v.redaction == envs.RedactionPolicyURNs {
			v.errors = append(v.errors, errors.New("URN scheme not allowed"))
		}
	} else {
		propType = PropertyTypeField
	}

	return newCondition(propType, propKey, comparator, value)
}

// expression : expression AND expression
func (v *Visitor) VisitCombinationAnd(ctx *gen.CombinationAndContext) interface{} {
	child1 := v.Visit(ctx.Expression(0)).(QueryNode)
	child2 := v.Visit(ctx.Expression(1)).(QueryNode)
	return NewBoolCombination(BoolOperatorAnd, child1, child2)
}

// expression : expression expression
func (v *Visitor) VisitCombinationImpicitAnd(ctx *gen.CombinationImpicitAndContext) interface{} {
	child1 := v.Visit(ctx.Expression(0)).(QueryNode)
	child2 := v.Visit(ctx.Expression(1)).(QueryNode)
	return NewBoolCombination(BoolOperatorAnd, child1, child2)
}

// expression : expression OR expression
func (v *Visitor) VisitCombinationOr(ctx *gen.CombinationOrContext) interface{} {
	child1 := v.Visit(ctx.Expression(0)).(QueryNode)
	child2 := v.Visit(ctx.Expression(1)).(QueryNode)
	return NewBoolCombination(BoolOperatorOr, child1, child2)
}

// expression : LPAREN expression RPAREN
func (v *Visitor) VisitExpressionGrouping(ctx *gen.ExpressionGroupingContext) interface{} {
	return v.Visit(ctx.Expression())
}

// literal : TEXT
func (v *Visitor) VisitTextLiteral(ctx *gen.TextLiteralContext) interface{} {
	return ctx.GetText()
}

// literal : STRING
func (v *Visitor) VisitStringLiteral(ctx *gen.StringLiteralContext) interface{} {
	value := ctx.GetText()
	value = value[1 : len(value)-1]
	return strings.Replace(value, `""`, `"`, -1) // unescape embedded quotes
}
