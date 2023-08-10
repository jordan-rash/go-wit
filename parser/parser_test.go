package parser

import (
	"strings"
	"testing"

	"text/template"

	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/lexer"
	"github.com/jordan-rash/go-wit/token"

	"github.com/stretchr/testify/assert"
)

type tests []struct {
	Input             string
	expectedType      any
	expectedValueType any
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
		{"type derp = foo", token.KEYWORD_TYPE, token.IDENTIFIER},
		// {"type derp = result", token.KEYWORD_TYPE, token.KEYWORD_RESULT},
		// {"type derp = result<string>", token.KEYWORD_TYPE, nil},
		// {"type derp = result<char, errno>", token.KEYWORD_TYPE, token.IDENTIFIER},
		// {"type derp = result<_, errno>", token.KEYWORD_TYPE, token.IDENTIFIER},
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
		{"list<foo>", token.KEYWORD_TYPE, token.IDENTIFIER},
		{"list<foo-bar>", token.KEYWORD_TYPE, token.IDENTIFIER},
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
		{"option<foo>", token.KEYWORD_TYPE, token.IDENTIFIER},
		{"option<foo-bar>", token.KEYWORD_TYPE, token.IDENTIFIER},
	}
	exportTests = []struct {
		Input             string
		expectedName      string
		expectedType      any
		expectedValueType any
	}{
		{"export derp: func() -> string", "derp", token.KEYWORD_EXPORT, nil},
		{"export derp: interface { derp: func() -> string }", "derp", token.KEYWORD_EXPORT, nil},
		{"export foo-bar", "foo-bar", token.KEYWORD_EXPORT, token.IDENTIFIER},
		{"export wasi:http/handler", "wasi:http/handler", token.KEYWORD_EXPORT, nil},
		{"export wasi:http/handler@1.0.0", "wasi:http/handler@1.0.0", token.KEYWORD_EXPORT, nil},
	}
	resultTests = []struct {
		input            string
		expectedOkValue  any
		expectedErrValue any
	}{
		{"result<string>", token.KEYWORD_STRING, nil},
		{"result<char, errno>", token.KEYWORD_CHAR, token.IDENTIFIER},
		{"result<_, u16>", nil, token.KEYWORD_U16},
		{"result<list<string>, errno>", token.KEYWORD_LIST, token.IDENTIFIER},
		{"result<_, option<string>>", nil, token.KEYWORD_OPTION},
	}
	tupleTests = []struct {
		input          string
		expectedsValue []any
	}{
		{"tuple<string>", []any{token.KEYWORD_STRING}},
		{"tuple<char, errno>", []any{token.KEYWORD_CHAR, token.IDENTIFIER}},
		{"tuple<list<string>>", []any{token.KEYWORD_LIST}},
		{"tuple<list<string>, option<string>>", []any{token.KEYWORD_LIST, token.KEYWORD_OPTION}},
	}
	packageTests = []struct {
		input     string
		namespace string
		name      string
		version   string
	}{
		{"package wasi:derp", "wasi", "derp", ""},
		{"package wasi:derp@0.1.0", "wasi", "derp", "0.1.0"},
	}
	funcTests = []struct {
		Input              string
		name               string
		expectedParamList  []string
		expectedResultList []string
	}{
		{"derp: func() -> string", "derp", []string{}, []string{token.KEYWORD_STRING}},
		{"derp: func() -> foo", "derp", []string{}, []string{token.IDENTIFIER}},
	}
)

func TestParsePingPong(t *testing.T) {
	input := `package jordan-rash:pingpong@0.1.0

interface types {
  type pong = string
}

interface pingpong {
  use types.{pong}
  ping: func() -> pong

  log: func(string) -> string
}

world ping-pong {
  export pingpong
  export types
}
`

	iFaceAnswers := []struct {
		name        string
		useLength   int
		tDefsLength int
		funcLength  int
	}{{"types", 0, 1, 0}, {"pingpong", 1, 0, 2}}

	p := New(lexer.NewLexer(input))
	a := p.Parse()

	assert.NoError(t, p.Errors())
	assert.NotNil(t, a)

	pkg, ok := a.Package.(*ast.Package)
	if assert.True(t, ok) {
		assert.Equal(t, "jordan-rash", pkg.Namespace)
		assert.Equal(t, "pingpong", pkg.Name)
		assert.Equal(t, "0.1.0", pkg.SemVer)
	}

	for idx, i := range a.Interfaces {
		ii, ok := i.(*ast.Interface)
		if assert.True(t, ok) {
			assert.Equal(t, iFaceAnswers[idx].name, ii.Name)
			assert.Len(t, ii.Items.UseItems, iFaceAnswers[idx].useLength)
			assert.Len(t, ii.Items.TypedefItems, iFaceAnswers[idx].tDefsLength)
			assert.Len(t, ii.Items.FuncItems, iFaceAnswers[idx].funcLength)
		}
	}

	//TODO add top level use to test
}

