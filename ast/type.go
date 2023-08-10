package ast

import (
	"github.com/jordan-rash/go-wit/token"
)

type Identifier struct {
	Location int

	Token token.Token
	Alias string
	Value string
}

func (t *Identifier) expressionNode()      {}
func (t *Identifier) Validate() bool       { return true }
func (t *Identifier) TokenLiteral() string { return t.Token.Literal }

type Ty struct {
	Name  *Identifier
	Token token.Token
	Value Expression
}

func (t *Ty) expressionNode()      {}
func (t *Ty) Validate() bool       { return true }
func (t *Ty) TokenLiteral() string { return t.Token.Literal }

// Root shapes

type TopUseShape struct {
	Token token.Token
	Name  *Identifier
	Value Expression //TODO this should be a Shape
}

func (t *TopUseShape) useNode()             {}
func (t *TopUseShape) TokenLiteral() string { return t.Token.Literal }

// Secondary interface shapes

type TypeDef struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (t *TypeDef) interfaceNode()       {}
func (t *TypeDef) TokenLiteral() string { return t.Token.Literal }

type TypeShape struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

type UseShape struct {
	Token token.Token
	Name  *Identifier
	Value Expression //TODO this should be a Shape
}

func (t *UseShape) interfaceNode()       {}
func (t *UseShape) TokenLiteral() string { return t.Token.Literal }

func (t *TypeShape) expressionNode()      {}
func (t *TypeShape) Validate() bool       { return true }
func (t *TypeShape) TokenLiteral() string { return t.Token.Literal }

type ListShape struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (t *ListShape) expressionNode()      {}
func (t *ListShape) Validate() bool       { return true }
func (t *ListShape) TokenLiteral() string { return t.Token.Literal }

type OptionShape struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (t *OptionShape) expressionNode()      {}
func (t *OptionShape) Validate() bool       { return true }
func (t *OptionShape) TokenLiteral() string { return t.Token.Literal }

type ResultShape struct {
	Token    token.Token
	Name     *Identifier
	OkValue  Expression
	ErrValue Expression
}

func (t *ResultShape) expressionNode()      {}
func (t *ResultShape) Validate() bool       { return true }
func (t *ResultShape) TokenLiteral() string { return t.Token.Literal }

type TupleShape struct {
	Token token.Token
	Name  *Identifier
	Value []Expression
}

func (t *TupleShape) expressionNode()      {}
func (t *TupleShape) Validate() bool       { return true }
func (t *TupleShape) TokenLiteral() string { return t.Token.Literal }

type ExportShape struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (t *ExportShape) worldNode()           {}
func (t *ExportShape) Validate() bool       { return true }
func (t *ExportShape) TokenLiteral() string { return t.Token.Literal }

type ImportShape struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (t *ImportShape) shapeNode()           {}
func (t *ImportShape) TokenLiteral() string { return t.Token.Literal }

type IncludeShape struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (t *IncludeShape) shapeNode()           {}
func (t *IncludeShape) TokenLiteral() string { return t.Token.Literal }

type FuncShape struct {
	Token  token.Token
	Name   *Identifier
	Static bool
	Value  Expression
}

func (t *FuncShape) expressionNode()      {}
func (t *FuncShape) Validate() bool       { return true }
func (t *FuncShape) TokenLiteral() string { return t.Token.Literal }

type FuncType struct {
	Token      token.Token
	Name       *Identifier
	ParamList  *ParamList
	ResultList *ResultList
}

func (t *FuncType) expressionNode()      {}
func (t *FuncType) Validate() bool       { return true }
func (t *FuncType) TokenLiteral() string { return t.Token.Literal }

type ParamList []Expression

func (t *ParamList) expressionNode()      {}
func (t *ParamList) Validate() bool       { return true }
func (t *ParamList) TokenLiteral() string { return "" }

type ResultList []Expression

func (t *ResultList) expressionNode()      {}
func (t *ResultList) Validate() bool       { return true }
func (t *ResultList) TokenLiteral() string { return "" }

type ResourceShape struct {
	Token token.Token
	Name  *Identifier
	Value []Expression
}

func (t *ResourceShape) interfaceNode()       {}
func (t *ResourceShape) Validate() bool       { return true }
func (t *ResourceShape) TokenLiteral() string { return t.Token.Literal }

type NamedType struct {
	Token token.Token
	Name  *Identifier

	Id Expression //
	Ty Expression
}

func (t *NamedType) expressionNode()      {}
func (t *NamedType) Validate() bool       { return true }
func (t *NamedType) TokenLiteral() string { return t.Token.Literal }

type EnumShape struct {
	Name  *Identifier
	Token token.Token

	Value []Expression
}

func (t *EnumShape) expressionNode()      {}
func (t *EnumShape) Validate() bool       { return true }
func (t *EnumShape) TokenLiteral() string { return t.Token.Literal }
