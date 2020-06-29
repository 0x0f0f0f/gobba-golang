package parser

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/lexer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input               string
		expectedIdentifiers []string
		expectedValues      []interface{}
	}{
		{"let x = 5;", []string{"x"}, []interface{}{5}},
		{"let x = 5 and y = 4;", []string{"x", "y"}, []interface{}{5, 4}},
		{"let y = true;", []string{"y"}, []interface{}{true}},
		{"let foobar = y;", []string{"foobar"}, []interface{}{"y"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		assert.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		testLetStatement(t, stmt, tt.expectedIdentifiers, tt.expectedValues)

	}
}

func testLetStatement(t *testing.T, s ast.Statement, names []string, values []interface{}) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	assert.True(t, ok, "casting to *ast.LetStatement")

	assert.Equal(t, len(letStmt.Assignments), len(names))
	assert.Equal(t, len(letStmt.Assignments), len(values))

	for i, ass := range letStmt.Assignments {
		assert.Equal(t, names[i], ass.Name.Value)
		assert.Equal(t, names[i], ass.Name.TokenLiteral())

		testLiteralExpression(t, ass.Value, values[i])
	}
	return true
}
