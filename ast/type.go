package ast

import (
	"github.com/jordan-rash/go-wit/token"
)

type InterfaceShape struct {
	Token    token.Token
	Name     *Identifier
	Children []Statement
}

func (t *InterfaceShape) shapeNode()           {}
func (t *InterfaceShape) TokenLiteral() string { return t.Token.Literal }

type WorldShape struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (t *WorldShape) shapeNode()           {}
func (t *WorldShape) TokenLiteral() string { return t.Token.Literal }

type UseShape struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (t *UseShape) shapeNode()           {}
func (t *UseShape) TokenLiteral() string { return t.Token.Literal }

type TypeStatement struct {
	Token token.Token
	Value Expression
}

func (t *TypeStatement) statementNode()       {}
func (t *TypeStatement) TokenLiteral() string { return t.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (t *Identifier) expressionNode()      {}
func (t *Identifier) TokenLiteral() string { return t.Token.Literal }
