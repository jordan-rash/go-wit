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
	default:
		rs.OkValue = p.parseTy()
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
	rs.ErrValue = p.parseTy()

	if !p.expectNextToken(token.OP_BRACKET_ANGLE_RIGHT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_ANGLE_RIGHT, got %s", p.peekToken.Type))
		return nil
	}

	return rs
}
