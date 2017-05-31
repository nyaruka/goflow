// Generated from excellent/gen/Excellent.g4 by ANTLR 4.7.

package gen // Excellent
import "github.com/antlr/antlr4/runtime/Go/antlr"

type BaseExcellentVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseExcellentVisitor) VisitParse(ctx *ParseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitDecimalLiteral(ctx *DecimalLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitDotLookup(ctx *DotLookupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitStringLiteral(ctx *StringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitFunctionCall(ctx *FunctionCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitTrue(ctx *TrueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitFalse(ctx *FalseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitArrayLookup(ctx *ArrayLookupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitContextReference(ctx *ContextReferenceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitParentheses(ctx *ParenthesesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitNegation(ctx *NegationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitComparison(ctx *ComparisonContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitConcatenation(ctx *ConcatenationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitMultiplicationOrDivision(ctx *MultiplicationOrDivisionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitAtomReference(ctx *AtomReferenceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitAdditionOrSubtraction(ctx *AdditionOrSubtractionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitEquality(ctx *EqualityContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitExponent(ctx *ExponentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitFnname(ctx *FnnameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseExcellentVisitor) VisitFunctionParameters(ctx *FunctionParametersContext) interface{} {
	return v.VisitChildren(ctx)
}
