package parser

import (
	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

// enum-items ::= 'enum' id '{' enum-cases '}'
//
// enum-cases ::= id
//              | id ',' enum-cases?

func (p *Parser) parseEnumShape() *ast.EnumShape {
	es := new(ast.EnumShape)
	es.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	es.Name = &ast.Identifier{Token: p.curToken}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_LEFT) {
		return nil
	}
	expectComma := false
	for p.peekToken.Type != token.OP_BRACKET_CURLY_RIGHT {
		switch p.peekToken.Type {
		case token.IDENTIFIER:
			es.Value = append(es.Value, p.parseTy())
			expectComma = true
		case token.OP_COMMA:
			if expectComma {
				if !p.expectNextToken(token.OP_COMMA) {
					return nil
				}

				expectComma = false
				break
			}
			fallthrough
		default:
			//invalid token
			p.nextToken()
		}
	}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_RIGHT) {
		return nil
	}

	return es
}
