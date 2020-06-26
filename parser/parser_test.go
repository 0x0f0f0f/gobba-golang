package parser

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/lexer"
	"strings"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10 and z = 3; 
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	numStat := len(strings.Split(input, ";")) - 1
	if len(program.Statements) != numStat {
		t.Fatalf("program.Statements does not contain %d statements, got=%d",
			numStat, len(program.Statements))

	}

	tests := []struct {
		expectedIdentifiers []string
	}{
		{[]string{"x"}},
		{[]string{"y", "z"}},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifiers) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, names []string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if len(letStmt.Assignments) < 1 {
		t.Errorf("letStmt contains %d assignments instead of %d",
			len(letStmt.Assignments), len(names))
		return false
	}

	for i, ass := range letStmt.Assignments {
		if ass.Name.Value != names[i] {
			t.Errorf("Expected identifier '%s'. got=%s",
				names[i], ass.Name.Value)
			return false
		}

		if ass.Name.TokenLiteral() != names[i] {
			t.Errorf("Expected identifier '%s'. got=%s",
				names[i], ass.Name.TokenLiteral())
			return false
		}
	}
	return true
}

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

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected a single statement")
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected an ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression is not an Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Fatalf("identifier is not foobar. got=%s", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("identifier is not foobar. got=%s", ident.TokenLiteral())
	}
}
