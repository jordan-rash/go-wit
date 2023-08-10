package parser

import (
	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

func (p *Parser) parseWorldShape() *ast.World {
	world := new(ast.World)
	world.Identifier = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}
	world.Name = p.curToken.Literal

	if !p.expectNextToken(token.OP_BRACKET_CURLY_LEFT) {
		return nil
	}

	// TODO: NEED TO HANDLE ALL POSSIBLE TYPES
	for p.peekToken.Type != token.OP_BRACKET_CURLY_RIGHT {
		switch p.peekToken.Type {
		case token.KEYWORD_EXPORT:
			if !p.expectNextToken(token.KEYWORD_EXPORT) {
				return nil
			}
			world.ExportItems = append(world.ExportItems, p.parseExportStatement())

		case token.KEYWORD_IMPORT:
			if !p.expectNextToken(token.KEYWORD_IMPORT) {
				return nil
			}
			world.ImportItems = append(world.ImportItems, p.parseImportStatement())

		}
	}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_RIGHT) {
		return nil
	}

	return world
}
