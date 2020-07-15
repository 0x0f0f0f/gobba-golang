package parser

import (
	// "fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/token"
	// "strconv"
)

// Default types map
var defaultTypes = map[token.TokenType]string{
	token.TBOOL: "bool",
	token.TINT:  "int",
}

// Parse a type value
func (p *Parser) parseTypeValue() ast.TypeValue {
	if p.curTokenIs(token.IDENT) {
		return &ast.VariableType{Identifier: ast.UniqueIdentifier{Value: p.curToken.Literal}}
	} else if p.curTokenIs(token.UNIT) {
		return &ast.UnitType{}
	}
	p.expectedType(p.curToken)
	return nil
}

// Parse a type annotation
func (p *Parser) parseFunArgAnnot() ast.Expression {
	name := p.parseIdentifier()
	iid, ok := name.(*ast.ExprIdentifier)
	if !ok {
		panic("fatal parsing error")
	}

	if !p.peekTokenIs(token.ANNOT) {
		return iid
	}
	p.nextToken()
	p.nextToken()

	ty := p.parseTypeValue()

	return &ast.ExprAnnot{
		Token: iid.Token,
		Body:  iid,
		Type:  ty,
	}

}
