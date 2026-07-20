package contactql

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/obfuscate"
	"github.com/shopspring/decimal"
)

// MaxConditions is the maximum number of conditions a query can contain. Every condition becomes a clause
// when the query is translated for a search backend, so this bounds how much work an untrusted query can
// create both here and downstream. It's far more than any hand written query needs.
const MaxConditions = 1000

// Operator is a comparison operation between two values in a condition
type Operator string

// supported operators
const (
	OpEqual              Operator = "="
	OpNotEqual           Operator = "!="
	OpContains           Operator = "~"
	OpGreaterThan        Operator = ">"
	OpLessThan           Operator = "<"
	OpGreaterThanOrEqual Operator = ">="
	OpLessThanOrEqual    Operator = "<="
)

// BoolOperator is a boolean operator (and or or)
type BoolOperator string

const (
	// BoolOperatorAnd is our constant for an AND operation
	BoolOperatorAnd BoolOperator = "and"

	// BoolOperatorOr is our constant for an OR operation
	BoolOperatorOr BoolOperator = "or"
)

// PropertyType is the type of the lefthand side of a condition
type PropertyType string

const (
	// PropertyTypeAttribute is builtin property
	PropertyTypeAttribute PropertyType = "attr"

	// PropertyTypeURN is a URN scheme
	PropertyTypeURN PropertyType = "urn"

	// PropertyTypeField is a custom contact field
	PropertyTypeField PropertyType = "field"
)

// name based contains conditions are tokenized but only tokens of at least 2 characters are used
const minNameTokenContainsLength = 2

// URN based contains conditions ust be at least 3 characters long as the ES implementation uses trigrams
const minURNContainsLength = 3

var isNumberRegex = regexp.MustCompile(`^\d+(\.\d+)?$`)

// QueryNode is the base for nodes in our query parse tree
type QueryNode interface {
	fmt.Stringer
	validate(envs.Environment, Resolver) error
	Simplify() QueryNode
}

// Condition represents a comparison between a keywed value on the contact and a provided value
type Condition struct {
	propType PropertyType
	propKey  string
	operator Operator
	value    string
}

func NewCondition(propType PropertyType, propKey string, operator Operator, value string) *Condition {
	return &Condition{
		propType: propType,
		propKey:  propKey,
		operator: operator,
		value:    value,
	}
}

// an implicit condition like +123-124-6546 or 1234 will be interpreted as a tel ~ condition
var implicitIsPhoneNumberRegex = regexp.MustCompile(`^\+?[\-\d]{4,}$`)

// used to strip formatting from phone number values
var cleanPhoneNumberRegex = regexp.MustCompile(`[^+\d]+`)

