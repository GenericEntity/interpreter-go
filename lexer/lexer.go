package lexer

import "github.com/GenericEntity/interpreter-go/monkey/token"

type Lexer struct {
	input        string
	position     int  // current position in input (after ch)
	readPosition int  // current reading position in input (after ch)
	ch           byte // current character under examination
}

func New(input string) *Lexer {
	lex := &Lexer{
		input: input,
	}

	// initialize lexer to prep first char
	lex.readChar()

	return lex
}

func (lex *Lexer) readChar() {
	// current implementation only supports ASCII chars
	// to extend to unicode and UTF-8
	//	1. change from byte to rune
	//	2. change way of reading characters (to support multi-byte runes)

	if lex.readPosition >= len(lex.input) {
		lex.ch = 0 // NUL byte to indicate end of file
	} else {
		lex.ch = lex.input[lex.readPosition]
	}
	lex.position = lex.readPosition
	lex.readPosition++
}

func (lex *Lexer) NextToken() token.Token {
	var tok token.Token

	switch lex.ch {
	case '=':
		tok = newToken(token.ASSIGN, lex.ch)
	case '+':
		tok = newToken(token.PLUS, lex.ch)
	case ';':
		tok = newToken(token.SEMICOLON, lex.ch)
	case '(':
		tok = newToken(token.LPAREN, lex.ch)
	case ')':
		tok = newToken(token.RPAREN, lex.ch)
	case '{':
		tok = newToken(token.LBRACE, lex.ch)
	case '}':
		tok = newToken(token.RBRACE, lex.ch)
	case ',':
		tok = newToken(token.COMMA, lex.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	lex.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
