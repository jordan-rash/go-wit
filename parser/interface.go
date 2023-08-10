// Parsing interfaces
// Documentation: https://github.com/WebAssembly/component-model/blob/main/design/mvp/WIT.md#item-interface
package parser

import (
	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

func (p *Parser) parseInterfaceShape() *ast.Interface {
	iFace := new(ast.Interface)
	iFace.Identifier = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	iFace.Name = p.curToken.Literal
	iFace.Items = *p.parseInterfaceItems()

	return iFace
}

func (p *Parser) parseInterfaceItems() *ast.InterfaceItems {
	ii := new(ast.InterfaceItems)

	if !p.expectNextToken(token.OP_BRACKET_CURLY_LEFT) {
		return nil
	}

	for p.peekToken.Type != token.OP_BRACKET_CURLY_RIGHT {
		switch p.peekToken.Type {
		// USE ITEMS -----------------------
		case token.KEYWORD_USE:
			if !p.expectNextToken(token.KEYWORD_USE) {
				return nil
			}
			// TODO: check this parser
			ii.UseItems = append(ii.UseItems, p.parseUseShape())

		// FUNC ITEMS ----------------------
		case token.IDENTIFIER: // derp: func() -> foo
			s := p.parseFuncItem()
			if s != nil {
				ii.FuncItems = append(ii.FuncItems, s)
			}

		// TYPEDEFS ------------------------
		default:
			ii.TypedefItems = append(ii.TypedefItems, p.parseTypeDef())
		}
	}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_RIGHT) {
		return nil
	}

	return ii
}
