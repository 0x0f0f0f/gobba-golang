package ast

import (
	"fmt"
)

// This file contains definitions of types
// See https://github.com/chrisnevers/bidirectional-typechecking/blob/master/lib/ast/type.ml

type TypeValue interface {
	typeValue()
	String() string
	FullString() string // Also display UID numbers
	IsMonotype() bool
}

// ======================================================================
// Definitions of types of values that compose an algorithmic type context
// ======================================================================

type UnitType struct{}

func (u *UnitType) typeValue()       {}
func (u *UnitType) IsMonotype() bool { return true }
func (u *UnitType) String() string {
	return "unit"
}
func (u *UnitType) FullString() string {
	return "unit"
}

type IntegerType struct{}

func (u *IntegerType) typeValue()       {}
func (u *IntegerType) IsMonotype() bool { return true }
func (u *IntegerType) String() string {
	return "int"
}
func (u *IntegerType) FullString() string {
	return "int"
}

type FloatType struct{}

func (u *FloatType) typeValue()       {}
func (u *FloatType) IsMonotype() bool { return true }
func (u *FloatType) String() string {
	return "float"
}
func (u *FloatType) FullString() string {
	return "float"
}

type ComplexType struct{}

func (u *ComplexType) typeValue()       {}
func (u *ComplexType) IsMonotype() bool { return true }
func (u *ComplexType) String() string {
	return "complex"
}
func (u *ComplexType) FullString() string {
	return "complex"
}

type BoolType struct{}

func (u *BoolType) typeValue()       {}
func (u *BoolType) IsMonotype() bool { return true }
func (u *BoolType) String() string {
	return "bool"
}
func (u *BoolType) FullString() string {
	return "bool"
}

type StringType struct{}

func (u *StringType) typeValue()       {}
func (u *StringType) IsMonotype() bool { return true }
func (u *StringType) String() string {
	return "string"
}
func (u *StringType) FullString() string {
	return "string"
}

type RuneType struct{}

func (u *RuneType) typeValue()       {}
func (u *RuneType) IsMonotype() bool { return true }
func (u *RuneType) String() string {
	return "rune"
}
func (u *RuneType) FullString() string {
	return "rune"
}

// Denoted with α in the paper
type VariableType struct {
	Identifier UniqueIdentifier
}

func (u *VariableType) typeValue()       {}
func (u *VariableType) IsMonotype() bool { return true }
func (u *VariableType) String() string {
	return "'" + u.Identifier.String()
}
func (u *VariableType) FullString() string {
	return "'" + u.Identifier.FullString()
}

// Denoted with ∀α. A in the paper
type ForAllType struct {
	Identifier UniqueIdentifier
	Type       TypeValue
}

func (u *ForAllType) typeValue()       {}
func (u *ForAllType) IsMonotype() bool { return false }
func (u *ForAllType) String() string {
	return fmt.Sprintf("∀%s.%s", u.Identifier.String(), u.Type.String())
}
func (u *ForAllType) FullString() string {
	return fmt.Sprintf("∀%s.%s", u.Identifier.FullString(), u.Type.String())
}

// Denoted with A → B in the paper
type LambdaType struct {
	Domain   TypeValue
	Codomain TypeValue
}

func (u *LambdaType) typeValue() {}
func (u *LambdaType) IsMonotype() bool {
	return u.Domain.IsMonotype() && u.Codomain.IsMonotype()
}
func (u *LambdaType) String() string {
	return fmt.Sprintf("%s -> %s", u.Domain.String(), u.Codomain.String())
}
func (u *LambdaType) FullString() string {
	return fmt.Sprintf("%s -> %s", u.Domain.FullString(), u.Codomain.FullString())
}

// Denoted with α^ in the paper
type ExistsType struct {
	Identifier UniqueIdentifier
}

func (u *ExistsType) typeValue()       {}
func (u *ExistsType) IsMonotype() bool { return true }
func (u *ExistsType) String() string {
	return "∃'" + u.Identifier.String()
}
func (u *ExistsType) FullString() string {
	return "∃'" + u.Identifier.FullString()
}

// TODO record types
