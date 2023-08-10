package lexer_test

import (
	"testing"

	"github.com/jordan-rash/go-wit/lexer"
	"github.com/jordan-rash/go-wit/token"
	"github.com/stretchr/testify/assert"
)

func TestOPNextToken(t *testing.T) {
	input := "@/<>{}():,=%-.+;*_ "

	l := lexer.NewLexer(input)

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.OP_AT, "@"},
		{token.OP_SLASH, "/"},
		{token.OP_BRACKET_ANGLE_LEFT, "<"},
		{token.OP_BRACKET_ANGLE_RIGHT, ">"},
		{token.OP_BRACKET_CURLY_LEFT, "{"},
		{token.OP_BRACKET_CURLY_RIGHT, "}"},
		{token.OP_BRACKET_PAREN_LEFT, "("},
		{token.OP_BRACKET_PAREN_RIGHT, ")"},
		{token.OP_COLON, ":"},
		{token.OP_COMMA, ","},
		{token.OP_EQUAL, "="},
		{token.OP_EXPLICIT_ID, "%"},
		{token.OP_MINUS, "-"},
		{token.OP_PERIOD, "."},
		{token.OP_PLUS, "+"},
		{token.OP_SEMICOLON, ";"},
		{token.OP_STAR, "*"},
		{token.OP_UNDERSCORE, "_"},
		{token.END_OF_FILE, "EOF"},
	}

	for i, tt := range tests {
		nTok := l.NextToken()

		assert.Equal(t, tt.expectedType, nTok.Type, i)
		assert.Equal(t, tt.expectedLiteral, nTok.Literal, i)
	}
}

func TestSimpleWorld(t *testing.T) {
	input := `package example:host

interface derp {
  type derp = result<_, errno>

  foo: func() -> u16
  bar: func() -> float32
}

world host {
  import print: func(msg: string)

  export run: func()
}
`

	l := lexer.NewLexer(input)

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.KEYWORD_PACKAGE, "package"},
		{token.IDENTIFIER, "example"},
		{token.OP_COLON, ":"},
		{token.IDENTIFIER, "host"},
		{token.KEYWORD_INTERFACE, "interface"},
		{token.IDENTIFIER, "derp"},
		{token.OP_BRACKET_CURLY_LEFT, "{"},
		{token.KEYWORD_TYPE, "type"},
		{token.IDENTIFIER, "derp"},
		{token.OP_EQUAL, "="},
		{token.KEYWORD_RESULT, "result"},
		{token.OP_BRACKET_ANGLE_LEFT, "<"},
		{token.OP_UNDERSCORE, "_"},
		{token.OP_COMMA, ","},
		{token.IDENTIFIER, "errno"},
		{token.OP_BRACKET_ANGLE_RIGHT, ">"},
		{token.IDENTIFIER, "foo"},
		{token.OP_COLON, ":"},
		{token.KEYWORD_FUNC, "func"},
		{token.OP_BRACKET_PAREN_LEFT, "("},
		{token.OP_BRACKET_PAREN_RIGHT, ")"},
		{token.OP_ARROW, "->"},
		{token.KEYWORD_U16, "u16"},
		{token.IDENTIFIER, "bar"},
		{token.OP_COLON, ":"},
		{token.KEYWORD_FUNC, "func"},
		{token.OP_BRACKET_PAREN_LEFT, "("},
		{token.OP_BRACKET_PAREN_RIGHT, ")"},
		{token.OP_ARROW, "->"},
		{token.KEYWORD_FLOAT32, "float32"},
		{token.OP_BRACKET_CURLY_RIGHT, "}"},
		{token.KEYWORD_WORLD, "world"},
		{token.IDENTIFIER, "host"},
		{token.OP_BRACKET_CURLY_LEFT, "{"},
		{token.KEYWORD_IMPORT, "import"},
		{token.IDENTIFIER, "print"},
		{token.OP_COLON, ":"},
		{token.KEYWORD_FUNC, "func"},
		{token.OP_BRACKET_PAREN_LEFT, "("},
		{token.IDENTIFIER, "msg"},
		{token.OP_COLON, ":"},
		{token.KEYWORD_STRING, "string"},
		{token.OP_BRACKET_PAREN_RIGHT, ")"},
		{token.KEYWORD_EXPORT, "export"},
		{token.IDENTIFIER, "run"},
		{token.OP_COLON, ":"},
		{token.KEYWORD_FUNC, "func"},
		{token.OP_BRACKET_PAREN_LEFT, "("},
		{token.OP_BRACKET_PAREN_RIGHT, ")"},
		{token.OP_BRACKET_CURLY_RIGHT, "}"},
		{token.END_OF_FILE, "EOF"},
	}

	for i, tt := range tests {
		nTok := l.NextToken()

		assert.Equal(t, tt.expectedType, nTok.Type, i)
		assert.Equal(t, tt.expectedLiteral, nTok.Literal, i)

	}
}

func TestKebabIdents(t *testing.T) {

	tests := []struct {
		expectedType    string
		expectedLiteral string
	}{
		{token.INT, "7"},
		{token.IDENTIFIER, "jordan-rash"},
		{token.IDENTIFIER, "a-1-a-2"},
		{token.IDENTIFIER, "a1-b2"},
		{token.IDENTIFIER, "1-2"},
	}

	for i, tt := range tests {
		l := lexer.NewLexer(tt.expectedLiteral)
		nTok := l.NextToken()

		assert.Equal(t, tt.expectedType, string(nTok.Type), i)
		assert.Equal(t, tt.expectedLiteral, nTok.Literal, i)

	}
}
