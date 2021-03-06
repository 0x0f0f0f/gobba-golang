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

type UnitType struct {
	Identifier UniqueIdentifier
}

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

// ADDITION: encoded with A U B
// type UnionType struct {
// 	Left  TypeValue
// 	Right TypeValue
// }

// Denoted with α^ in the paper
type ExistsType struct {
	Identifier UniqueIdentifier
}

func (u *UnitType) typeValue()     {}
func (u *VariableType) typeValue() {}
func (u *ForAllType) typeValue()   {}
func (u *LambdaType) typeValue()   {}

// func (u *UnionType) typeValue()    {}
func (u *ExistsType) typeValue() {}

func (u *UnitType) IsMonotype() bool     { return true }
func (u *VariableType) IsMonotype() bool { return true }
func (u *ForAllType) IsMonotype() bool   { return false }
func (u *ExistsType) IsMonotype() bool   { return true }
func (u *LambdaType) IsMonotype() bool   { return u.Domain.IsMonotype() && u.Codomain.IsMonotype() }

// func (u *UnionType) IsMonotype() bool    { return u.Left.IsMonotype() && u.Right.IsMonotype() }

// TODO record types

// Default variable types

func NewVariableType(name string) *VariableType {
	return &VariableType{Identifier: UniqueIdentifier{Value: name}}
}

// func NewUnionType(left, right TypeValue) TypeValue {
// 	if CompareTypeValues(left, right) {
// 		return left
// 	}
// 	return &UnionType{Left: left, Right: right}
// }

func CompareTypeValues(a, b TypeValue) bool {
	switch va := a.(type) {
	case *UnitType:
		_, ok := b.(*UnitType)
		return ok
	case *VariableType:
		vb, ok := b.(*VariableType)
		return ok && va.Identifier == vb.Identifier
	case *ExistsType:
		vb, ok := b.(*ExistsType)
		return ok && va.Identifier == vb.Identifier
	case *ForAllType:
		vb, ok := b.(*ForAllType)
		return ok && va.Identifier == vb.Identifier && CompareTypeValues(va.Type, vb.Type)
	// case *UnionType:
	// 	vb, ok := b.(*UnionType)
	// 	return ok && CompareTypeValues(va.Left, vb.Left) && CompareTypeValues(va.Right, vb.Right)
	case *LambdaType:
		vb, ok := b.(*LambdaType)
		return ok && CompareTypeValues(va.Domain, vb.Domain) && CompareTypeValues(va.Codomain, vb.Codomain)
	}
	return false
}
