package parser

import (
	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

func (p *Parser) parseTopUseShape() *ast.Use {
	u := new(ast.Use)

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	u.Identifier = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

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

	// stmt.Value = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_RIGHT) {
		return nil
	}

	return u
}
