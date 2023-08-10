package parser

import (
	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

func (p *Parser) parseFuncItem() *ast.FuncShape {
	fs := new(ast.FuncShape)
	fs.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	fs.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_COLON) {
		return nil
	}

	if !p.expectNextToken(token.KEYWORD_FUNC) {
		return nil
	}
	fs.Value = p.parseFuncType()

	return fs
}

func (p *Parser) parseFuncType() *ast.FuncType {
	ft := new(ast.FuncType)
	ft.Token = p.curToken

	ft.ParamList = p.parseParamList()

	if !p.expectNextToken(token.OP_ARROW) {
		ft.ResultList = nil
		return ft // result list can be empty
	}

	expectComma := false
	// result-list
	if p.expectNextToken(token.OP_BRACKET_PAREN_LEFT) {
		for p.expectNextToken(token.OP_BRACKET_PAREN_RIGHT) {
			switch p.peekToken.Type {
			case token.IDENTIFIER:
				*ft.ResultList = append(*ft.ResultList, p.parseNamedType())
				expectComma = true
			case token.OP_COMMA:
				if expectComma {
					if !p.expectNextToken(token.OP_COMMA) {
						return nil
					}
					expectComma = false
					p.nextToken() // eat comma
					break
				}
				fallthrough
			default:
				// invalid token
				p.nextToken()
			}
		}
	} else {
		ft.ResultList = &ast.ResultList{p.parseTy()}
	}

	return ft
}

func (p *Parser) parseNamedType() *ast.NamedType {
	nt := new(ast.NamedType)

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	nt.Id = &ast.Identifier{Token: p.curToken}

	if !p.expectNextToken(token.OP_COLON) {
		return nil
	}

	nt.Ty = p.parseTypeShape()

	return nt
}

func (p *Parser) parseParamList() *ast.ParamList {
	paramList := new(ast.ParamList)

	if !p.expectNextToken(token.OP_BRACKET_PAREN_LEFT) {
		return nil
	}

	expectComma := false

	for p.peekToken.Type != token.OP_BRACKET_PAREN_RIGHT {
		switch p.peekToken.Type {
		case token.IDENTIFIER:
			if !p.expectNextToken(token.IDENTIFIER) {
				return nil
			}

			*paramList = append(*paramList, p.parseNamedType())
			expectComma = true
		case token.OP_COMMA:
			if expectComma {
				if !p.expectNextToken(token.OP_COMMA) {
					return nil
				}
				expectComma = false
				p.nextToken() // eat comma
				break
			}
			fallthrough
		default:
			// invalid token
			p.nextToken()
		}
	}

	if !p.expectNextToken(token.OP_BRACKET_PAREN_RIGHT) {
		return nil
	}

	return paramList
}
