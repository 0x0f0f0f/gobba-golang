package parser

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/lexer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseExpression(LOWEST)
	CheckParserErrors(t, p)
	testUniqueIdentifier(t, program, ast.UniqueIdentifier{"foobar", 0})
}

func TestBooleanExpression(t *testing.T) {
	input := "true; false;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseExpression(LOWEST)
	CheckParserErrors(t, p)

	stmt, ok := program.(*ast.InfixExpression)
	assert.True(t, ok, "casting to *ast.InfixExpression")
	testBoolLiteral(t, stmt.Left, true)
	testBoolLiteral(t, stmt.Right, false)
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseExpression(LOWEST)
	CheckParserErrors(t, p)

	literal, ok := program.(*ast.IntegerLiteral)
	assert.True(t, ok, "casting to *ast.IntegerLiteral")

	assert.Equal(t, literal.Value, int64(5))
	assert.Equal(t, literal.TokenLiteral(), "5")
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"-15;", "-", 15},
		{"-15;", "-", 15},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseExpression(LOWEST)
		CheckParserErrors(t, p)

		exp, ok := program.(*ast.PrefixExpression)
		assert.True(t, ok, "casting to *ast.PrefixExpression")

		assert.Equal(t, exp.Operator, tt.operator)
		testLiteralExpression(t, exp.Right, tt.value)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5.243 - 2.23e2;", 5.243, "-", 2.23e2},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 :: 5;", 5, "::", 5},
		{"5 % 5;", 5, "%", 5},
		{"5 = 5;", 5, "=", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true = true;", true, "=", true},
		{"true != false;", true, "!=", false},
		{"false = false;", false, "=", false},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseExpression(LOWEST)
		CheckParserErrors(t, p)

		testInfixExpression(t, program, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestLetExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"let x = 5; x", "(λ x . x)(5)"},
		// {"let x = 5 and y = 4;", []string{"x", "y"}, []interface{}{5, 4}},
		// {"let y = true;", []string{"y"}, []interface{}{true}},
		// {"let foobar = y;", []string{"foobar"}, []interface{}{"y"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseExpression(LOWEST)
		CheckParserErrors(t, p)

		actual := program.String()
		assert.Equal(t, tt.expected, actual)
	}
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case float64:
		return testFloatLiteral(t, exp, v)
	case complex128:
		return testComplexLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	assert.True(t, ok, "casting to *ast.InfixExpression")
	testLiteralExpression(t, opExp.Left, left)
	assert.Equal(t, opExp.Operator, operator)
	testLiteralExpression(t, opExp.Right, right)
	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	assert.True(t, ok, "is *ast.IntegerLiteral")
	assert.Equal(t, integ.Value, value)
	assert.Equal(t, integ.TokenLiteral(), fmt.Sprintf("%d", value))
	return true
}

func testFloatLiteral(t *testing.T, il ast.Expression, value float64) bool {
	integ, ok := il.(*ast.FloatLiteral)
	assert.True(t, ok, "is *ast.IntegerLiteral")
	assert.Equal(t, integ.Value, value)
	return true
}

func testComplexLiteral(t *testing.T, il ast.Expression, value complex128) bool {
	integ, ok := il.(*ast.ComplexLiteral)
	assert.True(t, ok, "is *ast.IntegerLiteral")
	assert.Equal(t, integ.Value, value)
	return true
}

func testBoolLiteral(t *testing.T, il ast.Expression, value bool) bool {
	integ, ok := il.(*ast.BoolLiteral)
	assert.True(t, ok, "is *ast.BoolLiteral")
	assert.Equal(t, integ.Value, value)
	assert.Equal(t, integ.TokenLiteral(), fmt.Sprintf("%t", value))
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.IdentifierExpr)
	assert.True(t, ok, "casting to *ast.Identifier")
	assert.Equal(t, value, ident.Identifier.Value)
	return true
}

func testUniqueIdentifier(t *testing.T, exp ast.Expression, value ast.UniqueIdentifier) bool {
	ident, ok := exp.(*ast.IdentifierExpr)
	assert.True(t, ok, "casting to *ast.Identifier")
	assert.Equal(t, value, ident.Identifier)
	assert.Equal(t, value.Value, ident.TokenLiteral())
	return true
}
