package parser

import (
	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

// variant-items ::= 'variant' id '{' variant-cases '}'
//
// variant-cases ::= variant-case
//                 | variant-case ',' variant-cases?
//
// variant-case ::= id
//                | id '(' ty ')'

func (p *Parser) parseVariantShape() *ast.VariantShape {
	vs := new(ast.VariantShape)
	vs.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	vs.Identifier = &ast.Identifier{Token: p.curToken}
	if !p.expectNextToken(token.OP_BRACKET_CURLY_LEFT) {
		return nil
	}

	expectComma := false
	for p.peekToken.Type != token.OP_BRACKET_CURLY_RIGHT {
		switch p.peekToken.Type {
		case token.OP_COMMA:
			if expectComma {
				if !p.expectNextToken(token.OP_COMMA) {
					return nil
				}
				expectComma = false
				break
			}
		default:
			vs.Value = append(vs.Value, p.parseVariantCase())
			expectComma = true
		}
	}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_RIGHT) {
		return nil
	}

	return vs
}

func (p *Parser) parseVariantCase() *ast.VariantCase {
	vc := new(ast.VariantCase)
	vc.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	vc.Identifier = &ast.Identifier{Token: p.curToken}

	if !p.expectNextToken(token.OP_BRACKET_PAREN_LEFT) {
		vc.Value = nil
		return vc
	}

	vc.Value = p.parseTy()

	if !p.expectNextToken(token.OP_BRACKET_PAREN_RIGHT) {
		return nil
	}

	return vc
}
