package parser

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/token"
	"strconv"
)

func (p *Parser) parseExpression(prec int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMI) && prec < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// Parse a simple terminal symbol
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.IdentifierExpr{
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

// Parse an expression with a prefix operator
func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)

	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{}
	exp.Token = p.curToken
	exp.Operator = p.curToken.Literal
	exp.Left = left

	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)

	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	if p.curTokenIs(token.RPAREN) {
		return &ast.UnitLiteral{Token: p.curToken}
	}

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseDollarExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}

	p.nextToken()

	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.THEN) {
		return nil
	}

	p.nextToken()

	exp.Consequence = p.parseExpression(LOWEST)

	if !p.expectPeek(token.ELSE) {
		return nil
	}

	p.nextToken()

	exp.Alternative = p.parseExpression(LOWEST)
	return exp
}

// NOTE: Function literals hold a single parameter. Multi-parameter
// functions are composed of nested single parameter functions in the AST
// because this eases currying during evaluation
func (p *Parser) parseFunctionLiteral() ast.Expression {
	parent_fun := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	first_param := p.parseIdentifier().(*ast.IdentifierExpr)
	parent_fun.Param = first_param

	// Parameter list unrolling is done with a iterative loop
	cur_fun := parent_fun
	for p.peekTokenIs(token.IDENT) {
		p.nextToken()
		cur_param := p.parseIdentifier().(*ast.IdentifierExpr)

		child_fun := &ast.FunctionLiteral{
			Token: parent_fun.Token,
			Param: cur_param,
		}

		cur_fun.Body = child_fun
		cur_fun = child_fun
	}

	if !p.expectPeek(token.RARROW) {
		return nil
	}
	p.nextToken()
	cur_fun.Body = p.parseExpression(LOWEST)

	return parent_fun
}

func (p *Parser) parseApplyExpression(f ast.Expression) ast.Expression {
	exp := &ast.ApplyExpr{Token: p.curToken, Function: f}

	precedence := p.curPrecedence()

	exp.Arg = p.parseExpression(precedence)

	return exp
}

// Parse a let expression
// Use the let over lambda principle TODO review
// let x = 1 in x + 2 === (lambda x -> x + 2) 1
// let x = 1 and y = 2 in x + y === (lambda y -> (lambda x -> x + y)(1))(2)
func (p *Parser) parseLetExpression() ast.Expression {
	// exp := &ast.LetExpression{Token: p.curToken}
	inner_app := &ast.ApplyExpr{Token: p.curToken}
	inner_fun := &ast.FunctionLiteral{Token: p.curToken}
	inner_app.Function = inner_fun

	// Parse the first assignment
	ass := p.parseAssignment()
	if ass == nil {
		return nil
	}

	inner_fun.Param = ass.Name
	inner_app.Arg = ass.Value

	curr_app := inner_app
	for !p.peekTokenIs(token.IN) {
		p.expectPeek(token.AND)

		ass := p.parseAssignment()
		if ass == nil {
			return nil
		}
		curr_fun := &ast.FunctionLiteral{Token: p.curToken}
		curr_fun.Param = ass.Name
		curr_fun.Body = curr_app

		// Replace
		curr_app = &ast.ApplyExpr{Token: p.curToken}
		curr_app.Function = curr_fun
		curr_app.Arg = ass.Value

	}

	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken()

	inner_fun.Body = p.parseExpression(LOWEST)

	return curr_app
}
