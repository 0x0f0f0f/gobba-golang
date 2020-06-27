// Describes the structure of gobba's AST. Contains definitions
// of interfaces and Node types
package ast

import (
	"bytes"
	"github.com/0x0f0f0f/gobba-golang/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

// These interfaces contain dummy method
// but exist so they can be identified correctly by the go compiler
type Statement interface {
	Node
	statementNode()
}
type Expression interface {
	Node
	expressionNode()
}
type Program struct {
	Statements []Statement
}

// Get the first literal from a program. This is needed
// so that the Program struct implements the Node interface
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string {
	var b bytes.Buffer

	for _, s := range p.Statements {
		b.WriteString(s.String())
	}

	return b.String()
}

// ======================================================================
// AST nodes types definitions
// ======================================================================

// Contains a list of assignments without a body
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// Represents a symbol-value pair in the AST.
type Assignment struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (a *Assignment) expressionNode()      {}
func (a *Assignment) TokenLiteral() string { return a.Token.Literal }
func (a *Assignment) String() string {
	var b bytes.Buffer

	b.WriteString(a.Name.String() + " = ")
	b.WriteString(a.Value.String())

	return b.String()
}

// Contains a list of assignments without a body
type LetStatement struct {
	Token       token.Token
	Assignments []*Assignment
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var b bytes.Buffer

	b.WriteString(ls.TokenLiteral() + " ")
	for i, ass := range ls.Assignments {
		b.WriteString(ass.String())
		if i < len(ls.Assignments)-1 {
			b.WriteString(" and ")
		}
	}
	b.WriteString(";")

	return b.String()
}

// Represents `let var = val in body`
type LetExpression struct {
	Token       token.Token
	Assignments []Assignment
	Body        Expression
}

func (le *LetExpression) expressionNode()      {}
func (le *LetExpression) TokenLiteral() string { return le.Token.Literal }
func (le *LetExpression) String() string {
	var b bytes.Buffer

	b.WriteString(le.TokenLiteral() + " ")
	for i, ass := range le.Assignments {
		b.WriteString(ass.String())
		if i < len(le.Assignments)-1 {
			b.WriteString(" and ")
		}
	}
	b.WriteString(" in ")
	b.WriteString(le.Body.String())

	return b.String()
}

// Represents a prefix Expression
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }
func (p *PrefixExpression) String() string {
	var b bytes.Buffer
	b.WriteString("(")
	b.WriteString(p.Operator)
	b.WriteString(p.Right.String())
	b.WriteString(")")
	return b.String()
}

// Represents an infix expression
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (p *InfixExpression) expressionNode()      {}
func (p *InfixExpression) TokenLiteral() string { return p.Token.Literal }
func (p *InfixExpression) String() string {
	var b bytes.Buffer
	b.WriteString("(")
	b.WriteString(p.Left.String())
	b.WriteString(" " + p.Operator + " ")
	b.WriteString(p.Right.String())
	b.WriteString(")")
	return b.String()
}

// Represents a symbol or an identifier
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// Represents an integer literal
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }

// Represents a floating point literal
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (f *FloatLiteral) expressionNode()      {}
func (f *FloatLiteral) TokenLiteral() string { return f.Token.Literal }
func (f *FloatLiteral) String() string       { return f.Token.Literal }

// Represents a complex number literal
type ComplexLiteral struct {
	Token token.Token
	Value complex128
}

func (c *ComplexLiteral) expressionNode()      {}
func (c *ComplexLiteral) TokenLiteral() string { return c.Token.Literal }
func (c *ComplexLiteral) String() string       { return c.Token.Literal }

// Represents an imaginary part number literal
type ImagLiteral struct {
	Token token.Token
	Value float64
}

func (c *ImagLiteral) expressionNode()      {}
func (c *ImagLiteral) TokenLiteral() string { return c.Token.Literal }
func (c *ImagLiteral) String() string       { return c.Token.Literal }
