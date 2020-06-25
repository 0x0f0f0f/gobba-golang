// Provides utilities for parsing and printing gobba expressions
package parser

import (
    "text/scanner"
    "strconv"
    "github.com/alecthomas/participle"
    "github.com/alecthomas/participle/lexer"
    "github.com/alecthomas/repr"
)

type Expression struct {
    Binop     *BinopExpr     `  @@`
    Primitive *PrimitiveExpr `| @@`
}

type BinopInfo struct {
    RightAssociative bool
    Priority int
}

// Precedence climbing is an approach to parsing expressions that efficiently
// produces compact parse trees.
//
// In contrast, naive recursive descent expression parsers produce parse trees proportional in
// complexity to the number of operators supported. This impacts both readability and
// performance.
var binopPrecedence = map[string]BinopInfo {
    ">>": {Priority: 1},
    ">=>": {Priority: 2},
    "<=<": {Priority: 2},
    "=": {Priority: 3},
    ">": {Priority: 3},
    "<": {Priority: 3},
    ">=": {Priority: 3},
    "<=": {Priority: 3},
    "&&": {Priority: 4},
    "||": {Priority: 4},
    "+": {Priority: 5},
    "-": {Priority: 5},
    "*": {Priority: 6},
    "/": {Priority: 6},
    "^": {Priority: 7, RightAssociative: true},
}

// TODO priority
type BinopExpr struct {
    Terminal  *PrimitiveExpr
    Left      *BinopExpr
    Op        string
    Right     *BinopExpr
}

func (e *BinopExpr) Parse(lex *lexer.PeekingLexer) error {
	*e = *parseExpr(lex, 0)
	return nil
}

func parseExpr(lex *lexer.PeekingLexer, minPrec int) *BinopExpr {
	lhs := next(lex)
	for {
		op := peek(lex)
		if op == nil || binopPrecedence[op.Op].Priority < minPrec {
			break
		}
		nextMinPrec := binopPrecedence[op.Op].Priority
		if !binopPrecedence[op.Op].RightAssociative {
			nextMinPrec++
		}
		next(lex)
		rhs := parseExpr(lex, nextMinPrec)
		lhs = parseOp(op, lhs, rhs)
	}
	return lhs
}

func parseOp(op, lhs, rhs *BinopExpr) *BinopExpr {
	op.Left = lhs
	op.Right = rhs
	return op
}

func next(lex *lexer.PeekingLexer) *BinopExpr {
    e := peek(lex)
    if e == nil {
        return e
    }
    _, _ = lex.Next()
    switch e.Op {
    case "(":
        return next(lex)
    }
    return e
}

func peek(lex *lexer.PeekingLexer) *BinopExpr {
	t, err := lex.Peek(0)
	if err != nil {
		panic(err)
	}
	if t.EOF() {
		return nil
	}
	switch t.Type {
	case scanner.Int:
		n, err := strconv.ParseInt(t.Value, 10, 64)
		if err != nil {
			panic(err)
		}
		ni := int(n)
		return &BinopExpr{Terminal: &ni}

	case ')':
		_, _ = lex.Next()
		return nil

	default:
		return &BinopExpr{Op: t.Value}
	}
}


type PrimitiveExpr struct {
    Integer *int64       `  @Int`
    Float   *float64     `| @Float`
    Boolean bool         `| ( @"true" | "false" )`
    String  *string      `| @String`
    Unit    bool         `| @"nil"`
    Ident   *string      `| @Ident`
    SubExpression *Expression `| "(" @@ ")"`
}

type Parser struct {
    parser *participle.Parser
}

func NewParser() (*Parser, error) {
    parser, err := participle.Build(&Expression{})
    if err != nil {
        return nil, err 
    }
    p := &Parser{
        parser: parser,
    }
    return p, nil
}

func (p *Parser) ParseString(s string) (*Expression, error) {
    ast := &Expression{}
    err := p.parser.ParseString(s, ast)
    if err != nil {
        return nil, err
    }
    return ast, nil
}

func PrintExpression(e *Expression) {
   repr.Println(e) 
} 
