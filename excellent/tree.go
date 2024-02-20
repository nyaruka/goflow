package excellent

import (
	"fmt"
	"slices"
	"strings"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/operators"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
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
	String() string
}

// ContextReference is an identifier which is a function name or root variable in the context
type ContextReference struct {
	name string
}

func (x *ContextReference) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	value, exists := scope.Get(x.name)
	if !exists {
		return types.NewXErrorf("context has no property '%s'", x.name)
	}

	if !utils.IsNil(value) && value.Deprecated() != "" {
		warnings.deprecatedContext(value)
	}

	return value
}

func (x *ContextReference) String() string {
	return strings.ToLower(x.name)
}

type DotLookup struct {
	container Expression
	lookup    string
}

func (x *DotLookup) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	containerVal := x.container.Evaluate(env, scope, warnings)
	if types.IsXError(containerVal) {
		return containerVal
	}

	return resolveLookup(env, containerVal, types.NewXText(x.lookup), true, warnings)
}

func (x *DotLookup) String() string {
	return fmt.Sprintf("%s.%s", x.container.String(), x.lookup)
}

type ArrayLookup struct {
	container Expression
	lookup    Expression
}

func (x *ArrayLookup) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	containerVal := x.container.Evaluate(env, scope, warnings)
	if types.IsXError(containerVal) {
		return containerVal
	}

	lookupVal := x.lookup.Evaluate(env, scope, warnings)
	if types.IsXError(lookupVal) {
		return lookupVal
	}

	return resolveLookup(env, containerVal, lookupVal, false, warnings)
}

func (x *ArrayLookup) String() string {
	return fmt.Sprintf("%s[%s]", x.container.String(), x.lookup.String())
}

type FunctionCall struct {
	function Expression
	params   []Expression
}

func (x *FunctionCall) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	funcVal := x.function.Evaluate(env, scope, warnings)
	if types.IsXError(funcVal) {
		return funcVal
	}

	asFunction, isFunction := funcVal.(*types.XFunction)
	if !isFunction {
		return types.NewXErrorf("%s is not a function", x.function.String())
	}

	params := make([]types.XValue, len(x.params))
	for i := range x.params {
		params[i] = x.params[i].Evaluate(env, scope, warnings)
	}

	return asFunction.Call(env, params)
}

func (x *FunctionCall) String() string {
	params := make([]string, len(x.params))
	for i := range x.params {
		params[i] = x.params[i].String()
	}

	return fmt.Sprintf("%s(%s)", x.function.String(), strings.Join(params, ", "))
}

type AnonFunction struct {
	args []string
	body Expression
}

func (x *AnonFunction) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	// create an XFunction which wraps our body expression
	fn := func(env envs.Environment, args ...types.XValue) types.XValue {
		// create new context that includes the args
		argsMap := make(map[string]types.XValue, len(x.args))
		for i := range x.args {
			argsMap[x.args[i]] = args[i]
		}
		childScope := NewScope(types.NewXObject(argsMap), scope)

		return x.body.Evaluate(env, childScope, warnings)
	}

	return types.NewXFunction("", functions.NumArgsCheck(len(x.args), fn))
}

func (x *AnonFunction) String() string {
	return fmt.Sprintf("(%s) => %s", strings.Join(x.args, ", "), x.body)
}

type Concatenation struct {
	exp1 Expression
	exp2 Expression
}

func (x *Concatenation) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.Concatenate(env, x.exp1.Evaluate(env, scope, warnings), x.exp2.Evaluate(env, scope, warnings))
}

func (x *Concatenation) String() string {
	return fmt.Sprintf("%s & %s", x.exp1.String(), x.exp2.String())
}

type Addition struct {
	exp1 Expression
	exp2 Expression
}

func (x *Addition) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.Add(env, x.exp1.Evaluate(env, scope, warnings), x.exp2.Evaluate(env, scope, warnings))
}

func (x *Addition) String() string {
	return fmt.Sprintf("%s + %s", x.exp1.String(), x.exp2.String())
}

type Subtraction struct {
	exp1 Expression
	exp2 Expression
}

func (x *Subtraction) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.Subtract(env, x.exp1.Evaluate(env, scope, warnings), x.exp2.Evaluate(env, scope, warnings))
}

func (x *Subtraction) String() string {
	return fmt.Sprintf("%s - %s", x.exp1.String(), x.exp2.String())
}

type Multiplication struct {
	exp1 Expression
	exp2 Expression
}

func (x *Multiplication) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.Multiply(env, x.exp1.Evaluate(env, scope, warnings), x.exp2.Evaluate(env, scope, warnings))
}

