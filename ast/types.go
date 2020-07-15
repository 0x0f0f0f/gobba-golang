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

type TyUnit struct {
	Identifier UniqueIdentifier
}

// Denoted with α in the paper. Universal variable
type TyUnVar struct {
	Identifier UniqueIdentifier
}

// Denoted with ∀α. A in the paper
type TyForAll struct {
	Identifier UniqueIdentifier
	Type       TypeValue
}

// Denoted with α^ in the paper
type TyExVar struct {
	Identifier UniqueIdentifier
}

// Denoted with A → B in the paper
type TyLambda struct {
	Domain   TypeValue
	Codomain TypeValue
}

// Sum type: encoded with A + B
type TySum struct {
	Left  TypeValue
	Right TypeValue
}

// Product type: encoded with A + B
type TyProduct struct {
	Left  TypeValue
	Right TypeValue
}

func (u *TyUnit) typeValue()    {}
func (u *TyUnVar) typeValue()   {}
func (u *TyForAll) typeValue()  {}
func (u *TyLambda) typeValue()  {}
func (u *TyExVar) typeValue()   {}
func (u *TySum) typeValue()     {}
func (u *TyProduct) typeValue() {}

func (u *TyUnit) IsMonotype() bool    { return true }
func (u *TyUnVar) IsMonotype() bool   { return true }
func (u *TyForAll) IsMonotype() bool  { return false }
func (u *TyExVar) IsMonotype() bool   { return true }
func (u *TyLambda) IsMonotype() bool  { return u.Domain.IsMonotype() && u.Codomain.IsMonotype() }
func (u *TySum) IsMonotype() bool     { return u.Left.IsMonotype() && u.Right.IsMonotype() }
func (u *TyProduct) IsMonotype() bool { return u.Left.IsMonotype() && u.Right.IsMonotype() }

// TODO record types

// Default variable types

func NewTyUnVar(name string) *TyUnVar {
	return &TyUnVar{Identifier: UniqueIdentifier{Value: name}}
}

func CompareTypeValues(a, b TypeValue) bool {
	switch va := a.(type) {
	case *TyUnit:
		_, ok := b.(*TyUnit)
		return ok
	case *TyUnVar:
		vb, ok := b.(*TyUnVar)
		return ok && va.Identifier == vb.Identifier
	case *TyExVar:
		vb, ok := b.(*TyExVar)
		return ok && va.Identifier == vb.Identifier
	case *TyForAll:
		vb, ok := b.(*TyForAll)
		return ok && va.Identifier == vb.Identifier && CompareTypeValues(va.Type, vb.Type)
	case *TySum:
		vb, ok := b.(*TySum)
		return ok && CompareTypeValues(va.Left, vb.Left) && CompareTypeValues(va.Right, vb.Right)
	case *TyLambda:
		vb, ok := b.(*TyLambda)
		return ok && CompareTypeValues(va.Domain, vb.Domain) && CompareTypeValues(va.Codomain, vb.Codomain)
	}
	return false
}
