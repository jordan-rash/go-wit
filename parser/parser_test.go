package parser

import (
	"strings"
	"testing"
	"text/template"

	"github.com/jordan-rash/go-wit/lexer"
	"github.com/jordan-rash/go-wit/token"

	"github.com/stretchr/testify/assert"
)

type tests []struct {
	Input             string
	expectedType      token.TokenType
	expectedValueType token.TokenType
}

var (
	typeTests = tests{
		{"type derp = string", token.KEYWORD_TYPE, token.KEYWORD_STRING},
		{"type derp = bool", token.KEYWORD_TYPE, token.KEYWORD_BOOL},
		{"type derp = char", token.KEYWORD_TYPE, token.KEYWORD_CHAR},
		{"type derp = float32", token.KEYWORD_TYPE, token.KEYWORD_FLOAT32},
		{"type derp = float64", token.KEYWORD_TYPE, token.KEYWORD_FLOAT64},
		{"type derp = s8", token.KEYWORD_TYPE, token.KEYWORD_S8},
		{"type derp = s16", token.KEYWORD_TYPE, token.KEYWORD_S16},
		{"type derp = s32", token.KEYWORD_TYPE, token.KEYWORD_S32},
		{"type derp = s64", token.KEYWORD_TYPE, token.KEYWORD_S64},
		{"type derp = u8", token.KEYWORD_TYPE, token.KEYWORD_U8},
		{"type derp = u16", token.KEYWORD_TYPE, token.KEYWORD_U16},
		{"type derp = u32", token.KEYWORD_TYPE, token.KEYWORD_U32},
		{"type derp = u64", token.KEYWORD_TYPE, token.KEYWORD_U64},
	}

	listTests = tests{
		{"list<string>", token.KEYWORD_TYPE, token.KEYWORD_STRING},
		{"list<bool>", token.KEYWORD_TYPE, token.KEYWORD_BOOL},
		{"list<char>", token.KEYWORD_TYPE, token.KEYWORD_CHAR},
		{"list<float32>", token.KEYWORD_TYPE, token.KEYWORD_FLOAT32},
		{"list<float64>", token.KEYWORD_TYPE, token.KEYWORD_FLOAT64},
		{"list<s8>", token.KEYWORD_TYPE, token.KEYWORD_S8},
		{"list<s16>", token.KEYWORD_TYPE, token.KEYWORD_S16},
		{"list<s32>", token.KEYWORD_TYPE, token.KEYWORD_S32},
		{"list<s64>", token.KEYWORD_TYPE, token.KEYWORD_S64},
		{"list<u8>", token.KEYWORD_TYPE, token.KEYWORD_U8},
		{"list<u16>", token.KEYWORD_TYPE, token.KEYWORD_U16},
		{"list<u32>", token.KEYWORD_TYPE, token.KEYWORD_U32},
	}

	optionTests = tests{
		{"option<string>", token.KEYWORD_TYPE, token.KEYWORD_STRING},
		{"option<bool>", token.KEYWORD_TYPE, token.KEYWORD_BOOL},
		{"option<char>", token.KEYWORD_TYPE, token.KEYWORD_CHAR},
		{"option<float32>", token.KEYWORD_TYPE, token.KEYWORD_FLOAT32},
		{"option<float64>", token.KEYWORD_TYPE, token.KEYWORD_FLOAT64},
		{"option<s8>", token.KEYWORD_TYPE, token.KEYWORD_S8},
		{"option<s16>", token.KEYWORD_TYPE, token.KEYWORD_S16},
		{"option<s32>", token.KEYWORD_TYPE, token.KEYWORD_S32},
		{"option<s64>", token.KEYWORD_TYPE, token.KEYWORD_S64},
		{"option<u8>", token.KEYWORD_TYPE, token.KEYWORD_U8},
		{"option<u16>", token.KEYWORD_TYPE, token.KEYWORD_U16},
		{"option<u32>", token.KEYWORD_TYPE, token.KEYWORD_U32},
		{"option<u64>", token.KEYWORD_TYPE, token.KEYWORD_U64},
	}
)

func TestRootShapes(t *testing.T) {
	tests := []struct {
		input        string
		expectedType token.TokenType
	}{
		{"interface derp {}", token.KEYWORD_INTERFACE},
		{"world derp {}", token.KEYWORD_WORLD},
		{"use derp.{foo}", token.KEYWORD_USE},
	}

	for _, tt := range tests {
		p := New(lexer.NewLexer(tt.input))

		tree := p.Parse()
		assert.NoError(t, p.Errors())

		assert.NotNil(t, tree)
		assert.Len(t, tree.Shapes, 1)

		for _, tT := range tree.Shapes {
			assert.NotNil(t, tT)
			assert.Equal(t, strings.ToLower(string(tt.expectedType)), tT.TokenLiteral())
		}
	}
}

func TestNestedInterfaceShapes(t *testing.T) {
	tmpl, err := template.New("test").Parse("interface foo { {{ .Input }} }")
	assert.NoError(t, err)

	for _, tt := range typeTests {
		sb := strings.Builder{}
		err = tmpl.Execute(&sb, tt)
		assert.NoError(t, err)

		p := New(lexer.NewLexer(sb.String()))

		tree := p.Parse()
		assert.NotNil(t, tree)
		assert.NoError(t, p.Errors())

		assert.NotNil(t, tree)
		assert.Len(t, tree.Shapes, 1)
	}
}

func TestTypeShape(t *testing.T) {
	for _, tt := range typeTests {
		p := New(lexer.NewLexer(tt.Input))

		for p.peekToken.Type != token.END_OF_FILE {
			tempType := p.parseTypeStatement()
			assert.Equal(t, tt.expectedType, tempType.Token.Type)
			assert.Equal(t, strings.ToLower(string(tt.expectedValueType)), tempType.Value.TokenLiteral())
		}
	}
}

func TestTypeListShape(t *testing.T) {
	for _, tt := range listTests {
		p := New(lexer.NewLexer(tt.Input))

		for p.peekToken.Type != token.END_OF_FILE {
			tempType := p.parseListShape()
			assert.NoError(t, p.Errors())

			assert.Equal(t, token.KEYWORD_LIST, string(tempType.Name.Token.Type))
			assert.Equal(t, strings.ToLower(string(tt.expectedValueType)), tempType.Value.TokenLiteral())
		}
	}
}

func TestTypeOptionShape(t *testing.T) {
	for _, tt := range optionTests {
		p := New(lexer.NewLexer(tt.Input))

		for p.peekToken.Type != token.END_OF_FILE {
			tempType := p.parseOptionShape()
			assert.NoError(t, p.Errors())

			assert.Equal(t, token.KEYWORD_OPTION, string(tempType.Name.Token.Type))
			assert.Equal(t, strings.ToLower(string(tt.expectedValueType)), tempType.Value.TokenLiteral())
		}
	}
}
