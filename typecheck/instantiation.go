package typecheck

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// Defined in the Instantiation paragraph:
// α^ :=< B, instantiate α^ such that α^ <: B
func (c Context) InstantiateL(alpha ast.UniqueIdentifier, ty ast.TypeValue) Context {
	fmt.Println("InstantiateL", alpha.FullString(), ":=<", ty.FullString())
	exv := &ExistentialVariable{alpha, nil}
	leftc, rightc := c.SplitAt(exv)
	fmt.Println("Split in", leftc, "and", rightc)
	if ty.IsMonotype() && leftc.IsWellFormed(ty) {
		// Rule InstLSolve
		fmt.Println("\tApplying rule InstLSolve", c)
		solved := &ExistentialVariable{alpha, &ty}
		c = c.Insert(exv, []ContextValue{solved})
	}

	switch vty := ty.(type) {
	case *ast.LambdaType:
		// Rule InstLArr
		fmt.Println("\tApplying rule InstLArr", c)

		alpha1 := ast.GenUID("α")
		alpha2 := ast.GenUID("α")

		var arrow ast.TypeValue = &ast.LambdaType{
			Domain:   &ast.ExistsType{Identifier: alpha1},
			Codomain: &ast.ExistsType{Identifier: alpha2},
		}

		// First premise
		exv := &ExistentialVariable{alpha, nil}
		gamma := c.Insert(exv, []ContextValue{
			&ExistentialVariable{alpha2, nil},
			&ExistentialVariable{alpha1, nil},
			&ExistentialVariable{
				Identifier: alpha,
				Value:      &arrow,
			},
		})
		theta := gamma.InstantiateR(vty.Domain, alpha1)

		// Second premise, output context
		delta := theta.InstantiateL(alpha2, theta.Apply(vty.Codomain))
		return delta

	case *ast.ForAllType:
		// Rule InstLAllR
		fmt.Println("\tApplying rule InstLAllR", c)

		unv := &UniversalVariable{vty.Identifier}
		delta := c.InsertHead(unv).InstantiateL(alpha, vty.Type)
		return delta.Drop(unv)

	case *ast.ExistsType:
		// Rule InstLReach
		fmt.Println("\tApplying rule InstLReach", c)

		beta := vty.Identifier

		if rightc.IsWellFormed(ty) {
			exv := &ExistentialVariable{beta, nil}
			var vt ast.TypeValue = &ast.ExistsType{Identifier: alpha}
			return c.Insert(exv, []ContextValue{
				&ExistentialVariable{
					Identifier: vty.Identifier,
					Value:      &vt,
				},
			})
		}
	}

	return c
}

// A =<: α^, instantiate α^ such that α^ <: B
func (c Context) InstantiateR(a ast.TypeValue, alpha ast.UniqueIdentifier) Context {
	exv := &ExistentialVariable{alpha, nil}
	leftc, rightc := c.SplitAt(exv)
	nc := c
	if a.IsMonotype() && leftc.IsWellFormed(a) {
		// Rule InstRSolve
		solved := &ExistentialVariable{alpha, &a}
		nc = nc.Insert(exv, []ContextValue{solved})
	}

	switch va := a.(type) {
	case *ast.LambdaType:
		// Rule InstRArr
		alpha1 := ast.GenUID("α")
		alpha2 := ast.GenUID("α")

		var arrow ast.TypeValue = &ast.LambdaType{
			Domain:   &ast.ExistsType{Identifier: alpha1},
			Codomain: &ast.ExistsType{Identifier: alpha2},
		}

		gamma := nc.InsertHead(&ExistentialVariable{
			Identifier: alpha,
			Value:      &arrow,
		}).InsertHead(&ExistentialVariable{
			Identifier: alpha1,
			Value:      nil,
		}).InsertHead(&ExistentialVariable{
			Identifier: alpha2,
			Value:      nil,
		})

		theta := gamma.InstantiateL(alpha1, va.Domain)
		delta := gamma.InstantiateR(theta.Apply(va.Codomain), alpha2)

		return delta
	case *ast.ForAllType:
		// Rule InstRallL
		beta1 := ast.GenUID("β")
		marker := &Marker{Identifier: beta1}
		exv := &ExistentialVariable{
			Identifier: beta1,
			Value:      nil,
		}
		ext := &ast.ExistsType{
			Identifier: beta1,
		}

		gamma := nc.InsertHead(exv).InsertHead(marker)
		delta := gamma.InstantiateR(Substitution(va.Type, ext, va.Identifier), alpha)
		return delta.Drop(marker)

	case *ast.ExistsType:
		// Rule InstRReach
		var exv ast.TypeValue = &ast.ExistsType{Identifier: alpha}
		if rightc.IsWellFormed(a) {
			return nc.InsertHead(&ExistentialVariable{
				Identifier: va.Identifier,
				Value:      &exv,
			})
		}

	}
	return nc
}
