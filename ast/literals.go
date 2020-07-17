package ast

import (
	"fmt"

	"github.com/0x0f0f0f/gobba-golang/token"
)

// ======================================================================
// Terminal values: literals
// ======================================================================

// Represents an integer literal
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

// Represents a floating point literal
type FloatLiteral struct {
	Token token.Token
	Value float64
}

// Represents a complex number literal
type ComplexLiteral struct {
	Token token.Token
	Value complex128
}

// Represents a boolean value
type BoolLiteral struct {
	Token token.Token
	Value bool
}

// Represents an unit value
type UnitLiteral struct {
	Token token.Token
}

// Represents an empty vector literal
type EmptyVecLiteral struct {
	Token token.Token
}

// Represents a string value
type StringLiteral struct {
	Token token.Token
	Value string
}

// Represents an Unicode character value
type RuneLiteral struct {
	Token token.Token
	Value string
}

func (e *IntegerLiteral) expressionNode()  {}
func (e *FloatLiteral) expressionNode()    {}
func (e *ComplexLiteral) expressionNode()  {}
func (e *BoolLiteral) expressionNode()     {}
func (e *UnitLiteral) expressionNode()     {}
func (e *EmptyVecLiteral) expressionNode() {}
func (e *StringLiteral) expressionNode()   {}
func (e *RuneLiteral) expressionNode()     {}

func (e *IntegerLiteral) isPattern() bool  { return true }
func (e *FloatLiteral) isPattern() bool    { return true }
func (e *ComplexLiteral) isPattern() bool  { return true }
func (e *BoolLiteral) isPattern() bool     { return true }
func (e *UnitLiteral) isPattern() bool     { return true }
func (e *EmptyVecLiteral) isPattern() bool { return true }
func (e *StringLiteral) isPattern() bool   { return true }
func (e *RuneLiteral) isPattern() bool     { return true }

func (e *IntegerLiteral) TokenLiteral() string  { return e.Token.Literal }
func (e *FloatLiteral) TokenLiteral() string    { return e.Token.Literal }
func (e *ComplexLiteral) TokenLiteral() string  { return e.Token.Literal }
func (e *BoolLiteral) TokenLiteral() string     { return e.Token.Literal }
func (e *UnitLiteral) TokenLiteral() string     { return e.Token.Literal }
func (e *EmptyVecLiteral) TokenLiteral() string { return e.Token.Literal }
func (e *StringLiteral) TokenLiteral() string   { return e.Token.Literal }
func (e *RuneLiteral) TokenLiteral() string     { return e.Token.Literal }

func (e *IntegerLiteral) String() string  { return fmt.Sprintf("%d", e.Value) }
func (e *FloatLiteral) String() string    { return fmt.Sprintf("%g", e.Value) }
func (e *ComplexLiteral) String() string  { return fmt.Sprintf("%g", e.Value) }
func (e *BoolLiteral) String() string     { return fmt.Sprintf("%t", e.Value) }
func (e *UnitLiteral) String() string     { return token.UNIT }
func (e *EmptyVecLiteral) String() string { return token.EMPTYVEC }
func (e *StringLiteral) String() string   { return e.Token.Literal }
func (e *RuneLiteral) String() string     { return e.Token.Literal }
