package ast

import (
	"github.com/jordan-rash/go-wit/token"
)

type Identifier struct {
	Token token.Token
	Value string
}

func (t *Identifier) expressionNode()      {}
func (t *Identifier) TokenLiteral() string { return t.Token.Literal }

type Child struct {
	Token token.Token
	Value Shape
}

func (t *Child) expressionNode()      {}
func (t *Child) TokenLiteral() string { return t.Token.Literal }

// Root shapes

type InterfaceShape struct {
	Token    token.Token
	Name     *Identifier
	Children []Shape
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

// Secondary shapes

type TypeStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (t *TypeStatement) shapeNode()           {}
func (t *TypeStatement) TokenLiteral() string { return t.Token.Literal }

type ListShape struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (t *ListShape) shapeNode()           {}
func (t *ListShape) TokenLiteral() string { return t.Token.Literal }
