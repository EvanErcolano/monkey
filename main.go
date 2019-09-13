package main

import (
	"fmt"
	"monkey/repl"
	"os"
)

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func main() {
	fmt.Printf("Hello! This is the Monkey Programming Language!\n %s", MONKEY_FACE)

	fmt.Printf("This is the REPL. Type in some Monkey commands!\n")
	repl.Start(os.Stdin, os.Stdout)
}
