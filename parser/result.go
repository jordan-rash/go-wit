package parser

import (
	"errors"
	"fmt"

	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

func (p *Parser) parseResultShape() *ast.ResultShape {
	rs := new(ast.ResultShape)
	rs.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// type derp = result
	if p.peekToken.Type != token.OP_BRACKET_ANGLE_LEFT {
		rs.OkValue = nil
		rs.ErrValue = nil
		return rs
	}

	if !p.expectNextToken(token.OP_BRACKET_ANGLE_LEFT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_ANGLE_LEFT, got %s", p.peekToken.Type))
		return nil
	}

	// PARSER OK VALUE ---------------------------
	switch p.peekToken.Type {
	case token.OP_UNDERSCORE:
		if !p.expectNextToken(token.OP_UNDERSCORE) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected UNDERSCORE, got %s", p.peekToken.Type))
			return nil
		}

		rs.OkValue = nil

	case token.KEYWORD_STRING, token.KEYWORD_BOOL, token.KEYWORD_CHAR,
		token.KEYWORD_FLOAT32, token.KEYWORD_FLOAT64,
		token.KEYWORD_S8, token.KEYWORD_S16, token.KEYWORD_S32, token.KEYWORD_S64,
		token.KEYWORD_U8, token.KEYWORD_U16, token.KEYWORD_U32, token.KEYWORD_U64,
		token.IDENTIFIER:

		p.nextToken()
		rs.OkValue = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	case token.KEYWORD_LIST:
		if !p.expectNextToken(token.KEYWORD_LIST) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_LIST, got %s", p.peekToken.Type))
			return nil
		}

		c := &ast.Child{Token: p.curToken}
		c.Value = p.parseListShape()
		rs.OkValue = c

	case token.KEYWORD_OPTION:
		if !p.expectNextToken(token.KEYWORD_OPTION) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_OPTION, got %s", p.peekToken.Type))
			return nil
		}
		rs.OkValue = &ast.Child{Token: p.curToken, Value: p.parseOptionShape()}

	case token.KEYWORD_RESULT:
		if !p.expectNextToken(token.KEYWORD_RESULT) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_RESULT, got %s", p.peekToken.Type))
			return nil
		}

		rs.OkValue = &ast.Child{Token: p.curToken, Value: p.parseResultShape()}
	case token.KEYWORD_TUPLE:
		if !p.expectNextToken(token.KEYWORD_TUPLE) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_TUPLE, got %s", p.peekToken.Type))
			return nil
		}

		rs.OkValue = &ast.Child{Token: p.curToken, Value: p.parseTupleShape()}
	}

	if p.peekToken.Type == token.OP_BRACKET_ANGLE_RIGHT {
		if !p.expectNextToken(token.OP_BRACKET_ANGLE_RIGHT) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_TUPLE, got %s", p.peekToken.Type))
			return nil
		}

		rs.ErrValue = nil
		return rs
	}

	if !p.expectNextToken(token.OP_COMMA) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected COMMA, got %s", p.peekToken.Type))
		return nil
	}

	// PARSER ERR VALUE ---------------------------
	switch p.peekToken.Type {
	case token.KEYWORD_STRING, token.KEYWORD_BOOL, token.KEYWORD_CHAR,
		token.KEYWORD_FLOAT32, token.KEYWORD_FLOAT64,
		token.KEYWORD_S8, token.KEYWORD_S16, token.KEYWORD_S32, token.KEYWORD_S64,
		token.KEYWORD_U8, token.KEYWORD_U16, token.KEYWORD_U32, token.KEYWORD_U64,
		token.IDENTIFIER:

		p.nextToken()
		rs.ErrValue = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	case token.KEYWORD_LIST:
		if !p.expectNextToken(token.KEYWORD_LIST) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_LIST, got %s", p.peekToken.Type))
			return nil
		}

		rs.ErrValue = &ast.Child{Token: p.curToken, Value: p.parseListShape()}
	case token.KEYWORD_OPTION:
		if !p.expectNextToken(token.KEYWORD_OPTION) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_OPTION, got %s", p.peekToken.Type))
			return nil
		}

		c := &ast.Child{Token: p.curToken}
		c.Value = p.parseOptionShape()
		rs.ErrValue = c

	case token.KEYWORD_RESULT:
		if !p.expectNextToken(token.KEYWORD_RESULT) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_RESULT, got %s", p.peekToken.Type))
			return nil
		}

		rs.ErrValue = &ast.Child{Token: p.curToken, Value: p.parseResultShape()}
	case token.KEYWORD_TUPLE:
		if !p.expectNextToken(token.KEYWORD_TUPLE) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected KEYWORD_TUPLE, got %s", p.peekToken.Type))
			return nil
		}

		rs.ErrValue = &ast.Child{Token: p.curToken, Value: p.parseTupleShape()}
	}

	if !p.expectNextToken(token.OP_BRACKET_ANGLE_RIGHT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_ANGLE_RIGHT, got %s", p.peekToken.Type))
		return nil
	}

	return rs
}
