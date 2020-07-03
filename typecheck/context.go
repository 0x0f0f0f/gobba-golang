package typecheck

import (
	"github.com/0x0f0f0f/gobba-golang/ast"
)

// This file contains definition for algorithmic type contexts
// defined in "Complete and Easy Bidirectional Typechecking
// for Higher-Rank Polymorphism"
// https://www.cl.cam.ac.uk/~nk480/bidir.pdf
// https://github.com/chrisnevers/bidirectional-typechecking/blob/master/lib/ast/context.ml

type ContextValue interface {
	contextValue()
}

// ======================================================================
// Definitions of types of values that compose an algorithmic type context
// ======================================================================

// Denoted with α in the paper
type UniversalVariable struct {
	Identifier ast.UniqueIdentifier
}

func (v *UniversalVariable) contextValue() {}

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

// Denoted with |>α^ in the paper
type Marker struct {
	Identifier ast.UniqueIdentifier
}

func (v *Marker) contextValue() {}

// Denoted with |>α^ in the paper
type TypedVariable struct {
	Identifier ast.UniqueIdentifier
	Value      ast.TypeValue
}

func (v *TypedVariable) contextValue() {}

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

// Sorted insertion after element el in the context
// Return a new context after insertion
func (c Context) Insert(el ContextValue, values []ContextValue) Context {
	nc := NewContext()

	i := len(c.Contents)
	for j, s := range c.Contents {
		if s == el {
			i = j
			break
		}
	}

	nc.Contents = append(nc.Contents, c.Contents[:i]...)
	nc.Contents = append(nc.Contents, values...)
	nc.Contents = append(nc.Contents, c.Contents[i:]...)

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
		if old != el {
			nc.Contents = append(nc.Contents, old)
		}
	}
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

// Split a context in two left and right context when a value is encountered
func (c Context) SplitAt(el ContextValue) (Context, Context) {
	left := NewContext()
	right := NewContext()
	found := false
	for _, old := range c.Contents {
		if found {
			right.Contents = append(right.Contents, old)
			continue
		}

		left.Contents = append(left.Contents, old)
		if old == el {
			found = true
		}
	}
	return *left, *right
}
