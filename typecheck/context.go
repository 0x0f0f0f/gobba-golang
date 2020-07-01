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
	Identifier ast.Identifier
}

func (v *UniversalVariable) contextValue() {}

// Unsolved existential variable α^
type ExistentialVariable struct {
	Identifier ast.Identifier
	Value      *ast.TypeValue
}

func (v *ExistentialVariable) contextValue() {}
