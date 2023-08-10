package parser

import (
	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

// flags-items ::= 'flags' id '{' flags-fields '}'
//
// flags-fields ::= id
//                | id ',' flags-fields?

func (p *Parser) parseFlagShape() *ast.FlagShape {
	fs := new(ast.FlagShape)
	fs.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	fs.Name = &ast.Identifier{Token: p.curToken}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_LEFT) {
		return nil
	}
	expectComma := false
	for p.peekToken.Type != token.OP_BRACKET_CURLY_RIGHT {
		switch p.peekToken.Type {
		case token.IDENTIFIER:
			fs.Value = append(fs.Value, p.parseTy())
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

	return fs
}
