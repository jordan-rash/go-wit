package parser

import (
	"errors"
	"fmt"

	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

func (p *Parser) parseInterfaceShape() *ast.InterfaceShape {
	stmt := new(ast.InterfaceShape)
	stmt.Token = p.curToken

	if !p.expectNextType(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextType(token.OP_BRACKET_CURLY_LEFT) {
		return nil
	}

	switch p.peekToken.Type {
	case token.OP_BRACKET_CURLY_RIGHT:
		p.nextToken()
	case token.KEYWORD_TYPE:
		stmt.Children = append(stmt.Children, p.parseTypeStatement())
	//case token.KEYWORD_RECORD:
	//case token.KEYWORD_VARIANT:
	//case token.KEYWORD_UNION:
	//case token.KEYWORD_ENUM:
	default:
		p.errors = errors.Join(p.errors, fmt.Errorf("unexpected token: %s", p.peekToken.Type))
	}

	// TODO: NEED TO HANDLE WHAT IS INSIDE THE { .... } pg 52

	return stmt
}

func (p *Parser) parseWorldShape() *ast.WorldShape {
	stmt := new(ast.WorldShape)
	stmt.Token = p.curToken

	if !p.expectNextType(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextType(token.OP_BRACKET_CURLY_LEFT) {
		return nil
	}

	// TODO: NEED TO HANDLE WHAT IS INSIDE THE { .... } pg 52

	if !p.expectNextType(token.OP_BRACKET_CURLY_RIGHT) {
		return nil
	}

	return stmt
}

func (p *Parser) parseUseShape() *ast.UseShape {
	stmt := new(ast.UseShape)
	stmt.Token = p.curToken

	if !p.expectNextType(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextType(token.OP_PERIOD) {
		return nil
	}

	if !p.expectNextType(token.OP_BRACKET_CURLY_LEFT) {
		return nil
	}

	// TODO: this can be 0..n identifiers
	if !p.expectNextType(token.IDENTIFIER) {
		return nil
	}

	if !p.expectNextType(token.OP_BRACKET_CURLY_RIGHT) {
		return nil
	}

	return stmt
}
