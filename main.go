package main

import (
	"fmt"
	"os"
	"time"
	"github.com/mistium/raingoer/parser"
	"github.com/mistium/raingoer/interpreter"
	"github.com/mistium/raingoer/ast"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: raingoer <filename>")
		os.Exit(1)
	}

	fi, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	code := string(fi)

	p := parser.NewLineParser(code)
	program := p.Parse()

	if len(os.Args) > 2 && os.Args[2] == "--ast" {
		fmt.Println("AST:")
		fmt.Println(program.String())
		return
	}

	interp := interpreter.New()
	
	start := time.Now()

	for _, stmt := range program.Statements {
		interp.Interpret(&ast.Program{Statements: []ast.Statement{stmt}})
	}
	
	duration := time.Since(start)
	fmt.Printf("Execution time: %v\n", duration)
}
