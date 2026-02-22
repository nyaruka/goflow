package excellent

import (
	"fmt"
	"slices"
	"strings"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/operators"
	"github.com/nyaruka/goflow/excellent/types"
)

type Warnings struct {
	all []string
}

func (w *Warnings) add(m string) {
	if !slices.Contains(w.all, m) {
		w.all = append(w.all, m)
	}
}

func (w *Warnings) deprecatedContext(v types.XValue) {
	w.add("deprecated context value accessed: " + v.Deprecated())
}

// Expression is the base interface of all syntax elements
type Expression interface {
	Evaluate(envs.Environment, *Scope, *Warnings) types.XValue
	Visit(func(Expression))
	String() string
}

// ContextReference is an identifier which is a function name or root variable in the context
type ContextReference struct {
	Name string
}

func (x *ContextReference) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	value, exists := scope.Get(x.Name)
	if !exists {
		return types.NewXErrorf("context has no property '%s'", x.Name)
	}

	if !types.IsNil(value) && value.Deprecated() != "" {
		warnings.deprecatedContext(value)
	}

	return value
}

func (x *ContextReference) Visit(v func(Expression)) {
	v(x)
}

func (x *ContextReference) String() string {
	return strings.ToLower(x.Name)
}

type DotLookup struct {
	Container Expression
	Lookup    string
}

func (x *DotLookup) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	containerVal := x.Container.Evaluate(env, scope, warnings)
	if types.IsXError(containerVal) {
		return containerVal
	}

	return resolveLookup(env, containerVal, types.NewXText(x.Lookup), true, warnings)
}

func (x *DotLookup) Visit(v func(Expression)) {
	x.Container.Visit(v)
	v(x)
}

func (x *DotLookup) String() string {
	return fmt.Sprintf("%s.%s", x.Container.String(), x.Lookup)
}

type ArrayLookup struct {
	Container Expression
	Lookup    Expression
}

func (x *ArrayLookup) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	containerVal := x.Container.Evaluate(env, scope, warnings)
	if types.IsXError(containerVal) {
		return containerVal
	}

	lookupVal := x.Lookup.Evaluate(env, scope, warnings)
	if types.IsXError(lookupVal) {
		return lookupVal
	}

	return resolveLookup(env, containerVal, lookupVal, false, warnings)
}

func (x *ArrayLookup) Visit(v func(Expression)) {
	x.Container.Visit(v)
	x.Lookup.Visit(v)
	v(x)
}

func (x *ArrayLookup) String() string {
	return fmt.Sprintf("%s[%s]", x.Container.String(), x.Lookup.String())
}

type FunctionCall struct {
	Func   Expression
	Params []Expression
}

func (x *FunctionCall) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	funcVal := x.Func.Evaluate(env, scope, warnings)
	if types.IsXError(funcVal) {
		return funcVal
	}

	asFunction, isFunction := funcVal.(*types.XFunction)
	if !isFunction {
		return types.NewXErrorf("%s is not a function", x.Func.String())
	}

	params := make([]types.XValue, len(x.Params))
	for i := range x.Params {
		params[i] = x.Params[i].Evaluate(env, scope, warnings)
	}

	return asFunction.Call(env, params)
}

func (x *FunctionCall) Visit(v func(Expression)) {
	x.Func.Visit(v)
	for _, p := range x.Params {
		p.Visit(v)
	}
	v(x)
}

func (x *FunctionCall) String() string {
	params := make([]string, len(x.Params))
	for i := range x.Params {
		params[i] = x.Params[i].String()
	}

	return fmt.Sprintf("%s(%s)", x.Func.String(), strings.Join(params, ", "))
}

type AnonFunction struct {
	Args []string
	Body Expression
}

func (x *AnonFunction) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	// create an XFunction which wraps our body expression
	fn := func(env envs.Environment, args ...types.XValue) types.XValue {
		// create new context that includes the args
		argsMap := make(map[string]types.XValue, len(x.Args))
		for i := range x.Args {
			argsMap[x.Args[i]] = args[i]
		}
		childScope := NewScope(types.NewXObject(argsMap), scope)

		return x.Body.Evaluate(env, childScope, warnings)
	}

	return types.NewXFunction("", functions.NumArgsCheck(len(x.Args), fn))
}

func (x *AnonFunction) Visit(v func(Expression)) {
	x.Body.Visit(v)
	v(x)
}

func (x *AnonFunction) String() string {
	return fmt.Sprintf("(%s) => %s", strings.Join(x.Args, ", "), x.Body)
}

type Concatenation struct {
	Exp1 Expression
	Exp2 Expression
}

func (x *Concatenation) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.Concatenate(env, x.Exp1.Evaluate(env, scope, warnings), x.Exp2.Evaluate(env, scope, warnings))
}

func (x *Concatenation) Visit(v func(Expression)) {
	x.Exp1.Visit(v)
	x.Exp2.Visit(v)
	v(x)
}

