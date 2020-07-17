package typecheck

import (
	"bytes"

	"github.com/0x0f0f0f/gobba-golang/ast"
	// "reflect"
)

// This file contains definition for algorithmic type contexts
// defined in "Complete and Easy Bidirectional Typechecking for Higher-Rank Polymorphism"
// https://www.cl.cam.ac.uk/~nk480/bidir.pdf
// https://github.com/chrisnevers/bidirectional-typechecking/blob/master/lib/ast/context.ml

type CtxValue interface {
	contextValue()
	String() string
}

// ======================================================================
// Definitions of types of values that compose an algorithmic type context
// ======================================================================

// Universal variable sort. Denoted with α:κ in the paper
type CtxUnSort struct {
	Identifier ast.UniqueIdentifier
	Sort       ast.Sort
}

// Existential variable sort. α^:κ
type CtxExSort struct {
	Identifier ast.UniqueIdentifier
	Sort       ast.Sort
}

// Universal variable equality judgment α = t
type CtxUnEq struct {
	Identifier ast.UniqueIdentifier
	Term       ast.TypeValue
}

// Existential variable equality. α^:κ = τ
type CtxExEq struct {
	Identifier ast.UniqueIdentifier
	Sort       ast.Sort
	Value      ast.TypeValue
}

// Denoted with |>α in the paper
type CtxUnMarker struct {
	Identifier ast.UniqueIdentifier
}

// Denoted with |>α^ in the paper
type CtxExMarker struct {
	Identifier ast.UniqueIdentifier
}

// Expressions variable typings. x:Ap
type CtxVarType struct {
	Identifier   ast.UniqueIdentifier
	Type         ast.TypeValue
	Principality ast.Principality
}

func (v *CtxUnSort) contextValue()   {}
func (v *CtxExSort) contextValue()   {}
func (v *CtxUnEq) contextValue()     {}
func (v *CtxExEq) contextValue()     {}
func (v *CtxUnMarker) contextValue() {}
func (v *CtxExMarker) contextValue() {}
func (v *CtxVarType) contextValue()  {}

func (v *CtxUnSort) String() string {
	return v.Identifier.FullString() + ": " + v.Sort.String()
}
func (v *CtxExSort) String() string {
	return v.Identifier.FullString() + "^: " + v.Sort.String()
}

func (v *CtxUnEq) String() string {
	return v.Identifier.FullString() + "= " + v.Term.String()
}
func (v *CtxExEq) String() string {
	return v.Identifier.FullString() + "^: " + v.Sort.String() + "= " + v.Value.String()
}

func (v *CtxUnMarker) String() string {
	return "►" + v.Identifier.FullString()
}
func (v *CtxExMarker) String() string {
	return "►" + v.Identifier.FullString() + "^"
}

func (v *CtxVarType) String() string {
	return v.Identifier.FullString() + ": " + v.Type.FullString() + v.Principality.String()
}

// Returns true if two values implementing ContextValue are equal
func CompareContextValues(a, b CtxValue) bool {
	switch va := a.(type) {
	case *CtxUnSort:
		if vb, ok := b.(*CtxUnSort); ok {
			return *va == *vb
		}
	case *CtxExSort:
		if vb, ok := b.(*CtxExSort); ok {
			return *va == *vb
		}
	case *CtxUnEq:
		if vb, ok := b.(*CtxUnEq); ok {
			return *va == *vb
		}
	case *CtxExEq:
		if vb, ok := b.(*CtxExEq); ok {
			return *va == *vb
		}
	case *CtxUnMarker:
		if vb, ok := b.(*CtxUnMarker); ok {
			return *va == *vb
		}
	case *CtxExMarker:
		if vb, ok := b.(*CtxExMarker); ok {
			return *va == *vb
		}
	case *CtxVarType:
		if vb, ok := b.(*CtxVarType); ok {
			return *va == *vb
		}

	}

	return false
}

// ======================================================================
// Algorithmic Context Type: Γ, ∆, Θ
// Complete Contexts Ω do not contain unsolved variables
// ======================================================================

type Context struct {
	Contents []CtxValue
}

// Creates a new empty context
func NewContext() *Context {
	c := &Context{}
	c.Contents = make([]CtxValue, 0)
	return c
}

