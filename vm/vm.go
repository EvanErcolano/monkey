package vm

import (
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

const StackSize = 2048

// VM is our virtual machine utilizing a stack machine architecture. It holds a
// constant pool and the instructions emitted by our bytecode compiler.
type VM struct {
	constants    []object.Object   // the VM's constant pool
	instructions code.Instructions // the instructions given to the VM

	stack []object.Object
	sp    int // Always points to the next value. Top of stack is stack[sp-1]
}

// New initializes our Virtual machine with bytecode
func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

// Run initiates our VM's fetch-decode-execute cycle.
func (vm *VM) Run() error {
	// FETCH our instruction
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		// DECODE the instruction
		switch op {
		case code.OpConstant:
			// EXECUTE the instruction based off of the opcode
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			// push our constants onto the stack
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd:
			// pop last two instructions from the stack into mem
			right := vm.pop()
			left := vm.pop()
			rightValue := right.(*object.Integer).Value
			leftValue := left.(*object.Integer).Value

			result := leftValue + rightValue
			vm.push(&object.Integer{Value: result})
		}
	}
	return nil
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

// StackTop returns our stack's top element. The top element in the stack
// is the sp -1 position. Our stack pointer, sp, always points to the next
// spot available.
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}
