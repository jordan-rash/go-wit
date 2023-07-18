package parser

import (
	"strings"
	"testing"
	"text/template"

	"github.com/jordan-rash/go-wit/lexer"
	"github.com/jordan-rash/go-wit/token"

	"github.com/stretchr/testify/assert"
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
	tests := []struct {
		Input              string
		expectedType       token.TokenType
		expectedNestedType token.TokenType
	}{
		{"type derp = string", token.KEYWORD_TYPE, token.KEYWORD_STRING},
		{"type derp = list<string>", token.KEYWORD_TYPE, token.KEYWORD_STRING},
		{"type derp = list<list<string>>", token.KEYWORD_TYPE, token.KEYWORD_STRING},
		// {"type derp = list<tuple<string,string>>", token.KEYWORD_TYPE},
		// {"use derp.{bar}", token.KEYWORD_USE},
	}

	tmpl, err := template.New("test").Parse("interface foo { {{ .Input }} }")
	assert.NoError(t, err)

	for _, tt := range tests {

		sb := strings.Builder{}
		err = tmpl.Execute(&sb, tt)
		assert.NoError(t, err)

		p := New(lexer.NewLexer(sb.String()))

		tree := p.Parse()
		assert.NoError(t, p.Errors())

		assert.NotNil(t, tree)
		assert.Len(t, tree.Shapes, 1)

		// for _, tT := range tree.Shapes {
		// 	tShape, ok := tT.(*ast.InterfaceShape)
		// 	if assert.True(t, ok) {
		// 		if assert.Len(t, tShape.Children, 1) {
		// 			for _, tC := range tShape.Children {
		// 				switch tC := tC.(type) {
		// 				case *ast.TypeStatement:
		// 					t.Log("+++", tC.Name, tC.Value)
		// 					assert.Equal(t, strings.ToLower(string(tt.expectedNestedType)), tC.Value.TokenLiteral(), i)
		// 				default:
		// 					t.Error("invalid interface child type")
		// 				}
		// 			}
		// 		}
		// 	}
		// }
	}
}

func TestTypeShape(t *testing.T) {
	tests := []struct {
		Input             string
		expectedType      token.TokenType
		expectedValueType token.TokenType
	}{
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

	for _, tt := range tests {
		p := New(lexer.NewLexer(tt.Input))

		for p.peekToken.Type != token.END_OF_FILE {
			tempType := p.parseTypeStatement()
			assert.Equal(t, tt.expectedType, tempType.Token.Type)
			assert.Equal(t, strings.ToLower(string(tt.expectedValueType)), tempType.Value.TokenLiteral())
		}
	}
}