func (x *Concatenation) String() string {
	return fmt.Sprintf("%s & %s", x.Exp1.String(), x.Exp2.String())
}

type Addition struct {
	Exp1 Expression
	Exp2 Expression
}

func (x *Addition) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.Add(env, x.Exp1.Evaluate(env, scope, warnings), x.Exp2.Evaluate(env, scope, warnings))
}

func (x *Addition) Visit(v func(Expression)) {
	x.Exp1.Visit(v)
	x.Exp2.Visit(v)
	v(x)
}

func (x *Addition) String() string {
	return fmt.Sprintf("%s + %s", x.Exp1.String(), x.Exp2.String())
}

type Subtraction struct {
	Exp1 Expression
	Exp2 Expression
}

func (x *Subtraction) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.Subtract(env, x.Exp1.Evaluate(env, scope, warnings), x.Exp2.Evaluate(env, scope, warnings))
}

func (x *Subtraction) Visit(v func(Expression)) {
	x.Exp1.Visit(v)
	x.Exp2.Visit(v)
	v(x)
}

func (x *Subtraction) String() string {
	return fmt.Sprintf("%s - %s", x.Exp1.String(), x.Exp2.String())
}

type Multiplication struct {
	Exp1 Expression
	Exp2 Expression
}

func (x *Multiplication) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.Multiply(env, x.Exp1.Evaluate(env, scope, warnings), x.Exp2.Evaluate(env, scope, warnings))
}

func (x *Multiplication) Visit(v func(Expression)) {
	x.Exp1.Visit(v)
	x.Exp2.Visit(v)
	v(x)
}

func (x *Multiplication) String() string {
	return fmt.Sprintf("%s * %s", x.Exp1.String(), x.Exp2.String())
}

type Division struct {
	Exp1 Expression
	Exp2 Expression
}

func (x *Division) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.Divide(env, x.Exp1.Evaluate(env, scope, warnings), x.Exp2.Evaluate(env, scope, warnings))
}

func (x *Division) Visit(v func(Expression)) {
	x.Exp1.Visit(v)
	x.Exp2.Visit(v)
	v(x)
}

func (x *Division) String() string {
	return fmt.Sprintf("%s / %s", x.Exp1.String(), x.Exp2.String())
}

type Exponent struct {
	Expression Expression
	Exponent   Expression
}

func (x *Exponent) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.Exponent(env, x.Expression.Evaluate(env, scope, warnings), x.Exponent.Evaluate(env, scope, warnings))
}

func (x *Exponent) Visit(v func(Expression)) {
	x.Expression.Visit(v)
	x.Exponent.Visit(v)
	v(x)
}

func (x *Exponent) String() string {
	return fmt.Sprintf("%s ^ %s", x.Expression.String(), x.Exponent.String())
}

type Negation struct {
	Exp Expression
}

func (x *Negation) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.Negate(env, x.Exp.Evaluate(env, scope, warnings))
}

func (x *Negation) Visit(v func(Expression)) {
	x.Exp.Visit(v)
	v(x)
}

func (x *Negation) String() string {
	return fmt.Sprintf("-%s", x.Exp.String())
}

type Equality struct {
	Exp1 Expression
	Exp2 Expression
}

func (x *Equality) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.Equal(env, x.Exp1.Evaluate(env, scope, warnings), x.Exp2.Evaluate(env, scope, warnings))
}

func (x *Equality) Visit(v func(Expression)) {
	x.Exp1.Visit(v)
	x.Exp2.Visit(v)
	v(x)
}

func (x *Equality) String() string {
	return fmt.Sprintf("%s = %s", x.Exp1.String(), x.Exp2.String())
}

type InEquality struct {
	Exp1 Expression
	Exp2 Expression
}

func (x *InEquality) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.NotEqual(env, x.Exp1.Evaluate(env, scope, warnings), x.Exp2.Evaluate(env, scope, warnings))
}

func (x *InEquality) Visit(v func(Expression)) {
	x.Exp1.Visit(v)
	x.Exp2.Visit(v)
	v(x)
}

func (x *InEquality) String() string {
	return fmt.Sprintf("%s != %s", x.Exp1.String(), x.Exp2.String())
}

type LessThan struct {
	Exp1 Expression
	Exp2 Expression
}

func (x *LessThan) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.LessThan(env, x.Exp1.Evaluate(env, scope, warnings), x.Exp2.Evaluate(env, scope, warnings))
}

func (x *LessThan) Visit(v func(Expression)) {
	x.Exp1.Visit(v)
	x.Exp2.Visit(v)
	v(x)
}

func (x *LessThan) String() string {
	return fmt.Sprintf("%s < %s", x.Exp1.String(), x.Exp2.String())
}

type LessThanOrEqual struct {
	Exp1 Expression
	Exp2 Expression
}

