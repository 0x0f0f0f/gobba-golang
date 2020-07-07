package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// TODO document
func (c Context) CheckAgainst(expr ast.Expression, ty ast.TypeValue) (Context, error) {
	c.debugSection("check", expr.String(), "<=", ty.FullString())
	if !c.IsWellFormed(ty) {
		return c, c.malformedError(ty)
	}

	switch vexpr := expr.(type) {
	case *ast.UnitLiteral:
		// Rule 1l
		c.debugRuleOut("1I")

		if _, ok := ty.(*ast.UnitType); ok {
			return c, nil
		}

	case *ast.BoolLiteral:
		// Rule booll
		c.debugRuleOut("boolI")

		if _, ok := ty.(*ast.BoolType); ok {
			return c, nil
		}
	case *ast.FloatLiteral:
		// Rule floatl
		c.debugRuleOut("floatI")

		if _, ok := ty.(*ast.FloatType); ok {
			return c, nil
		}
	case *ast.ComplexLiteral:
		// Rule complexl
		c.debugRuleOut("complexI")

		if _, ok := ty.(*ast.ComplexType); ok {
			return c, nil
		}
	case *ast.IntegerLiteral:
		// Rule intl
		c.debugRuleOut("intI")

		if _, ok := ty.(*ast.IntegerType); ok {
			return c, nil
		}
	case *ast.StringLiteral:
		// Rule stringl
		c.debugRuleOut("stringI")

		if _, ok := ty.(*ast.StringType); ok {
			return c, nil
		}
	case *ast.RuneLiteral:
		// Rule runel
		c.debugRuleOut("runeI")

		if _, ok := ty.(*ast.RuneType); ok {
			return c, nil
		}
	case *ast.FunctionLiteral:
		// Rule ->l
		c.debugRule("->l")

		if lty, ok := ty.(*ast.LambdaType); ok {
			typedvar := &TypeAnnotation{
				Identifier: vexpr.Param.Identifier,
				Value:      lty.Domain,
			}
			nc := c.InsertHead(typedvar)
			subcheck, err := nc.CheckAgainst(vexpr.Body, lty.Codomain)
			if err != nil {
				return c, err
			}
			return subcheck.Drop(typedvar), nil

		}

	}

	if fty, ok := ty.(*ast.ForAllType); ok {
		// Rule ∀l
		c.debugRule("∀I")

		uv := &UniversalVariable{Identifier: fty.Identifier}
		nc := c.InsertHead(uv)
		subcheck, err := nc.CheckAgainst(expr, fty.Type)
		if err != nil {
			c.debugRuleFail("∀I")
			return c, err
		}

		c.debugRuleOut("∀I")
		return subcheck.Drop(uv), nil
	}
	// Rule Sub
	c.debugRule("Sub")

	a, theta, err := c.SynthesizesTo(expr)
	if err != nil {
		c.debugRuleFail("Sub")

		return c, err
	}

	c.debugRuleOut("Sub")
	return theta.Subtype(theta.Apply(a), theta.Apply(ty))

}
