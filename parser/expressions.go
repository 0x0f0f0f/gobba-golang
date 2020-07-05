package parser

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/token"
	"strconv"
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
	exp.Right = p.ParseExpression(PREFIX)

	return exp
}

// Parse an infix expression given the left branch
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{}
	exp.Token = p.curToken
	exp.Operator = p.curToken.Literal
	exp.Left = left

	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.ParseExpression(precedence)

	return exp
}

// Parse a sequencing expression
func (p *Parser) parseInfixSequence(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{}
	exp.Token = p.curToken
	exp.Operator = p.curToken.Literal
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
	exp := &ast.InfixExpression{}
	exp.Token = p.curToken
	exp.Operator = p.curToken.Literal
	exp.Left = left

	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.ParseExpression(precedence - 1)

	return exp
}

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
	exp := &ast.IfExpression{Token: p.curToken}

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

// NOTE: Function literals hold a single parameter. Multi-parameter
// functions are composed of nested single parameter functions in the AST
// because this eases currying during evaluation
func (p *Parser) parseFunctionLiteral() ast.Expression {
	parent_fun := &ast.FunctionLiteral{Token: p.curToken}
	cur_fun := parent_fun

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	args := p.parseDefinitionArguments()

	if len(args) == 0 {
		parent_fun.Param = &ast.IdentifierExpr{
			Token:      p.curToken,
			Identifier: ast.UniqueIdentifier{"_", 0},
		}

	} else {
		parent_fun.Param = args[0]

		// Parameter list unrolling is done with a iterative loop
		for _, arg := range args[1:] {
			child_fun := &ast.FunctionLiteral{
				Token: parent_fun.Token,
				Param: arg,
			}

			cur_fun.Body = child_fun
			cur_fun = child_fun
		}
	}

	if !p.expectPeek(token.LBRACKET) {
		return nil
	}
	cur_fun.Body = p.parseBraceGroupedExpression()

	return parent_fun
}

// Parse the arguments of a function definition/literal (TODO allow type annotations)
// and return them as a slice
func (p *Parser) parseDefinitionArguments() []*ast.IdentifierExpr {
	args := []*ast.IdentifierExpr{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	id := p.parseIdentifier()
	iid, ok := id.(*ast.IdentifierExpr)
	if !ok {
		panic("fatal parsing error")
	}
	args = append(args, iid)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		id := p.parseIdentifier()
		iid, ok := id.(*ast.IdentifierExpr)
		if !ok {
			panic("fatal parsing error")
		}
		args = append(args, iid)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseApplyExpression(f ast.Expression) ast.Expression {
	inner_expr := &ast.ApplyExpr{Token: p.curToken, Function: f}

	args := p.parseApplyArguments()

	// TODO unroll
	if len(args) == 0 {
		inner_expr.Arg = &ast.UnitLiteral{}
		return inner_expr
	}

	curr_expr := inner_expr
	curr_expr.Arg = args[0]

	for _, arg := range args[1:] {
		outer_app := &ast.ApplyExpr{Token: inner_expr.Token, Function: curr_expr}
		outer_app.Arg = arg
		curr_expr = outer_app

	}

	return curr_expr
}

// Parse the arguments of a function call and return them as a slice
func (p *Parser) parseApplyArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.ParseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.ParseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
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
	for !p.peekTokenIs(token.SEMI) && !p.peekTokenIs(token.EOF) {
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

	if !p.expectPeek(token.SEMI) {
		return nil
	}

	if p.peekTokenIs(token.EOF) {
		fmt.Println("EOF")
		inner_fun.Body = &ast.UnitLiteral{}
		return curr_app
	}

	p.nextToken()

	inner_fun.Body = p.ParseExpression(LOWEST)

	return curr_app
}
