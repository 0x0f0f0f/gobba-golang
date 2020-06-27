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
	IDENT  = "IDENT"
	UNIT   = "UNIT"
	INT    = "INT"
	FLOAT  = "FLOAT"
	IMAG   = "IMAG" // Imaginary part
	STRING = "STRING"

	// Operators
	PLUS      = "+"
	MINUS     = "-"
	TIMES     = "*"
	TOPOW     = "^"
	DIVIDE    = "/"
	EQUALS    = "="
	DIFFERS   = "!="
	LESS      = "<"
	GREATER   = ">"
	LESSEQ    = "<="
	GREATEREQ = ">="
	LARROW    = "->"
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
	ACCESS    = ":"

	// Delimiters
	SEMI        = ";"
	LPAREN      = "("
	RPAREN      = ")"
	LINECOMMENT = "//"
	LCOMMENT    = "/*"
	RCOMMENT    = "*/"

	// Keywords
	LAMBDA = "LAMBDA"
	LET    = "LET"
	IN     = "IN"
	AND    = "AND"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	NOT    = "NOT"
	IF     = "IF"
	THEN   = "THEN"
	ELSE   = "ELSE"
)

// Table of internal keywords
var keywords = map[string]TokenType{
	"lambda": LAMBDA,
	"fun":    LAMBDA,
	"let":    LET,
	"in":     IN,
	"and":    AND,
	"not":    NOT,
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
