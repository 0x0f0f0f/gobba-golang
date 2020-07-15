// Describes the structure of gobba's AST. Contains definitions
// of interfaces and Node types
package ast

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

// These interfaces contain dummy method
// but exist so they can be identified correctly by the go compiler
// type Statement interface {
// 	Node
// 	statementNode()
// }
//
type Expression interface {
	Node
	expressionNode()
	isPattern() bool
}

// Represents a symbol-value pair in the AST.
type Assignment struct {
	Token token.Token
	Name  *ExprIdentifier
	Value Expression
}

// Contains a list of assignments without a body
type LetStatement struct {
	Token       token.Token
	Assignments []*Assignment
}

// Represents `let var = val in body`
type ExprLet struct {
	Token      token.Token
	Assignment Assignment
	Body       Expression
}

// Represents an if-then-else expression
type ExprIf struct {
	Token       token.Token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

// Represents a function definition literal
type ExprLambda struct {
	Token token.Token
	Param *ExprIdentifier
	Body  Expression
}

// Represent a function application
type ExprApply struct {
	Token    token.Token
	Function Expression
	Arg      Expression
}

// Represents a prefix Expression
type ExprPrefix struct {
	Token    token.Token
	Operator string
	Right    Expression
}

// Represents an infix expression
type ExprInfix struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

// Represents a symbol or an identifier
type ExprIdentifier struct {
	Token      token.Token
	Identifier UniqueIdentifier
}

// Represents a type annotation
type ExprAnnot struct {
	Token token.Token
	Body  Expression
	Type  TypeValue
}

// Represents a fixed point combinator
type ExprFix struct {
	Token token.Token
	Param ExprIdentifier
	Body  Expression
}

type ExprPair struct {
	Token token.Token
	Left  Expression
	Right Expression
}

type ExprInj struct {
	Token   token.Token
	IsRight bool
	Expr    Expression
}

// Pattern matching branch
// type Branch {}

func (a *Assignment) expressionNode()     {}
func (ls *LetStatement) statementNode()   {}
func (le *ExprLet) expressionNode()       {}
func (i *ExprIf) expressionNode()         {}
func (f *ExprLambda) expressionNode()     {}
func (f *ExprApply) expressionNode()      {}
func (p *ExprPrefix) expressionNode()     {}
func (p *ExprInfix) expressionNode()      {}
func (i *ExprIdentifier) expressionNode() {}
func (i *ExprAnnot) expressionNode()      {}
func (i *ExprFix) expressionNode()        {}
func (i *ExprPair) expressionNode()       {}
func (i *ExprInj) expressionNode()        {}

func (a *Assignment) isPattern() bool     { return false }
func (le *ExprLet) isPattern() bool       { return false }
func (i *ExprIf) isPattern() bool         { return false }
func (f *ExprLambda) isPattern() bool     { return false }
func (f *ExprApply) isPattern() bool      { return false }
func (p *ExprPrefix) isPattern() bool     { return false }
func (p *ExprInfix) isPattern() bool      { return false }
func (i *ExprIdentifier) isPattern() bool { return true }
func (i *ExprAnnot) isPattern() bool      { return false }
func (i *ExprFix) isPattern() bool        { return false }
func (i *ExprPair) isPattern() bool       { return i.Left.isPattern() && i.Right.isPattern() }
func (i *ExprInj) isPattern() bool        { return i.Expr.isPattern() }

func (a *Assignment) TokenLiteral() string     { return a.Token.Literal }
func (ls *LetStatement) TokenLiteral() string  { return ls.Token.Literal }
func (le *ExprLet) TokenLiteral() string       { return le.Token.Literal }
func (i *ExprIf) TokenLiteral() string         { return i.Token.Literal }
func (f *ExprLambda) TokenLiteral() string     { return f.Token.Literal }
func (f *ExprApply) TokenLiteral() string      { return f.Token.Literal }
func (p *ExprPrefix) TokenLiteral() string     { return p.Token.Literal }
func (p *ExprInfix) TokenLiteral() string      { return p.Token.Literal }
func (i *ExprIdentifier) TokenLiteral() string { return i.Token.Literal }
func (i *ExprAnnot) TokenLiteral() string      { return i.Token.Literal }
func (i *ExprFix) TokenLiteral() string        { return i.Token.Literal }

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

func (i *IntegerLiteral) expressionNode() {}
func (f *FloatLiteral) expressionNode()   {}
func (c *ComplexLiteral) expressionNode() {}
func (c *BoolLiteral) expressionNode()    {}
func (c *UnitLiteral) expressionNode()    {}
func (c *StringLiteral) expressionNode()  {}
func (c *RuneLiteral) expressionNode()    {}

func (i *IntegerLiteral) isPattern() bool { return true }
func (f *FloatLiteral) isPattern() bool   { return true }
func (c *ComplexLiteral) isPattern() bool { return true }
func (c *BoolLiteral) isPattern() bool    { return true }
func (c *UnitLiteral) isPattern() bool    { return true }
func (c *StringLiteral) isPattern() bool  { return true }
func (c *RuneLiteral) isPattern() bool    { return true }

func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (f *FloatLiteral) TokenLiteral() string   { return f.Token.Literal }
func (c *ComplexLiteral) TokenLiteral() string { return c.Token.Literal }
func (c *BoolLiteral) TokenLiteral() string    { return c.Token.Literal }
func (c *UnitLiteral) TokenLiteral() string    { return c.Token.Literal }
func (c *StringLiteral) TokenLiteral() string  { return c.Token.Literal }
func (c *RuneLiteral) TokenLiteral() string    { return c.Token.Literal }

func (i *IntegerLiteral) String() string { return i.Token.Literal }
func (f *FloatLiteral) String() string   { return f.Token.Literal }
func (c *ComplexLiteral) String() string {
	return fmt.Sprintf("%g", c.Value)
}
func (c *BoolLiteral) String() string   { return c.Token.Literal }
func (c *UnitLiteral) String() string   { return "()" }
func (c *StringLiteral) String() string { return c.Token.Literal }
func (c *RuneLiteral) String() string   { return c.Token.Literal }
