// Contains the Parser class, which builds an AST out of
// a program. TODO: specialize errors with line and column number
package parser

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/lexer"
	"github.com/0x0f0f0f/gobba-golang/token"
)

// Precedence levels
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or < or <= or >=
	CALL        // function call
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X
)

// A Pratt parser consists in semantic code, or the association
// of parsing functions with token types. There can be two types
// of parsing functions, either for infix or prefix operators
type prefixParseFn func() ast.Expression

// The argument for the infix type of functions is the
// left side of the infix operator
type infixParseFn func(ast.Expression) ast.Expression

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
	// These two maps are needed as lookup table for
	// operators either found in prefix or infix position
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// Create a new parser from a given Lexer
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)

	// Read two tokens so that curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// Simply get the list of parsing errors
func (p *Parser) Errors() []string {
	return p.errors
}

// Register a prefix parse function
func (p *Parser) registerPrefix(tt token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tt] = fn
}

// Register a infix parse function
func (p *Parser) registerInfix(tt token.TokenType, fn infixParseFn) {
	p.infixParseFns[tt] = fn
}

// Add an error when a peekToken is not the expected one
func (p *Parser) peekError(tt token.TokenType, t token.Token) {
	msg := fmt.Sprintf("line %d column %d: expected '%s'. got '%s' instead",
		t.Line, t.Column,
		tt, p.peekToken.Literal)
	p.errors = append(p.errors, msg)
}

// Advance parsing by a token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// If the peek token matches with expectation, advance
// and return true, false otherwise. Enforce correctness
// of the order of tokens by checking the type of the next token.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t, p.peekToken)
		return false
	}
}

// ==================================================================
// Actual parsing of nodes
// ==================================================================

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

// Parse an assignment
func (p *Parser) parseAssignment() *ast.Assignment {
	// TODO nil checks
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	ass := &ast.Assignment{Token: p.curToken}
	ass.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// TODO nil checks
	if !p.expectPeek(token.EQUALS) {
		return nil
	}

	// TODO: skipping value expression until end of the statement
	for !p.curTokenIs(token.SEMI) && !p.curTokenIs(token.AND) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}
	return ass
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

// Parse a single expression statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMI) {
		p.nextToken()
	}
	return stmt
}

// TODO fix!!!
func (p *Parser) parseExpression(prec int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}
