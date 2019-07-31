package contactql

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/nyaruka/goflow/contactql/gen"
	"github.com/nyaruka/goflow/envs"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
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

// QueryNode is the base for nodes in our query parse tree
type QueryNode interface {
	fmt.Stringer
	Evaluate(envs.Environment, Queryable) (bool, error)
}

// Condition represents a comparison between a keywed value on the contact and a provided value
type Condition struct {
	propType   PropertyType
	propKey    string
	comparator string
	value      string
}

func newCondition(propType PropertyType, propKey string, comparator string, value string) *Condition {
	return &Condition{propType: propType, propKey: propKey, comparator: comparator, value: value}
}

// PropertyKey returns the key for the property being queried
func (c *Condition) PropertyKey() string { return c.propKey }

// PropertyType returns the type (attribute, scheme, field)
func (c *Condition) PropertyType() PropertyType { return c.propType }

// Comparator returns the type of comparison being made
func (c *Condition) Comparator() string { return c.comparator }

// Value returns the value being compared against
func (c *Condition) Value() string { return c.value }

// Evaluate evaluates this condition against the queryable contact
func (c *Condition) Evaluate(env envs.Environment, queryable Queryable) (bool, error) {
	// contacts can return multiple values per key, e.g. multiple phone numbers in a "tel = x" condition
	vals := queryable.QueryProperty(env, c.PropertyKey(), c.PropertyType())

	// is this an existence check?
	if c.value == "" {
		if c.comparator == "=" {
			return len(vals) == 0, nil // x = "" is true if x doesn't exist
		} else if c.comparator == "!=" {
			return len(vals) > 0, nil // x != "" is false if x doesn't exist (i.e. true if x does exist)
		}
	}

	// if keyed value doesn't exist on our contact then all other comparisons at this point are false
	if len(vals) == 0 {
		return false, nil
	}

	// check each resolved value
	for _, val := range vals {
		res, err := c.evaluateValue(env, val)
		if err != nil {
			return false, err
		}
		if res {
			return true, nil
		}
	}

	return false, nil
}

func (c *Condition) evaluateValue(env envs.Environment, val interface{}) (bool, error) {
	switch val.(type) {
	case string:
		return textComparison(val.(string), c.comparator, c.value)

	case decimal.Decimal:
		asDecimal, err := decimal.NewFromString(c.value)
		if err != nil {
			return false, errors.Errorf("can't convert '%s' to a number", c.value)
		}
		return numberComparison(val.(decimal.Decimal), c.comparator, asDecimal)

	case time.Time:
		asDate, err := envs.DateTimeFromString(env, c.value, false)
		if err != nil {
			return false, err
		}
		return dateComparison(val.(time.Time), c.comparator, asDate)

	default:
		return false, errors.Errorf("unsupported query data type: %+v", reflect.TypeOf(val))
	}
}

func (c *Condition) String() string {
	var value string
	if c.value == "" {
		value = `""`
	} else {
		value = c.value
	}
	return fmt.Sprintf("%s%s%s", c.propKey, c.comparator, value)
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
	return fmt.Sprintf("%s(%s)", strings.ToUpper(string(b.op)), strings.Join(children, ", "))
}

type ContactQuery struct {
	root QueryNode
}

func (q *ContactQuery) Root() QueryNode { return q.root }

func (q *ContactQuery) Evaluate(env envs.Environment, queryable Queryable) (bool, error) {
	return q.root.Evaluate(env, queryable)
}

func (q *ContactQuery) String() string {
	return q.root.String()
}

// ParseQuery parses a ContactQL query from the given input
func ParseQuery(text string, redaction envs.RedactionPolicy) (*ContactQuery, error) {
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

	visitor := newVisitor(redaction)
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
	return errors.Errorf(strings.Join(l.messages, "\n"))
}

func (l *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	l.messages = append(l.messages, msg)
}