func (x *LessThanOrEqual) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.LessThanOrEqual(env, x.Exp1.Evaluate(env, scope, warnings), x.Exp2.Evaluate(env, scope, warnings))
}

func (x *LessThanOrEqual) Visit(v func(Expression)) {
	x.Exp1.Visit(v)
	x.Exp2.Visit(v)
	v(x)
}

func (x *LessThanOrEqual) String() string {
	return fmt.Sprintf("%s <= %s", x.Exp1.String(), x.Exp2.String())
}

type GreaterThan struct {
	Exp1 Expression
	Exp2 Expression
}

func (x *GreaterThan) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.GreaterThan(env, x.Exp1.Evaluate(env, scope, warnings), x.Exp2.Evaluate(env, scope, warnings))
}

func (x *GreaterThan) Visit(v func(Expression)) {
	x.Exp1.Visit(v)
	x.Exp2.Visit(v)
	v(x)
}

func (x *GreaterThan) String() string {
	return fmt.Sprintf("%s > %s", x.Exp1.String(), x.Exp2.String())
}

type GreaterThanOrEqual struct {
	Exp1 Expression
	Exp2 Expression
}

func (x *GreaterThanOrEqual) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.GreaterThanOrEqual(env, x.Exp1.Evaluate(env, scope, warnings), x.Exp2.Evaluate(env, scope, warnings))
}

func (x *GreaterThanOrEqual) Visit(v func(Expression)) {
	x.Exp1.Visit(v)
	x.Exp2.Visit(v)
	v(x)
}

func (x *GreaterThanOrEqual) String() string {
	return fmt.Sprintf("%s >= %s", x.Exp1.String(), x.Exp2.String())
}

type Parentheses struct {
	Exp Expression
}

func (x *Parentheses) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return x.Exp.Evaluate(env, scope, warnings)
}

func (x *Parentheses) Visit(v func(Expression)) {
	x.Exp.Visit(v)
	v(x)
}

func (x *Parentheses) String() string {
	return fmt.Sprintf("(%s)", x.Exp.String())
}

type TextLiteral struct {
	Value *types.XText
}

func (x *TextLiteral) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return x.Value
}

func (x *TextLiteral) Visit(v func(Expression)) {
	v(x)
}

func (x *TextLiteral) String() string {
	return x.Value.Describe()
}

// NumberLiteral is a literal number like 123 or 1.5
type NumberLiteral struct {
	Value *types.XNumber
}

func (x *NumberLiteral) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return x.Value
}

func (x *NumberLiteral) Visit(v func(Expression)) {
	v(x)
}

func (x *NumberLiteral) String() string {
	return x.Value.Describe()
}

// ErrorLiteral is a literal error value
type ErrorLiteral struct {
	Err *types.XError
}

func (x *ErrorLiteral) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return x.Err
}

func (x *ErrorLiteral) Visit(v func(Expression)) {
	v(x)
}

func (x *ErrorLiteral) String() string {
	return "ERROR"
}

// BooleanLiteral is a literal bool
type BooleanLiteral struct {
	Value *types.XBoolean
}

func (x *BooleanLiteral) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return x.Value
}

func (x *BooleanLiteral) Visit(v func(Expression)) {
	v(x)
}

func (x *BooleanLiteral) String() string {
	return x.Value.Describe()
}

type NullLiteral struct{}

func (x *NullLiteral) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return nil
}

func (x *NullLiteral) Visit(v func(Expression)) {
	v(x)
}

func (x *NullLiteral) String() string {
	return "null"
}

func resolveLookup(env envs.Environment, container types.XValue, lookup types.XValue, dotNotation bool, warnings *Warnings) types.XValue {
	array, isArray := container.(*types.XArray)
	object, isObject := container.(*types.XObject)
	var resolved types.XValue

	if isArray && array != nil {
		// if left-hand side is an array, then this is an index
		index, xerr := types.ToInteger(env, lookup)
		if xerr != nil {
			return xerr
		}

		if index >= array.Count() || index < -array.Count() {
			return types.NewXErrorf("index %d out of range for %d items", index, array.Count())
		}
		if index < 0 {
			index += array.Count()
		}

		resolved = array.Get(index)

	} else if isObject && object != nil {
		// if left-hand side is an object, then this is a property lookup
		property, xerr := types.ToXText(env, lookup)
		if xerr != nil {
			return xerr
		}

		value, exists := object.Get(property.Native())

		// [] notation doesn't error for non-existent properties, . does
		if !exists && dotNotation {
			return types.NewXErrorf("%s has no property '%s'", types.Describe(container), property.Native())
		}

		resolved = value

	} else {
		return types.NewXErrorf("%s doesn't support lookups", types.Describe(container))
	}

	if !types.IsNil(resolved) && resolved.Deprecated() != "" {
		warnings.deprecatedContext(resolved)
	}

	return resolved
}
