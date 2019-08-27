package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// bytecode instructions
type Instructions []byte

// String outputs our byte slice instructions in a more readable form
// formatted as the it's index in our instructions and the instruction
func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

// fmtInstruction formats our bytecode instructions
func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

type Opcode byte

const (
	OpConstant Opcode = iota

	OpTrue
	OpFalse

	OpAdd
	OpSub
	OpMul
	OpDiv

	OpEqual
	OpNotEqual
	OpGreaterThan

	OpBang
	OpMinus

	OpJumpNotTruthy
	OpJump
	OpPop
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
	OpConstant: &Definition{"OpConstant", []int{2}}, // only operand is 2 bytes wide
	OpTrue:     {"OpTrue", []int{}},
	OpFalse:    {"OpFalse", []int{}},

	// Infix Expressions
	OpAdd: {"OpAdd", []int{}},
	OpSub: {"OpSub", []int{}},
	OpMul: {"OpMul", []int{}},
	OpDiv: {"OpDiv", []int{}},

	OpEqual:       {"OpEqual", []int{}},
	OpNotEqual:    {"OpNotEqual", []int{}},
	OpGreaterThan: {"OpGreaterThan", []int{}},

	OpBang:  {"OpBang", []int{}},
	OpMinus: {"OpMinus", []int{}},

	OpPop: {"OpPop", []int{}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpJump: {"OpJump", []int{2}},
}

// Lookup looks up an Opcode definition via our definition map
// Either returns a Definition or an opcode undefined error
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

// Make takes an opcode byte and int operands and converts
// them into a slice of bytes repr'ing our bytecode
// Assembling from human readable instructions to our VM bytecode
// Ex. Opconstant will contain:
// [bytecode instruction, 1/2 constant byte, 2/2 constant byte]
// If we add the constant 1 our returned bytecode will be:
// [0, 00000000, 00000001]
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1 // index 0 contains our Opcode bytecode
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	// create a byte slice big enough for our instruction length
	instruction := make([]byte, instructionLen)
	// place our Opcode bytcode
	instruction[0] = byte(op)

	// append our operands as bytecode onto the byte slice
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

// ReadOperands takes an instruction definition and some VM bytecode
// and uses that definition to convert the bytecode back to its
// integer representation. Disassembling back to human readable form.
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

// ReadUint16 reads a byteslice and converts it to an integer.
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
