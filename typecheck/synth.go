package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// This file contains definitions for synthesization rules

// TODO SynthesizesTo
func (c Context) SynthesizesTo(exp ast.Expression) (ast.TypeValue, *Context) {
	return nil, nil
}

// TODO accept statements
func (c Context) SynthExpr(exp ast.Expression) ast.TypeValue {
	t, nc := c.SynthesizesTo(exp)
	return nc.Apply(t)
}

func (c Context) SynthStatement(st ast.Statement) ast.TypeValue {
	switch vs := st.(type) {
	case *ast.ExpressionStatement:
		return c.SynthExpr(vs.Expression)
	default:
		panic("not implemented yet!")
	}

}

func (c Context) SynthProgram(p *ast.Program) []ast.TypeValue {
	// TODO errors, everything else
	types := make([]ast.TypeValue, 0)
	for _, st := range p.Statements {
		types = append(types, c.SynthStatement(st))
	}
	return types
}
