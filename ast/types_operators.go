package ast

import (
	"github.com/0x0f0f0f/gobba-golang/token"
)

type InfixOperatorType struct {
	Left   TypeValue
	Right  TypeValue
	Result TypeValue
}

func NewInfixOperatorType(left, right, result TypeValue) *InfixOperatorType {
	return &InfixOperatorType{
		Left:   left,
		Right:  right,
		Result: result,
	}
}

type PrefixOperatorType struct {
	Right  TypeValue
	Result TypeValue
}

func NewPrefixOperatorType(right, result TypeValue) *PrefixOperatorType {
	return &PrefixOperatorType{
		Right:  right,
		Result: result,
	}
}

var TINT = NewVariableType(token.TINT)
var TFLOAT = NewVariableType(token.TFLOAT)
var TCOMPLEX = NewVariableType(token.TCOMPLEX)
var TBOOL = NewVariableType(token.TBOOL)
var TRUNE = NewVariableType(token.TRUNE)
var TSTRING = NewVariableType(token.TSTRING)

var DefaultVariableTypes map[string]*VariableType = map[string]*VariableType{
	token.TINT:     TINT,
	token.TFLOAT:   TFLOAT,
	token.TCOMPLEX: TCOMPLEX,
	token.TBOOL:    TBOOL,
	token.TRUNE:    TRUNE,
	token.TSTRING:  TSTRING,
}

var InfixOperatorTypes map[string]*InfixOperatorType = map[string]*InfixOperatorType{
	token.PLUS:    NewInfixOperatorType(TINT, TINT, TINT),
	token.MINUS:   NewInfixOperatorType(TINT, TINT, TINT),
	token.TIMES:   NewInfixOperatorType(TINT, TINT, TINT),
	token.TOPOW:   NewInfixOperatorType(TINT, TINT, TINT),
	token.DIVIDE:  NewInfixOperatorType(TINT, TINT, TINT),
	token.MODULO:  NewInfixOperatorType(TINT, TINT, TINT),
	token.FPLUS:   NewInfixOperatorType(TFLOAT, TFLOAT, TFLOAT),
	token.FMINUS:  NewInfixOperatorType(TFLOAT, TFLOAT, TFLOAT),
	token.FTIMES:  NewInfixOperatorType(TFLOAT, TFLOAT, TFLOAT),
	token.FTOPOW:  NewInfixOperatorType(TFLOAT, TFLOAT, TFLOAT),
	token.FDIVIDE: NewInfixOperatorType(TFLOAT, TFLOAT, TFLOAT),
	token.CPLUS:   NewInfixOperatorType(TCOMPLEX, TCOMPLEX, TCOMPLEX),
	token.CMINUS:  NewInfixOperatorType(TCOMPLEX, TCOMPLEX, TCOMPLEX),
	token.CTIMES:  NewInfixOperatorType(TCOMPLEX, TCOMPLEX, TCOMPLEX),
	token.CTOPOW:  NewInfixOperatorType(TCOMPLEX, TCOMPLEX, TCOMPLEX),
	token.CDIVIDE: NewInfixOperatorType(TCOMPLEX, TCOMPLEX, TCOMPLEX),
	token.LAND:    NewInfixOperatorType(TBOOL, TBOOL, TBOOL),
	token.OR:      NewInfixOperatorType(TBOOL, TBOOL, TBOOL),
}

var PrefixOperatorTypes map[string]*PrefixOperatorType = map[string]*PrefixOperatorType{
	token.MINUS:  NewPrefixOperatorType(TINT, TINT),
	token.FMINUS: NewPrefixOperatorType(TFLOAT, TFLOAT),
	token.CMINUS: NewPrefixOperatorType(TCOMPLEX, TCOMPLEX),
	token.NOT:    NewPrefixOperatorType(TBOOL, TBOOL),
}
