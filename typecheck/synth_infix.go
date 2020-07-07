package typecheck

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
)

func (c Context) synthComparison(leftt, rightt ast.TypeValue) (ast.TypeValue, Context, error) {
	// TODO, when interfaces are implemented, check if exp
	// implements the `comparable` interface
	theta, err := c.Subtype(leftt, rightt)
	if err != nil {
		return nil, c, c.expectedSameTypeComparison(leftt, rightt)
	}

	theta1, err := theta.Subtype(rightt, leftt)
	if err != nil {
		return nil, c, c.expectedSameTypeComparison(leftt, rightt)
	}

	return &ast.BoolType{}, theta1, nil

}

func (c Context) synthInfixExpr(exp *ast.InfixExpression) (ast.TypeValue, Context, error) {
	// Synthesize types for operands
	leftt, gamma, err := c.SynthesizesTo(exp.Left)
	if err != nil {
		return nil, c, err
	}
	rightt, gamma1, err := gamma.SynthesizesTo(exp.Right)
	if err != nil {
		return nil, c, err
	}

	switch exp.Operator {
	// ======================================================================
	// Comparison Operators
	// ======================================================================
	case "=":
		return gamma1.synthComparison(leftt, rightt)
	case ">":
		return gamma1.synthComparison(leftt, rightt)
	case "<":
		return gamma1.synthComparison(leftt, rightt)
	case "<=":
		return gamma1.synthComparison(leftt, rightt)
	case ">=":
		return gamma1.synthComparison(leftt, rightt)

	case "+":
		// FIXME

		fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", leftt.FullString())

		// If the left operand is an existential variable

		// if leftext, ok := leftt.(*ast.ExistsType); ok {

		// }

		gamma1.debugRule("AAAAAAAAAAAAA" + leftt.FullString() + " " + rightt.FullString())

		theta, err := gamma1.Subtype(leftt, &ast.NumberType{})
		if err != nil {
			return nil, c, err
		}

		theta1, err := theta.Subtype(rightt, &ast.NumberType{})
		if err != nil {
			return nil, c, err
		}

		var delta Context
		leftapp := theta1.Apply(leftt)
		rightapp := theta1.Apply(rightt)

		theta1.debugRule("BBBBBBBBBBBB " + leftapp.FullString() + " " + rightapp.FullString())

		// Try to see if left <: right
		delta, err = theta1.Subtype(leftapp, rightapp)
		if err != nil {
			// Try the other way around
			delta, err = theta1.Subtype(rightapp, leftapp)
			if err != nil {
				return nil, c, err
			}
			// Rule ◦RSubL=>
			c.debugRule("◦RSubL=>")
			return leftt, delta, nil
		}
		// Rule ◦LSubR=>
		c.debugRule("◦LSubR=>")
		return rightt, delta, nil

	default:
		// TODO
		panic("Type synthesis Not yet implemented for expression " + exp.String())

	}

	return nil, c, c.synthError(exp)
}
