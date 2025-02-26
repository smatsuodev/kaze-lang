package lexer

import "kaze/token"

type Lexer struct {
	input   string
	pos     int
	nextPos int
	ch      byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.nextPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.nextPos]
	}
	l.pos = l.nextPos
	l.nextPos++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	for l.ch == '/' && l.peekChar() == '/' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
		l.skipWhitespace()
	}

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			tok = token.Token{Type: token.EQ, Literal: "=="}
			l.readChar()
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '!':
		if l.peekChar() == '=' {
			tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
			l.readChar()
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			tok = token.Token{Type: token.AND, Literal: "&&"}
			l.readChar()
		} else {
			tok = newToken(token.UNKNOWN, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			tok = token.Token{Type: token.OR, Literal: "||"}
			l.readChar()
		} else {
			tok = newToken(token.UNKNOWN, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			tok = token.Token{Type: token.LE, Literal: "<="}
			l.readChar()
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			tok = token.Token{Type: token.GE, Literal: ">="}
			l.readChar()
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '#':
		tok = newToken(token.HASH, l.ch)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	case '"':
		literal, ok := l.readString()
		tok.Literal = literal
		if ok {
			tok.Type = token.STRING
		} else {
			tok.Type = token.UNKNOWN
		}
		return tok
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		}
		if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readInteger()
			return tok
		}
		tok = newToken(token.UNKNOWN, l.ch)
	}

	l.readChar()
	return tok
}

func (l *Lexer) peekChar() byte {
	if l.nextPos >= len(l.input) {
		return 0
	}
	return l.input[l.nextPos]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readInteger() string {
	pos := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readString() (string, bool) {
	l.readChar()
	pos := l.pos
	for l.ch != '"' && l.ch != 0 && l.ch != '\n' {
		l.readChar()
	}
	if l.ch != '"' {
		println("string literal not terminated")
		return "", false
	}
	result := l.input[pos:l.pos]
	l.readChar()
	return result, true
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
