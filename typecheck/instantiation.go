package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// Defined in the Instantiation paragraph:
// α^ :=< A, instantiate α^ such that α^ <: A
func (c Context) InstantiateL(alpha ast.UniqueIdentifier, ty ast.TypeValue) Context {
	c.debugSection("InstantiateL", alpha.FullString(), ":=<", ty.FullString())
	exv := &ExistentialVariable{alpha, nil}
	leftc, rightc := c.SplitAt(exv)

	if ty.IsMonotype() && leftc.IsWellFormed(ty) {
		// Rule InstLSolve
		c.debugRule("InstLSolve")

		solved := &ExistentialVariable{alpha, &ty}
		c = c.Insert(exv, []ContextValue{solved})
		c.debugRuleOut("InstLSolve")
		// return c
	}

	switch vty := ty.(type) {
	case *ast.TyLambda:
		// Rule InstLArr
		c.debugRule("InstLArr")

		alpha1 := ast.GenUID("α")
		alpha2 := ast.GenUID("α")

		var arrow ast.TypeValue = &ast.TyLambda{
			Domain:   &ast.TyExVar{Identifier: alpha1},
			Codomain: &ast.TyExVar{Identifier: alpha2},
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
		delta.debugRuleOut("InstLArr")
		return delta

	case *ast.TyForAll:
		// Rule InstLAllR
		c.debugRule("InstLAllR")

		unv := &UniversalVariable{vty.Identifier}
		delta := c.InsertHead(unv).InstantiateL(alpha, vty.Type)

		delta.debugRuleOut("InstLAllR")
		return delta.Drop(unv)

	case *ast.TyExVar:
		// Rule InstLReach
		c.debugRule("InstLReach")

		beta := vty.Identifier

		if rightc.IsWellFormed(ty) {
			exv := &ExistentialVariable{beta, nil}
			var vt ast.TypeValue = &ast.TyExVar{Identifier: alpha}

			outc := c.Insert(exv, []ContextValue{
				&ExistentialVariable{
					Identifier: vty.Identifier,
					Value:      &vt,
				},
			})
			outc.debugRuleOut("InstLReach")
			return outc
		}
	}

	return c
}

// A =<: α^, instantiate α^ such that A <: α^
func (c Context) InstantiateR(ty ast.TypeValue, alpha ast.UniqueIdentifier) Context {

	exv := &ExistentialVariable{alpha, nil}
	leftc, rightc := c.SplitAt(exv)
	if ty.IsMonotype() && leftc.IsWellFormed(ty) {
		// Rule InstRSolve
		c.debugRule("InstRSolve")

		solved := &ExistentialVariable{alpha, &ty}
		c = c.Insert(exv, []ContextValue{solved})
		c.debugRuleOut("InstRSolve")
		// return c
	}

	switch va := ty.(type) {
	case *ast.TyLambda:
		// Rule InstRArr
		c.debugRule("InstRArr")

		alpha1 := ast.GenUID("α")
		alpha2 := ast.GenUID("α")

		var arrow ast.TypeValue = &ast.TyLambda{
			Domain:   &ast.TyExVar{Identifier: alpha1},
			Codomain: &ast.TyExVar{Identifier: alpha2},
		}

		gamma := c.InsertHead(&ExistentialVariable{
			Identifier: alpha,
			Value:      &arrow,
		}).InsertHead(&ExistentialVariable{
			Identifier: alpha1,
		}).InsertHead(&ExistentialVariable{
			Identifier: alpha2,
		})

		theta := gamma.InstantiateL(alpha1, va.Domain)
		delta := theta.InstantiateR(theta.Apply(va.Codomain), alpha2)
		delta.debugRuleOut("InstRArr")
		return delta
	case *ast.TyForAll:
		// Rule InstRAllL
		c.debugRule("InstRAllL")

		beta1 := ast.GenUID("β")
		marker := &Marker{Identifier: beta1}
		beta1exv := &ExistentialVariable{
			Identifier: beta1,
		}
		ext := &ast.TyExVar{
			Identifier: beta1,
		}

		gamma := c.InsertHead(beta1exv).InsertHead(marker)
		delta := gamma.InstantiateR(Substitution(va.Type, ext, va.Identifier), alpha)

		delta.debugRuleOut("InstRAllL")
		return delta.Drop(marker)

	case *ast.TyExVar:
		// Rule InstRReach
		c.debugRule("InstRReach")

		var exv ast.TypeValue = &ast.TyExVar{Identifier: alpha}
		if rightc.IsWellFormed(ty) {
			outc := c.InsertHead(&ExistentialVariable{
				Identifier: va.Identifier,
				Value:      &exv,
			})
			outc.debugRuleOut("InstRReach")
			return outc
		}

	}
	return c
}
