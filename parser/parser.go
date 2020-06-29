// Contains the Parser class, which builds an AST out of
// a lexer. TODO: specialize errors with line and column number
package parser

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/lexer"
	"github.com/0x0f0f0f/gobba-golang/token"
	"runtime/debug"
)

// Precedence levels for operators
const (
	_           int = iota
	LOWEST          // Terminal expression
	COMPOSITION     // >=> or <=<
	SEQUENCING      // >>
	OR              // ||
	AND             // &&
	EQUALS          // = or !=
	COMPARISON      // > or < or <= or >=
	CONS            // ::
	CALL            // function application TODO sort this out
	SUM             // + or -
	PRODUCT         // * and /
	MODULO          // %
	POWER           // ^
	ACCESS          // @ and :
	PREFIX          // -X or !X (and function call?)
)

var precedences = map[token.TokenType]int{
	token.COMPOSE:   COMPOSITION,
	token.PIPE:      COMPOSITION,
	token.SEQUENCE:  SEQUENCING,
	token.OR:        OR,
	token.LAND:      AND,
	token.EQUALS:    EQUALS,
	token.DIFFERS:   EQUALS,
	token.LESS:      COMPARISON,
	token.LESSEQ:    COMPARISON,
	token.GREATER:   COMPARISON,
	token.GREATEREQ: COMPARISON,
	token.CONS:      CONS,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.TIMES:     PRODUCT,
	token.DIVIDE:    PRODUCT,
	token.MODULO:    MODULO,
	token.TOPOW:     POWER,
	token.ACCESS:    ACCESS,
	token.AT:        ACCESS,
}

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
	TraceOnError   bool // Debug option to print a stack trace on error
}

// Create a new parser from a given Lexer
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Registration of prefix operators
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.IMAG, p.parseImagLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.DOLLAR, p.parseDollarExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.LAMBDA, p.parseFunctionLiteral)
	// TODO p.registerPrefix(token.LET, p.parseLetExpression)

	// Registration of infix operators
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.COMPOSE, p.parseInfixExpression)
	p.registerInfix(token.PIPE, p.parseInfixExpression)
	p.registerInfix(token.SEQUENCE, p.parseInfixExpression)
	p.registerInfix(token.LAND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.EQUALS, p.parseInfixExpression)
	p.registerInfix(token.DIFFERS, p.parseInfixExpression)
	p.registerInfix(token.LESS, p.parseInfixExpression)
	p.registerInfix(token.LESSEQ, p.parseInfixExpression)
	p.registerInfix(token.GREATER, p.parseInfixExpression)
	p.registerInfix(token.GREATEREQ, p.parseInfixExpression)
	p.registerInfix(token.CONS, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.TIMES, p.parseInfixExpression)
	p.registerInfix(token.DIVIDE, p.parseInfixExpression)
	p.registerInfix(token.MODULO, p.parseInfixExpression)
	p.registerInfix(token.TOPOW, p.parseInfixExpression)
	p.registerInfix(token.ACCESS, p.parseInfixExpression)
	p.registerInfix(token.AT, p.parseInfixExpression)

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
	msg := fmt.Sprintf("syntax error at line %d column %d: expected '%s'. got '%s' instead",
		t.Line, t.Column,
		tt, p.peekToken.Literal)
	if p.TraceOnError {
		fmt.Printf("ERRORING HERE: %s\n", msg)
		debug.PrintStack()
	}
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.Token) {
	msg := fmt.Sprintf("syntax error at line %d column %d: unexpected token %s (no prefix parser)",
		t.Line, t.Column, t.Literal)
	if p.TraceOnError {
		fmt.Printf("ERRORING HERE: %s\n", msg)
		debug.PrintStack()
	}
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

// Get the next token precedence level
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

// Get the current token precedence level
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

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

// ======================================================================
// Actual parsing functions
// ======================================================================

// Parse an assignment
func (p *Parser) parseAssignment() *ast.Assignment {
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	ass := &ast.Assignment{Token: p.curToken}
	ass.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.EQUALS) {
		return nil
	}

	p.nextToken()
	ass.Value = p.parseExpression(LOWEST)
	return ass
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
func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token: p.curToken}

	stmt.Assignments = make([]*ast.Assignment, 0)

	// Parse the first assignment
	ass := p.parseAssignment()
	if ass == nil {
		return nil
	}
	stmt.Assignments = append(stmt.Assignments, ass)

	for !p.peekTokenIs(token.SEMI) {
		if p.peekTokenIs(token.IN) {
			// This is not a let statement but a let expression
			p.nextToken()
			p.nextToken()
			exp := &ast.LetExpression{
				Token:       stmt.Token,
				Assignments: stmt.Assignments,
			}
			exp.Body = p.parseExpression(LOWEST)
			if p.peekTokenIs(token.SEMI) {
				p.nextToken()
			}
			return &ast.ExpressionStatement{
				Token:      exp.Token,
				Expression: exp,
			}
		} else if p.peekTokenIs(token.AND) {
			p.nextToken()
		}

		ass := p.parseAssignment()
		stmt.Assignments = append(stmt.Assignments, ass)
	}

	if !p.expectPeek(token.SEMI) {
		return nil
	}

	if len(stmt.Assignments) == 0 {
		return nil
	}

	return stmt
}

// Parse a single expression statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMI) {
		return nil
	}
	return stmt
}
