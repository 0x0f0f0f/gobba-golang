package parser

import (
	"testing"

	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/lexer"
	"github.com/stretchr/testify/assert"
)

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fun (x: int) int {x};", expectedParams: []string{"x"}},
		{input: "fun () {};", expectedParams: []string{"_"}},
		{input: "fun (x: int, y: int, z: int) {x + y + z};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		CheckParserErrors(t, p)

		assert.Len(t, p.Errors(), 0)
		ann, ok := program.(*ast.ExprAnnot)
		assert.True(t, ok, "casting to *ast.ExprAnnot")
		f, ok := ann.Body.(*ast.ExprLambda)
		assert.True(t, ok, "casting to *ast.ExprLambda")

		for i, par := range tt.expectedParams {
			testLiteralExpression(t, f.Param, par)
			if i != len(tt.expectedParams)-1 {
				f, ok = f.Body.(*ast.ExprLambda)
				assert.True(t, ok, "casting to *ast.ExprLambda")
			}
		}
	}
}

func TestFunctionParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"fun () {}",
			"((λ _ . ()): unit -> unit)",
		},
		{
			"fun (x:int ,y : int, z : float) float { (x + y) +. z }",
			"((λ x . (λ y . (λ z . ((x + y) +. z)))): int -> int -> float -> float)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		CheckParserErrors(t, p)
		actual := program.String()
		assert.Equal(t, tt.expected, actual)
	}
}
