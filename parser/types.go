package parser

import (
	// "fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/token"
	// "strconv"
)

// Default types map

var defaultTypes = map[token.TokenType]ast.TypeValue{
	token.TBOOL: &ast.BoolType{},
	token.TINT:  &ast.IntegerType{},
}

// Parse a type value
func (p *Parser) parseTypeValue() ast.TypeValue {
	if ty, ok := defaultTypes[p.curToken.Type]; ok {
		return ty
	}
	p.expectedType(p.curToken)
	return nil
}

// Parse a type annotation
func (p *Parser) parseFunArgAnnot() ast.Expression {
	name := p.parseIdentifier()
	iid, ok := name.(*ast.IdentifierExpr)
	if !ok {
		panic("fatal parsing error")
	}

	if !p.peekTokenIs(token.ANNOT) {
		return iid
	}
	p.nextToken()
	p.nextToken()

	ty := p.parseTypeValue()

	return &ast.AnnotExpr{
		Token: iid.Token,
		Body:  iid,
		Type:  ty,
	}

}
