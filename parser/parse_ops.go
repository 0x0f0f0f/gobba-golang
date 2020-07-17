package parser

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/token"
)

var OperatorCanBePattern map[string]bool = map[string]bool{
	token.CONS: true,
}

// ======================================================================
// Infix and prefix expressions
// ======================================================================

// Parse an expression with a prefix operator
func (p *Parser) parsePrefixExpression() ast.Expression {
	var isPattern bool
	isPattern, ok := OperatorCanBePattern[p.curToken.Literal]
	if !ok {
		isPattern = false
	}

	exp := &ast.ExprPrefix{
		Token: p.curToken,
		Operator: ast.Operator{
			IsPattern: isPattern,
			Kind:      p.curToken.Literal,
		},
	}

	p.nextToken()
	exp.Right = p.ParseExpression(PREFIX)

	return exp
}

// Parse an infix expression given the left branch
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	var isPattern bool
	isPattern, ok := OperatorCanBePattern[p.curToken.Literal]
	if !ok {
		isPattern = false
	}

	exp := &ast.ExprInfix{
		Token: p.curToken,
		Operator: ast.Operator{
			IsPattern: isPattern,
			Kind:      p.curToken.Literal,
		},
	}
	exp.Left = left

	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.ParseExpression(precedence)

	return exp
}

// Parse a sequencing expression
func (p *Parser) parseInfixSequence(left ast.Expression) ast.Expression {
	exp := &ast.ExprInfix{}
	exp.Token = p.curToken
	exp.Operator = ast.Operator{IsPattern: false, Kind: p.curToken.Literal}
	exp.Left = left

	precedence := p.curPrecedence()

	if p.peekTokenIs(token.EOF) {
		return left
	}
	p.nextToken()
	// Sequencing is right associative
	exp.Right = p.ParseExpression(precedence - 1)

	return exp
}

// Parse a right-associative operator
func (p *Parser) parseInfixRightAssocExpression(left ast.Expression) ast.Expression {
	var isPattern bool
	isPattern, ok := OperatorCanBePattern[p.curToken.Literal]
	if !ok {
		isPattern = false
	}

	exp := &ast.ExprInfix{
		Token: p.curToken,
		Operator: ast.Operator{
			IsPattern: isPattern,
			Kind:      p.curToken.Literal,
		},
	}
	exp.Left = left

	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.ParseExpression(precedence - 1)

	return exp
}