func (x *Multiplication) String() string {
	return fmt.Sprintf("%s * %s", x.exp1.String(), x.exp2.String())
}

type Division struct {
	exp1 Expression
	exp2 Expression
}

func (x *Division) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.Divide(env, x.exp1.Evaluate(env, scope, warnings), x.exp2.Evaluate(env, scope, warnings))
}

func (x *Division) String() string {
	return fmt.Sprintf("%s / %s", x.exp1.String(), x.exp2.String())
}

type Exponent struct {
	expression Expression
	exponent   Expression
}

func (x *Exponent) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.Exponent(env, x.expression.Evaluate(env, scope, warnings), x.exponent.Evaluate(env, scope, warnings))
}

func (x *Exponent) String() string {
	return fmt.Sprintf("%s ^ %s", x.expression.String(), x.exponent.String())
}

type Negation struct {
	exp Expression
}

func (x *Negation) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.Negate(env, x.exp.Evaluate(env, scope, warnings))
}

func (x *Negation) String() string {
	return fmt.Sprintf("-%s", x.exp.String())
}

type Equality struct {
	exp1 Expression
	exp2 Expression
}

func (x *Equality) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.Equal(env, x.exp1.Evaluate(env, scope, warnings), x.exp2.Evaluate(env, scope, warnings))
}

func (x *Equality) String() string {
	return fmt.Sprintf("%s = %s", x.exp1.String(), x.exp2.String())
}

type InEquality struct {
	exp1 Expression
	exp2 Expression
}

func (x *InEquality) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.NotEqual(env, x.exp1.Evaluate(env, scope, warnings), x.exp2.Evaluate(env, scope, warnings))
}

func (x *InEquality) String() string {
	return fmt.Sprintf("%s != %s", x.exp1.String(), x.exp2.String())
}

type LessThan struct {
	exp1 Expression
	exp2 Expression
}

func (x *LessThan) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.LessThan(env, x.exp1.Evaluate(env, scope, warnings), x.exp2.Evaluate(env, scope, warnings))
}

func (x *LessThan) String() string {
	return fmt.Sprintf("%s < %s", x.exp1.String(), x.exp2.String())
}

type LessThanOrEqual struct {
	exp1 Expression
	exp2 Expression
}

func (x *LessThanOrEqual) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.LessThanOrEqual(env, x.exp1.Evaluate(env, scope, warnings), x.exp2.Evaluate(env, scope, warnings))
}

func (x *LessThanOrEqual) String() string {
	return fmt.Sprintf("%s <= %s", x.exp1.String(), x.exp2.String())
}

type GreaterThan struct {
	exp1 Expression
	exp2 Expression
}

func (x *GreaterThan) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.GreaterThan(env, x.exp1.Evaluate(env, scope, warnings), x.exp2.Evaluate(env, scope, warnings))
}

func (x *GreaterThan) String() string {
	return fmt.Sprintf("%s > %s", x.exp1.String(), x.exp2.String())
}

type GreaterThanOrEqual struct {
	exp1 Expression
	exp2 Expression
}

func (x *GreaterThanOrEqual) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return operators.GreaterThanOrEqual(env, x.exp1.Evaluate(env, scope, warnings), x.exp2.Evaluate(env, scope, warnings))
}

func (x *GreaterThanOrEqual) String() string {
	return fmt.Sprintf("%s >= %s", x.exp1.String(), x.exp2.String())
}

type Parentheses struct {
	exp Expression
}

func (x *Parentheses) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return x.exp.Evaluate(env, scope, warnings)
}

func (x *Parentheses) String() string {
	return fmt.Sprintf("(%s)", x.exp.String())
}

type TextLiteral struct {
	val *types.XText
}

func (x *TextLiteral) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return x.val
}

func (x *TextLiteral) String() string {
	return x.val.Describe()
}

// NumberLiteral is a literal number like 123 or 1.5
type NumberLiteral struct {
	val *types.XNumber
}

func (x *NumberLiteral) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return x.val
}

func (x *NumberLiteral) String() string {
	return x.val.Describe()
}

// BooleanLiteral is a literal bool
type BooleanLiteral struct {
	val *types.XBoolean
}

func (x *BooleanLiteral) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return x.val
}

func (x *BooleanLiteral) String() string {
	return x.val.Describe()
}

type NullLiteral struct{}

func (x *NullLiteral) Evaluate(env envs.Environment, scope *Scope, warnings *Warnings) types.XValue {
	return nil
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

	if !utils.IsNil(resolved) && resolved.Deprecated() != "" {
		warnings.deprecatedContext(resolved)
	}

	return resolved
}