func (c Context) String() string {
	var b bytes.Buffer

	b.WriteString("[")

	for _, v := range c.Contents {
		b.WriteString(v.String())
		b.WriteString(", ")
	}
	b.WriteString("]")
	return b.String()
}

// ======================================================================
// Algorithmic contex insertion and deletion
// ======================================================================

// Sorted insertion in place of  element el in the context
// Return a new context after insertion
func (Γ Context) Insert(el CtxValue, values []CtxValue) Context {
	nc := NewContext()

	i := len(Γ.Contents)
	for j, s := range Γ.Contents {
		if CompareContextValues(s, el) {
			i = j
			break
		}
	}

	nc.Contents = append(nc.Contents, Γ.Contents[:i]...)
	nc.Contents = append(nc.Contents, values...)
	if i < len(Γ.Contents) {
		nc.Contents = append(nc.Contents, Γ.Contents[i+1:]...)
	}

	return *nc
}

// Insert at head and return a new context
func (Γ Context) InsertHead(el CtxValue) Context {
	nc := NewContext()
	nc.Contents = append(nc.Contents, el)
	nc.Contents = append(nc.Contents, Γ.Contents...)
	return *nc
}

// Remove an element from a context and return a new one
func (Γ Context) Drop(el CtxValue) Context {
	nc := NewContext()
	for _, old := range Γ.Contents {
		if !CompareContextValues(old, el) {
			nc.Contents = append(nc.Contents, old)
		}
	}
	return *nc
}

func (Γ Context) Concat(rc Context) Context {
	nc := NewContext()
	copy(nc.Contents, Γ.Contents)
	nc.Contents = append(nc.Contents, rc.Contents...)

	return *nc
}

// Split a context in two left and right context when a value is encountered
func (Γ Context) SplitAt(el CtxValue) (Context, Context) {
	left := NewContext()
	right := NewContext()
	found := false
	for _, old := range Γ.Contents {
		if CompareContextValues(old, el) {
			found = true
		}

		if found {
			right.Contents = append(right.Contents, old)
			continue
		}

		left.Contents = append(left.Contents, old)
	}
	return *left, *right
}

// ======================================================================
// Retrieving values from a Context
// ======================================================================

// Used in [Γ]α
func (Γ Context) GetUnEq(α ast.UniqueIdentifier) ast.TypeValue {
	for _, el := range Γ.Contents {
		if uneq, ok := el.(*CtxUnEq); ok {
			if uneq.Identifier == α {
				return uneq.Term
			}
		}
	}
	return nil
}

// ======================================================================
// Context as a substitution
// ======================================================================

// Apply a context as a substitution to a solved existential variable.
func (Γ Context) ApplyTypeValue(α ast.TypeValue) ast.TypeValue {
	Γ.debugSection("apply", α.FullString(), "=", α.FullString())
	switch va := α.(type) {
	// [Γ]α
	case *ast.TyUnVar:
		if τ := Γ.GetUnEq(va.Identifier); τ != nil {
			return τ
		} else {
			return α
		}
	// [Γ](P ⊃ A)
	case *ast.TyGuarded:
		return &ast.TyGuarded{
			Type:  Γ.ApplyTypeValue(va.Type),
			Guard: Γ.ApplyProp(va.Guard),
		}
	// [Γ](A ∧ P)
	case *ast.TyAsserting:
		return &ast.TyAsserting{
			Type:  Γ.ApplyTypeValue(va.Type),
			Guard: Γ.ApplyProp(va.Guard),
		}
	// [Γ](A -> B)
	case *ast.TyLambda:
		return &ast.TyLambda{
			Domain:   Γ.ApplyTypeValue(va.Domain),
			Codomain: Γ.ApplyTypeValue(va.Codomain),
		}
	// [Γ](A + B)
	case *ast.TySum:
		return &ast.TySum{
			Left:  Γ.ApplyTypeValue(va.Left),
			Right: Γ.ApplyTypeValue(va.Right),
		}
	// [Γ](A * B)
	case *ast.TyProduct:
		return &ast.TyProduct{
			Left:  Γ.ApplyTypeValue(va.Left),
			Right: Γ.ApplyTypeValue(va.Right),
		}
		// TODO case vec
	}
	return α
}

func (Γ Context) ApplyProp(p ast.Prop) ast.Prop {
	return ast.Prop{
		Left:  Γ.ApplyTypeValue(p.Left),
		Right: Γ.ApplyTypeValue(p.Left),
	}
}
