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

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}

	p.nextToken()
	p.nextToken()

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
	for p.curToken.Type != token.END_OF_FILE {
		stmt := p.parseShape()
		if stmt != nil {
			tree.Shapes = append(tree.Shapes, stmt)
		}
		p.nextToken()
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
	default:
		// TODO: log error
		return nil
	}
}

func (p *Parser) expectNextType(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}

	p.errors = errors.Join(p.errors, WRONG_NEXT_TYPE_ERROR(string(p.peekToken.Type), string(t)))
	return false
}
