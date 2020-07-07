package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
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
	case *ast.UnitType:
		// Rule <:Unit
		c.debugRule("<:Unit")

		if _, ok := b.(*ast.UnitType); ok {
			return c, nil
		}
	case *ast.BoolType:
		// Rule <:bool
		c.debugRule("<:bool")

		if _, ok := b.(*ast.BoolType); ok {
			return c, nil
		}

		// =============================================================
		// Numerical Subtyping Rules
		// =============================================================

	case *ast.IntegerType:
		switch b.(type) {
		case *ast.IntegerType:
			// Rule <:int
			c.debugRule("<:int")
			return c, nil
		case *ast.FloatType:
			// Rule int<:float
			c.debugRule("int<:float")
			return c, nil
		case *ast.ComplexType:
			// Rule int<:complex
			c.debugRule("int<:complex")
			return c, nil
		case *ast.NumberType:
			// Rule int<:number
			c.debugRule("int<:number")
			return c, nil

		}
	case *ast.FloatType:
		switch b.(type) {
		case *ast.FloatType:
			// Rule <:float
			c.debugRule("<:float")
			return c, nil
		case *ast.ComplexType:
			// Rule float<:complex
			c.debugRule("float<:complex")
			return c, nil
		case *ast.NumberType:
			// Rule float<:number
			c.debugRule("float<:number")
			return c, nil
		}
	case *ast.ComplexType:
		switch b.(type) {
		case *ast.ComplexType:
			// Rule <:complex
			c.debugRule("<:complex")
			return c, nil
		case *ast.NumberType:
			// Rule complex<:number
			c.debugRule("complex<:number")
			return c, nil

		}
	case *ast.NumberType:
		switch b.(type) {
		case *ast.NumberType:
			// Rule <:number
			c.debugRule("<:number")
			return c, nil
		}

		// =============================================================
		// Other Primitive Subtyping Rules
		// =============================================================

	case *ast.StringType:
		switch b.(type) {
		case *ast.StringType:
			// Rule <:string
			c.debugRule("<:string")
			return c, nil
		}

	case *ast.VariableType:
		if vb, ok := b.(*ast.VariableType); ok {
			// Rule <:Var
			c.debugRule("<:Var")

			if va.Identifier == vb.Identifier {
				return c, nil
			}
		}

	case *ast.ExistsType:
		if vb, ok := b.(*ast.ExistsType); ok {
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

	case *ast.LambdaType:
		if vb, ok := b.(*ast.LambdaType); ok {
			// Rule <:->
			c.debugRule("<:->")

			theta, err := c.Subtype(va.Domain, vb.Domain)
			if err != nil {
				return c, err
			}
			return theta.Subtype(theta.Apply(va.Codomain),
				theta.Apply(vb.Codomain))
		}

	case *ast.ForAllType:
		// Rule <:∀L
		c.debugRule("<:∀L")

		r1 := ast.GenUID("α")
		marker := &Marker{r1}
		exv := &ExistentialVariable{r1, nil}
		ext := &ast.ExistsType{Identifier: r1}
		gamma := c.InsertHead(exv).InsertHead(marker)
		sub_a := Substitution(va.Type, ext, va.Identifier)
		delta, err := gamma.Subtype(sub_a, b)
		if err != nil {
			return c, err
		}
		return delta.Drop(marker), nil

	}

	if vb, ok := b.(*ast.ExistsType); ok {
		if !OccursIn(vb.Identifier, a) {
			// Rule <:InstantiateR
			c.debugRule("<:InstantiateR")

			res := c.InstantiateR(a, vb.Identifier)
			return res, nil
		}

	}

	if vb, ok := b.(*ast.ForAllType); ok {
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
