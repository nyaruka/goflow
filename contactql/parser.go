package contactql

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/nyaruka/goflow/contactql/gen"
	"github.com/nyaruka/goflow/utils"
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
	val := queryable.ResolveQueryKey(c.key)

	// is this an existence check?
	if c.value == "" {
		if c.comparator == "=" {
			return utils.IsNil(val), nil // x = "" is true if x doesn't exist
		} else if c.comparator == "!=" {
			return !utils.IsNil(val), nil // x != "" is false if x doesn't exist (i.e. true if x does exist)
		}
	}

	// if keyed value doesn't exist on our contact then all other comparisons at this point are false
	if utils.IsNil(val) {
		return false, nil
	}

	if c.key == ImplicitKey {
		return implicitComparison(val.([]string), c.value), nil
	}

	switch val.(type) {
	case string:
		return stringComparison(val.(string), c.comparator, c.value)

	case decimal.Decimal:
		asDecimal, err := utils.ToDecimal(env, c.value)
		if err != nil {
			return false, err
		}
		return decimalComparison(val.(decimal.Decimal), c.comparator, asDecimal)

	case time.Time:
		asDate, err := utils.ToDate(env, c.value)
		if err != nil {
			return false, err
		}
		return dateComparison(val.(time.Time), c.comparator, asDate)
	}

	// TODO locations

	return false, fmt.Errorf("unsupported query data type %+v", reflect.TypeOf(val))
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

type BoolCombination struct {
	op       boolOp
	children []QueryNode
}

func NewBoolCombination(op boolOp, children ...QueryNode) *BoolCombination {
	return &BoolCombination{op: op, children: children}
}

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

func ParseQuery(text string) (*ContactQuery, error) {
	errors := newErrorListener()

	input := antlr.NewInputStream(text)
	lexer := gen.NewContactQLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := gen.NewContactQLParser(stream)
	tree := p.Parse()

	// if we ran into errors parsing, bail
	if errors.HasErrors() {
		return nil, fmt.Errorf(errors.Errors())
	}

	visitor := NewVisitor()
	rootNode := visitor.Visit(tree).(QueryNode)

	return &ContactQuery{root: rootNode}, nil
}

type errorListener struct {
	errors bytes.Buffer
	*antlr.DefaultErrorListener
}

func newErrorListener() *errorListener {
	return &errorListener{}
}

func (l *errorListener) HasErrors() bool {
	return l.errors.Len() > 0
}

func (l *errorListener) Errors() string {
	return l.errors.String()
}

func (l *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	l.errors.WriteString(fmt.Sprintln("line " + strconv.Itoa(line) + ":" + strconv.Itoa(column) + " " + msg))
}
