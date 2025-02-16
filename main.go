package main

import (
	"kaze/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
