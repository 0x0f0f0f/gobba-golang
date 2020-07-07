package token

type TokenType string

type Token struct {
	Type     TokenType
	Line     int
	Column   int
	Position int
	Literal  string
}

// func (t Token) String() string {
// 	return t.Literal
// }

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENT   = "identifier"
	UNIT    = "()"
	INT     = "integer"
	FLOAT   = "float"
	COMPLEX = "complex number" // Imaginary part
	STRING  = "string"

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
	ANNOT     = ":"
	LAND      = "&&"
	OR        = "||"
	ACCESS    = "."

	// Delimiters
	COMMA       = ","
	SEMI        = ";"
	LPAREN      = "("
	RPAREN      = ")"
	LBRACKET    = "{"
	RBRACKET    = "}"
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
	// Keyword types
	TBOOL    = "bool"
	TINT     = "int"
	TFLOAT   = "float"
	TCOMPLEX = "complex"
	TNUMBER  = "number"
	TRUNE    = "rune"
	TSTRING  = "string"
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
	// Keyword types
	"bool":    TBOOL,
	"int":     TINT,
	"float":   TFLOAT,
	"complex": TCOMPLEX,
	"number":  TNUMBER,
	"rune":    TRUNE,
	"string":  TSTRING,
}

// Check the keywords table to see whether the given
// identifier is a reserved keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
