package main

import (
	"fmt"
	"monkey/repl"
	"os"
)

func main() {
	fmt.Printf("Hello! This is the Monkey Programming Language!\n")

	fmt.Printf("Type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
