package ast

import (
	"fmt"
)

// This file contains definitions of types
// See https://github.com/chrisnevers/bidirectional-typechecking/blob/master/lib/ast/type.ml and
// https://github.com/evertedsphere/sound-and-complete/blob/master/src/Types.hs

// ======================================================================
// Polarity and principality
// ======================================================================

type Polarity int

const (
	NONPOLAR = iota
	NEGATIVE
	POSITIVE
)

type Principality int

const (
	PRINCIPAL = iota
	NONPRINCIPAL
)

func (p Principality) String() string {
	if p == PRINCIPAL {
		return "!"
	} else {
		return "/!"
	}
}

// ======================================================================
// Sorts: ★ or ℕ
// ======================================================================
type Sort interface {
	sortType()
	String() string
}

type Star struct{}
type Nat struct{}

func (s Star) sortType() {}
func (s Nat) sortType()  {}

func (s Star) String() string { return "★" }
func (s Nat) String() string  { return "ℕ" }

// ======================================================================
// Propositions and terms
// ======================================================================

// Denoted with zero or succ(t) in the paper
type Term struct {
	Num uint16 // Zero or succ
}

// Equalities
type Prop struct {
	Left  TypeValue
	Right TypeValue
}

func (p Prop) String() string {
	return fmt.Sprintf("%s = %s", p.Left.String(), p.Right.String())
}

// ======================================================================
// Definitions of types and terms
// ======================================================================

type TypeValue interface {
	typeValue()
	IsMonotype() bool
	String() string
	FullString() string // Also display UID numbers
	// prints like a -> b -> c for polymorphic types
	// instead of the full UID
	FancyString(map[UniqueIdentifier]int) string
}

// Denoted with 1 in the paper.
type TyUnit struct {
	Identifier UniqueIdentifier
}

// Denoted with α in the paper. Represents Universal variables.
type TyUnVar struct {
	Identifier UniqueIdentifier
}

// Denoted with α^ in the paper. Exisential variables.
type TyExVar struct {
	Identifier UniqueIdentifier
}

// Denoted with ∀α:κ. A in the paper
type TyForAll struct {
	Identifier UniqueIdentifier
	Sort       Sort
	Type       TypeValue
}

// Denoted with ∃α:κ. A in the paper
type TyExists struct {
	Identifier UniqueIdentifier
	Sort       Sort
	Type       TypeValue
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

// Implication type: encoded with P ⊃ A
// P implies A. Corresponds to type A provided P holds
type TyGuarded struct {
	Type  TypeValue
	Guard Prop
}

// Dual of the implication type: encoded with A ∧ P, read A with P
type TyAsserting struct {
	Type  TypeValue
	Guard Prop
}

// Vec t A
type TyVector struct {
	Term TypeValue
	Type TypeValue
}

// ======================================================================
// Methods required to satisfy the TypeValue interface
// ======================================================================

func (u Term) typeValue()        {}
func (u TyUnit) typeValue()      {}
func (u TyUnVar) typeValue()     {}
func (u TyForAll) typeValue()    {}
func (u TyExists) typeValue()    {}
func (u TyLambda) typeValue()    {}
func (u TyExVar) typeValue()     {}
func (u TySum) typeValue()       {}
func (u TyProduct) typeValue()   {}
func (u TyGuarded) typeValue()   {}
func (u TyAsserting) typeValue() {}

func (u Term) IsMonotype() bool        { return true }
func (u TyUnit) IsMonotype() bool      { return true }
func (u TyUnVar) IsMonotype() bool     { return true }
func (u TyExVar) IsMonotype() bool     { return true }
func (u TyForAll) IsMonotype() bool    { return false }
func (u TyExists) IsMonotype() bool    { return false }
func (u TyLambda) IsMonotype() bool    { return u.Domain.IsMonotype() && u.Codomain.IsMonotype() }
func (u TySum) IsMonotype() bool       { return u.Left.IsMonotype() && u.Right.IsMonotype() }
func (u TyProduct) IsMonotype() bool   { return u.Left.IsMonotype() && u.Right.IsMonotype() }
func (u TyGuarded) IsMonotype() bool   { return false }
func (u TyAsserting) IsMonotype() bool { return false }

// TODO record types

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
	case *TyExists:
		vb, ok := b.(*TyExists)
		return ok && va.Identifier == vb.Identifier && CompareTypeValues(va.Type, vb.Type)
	case *TySum:
		vb, ok := b.(*TySum)
		return ok && CompareTypeValues(va.Left, vb.Left) && CompareTypeValues(va.Right, vb.Right)
	case *TyLambda:
		vb, ok := b.(*TyLambda)
		return ok && CompareTypeValues(va.Domain, vb.Domain) && CompareTypeValues(va.Codomain, vb.Codomain)
	default:
		panic("FATAL. type comparison not implemented yet")
	}
}

// ======================================================================
// String representation of type values
// ======================================================================

