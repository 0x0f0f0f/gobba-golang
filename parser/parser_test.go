package parser

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/lexer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b;",
			"((-a) * b);",
		},
		{
			"-a;",
			"(-a);",
		},
		{
			"a + b + c;",
			"((a + b) + c);",
		},
		{
			"a + b - c;",
			"((a + b) - c);",
		},
		{
			"a * b * c;",
			"((a * b) * c);",
		},
		{
			"a * b / c;",
			"((a * b) / c);",
		},
		{
			"a + b / c;",
			"(a + (b / c));",
		},
		{
			"a + b * c + d / e - f;",
			"(((a + (b * c)) + (d / e)) - f);",
		},
		{
			"3 + 4 - -5 * 5;",
			"((3 + 4) - ((-5) * 5));",
		},
		{
			"5 > 4 = 3 < 4;",
			"((5 > 4) = (3 < 4));",
		},
		{
			"5 < 4 != 3 > 4;",
			"((5 < 4) != (3 > 4));",
		},
		{
			"3 + 4 * 5 = 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) = ((3 * 1) + (4 * 5)));",
		},
		{
			"3 + 4 * 5 = 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) = ((3 * 1) + (4 * 5)));",
		},
		{
			"5 + 3 >> 2 >=> 4 * 2 - 1;",
			"(((5 + 3) >> 2) >=> ((4 * 2) - 1));",
		},
		{
			"true;",
			"true;",
		},
		{
			"false;",
			"false;",
		},
		{
			"3 > 5 = false;",
			"((3 > 5) = false);",
		},
		{
			"3 < 5 = true;",
			"((3 < 5) = true);",
		},
		{
			"1 + (2 + 3) + 4;",
			"((1 + (2 + 3)) + 4);",
		},
		{
			"(5 + 5) * 2;",
			"((5 + 5) * 2);",
		},
		{
			"2 / (5 + 5);",
			"(2 / (5 + 5));",
		},
		{
			"2 / $ 5 + 5;",
			"(2 / (5 + 5));",
		},
		{
			"-(5 + 5);",
			"(-(5 + 5));",
		},
		{
			"!(true = true);",
			"(!(true = true));",
		},
		{
			"a + add (b * c) + d;",
			"((a + add((b * c))) + d);",
		},
		{
			"add a b 1 (2 * 3) (4 + 5) $ add 6 7 * 8;",
			"add(a)(b)(1)((2 * 3))((4 + 5))((add(6)(7) * 8));",
		},
		{
			"add (a + b + c * d / f + g);",
			"add((((a + b) + ((c * d) / f)) + g));",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		assert.Equal(t, tt.expected, actual)
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fun x -> x;", expectedParams: []string{"x"}},
		{input: "fun x y z -> x + y + z;", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		assert.Len(t, p.Errors(), 0)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		f, ok := stmt.Expression.(*ast.FunctionLiteral)
		assert.True(t, ok, "casting to *ast.FunctionLiteral")

		for i, par := range tt.expectedParams {
			testLiteralExpression(t, f.Param, par)
			if i != len(tt.expectedParams)-1 {
				f, ok = f.Body.(*ast.FunctionLiteral)
				assert.True(t, ok, "casting to *ast.FunctionLiteral")
			}
		}
	}
}
