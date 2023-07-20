package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jordan-rash/go-wit/lexer"
	"github.com/jordan-rash/go-wit/parser"
)

func main() {
	done := make(chan struct{}, 1)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	go parseWit(done)

	select {
	case <-ctxTimeout.Done():
		panic(fmt.Errorf("time limit"))
	case <-done:
		fmt.Println("hit done")
	}
}

func parseWit(done chan struct{}) {
	b, err := os.ReadFile("./core.wit")
	if err != nil {
		panic(err)
	}

	p := parser.New(lexer.NewLexer(string(b)))
	p.Parse()

	if p.Errors() != nil {
		panic(p.Errors())
	}

	done <- struct{}{}
}