func TestRootShapes(t *testing.T) {
	tests := []struct {
		input        string
		expectedType token.TokenType
	}{
		{"interface derp {}", token.KEYWORD_INTERFACE},
		{"world derp {}", token.KEYWORD_WORLD},
		{"use derp.{foo}", token.KEYWORD_USE},
		{"package jordan-rash:pingpong@0.1.0", token.KEYWORD_PACKAGE},
	}

	for _, tt := range tests {
		p := New(lexer.NewLexer(tt.input))

		tree := p.Parse()
		assert.NoError(t, p.Errors())

		assert.NotNil(t, tree)

		switch tt.expectedType {
		case token.KEYWORD_INTERFACE:
			i, ok := tree.Interfaces[0].(*ast.Interface)

			assert.True(t, ok)
			assert.Equal(t, "derp", i.Name)
		case token.KEYWORD_WORLD:
			w, ok := tree.World.(*ast.World)

			assert.True(t, ok)
			assert.Equal(t, "derp", w.Name)
		case token.KEYWORD_USE:
			u, ok := tree.Uses[0].(*ast.Use)

			assert.True(t, ok)
			assert.Equal(t, "derp", u.Identifier.TokenLiteral())
		case token.KEYWORD_PACKAGE:
			p, ok := tree.Package.(*ast.Package)

			assert.True(t, ok)
			assert.Equal(t, "jordan-rash", p.Namespace)
			assert.Equal(t, "pingpong", p.Name)
			assert.Equal(t, "0.1.0", p.SemVer)
		}
	}
}

func TestNestedInterfaceShapes(t *testing.T) {
	tmpl, err := template.New("test").Parse("interface foo { {{ .Input }} }")
	assert.NoError(t, err)

	for i, tt := range typeTests {
		sb := strings.Builder{}
		err = tmpl.Execute(&sb, tt)
		assert.NoError(t, err)

		p := New(lexer.NewLexer(sb.String()))

		tree := p.Parse()
		assert.NotNil(t, tree)
		assert.NoError(t, p.Errors(), i)

		assert.NotNil(t, tree)
		assert.Len(t, tree.Interfaces, 1)
	}

	for _, tt := range funcTests {
		sb := strings.Builder{}
		err = tmpl.Execute(&sb, tt)
		assert.NoError(t, err)

		p := New(lexer.NewLexer(sb.String()))

		tree := p.Parse()
		assert.NotNil(t, tree)
		assert.NoError(t, p.Errors())

		assert.NotNil(t, tree)
		assert.Len(t, tree.Interfaces, 1)
	}
}

func TestNestedWorldShapes(t *testing.T) {
	tmpl, err := template.New("test").Parse("world foo { {{ .Input }} }")
	assert.NoError(t, err)

	for _, tt := range exportTests {
		sb := strings.Builder{}
		err = tmpl.Execute(&sb, tt)
		t.Log("TESTING -> ", sb.String())
		assert.NoError(t, err)

		p := New(lexer.NewLexer(sb.String()))
		tree := p.Parse()

		assert.NotNil(t, tree)
		assert.NoError(t, p.Errors())

		assert.NotNil(t, tree.World)

		w, ok := tree.World.(*ast.World)
		assert.True(t, ok)
		assert.Equal(t, "foo", w.Name)
		assert.Len(t, w.ExportItems, 1)
	}
}

func TestTypeShape(t *testing.T) {
	for i, tt := range typeTests {
		p := New(lexer.NewLexer(tt.Input))
		t.Log("TESTING ->", tt.Input)
		for p.peekToken.Type != token.END_OF_FILE {
			assert.True(t, p.expectNextToken(token.KEYWORD_TYPE))
			e := p.parseTypeShape()
			assert.NoError(t, p.Errors())

			ts, ok := e.(*ast.TypeShape)
			assert.True(t, ok, i)

			assert.Equal(t, tt.expectedType, string(ts.Token.Type))

			ty, ok := ts.Value.(*ast.Ty)
			assert.True(t, ok, i)
			assert.Equal(t, tt.expectedValueType, string(ty.Token.Type))
		}
	}
}

