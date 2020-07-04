package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/lexer"
	"github.com/0x0f0f0f/gobba-golang/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSynthExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected ast.TypeValue
	}{
		{"();", &ast.UnitType{}},
		{"true;", &ast.BoolType{}},
		{"false;", &ast.BoolType{}},
		{"4;", &ast.IntegerType{}},
		{"4.5;", &ast.FloatType{}},
		{"4.5+3.2e-2i;", &ast.ComplexType{}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		// parser.CheckParseErrors(t, p)
		assert.Len(t, p.Errors(), 0)
		assert.Len(t, program.Statements, 1)

		exprst, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			assert.Fail(t, "not an expression")
		}

		ctx := NewContext()
		ast.ResetUIDCounter()
		ty, err := ctx.SynthExpr(exprst.Expression)
		if assert.Nil(t, err) {
			assert.Equal(t, tt.expected, ty)
		}
	}
}
