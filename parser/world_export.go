package parser

import (
	"strings"

	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

// export-item ::= 'export' id ':' extern-type
//               | 'export' interface
//
//interface ::= id
//            | id ':' id '/' id ('@' valid-semver)?
//
//extern-type ::= func-type | 'interface' '{' interface-items* '}'
//
// 'export' id ':' 'func' param-list result-list
// 'export' id ':' 'interface' '{' interface-items* '}'
// 'export' id
// 'export' id ':' id '/' id ('@' valid-semver)?

func (p *Parser) parseExportStatement() *ast.ExportShape {
	es := new(ast.ExportShape)
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
		if !p.expectNextToken(token.KEYWORD_FUNC) {
			return nil
		}
		es.Value = p.parseFuncType()
	case token.KEYWORD_INTERFACE:
		if !p.expectNextToken(token.KEYWORD_INTERFACE) {
			return nil
		}

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
