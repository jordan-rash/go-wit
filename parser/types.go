package parser

import (
	"errors"
	"fmt"

	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

func (p *Parser) parseTypeShape() ast.Expression {
	ts := new(ast.TypeShape)
	ts.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected IDENT, got %s", p.peekToken.Type))
		return nil
	}

	ts.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_EQUAL) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected EQUAL, got %s", p.peekToken.Type))
		return nil
	}

	ts.Value = p.parseTy()

	return ts
}

func (p *Parser) parseTy() *ast.Ty {
	i := new(ast.Ty)
	i.Token = p.peekToken

	switch p.peekToken.Type {
	case token.KEYWORD_STRING, token.KEYWORD_BOOL, token.KEYWORD_CHAR,
		token.KEYWORD_FLOAT32, token.KEYWORD_FLOAT64,
		token.KEYWORD_S8, token.KEYWORD_S16, token.KEYWORD_S32, token.KEYWORD_S64,
		token.KEYWORD_U8, token.KEYWORD_U16, token.KEYWORD_U32, token.KEYWORD_U64,
		token.IDENTIFIER:

		p.nextToken()

		i.Value = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case token.KEYWORD_LIST:
		if !p.expectNextToken(token.KEYWORD_LIST) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_LIST, got %s", p.peekToken.Type))
			return nil
		}

		c := &ast.TypeShape{Token: p.curToken}
		c.Value = p.parseListShape()
		i.Value = c

	case token.KEYWORD_OPTION:
		if !p.expectNextToken(token.KEYWORD_OPTION) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_OPTION, got %s", p.peekToken.Type))
			return nil
		}

		c := &ast.TypeShape{Token: p.curToken}
		c.Value = p.parseOptionShape()
		i.Value = c

	case token.KEYWORD_RESULT:
		if !p.expectNextToken(token.KEYWORD_RESULT) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected RESULT_LIST, got %s", p.peekToken.Type))
			return nil
		}

		c := &ast.TypeShape{Token: p.curToken}
		c.Value = p.parseResultShape()
		i.Value = c

	case token.KEYWORD_TUPLE:
		if !p.expectNextToken(token.KEYWORD_TUPLE) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_TUPLE, got %s", p.peekToken.Type))
			return nil
		}

		c := &ast.TypeShape{Token: p.curToken}
		c.Value = p.parseTupleShape()
		i.Value = c

	default:
		p.errors = errors.Join(p.errors, fmt.Errorf("unexpected token: %s", p.peekToken.Type))
		p.nextToken()
		i.Value = &ast.Identifier{Token: token.Token{Type: token.ILLEGAL, Literal: ""}, Value: p.curToken.Literal}
	}

	return i
}
