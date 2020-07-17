package parser

import (
	"testing"

	"github.com/0x0f0f0f/gobba-golang/lexer"
	"github.com/stretchr/testify/assert"
)

func CheckParserErrors(t *testing.T, p *Parser) {
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

func TestFailures(t *testing.T) {
	tests := []string{
		"4 32",
		"9283c 9n8f29n3f890jn29083fn=-=-dpvp3;r=;2c./23.c",
	}
	for _, tt := range tests {
		l := lexer.New(tt)
		p := New(l)
		_ = p.ParseProgram()
		assert.NotEqual(t, 0, p.errors)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"-a",
			"(-a)",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4 - -5 * 5",
			"((3 + 4) - ((-5) * 5))",
		},
		{
			"5 > 4 = 3 < 4",
			"((5 > 4) = (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 = 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) = ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 = 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) = ((3 * 1) + (4 * 5)))",
		},
		{
			"5 + 3 ; 2 >=> 4 * 2 - 1",
			"(((5 + 3) ; 2) >=> ((4 * 2) - 1))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 = false",
			"((3 > 5) = false)",
		},
		{
			"3 < 5 = true",
			"((3 < 5) = true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"2 / $ 5 + 5",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true = true)",
			"(!(true = true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, (2 * 3), (4 + 5), $ add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add (a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"add()",
			"add(())",
		},
		{
			"h.f :: a :: b ++ c",
			"((h . f) :: (a :: (b ++ c)))",
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
