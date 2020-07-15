package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/token"
)

// This file contains definitions for the Algorithmic Subtyping rules

// Helper function
// func sameType

func (c Context) Subtype(a, b ast.TypeValue) (Context, error) {
	c.debugSection("subtype", a.FullString(), "<:", b.FullString())
	if !c.IsWellFormed(a) {
		return c, c.malformedError(a)
	}
	if !c.IsWellFormed(b) {
		return c, c.malformedError(b)
	}

	switch va := a.(type) {
	case *ast.TyUnit:
		// Rule <:Unit
		c.debugRule("<:Unit")

		if _, ok := b.(*ast.TyUnit); ok {
			return c, nil
		}
	case *ast.TyUnVar:
		switch vb := b.(type) {
		case *ast.TyUnVar:
			if va.Identifier == vb.Identifier {
				// Rule <:Var
				c.debugRule("<:Var")

				return c, nil
			}
			switch va.Identifier.Value {
			case token.TINT:
				switch vb.Identifier.Value {
				case token.TFLOAT:
					// Rule int<:float
					c.debugRule("int<:float")
					return c, nil
				case token.TCOMPLEX:
					// Rule int<:complex
					c.debugRule("int<:complex")
					return c, nil

				}
			case token.TFLOAT:
				switch vb.Identifier.Value {
				case token.TCOMPLEX:
					// Rule float<:complex
					c.debugRule("float<:complex")
					return c, nil
				}
			}

		}

		// =============================================================
		// Other Primitive Subtyping Rules
		// =============================================================

	case *ast.TyExVar:
		if vb, ok := b.(*ast.TyExVar); ok {
			if va.Identifier == vb.Identifier {
				// Rule <:Exvar
				c.debugRule("<:Exvar")

				return c, nil
			}
		}
		if !OccursIn(va.Identifier, b) {
			// Rule <:InstantiateL
			c.debugRule("<:InstantiateL")

			res := c.InstantiateL(va.Identifier, b)
			return res, nil
		}

	case *ast.TyLambda:
		switch vb := b.(type) {
		case *ast.TyLambda:
			// Rule <:->
			c.debugRule("<:->")

			theta, err := c.Subtype(va.Domain, vb.Domain)
			if err != nil {
				return c, err
			}
			return theta.Subtype(theta.Apply(va.Codomain),
				theta.Apply(vb.Codomain))
		}

	case *ast.TyForAll:
		// Rule <:∀L
		c.debugRule("<:∀L")

		r1 := ast.GenUID("α")
		marker := &Marker{r1}
		exv := &ExistentialVariable{r1, nil}
		ext := &ast.TyExVar{Identifier: r1}
		gamma := c.InsertHead(exv).InsertHead(marker)
		sub_a := Substitution(va.Type, ext, va.Identifier)
		delta, err := gamma.Subtype(sub_a, b)
		if err != nil {
			return c, err
		}
		return delta.Drop(marker), nil

	}

	if vb, ok := b.(*ast.TyExVar); ok {
		if !OccursIn(vb.Identifier, a) {
			// Rule <:InstantiateR
			c.debugRule("<:InstantiateR")

			res := c.InstantiateR(a, vb.Identifier)
			return res, nil
		}

	}

	if vb, ok := b.(*ast.TyForAll); ok {
		// Rule <:∀R
		c.debugRule("<:∀R")

		u := &UniversalVariable{vb.Identifier}
		theta := c.InsertHead(u)
		delta, err := theta.Subtype(a, vb)
		if err != nil {
			return c, err
		}
		return delta.Drop(u), nil

	}

	return c, c.subtypeError(a, b)

}
