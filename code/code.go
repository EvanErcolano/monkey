package code

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte // aliasing bytes to type opcode

const (
	OpConstant Opcode = iota
)

// Definition helps make our opcodes readable and
// provides info on # of bytes each operand takes up
type Definition struct {
	Name          string
	OperandWidths []int
}

// Lookup for how many operands an opcode has and
// what its human-readable name is
var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}}, // only operand is 2 bytes wide
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1 // for the initial opcode
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}
	return instruction
}
