package typecheck

import (
	"bytes"
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// This file contains definition for algorithmic type contexts
// defined in "Complete and Easy Bidirectional Typechecking for Higher-Rank Polymorphism"
// https://www.cl.cam.ac.uk/~nk480/bidir.pdf
// https://github.com/chrisnevers/bidirectional-typechecking/blob/master/lib/ast/context.ml

type ContextValue interface {
	contextValue()
	String() string
}

// ======================================================================
// Definitions of types of values that compose an algorithmic type context
// ======================================================================

// Denoted with α in the paper
type UniversalVariable struct {
	Identifier ast.UniqueIdentifier
}

func (v *UniversalVariable) contextValue() {}
func (v *UniversalVariable) String() string {
	return v.Identifier.FullString()
}

// Unsolved existential variable α^
// solved when Value is not nil
type ExistentialVariable struct {
	Identifier ast.UniqueIdentifier
	Value      *ast.TypeValue
}

func (v *ExistentialVariable) contextValue() {}
func (v *ExistentialVariable) solved() bool {
	return v.Value != nil
}
func (v *ExistentialVariable) String() string {
	var b bytes.Buffer

	b.WriteString(v.Identifier.FullString())

	if v.Value != nil {
		b.WriteString("=")
		b.WriteString((*v.Value).FullString())
	}
	return b.String()
}

// Denoted with |>α^ in the paper
type Marker struct {
	Identifier ast.UniqueIdentifier
}

func (v *Marker) contextValue() {}
func (v *Marker) String() string {
	return "►" + v.Identifier.FullString()
}

// Denoted with x : A in the paper
type TypeAnnotation struct {
	Identifier ast.UniqueIdentifier
	Value      ast.TypeValue
}

func (v *TypeAnnotation) contextValue() {}
func (v *TypeAnnotation) String() string {
	return v.Identifier.FullString() + ": " + v.Value.FullString()
}

// Returns true if two values implementing ContextValue are equal
func CompareContextValues(a, b ContextValue) bool {
	switch va := a.(type) {
	case *UniversalVariable:
		if vb, ok := b.(*UniversalVariable); ok {
			return *va == *vb
		}
	case *ExistentialVariable:
		if vb, ok := b.(*ExistentialVariable); ok {
			return *va == *vb
		}
	case *Marker:
		if vb, ok := b.(*Marker); ok {
			return *va == *vb
		}
	case *TypeAnnotation:
		if vb, ok := b.(*TypeAnnotation); ok {
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
	Contents []ContextValue
}

func NewContext() *Context {
	c := &Context{}
	c.Contents = make([]ContextValue, 0)
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

// Sorted insertion after element el in the context
// Return a new context after insertion
func (c Context) Insert(el ContextValue, values []ContextValue) Context {
	nc := NewContext()

	i := len(c.Contents)
	for j, s := range c.Contents {
		if CompareContextValues(s, el) {
			i = j
			break
		}
	}

	nc.Contents = append(nc.Contents, c.Contents[:i]...)
	nc.Contents = append(nc.Contents, values...)
	if i < len(c.Contents) {
		nc.Contents = append(nc.Contents, c.Contents[i+1:]...)
	}

	return *nc
}

// Insert at head and return a new context
func (c Context) InsertHead(el ContextValue) Context {
	nc := NewContext()
	nc.Contents = append(nc.Contents, el)
	nc.Contents = append(nc.Contents, c.Contents...)
	return *nc
}

// Remove an element from a context and return a new one
func (c Context) Drop(el ContextValue) Context {
	nc := NewContext()
	for _, old := range c.Contents {
		if !CompareContextValues(old, el) {
			nc.Contents = append(nc.Contents, old)
		}
	}
	return *nc
}

func (c Context) Concat(rc Context) Context {
	nc := NewContext()
	copy(nc.Contents, c.Contents)
	nc.Contents = append(nc.Contents, rc.Contents...)

	return *nc
}

// True if the context contain the universal variable with the given identifier
func (c Context) HasUniversalVariable(alpha ast.UniqueIdentifier) bool {
	for _, c := range c.Contents {
		if v, ok := c.(*UniversalVariable); ok {
			if v.Identifier.Value == alpha.Value {
				return true
			}
		}
	}
	return false
}

// True if the context contains an unsolved universal variable with the given identifier
func (c Context) HasExistentialVariable(alpha ast.UniqueIdentifier) bool {
	for _, c := range c.Contents {
		if v, ok := c.(*ExistentialVariable); ok {
			if v.Identifier == alpha {
				return !v.solved()
			}
		}
	}
	return false
}

// If the context contains a solved universal variable with the given identifier
// just return the corresponding monotype
func (c Context) GetSolvedVariable(alpha ast.UniqueIdentifier) *ast.TypeValue {
	for _, c := range c.Contents {
		if v, ok := c.(*ExistentialVariable); ok {
			if v.Identifier == alpha && v.solved() {
				return v.Value
			}
		}
	}
	return nil
}

// Return the type of a type annotation.
func (c Context) GetAnnotation(alpha ast.UniqueIdentifier) *ast.TypeValue {
	for _, c := range c.Contents {
		if v, ok := c.(*TypeAnnotation); ok {
			if v.Identifier == alpha {
				return &v.Value
			}
		}
	}
	return nil
}

// Split a context in two left and right context when a value is encountered
func (c Context) SplitAt(el ContextValue) (Context, Context) {
	left := NewContext()
	right := NewContext()
	found := false
	for _, old := range c.Contents {
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