func TestTypeListShape(t *testing.T) {
	for i, tt := range listTests {
		p := New(lexer.NewLexer(tt.Input))
		t.Log("TESTING ->", tt.Input)
		for p.peekToken.Type != token.END_OF_FILE {
			assert.True(t, p.expectNextToken(token.KEYWORD_LIST))
			tempType := p.parseListShape()
			assert.NoError(t, p.Errors(), i)

			x, ok := tempType.Value.(*ast.Ty)
			assert.True(t, ok)

			assert.Equal(t, token.KEYWORD_LIST, string(tempType.Name.Token.Type))
			assert.Equal(t, tt.expectedValueType, string(x.Token.Type))
		}
	}
}

func TestTypeOptionShape(t *testing.T) {
	for _, tt := range optionTests {
		p := New(lexer.NewLexer(tt.Input))
		t.Log("TESTING ->", tt.Input)

		for p.peekToken.Type != token.END_OF_FILE {
			assert.True(t, p.expectNextToken(token.KEYWORD_OPTION))
			tempType := p.parseOptionShape()
			assert.NoError(t, p.Errors())

			x, ok := tempType.Value.(*ast.Ty)
			assert.True(t, ok)

			assert.Equal(t, token.KEYWORD_OPTION, string(tempType.Name.Token.Type))
			assert.Equal(t, tt.expectedValueType, string(x.Token.Type))
		}
	}
}

func TestTypeTupleShape(t *testing.T) {
	for _, tt := range tupleTests {
		p := New(lexer.NewLexer(tt.input))
		t.Log("TESTING ->", tt.input)
		for p.peekToken.Type != token.END_OF_FILE {
			assert.True(t, p.expectNextToken(token.KEYWORD_TUPLE))
			tempType := p.parseTupleShape()
			assert.NoError(t, p.Errors())

			assert.Equal(t, token.KEYWORD_TUPLE, string(tempType.Name.Token.Type))
			for i, x := range tempType.Value {
				xx, ok := x.(*ast.Ty)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedsValue[i], string(xx.Token.Type))
			}

		}
	}
}

func TestTypeResultShape(t *testing.T) {
	for _, tt := range resultTests {
		p := New(lexer.NewLexer(tt.input))
		t.Log("TESTING ->", tt.input)
		for p.peekToken.Type != token.END_OF_FILE {
			assert.True(t, p.expectNextToken(token.KEYWORD_RESULT))
			tempType := p.parseResultShape()
			assert.NoError(t, p.Errors())
			assert.Equal(t, token.KEYWORD_RESULT, string(tempType.Name.Token.Type))

			if tempType.OkValue != nil {
				ty, ok := tempType.OkValue.(*ast.Ty)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedOkValue, string(ty.Token.Type))
			} else {
				assert.Nil(t, tempType.OkValue)
			}

			if tempType.ErrValue != nil {
				ty, ok := tempType.ErrValue.(*ast.Ty)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedErrValue, string(ty.Token.Type))
			} else {
				assert.Nil(t, tempType.ErrValue)
			}
		}
	}
}

func TestTypePackageShape(t *testing.T) {
	for _, tt := range packageTests {
		p := New(lexer.NewLexer(tt.input))
		t.Log("TESTING ->", tt.input)
		for p.peekToken.Type != token.END_OF_FILE {
			tree := p.Parse()
			assert.NoError(t, p.Errors())

			pkg, ok := tree.Package.(*ast.Package)
			assert.True(t, ok)

			assert.Equal(t, tt.namespace, pkg.Namespace)
			assert.Equal(t, tt.name, pkg.Name)
			assert.Equal(t, tt.version, pkg.SemVer)

			assert.Equal(t, token.KEYWORD_PACKAGE, string(strings.ToUpper(pkg.Identifier.TokenLiteral())))
		}
	}
}

func TestExportShape(t *testing.T) {
	for _, tt := range exportTests {
		p := New(lexer.NewLexer(tt.Input))
		t.Log("TESTING ->", tt.Input)
		for p.peekToken.Type != token.END_OF_FILE {
			assert.True(t, p.expectNextToken(token.KEYWORD_EXPORT))
			es := p.parseExportStatement()
			assert.NoError(t, p.Errors())
			assert.NotNil(t, es)
			assert.Equal(t, tt.expectedType, string(es.Token.Type))
			assert.Equal(t, tt.expectedName, es.Name.Value)
		}
	}
}

