package parser

import (
	"errors"
	"fmt"

	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

func (p *Parser) parseTypeStatement() *ast.TypeStatement {
	stmt := new(ast.TypeStatement)
	stmt.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected IDENT, got %s", p.peekToken.Type))
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_EQUAL) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected EQUAL, got %s", p.peekToken.Type))
		return nil
	}

	switch p.peekToken.Type {
	case token.KEYWORD_STRING, token.KEYWORD_BOOL, token.KEYWORD_CHAR,
		token.KEYWORD_FLOAT32, token.KEYWORD_FLOAT64,
		token.KEYWORD_S8, token.KEYWORD_S16, token.KEYWORD_S32, token.KEYWORD_S64,
		token.KEYWORD_U8, token.KEYWORD_U16, token.KEYWORD_U32, token.KEYWORD_U64:

		p.nextToken()
		stmt.Value = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	case token.KEYWORD_LIST:
		if !p.expectNextToken(token.KEYWORD_LIST) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_LIST, got %s", p.peekToken.Type))
			return nil
		}
		stmt.Value = &ast.Child{Token: p.curToken, Value: p.parseListShape()}
	case token.KEYWORD_OPTION:
		if !p.expectNextToken(token.KEYWORD_OPTION) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_OPTION, got %s", p.peekToken.Type))
			return nil
		}
		stmt.Value = &ast.Child{Token: p.curToken, Value: p.parseOptionShape()}
	case token.KEYWORD_RESULT:
		// p.parseResultShape()
	case token.KEYWORD_TUPLE:
		// p.parseTupleShape()
	default:
		// stmt.Value = &ast.Child{Token: token.ILLEGAL, Value: p.curToken.Literal}
		p.errors = errors.Join(p.errors, fmt.Errorf("unexpected token: %s", p.peekToken.Type))
		return nil
	}

	return stmt
}

func (p *Parser) parseOptionShape() *ast.OptionShape {
	os := new(ast.OptionShape)
	os.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_BRACKET_ANGLE_LEFT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_ANGLE_LEFT, got %s", p.peekToken.Type))
		return nil
	}

	switch p.peekToken.Type {
	case token.KEYWORD_STRING, token.KEYWORD_BOOL, token.KEYWORD_CHAR,
		token.KEYWORD_FLOAT32, token.KEYWORD_FLOAT64,
		token.KEYWORD_S8, token.KEYWORD_S16, token.KEYWORD_S32, token.KEYWORD_S64,
		token.KEYWORD_U8, token.KEYWORD_U16, token.KEYWORD_U32, token.KEYWORD_U64:

		p.nextToken()
		os.Value = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	case token.KEYWORD_LIST:
		if !p.expectNextToken(token.KEYWORD_LIST) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_LIST, got %s", p.peekToken.Type))
			return nil
		}

		os.Value = &ast.Child{Token: p.curToken, Value: p.parseListShape()}
	case token.KEYWORD_OPTION:
		if !p.expectNextToken(token.KEYWORD_OPTION) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_OPTION, got %s", p.peekToken.Type))
			return nil
		}

		os.Value = &ast.Child{Token: p.curToken, Value: p.parseOptionShape()}
	case token.KEYWORD_RESULT:
		// p.parseResultShape()
	case token.KEYWORD_TUPLE:
	// p.parseTupleShape()
	default:
	}

	if !p.expectNextToken(token.OP_BRACKET_ANGLE_RIGHT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_ANGLE_RIGHT, got %s", p.peekToken.Type))
		return nil
	}

	return os
}

func (p *Parser) parseListShape() *ast.ListShape {
	ls := new(ast.ListShape)
	ls.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_BRACKET_ANGLE_LEFT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_ANGLE_LEFT, got %s", p.peekToken.Type))
		return nil
	}

	switch p.peekToken.Type {
	case token.KEYWORD_STRING, token.KEYWORD_BOOL, token.KEYWORD_CHAR,
		token.KEYWORD_FLOAT32, token.KEYWORD_FLOAT64,
		token.KEYWORD_S8, token.KEYWORD_S16, token.KEYWORD_S32, token.KEYWORD_S64,
		token.KEYWORD_U8, token.KEYWORD_U16, token.KEYWORD_U32, token.KEYWORD_U64:

		p.nextToken()
		ls.Value = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	case token.KEYWORD_LIST:
		if !p.expectNextToken(token.KEYWORD_LIST) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_LIST, got %s", p.peekToken.Type))
			return nil
		}

		ls.Value = &ast.Child{Token: p.curToken, Value: p.parseListShape()}
	case token.KEYWORD_OPTION:
		if !p.expectNextToken(token.KEYWORD_OPTION) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_OPTION, got %s", p.peekToken.Type))
			return nil
		}

		ls.Value = &ast.Child{Token: p.curToken, Value: p.parseOptionShape()}
	case token.KEYWORD_RESULT:
		// p.parseResultShape()
	case token.KEYWORD_TUPLE:
	// p.parseTupleShape()
	default:
	}

	if !p.expectNextToken(token.OP_BRACKET_ANGLE_RIGHT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_ANGLE_RIGHT, got %s", p.peekToken.Type))
		return nil
	}

	return ls
}
