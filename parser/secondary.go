package parser

import (
	"errors"
	"fmt"

	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

func (p *Parser) parseTypeStatement() *ast.TypeStatement {
	stmt := new(ast.TypeStatement)
	stmt.Token = p.curToken

	if !p.expectNextType(token.IDENTIFIER) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected IDENT, got %s", p.peekToken.Type))
		return nil
	}

	if !p.expectNextType(token.OP_EQUAL) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected EQUAL, got %s", p.peekToken.Type))
		return nil
	}

	p.nextToken()

	switch p.curToken.Type {
	case token.KEYWORD_STRING, token.KEYWORD_BOOL, token.KEYWORD_CHAR,
		token.KEYWORD_FLOAT32, token.KEYWORD_FLOAT64,
		token.KEYWORD_S8, token.KEYWORD_S16, token.KEYWORD_S32, token.KEYWORD_S64,
		token.KEYWORD_U8, token.KEYWORD_U16, token.KEYWORD_U32, token.KEYWORD_U64:

		stmt.Value = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case token.KEYWORD_LIST:
		// p.parseListShape()
	case token.KEYWORD_OPTION:
		// p.parseOptionShape()
	case token.KEYWORD_RESULT:
		// p.parseResultShape()
	case token.KEYWORD_TUPLE:
		// p.parseTupleShape()
	default:
		p.errors = errors.Join(p.errors, fmt.Errorf("unexpected token: %s", p.peekToken.Type))
		return nil
	}

	p.nextToken() // eat switch statement token
	return stmt
}
