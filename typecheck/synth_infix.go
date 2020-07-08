package typecheck

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
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

	return ast.NewVariableType("bool"), theta1, nil

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

		theta.debugRule("AAAAAAAAAAAAA" + leftt.FullString() + " " + rightt.FullString())

		theta1, err := theta.Subtype(leftt, ast.NewVariableType("number"))
		if err != nil {
			return nil, c, err
		}

		delta, err := theta1.Subtype(rightt, ast.NewVariableType("number"))
		if err != nil {
			return nil, c, err
		}

		leftapp := delta.Apply(leftt)
		rightapp := delta.Apply(rightt)

		delta.debugRule("BBBBBBBBBBBB " + leftapp.FullString() + " " + rightapp.FullString())

		// Try to see if left <: right
		_, err = delta.Subtype(leftapp, rightapp)
		if err != nil {
			// Try the other way around
			_, err = delta.Subtype(rightapp, leftapp)
			if err != nil {
				return nil, c, err
			}
			// Rule ◦RSubL=>
			c.debugRule("◦RSubL=>")
			return leftt, theta, nil
		}
		// Rule ◦LSubR=>
		c.debugRule("◦LSubR=>")
		return rightt, theta, nil

	default:
		// TODO
		panic("Type synthesis Not yet implemented for expression " + exp.String())

	}

	return nil, c, c.synthError(exp)
}
