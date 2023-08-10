package parser

import (
	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

func (p *Parser) parseTypeDef() *ast.TypeDef {
	td := new(ast.TypeDef)
	td.Token = p.curToken

	switch p.peekToken.Type {
	case token.KEYWORD_RESOURCE:
		if !p.expectNextToken(token.KEYWORD_RESOURCE) {
			return nil
		}
		//TODO: IMPLEMENT

	case token.KEYWORD_VARIANT:
		if !p.expectNextToken(token.KEYWORD_VARIANT) {
			return nil
		}
		//TODO: IMPLEMENT

	case token.KEYWORD_RECORD:
		if !p.expectNextToken(token.KEYWORD_RECORD) {
			return nil
		}
		//TODO: IMPLEMENT

	case token.KEYWORD_UNION:
		if !p.expectNextToken(token.KEYWORD_UNION) {
			return nil
		}
		//TODO: IMPLEMENT

	case token.KEYWORD_FLAGS:
		if !p.expectNextToken(token.KEYWORD_FLAGS) {
			return nil
		}
		//TODO: IMPLEMENT

	case token.KEYWORD_ENUM:
		if !p.expectNextToken(token.KEYWORD_ENUM) {
			return nil
		}
		//TODO: IMPLEMENT

	case token.KEYWORD_TYPE:
		if !p.expectNextToken(token.KEYWORD_TYPE) {
			return nil
		}

		td.Value = p.parseTypeShape()

	default:
		// unknown token error
		p.nextToken()
	}

	return td
}
