package token

type TokenType string // not great performance to use string, but is simple and versatile

type Token struct {
	Type    TokenType
	Literal string
}

// Possible token types
const (
	ILLEGAL = "ILLEGAL" // Unrecognized token
	EOF     = "EOF"     // We can stop parsing

	// Identifiers + literals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN = "="
	PLUS   = "+"

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

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

// LookupIdent checks if ident is a keyword, and returns the appropriate TokenType.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
