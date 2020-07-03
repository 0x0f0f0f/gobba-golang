package parser

import (
	"fmt"
	"github.com/0x0f0f0f/gobba-golang/token"
	"runtime/debug"
)

// This file contains parsing error methods and definitions

// ======================================================================

type ParserError struct {
	Line, Column int
	Token        token.Token
	Expected     *token.TokenType
	Msg          string
}

func (pe ParserError) Error() string {
	s := fmt.Sprintf("syntax error at line %d column %d: ", pe.Line, pe.Column)
	s += fmt.Sprintf("unexpected token %s. ", pe.Token.Literal)
	if pe.Expected != nil {
		s += fmt.Sprintf("expected %s. ", *pe.Expected)
	}
	if len(pe.Msg) > 0 {
		s += fmt.Sprintf("\n\t %s", pe.Msg)
	}
	return s
}

func (p *Parser) appendError(err ParserError) {
	if p.TraceOnError {
		fmt.Printf("ERRORING HERE: %s\n", err)
		debug.PrintStack()
	}
	p.errors = append(p.errors, err)
}

func (p *Parser) Errors() []string {
	s := make([]string, len(p.errors))
	for i, err := range p.errors {
		s[i] = err.Error()
	}
	return s
}

// Add an error when a peekToken is not the expected one
func (p *Parser) peekError(expected token.TokenType, t token.Token) {
	e := ParserError{t.Line, t.Column, t, &expected, ""}
	p.appendError(e)
}

// Add a custom error
func (p *Parser) customError(expected *token.TokenType, t token.Token, msg string) {
	e := ParserError{t.Line, t.Column, t, expected, msg}
	p.appendError(e)
}

func (p *Parser) noPrefixParseFnError(t token.Token) {
	e := ParserError{t.Line, t.Column, t, nil, ""}

	p.errors = append(p.errors, e)
}

// ======================================================================
