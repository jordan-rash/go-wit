package parser

import (
	"errors"
	"fmt"

	"github.com/jordan-rash/go-wit/token"
)

type semVer struct {
	Major token.Token
	Minor token.Token
	Patch token.Token
}

func (s semVer) String() string {
	return s.Major.Literal + "." + s.Minor.Literal + "." + s.Patch.Literal
}

// pulled from https://semver.org/#backusnaur-form-grammar-for-valid-semver-versions

func (p *Parser) parseSemVer() *semVer {
	sv := new(semVer)

	if !p.expectNextToken(token.INT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("failed to parse semver major version: %v", p.peekToken))
	}
	sv.Major = p.curToken

	if !p.expectNextToken(token.OP_PERIOD) {
		p.errors = errors.Join(p.errors, fmt.Errorf("failed to parse period: %v", p.peekToken))
	}
	if !p.expectNextToken(token.INT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("failed to parse semver minor version: %v", p.peekToken))
	}
	sv.Minor = p.curToken

	if !p.expectNextToken(token.OP_PERIOD) {
		p.errors = errors.Join(p.errors, fmt.Errorf("failed to parse period: %v", p.peekToken))
	}
	if !p.expectNextToken(token.INT) {
		p.errors = errors.Join(p.errors, fmt.Errorf("failed to parse semver patch version: %v", p.peekToken))
	}

	sv.Patch = p.curToken

	// TODO: got lazy, this needs to be completed for pre/build releases
	// switch p.peekToken.Type {
	// case token.OP_MINUS:
	// case token.OP_PLUS:
	// default:
	// 	return sv
	// }

	p.nextToken()
	return sv
}
