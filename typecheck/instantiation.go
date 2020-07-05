package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// Defined in the Instantiation paragraph:
// α^ :=< B, instantiate α^ such that α^ <: B
func (c Context) InstantiateL(alpha ast.UniqueIdentifier, ty ast.TypeValue) Context {
	c.debugSection("InstantiateL", alpha.FullString(), ":=<", ty.FullString())
	exv := &ExistentialVariable{alpha, nil}
	leftc, rightc := c.SplitAt(exv)

	if ty.IsMonotype() && leftc.IsWellFormed(ty) {
		// Rule InstLSolve
		c.debugRule("InstLSolve")

		solved := &ExistentialVariable{alpha, &ty}
		c = c.Insert(exv, []ContextValue{solved})
	}

	switch vty := ty.(type) {
	case *ast.LambdaType:
		// Rule InstLArr
		c.debugRule("InstLArr")

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
		c.debugRule("InstLAllR")

		unv := &UniversalVariable{vty.Identifier}
		delta := c.InsertHead(unv).InstantiateL(alpha, vty.Type)
		return delta.Drop(unv)

	case *ast.ExistsType:
		// Rule InstLReach
		c.debugRule("InstLReach")

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
func (c Context) InstantiateR(ty ast.TypeValue, alpha ast.UniqueIdentifier) Context {
	exv := &ExistentialVariable{alpha, nil}
	leftc, rightc := c.SplitAt(exv)
	nc := c
	if ty.IsMonotype() && leftc.IsWellFormed(ty) {
		// Rule InstRSolve
		c.debugRule("InstRSolve")

		solved := &ExistentialVariable{alpha, &ty}
		nc = nc.Insert(exv, []ContextValue{solved})
	}

	switch va := ty.(type) {
	case *ast.LambdaType:
		// Rule InstRArr
		c.debugRule("InstRArr")

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
		}).InsertHead(&ExistentialVariable{
			Identifier: alpha2,
		})

		theta := gamma.InstantiateL(alpha1, va.Domain)
		delta := theta.InstantiateR(theta.Apply(va.Codomain), alpha2)

		return delta
	case *ast.ForAllType:
		// Rule InstRAllL
		c.debugRule("InstRAllL")

		beta1 := ast.GenUID("β")
		marker := &Marker{Identifier: beta1}
		beta1exv := &ExistentialVariable{
			Identifier: beta1,
		}
		ext := &ast.ExistsType{
			Identifier: beta1,
		}

		gamma := nc.InsertHead(beta1exv).InsertHead(marker)
		delta := gamma.InstantiateR(Substitution(va.Type, ext, va.Identifier), alpha)
		return delta.Drop(marker)

	case *ast.ExistsType:
		// Rule InstRReach
		c.debugRule("InstRReach")

		var exv ast.TypeValue = &ast.ExistsType{Identifier: alpha}
		if rightc.IsWellFormed(ty) {
			return nc.InsertHead(&ExistentialVariable{
				Identifier: va.Identifier,
				Value:      &exv,
			})
		}

	}
	return nc
}
