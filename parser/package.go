package parser

import (
	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

func (p *Parser) parsePackageShape() *ast.Package {
	pkg := new(ast.Package)

	pkg.Identifier = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	pkg.Namespace = p.curToken.Literal

	if !p.expectNextToken(token.OP_COLON) {
		return nil
	}

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	pkg.Name = p.curToken.Literal

	if p.peekToken.Literal != token.OP_AT {
		pkg.SemVer = ""
		return pkg
	}

	if !p.expectNextToken(token.OP_AT) {
		return nil
	}

	sv := p.parseSemVer()
	pkg.SemVer = sv.String()

	return pkg
}
