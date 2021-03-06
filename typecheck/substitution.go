package typecheck

import (
	// "fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// Returns true if a type variable identifier occurs in a given type
func OccursIn(alpha ast.UniqueIdentifier, a ast.TypeValue) bool {
	switch va := a.(type) {
	case *ast.VariableType:
		return va.Identifier == alpha
	case *ast.ExistsType:
		return va.Identifier == alpha
	case *ast.LambdaType:
		return OccursIn(alpha, va.Domain) || OccursIn(alpha, va.Codomain)
	case *ast.ForAllType:
		return va.Identifier == alpha || OccursIn(alpha, va.Type)
	default:
		// Type variables do not occur in monotypes
		return false
	}

}

// TODO document
func Substitution(a, b ast.TypeValue, alpha ast.UniqueIdentifier) ast.TypeValue {
	switch va := a.(type) {
	case *ast.VariableType:
		if va.Identifier == alpha {
			return b
		} else {
			return a
		}
	case *ast.ExistsType:
		if va.Identifier == alpha {
			return b
		} else {
			return a
		}
	case *ast.ForAllType:
		if va.Identifier == alpha {
			return &ast.ForAllType{va.Identifier, b}
		} else {
			return &ast.ForAllType{
				Identifier: va.Identifier,
				Type:       Substitution(va.Type, b, alpha),
			}
		}
	case *ast.LambdaType:
		return &ast.LambdaType{
			Domain:   Substitution(va.Domain, b, alpha),
			Codomain: Substitution(va.Codomain, b, alpha),
		}
	default:
		return a

	}

}

// Apply a context as a substitution to a value
func (c *Context) Apply(a ast.TypeValue) ast.TypeValue {
	switch va := a.(type) {
	case *ast.ExistsType:
		tau := c.GetSolvedVariable(va.Identifier)
		if tau == nil {
			c.debugSection("apply", a.FullString(), "=", a.FullString())
			return a
		} else {
			ret := c.Apply(*tau)
			c.debugSection("apply", a.FullString(), "=", ret.FullString())
			return ret
		}
	case *ast.LambdaType:
		ret := &ast.LambdaType{
			Domain:   c.Apply(va.Domain),
			Codomain: c.Apply(va.Codomain),
		}
		c.debugSection("apply", a.FullString(), "=", ret.FullString())
		return ret
	case *ast.ForAllType:
		ret := &ast.ForAllType{
			Identifier: va.Identifier,
			Type:       c.Apply(va.Type),
		}
		c.debugSection("apply", a.FullString(), "=", ret.FullString())
		return ret
	}
	c.debugSection("apply", a.FullString(), "=", a.FullString())
	return a
}