func (t Term) String() string    { return fmt.Sprintf("%d", t.Num) }
func (u TyUnit) String() string  { return "unit" }
func (u TyExVar) String() string { return "∃'" + u.Identifier.String() }
func (u TyUnVar) String() string { return u.Identifier.String() }
func (u TyForAll) String() string {
	return fmt.Sprintf("∀%s.%s", u.Identifier.String(), u.Type.String())
}
func (u TyExists) String() string {
	return fmt.Sprintf("∀%s.%s", u.Identifier.String(), u.Type.String())
}
func (u TyLambda) String() string {
	return fmt.Sprintf("%s -> %s", u.Domain.String(), u.Codomain.String())
}
func (u TySum) String() string {
	return fmt.Sprintf("%s | %s", u.Left.String(), u.Right.String())
}
func (u TyProduct) String() string {
	return fmt.Sprintf("%s * %s", u.Left.String(), u.Right.String())
}
func (u *TyGuarded) String() string {
	return fmt.Sprintf("%s ⊃ %s", u.Guard.String(), u.Type.String())
}
func (u *TyAsserting) String() string {
	return fmt.Sprintf("%s ∧ %s", u.Type.String(), u.Guard.String())
}

// ======================================================================
// FullString representation of type values: also show UID numbers
// ======================================================================

func (t Term) FullString() string    { return t.String() }
func (u TyUnit) FullString() string  { return u.String() }
func (u TyUnVar) FullString() string { return u.Identifier.FullString() }
func (u TyExVar) FullString() string { return "∃'" + u.Identifier.FullString() }
func (u TyForAll) FullString() string {
	return fmt.Sprintf("∀%s.%s", u.Identifier.FullString(), u.Type.String())
}
func (u TyExists) FullString() string {
	return fmt.Sprintf("∀%s.%s", u.Identifier.FullString(), u.Type.String())
}
func (u TyLambda) FullString() string {
	return fmt.Sprintf("%s -> %s", u.Domain.FullString(), u.Codomain.FullString())
}
func (u TySum) FullString() string {
	return fmt.Sprintf("%s | %s", u.Left.FullString(), u.Right.FullString())
}
func (u TyProduct) FullString() string {
	return fmt.Sprintf("%s * %s", u.Left.FullString(), u.Right.FullString())
}
func (u *TyGuarded) FullString() string {
	return fmt.Sprintf("%s ⊃ %s", u.Guard.String(), u.Type.FullString())
}
func (u *TyAsserting) FullString() string {
	return fmt.Sprintf("%s ∧ %s", u.Type.FullString(), u.Guard.String())
}

// ======================================================================
// FancyString representation of type values: ocaml style like "'a -> 'b"
// ======================================================================

// helper for generating fancy type names in OCaml style
func genFancy(occ map[UniqueIdentifier]int, id UniqueIdentifier) string {
	if num, ok := occ[id]; ok {
		return string(rune(num + 97))
	}

	// FIXME generate decent names
	max := -1
	for _, v := range occ {
		if v > max {
			max = v
		}
	}

	occ[id] = max + 1
	return string(rune(max + 1 + 97))

}

func (u Term) FancyString(occ map[UniqueIdentifier]int) string   { return u.String() }
func (u TyUnit) FancyString(occ map[UniqueIdentifier]int) string { return "unit" }
func (u TyExVar) FancyString(occ map[UniqueIdentifier]int) string {
	return "'" + genFancy(occ, u.Identifier)
}
func (u TyUnVar) FancyString(occ map[UniqueIdentifier]int) string {
	return u.String()
}
func (u TyForAll) FancyString(occ map[UniqueIdentifier]int) string {
	return fmt.Sprintf("∀%s:%s.%s", genFancy(occ, u.Identifier), u.Sort.String(), u.Type.FancyString(occ))
}
func (u *TyExists) FancyString(occ map[UniqueIdentifier]int) string {
	return fmt.Sprintf("∃%s:%s.%s", genFancy(occ, u.Identifier), u.Sort.String(), u.Type.FancyString(occ))
}
func (u *TyLambda) FancyString(occ map[UniqueIdentifier]int) string {
	return fmt.Sprintf("%s -> %s", u.Domain.FancyString(occ), u.Codomain.FancyString(occ))
}
func (u *TySum) FancyString(occ map[UniqueIdentifier]int) string {
	return fmt.Sprintf("%s | %s", u.Left.FancyString(occ), u.Right.FancyString(occ))
}
func (u *TyProduct) FancyString(occ map[UniqueIdentifier]int) string {
	return fmt.Sprintf("%s * %s", u.Left.FancyString(occ), u.Right.FancyString(occ))
}
func (u *TyGuarded) FancyString(occ map[UniqueIdentifier]int) string {
	return fmt.Sprintf("%s ⊃ %s", u.Guard.String(), u.Type.FancyString(occ))
}
func (u *TyAsserting) FancyString(occ map[UniqueIdentifier]int) string {
	return fmt.Sprintf("%s ∧ %s", u.Type.FancyString(occ), u.Guard.String())
}
