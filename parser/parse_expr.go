package parser

import (
	"fmt"
	"strconv"

	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/token"
)

func (p *Parser) ParseExpression(prec int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken)
		return nil
	}
	leftExp := prefix()

	// for !p.peekTokenIs(token.SEMI) && prec < p.peekPrecedence() {
	for !p.peekTokenIs(token.EOF) && prec < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// ======================================================================
// Literals
// ======================================================================

// Parse a simple terminal symbol
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.ExprIdentifier{
		Token:      p.curToken,
		Identifier: ast.UniqueIdentifier{Value: p.curToken.Literal},
	}
}

// Parse an integer literal
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.customError(nil, p.curToken, "could not parse as integer")
		return nil
	}

	lit.Value = value
	return lit
}

// Parse a floating point literal
func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.customError(nil, p.curToken, "could not parse as float")
		return nil
	}

	lit.Value = value
	return lit
}

// Parse a floating point literal
func (p *Parser) parseComplexLiteral() ast.Expression {
	lit := &ast.ComplexLiteral{Token: p.curToken}

	value := 0 + 0i
	_, err := fmt.Sscanf(p.curToken.Literal, "%f", &value)
	if err != nil {
		p.customError(nil, p.curToken, "could not parse as complex")
		return nil
	}

	lit.Value = value
	return lit
}

// Parse a string literal
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// Parse a boolean literal
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BoolLiteral{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}

// Parse a subexpression grouped by ()
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	if p.curTokenIs(token.RPAREN) {
		return &ast.UnitLiteral{Token: p.curToken}
	}

	exp := p.ParseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// Parse an expression grouped by {}
func (p *Parser) parseBraceGroupedExpression() ast.Expression {
	p.nextToken()

	if p.curTokenIs(token.RBRACKET) {
		return &ast.UnitLiteral{Token: p.curToken}
	}

	exp := p.ParseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseDollarExpression() ast.Expression {
	p.nextToken()
	exp := p.ParseExpression(LOWEST)
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.ExprIf{Token: p.curToken}

	p.nextToken()

	exp.Condition = p.ParseExpression(LOWEST)

	if !p.expectPeek(token.THEN) {
		return nil
	}

	p.nextToken()

	exp.Consequence = p.ParseExpression(LOWEST)

	if !p.expectPeek(token.ELSE) {
		return nil
	}

	p.nextToken()

	exp.Alternative = p.ParseExpression(LOWEST)
	return exp
}
