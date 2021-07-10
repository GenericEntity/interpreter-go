package lexer

import (
	"fmt"
	"strings"

	"github.com/GenericEntity/interpreter-go/monkey/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (after ch)
	readPosition int  // current reading position in input (after ch)
	ch           byte // current character under examination
}

// TODO:
//  floating point numbers
//  hex notation, octal notation, binary notation
//  identifiers with digits
//  comments
//  &&, ||

var escapeCharacterMap = map[byte]byte{
	'\'': '\'',
	'"':  '"',
	'\\': '\\',
	'a':  '\a',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'v':  '\v',
}

func New(input string) *Lexer {
	lex := &Lexer{
		input: input,
	}

	// initialize lexer to prep first char
	lex.readChar()

	return lex
}

func (lex *Lexer) NextToken() token.Token {
	var tok token.Token

	// some languages check newlines. if so, then they can't be skipped
	lex.skipWhitespace()

	switch lex.ch {
	case '=':
		if lex.peekChar() == '=' {
			ch := lex.ch
			lex.readChar()
			literal := string(ch) + string(lex.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, lex.ch)
		}
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
	case '!':
		if lex.peekChar() == '=' {
			firstChar := lex.ch
			lex.readChar()
			literal := string(firstChar) + string(lex.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, lex.ch)
		}
	case '-':
		tok = newToken(token.MINUS, lex.ch)
	case '/':
		tok = newToken(token.SLASH, lex.ch)
	case '*':
		tok = newToken(token.ASTERISK, lex.ch)
	case '<':
		tok = newToken(token.LT, lex.ch)
	case '>':
		tok = newToken(token.GT, lex.ch)

	case '"':
		var err error
		tok.Type = token.STRING
		tok.Literal, err = lex.readString()
		if err != nil {
			tok.Type = token.ILLEGAL
		}

	case '[':
		tok = newToken(token.LBRACKET, lex.ch)
	case ']':
		tok = newToken(token.RBRACKET, lex.ch)

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isIdentifierChar(lex.ch) {
			tok.Literal = lex.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(lex.ch) {
			tok.Literal = lex.readInteger()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, lex.ch)
		}
	}

	lex.readChar()
	return tok
}

func (lex *Lexer) peekChar() byte {
	if lex.readPosition >= len(lex.input) {
		return 0
	}
	return lex.input[lex.readPosition]
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

func (lex *Lexer) readIdentifier() string {
	start := lex.position
	for isIdentifierChar(lex.ch) {
		lex.readChar()
	}
	return lex.input[start:lex.position]
}

func isIdentifierChar(ch byte) bool {
	return isLetter(ch) || ch == '_'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (lex *Lexer) skipWhitespace() {
	for isWhitespace(lex.ch) {
		lex.readChar()
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (lex *Lexer) readInteger() string {
	start := lex.position
	for isDigit(lex.ch) {
		lex.readChar()
	}
	return lex.input[start:lex.position]
}

func isEscapeCharacter(ch byte) bool {
	return ch == '\\'
}

func (lex *Lexer) readString() (string, error) {
	var str strings.Builder
	for {
		lex.readChar()
		ch := lex.ch

		// handle unexpected end of file
		if lex.ch == 0 {
			return "", fmt.Errorf("unexpected EOF encountered while reading string")
		}

		// handle end of string
		if lex.ch == '"' {
			break
		}

		// handle escape sequences AFTER testing for EOF or end of string
		if isEscapeCharacter(lex.ch) {
			lex.readChar()
			escChar, ok := escapeCharacterMap[lex.ch]
			if !ok {
				return "", fmt.Errorf("unknown escape sequence: '\\%v'", lex.ch)
			}
			ch = escChar
		}

		str.WriteByte(ch)
	}
	return str.String(), nil
}
