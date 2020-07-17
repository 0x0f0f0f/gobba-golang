package parser

import (
	// "fmt"

	"fmt"

	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/token"
	// "strconv"
)

// Parse the arguments of a function definition/literal (TODO allow type annotations)
// and return them as a slice
func (p *Parser) parseFunArgs() []ast.ExprAnnot {
	args := []ast.ExprAnnot{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	eannot := p.parseFunArgAnnot()
	args = append(args, *eannot)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		eannot := p.parseFunArgAnnot()
		args = append(args, *eannot)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

// Parse a type annotation
func (p *Parser) parseFunArgAnnot() *ast.ExprAnnot {
	name := p.parseIdentifier()
	iid, ok := name.(*ast.ExprIdentifier)
	if !ok {
		panic("fatal parsing error")
	}

	if !p.expectPeek(token.ANNOT) {
		return nil
	}
	p.nextToken()

	ty := p.parseTypeValue()

	return &ast.ExprAnnot{
		Token: iid.Token,
		Body:  iid,
		Type:  ty,
	}

}

// NOTE: Function literals hold a single parameter. Multi-parameter
// functions are composed of nested single parameter functions in the AST.
// Type annotations are required by now. Return value is after the function arguments
func (p *Parser) parseExprLambda() ast.Expression {
	//start_token := p.curToken

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// Parse arguments with annotations
	annots := p.parseFunArgs()
	if len(annots) == 0 {
		annots = append(annots, ast.ExprAnnot{
			Body: &ast.ExprIdentifier{
				Token:      p.curToken,
				Identifier: ast.UniqueIdentifier{Value: "_", Id: 0},
			},
			Type: ast.TyUnit{},
		})
	}

	// Unzip function arguments and types
	args := []*ast.ExprIdentifier{}
	types := []ast.TypeValue{}
	for _, annot := range annots {
		iid, ok := annot.Body.(*ast.ExprIdentifier)
		if !ok {
			return nil
		}

		args = append(args, iid)
		types = append(types, annot.Type)
	}

	// Parse the return type
	var return_type ast.TypeValue

	if p.peekTokenIs(token.LBRACKET) {
		return_type = ast.TyUnit{}
	} else {
		p.nextToken()
		return_type = p.parseTypeValue()
	}

	if !p.expectPeek(token.LBRACKET) {
		return nil
	}

	body := p.parseBraceGroupedExpression()

	var cur_fun *ast.ExprLambda
	var cur_ty *ast.TyLambda

	// Parameter list unrolling is done with a iterative loop
	for k := range annots {
		fmt.Println(len(annots))
		k = len(annots) - 1 - k
		//el := args[k]

		// Build the function abstraction
		new_fun := &ast.ExprLambda{}
		new_fun.Param = args[k]

		// Build the type annotation
		new_ty := &ast.TyLambda{}
		new_ty.Domain = types[k]

		if k == len(args)-1 {
			new_fun.Body = body
			new_ty.Codomain = return_type
		} else {
			new_fun.Body = cur_fun
			new_ty.Codomain = cur_ty
		}

		cur_fun = new_fun
		cur_ty = new_ty
	}

	return &ast.ExprAnnot{
		Body: cur_fun,
		Type: cur_ty,
	}
}

// If the parsed parameter is an annotation, return a wrapper abstraction
// that applies the original function
// to the type annotation
func replaceTypedFun(old_fun *ast.ExprLambda, param ast.Expression) *ast.ExprLambda {
	annot, ok := param.(*ast.ExprAnnot)
	if !ok {
		id, ok := param.(*ast.ExprIdentifier)
		if !ok || id == nil {
			return nil
			//panic("expected an identifier or an annotation")
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
	appl_annot_expr := &ast.ExprApplySpine{
		Token:    annot.Token,
		Function: old_fun,
		Spine:    []ast.Expression{annot},
	}

	new_fun := &ast.ExprLambda{
		Token: appl_annot_expr.Token,
		Param: iid,
		Body:  appl_annot_expr,
	}

	return new_fun

}

// Parse a function application using spines
func (p *Parser) parseExprApplySpine(f ast.Expression) ast.Expression {
	appl_expr := &ast.ExprApplySpine{Token: p.curToken, Function: f}

	args := p.parseApplyArguments()

	if len(args) == 0 {
		appl_expr.Spine = []ast.Expression{&ast.UnitLiteral{}}
		return appl_expr
	}

	appl_expr.Spine = args

	return appl_expr
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
// let x = 1 and y = 2 in x + y === (lambda y -> (lambda x -> x + y))(1,2)
func (p *Parser) parseLetExpression() ast.Expression {
	curr_fun := &ast.ExprLambda{Token: p.curToken}
	app := &ast.ExprApplySpine{Token: p.curToken, Function: curr_fun}

	// Parse the first assignment
	ass := p.parseAssignment()
	if ass == nil {
		return nil
	}

	curr_fun.Param = ass.Name
	app.Spine = []ast.Expression{ass.Value}

	for !p.peekTokenIs(token.SEMI) && !p.peekTokenIs(token.EOF) {
		p.expectPeek(token.AND)

		ass := p.parseAssignment()
		if ass == nil {
			return nil
		}
		new_fun := &ast.ExprLambda{Token: p.curToken}
		new_fun.Param = ass.Name
		curr_fun.Body = new_fun
		curr_fun = new_fun

		app.Spine = append(app.Spine, ass.Value)

	}

	// Parse the let expression body
	if p.peekTokenIs(token.EOF) {
		curr_fun.Body = &ast.UnitLiteral{}
		return app
	}

	if !p.expectPeek(token.SEMI) {
		return nil
	}

	if p.peekTokenIs(token.EOF) {
		curr_fun.Body = &ast.UnitLiteral{}
		return app
	}

	p.nextToken()

	curr_fun.Body = p.ParseExpression(LOWEST)

	return app
}
