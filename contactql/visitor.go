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

var operatorAliases = map[string]Operator{
	"has": OpContains,
	"is":  OpEqual,
}

// Fixed attributes that can be searched
const (
	AttributeUUID       = "uuid"
	AttributeID         = "id"
	AttributeName       = "name"
	AttributeLanguage   = "language"
	AttributeURN        = "urn"
	AttributeGroup      = "group"
	AttributeCreatedOn  = "created_on"
	AttributeLastSeenOn = "last_seen_on"
)

var attributes = map[string]assets.FieldType{
	AttributeUUID:       assets.FieldTypeText,
	AttributeID:         assets.FieldTypeText,
	AttributeName:       assets.FieldTypeText,
	AttributeLanguage:   assets.FieldTypeText,
	AttributeURN:        assets.FieldTypeText,
	AttributeGroup:      assets.FieldTypeText,
	AttributeCreatedOn:  assets.FieldTypeDatetime,
	AttributeLastSeenOn: assets.FieldTypeDatetime,
}

// Resolver provides functions for resolving fields and groups referenced in queries
type Resolver interface {
	ResolveField(key string) assets.Field
	ResolveGroup(name string) assets.Group
}

type visitor struct {
	gen.BaseContactQLVisitor

	env      envs.Environment
	resolver Resolver

	errors []error
}

// creates a new ContactQL visitor
func newVisitor(env envs.Environment, resolver Resolver) *visitor {
	return &visitor{env: env, resolver: resolver}
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

	if v.env.RedactionPolicy() == envs.RedactionPolicyURNs {
		num, err := strconv.Atoi(value)
		if err == nil {
			return newCondition(AttributeID, PropertyTypeAttribute, nil, OpEqual, strconv.Itoa(num), attributes[AttributeID])
		}
	} else if asURN != urns.NilURN {
		scheme, path, _, _ := asURN.ToParts()

		return newCondition(scheme, PropertyTypeScheme, nil, OpEqual, path, assets.FieldTypeText)

	} else if implicitIsPhoneNumberRegex.MatchString(value) {
		value = cleanPhoneNumberRegex.ReplaceAllLiteralString(value, "")

		return newCondition(urns.TelScheme, PropertyTypeScheme, nil, OpContains, value, assets.FieldTypeText)
	}

	// convert to contains condition only if we have the right tokens, otherwise make equals check
	operator := OpContains
	if len(tokenizeNameValue(value)) == 0 {
		operator = OpEqual
	}

	condition := newCondition(AttributeName, PropertyTypeAttribute, nil, operator, value, attributes[AttributeName])

	if err := condition.Validate(v.env, v.resolver); err != nil {
		v.addError(err)
	}

	return condition
}

// expression : TEXT COMPARATOR literal
func (v *visitor) VisitCondition(ctx *gen.ConditionContext) interface{} {
	propKey := strings.ToLower(ctx.TEXT().GetText())
	operatorText := strings.ToLower(ctx.COMPARATOR().GetText())
	value := v.Visit(ctx.Literal()).(string)

	operator, isAlias := operatorAliases[operatorText]
	if !isAlias {
		operator = Operator(operatorText)
	}

	var propType PropertyType
	var propField assets.Field

	// first try to match a fixed attribute
	valueType, isAttribute := attributes[propKey]
	if isAttribute {
		propType = PropertyTypeAttribute

		if propKey == AttributeURN && v.env.RedactionPolicy() == envs.RedactionPolicyURNs && value != "" {
			v.addError(NewQueryError(ErrRedactedURNs, "cannot query on redacted URNs"))
		}

	} else if urns.IsValidScheme(propKey) {
		// second try to match a URN scheme
		propType = PropertyTypeScheme
		valueType = assets.FieldTypeText

		if v.env.RedactionPolicy() == envs.RedactionPolicyURNs && value != "" {
			v.addError(NewQueryError(ErrRedactedURNs, "cannot query on redacted URNs"))
		}
	} else {
		field := v.resolver.ResolveField(propKey)
		if field != nil {
			propType = PropertyTypeField
			propField = field
			valueType = field.Type()
		} else {
			v.addError(NewQueryError(ErrUnknownProperty, "can't resolve '%s' to attribute, scheme or field", propKey).withExtra("property", propKey))
		}
	}

	condition := newCondition(propKey, propType, propField, operator, value, valueType)

	if err := condition.Validate(v.env, v.resolver); err != nil {
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
