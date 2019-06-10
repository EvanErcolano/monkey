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

	// Section for operators

	// ASSIGN represents an assignment operator
	ASSIGN = "="
	// PLUS represents a PLUS operator
	PLUS = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
