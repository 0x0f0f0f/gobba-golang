// Describes the structure of gobba's AST. Contains definitions
// of interfaces and Node types
package ast

import (
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

// Represents e symbol-value pair in the AST.
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
//type ExprApply struct {
//	Token    token.Token
//	Function Expression
//	Arg      Expression
//}

// Spines are lists of arguments passed to a function application.
// passed to functions and are needed for the DK16 typechecker
type ExprApplySpine struct {
	Token    token.Token
	Function Expression
	Spine    []Expression
}

// Represents a prefix Expression
type ExprPrefix struct {
	Token    token.Token
	Operator Operator
	Right    Expression
}

type Operator struct {
	IsPattern bool // True if the operator is allowed in patterns
	Kind      string
}

func (o Operator) String() string {
	return string(o.Kind)
}

// Represents an infix expression
type ExprInfix struct {
	Token    token.Token
	Left     Expression
	Operator Operator
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
type ExprRec struct {
	Token token.Token
	Name  ExprIdentifier
	Body  Expression
}

type ExprPair struct {
	Token token.Token
	Left  Expression
	Right Expression
}

// Sum type injections inj_l : A -> A + B, inj_r : B -> A + B
type ExprInj struct {
	Token   token.Token
	IsRight bool
	Expr    Expression
}

// Pattern Matching branch
type MatchBranch struct {
	Patterns []Expression // TODO check that foreach isPattern() == true
	Body     Expression
}

// Pattern matching expression
type ExprMatch struct {
	Token    token.Token
	Expr     Expression
	Branches []MatchBranch
}

func (e *Assignment) expressionNode()     {}
func (e *LetStatement) statementNode()    {}
func (e *ExprLet) expressionNode()        {}
func (e *ExprIf) expressionNode()         {}
func (e *ExprLambda) expressionNode()     {}
func (e *ExprApplySpine) expressionNode() {}
func (e *ExprPrefix) expressionNode()     {}
func (e *ExprInfix) expressionNode()      {}
func (e *ExprIdentifier) expressionNode() {}
func (e *ExprAnnot) expressionNode()      {}
func (e *ExprRec) expressionNode()        {}
func (e *ExprPair) expressionNode()       {}
func (e *ExprInj) expressionNode()        {}
func (e *ExprMatch) expressionNode()      {}

func (e *ExprIdentifier) isPattern() bool { return true }
func (e *ExprPair) isPattern() bool       { return e.Left.isPattern() && e.Right.isPattern() }
func (e *ExprInj) isPattern() bool        { return e.Expr.isPattern() }
func (e *Assignment) isPattern() bool     { return false }
func (e *ExprLet) isPattern() bool        { return false }
func (e *ExprIf) isPattern() bool         { return false }
func (e *ExprLambda) isPattern() bool     { return false }
func (e *ExprApplySpine) isPattern() bool { return false }
func (e *ExprPrefix) isPattern() bool {
	if e.Operator.IsPattern {
		return e.Right.isPattern()
	}
	return false
}
func (e *ExprInfix) isPattern() bool {
	if e.Operator.IsPattern {
		return e.Left.isPattern() && e.Right.isPattern()
	}
	return false
}
func (e *ExprAnnot) isPattern() bool { return false }
func (e *ExprRec) isPattern() bool   { return false }
func (e *ExprMatch) isPattern() bool { return false }

func (e *Assignment) TokenLiteral() string     { return e.Token.Literal }
func (e *LetStatement) TokenLiteral() string   { return e.Token.Literal }
func (e *ExprLet) TokenLiteral() string        { return e.Token.Literal }
func (e *ExprIf) TokenLiteral() string         { return e.Token.Literal }
func (e *ExprLambda) TokenLiteral() string     { return e.Token.Literal }
func (e *ExprApplySpine) TokenLiteral() string { return e.Token.Literal }
func (e *ExprPrefix) TokenLiteral() string     { return e.Token.Literal }
func (e *ExprInfix) TokenLiteral() string      { return e.Token.Literal }
func (e *ExprIdentifier) TokenLiteral() string { return e.Token.Literal }
func (e *ExprAnnot) TokenLiteral() string      { return e.Token.Literal }
func (e *ExprRec) TokenLiteral() string        { return e.Token.Literal }
func (e *ExprInj) TokenLiteral() string        { return e.Token.Literal }
func (e *ExprMatch) TokenLiteral() string      { return e.Token.Literal }
