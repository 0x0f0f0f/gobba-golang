package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// Check if a gobba type is
// well formed under an algorithmic context
func (c *Context) IsWellFormed(t ast.TypeValue) bool {
	switch v := t.(type) {
	// Rule UvarWF
	case *ast.VariableType:
		return c.HasTypeVar(v.Identifier)
	// Rule ArrowWF
	case *ast.LambdaType:
		return c.IsWellFormed(v.Domain) && c.IsWellFormed(v.Codomain)
	// Rule ForallWF
	case *ast.ForAllType:
		nc := c.InsertHead(&UniversalVariable{v.Identifier})
		return nc.IsWellFormed(v.Type)
	// Rules EvarWF and SolvedEvarWF
	case *ast.ExistsType:
		return c.HasExistentialVariable(v.Identifier) || nil != c.GetSolvedVariable(v.Identifier)
	default:
		// Primitive types are well formed, rules UnitWF
		return true
	}
}
