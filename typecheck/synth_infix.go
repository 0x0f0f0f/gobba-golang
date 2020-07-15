package typecheck

import (
	// "fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/token"
)

func (Γ Context) synthComparison(leftt, rightt ast.TypeValue) (ast.TypeValue, Context, error) {
	// TODO, when interfaces are implemented, check if exp
	// implements the `comparable` interface
	Γ1, err := Γ.Subtype(leftt, rightt)
	if err != nil {
		return nil, Γ, Γ.expectedSameTypeComparison(leftt, rightt)
	}

	Γ2, err := Γ1.Subtype(rightt, leftt)
	if err != nil {
		return nil, Γ, Γ.expectedSameTypeComparison(leftt, rightt)
	}

	return ast.TBOOL, Γ2, nil

}

func (Γ Context) synthInfixExpr(exp *ast.ExprInfix) (ast.TypeValue, Context, error) {
	// Synthesize types for operands
	leftt, Γ1, err := Γ.SynthesizesTo(exp.Left)
	if err != nil {
		return nil, Γ, err
	}
	rightt, Θ, err := Γ1.SynthesizesTo(exp.Right)
	if err != nil {
		return nil, Γ, err
	}

	if resultt, ok := ast.InfixOperatorTypes[exp.Operator]; ok {
		Θ1, err := Θ.Subtype(leftt, resultt.Left)
		if err != nil {
			return nil, Γ, err
		}

		Δ, err := Θ1.Subtype(rightt, resultt.Right)
		if err != nil {
			return nil, Γ, err
		}

		return resultt.Result, Δ, err

	}

	switch exp.Operator {
	// ======================================================================
	// Comparison Operators
	// ======================================================================
	case token.EQUALS:
		return Γ1.synthComparison(leftt, rightt)
	case token.GREATER:
		return Γ1.synthComparison(leftt, rightt)
	case token.LESS:
		return Γ1.synthComparison(leftt, rightt)
	case token.LESSEQ:
		return Γ1.synthComparison(leftt, rightt)
	case token.GREATEREQ:
		return Γ1.synthComparison(leftt, rightt)
	}

	return nil, Γ, Γ.synthError(exp)
}

func (Γ Context) synthPrefixExpr(exp *ast.ExprPrefix) (ast.TypeValue, Context, error) {
	if resultt, ok := ast.PrefixOperatorTypes[exp.Operator]; ok {
		Δ, err := Γ.CheckAgainst(exp.Right, resultt.Right)
		if err != nil {
			return nil, Γ, err
		}
		return resultt.Result, Δ, err
	}
	return nil, Γ, Γ.synthError(exp)
}
