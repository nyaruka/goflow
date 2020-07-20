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

// Comparator is a way of comparing two values in a condition
type Comparator string

// supported comparators
const (
	ComparatorEqual              Comparator = "="
	ComparatorNotEqual           Comparator = "!="
	ComparatorContains           Comparator = "~"
	ComparatorGreaterThan        Comparator = ">"
	ComparatorLessThan           Comparator = "<"
	ComparatorGreaterThanOrEqual Comparator = ">="
	ComparatorLessThanOrEqual    Comparator = "<="
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
	propType   PropertyType
	propKey    string
	comparator Comparator
	value      string
	valueType  assets.FieldType
	reference  assets.Reference
}

func newCondition(propType PropertyType, propKey string, comparator Comparator, value string, valueType assets.FieldType, reference assets.Reference) *Condition {
	return &Condition{
		propType:   propType,
		propKey:    propKey,
		comparator: comparator,
		value:      value,
		valueType:  valueType,
		reference:  reference,
	}
}

// PropertyKey returns the key for the property being queried
func (c *Condition) PropertyKey() string { return c.propKey }

// PropertyType returns the type (attribute, scheme, field)
func (c *Condition) PropertyType() PropertyType { return c.propType }

// Comparator returns the type of comparison being made
func (c *Condition) Comparator() Comparator { return c.comparator }

// Value returns the value being compared against
func (c *Condition) Value() string { return c.value }

// Validate checks that this condition is valid (and thus can be evaluated)
func (c *Condition) Validate(resolver Resolver) error {
	switch c.comparator {
	case ComparatorContains:
		if c.propKey == AttributeName {
			if len(tokenizeNameValue(c.value)) == 0 {
				return NewQueryErrorf("value must contain a word of at least %d characters long for a contains condition on name", minNameTokenContainsLength)
			}
		} else if c.propKey == AttributeURN || c.propType == PropertyTypeScheme {
			if len(c.value) < minURNContainsLength {
				return NewQueryErrorf("value must be least %d characters long for a contains condition on a URN", minURNContainsLength)
			}
		} else {
			// ~ can only be used with the name/urn attributes or actual URNs
			return NewQueryErrorf("contains conditions can only be used with name or URN values")
		}

	case ComparatorGreaterThan, ComparatorGreaterThanOrEqual, ComparatorLessThan, ComparatorLessThanOrEqual:
		if c.valueType != assets.FieldTypeNumber && c.valueType != assets.FieldTypeDatetime {
			return NewQueryErrorf("comparisons with %s can only be used with date and number fields", c.comparator)
		}
	}

	// if existence check, disallow certain attributes
	if c.value == "" {
		switch c.propKey {
		case AttributeUUID, AttributeID, AttributeCreatedOn, AttributeGroup:
			return NewQueryErrorf("can't check whether '%s' is set or not set", c.propKey)
		}
	} else {
		// check values are valid for the attribute type
		switch c.propKey {
		case AttributeGroup:
			group := resolver.ResolveGroup(c.value)
			if group == nil {
				return NewQueryErrorf("'%s' is not a valid group name", c.value)
			}
			c.value = group.Name()
			c.reference = assets.NewGroupReference(group.UUID(), group.Name())
		case AttributeLanguage:
			if c.value != "" {
				_, err := envs.ParseLanguage(c.value)
				if err != nil {
					return NewQueryErrorf("'%s' is not a valid language code", c.value)
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
		if c.comparator == ComparatorEqual {
			return len(vals) == 0, nil // x = "" is true if x doesn't exist
		} else if c.comparator == ComparatorNotEqual {
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
		res, err := c.evaluateValue(env, val)
		if err != nil {
			return false, err
		}
		if res {
			anyTrue = true
		} else {
			allTrue = false
		}
	}

	// foo != x is only true if all values of foo are not x
	if c.comparator == ComparatorNotEqual {
		return allTrue, nil
	}

	// foo = x is true if any value of foo is x
	return anyTrue, nil
}

func (c *Condition) evaluateValue(env envs.Environment, val interface{}) (bool, error) {
	switch val.(type) {
	case string:
		isName := c.propKey == AttributeName // needs to be handled as special case

		return textComparison(val.(string), c.comparator, c.value, isName)

	case decimal.Decimal:
		asDecimal, err := decimal.NewFromString(c.value)
		if err != nil {
			return false, NewQueryErrorf("can't convert '%s' to a number", c.value)
		}
		return numberComparison(val.(decimal.Decimal), c.comparator, asDecimal)

	case time.Time:
		asDate, err := envs.DateTimeFromString(env, c.value, false)
		if err != nil {
			return false, NewQueryErrorf("can't convert '%s' to a date", c.value)
		}
		return dateComparison(val.(time.Time), c.comparator, asDate)
	}

	panic(fmt.Sprintf("unsupported query data type: %T", val))
}

func (c *Condition) String() string {
	value := c.value

	if !isNumberRegex.MatchString(value) {
		// if not a decimal then quote
		value = strconv.Quote(value)
	}

	return fmt.Sprintf(`%s %s %s`, c.propKey, c.comparator, value)
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
func ParseQuery(text string, redaction envs.RedactionPolicy, country envs.Country, resolver Resolver) (*ContactQuery, error) {
	// preprocess text before parsing
	text = strings.TrimSpace(text)

	// if query is a valid number, rewrite as a tel = query
	if redaction != envs.RedactionPolicyURNs {
		if number := utils.ParsePhoneNumber(text, string(country)); number != "" {
			text = fmt.Sprintf(`tel = %s`, number)
		}
	}

	errListener := NewErrorListener()
	input := antlr.NewInputStream(text)
	lexer := gen.NewContactQLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := gen.NewContactQLParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(errListener)
	tree := p.Parse()

	// if we ran into errors parsing, bail
	if errListener.HasErrors() {
		return nil, errListener.Error()
	}

	visitor := newVisitor(redaction, resolver)
	rootNode := visitor.Visit(tree).(QueryNode)

	if len(visitor.errors) > 0 {
		return nil, visitor.errors[0]
	}

	return &ContactQuery{root: rootNode}, nil
}

type errorListener struct {
	*antlr.DefaultErrorListener

	messages []string
}

func NewErrorListener() *errorListener {
	return &errorListener{}
}

func (l *errorListener) HasErrors() bool {
	return len(l.messages) > 0
}

func (l *errorListener) Error() error {
	return NewQueryErrorf(strings.Join(l.messages, "\n"))
}

func (l *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	l.messages = append(l.messages, msg)
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
