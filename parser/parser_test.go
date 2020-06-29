package parser

import (
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
