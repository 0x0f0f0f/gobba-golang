package ast

// This file contains definitions of types
// See https://github.com/chrisnevers/bidirectional-typechecking/blob/master/lib/ast/type.ml

type TypeValue interface {
	typeValue()
}

// ======================================================================
// Definitions of types of values that compose an algorithmic type context
// ======================================================================

type UnitType struct{}

func (u *UnitType) typeValue() {}

type IntegerType struct{}

func (u *IntegerType) typeValue() {}

type FloatType struct{}

func (u *FloatType) typeValue() {}

type BoolType struct{}

func (u *BoolType) typeValue() {}

type StringType struct{}

func (u *StringType) typeValue() {}

type RuneType struct{}

func (u *RuneType) typeValue() {}

// Denoted with α in the paper
type VariableType struct {
	Identifier Identifier
}

func (u *VariableType) typeValue() {}

// Denoted with ∀α. A in the paper
type ForAllType struct {
	Identifier Identifier
	Type       TypeValue
}

func (u *ForAllType) typeValue() {}

// Denoted with A → B in the paper
type LambdaType struct {
	Domain   TypeValue
	Codomain TypeValue
}

func (u *LambdaType) typeValue() {}

// TODO record types
