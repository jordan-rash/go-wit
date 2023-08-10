package parser

import (
	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

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

	stmt.Value = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_RIGHT) {
		return nil
	}

	return stmt
}
