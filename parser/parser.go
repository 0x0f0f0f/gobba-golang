// Contains the Parser class, which builds an AST out of
// a lexer.
package parser

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/lexer"
	"github.com/0x0f0f0f/gobba-golang/token"
)

// Precedence levels for operators
const (
	_           int = iota
	LOWEST          // Terminal expression
	COMPOSITION     // >=> or <=<
	SEQUENCING      // ;
	OR              // ||
	AND             // &&
	EQUALS          // = or !=
	COMPARISON      // > or < or <= or >=
	CONS            // ::
	SUM             // + or -
	PRODUCT         // * and /
	MODULO          // %
	POWER           // ^
	CALL            // function application f x y
	ACCESS          // @ and .
	PREFIX          // -X or !X (and function call?)
)

var precedences = map[token.TokenType]int{
	token.COMPOSE:   COMPOSITION,
	token.PIPE:      COMPOSITION,
	token.SEMI:      SEQUENCING,
	token.OR:        OR,
	token.LAND:      AND,
	token.EQUALS:    EQUALS,
	token.DIFFERS:   EQUALS,
	token.LESS:      COMPARISON,
	token.LESSEQ:    COMPARISON,
	token.GREATER:   COMPARISON,
	token.GREATEREQ: COMPARISON,
	token.CONS:      CONS,
	token.CONCAT:    CONS,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.FPLUS:     SUM,
	token.FMINUS:    SUM,
	token.CPLUS:     SUM,
	token.CMINUS:    SUM,
	token.TIMES:     PRODUCT,
	token.DIVIDE:    PRODUCT,
	token.FTIMES:    PRODUCT,
	token.FDIVIDE:   PRODUCT,
	token.CTIMES:    PRODUCT,
	token.CDIVIDE:   PRODUCT,
	token.MODULO:    MODULO,
	token.TOPOW:     POWER,
	token.FTOPOW:    POWER,
	token.CTOPOW:    POWER,
	token.ACCESS:    ACCESS,
	token.AT:        ACCESS,
	// function application
	token.LPAREN: CALL,
}

var rightAssociative = map[token.TokenType]bool{
	token.CONS: true,
}

// ======================================================================

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
	errors    []ParserError
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
		errors: []ParserError{},
	}

	// Registration of prefix operators
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.COMPLEX, p.parseComplexLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	// TODO rune
	// TODO vectors ???
	// TODO lists
	// TODO records

	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.DOLLAR, p.parseDollarExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.LAMBDA, p.parseExprLambda)
	p.registerPrefix(token.LET, p.parseLetExpression)

	// Registration of infix operators
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.COMPOSE, p.parseInfixExpression)
	p.registerInfix(token.PIPE, p.parseInfixExpression)
	p.registerInfix(token.LAND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.EQUALS, p.parseInfixExpression)
	p.registerInfix(token.DIFFERS, p.parseInfixExpression)
	p.registerInfix(token.LESS, p.parseInfixExpression)
	p.registerInfix(token.LESSEQ, p.parseInfixExpression)
	p.registerInfix(token.GREATER, p.parseInfixExpression)
	p.registerInfix(token.GREATEREQ, p.parseInfixExpression)
	p.registerInfix(token.CONCAT, p.parseInfixExpression)
	// Arithmetical Operartors
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.TIMES, p.parseInfixExpression)
	p.registerInfix(token.DIVIDE, p.parseInfixExpression)
	p.registerInfix(token.MODULO, p.parseInfixExpression)
	p.registerInfix(token.TOPOW, p.parseInfixExpression)
	p.registerInfix(token.FPLUS, p.parseInfixExpression)
	p.registerInfix(token.FMINUS, p.parseInfixExpression)
	p.registerInfix(token.FTIMES, p.parseInfixExpression)
	p.registerInfix(token.FDIVIDE, p.parseInfixExpression)
	p.registerInfix(token.FTOPOW, p.parseInfixExpression)
	p.registerInfix(token.CPLUS, p.parseInfixExpression)
	p.registerInfix(token.CMINUS, p.parseInfixExpression)
	p.registerInfix(token.CTIMES, p.parseInfixExpression)
	p.registerInfix(token.CDIVIDE, p.parseInfixExpression)
	p.registerInfix(token.CTOPOW, p.parseInfixExpression)

	p.registerInfix(token.ACCESS, p.parseInfixExpression)
	p.registerInfix(token.AT, p.parseInfixExpression)

	p.registerInfix(token.CONS, p.parseInfixRightAssocExpression)

	p.registerInfix(token.SEMI, p.parseInfixSequence)

	// function application
	p.registerInfix(token.LPAREN, p.parseExprApply)

	// Read two tokens so that curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// Register a prefix parse function
func (p *Parser) registerPrefix(tt token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tt] = fn
}

// Register a infix parse function
func (p *Parser) registerInfix(tt token.TokenType, fn infixParseFn) {
	p.infixParseFns[tt] = fn
}

// Advance parsing by a token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) resetToken(t token.Token) {
	p.l.ResetPosition(t.Position)
	p.curToken = t
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// Returns true if a token is right associative
func (p *Parser) isRightAssociative(t token.Token) bool {
	if v, ok := rightAssociative[t.Type]; ok {
		return v
	}
	return false
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

// ======================================================================
// Actual parsing functions
// ======================================================================

// Parse an assignment
func (p *Parser) parseAssignment() *ast.Assignment {
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	ass := &ast.Assignment{Token: p.curToken}
	ass.Name = &ast.ExprIdentifier{
		Token: p.curToken,
		Identifier: ast.UniqueIdentifier{
			Value: p.curToken.Literal,
		},
	}

	if !p.expectPeek(token.EQUALS) {
		return nil
	}

	p.nextToken()
	ass.Value = p.ParseExpression(SEQUENCING)

	if f, ok := ass.Value.(*ast.ExprLambda); ok {
		// Combinator for recursion
		fix := &ast.ExprFix{
			Token: f.Token,
			Param: *ass.Name,
			Body:  f,
		}

		ass.Value = fix

		return ass

	}
	return ass
}

func (p *Parser) ParseProgram() ast.Expression {
	ast.ResetUIDCounter()
	expr := p.ParseExpression(LOWEST)
	if !p.expectPeek(token.EOF) {
		return nil
	}
	return expr
}
