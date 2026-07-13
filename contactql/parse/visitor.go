package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/nyaruka/gocommon/urns"
	gen "github.com/nyaruka/goflow/antlr/gen/contactql"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"
)

var operatorAliases = map[string]contactql.Operator{
	"has": contactql.OpContains,
	"is":  contactql.OpEqual,
}

type visitor struct {
	gen.BaseContactQLVisitor

	env    envs.Environment
	errors []error
}

// creates a new ContactQL visitor
func newVisitor(env envs.Environment) *visitor {
	return &visitor{env: env}
}

// Visit the top level parse tree
func (v *visitor) Visit(tree antlr.ParseTree) any {
	return tree.Accept(v)
}

// parse: expression
func (v *visitor) VisitParse(ctx *gen.ParseContext) any {
	return v.Visit(ctx.Expression())
}

// expression : TEXT
func (v *visitor) VisitImplicitCondition(ctx *gen.ImplicitConditionContext) any {
	value := v.Visit(ctx.Literal()).(string)

	return contactql.NewImplicitCondition(v.env, value)
}

// expression : PROPERTY COMPARATOR literal
func (v *visitor) VisitCondition(ctx *gen.ConditionContext) any {
	propText := strings.ToLower(ctx.PROPERTY().GetText())
	operatorText := strings.ToLower(ctx.COMPARATOR().GetText())
	value := v.Visit(ctx.Literal()).(string)

	operator, isAlias := operatorAliases[operatorText]
	if !isAlias {
		operator = contactql.Operator(operatorText)
	}

	var propType contactql.PropertyType
	var propKey string

	// check if property type is specified as prefix
	if strings.Contains(propText, ".") {
		parts := strings.SplitN(propText, ".", 2)

		if parts[0] == "fields" {
			propType = contactql.PropertyTypeField
			propKey = parts[1]
		} else if parts[0] == "urns" {
			propType = contactql.PropertyTypeURN
			propKey = parts[1]
		} else {
			v.addError(contactql.NewQueryError(contactql.ErrUnknownPropertyType, fmt.Sprintf("unknown property type '%s'", parts[0])).WithExtra("type", parts[0]))
		}
	} else {
		propKey = propText

		// first try to match a fixed attribute
		_, isAttribute := contactql.Attributes[propKey]
		if isAttribute {
			propType = contactql.PropertyTypeAttribute

			if propKey == contactql.AttributeURN && v.env.RedactionPolicy() == envs.RedactionPolicyURNs && value != "" {
				v.addError(contactql.NewQueryError(contactql.ErrRedactedURNs, "cannot query on redacted URNs"))
			}

		} else if urns.IsValidScheme(propKey) {
			// second try to match a URN scheme
			propType = contactql.PropertyTypeURN

			if v.env.RedactionPolicy() == envs.RedactionPolicyURNs && value != "" {
				v.addError(contactql.NewQueryError(contactql.ErrRedactedURNs, "cannot query on redacted URNs"))
			}
		} else {
			propType = contactql.PropertyTypeField
		}
	}

	return contactql.NewCondition(propType, propKey, operator, value)
}

// expression : expression AND expression
func (v *visitor) VisitCombinationAnd(ctx *gen.CombinationAndContext) any {
	child1 := v.Visit(ctx.Expression(0)).(contactql.QueryNode)
	child2 := v.Visit(ctx.Expression(1)).(contactql.QueryNode)
	return contactql.NewBoolCombination(contactql.BoolOperatorAnd, child1, child2)
}

// expression : expression expression
func (v *visitor) VisitCombinationImpicitAnd(ctx *gen.CombinationImpicitAndContext) any {
	child1 := v.Visit(ctx.Expression(0)).(contactql.QueryNode)
	child2 := v.Visit(ctx.Expression(1)).(contactql.QueryNode)
	return contactql.NewBoolCombination(contactql.BoolOperatorAnd, child1, child2)
}

// expression : expression OR expression
func (v *visitor) VisitCombinationOr(ctx *gen.CombinationOrContext) any {
	child1 := v.Visit(ctx.Expression(0)).(contactql.QueryNode)
	child2 := v.Visit(ctx.Expression(1)).(contactql.QueryNode)
	return contactql.NewBoolCombination(contactql.BoolOperatorOr, child1, child2)
}

// expression : LPAREN expression RPAREN
func (v *visitor) VisitExpressionGrouping(ctx *gen.ExpressionGroupingContext) any {
	return v.Visit(ctx.Expression())
}

// literal : TEXT | NAME
func (v *visitor) VisitTextLiteral(ctx *gen.TextLiteralContext) any {
	return ctx.GetText()
}

// literal : STRING
func (v *visitor) VisitStringLiteral(ctx *gen.StringLiteralContext) any {
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
