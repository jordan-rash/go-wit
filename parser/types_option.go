package parser

import (
	"errors"
	"fmt"

	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/token"
)

func (p *Parser) parseOptionShape() *ast.OptionShape {
	os := new(ast.OptionShape)
	os.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.OP_BRACKET_ANGLE_LEFT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_ANGLE_LEFT, got %s", p.peekToken.Type))
		return nil
	}

	os.Value = p.parseTy()

	if !p.expectNextToken(token.OP_BRACKET_ANGLE_RIGHT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("expected BRACKET_ANGLE_RIGHT, got %s", p.peekToken.Type))
		return nil
	}

	return os
}
