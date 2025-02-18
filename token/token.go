package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	UNKNOWN = "UNKNOWN"
	EOF     = "EOF"

	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	BANG     = "!"

	EQ     = "=="
	NOT_EQ = "!="

	LT = "<"
	GT = ">"

	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
	COLON     = ":"
	SEMICOLON = ";"
	COMMA     = ","
	HASH      = "#"

	VAR      = "VAR"
	FUN      = "FUN"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	WHILE    = "WHILE"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
)

var keywords = map[string]TokenType{
	"var":      VAR,
	"fun":      FUN,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"true":     TRUE,
	"false":    FALSE,
	"while":    WHILE,
	"break":    BREAK,
	"continue": CONTINUE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
