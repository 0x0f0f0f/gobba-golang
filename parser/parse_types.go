package parser

import (
	// "fmt"
	"github.com/0x0f0f0f/gobba-golang/ast"
	"github.com/0x0f0f0f/gobba-golang/token"
	// "strconv"
)

// Parse a type value
func (p *Parser) parseTypeValue() ast.TypeValue {
	if p.curTokenIs(token.IDENT) {
		return &ast.TyUnVar{Identifier: ast.UniqueIdentifier{Value: p.curToken.Literal}}
	} else if p.curTokenIs(token.UNIT) {
		return &ast.TyUnit{}
	}
	p.expectedType(p.curToken)
	return nil
}
