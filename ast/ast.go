package ast

type Node interface {
	TokenLiteral() string
}

type Shape interface {
	Node
	shapeNode()
}

type Expression interface {
	Node
	expressionNode()
}

type AST struct {
	Shapes []Shape
}

func (a *AST) TokenLiteral() string {
	if len(a.Shapes) > 0 {
		return a.Shapes[0].TokenLiteral()
	}
	return ""
}
