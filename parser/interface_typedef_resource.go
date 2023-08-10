// Parsing interface typedef: resource
// Documentation: https://github.com/WebAssembly/component-model/blob/main/design/mvp/WIT.md#item-resource
package parser

import (
	"errors"
	"fmt"

	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

// resource-item ::= 'resource' id resource-methods?
// resource-methods ::= '{' resource-method* '}'
// resource-method ::= func-item
//                   | id ':' 'static' func-type
//                   | 'constructor' param-list
//
// param-list ::= '(' named-type-list ')'
//
//
// func-item ::= id ':' func-type
// func-type ::= 'func' param-list result-list

func (p *Parser) parseResourceShape() *ast.ResourceShape {
	resource := new(ast.ResourceShape)
	resource.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected IDENT, got %s", p.peekToken.Type))
		return nil
	}

	resource.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_LEFT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_CURLY_LEFT, got %s", p.peekToken.Type))
		return nil
	}

	for p.peekToken.Type != token.OP_BRACKET_CURLY_RIGHT {
		switch p.peekToken.Type {
		case token.IDENTIFIER:
			if !p.expectNextToken(token.IDENTIFIER) {
				return nil
			}
			if !p.expectNextToken(token.OP_COLON) {
				return nil
			}
			switch p.peekToken.Type {
			case token.KEYWORD_FUNC:
				resource.Value = append(resource.Value, p.parseFuncType())

			case token.KEYWORD_STATIC:
				if !p.expectNextToken(token.KEYWORD_STATIC) {
					return nil
				}
				if !p.expectNextToken(token.KEYWORD_FUNC) {
					return nil
				}
				resource.Value = append(resource.Value, p.parseFuncType())

			default:
				// invalid token
				p.nextToken()
			}

		case token.KEYWORD_CONSTRUCTOR:
			if !p.expectNextToken(token.KEYWORD_CONSTRUCTOR) {
				return nil
			}

			resource.Value = append(resource.Value, p.parseParamList())

		default:
			// invalid token
			p.nextToken()
		}
	}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_RIGHT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_CURLY_RIGHT, got %s", p.peekToken.Type))
		return nil
	}

	p.nextToken() // eat RIGHT CURLY
	return resource
}
