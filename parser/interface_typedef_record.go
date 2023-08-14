package parser

import (
	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

// record-item ::= 'record' id '{' record-fields '}'
//
// record-fields ::= record-field
//                 | record-field ',' record-fields?
//
// record-field ::= id ':' ty

func (p *Parser) parseRecordShape() *ast.RecordShape {
	rs := new(ast.RecordShape)
	rs.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	rs.Identifier = &ast.Identifier{Token: p.curToken}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_LEFT) {
		return nil
	}

	expectComma := false
	for p.peekToken.Type != token.OP_BRACKET_CURLY_RIGHT {
		switch p.peekToken.Type {
		case token.IDENTIFIER:
			rf := new(ast.RecordField)
			rf.Token = p.curToken

			if !p.expectNextToken(token.IDENTIFIER) {
				return nil
			}
			rf.Identifier = &ast.Identifier{Token: p.curToken}

			if !p.expectNextToken(token.OP_COLON) {
				return nil
			}

			rf.Ty = p.parseTy()
			rs.Value = append(rs.Value, rf)
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
			// invalid token
			p.nextToken()
		}
	}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_RIGHT) {
		return nil
	}
	return rs
}
