package token

type TokenType string

type Token struct {
	Type    TokenType
	Line    int
	Column  int
	Literal string
}

// func (t Token) String() string {
// 	return t.Literal
// }

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENT  = "identifier"
	UNIT   = "()"
	INT    = "integer"
	FLOAT  = "float"
	IMAG   = "imaginary number" // Imaginary part
	STRING = "string"

	// Operators
	PLUS      = "+"
	MINUS     = "-"
	TIMES     = "*"
	TOPOW     = "^"
	DIVIDE    = "/"
	EQUALS    = "="
	DIFFERS   = "!="
	NOT       = "!"
	LESS      = "<"
	GREATER   = ">"
	LESSEQ    = "<="
	GREATEREQ = ">="
	RARROW    = "->"
	SEQUENCE  = ">>"
	COMPOSE   = "<=<"
	PIPE      = ">=>"
	MODULO    = "%"
	CONCAT    = "++"
	AT        = "@"
	DOLLAR    = "$"
	CONS      = "::"
	LAND      = "&&"
	OR        = "||"
	ACCESS    = "."

	// Delimiters
	SEMI        = ";"
	LPAREN      = "("
	RPAREN      = ")"
	LINECOMMENT = "//"
	LCOMMENT    = "/*"
	RCOMMENT    = "*/"

	// Keywords
	LAMBDA = "lambda"
	LET    = "let"
	IN     = "in"
	AND    = "and"
	TRUE   = "true"
	FALSE  = "false"
	IF     = "if"
	THEN   = "then"
	ELSE   = "else"
)

// Table of internal keywords
var keywords = map[string]TokenType{
	"lambda": LAMBDA,
	"fun":    LAMBDA,
	"let":    LET,
	"in":     IN,
	"and":    AND,
	"if":     IF,
	"then":   THEN,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
}

// Check the keywords table to see whether the given
// identifier is a reserved keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
