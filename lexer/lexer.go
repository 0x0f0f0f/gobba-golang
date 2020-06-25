package lexer

import (
	"github.com/0x0f0f0f/gobba-golang/token"
)

// Represent the current state of a lexer
// TODO move to UTF8
type Lexer struct {
	input        string
	position     int  // Current position in input
	readPosition int  // Current reading position
	ch           byte // Current char
}

// Creates a new lexer on a given input
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Advances the current character in input
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// Creates a token from a type and a character
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// Checks if a character is a letter
// TODO alphanumeric identifiers
func isIdentifier(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// Reads an identifier
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isIdentifier(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// Read the next character without incrementing position
func (l *Lexer) peekChar() byte {
    if l.readPosition >= len(l.input) {
        return 0
    } else {
        return l.input[l.readPosition]
    }
}

// Scans and return a token, advancing by a rune
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	// Single character operators
	case '+':
        	if l.peekChar() == '+' {
                	ch := l.ch
                	l.readChar()
                	tok = token.Token{
                        	Type: token.CONCAT,
                        	Literal: string(ch) + string(l.ch),
                    	}
        	} else {
        		tok = newToken(token.PLUS, l.ch)
            	}
	case '-':
        	if l.peekChar() == '>' {
                	ch := l.ch
                	l.readChar()
                	tok = token.Token{
                        	Type: token.LARROW,
                        	Literal: string(ch) + string(l.ch),
                	}
        	} else {
        		tok = newToken(token.MINUS, l.ch)
        	}
	case '*':
		tok = newToken(token.TIMES, l.ch)
	case '^':
		tok = newToken(token.TOPOW, l.ch)
	case '=':
		tok = newToken(token.EQUALS, l.ch)
	case '%':
		tok = newToken(token.MODULO, l.ch)
	case '@':
		tok = newToken(token.AT, l.ch)
	case '$':
		tok = newToken(token.DOLLAR, l.ch)
	case '!':
        	if l.peekChar() == '=' {
                	ch := l.ch
                	l.readChar()
                	tok = token.Token{
                        	Type: token.DIFFERS,
                        	Literal: string(ch) + string(l.ch),
                	}
        	} else {
                	tok = newToken(token.ILLEGAL, l.ch)
        	}
	case '&':
        	if l.peekChar() == '&' {
                	ch := l.ch
                	l.readChar()
                	tok = token.Token{
                        	Type: token.LAND,                       	Literal: string(ch) + string(l.ch),
                	}
        	} else {
                	tok = newToken(token.ILLEGAL, l.ch)
        	}
	case '|':
        	if l.peekChar() == '|' {
                	ch := l.ch
                	l.readChar()
                	tok = token.Token{
                        	Type: token.LAND,                       	Literal: string(ch) + string(l.ch),
                	}
        	} else {
                	tok = newToken(token.ILLEGAL, l.ch)
        	}
	case ':':
        	if l.peekChar() == '+' {
                	ch := l.ch
                	l.readChar()
                	tok = token.Token{
                        	Type: token.COMPLEX,
                        	Literal: string(ch) + string(l.ch),
                	}
        	} else if l.peekChar() == ':' {
                	ch := l.ch
                	l.readChar()
                	tok = token.Token{
                        	Type: token.CONS,
                        	Literal: string(ch) + string(l.ch),
                	}
        	} else {
                	tok = newToken(token.ACCESS, l.ch)
        	}
	case '<':
        	if l.peekChar() == '=' {
                	ch := l.ch
                	l.readChar()
                	if l.peekChar() == '<' {
                        	ch2 := l.ch
                        	l.readChar()
                        	tok = token.Token{
                                	Type: token.COMPOSE,
                                	Literal: (string(ch) + string(ch2) + string(l.ch)),
                        	}
                	} else {
                        	tok = token.Token{
                                	Type: token.LESSEQ,
                                	Literal: string(ch) + string(l.ch),
                        	}
                	}
        	} else {
        		tok = newToken(token.LESS, l.ch)
            	
        	}
	case '>':
        	if l.peekChar() == '=' {
                	ch := l.ch
                	l.readChar()
                	if l.peekChar() == '>' {
                        	ch2 := l.ch
                        	l.readChar()
                        	tok = token.Token{
                                	Type: token.PIPE,
                                	Literal: string(ch) + string(ch2) + string(l.ch),
                        	}
                	} else if l.peekChar() == '>' {
                        	tok = token.Token{
                                	Type: token.BIND,
                                	Literal: string(ch) + string(l.ch),
                        	}
                	} else {
                        	tok = token.Token{
                                	Type: token.GREATEREQ,
                                	Literal: string(ch) + string(l.ch),
                        	}
                	}
        	} else {
        		tok = newToken(token.LESS, l.ch)
            	
        	}
	case ';':
		tok = newToken(token.SEMI, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	// Now the next token must either be and identifier
	// Or an invalid token
	default:
		if isIdentifier(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	// Advance by a character
	l.readChar()
	return tok
}
