package contactql

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/nyaruka/goflow/contactql/gen"
	"github.com/nyaruka/goflow/utils"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type boolOp string

const (
	boolOpAnd boolOp = "and"
	boolOpOr  boolOp = "or"
)

// QueryNode is the base for nodes in our query parse tree
type QueryNode interface {
	fmt.Stringer

	Evaluate(utils.Environment, Queryable) (bool, error)
}

// Condition represents a comparison between a keywed value on the contact and a provided value
type Condition struct {
	key        string
	comparator string
	value      string
}

// Evaluate evaluates this condition against the queryable contact
func (c *Condition) Evaluate(env utils.Environment, queryable Queryable) (bool, error) {
	if c.key == ImplicitKey {
		return false, errors.Errorf("dynamic group queries can't contain implicit conditions")
	}

	// contacts can return multiple values per key, e.g. multiple phone numbers in a "tel = x" condition
	vals := queryable.ResolveQueryKey(env, c.key)

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

func (c *Condition) evaluateValue(env utils.Environment, val interface{}) (bool, error) {
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
		asDate, err := utils.DateTimeFromString(env, c.value, false)
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
	return fmt.Sprintf("%s%s%s", c.key, c.comparator, value)
}

// BoolCombination is a AND or OR combination of multiple conditions
type BoolCombination struct {
	op       boolOp
	children []QueryNode
}

// NewBoolCombination creates a new boolean combination
func NewBoolCombination(op boolOp, children ...QueryNode) *BoolCombination {
	return &BoolCombination{op: op, children: children}
}

// Evaluate returns whether this combination evaluates to true or false
func (b *BoolCombination) Evaluate(env utils.Environment, queryable Queryable) (bool, error) {
	var childRes bool
	var err error

	if b.op == boolOpAnd {
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
	for c := range b.children {
		children[c] = b.children[c].String()
	}
	return fmt.Sprintf("%s(%s)", strings.ToUpper(string(b.op)), strings.Join(children, ", "))
}

type ContactQuery struct {
	root QueryNode
}

func (q *ContactQuery) Evaluate(env utils.Environment, queryable Queryable) (bool, error) {
	return q.root.Evaluate(env, queryable)
}

func (q *ContactQuery) String() string {
	return q.root.String()
}

// ParseQuery parses a ContactQL query from the given input
func ParseQuery(text string) (*ContactQuery, error) {
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

	visitor := NewVisitor()
	rootNode := visitor.Visit(tree).(QueryNode)

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
