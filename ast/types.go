package ast

// This file contains definitions of types
// See https://github.com/chrisnevers/bidirectional-typechecking/blob/master/lib/ast/type.ml

type TypeValue interface {
	typeValue()
	String() string
	FullString() string // Also display UID numbers
	IsMonotype() bool
	// prints like a -> b -> c for polymorphic types
	// instead of the full UID
	FancyString(map[UniqueIdentifier]int) string
}

// ======================================================================
// Definitions of types of values that compose an algorithmic type context
// ======================================================================

type UnitType struct{}
type IntegerType struct{}
type FloatType struct{}
type ComplexType struct{}
type NumberType struct{}
type BoolType struct{}
type StringType struct{}
type RuneType struct{}

// Denoted with α in the paper
type VariableType struct {
	Identifier UniqueIdentifier
}

// Denoted with ∀α. A in the paper
type ForAllType struct {
	Identifier UniqueIdentifier
	Type       TypeValue
}

// Denoted with A → B in the paper
type LambdaType struct {
	Domain   TypeValue
	Codomain TypeValue
}

// Denoted with α^ in the paper
type ExistsType struct {
	Identifier UniqueIdentifier
}

func (u *UnitType) typeValue()     {}
func (u *IntegerType) typeValue()  {}
func (u *FloatType) typeValue()    {}
func (u *ComplexType) typeValue()  {}
func (u *NumberType) typeValue()   {}
func (u *BoolType) typeValue()     {}
func (u *StringType) typeValue()   {}
func (u *RuneType) typeValue()     {}
func (u *VariableType) typeValue() {}
func (u *ForAllType) typeValue()   {}
func (u *LambdaType) typeValue()   {}
func (u *ExistsType) typeValue()   {}

func (u *UnitType) IsMonotype() bool     { return true }
func (u *IntegerType) IsMonotype() bool  { return true }
func (u *FloatType) IsMonotype() bool    { return true }
func (u *ComplexType) IsMonotype() bool  { return true }
func (u *NumberType) IsMonotype() bool   { return true }
func (u *BoolType) IsMonotype() bool     { return true }
func (u *StringType) IsMonotype() bool   { return true }
func (u *RuneType) IsMonotype() bool     { return true }
func (u *VariableType) IsMonotype() bool { return true }
func (u *ForAllType) IsMonotype() bool   { return false }
func (u *ExistsType) IsMonotype() bool   { return true }
func (u *LambdaType) IsMonotype() bool   { return u.Domain.IsMonotype() && u.Codomain.IsMonotype() }

// TODO record types
