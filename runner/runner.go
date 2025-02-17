package runner

import (
	"kaze/eval"
	"kaze/lexer"
	"kaze/object"
	"kaze/parser"
	"os"
)

func RunFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	l := lexer.New(string(bytes))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			println(err)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	println(eval.Eval(program, env).Inspect())
}
