package lexer

import (
	"github.com/0x0f0f0f/gobba-golang/token"
)

// Represent the current state of a lexer
// TODO move to UTF8
type Lexer struct {
	input        string
	line         int  // Current line in the file
	column       int  // Current column in the line
	position     int  // Current position in input
	readPosition int  // Current reading position
	ch           byte // Current char
}

// Creates a new lexer on a given input
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	l.line = 1
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
	if l.ch == '\n' {
		l.column = 0
		l.line++
	} else {
		l.column++
	}
}

// Creates a token from a type and a character
func (l *Lexer) newToken(tokenType token.TokenType, lit string) token.Token {
	return token.Token{
		Type:    tokenType,
		Line:    l.line,
		Column:  l.column,
		Literal: lit,
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	if l.ch == '/' {
		// skip a line comment
		if l.peekChar() == '/' {
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
		}
	}
}

// Checks if a character is a letter
func isIdentifier(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// Reads an identifier
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isIdentifier(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isCharAllowedInNumber(ch byte) bool {
	return isDigit(ch) || ch == '-' || ch == '+' || ch == 'e' || ch == 'i' || ch == '.'
}

// Examples of valid numbers
// 12
// 12.2
// 12e-1
// 12e+1
// 12e1
// 12i
// 12.2i
// 12e-1i
// 12e+1i
// 12e1i
func (l *Lexer) readNumber() (token.TokenType, string) {
	position := l.position
	var kind token.TokenType = token.INT
	hasDot := false
	hasExponent := false

	for isCharAllowedInNumber(l.ch) {
		// End of the expression. No exponent
		if l.ch == '+' || l.ch == '-' {
			return kind, string(l.input[position:l.position])
		}

		// If the number contains e we're in scientific notation
		if l.ch == 'e' {
			if hasExponent {
				return token.ILLEGAL, string(l.input[position:l.position])
			}

			hasExponent = true
			peeked := l.peekChar()
			if peeked == '+' || peeked == '-' {
				l.readChar()
				peeked = l.peekChar()
			}
			if !isDigit(peeked) {
				return token.ILLEGAL, string(l.input[position:l.position])
			}
			kind = token.FLOAT
		}

		// If the number ends with an 'i', it is an imaginary part number
		if l.ch == 'i' {
			kind = token.IMAG
			l.readChar()
			return kind, string(l.input[position:l.position])
		}

		// If the number contains two dots, that's a problem
		// Floating point exponents are not allowed.
		if l.ch == '.' {
			if !isDigit(l.peekChar()) {
				return token.ILLEGAL, string(l.input[position:l.position])
			}

			if hasDot || hasExponent {
				return token.ILLEGAL, string(l.input[position:l.position])
			}
			hasDot = true
			kind = token.FLOAT
		}
		l.readChar()
	}

	// The number should not end with the exponent character
	if l.input[l.position-1] == 'e' {
		return token.ILLEGAL, string(l.input[position:l.position])
	}

	return kind, l.input[position:l.position]
}

// Read the next character without incrementing position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// Read a string and return either a valid token
// or an "ILLEGAL" token if string is terminated early
func (l *Lexer) readString() (token.TokenType, string) {
	position := l.position + 1
	var kind token.TokenType = token.STRING
	for {
		l.readChar()
		if l.ch == 0 {
			kind = token.ILLEGAL
			break
		}
		if l.ch == '"' {
			break
		}
	}
	return kind, l.input[position:l.position]
}

// Scans and return a token, advancing by a rune
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()
	l.skipComment()
	l.skipWhitespace()

	switch l.ch {
	case '"':
		kind, literal := l.readString()
		tok = l.newToken(kind, literal)
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.CONCAT, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.PLUS, string(l.ch))
		}
	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.LARROW, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.MINUS, string(l.ch))
		}
	case '(':
		tok = l.newToken(token.LPAREN, string(l.ch))
	case ')':
		tok = l.newToken(token.RPAREN, string(l.ch))
	case '*':
		tok = l.newToken(token.TIMES, string(l.ch))
	case '/':
		tok = l.newToken(token.DIVIDE, string(l.ch))
	case '^':
		tok = l.newToken(token.TOPOW, string(l.ch))
	case '=':
		tok = l.newToken(token.EQUALS, string(l.ch))
	case '%':
		tok = l.newToken(token.MODULO, string(l.ch))
	case '@':
		tok = l.newToken(token.AT, string(l.ch))
	case '$':
		tok = l.newToken(token.DOLLAR, string(l.ch))
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.DIFFERS, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.ch))
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.LAND, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.ch))
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.OR, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.ch))
		}
	case ':':
		if l.peekChar() == ':' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.CONS, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.ACCESS, string(l.ch))
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '<' {
				ch2 := l.ch
				l.readChar()
				tok = l.newToken(token.COMPOSE, string(ch)+string(ch2)+string(l.ch))
			} else {
				tok = l.newToken(token.LESSEQ, string(ch)+string(l.ch))
			}
		} else {
			tok = l.newToken(token.LESS, string(l.ch))

		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '>' {
				ch2 := l.ch
				l.readChar()
				tok = l.newToken(token.PIPE, string(ch)+string(ch2)+string(l.ch))
			} else {
				tok = l.newToken(token.GREATEREQ, string(ch)+string(l.ch))
			}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.SEQUENCE, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.GREATER, string(l.ch))

		}
	case ';':
		tok = l.newToken(token.SEMI, string(l.ch))
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	// Now the next token must either be and identifier
	// Or an invalid token
	default:
		if isIdentifier(l.ch) {
			lit := l.readIdentifier()
			tok = l.newToken(token.LookupIdent(lit), lit)
			return tok
		} else if isDigit(l.ch) {
			kind, value := l.readNumber()
			tok = l.newToken(kind, value)
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.ch))
		}
	}

	// Advance by a character
	l.readChar()
	return tok
}