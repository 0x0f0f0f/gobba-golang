package parser

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/token"
)

func (p *Parser) ParseProgram() *ast.Program {
	// Allocate AST root
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// Parse a toplevel statement
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	// TODO parse directives
	default:
		return p.parseExpressionStatement()
	}
}

// Parse a let statement (not a let expression)
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	stmt.Assignments = make([]*ast.Assignment, 0)

	for !p.curTokenIs(token.SEMI) {
		ass := p.parseAssignment()
		// TODO nil checks
		if ass == nil {
			return nil
		}
		stmt.Assignments = append(stmt.Assignments, ass)
	}

	return stmt
}
