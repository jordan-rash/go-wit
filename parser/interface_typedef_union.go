package parser

import (
	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

// union-items ::= 'union' id '{' union-cases '}'
//
// union-cases ::= ty
//               | ty ',' union-cases?

func (p *Parser) parseUnionShape() *ast.UnionShape {
	us := new(ast.UnionShape)

	us.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	us.Name = &ast.Identifier{Token: p.curToken}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_LEFT) {
		return nil
	}

	expectComma := false
	for p.peekToken.Type != token.OP_BRACKET_CURLY_RIGHT {
		switch p.peekToken.Type {
		case token.KEYWORD_STRING, token.KEYWORD_BOOL, token.KEYWORD_CHAR,
			token.KEYWORD_FLOAT32, token.KEYWORD_FLOAT64,
			token.KEYWORD_S8, token.KEYWORD_S16, token.KEYWORD_S32, token.KEYWORD_S64,
			token.KEYWORD_U8, token.KEYWORD_U16, token.KEYWORD_U32, token.KEYWORD_U64,
			token.IDENTIFIER, token.KEYWORD_LIST, token.KEYWORD_OPTION,
			token.KEYWORD_RESULT, token.KEYWORD_TUPLE:

			us.Value = append(us.Value, p.parseTy())
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

	return us
}
