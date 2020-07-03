package ast

import (
	"fmt"
)

// This file contains definitions of types
// See https://github.com/chrisnevers/bidirectional-typechecking/blob/master/lib/ast/type.ml

type TypeValue interface {
	typeValue()
	String() string
}

// ======================================================================
// Definitions of types of values that compose an algorithmic type context
// ======================================================================

type UnitType struct{}

func (u *UnitType) typeValue() {}
func (u *UnitType) String() string {
	return "unit"
}

type IntegerType struct{}

func (u *IntegerType) typeValue() {}
func (u *IntegerType) String() string {
	return "int"
}

type FloatType struct{}

func (u *FloatType) typeValue() {}
func (u *FloatType) String() string {
	return "float"
}

type BoolType struct{}

func (u *BoolType) typeValue() {}
func (u *BoolType) String() string {
	return "bool"
}

type StringType struct{}

func (u *StringType) typeValue() {}
func (u *StringType) String() string {
	return "string"
}

type RuneType struct{}

func (u *RuneType) typeValue() {}
func (u *RuneType) String() string {
	return "rune"
}

// Denoted with α in the paper
type VariableType struct {
	Identifier UniqueIdentifier
}

func (u *VariableType) typeValue() {}
func (u *VariableType) String() string {
	return "'" + u.Identifier.Value
}

// Denoted with ∀α. A in the paper
type ForAllType struct {
	Identifier UniqueIdentifier
	Type       TypeValue
}

func (u *ForAllType) typeValue() {}
func (u *ForAllType) String() string {
	return fmt.Sprintf("∀%s.%s", u.Identifier.Value, u.Type.String())
}

// Denoted with A → B in the paper
type LambdaType struct {
	Domain   TypeValue
	Codomain TypeValue
}

func (u *LambdaType) typeValue() {}
func (u *LambdaType) String() string {
	return fmt.Sprintf("%s -> %s", u.Domain.String(), u.Codomain.String())
}

// Denoted with α^ in the paper
type ExistsType struct {
	Identifier UniqueIdentifier
}

func (u *ExistsType) typeValue() {}
func (u *ExistsType) String() string {
	return "∃'" + u.Identifier.Value
}

// TODO record types
