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
		switch p.peekToken.Type {
		case token.KEYWORD_INTERFACE:
			if !p.expectNextToken(token.KEYWORD_INTERFACE) {
				return nil
			}
			tree.Interfaces = append(tree.Interfaces, p.parseInterfaceShape())
		case token.KEYWORD_WORLD:
			if !p.expectNextToken(token.KEYWORD_WORLD) {
				return nil
			}
			tree.World = p.parseWorldShape()
		case token.KEYWORD_USE:
			if !p.expectNextToken(token.KEYWORD_USE) {
				return nil
			}
			tree.Uses = append(tree.Uses, p.parseTopUseShape())
		case token.KEYWORD_PACKAGE:
			if !p.expectNextToken(token.KEYWORD_PACKAGE) {
				return nil
			}
			tree.Package = p.parsePackageShape()
		default:
			p.errors = errors.Join(p.errors, errors.New("invalid token: "+p.peekToken.Literal+" ["+string(p.peekToken.Type)+"]"))
			p.nextToken()
		}
	}

	return tree
}

func (p *Parser) expectNextToken(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}

	// p.errors = errors.Join(p.errors, WRONG_NEXT_TYPE_ERROR(string(p.peekToken.Type), string(t)))
	return false
}