func TestFuncShape(t *testing.T) {
	for i, tt := range funcTests {
		p := New(lexer.NewLexer(tt.Input))
		t.Log("TESTING ->", tt.Input)
		for p.peekToken.Type != token.END_OF_FILE {
			tempType := p.parseFuncItem()
			assert.NoError(t, p.Errors())

			ft, ok := tempType.Value.(*ast.FuncType)
			assert.True(t, ok)

			assert.Equal(t, tt.name, tempType.Name.Token.Literal, i)
			assert.Len(t, *ft.ParamList, len(tt.expectedParamList))
			assert.Len(t, *ft.ResultList, len(tt.expectedResultList))

			for i, param := range *ft.ParamList {
				ty, ok := param.(*ast.Ty)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedParamList[i], string(ty.Token.Type))
			}

			for i, res := range *ft.ResultList {
				ty, ok := res.(*ast.Ty)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedResultList[i], string(ty.Token.Type))
			}
		}
	}
}

func TestResourceShape(t *testing.T) {
	resourceTest := struct {
		input         string
		name          string
		expectedValue []ast.Expression
	}{
		input: `resource blob {
    constructor(init: list<u8>)
    write: func(bytes: list<u8>)
    read: func(n: u32) -> list<u8>
    merge: static func(lhs: borrow<blob>, rhs: borrow<blob>) -> blob
  }`,
		name:          "blob",
		expectedValue: []ast.Expression{},
	}

	p := New(lexer.NewLexer(resourceTest.input))
	t.Log("TESTING ->", resourceTest.input)

	for p.peekToken.Type != token.END_OF_FILE {
		assert.True(t, p.expectNextToken(token.KEYWORD_RESOURCE))
		tempType := p.parseResourceShape()
		assert.NoError(t, p.Errors())

		assert.Equal(t, "blob", tempType.Name.Token.Literal)
		assert.Equal(t, "RESOURCE", string(tempType.Token.Type))

		assert.Len(t, tempType.Value, 4)
	}
}

func TestEnumShape(t *testing.T) {
	enumTest := struct {
		input         string
		expectedName  string
		expectedValue []string
	}{
		input: `enum color {
    red,
    green,
    blue,
    yellow,
    other,
   }`,
		expectedName:  "color",
		expectedValue: []string{"red", "green", "blue", "yellow", "other"},
	}

	p := New(lexer.NewLexer(enumTest.input))
	t.Log("TESTING ->", enumTest.input)

	for p.peekToken.Type != token.END_OF_FILE {
		assert.True(t, p.expectNextToken(token.KEYWORD_ENUM))
		tempType := p.parseEnumShape()
		assert.NoError(t, p.Errors())

		assert.Equal(t, enumTest.expectedName, tempType.Name.Token.Literal)
		assert.Equal(t, token.KEYWORD_ENUM, string(tempType.Token.Type))

		assert.Len(t, tempType.Value, len(enumTest.expectedValue))

		for i, v := range tempType.Value {
			ty, ok := v.(*ast.Ty)
			assert.True(t, ok)
			assert.Equal(t, enumTest.expectedValue[i], ty.Value.TokenLiteral())
		}
	}
}

func TestFlagShape(t *testing.T) {
	flagTest := struct {
		input         string
		expectedName  string
		expectedValue []string
	}{
		input: `flags properties {
    lego,
    marvel-superhero,
    supervillan,
   }
   `,
		expectedName:  "properties",
		expectedValue: []string{"lego", "marvel-superhero", "supervillan"},
	}

	p := New(lexer.NewLexer(flagTest.input))
	t.Log("TESTING ->", flagTest.input)

	for p.peekToken.Type != token.END_OF_FILE {
		assert.True(t, p.expectNextToken(token.KEYWORD_FLAGS))
		tempType := p.parseFlagShape()
		assert.NoError(t, p.Errors())

		assert.Equal(t, flagTest.expectedName, tempType.Name.Token.Literal)
		assert.Equal(t, token.KEYWORD_FLAGS, string(tempType.Token.Type))

		assert.Len(t, tempType.Value, len(flagTest.expectedValue))

		for i, v := range tempType.Value {
			ty, ok := v.(*ast.Ty)
			assert.True(t, ok)
			assert.Equal(t, flagTest.expectedValue[i], ty.Value.TokenLiteral())
		}
	}
}
