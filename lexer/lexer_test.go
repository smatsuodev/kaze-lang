package lexer

import (
	"kaze/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `1 + 2 - 3 * 4 / 5 * (6 + 7) - 8;
1234567890;`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "1"},
		{token.PLUS, "+"},
		{token.INT, "2"},
		{token.MINUS, "-"},
		{token.INT, "3"},
		{token.ASTERISK, "*"},
		{token.INT, "4"},
		{token.SLASH, "/"},
		{token.INT, "5"},
		{token.ASTERISK, "*"},
		{token.LPAREN, "("},
		{token.INT, "6"},
		{token.PLUS, "+"},
		{token.INT, "7"},
		{token.RPAREN, ")"},
		{token.MINUS, "-"},
		{token.INT, "8"},
		{token.SEMICOLON, ";"},
		{token.INT, "1234567890"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
