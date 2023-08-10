package parser

import (
	"fmt"
	"strings"

	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

func (p *Parser) parseImportStatement() *ast.ImportShape {
	es := new(ast.ImportShape)

	es.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	es.Name = &ast.Identifier{Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_COLON) {
		es.Value = nil
		return es
	}

	switch p.peekToken.Type {
	case token.KEYWORD_FUNC:
		es.Value = p.parseFuncType()
	case token.KEYWORD_INTERFACE:
		if !p.expectNextToken(token.KEYWORD_INTERFACE) {
			return nil
		}

		fmt.Println("***", p.peekToken.Literal)
		es.Value = p.parseInterfaceItems()
	case token.IDENTIFIER:
		sb := strings.Builder{}
		sb.WriteString(es.Name.Value + ":")

		if !p.expectNextToken(token.IDENTIFIER) {
			return nil
		}
		sb.WriteString(p.curToken.Literal)

		if !p.expectNextToken(token.OP_SLASH) {
			return nil
		}
		sb.WriteString(p.curToken.Literal)

		if !p.expectNextToken(token.IDENTIFIER) {
			return nil
		}
		sb.WriteString(p.curToken.Literal)

		if p.peekToken.Literal == token.OP_AT {
			if !p.expectNextToken(token.OP_AT) {
				return nil
			}
			sb.WriteString(p.curToken.Literal)
			sv := p.parseSemVer()
			sb.WriteString(sv.String())
		}

		es.Name = &ast.Identifier{Value: sb.String()}
	default:
	}

	return es
}
