package parser

import (
	"errors"
	"fmt"

	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/lexer"
	"github.com/jordan-rash/go-wit/token"
)

var (
	WRONG_NEXT_TYPE_ERROR func(got, expected string) error = func(got, expected string) error {
		return fmt.Errorf("unexpected next type.  Got: %s\tExpected: %s", got, expected)
	}
)

type Parser struct {
	lexer *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors error
}

// TODO: change input to string and create lexer
func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}

	if p.peekToken.Literal != token.END_OF_FILE {
		p.nextToken()
	}

	if p.peekToken.Literal != token.END_OF_FILE {
		p.nextToken()
	}

	return p
}

func (p Parser) Errors() error {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) Parse() *ast.AST {
	tree := new(ast.AST)

	for p.peekToken.Type != token.END_OF_FILE {
		stmt := p.parseShape()
		if stmt != nil {
			tree.Shapes = append(tree.Shapes, stmt)
		}
	}

	return tree
}

func (p *Parser) parseShape() ast.Shape {

	switch p.curToken.Type {
	case token.KEYWORD_INTERFACE:
		return p.parseInterfaceShape()
	case token.KEYWORD_WORLD:
		return p.parseWorldShape()
	case token.KEYWORD_USE:
		return p.parseUseShape()
	case token.KEYWORD_PACKAGE:
		return p.parsePackageShape()
	}

	return nil
}

func (p *Parser) expectNextToken(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}

	p.errors = errors.Join(p.errors, WRONG_NEXT_TYPE_ERROR(string(p.peekToken.Type), string(t)))
	return false
}

func (p *Parser) makeExpression(tok token.Token) ast.Expression {
	switch tok.Type {
	case token.KEYWORD_STRING, token.KEYWORD_BOOL, token.KEYWORD_CHAR,
		token.KEYWORD_FLOAT32, token.KEYWORD_FLOAT64,
		token.KEYWORD_S8, token.KEYWORD_S16, token.KEYWORD_S32, token.KEYWORD_S64,
		token.KEYWORD_U8, token.KEYWORD_U16, token.KEYWORD_U32, token.KEYWORD_U64,
		token.IDENTIFIER:

		p.nextToken()
		return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	case token.KEYWORD_LIST:
		if !p.expectNextToken(token.KEYWORD_LIST) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_LIST, got %s", p.peekToken.Type))
			return nil
		}
		c := &ast.Child{Token: p.curToken}
		c.Value = p.parseListShape()
		return c

	case token.KEYWORD_OPTION:
		if !p.expectNextToken(token.KEYWORD_OPTION) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_OPTION, got %s", p.peekToken.Type))
			return nil
		}

		c := &ast.Child{Token: p.curToken}
		c.Value = p.parseOptionShape()
		return c

	case token.KEYWORD_RESULT:
		if !p.expectNextToken(token.KEYWORD_RESULT) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected RESULT_LIST, got %s", p.peekToken.Type))
			return nil
		}

		c := &ast.Child{Token: p.curToken}
		c.Value = p.parseResultShape()
		return c

	case token.KEYWORD_TUPLE:
		if !p.expectNextToken(token.KEYWORD_TUPLE) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_TUPLE, got %s", p.peekToken.Type))
			return nil
		}

		c := &ast.Child{Token: p.curToken}
		c.Value = p.parseTupleShape()
		return c

	default:
		p.errors = errors.Join(p.errors, fmt.Errorf("unexpected token: %s", p.peekToken.Type))
		p.nextToken()
		return &ast.Identifier{Token: token.Token{Type: token.ILLEGAL, Literal: ""}, Value: p.curToken.Literal}
	}
}
