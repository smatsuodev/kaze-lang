package main

import (
	"kaze/repl"
	"kaze/runner"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		runner.RunFile(os.Args[1])
		return
	}
	repl.Start(os.Stdin, os.Stdout)
}
