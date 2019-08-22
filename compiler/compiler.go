package compiler

import (
	"monkey/ast"
	"monkey/code"
	"monkey/object"
)

// Compiler is a bytecode compiler
type Compiler struct {
	instructions code.Instructions
	constants    []object.Object // our constant pool
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

// Compile our ast down to it's respective bytecode representation.
func (c *Compiler) Compile(node ast.Node) error {
	return nil
}

func (c *Compiler) ByteCode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
