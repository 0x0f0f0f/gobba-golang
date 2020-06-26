package lexer

import (
	"github.com/0x0f0f0f/gobba-golang/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5 and ten = 10 and add = // i should be ignored
lambda x y -> x + y in add five ten ;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.EQUALS, "="},
		{token.INT, "5"},
		{token.AND, "and"},
		{token.IDENT, "ten"},
		{token.EQUALS, "="},
		{token.INT, "10"},
		{token.AND, "and"},
		{token.IDENT, "add"},
		{token.EQUALS, "="},
		{token.LAMBDA, "lambda"},
		{token.IDENT, "x"},
		{token.IDENT, "y"},
		{token.LARROW, "->"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.IN, "in"},
		{token.IDENT, "add"},
		{token.IDENT, "five"},
		{token.IDENT, "ten"},
		{token.SEMI, ";"},
	}

	// Create a new Lexer
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		t.Logf("%+v\n", tok)

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong, expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
