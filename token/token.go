package token

const (
	// ILLEGAL describes a Token or character we don't about.
	ILLEGAL = "ILLEGAL"
	// EOF stands for end of file
	EOF = "EOF"

	// Section for Identifiers + Literals

	// IDENT are identifiers ex. x, y, foo, user (user defined)
	IDENT = "IDENT"
	// INT represents integers
	INT = "INT"

	STRING = "STRING"

	// OPERATORS
	ASSIGN   = "="
	BANG     = "!"
	PLUS     = "+"
	MINUS    = "-"
	SLASH    = "/"
	ASTERISK = "*"

	GT     = ">"
	LT     = "<"
	EQ     = "=="
	NOT_EQ = "!="

	// DELIMITERS
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// KEYWORDS
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
