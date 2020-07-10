package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
	// "github.com/0x0f0f0f/gobba-golang/token"
)

// Helper functions that checks if a VariableType has a given
// identifier name
func (c Context) checkVariable(ty ast.TypeValue, name string) bool {
	vty, ok := ty.(*ast.VariableType)
	return ok && vty.Identifier.Value == name
}

// TODO document
func (c Context) CheckAgainst(expr ast.Expression, ty ast.TypeValue) (Context, error) {
	c.debugSection("check", expr.String(), "<=", ty.FullString())
	if !c.IsWellFormed(ty) {
		return c, c.malformedError(ty)
	}

	switch vexpr := expr.(type) {
	case *ast.UnitLiteral:
		// Rule 1l
		c.debugRule("1I")

		if _, ok := ty.(*ast.UnitType); ok {
			c.debugRuleOut("1I")
			return c, nil
		}

	case *ast.BoolLiteral:
		// Rule booll
		c.debugRule("boolI")

		if c.checkVariable(ty, "bool") {
			c.debugRuleOut("boolI")
			return c, nil
		}
	case *ast.FloatLiteral:
		// Rule floatl
		c.debugRule("floatI")

		if c.checkVariable(ty, "float") {
			c.debugRuleOut("floatI")
			return c, nil
		}

	case *ast.ComplexLiteral:
		// Rule complexl
		c.debugRule("complexI")

		if c.checkVariable(ty, "complex") {
			c.debugRuleOut("complexI")
			return c, nil
		}

	case *ast.IntegerLiteral:
		// Rule intl
		c.debugRule("intI")

		if c.checkVariable(ty, "int") {
			c.debugRuleOut("intI")
			return c, nil
		}

	case *ast.StringLiteral:
		// Rule stringl
		c.debugRuleOut("stringI")

		if c.checkVariable(ty, "string") {
			c.debugRuleOut("stringI")
			return c, nil
		}
	case *ast.RuneLiteral:
		// Rule runel
		c.debugRuleOut("runeI")

		if c.checkVariable(ty, "rune") {
			c.debugRuleOut("runeI")
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
				c.debugRuleFail("->l")
				return c, err
			}
			// outc := subcheck.Drop(typedvar)
			c.debugRuleOut("->l")
			return subcheck, nil
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
