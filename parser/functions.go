package parser

import (
	// "fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/token"
	// "strconv"
)

// Parse the arguments of a function definition/literal (TODO allow type annotations)
// and return them as a slice
func (p *Parser) parseFunArgs() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	eannot := p.parseFunArgAnnot()
	args = append(args, eannot)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		eannot := p.parseFunArgAnnot()
		args = append(args, eannot)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

// NOTE: Function literals hold a single parameter. Multi-parameter
// functions are composed of nested single parameter functions in the AST
// because this eases evaluation
func (p *Parser) parseExprLambda() ast.Expression {
	start_token := p.curToken

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	args := p.parseFunArgs()

	if !p.expectPeek(token.LBRACKET) {
		return nil
	}
	body := p.parseBraceGroupedExpression()

	if len(args) == 0 {
		return &ast.ExprLambda{
			Token: start_token,
			Param: &ast.ExprIdentifier{
				Token:      p.curToken,
				Identifier: ast.UniqueIdentifier{"_", 0},
			},
			Body: body,
		}

	}

	var cur_fun *ast.ExprLambda

	// Parameter list unrolling is done with a iterative loop
	for k := range args {
		k = len(args) - 1 - k
		el := args[k]

		new_fun := &ast.ExprLambda{}
		if k == len(args)-1 {
			new_fun.Body = body
		} else {
			new_fun.Body = cur_fun
		}

		new_fun = replaceTypedFun(new_fun, el)

		cur_fun = new_fun
	}

	return cur_fun
}

// If the parsed parameter is an annotation, return a wrapper abstraction
// that applies the original function
// to the type annotation
func replaceTypedFun(old_fun *ast.ExprLambda, param ast.Expression) *ast.ExprLambda {
	annot, ok := param.(*ast.ExprAnnot)
	if !ok {
		id, ok := param.(*ast.ExprIdentifier)
		if !ok || id == nil {
			panic("expected an identifier or an annotation")
		}
		// The parameter is an identifier, nothing special
		old_fun.Param = id
		old_fun.Token = id.Token
		return old_fun
	}
	iid, ok := annot.Body.(*ast.ExprIdentifier)
	if !ok {
		panic("fatal error. annotated parameter is not an identifier")
	}

	old_fun.Param = iid

	// If the parameter is an annotation
	appl_annot_expr := &ast.ExprApply{
		Token:    annot.Token,
		Function: old_fun,
		Arg:      annot,
	}

	new_fun := &ast.ExprLambda{
		Token: appl_annot_expr.Token,
		Param: iid,
		Body:  appl_annot_expr,
	}

	return new_fun

}
