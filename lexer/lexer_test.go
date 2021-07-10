package lexer

import (
	"testing"

	"github.com/GenericEntity/interpreter-go/monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10;
	
	let add = fn(x, y) {
		x + y;
	};
	
	let result = add(five, ten);

	if ( !((ten - five)/2 != 2*1)) {
		let x = 2 < 3;
		return true == true;
	} else {
		let x = 3 > 2;
		return !!false;
	}

	"foobar"
	"foo bar"
	"\'\"\\\a\b\f\n\r\t\v"
	"Hello\t\"WORLD\"\n"
	[2,3,4]
	arr[true, 3]
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.BANG, "!"},
		{token.LPAREN, "("},
		{token.LPAREN, "("},
		{token.IDENT, "ten"},
		{token.MINUS, "-"},
		{token.IDENT, "five"},
		{token.RPAREN, ")"},
		{token.SLASH, "/"},
		{token.INT, "2"},
		{token.NOT_EQ, "!="},
		{token.INT, "2"},
		{token.ASTERISK, "*"},
		{token.INT, "1"},
		{token.RPAREN, ")"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "2"},
		{token.LT, "<"},
		{token.INT, "3"},
		{token.SEMICOLON, ";"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.EQ, "=="},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "3"},
		{token.GT, ">"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.RETURN, "return"},
		{token.BANG, "!"},
		{token.BANG, "!"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.STRING, "'\"\\\a\b\f\n\r\t\v"},
		{token.STRING, "Hello\t\"WORLD\"\n"},

		{token.LBRACKET, "["},
		{token.INT, "2"},
		{token.COMMA, ","},
		{token.INT, "3"},
		{token.COMMA, ","},
		{token.INT, "4"},
		{token.RBRACKET, "]"},
		{token.IDENT, "arr"},
		{token.LBRACKET, "["},
		{token.TRUE, "true"},
		{token.COMMA, ","},
		{token.INT, "3"},
		{token.RBRACKET, "]"},

		{token.EOF, ""},
	}

	lex := New(input)

	for i, tt := range tests {
		tok := lex.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
