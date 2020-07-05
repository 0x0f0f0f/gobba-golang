package typecheck

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// This file contains definitions for the Algorithmic Subtyping rules

// Helper function
// func sameType

func (c Context) Subtype(a, b ast.TypeValue) (Context, *TypeError) {
	fmt.Println("subtype", a, "<:", b)
	if !c.IsWellFormed(a) {
		return c, c.malformedError(a)
	}
	if !c.IsWellFormed(b) {
		return c, c.malformedError(b)
	}

	switch va := a.(type) {
	case *ast.UnitType:
		// Rule <:Unit
		fmt.Println("\tApplying rule <:Unit", c.String())
		if _, ok := b.(*ast.UnitType); ok {
			return c, nil
		}
	case *ast.IntegerType:
		// Rule <:Unit
		fmt.Println("\tApplying rule <:Int", c.String())
		if _, ok := b.(*ast.IntegerType); ok {
			return c, nil
		}

	case *ast.VariableType:
		if vb, ok := b.(*ast.VariableType); ok {
			// Rule <:Var
			fmt.Println("\tApplying rule <:Var", c.String())
			if va.Identifier == vb.Identifier {
				return c, nil
			}
		}

	case *ast.ExistsType:
		if vb, ok := b.(*ast.ExistsType); ok {
			if va.Identifier == vb.Identifier {
				// Rule <:Exvar
				fmt.Println("\tApplying rule <:Exvar", c.String())
				return c, nil
			} else if !OccursIn(va.Identifier, b) {
				// Rule <:InstantiateL
				fmt.Println("\tApplying rule <:InstantiateL", c.String())
				res := c.InstantiateL(va.Identifier, b)
				return res, nil
			}
		}

	case *ast.LambdaType:
		if vb, ok := b.(*ast.LambdaType); ok {
			// Rule <:->
			fmt.Println("\tApplying rule <:->", c.String())
			theta, err := c.Subtype(va.Domain, vb.Domain)
			if err != nil {
				return c, err
			}
			return theta.Subtype(theta.Apply(va.Codomain),
				theta.Apply(vb.Codomain))
		}

	case *ast.ForAllType:
		// Rule <:∀L
		fmt.Println("\tApplying rule <:∀L", c.String())
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
			fmt.Println("\tApplying rule <:InstantiateR", c.String())
			res := c.InstantiateR(a, vb.Identifier)
			return res, nil
		}

	}

	if vb, ok := b.(*ast.ForAllType); ok {
		// Rule <:∀R
		fmt.Println("\tApplying rule <:∀R", c.String())
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
