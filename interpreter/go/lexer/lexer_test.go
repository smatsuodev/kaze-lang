package lexer

import (
	"kaze/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `1 + 2 - 3 * 4 / 5 * (6 + 7) - 8;
1234567890;
var hoge = 1;
fun fuga(x, y) {
    if (true) {
		return x;
	} else if (false) {
		return x + 1;
	}
}
1 == 1;
1 != 1;
1 < 1;
1 > 1;
1 <= 2 >= 1;
!true;
while true {
	if true {
		break;
	}
	if false {
		continue;
	}
}
"hoge"[0];
#{"foo": "bar"};
[1,2,3][0];
true && true || true;
// comment
1; // comment
null;
`
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
		{token.VAR, "var"},
		{token.IDENT, "hoge"},
		{token.ASSIGN, "="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.FUN, "fun"},
		{token.IDENT, "fuga"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.TRUE, "true"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.FALSE, "false"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
		{token.INT, "1"},
		{token.EQ, "=="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.INT, "1"},
		{token.NOT_EQ, "!="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.INT, "1"},
		{token.LT, "<"},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.INT, "1"},
		{token.GT, ">"},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.INT, "1"},
		{token.LE, "<="},
		{token.INT, "2"},
		{token.GE, ">="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.WHILE, "while"},
		{token.TRUE, "true"},
		{token.LBRACE, "{"},
		{token.IF, "if"},
		{token.TRUE, "true"},
		{token.LBRACE, "{"},
		{token.BREAK, "break"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.IF, "if"},
		{token.FALSE, "false"},
		{token.LBRACE, "{"},
		{token.CONTINUE, "continue"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
		{token.STRING, "hoge"},
		{token.LBRACKET, "["},
		{token.INT, "0"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.HASH, "#"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.COMMA, ","},
		{token.INT, "3"},
		{token.RBRACKET, "]"},
		{token.LBRACKET, "["},
		{token.INT, "0"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.TRUE, "true"},
		{token.AND, "&&"},
		{token.TRUE, "true"},
		{token.OR, "||"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.NULL, "null"},
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
