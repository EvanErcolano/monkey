package vm

import (
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

const StackSize = 2048
const GlobalsSize = 65536
const MaxFrames = 1024

var (
	// True , False, Null are Immutable unique values, so we define them globally.
	// No need to create multiple boolean objects when we can reference one instance.
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
	Null  = &object.Null{}
)

// VM is our virtual machine utilizing a stack machine architecture. It holds a
// constant pool and the instructions emitted by our bytecode compiler.
type VM struct {
	constants []object.Object // the VM's constant pool

	stack []object.Object // Our data stack
	sp    int             // Always points to the next value. Top of stack is stack[sp-1]

	globals []object.Object

	frames      []*Frame // Our frame/call stack
	framesIndex int
}

// New initializes our Virtual machine with bytecode
func New(bytecode *compiler.Bytecode) *VM {
	// Insert our bytecode into the first stack frame
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	return &VM{
		constants: bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: make([]object.Object, GlobalsSize),

		frames:      frames,
		framesIndex: 1,
	}
}

// NewWithGlobalsStore Allows us to create a new VM with a pre-existing
// global store. Used by our REPL.
func NewWithGlobalsStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = s
	return vm
}

// currentFrame peeks the current frame on our call stack
func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}

// Run initiates our VM's fetch-decode-execute cycle.
func (vm *VM) Run() error {
	var ip int
	var ins code.Instructions
	var op code.Opcode

	// FETCH our instruction at our current frame
	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++

		ip = vm.currentFrame().ip
		ins = vm.currentFrame().Instructions()
		op = code.Opcode(ins[ip])

		// DECODE the instruction
		switch op {
		case code.OpConstant:
			// EXECUTE the instruction based off of the opcode
			constIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

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
		case code.OpBang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}
		case code.OpMinus:
			err := vm.executeMinusOperator()
			if err != nil {
				return err
			}
		case code.OpJump:
			// decode the integer address ahead of the opcode
			pos := int(code.ReadUint16(ins[ip+1:]))
			// set the instruction pointer to the target of our jump
			// since the loop will inc ip with each iteration, we set
			// the ip to the position directly before where we want
			vm.currentFrame().ip = pos - 1
		case code.OpJumpNotTruthy:
			// decode the target address of our jumpIfNotTruthy
			pos := int(code.ReadUint16(ins[ip+1:]))
			// inc IP by 2 so we skip over the 16bit operand we decoded above
			// this will put us at the conditional result which we want to test
			vm.currentFrame().ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				vm.currentFrame().ip = pos - 1
			}
		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			vm.globals[globalIndex] = vm.pop()
		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		case code.OpSetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip++

			frame := vm.currentFrame()

			// base ptr is right after our func call, localindex is the unique
			// 0 based index for our local variables.
			vm.stack[frame.basePointer+int(localIndex)] = vm.pop()
		case code.OpGetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip++

			currFrame := vm.currentFrame()

			err := vm.push(vm.stack[currFrame.basePointer+int(localIndex)])
			if err != nil {
				return err
			}

		case code.OpArray:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements

			err := vm.push(array)
			if err != nil {
				return err
			}
		case code.OpHash:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			hash, err := vm.buildHash(vm.sp-numElements, vm.sp)
			if err != nil {
				return err
			}
			vm.sp = vm.sp - numElements

			err = vm.push(hash)
			if err != nil {
				return err
			}
		case code.OpIndex:
			index := vm.pop()
			left := vm.pop()

			err := vm.executeIndexExpression(left, index)
			if err != nil {
				return err
			}
		case code.OpReturnValue:
			returnValue := vm.pop()

			// pop the frame associated w/ the func we executed off the stack
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(returnValue)
			if err != nil {
				return err
			}
		case code.OpReturn:
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpCall:
			vm.currentFrame().ip += 1
			fn, ok := vm.stack[vm.sp-1].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("attempting to call non-function")
			}
			// pass our curr location, "base ptr" / "frame ptr". This is the pos
			// after our function call. We will use this to clean the stack up
			// later, after the func is done executing.
			frame := NewFrame(fn, vm.sp)
			vm.pushFrame(frame)
			// increase our stack ptr to create a "hole" for all of our local
			// variables to fit into.
			vm.sp = frame.basePointer + fn.NumLocals
		}
	}
	return nil
}

// executeIndexExpression indexes into either an array or hash
// TODO: I don't like how many unsafe operations we're just solving with Null...
func (vm *VM) executeIndexExpression(left, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return vm.executeArrayIndex(left, index)
	case left.Type() == object.HASH_OBJ:
		return vm.executeHashIndex(left, index)
	default:
		return fmt.Errorf("index operator not supported: %s", left.Type())
	}
}

func (vm *VM) executeArrayIndex(left, index object.Object) error {
	array := left.(*object.Array)
	i := index.(*object.Integer).Value
	max := int64(len(array.Elements) - 1)

	// pop a Null result onto the stack, if we're out of bounds
	if i < 0 || i > max {
		return vm.push(Null)
	}

	return vm.push(array.Elements[i])
}

func (vm *VM) executeHashIndex(left, index object.Object) error {
	hash := left.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return fmt.Errorf("unusable as hash key: %s", index.Type())
	}

	pair, ok := hash.Pairs[key.HashKey()]
	if !ok {
		return vm.push(Null)
	}

	return vm.push(pair.Value)
}

// buildHash reads a HashMap off of the stack. It reads key first followed by value.
func (vm *VM) buildHash(startIndex, endIndex int) (object.Object, error) {
	hashedPairs := make(map[object.HashKey]object.HashPair)

	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]

		pair := object.HashPair{Key: key, Value: value}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("unusable as hash key: %s", key.Type())
		}

		hashedPairs[hashKey.HashKey()] = pair
	}
	return &object.Hash{Pairs: hashedPairs}, nil
}

func (vm *VM) buildArray(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return &object.Array{Elements: elements}
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		// everything other than False and Null is truthy
		return vm.push(False)
	}
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()
	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}
	value := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
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
	switch {
	case leftType == object.STRING_OBJ && rightType == object.STRING_OBJ:
		return vm.executeBinaryStringOperation(op, left, right)
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return vm.executeBinaryIntegerOperation(op, left, right)
	default:
		return fmt.Errorf("unsupported types for binary operation: %s %s",
			leftType, rightType)
	}
}

func (vm *VM) executeBinaryStringOperation(
	op code.Opcode,
	left, right object.Object,
) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operator :%d", op)
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return vm.push(&object.String{Value: leftValue + rightValue})
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

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}
