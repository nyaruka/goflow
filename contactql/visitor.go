package contactql

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/nyaruka/gocommon/urns"
	gen "github.com/nyaruka/goflow/antlr/gen/contactql"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils/obfuscate"
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
	AttributeID         = "id" // deprecated in favor of ref
	AttributeRef        = "ref"
	AttributeName       = "name"
	AttributeStatus     = "status"
	AttributeLanguage   = "language"
	AttributeURN        = "urn"
	AttributeGroup      = "group"
	AttributeFlow       = "flow"
	AttributeHistory    = "history"
	AttributeTickets    = "tickets"
	AttributeCreatedOn  = "created_on"
	AttributeLastSeenOn = "last_seen_on"
)

var attributes = map[string]assets.FieldType{
	AttributeUUID:       assets.FieldTypeText,
	AttributeID:         assets.FieldTypeText,
	AttributeRef:        assets.FieldTypeText,
	AttributeName:       assets.FieldTypeText,
	AttributeStatus:     assets.FieldTypeText,
	AttributeLanguage:   assets.FieldTypeText,
	AttributeURN:        assets.FieldTypeText,
	AttributeGroup:      assets.FieldTypeText,
	AttributeFlow:       assets.FieldTypeText,
	AttributeHistory:    assets.FieldTypeText,
	AttributeTickets:    assets.FieldTypeNumber,
	AttributeCreatedOn:  assets.FieldTypeDatetime,
	AttributeLastSeenOn: assets.FieldTypeDatetime,
}

// Resolver provides functions for resolving assets referenced in queries
type Resolver interface {
	ResolveField(key string) assets.Field
	ResolveGroup(name string) assets.Group
	ResolveFlow(name string) assets.Flow
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

	asURN, _ := urns.Parse(value)

	if v.env.RedactionPolicy() == envs.RedactionPolicyURNs {
		if obfuscate.WasID(value) {
			return NewCondition(PropertyTypeAttribute, AttributeRef, OpEqual, value)
		}
	} else if asURN != urns.NilURN {
		scheme, path, _, _ := asURN.ToParts()

		return NewCondition(PropertyTypeURN, scheme, OpEqual, path)

	} else if implicitIsPhoneNumberRegex.MatchString(value) {
		value = cleanPhoneNumberRegex.ReplaceAllLiteralString(value, "")

		return NewCondition(PropertyTypeURN, urns.Phone.Prefix, OpContains, value)
	}

	// convert to contains condition only if we have the right tokens, otherwise make equals check
	operator := OpContains
	if len(tokenizeNameValue(value)) == 0 {
		operator = OpEqual
	}

	return NewCondition(PropertyTypeAttribute, AttributeName, operator, value)
}

// expression : PROPERTY COMPARATOR literal
func (v *visitor) VisitCondition(ctx *gen.ConditionContext) any {
	propText := strings.ToLower(ctx.PROPERTY().GetText())
	operatorText := strings.ToLower(ctx.COMPARATOR().GetText())
	value := v.Visit(ctx.Literal()).(string)

	operator, isAlias := operatorAliases[operatorText]
	if !isAlias {
		operator = Operator(operatorText)
	}

	var propType PropertyType
	var propKey string

	// check if property type is specified as prefix
	if strings.Contains(propText, ".") {
		parts := strings.SplitN(propText, ".", 2)

		if parts[0] == "fields" {
			propType = PropertyTypeField
			propKey = parts[1]
		} else if parts[0] == "urns" {
			propType = PropertyTypeURN
			propKey = parts[1]
		} else {
			v.addError(NewQueryError(ErrUnknownPropertyType, fmt.Sprintf("unknown property type '%s'", parts[0])).withExtra("type", parts[0]))
		}
	} else {
		propKey = propText

		// first try to match a fixed attribute
		_, isAttribute := attributes[propKey]
		if isAttribute {
			propType = PropertyTypeAttribute

			if propKey == AttributeURN && v.env.RedactionPolicy() == envs.RedactionPolicyURNs && value != "" {
				v.addError(NewQueryError(ErrRedactedURNs, "cannot query on redacted URNs"))
			}

		} else if urns.IsValidScheme(propKey) {
			// second try to match a URN scheme
			propType = PropertyTypeURN

			if v.env.RedactionPolicy() == envs.RedactionPolicyURNs && value != "" {
				v.addError(NewQueryError(ErrRedactedURNs, "cannot query on redacted URNs"))
			}
		} else {
			propType = PropertyTypeField
		}
	}

	return NewCondition(propType, propKey, operator, value)
}

// expression : expression AND expression
func (v *visitor) VisitCombinationAnd(ctx *gen.CombinationAndContext) any {
	child1 := v.Visit(ctx.Expression(0)).(QueryNode)
	child2 := v.Visit(ctx.Expression(1)).(QueryNode)
	return NewBoolCombination(BoolOperatorAnd, child1, child2)
}

// expression : expression expression
func (v *visitor) VisitCombinationImpicitAnd(ctx *gen.CombinationImpicitAndContext) any {
	child1 := v.Visit(ctx.Expression(0)).(QueryNode)
	child2 := v.Visit(ctx.Expression(1)).(QueryNode)
	return NewBoolCombination(BoolOperatorAnd, child1, child2)
}

// expression : expression OR expression
func (v *visitor) VisitCombinationOr(ctx *gen.CombinationOrContext) any {
	child1 := v.Visit(ctx.Expression(0)).(QueryNode)
	child2 := v.Visit(ctx.Expression(1)).(QueryNode)
	return NewBoolCombination(BoolOperatorOr, child1, child2)
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
