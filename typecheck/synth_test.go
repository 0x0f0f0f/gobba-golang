package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/alpha"
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/lexer"
	"github.com/0x0f0f0f/gobba-golang/parser"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSynthExpr(t *testing.T) {
	alphaid := ast.UniqueIdentifier{
		Value: "α",
		Id:    1,
	}

	alphaext := ast.ExistsType{alphaid}

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
		{"fun (x) {x};", &ast.LambdaType{Domain: &alphaext, Codomain: &alphaext}},
		{"fun (x) {if x then 3 else 4});", &ast.LambdaType{Domain: &ast.BoolType{}, Codomain: &ast.IntegerType{}}},
		{"fun (x) {if x then x else x}(true);", &ast.BoolType{}},
		{"fun (x) {if x then x else false}(true);", &ast.BoolType{}},
		{"fun (x) {x}(2);", &ast.IntegerType{}},
		{"fun (x) {x}(2.2);", &ast.FloatType{}},
		{"let x = 4 and y = 3.2 and f = fun(x,y) {x}; f(y)", &ast.LambdaType{
			Domain:   &ast.ExistsType{ast.UniqueIdentifier{"α", 11}},
			Codomain: &ast.FloatType{},
		}},
	}

	log.SetLevel(log.DebugLevel)

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		// parser.CheckParseErrors(t, p)
		assert.Len(t, p.Errors(), 0)
		alphaconv_program, err := alpha.ProgramAlphaConversion(program)
		if err != nil {
			assert.Fail(t, "could not α-convert expression")
			return
		}

		ctx := NewContext()
		ast.ResetUIDCounter()
		ty, err := ctx.SynthExpr(*alphaconv_program)
		if assert.Nil(t, err) {
			assert.Equal(t, tt.expected, ty)
		}
	}
}
