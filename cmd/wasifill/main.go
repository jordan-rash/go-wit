package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/lexer"
	"github.com/jordan-rash/go-wit/parser"

	_ "embed"
)

//go:embed "rpc.tmpl"
var rpc string

const witfile string = `package jordan-rash:pingpong@0.1.0

interface types {
  type pong = string
}

interface pingpong {
  use types.{pong}
  ping: func() -> pong
}

world ping-pong {
  export pingpong
}
`

type wasifill struct {
	PackageNamespace string
	PackageContract  string
	Version          string
	Types            []wftype
	Funcs            []wffunc
	Exports          []wfexports
}

type wftype struct {
	Interface string
	Name      string
	Type      string
}

type wffunc struct {
	Interface string
	Name      string
	Input     string
	Output    string
}

func (w wffunc) PubName() string {
	return strings.Title(w.Name)
}

func (w wffunc) PubOutput() string {
	return strings.Title(w.Output)
}

func (w wftype) PubName() string {
	return strings.Title(w.Name)
}

func (w wfexports) PubName() string {
	return strings.Title(w.Name)
}

func (w wasifill) PubName() string {
	return strings.Title(w.PackageContract)
}

type wfexports struct {
	Type string
	Name string
}

func main() {
	p := parser.New(lexer.NewLexer(witfile))
	t := p.Parse()

	if p.Errors() != nil {
		fmt.Printf("parser errors: %s", p.Errors().Error())
		return
	}

	wf := new(wasifill)

	for _, s := range t.Shapes {
		switch s.TokenLiteral() {
		case "package":
			tS, ok := s.(*ast.PackageShape)
			if !ok {
				fmt.Println("package error")
				return
			}

			wf.PackageNamespace = strings.Split(tS.Value, ":")[0]
			wf.PackageContract = strings.Split(tS.Value, ":")[1]
			wf.Version = tS.Version
		case "interface":
			tS, ok := s.(*ast.InterfaceShape)
			if !ok {
				fmt.Println("interface error")
				return
			}
			for _, c := range tS.Children {
				switch c.TokenLiteral() {
				case "type":
					tC, ok := c.(*ast.TypeStatement)
					if !ok {
						fmt.Println("interface type error")
						return
					}

					tT := wftype{
						Interface: tS.Name.Value,
						Name:      tC.Name.TokenLiteral(),
						Type:      tC.Value.TokenLiteral(),
					}

					wf.Types = append(wf.Types, tT)
				case "func":
					tC, ok := c.(*ast.FuncShape)
					if !ok {
						fmt.Println("interface use error")
						return
					}

					tF := wffunc{
						Interface: tS.Name.Value,
						Name:      tC.Name.TokenLiteral(),
						Input:     "",
						Output:    tC.Value.TokenLiteral(),
					}

					wf.Funcs = append(wf.Funcs, tF)
				case "use":
					_, ok := c.(*ast.UseShape)
					if !ok {
						fmt.Println("interface use error")
						return
					}
					// fmt.Println("\t", tC.TokenLiteral(), tC.Value.TokenLiteral())
				default:
					fmt.Printf("interface child error: %T\n", c)
					return
				}
			}

		case "world":
			tS, ok := s.(*ast.WorldShape)
			if !ok {
				fmt.Println("world error")
				return
			}
			for _, c := range tS.Children {
				switch c.TokenLiteral() {
				case "export":
					tC, ok := c.(*ast.ExportShape)
					if !ok {
						fmt.Println("interface use error")
						return
					}
					tE := wfexports{
						Type: "function",
						Name: tC.Value.TokenLiteral(),
					}
					wf.Exports = append(wf.Exports, tE)
				}
			}

		default:
			fmt.Println("error: invalid root token")
			return
		}

	}

	err := generateFiles(wf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func generateFiles(wf *wasifill) error {
	_, err := os.Stat("gen")
	if os.IsNotExist(err) {
		if err := os.Mkdir("gen", os.ModePerm); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	tmpl, err := template.New("rpc.tmpl").Parse(rpc)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("./gen/gen.go")
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, wf)
	if err != nil {
		panic(err)
	}

	return nil
}
