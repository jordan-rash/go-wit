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

	stmt.Value = p.makeExpression(p.peekToken)

	return stmt
}

func (p *Parser) parseOptionShape() *ast.OptionShape {
	os := new(ast.OptionShape)
	os.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_BRACKET_ANGLE_LEFT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_ANGLE_LEFT, got %s", p.peekToken.Type))
		return nil
	}

	os.Value = p.makeExpression(p.peekToken)

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

	ls.Value = p.makeExpression(p.peekToken)

	if !p.expectNextToken(token.OP_BRACKET_ANGLE_RIGHT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_ANGLE_RIGHT, got %s", p.peekToken.Type))
		return nil
	}

	return ls
}

func (p *Parser) parseTupleShape() *ast.TupleShape {
	ts := new(ast.TupleShape)
	ts.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_BRACKET_ANGLE_LEFT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_ANGLE_LEFT, got %s", p.peekToken.Type))
		return nil
	}

	ts.Value = append(ts.Value, p.makeExpression(p.peekToken))

	for p.peekToken.Type == token.OP_COMMA {
		if !p.expectNextToken(token.OP_COMMA) {
			p.errors = errors.Join(p.errors, fmt.Errorf("expected COMMA, got %s", p.peekToken.Type))
			return nil
		}

		ts.Value = append(ts.Value, p.makeExpression(p.peekToken))

	}

	if !p.expectNextToken(token.OP_BRACKET_ANGLE_RIGHT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_ANGLE_RIGHT, got %s", p.peekToken.Type))
		return nil
	}

	return ts
}

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
		rs.OkValue = p.makeExpression(p.peekToken)
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
	rs.ErrValue = p.makeExpression(p.peekToken)

	if !p.expectNextToken(token.OP_BRACKET_ANGLE_RIGHT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_ANGLE_RIGHT, got %s", p.peekToken.Type))
		return nil
	}

	return rs
}

func (p *Parser) parseExportStatement() *ast.ExportShape {
	es := new(ast.ExportShape)
	es.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	es.Token = p.curToken

	// TODO: this is lazy POC logic.  an export can be more than this
	if !p.expectNextToken(token.IDENTIFIER) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_ANGLE_LEFT, got %s", p.peekToken.Type))
		return nil
	}

	es.Value = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return es
}

func (p *Parser) parseFuncLine() *ast.FuncShape {
	fs := new(ast.FuncShape)
	fs.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	fs.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_COLON) {
		return nil
	}
	if !p.expectNextToken(token.KEYWORD_FUNC) {
		return nil
	}
	if !p.expectNextToken(token.OP_BRACKET_PAREN_LEFT) {
		return nil
	}
	if !p.expectNextToken(token.OP_BRACKET_PAREN_RIGHT) {
		return nil
	}
	if !p.expectNextToken(token.OP_ARROW) {
		return nil
	}

	fs.Value = p.makeExpression(p.peekToken)

	return fs
}
