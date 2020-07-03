package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// This file contains definitions for the Algorithmic Subtyping rules

// Helper function
// func sameType

// TODO precedence is fucked up
func (c *Context) Subtype(a, b ast.TypeValue) (*Context, *TypeError) {
	if !c.IsWellFormed(a) {
		return c, c.malformedError(a)
	}
	if !c.IsWellFormed(b) {
		return c, c.malformedError(b)
	}

	switch va := a.(type) {
	case *ast.VariableType:
		switch vb := b.(type) {
		case *ast.VariableType:
			// Rule <:Var
			if va.Identifier == vb.Identifier {
				return c, nil
			}
		case *ast.ForAllType:
			// Rule <:∀R
			return c.ruleSubtypeForAllRight(a, vb)
		case *ast.ExistsType:
			// TODO Rule <:InstantiateR

		}

	case *ast.ExistsType:
		switch vb := b.(type) {
		case *ast.ExistsType:
			if va.Identifier == vb.Identifier {
				// Rule <: Exvar
				return c, nil
			} else {
				// TODO Rule <:InstantiateL
			}
		case *ast.ForAllType:
			// Rule <:∀R
			return c.ruleSubtypeForAllRight(a, vb)
		default:
			// TODO Rule <:InstantiateR
		}

	case *ast.LambdaType:
		switch vb := b.(type) {
		case *ast.LambdaType:
			// Rule <:->
			theta, err := c.Subtype(va.Domain, vb.Domain)
			if err != nil {
				return nil, err
			}
			// TODO apply_context
			return theta.Subtype(theta.ApplyContext(va.Codomain),
				theta.ApplyContext(vb.Codomain))
		case *ast.ForAllType:
			// Rule <:∀R
			return c.ruleSubtypeForAllRight(a, vb)
		case *ast.ExistsType:
			// TODO Rule <:InstantiateR

		}

	case *ast.ForAllType:
		// Rule <:∀L
		return c.ruleSubtypeForAllLeft(va, b)
	default:
		return nil, c.subtypeError(a, b)
	}

}

// Rule <:∀R
func (c *Context) ruleSubtypeForAllRight(a ast.TypeValue, b *ast.ForAllType) (*Context, *TypeError) {
	u := &UniversalVariable{b.Identifier}
	theta := c.InsertHead(u)
	delta, err := theta.Subtype(a, b)
	if err != nil {
		return nil, err
	}
	return delta.Drop(u), nil
}

// Rule <:∀L
func (c *Context) ruleSubtypeForAllLeft(a *ast.ForAllType, b ast.TypeValue) (*Context, *TypeError) {
	r1 := ast.GenUID("alpha")
	marker := &Marker{r1}
	exv := &ExistentialVariable{r1, nil}
	ext := &ast.ExistsType{r1}
	gamma := c.InsertHead(exv).InsertHead(marker)
	sub_a := Substitution(a.Type, ext, a.Identifier)
	delta, err := gamma.Subtype(sub_a, b)
	if err != nil {
		return nil, err
	}
	return delta.Drop(marker), nil
}
