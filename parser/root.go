package parser

import (
	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

func (p *Parser) parseInterfaceShape() *ast.InterfaceShape {
	iFace := new(ast.InterfaceShape)
	iFace.Token = p.curToken
	iFace.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_LEFT) {
		return nil
	}

	switch p.peekToken.Type {
	case token.KEYWORD_TYPE:
		if !p.expectNextToken(token.KEYWORD_TYPE) {
			return nil
		}

		iFace.Children = append(iFace.Children, p.parseTypeStatement())

	case token.KEYWORD_RECORD:
		if !p.expectNextToken(token.KEYWORD_RECORD) {
			return nil
		}
	case token.KEYWORD_VARIANT:
		if !p.expectNextToken(token.KEYWORD_VARIANT) {
			return nil
		}
	case token.KEYWORD_UNION:
		if !p.expectNextToken(token.KEYWORD_UNION) {
			return nil
		}
	case token.KEYWORD_ENUM:
		if !p.expectNextToken(token.KEYWORD_ENUM) {
			return nil
		}
	}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_RIGHT) {
		return nil
	}

	return iFace
}

func (p *Parser) parseWorldShape() *ast.WorldShape {
	stmt := new(ast.WorldShape)
	stmt.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_LEFT) {
		return nil
	}

	// TODO: NEED TO HANDLE WHAT IS INSIDE THE { .... } pg 52

	if !p.expectNextToken(token.OP_BRACKET_CURLY_RIGHT) {
		return nil
	}

	return stmt
}

func (p *Parser) parseUseShape() *ast.UseShape {
	stmt := new(ast.UseShape)
	stmt.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_PERIOD) {
		return nil
	}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_LEFT) {
		return nil
	}

	// TODO: this can be 0..n identifiers
	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_RIGHT) {
		return nil
	}

	return stmt
}
