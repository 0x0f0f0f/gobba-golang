package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// TODO document
func (c Context) CheckAgainst(expr ast.Expression, ty ast.TypeValue) (Context, *TypeError) {
	if !c.IsWellFormed(ty) {
		return c, c.malformedError(ty)
	}

	switch vexpr := expr.(type) {
	case *ast.UnitLiteral:
		// Rule 1l
		if _, ok := ty.(*ast.UnitType); ok {
			return c, nil
		}
	case *ast.BoolLiteral:
		if _, ok := ty.(*ast.BoolType); ok {
			return c, nil
		}
	case *ast.FloatLiteral:
		if _, ok := ty.(*ast.FloatType); ok {
			return c, nil
		}
	case *ast.ComplexLiteral:
		if _, ok := ty.(*ast.ComplexType); ok {
			return c, nil
		}
	case *ast.IntegerLiteral:
		if _, ok := ty.(*ast.IntegerType); ok {
			return c, nil
		}
	case *ast.StringLiteral:
		if _, ok := ty.(*ast.StringType); ok {
			return c, nil
		}
	case *ast.RuneLiteral:
		if _, ok := ty.(*ast.RuneType); ok {
			return c, nil
		}
	case *ast.FunctionLiteral:
		// Rule ->l
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
		uv := &UniversalVariable{Identifier: fty.Identifier}
		nc := c.InsertHead(uv)
		subcheck, err := nc.CheckAgainst(expr, fty.Type)
		if err != nil {
			return c, err
		}
		return subcheck.Drop(uv), nil
	} else {
		// Rule sub
		a, theta, err := c.SynthesizesTo(expr)
		if err != nil {
			return c, err
		}
		theta.Subtype(theta.Apply(a), theta.Apply(ty))
	}

	return c, nil

}
