package vm

import (
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

const StackSize = 2048

var (
	// True + False are Immutable unique values, so we define them globally here.
	// No need to create multiple boolean objects when we can just reference
	// these instead.
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

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
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpPop:
			// instruction we use after expression statements
			// to keep our stack cleaned up if the expr result isn't used
			vm.pop()
		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}
		case code.OpGreaterThan, code.OpNotEqual, code.OpEqual:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}

		}

	}
	return nil
}

// executeComparison executes a comparison based on the type of the operands
// sitting on the top of the stack. It supports integer and boolean comparison
func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.INTEGER_OBJ || right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(left == right))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(left != right))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)",
			op, left.Type(), right.Type())
	}
}

// executeIntegerComparison compares two integers and returns the result
func (vm *VM) executeIntegerComparison(
	op code.Opcode,
	left, right object.Object,
) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue == leftValue))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue != leftValue))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}
	return False
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	rightType := right.Type()
	leftType := left.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

// executeBinaryIntegerOperation takes an opcode + it's operands and performs
// the correct operation on those operands. The result will be pushed to the stack.
func (vm *VM) executeBinaryIntegerOperation(
	op code.Opcode,
	left, right object.Object,
) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	}

	return vm.push(&object.Integer{Value: result})
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

// LastPoppedStackElem is a test only method which returns the top of the stack.
// Since we don't explicitly clear the stack off after we use it, we can see
// the item that was last popped off the stack here.
func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}
