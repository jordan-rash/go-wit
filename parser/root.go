package parser

import (
	"strings"

	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

func (p *Parser) parseInterfaceShape() *ast.InterfaceShape {
	iFace := new(ast.InterfaceShape)
	iFace.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	iFace.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_LEFT) {
		return nil
	}

CHILDLOOP:
	for {
		switch p.peekToken.Type {
		case token.KEYWORD_TYPE:
			if !p.expectNextToken(token.KEYWORD_TYPE) {
				return nil
			}
			iFace.Children = append(iFace.Children, p.parseTypeStatement())
		case token.KEYWORD_USE:
			if !p.expectNextToken(token.KEYWORD_USE) {
				return nil
			}
			iFace.Children = append(iFace.Children, p.parseUseShape())
		case token.KEYWORD_RECORD:
			if !p.expectNextToken(token.KEYWORD_RECORD) {
				return nil
			}
		case token.KEYWORD_VARIANT:
			if !p.expectNextToken(token.KEYWORD_VARIANT) {
				return nil
			}
		case token.KEYWORD_UNION:
			if !p.expectNextToken(token.KEYWORD_UNION) {
				return nil
			}
		case token.KEYWORD_ENUM:
			if !p.expectNextToken(token.KEYWORD_ENUM) {
				return nil
			}
		case token.IDENTIFIER: // derp: func() -> foo
			s := p.parseFuncLine()
			if s != nil {
				iFace.Children = append(iFace.Children, s)
			}
		case token.OP_BRACKET_CURLY_RIGHT:
			if !p.expectNextToken(token.OP_BRACKET_CURLY_RIGHT) {
				return nil
			}
			break CHILDLOOP
		}
	}

	p.nextToken()
	return iFace
}

func (p *Parser) parseWorldShape() *ast.WorldShape {
	world := new(ast.WorldShape)
	world.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	world.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_LEFT) {
		return nil
	}

	// TODO: NEED TO HANDLE ALL POSSIBLE TYPES
	switch p.peekToken.Type {
	case token.KEYWORD_EXPORT:
		if !p.expectNextToken(token.KEYWORD_EXPORT) {
			return nil
		}

		world.Children = append(world.Children, p.parseExportStatement())
	}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_RIGHT) {
		return nil
	}

	return world
}

func (p *Parser) parseUseShape() *ast.UseShape {
	stmt := new(ast.UseShape)
	stmt.Token = p.curToken

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_PERIOD) {
		return nil
	}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_LEFT) {
		return nil
	}

	// TODO: this can be 0..n identifiers
	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	stmt.Value = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_BRACKET_CURLY_RIGHT) {
		return nil
	}

	return stmt
}

func (p *Parser) parsePackageShape() *ast.PackageShape {
	pkg := new(ast.PackageShape)
	pkg.Token = p.curToken
	pkg.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	value := strings.Builder{}

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	value.WriteString(p.curToken.Literal)

	if !p.expectNextToken(token.OP_COLON) {
		return nil
	}

	value.WriteString(p.curToken.Literal)

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	value.WriteString(p.curToken.Literal)

	pkg.Value = value.String()

	if p.peekToken.Literal != token.OP_AT {
		pkg.Version = ""
		return pkg
	}

	if !p.expectNextToken(token.OP_AT) {
		return nil
	}

	sv := p.parseSemVer()
	pkg.Version = sv.String()

	return pkg
}
