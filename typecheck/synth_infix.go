package typecheck

import (
	// "fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/token"
)

func (c Context) synthComparison(leftt, rightt ast.TypeValue) (ast.TypeValue, Context, error) {
	// TODO, when interfaces are implemented, check if exp
	// implements the `comparable` interface
	gamma1, err := c.Subtype(leftt, rightt)
	if err != nil {
		return nil, c, c.expectedSameTypeComparison(leftt, rightt)
	}

	gamma2, err := gamma1.Subtype(rightt, leftt)
	if err != nil {
		return nil, c, c.expectedSameTypeComparison(leftt, rightt)
	}

	return ast.TBOOL, gamma2, nil

}

func (c Context) synthInfixExpr(exp *ast.InfixExpression) (ast.TypeValue, Context, error) {
	// Synthesize types for operands
	leftt, gamma1, err := c.SynthesizesTo(exp.Left)
	if err != nil {
		return nil, c, err
	}
	rightt, theta, err := gamma1.SynthesizesTo(exp.Right)
	if err != nil {
		return nil, c, err
	}

	if resultt, ok := ast.OperatorTypes[exp.Operator]; ok {
		psi, err := theta.Subtype(leftt, resultt.Left)
		if err != nil {
			return nil, c, err
		}

		delta, err := psi.Subtype(rightt, resultt.Right)
		if err != nil {
			return nil, c, err
		}

		return resultt.Result, delta, err

	}

	switch exp.Operator {
	// ======================================================================
	// Comparison Operators
	// ======================================================================
	case token.EQUALS:
		return gamma1.synthComparison(leftt, rightt)
	case token.GREATER:
		return gamma1.synthComparison(leftt, rightt)
	case token.LESS:
		return gamma1.synthComparison(leftt, rightt)
	case token.LESSEQ:
		return gamma1.synthComparison(leftt, rightt)
	case token.GREATEREQ:
		return gamma1.synthComparison(leftt, rightt)
	default:
		// TODO
		return nil, c, c.synthError(exp)
	}

	return nil, c, c.synthError(exp)
}
