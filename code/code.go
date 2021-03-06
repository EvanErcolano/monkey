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
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

type Opcode byte

const (
	OpConstant Opcode = iota
	OpNull
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

	OpSetGlobal
	OpGetGlobal
	OpSetLocal
	OpGetLocal

	OpArray
	OpHash
	OpIndex

	OpCall
	OpReturnValue
	OpReturn // nothing to return -> return null

	OpGetBuiltin
	OpClosure
	OpGetFree

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
	OpNull:     {"OpNull", []int{}},
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

	OpPop:           {"OpPop", []int{}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpJump:          {"OpJump", []int{2}},

	OpSetGlobal: {"OpSetGlobal", []int{2}},
	OpGetGlobal: {"OpGetGlobal", []int{2}},
	OpSetLocal:  {"OpSetLocal", []int{1}}, // 256 locals possible per func
	OpGetLocal:  {"OpGetLocal", []int{1}},

	OpArray: {"OpArray", []int{2}}, // max len of list is 2^16
	OpHash:  {"OpHash", []int{2}},
	OpIndex: {"OpIndex", []int{}},

	OpCall:        {"OpCall", []int{1}},
	OpReturnValue: {"OpReturnValue", []int{}},
	OpReturn:      {"OpReturn", []int{}},

	OpGetBuiltin: {"OpGetBuiltin", []int{1}}, // 256 possible builtins

	// OpClosure has 2 operands:
	// 1) constant index which specifies where in the constant pool the func is
	// 2) number of free variables needed for the closure
	OpClosure: {"OpClosure", []int{2, 1}},
	OpGetFree: {"OpGetFree", []int{1}},
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
		case 1:
			instruction[offset] = byte(o)
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
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

// ReadUint16 reads a byteslice and converts it to an integer.
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func ReadUint8(ins Instructions) uint8 {
	return uint8(ins[0])
}
