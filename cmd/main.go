package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jordan-rash/go-wit/ast"
	"github.com/jordan-rash/go-wit/lexer"
	"github.com/jordan-rash/go-wit/parser"
)

func main() {
	done := make(chan *ast.AST, 1)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	go parseWit(done)

	var tree *ast.AST

	select {
	case <-ctxTimeout.Done():
		panic(fmt.Errorf("time limit"))
	case tree = <-done:
		// fmt.Println("got tree")
	}

	for _, s := range tree.Shapes {
		switch s.TokenLiteral() {
		case "package":
			tS, ok := s.(*ast.PackageShape)
			if !ok {
				// error
			}
			fmt.Println("Package: ", tS.Value)
			fmt.Println("Version: ", tS.Version)
		case "interface":
			tS, ok := s.(*ast.InterfaceShape)
			if !ok {
				// error
			}
			fmt.Println("Interface: ", tS.Name.Value)
			for _, c := range tS.Children {
				switch c.TokenLiteral() {
				case "type":
					tC, ok := c.(*ast.TypeStatement)
					if !ok {
						// error
					}
					fmt.Println("\t", tC.Name.TokenLiteral(), tC.Value.TokenLiteral())

				case "use":
					tC, ok := c.(*ast.UseShape)
					if !ok {
						// error
					}
					fmt.Println("\t", tC.TokenLiteral(), tC.Value.TokenLiteral())
				default:
					// error
				}
			}
		case "world":
			tS, ok := s.(*ast.WorldShape)
			if !ok {
				// error
			}
			fmt.Println("World: ", tS.Name.Value)
			for _, c := range tS.Children {
				switch c.TokenLiteral() {
				case "export":
					tC, ok := c.(*ast.ExportShape)
					if !ok {
						// error
					}
					fmt.Println("\t", tC.Name.TokenLiteral(), tC.Value.TokenLiteral())

				default:
					// error
				}
			}

		default:
			// error
		}
	}
}

func parseWit(done chan *ast.AST) {
	b, err := os.ReadFile("./pingpong.wit")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsing the following file:\n\n```\n%s\n```\n\n", string(b))

	p := parser.New(lexer.NewLexer(string(b)))
	t := p.Parse()

	if p.Errors() != nil {
		panic(p.Errors())
	}

	done <- t
}
