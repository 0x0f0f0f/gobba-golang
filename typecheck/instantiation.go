package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// Defined in the Instantiation paragraph:
// α^ :=< B, instantiate α^ such that α^ <: B
func (c Context) InstantiateL(alpha ast.UniqueIdentifier, b ast.TypeValue) *Context {
	exv := &ExistentialVariable{alpha, nil}
	leftc, rightc := c.SplitAt(exv)
	nc := &c
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
			Domain:   &ast.ExistsType{alpha1},
			Codomain: &ast.ExistsType{alpha2},
		}

		// First premise
		gamma := nc.Insert(exv, []ContextValue{
			&ExistentialVariable{alpha1, nil},
			&ExistentialVariable{alpha2, nil},
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
		unv := &UniversalVariable{vb.Identifier}
		if rightc.IsWellFormed(b) {
			var vt ast.TypeValue = &ast.ExistsType{alpha}
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

func (c *Context) InstantiateR(a ast.TypeValue, alpha ast.UniqueIdentifier) *Context {
	return c
}