// NewImplicitCondition interprets a bare literal as a condition, e.g. as a URN, phone number or name.
func NewImplicitCondition(env envs.Environment, value string) *Condition {
	asURN, _ := urns.Parse(value)

	if env.RedactionPolicy() == envs.RedactionPolicyURNs {
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

// PropertyType returns the type (attribute, scheme, field)
func (c *Condition) PropertyType() PropertyType { return c.propType }

// PropertyKey returns the key for the property being queried
func (c *Condition) PropertyKey() string { return c.propKey }

// Operator returns the type of comparison being made
func (c *Condition) Operator() Operator { return c.operator }

// Value returns the value being compared against
func (c *Condition) Value() string { return c.value }

// ValueAsNumber returns the value as a number if possible, or an error if not
func (c *Condition) ValueAsNumber() (decimal.Decimal, error) {
	return decimal.NewFromString(c.value)
}

// ValueAsDate returns the value as a date if possible, or an error if not
func (c *Condition) ValueAsDate(env envs.Environment) (time.Time, error) {
	return envs.DateTimeFromString(env, c.value, false)
}

// ValueAsGroup returns the value as a group if possible
func (c *Condition) ValueAsGroup(resolver Resolver) assets.Group {
	return resolver.ResolveGroup(c.value)
}

// ValueAsFlow returns the value as a flow if possible
func (c *Condition) ValueAsFlow(resolver Resolver) assets.Flow {
	return resolver.ResolveFlow(c.value)
}

func (c *Condition) resolveValueType(resolver Resolver) assets.FieldType {
	switch c.propType {
	case PropertyTypeAttribute:
		return attributes[c.propKey]
	case PropertyTypeURN:
		return assets.FieldTypeText
	case PropertyTypeField:
		field := resolver.ResolveField(c.propKey)
		if field != nil {
			return field.Type()
		}
	}
	return ""
}

// Validate checks that this condition is valid (and thus can be evaluated)
func (c *Condition) validate(env envs.Environment, resolver Resolver) error {
	// if URNs are redacted, block any conditions on them
	if (c.propType == PropertyTypeURN || c.propKey == AttributeURN) && env.RedactionPolicy() == envs.RedactionPolicyURNs && c.value != "" {
		return NewQueryError(ErrRedactedURNs, "cannot query on redacted URNs")
	}

	// if our property is a field and we don't have a resolver, we can't validate because we don't know the value type
	if c.propType == PropertyTypeField && resolver == nil {
		return nil
	}

	valueType := c.resolveValueType(resolver)
	if valueType == "" {
		return NewQueryError(ErrUnknownProperty, fmt.Sprintf("can't resolve '%s' to attribute, scheme or field", c.propKey)).WithExtra("property", c.propKey)
	}

	switch c.operator {
	case OpContains:
		if c.propKey == AttributeName {
			if len(tokenizeNameValue(c.value)) == 0 {
				return NewQueryError(ErrInvalidPartialName, fmt.Sprintf("contains operator on name requires token of minimum length %d", minNameTokenContainsLength)).WithExtra("min_token_length", strconv.Itoa(minNameTokenContainsLength))
			}
		} else if c.propKey == AttributeURN || c.propType == PropertyTypeURN {
			if len(c.value) < minURNContainsLength {
				return NewQueryError(ErrInvalidPartialURN, fmt.Sprintf("contains operator on URN requires value of minimum length %d", minURNContainsLength)).WithExtra("min_value_length", strconv.Itoa(minURNContainsLength))
			}
		} else {
			// ~ can only be used with the name/urn attributes or actual URNs
			return NewQueryError(ErrUnsupportedContains, "contains conditions can only be used with name or URN values").WithExtra("property", c.propKey)
		}

	case OpGreaterThan, OpGreaterThanOrEqual, OpLessThan, OpLessThanOrEqual:
		if valueType != assets.FieldTypeNumber && valueType != assets.FieldTypeDatetime {
			return NewQueryError(ErrUnsupportedComparison, fmt.Sprintf("comparisons with %s can only be used with date and number fields", c.operator)).WithExtra("property", c.propKey).WithExtra("operator", string(c.operator))
		}
	}

	// if existence check, disallow certain attributes
	if (c.operator == OpEqual || c.operator == OpNotEqual) && c.value == "" {
		switch c.propKey {
		case AttributeUUID, AttributeID, AttributeRef, AttributeStatus, AttributeCreatedOn, AttributeTickets:
			return NewQueryError(ErrUnsupportedSetCheck, fmt.Sprintf("can't check whether '%s' is set or not set", c.propKey)).WithExtra("property", c.propKey).WithExtra("operator", string(c.operator))
		}
	} else {
		// check values are valid for the value type
		if valueType == assets.FieldTypeNumber {
			_, err := c.ValueAsNumber()
			if err != nil {
				return NewQueryError(ErrInvalidNumber, fmt.Sprintf("can't convert '%s' to a number", c.value)).WithExtra("value", c.value)
			}
		} else if valueType == assets.FieldTypeDatetime {
			_, err := c.ValueAsDate(env)
			if err != nil {
				return NewQueryError(ErrInvalidDate, fmt.Sprintf("can't convert '%s' to a date", c.value)).WithExtra("value", c.value)
			}
		}

		// for some text attributes, do some additional validation
		if c.propType == PropertyTypeAttribute {
			if c.propKey == AttributeGroup && resolver != nil {
				group := c.ValueAsGroup(resolver)
				if group == nil {
					return NewQueryError(ErrInvalidGroup, fmt.Sprintf("'%s' is not a valid group name", c.value)).WithExtra("value", c.value)
				}
			} else if (c.propKey == AttributeFlow || c.propKey == AttributeHistory) && resolver != nil {
				flow := c.ValueAsFlow(resolver)
				if flow == nil {
					return NewQueryError(ErrInvalidFlow, fmt.Sprintf("'%s' is not a valid flow name", c.value)).WithExtra("value", c.value)
				}
			} else if c.propKey == AttributeStatus {
				val := strings.ToLower(c.value)
				if val != "active" && val != "blocked" && val != "stopped" && val != "archived" {
					return NewQueryError(ErrInvalidStatus, fmt.Sprintf("'%s' is not a valid contact status", c.value)).WithExtra("value", c.value)
				}
			} else if c.propKey == AttributeLanguage {
				if c.value != "" {
					_, err := i18n.ParseLanguage(c.value)
					if err != nil {
						return NewQueryError(ErrInvalidLanguage, fmt.Sprintf("'%s' is not a valid language code", c.value)).WithExtra("value", c.value)
					}
				}
			}
		}
	}

	return nil
}

func (c *Condition) Simplify() QueryNode {
	return c
}

func (c *Condition) String() string {
	property := c.propKey
	value := c.value

	// add prefix for fields and URNs
	if c.propType == PropertyTypeField {
		property = fmt.Sprintf(`fields.%s`, property)
	} else if c.propType == PropertyTypeURN {
		property = fmt.Sprintf(`urns.%s`, property)
	}

	if !isNumberRegex.MatchString(value) {
		// if not a decimal then quote
		value = strconv.Quote(value)
	}

	return fmt.Sprintf(`%s %s %s`, property, c.operator, value)
}

// BoolCombination is a AND or OR combination of multiple conditions
type BoolCombination struct {
	op       BoolOperator
	children []QueryNode
}

// Operator returns the type of boolean operator this combination is
func (b *BoolCombination) Operator() BoolOperator { return b.op }

// Children returns the children of this boolean combination
func (b *BoolCombination) Children() []QueryNode { return b.children }

// NewBoolCombination creates a new boolean combination
func NewBoolCombination(op BoolOperator, children ...QueryNode) *BoolCombination {
	return &BoolCombination{op: op, children: children}
}

// Validate validates this node
func (b *BoolCombination) validate(env envs.Environment, resolver Resolver) error {
	for _, child := range b.children {
		if err := child.validate(env, resolver); err != nil {
			return err
		}
	}
	return nil
}

func (b *BoolCombination) Simplify() QueryNode {
	var newChildren []QueryNode

	// simplify by promoting grand children to children if they're combined with same op
	for _, child := range b.children {
		// let children remove themselves by simplifying to nil
		sc := child.Simplify()
		if sc == nil {
			continue
		}

		switch typed := sc.(type) {
		case *BoolCombination:
			if typed.op == b.op {
				// adopt the grand children slice the first time we promote rather than copying into a new
				// one. Parsing a chain like a AND b AND c gives a left leaning tree, so copying at every
				// level would make flattening it quadratic in the length of the chain. The slice we adopt
				// was freshly built by the child's own Simplify and isn't referenced anywhere else.
				if newChildren == nil {
					newChildren = typed.children
				} else {
					newChildren = append(newChildren, typed.children...)
				}
			} else {
				newChildren = append(newChildren, typed)
			}
		case *Condition:
			newChildren = append(newChildren, typed)
		}
	}

	// you can't parse a boolean combination with less than 2 children but you can construct one
	if len(newChildren) == 0 {
		return nil
	} else if len(newChildren) == 1 {
		return newChildren[0]
	}

	return &BoolCombination{op: b.op, children: newChildren}
}

func (b *BoolCombination) String() string {
	children := make([]string, len(b.children))
	for i := range b.children {
		children[i] = b.children[i].String()
	}
	return fmt.Sprintf("(%s)", strings.Join(children, fmt.Sprintf(" %s ", strings.ToUpper(string(b.op)))))
}

// ContactQuery is a parsed contact QL query
type ContactQuery struct {
	root     QueryNode
	resolver Resolver
}

// Root returns the root node of this query
func (q *ContactQuery) Root() QueryNode { return q.root }

// Resolver returns the optional resolver this query was parsed with
func (q *ContactQuery) Resolver() Resolver { return q.resolver }

// String returns the pretty formatted version of this query
func (q *ContactQuery) String() string {
	return Stringify(q.root)
}

// NewContactQuery creates a new query from the given root node, validating it against the given resolver
// (or as much as possible if none is provided) and simplifying it.
func NewContactQuery(env envs.Environment, root QueryNode, resolver Resolver) (*ContactQuery, error) {
	// bound overall complexity before doing any per-condition work. Each condition becomes a clause when
	// the query is translated for a search backend, so an oversized query is expensive for every consumer
	// and not just for us.
	numConditions := 0
	walk(root, func(*Condition) { numConditions++ })
	if numConditions > MaxConditions {
		return nil, NewQueryError(ErrTooComplex, fmt.Sprintf("query contains more than %d conditions", MaxConditions))
	}

	if err := root.validate(env, resolver); err != nil {
		return nil, err
	}

	return &ContactQuery{root: root.Simplify(), resolver: resolver}, nil
}

// EscapeValue escapes a value for inclusion in a query, e.g. as the output of an evaluated expression
func EscapeValue(s string) string {
	return strconv.Quote(s)
}

// Stringify converts a query node to a string
func Stringify(n QueryNode) string {
	// since simplfying can remove nodes and potentially generate a nil query
	if n == nil {
		return ""
	}

	s := n.String()

	// bool combinations are wrapped in parentheses but the top level doesn't need to be
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		s = s[1 : len(s)-1]
	}
	return s
}

func tokenizeNameValue(value string) []string {
	tokens := make([]string, 0)
	for _, token := range utils.TokenizeStringByUnicodeSeg(value) {
		if len(token) >= minNameTokenContainsLength {
			tokens = append(tokens, token)
		}
	}
	return tokens
}
