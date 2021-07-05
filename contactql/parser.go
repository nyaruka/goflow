package contactql

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql/gen"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"

	"github.com/antlr/antlr4/runtime/Go/antlr"
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
	Evaluate(envs.Environment, Queryable) (bool, error)
}

// Condition represents a comparison between a keywed value on the contact and a provided value
type Condition struct {
	propKey       string
	propType      PropertyType
	propField     assets.Field
	operator      Operator
	value         string
	valueAsNumber decimal.Decimal
	valueAsDate   time.Time
	valueAsGroup  assets.Group
	valueType     assets.FieldType
}

func newCondition(propKey string, propType PropertyType, propField assets.Field, operator Operator, value string, valueType assets.FieldType) *Condition {
	return &Condition{
		propKey:   propKey,
		propType:  propType,
		propField: propField,
		operator:  operator,
		value:     value,
		valueType: valueType,
	}
}

// PropertyKey returns the key for the property being queried
func (c *Condition) PropertyKey() string { return c.propKey }

// PropertyType returns the type (attribute, scheme, field)
func (c *Condition) PropertyType() PropertyType { return c.propType }

// PropertyField returns the field for the property being queried if it's a field
func (c *Condition) PropertyField() assets.Field { return c.propField }

// Operator returns the type of comparison being made
func (c *Condition) Operator() Operator { return c.operator }

// Value returns the value being compared against
func (c *Condition) Value() string { return c.value }

// ValueAsNumber returns the value as a number if value type is number
func (c *Condition) ValueAsNumber() decimal.Decimal { return c.valueAsNumber }

// ValueAsDate returns the value as a date if condition is datetime
func (c *Condition) ValueAsDate() time.Time { return c.valueAsDate }

// ValueAsGroup returns the value as a group if condition is on the group attribute
func (c *Condition) ValueAsGroup() assets.Group { return c.valueAsGroup }

// Validate checks that this condition is valid (and thus can be evaluated)
func (c *Condition) Validate(env envs.Environment, resolver Resolver) error {
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
		if c.valueType != assets.FieldTypeNumber && c.valueType != assets.FieldTypeDatetime {
			return NewQueryError(ErrUnsupportedComparison, "comparisons with %s can only be used with date and number fields", c.operator).withExtra("property", c.propKey).withExtra("operator", string(c.operator))
		}
	}

	// if existence check, disallow certain attributes
	if (c.operator == OpEqual || c.operator == OpNotEqual) && c.value == "" {
		switch c.propKey {
		case AttributeUUID, AttributeID, AttributeCreatedOn, AttributeGroup:
			return NewQueryError(ErrUnsupportedSetCheck, "can't check whether '%s' is set or not set", c.propKey).withExtra("property", c.propKey).withExtra("operator", string(c.operator))
		}
	} else {
		// check values are valid for the attribute type
		if c.valueType == assets.FieldTypeNumber {
			asDecimal, err := decimal.NewFromString(c.value)
			if err != nil {
				return NewQueryError(ErrInvalidNumber, "can't convert '%s' to a number", c.value).withExtra("value", c.value)
			}
			c.valueAsNumber = asDecimal

		} else if c.valueType == assets.FieldTypeDatetime {
			asDate, err := envs.DateTimeFromString(env, c.value, false)
			if err != nil {
				return NewQueryError(ErrInvalidDate, "can't convert '%s' to a date", c.value).withExtra("value", c.value)
			}
			c.valueAsDate = asDate

		} else if c.propKey == AttributeGroup {
			group := resolver.ResolveGroup(c.value)
			if group == nil {
				return NewQueryError(ErrInvalidGroup, "'%s' is not a valid group name", c.value).withExtra("value", c.value)
			}
			c.value = group.Name()
			c.valueAsGroup = group

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

// Evaluate evaluates this condition against the queryable contact
func (c *Condition) Evaluate(env envs.Environment, queryable Queryable) (bool, error) {
	// contacts can return multiple values per key, e.g. multiple phone numbers in a "tel = x" condition
	vals := queryable.QueryProperty(env, c.PropertyKey(), c.PropertyType())

	// is this an existence check?
	if c.value == "" {
		if c.operator == OpEqual {
			return len(vals) == 0, nil // x = "" is true if x doesn't exist
		} else if c.operator == OpNotEqual {
			return len(vals) > 0, nil // x != "" is false if x doesn't exist (i.e. true if x does exist)
		}
	}

	// if keyed value doesn't exist on our contact then all other comparisons at this point are false
	if len(vals) == 0 {
		return false, nil
	}

	// evaluate condition against each resolved value
	anyTrue := false
	allTrue := true
	for _, val := range vals {
		if c.evaluateValue(env, val) {
			anyTrue = true
		} else {
			allTrue = false
		}
	}

	// foo != x is only true if all values of foo are not x
	if c.operator == OpNotEqual {
		return allTrue, nil
	}

	// foo = x is true if any value of foo is x
	return anyTrue, nil
}

func (c *Condition) evaluateValue(env envs.Environment, val interface{}) bool {
	switch typed := val.(type) {
	case string:
		isName := c.propKey == AttributeName // needs to be handled as special case

		return textComparison(typed, c.operator, c.value, isName)

	case decimal.Decimal:
		return numberComparison(typed, c.operator, c.valueAsNumber)

	case time.Time:
		return dateComparison(typed, c.operator, c.valueAsDate)
	}

	panic(fmt.Sprintf("unsupported query data type: %T", val))
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

// Evaluate returns whether this combination evaluates to true or false
func (b *BoolCombination) Evaluate(env envs.Environment, queryable Queryable) (bool, error) {
	var childRes bool
	var err error

	if b.op == BoolOperatorAnd {
		for _, child := range b.children {
			if childRes, err = child.Evaluate(env, queryable); err != nil {
				return false, err
			}
			if !childRes {
				return false, nil
			}
		}
		return true, nil
	}

	for _, child := range b.children {
		if childRes, err = child.Evaluate(env, queryable); err != nil {
			return false, err
		}
		if childRes {
			return true, nil
		}
	}
	return false, nil
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
	root QueryNode
}

// Root returns the root node of this query
func (q *ContactQuery) Root() QueryNode { return q.root }

// Evaluate returns whether the given queryable matches this query
func (q *ContactQuery) Evaluate(env envs.Environment, queryable Queryable) (bool, error) {
	return q.root.Evaluate(env, queryable)
}

// String returns the pretty formatted version of this query
func (q *ContactQuery) String() string {
	s := q.root.String()

	// strip extra parentheses if not needed
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		s = s[1 : len(s)-1]
	}
	return s
}

// ParseQuery parses a ContactQL query from the given input
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

	visitor := newVisitor(env, resolver)
	rootNode := visitor.Visit(tree).(QueryNode)

	if len(visitor.errors) > 0 {
		return nil, visitor.errors[0]
	}

	return &ContactQuery{root: rootNode}, nil
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
