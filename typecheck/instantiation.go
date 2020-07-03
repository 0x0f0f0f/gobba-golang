package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// Defined in the Instantiation paragraph:
// α^ :=< B, instantiate α^ such that α^ <: B
func (c Context) InstantiateL(alpha ast.UniqueIdentifier, b ast.TypeValue) Context {
	exv := &ExistentialVariable{alpha, nil}
	leftc, rightc := c.SplitAt(exv)
	nc := c
	if b.IsMonotype() && leftc.IsWellFormed(b) {
		// Rule InstLSolve
		solved := &ExistentialVariable{alpha, &b}
		nc = nc.Insert(exv, []ContextValue{solved})
	}

	switch vb := b.(type) {
	case *ast.LambdaType:
		// Rule InstLArr
		alpha1 := ast.GenUID("α")
		alpha2 := ast.GenUID("α")

		var arrow ast.TypeValue = &ast.LambdaType{
			Domain:   &ast.ExistsType{Identifier: alpha1},
			Codomain: &ast.ExistsType{Identifier: alpha2},
		}

		// First premise
		gamma := nc.Insert(exv, []ContextValue{
			&ExistentialVariable{alpha2, nil},
			&ExistentialVariable{alpha1, nil},
			&ExistentialVariable{
				Identifier: alpha,
				Value:      &arrow,
			},
		})
		theta := gamma.InstantiateR(vb.Domain, alpha1)

		// Second premise, output context
		delta := theta.InstantiateL(alpha2, theta.Apply(vb.Codomain))
		return delta

	case *ast.ForAllType:
		// Rule InstLAllR
		unv := &UniversalVariable{vb.Identifier}
		delta := nc.InsertHead(unv).InstantiateL(alpha, vb.Type)
		return delta.Drop(unv)

	case *ast.ExistsType:
		// Rule InstLReach
		unv := &UniversalVariable{vb.Identifier}
		if rightc.IsWellFormed(b) {
			var vt ast.TypeValue = &ast.ExistsType{Identifier: alpha}
			return nc.Insert(unv, []ContextValue{
				&ExistentialVariable{
					Identifier: vb.Identifier,
					Value:      &vt,
				},
			})
		}
	}

	return nc
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
