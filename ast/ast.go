package ast

type AST struct {
	Package    PackageNode
	World      WorldNode
	Uses       []UseNode
	Interfaces []InterfaceNode
}

func (a *AST) String() string {
	// ret := a.Package.Namespace + "/" + a.Package.Name
	// if a.Package.SemVer != "" {
	// 	ret += "@" + a.Package.SemVer
	// }
	// return ret
	return ""
}

type Node interface {
	TokenLiteral() string
	Validate() bool
}

type InterfaceNode interface {
	Node
	interfaceNode()
}

type WorldNode interface {
	Node
	worldNode()
}

type UseNode interface {
	Node
	useNode()
}

type PackageNode interface {
	Node
	packageNode()
}

// ------- Delete these

type Expression interface {
	Node
	expressionNode()
}

// --------

type Package struct {
	Identifier *Identifier

	Namespace string
	Name      string
	SemVer    string
}

func (p *Package) packageNode()         {}
func (p *Package) Validate() bool       { return true }
func (p *Package) TokenLiteral() string { return p.Identifier.Token.Literal }

type World struct {
	Identifier *Identifier

	Name string

	ExportItems  []*ExportShape
	ImportItems  []*ImportShape
	UseItems     []*UseShape
	TypedefItems []*TypeDef
	IncludeItems []*IncludeShape
}

func (w *World) worldNode()           {}
func (w *World) Validate() bool       { return true }
func (w *World) TokenLiteral() string { return w.Identifier.Token.Literal }

type Interface struct {
	Identifier *Identifier

	Name string

	Items InterfaceItems
}

func (i *Interface) interfaceNode()       {}
func (i *Interface) Validate() bool       { return true }
func (i *Interface) TokenLiteral() string { return i.Identifier.Token.Literal }

type InterfaceItems struct {
	UseItems     []*UseShape
	TypedefItems []*TypeDef
	FuncItems    []*FuncShape
}

func (i *InterfaceItems) expressionNode()      {}
func (i *InterfaceItems) Validate() bool       { return true }
func (i *InterfaceItems) TokenLiteral() string { return "" }

type Use struct {
	Identifier *Identifier

	UseInterface struct {
		Interface Interface
		Items     []Identifier
	}
}

func (u *Use) useNode()             {}
func (u *Use) Validate() bool       { return true }
func (u *Use) TokenLiteral() string { return u.Identifier.Token.Literal }
