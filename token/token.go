package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	UNKNOWN = "UNKNOWN"
	EOF     = "EOF"

	INT = "INT"

	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"

	LPAREN    = "("
	RPAREN    = ")"
	SEMICOLON = ";"
)
