package contactql

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql/gen"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
	"github.com/shopspring/decimal"
)

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
	PropertyTypeAttribute PropertyType = "attribute"

	// PropertyTypeScheme is a URN scheme
	PropertyTypeScheme PropertyType = "scheme"

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
	propKey  string
	propType PropertyType
	operator Operator
	value    string
}

func NewCondition(propKey string, propType PropertyType, operator Operator, value string) *Condition {
	return &Condition{
		propKey:  propKey,
		propType: propType,
		operator: operator,
		value:    value,
	}
}

// PropertyKey returns the key for the property being queried
func (c *Condition) PropertyKey() string { return c.propKey }

// PropertyType returns the type (attribute, scheme, field)
func (c *Condition) PropertyType() PropertyType { return c.propType }

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
	case PropertyTypeScheme:
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
	// if our property is a field and we don't have a resolver, we can't validate because we don't know the value type
	if c.propType == PropertyTypeField && resolver == nil {
		return nil
	}

	valueType := c.resolveValueType(resolver)
	if valueType == "" {
		return NewQueryError(ErrUnknownProperty, "can't resolve '%s' to attribute, scheme or field", c.propKey).withExtra("property", c.propKey)
	}

	switch c.operator {
	case OpContains:
		if c.propKey == AttributeName {
			if len(tokenizeNameValue(c.value)) == 0 {
				return NewQueryError(ErrInvalidPartialName, "contains operator on name requires token of minimum length %d", minNameTokenContainsLength).withExtra("min_token_length", strconv.Itoa(minNameTokenContainsLength))
			}
		} else if c.propKey == AttributeURN || c.propType == PropertyTypeScheme {
			if len(c.value) < minURNContainsLength {
				return NewQueryError(ErrInvalidPartialURN, "contains operator on URN requires value of minimum length %d", minURNContainsLength).withExtra("min_value_length", strconv.Itoa(minURNContainsLength))
			}
		} else {
			// ~ can only be used with the name/urn attributes or actual URNs
			return NewQueryError(ErrUnsupportedContains, "contains conditions can only be used with name or URN values").withExtra("property", c.propKey)
		}

	case OpGreaterThan, OpGreaterThanOrEqual, OpLessThan, OpLessThanOrEqual:
		if valueType != assets.FieldTypeNumber && valueType != assets.FieldTypeDatetime {
			return NewQueryError(ErrUnsupportedComparison, "comparisons with %s can only be used with date and number fields", c.operator).withExtra("property", c.propKey).withExtra("operator", string(c.operator))
		}
	}

	// if existence check, disallow certain attributes
	if (c.operator == OpEqual || c.operator == OpNotEqual) && c.value == "" {
		switch c.propKey {
		case AttributeUUID, AttributeID, AttributeStatus, AttributeCreatedOn, AttributeTickets:
			return NewQueryError(ErrUnsupportedSetCheck, "can't check whether '%s' is set or not set", c.propKey).withExtra("property", c.propKey).withExtra("operator", string(c.operator))
		}
	} else {
		// check values are valid for the property type
		if valueType == assets.FieldTypeNumber {
			_, err := c.ValueAsNumber()
			if err != nil {
				return NewQueryError(ErrInvalidNumber, "can't convert '%s' to a number", c.value).withExtra("value", c.value)
			}
		} else if valueType == assets.FieldTypeDatetime {
			_, err := c.ValueAsDate(env)
			if err != nil {
				return NewQueryError(ErrInvalidDate, "can't convert '%s' to a date", c.value).withExtra("value", c.value)
			}

		} else if c.propKey == AttributeGroup && resolver != nil {
			group := c.ValueAsGroup(resolver)
			if group == nil {
				return NewQueryError(ErrInvalidGroup, "'%s' is not a valid group name", c.value).withExtra("value", c.value)
			}
		} else if (c.propKey == AttributeFlow || c.propKey == AttributeHistory) && resolver != nil {
			flow := c.ValueAsFlow(resolver)
			if flow == nil {
				return NewQueryError(ErrInvalidFlow, "'%s' is not a valid flow name", c.value).withExtra("value", c.value)
			}
		} else if c.propKey == AttributeStatus {
			val := strings.ToLower(c.value)
			if val != "active" && val != "blocked" && val != "stopped" && val != "archived" {
				return NewQueryError(ErrInvalidStatus, "'%s' is not a valid contact status", c.value).withExtra("value", c.value)
			}
		} else if c.propKey == AttributeLanguage {
			if c.value != "" {
				_, err := envs.ParseLanguage(c.value)
				if err != nil {
					return NewQueryError(ErrInvalidLanguage, "'%s' is not a valid language code", c.value).withExtra("value", c.value)
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
	value := c.value

	if !isNumberRegex.MatchString(value) {
		// if not a decimal then quote
		value = strconv.Quote(value)
	}

	return fmt.Sprintf(`%s %s %s`, c.propKey, c.operator, value)
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
		err := child.validate(env, resolver)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BoolCombination) Simplify() QueryNode {
	simplifiedChildren := make([]QueryNode, 0, len(b.children))

	for _, child := range b.children {
		// let children remove themselves by simplifying to nil
		sc := child.Simplify()
		if sc != nil {
			simplifiedChildren = append(simplifiedChildren, sc)
		}
	}

	newChildren := make([]QueryNode, 0, 2*len(simplifiedChildren))

	// simplify by promoting grand children to children if they're combined with same op
	for _, child := range simplifiedChildren {
		switch typed := child.(type) {
		case *BoolCombination:
			if typed.op == b.op {
				newChildren = append(newChildren, typed.children...)
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

// ParseQuery parses a ContactQL query from the given input. If resolver is provided then we validate against it
// to ensure that fields and groups exist. If not provided then still validate what we can.
func ParseQuery(env envs.Environment, text string, resolver Resolver) (*ContactQuery, error) {
	// preprocess text before parsing
	text = strings.TrimSpace(text)

	// if query is a valid number, rewrite as a tel = query
	if env.RedactionPolicy() != envs.RedactionPolicyURNs {
		if number := utils.ParsePhoneNumber(text, string(env.DefaultCountry())); number != "" {
			text = fmt.Sprintf(`tel = %s`, number)
		}
	}

	errListener := &errorListener{}
	input := antlr.NewInputStream(text)
	lexer := gen.NewContactQLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := gen.NewContactQLParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(errListener)
	tree := p.Parse()

	// if we ran into errors parsing, bail
	err := errListener.Error()
	if err != nil {
		return nil, err
	}

	visitor := newVisitor(env)
	rootNode := visitor.Visit(tree).(QueryNode)

	if len(visitor.errors) > 0 {
		return nil, visitor.errors[0]
	}

	if err := rootNode.validate(env, resolver); err != nil {
		return nil, err
	}

	rootNode = rootNode.Simplify()

	return &ContactQuery{root: rootNode, resolver: resolver}, nil
}

type errorListener struct {
	*antlr.DefaultErrorListener

	errs []*QueryError
}

func (l *errorListener) Error() error {
	if len(l.errs) > 0 {
		return l.errs[0]
	}
	return nil
}

func (l *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	var err *QueryError
	switch typed := e.(type) {
	case *antlr.InputMisMatchException:
		token := typed.GetOffendingToken().GetText()
		err = NewQueryError(ErrUnexpectedToken, msg).withExtra("token", token)
	default:
		err = NewQueryError("", msg)
	}

	l.errs = append(l.errs, err)
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
