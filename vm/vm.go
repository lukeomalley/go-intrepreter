package vm

import (
	"fmt"

	"github.com/lukeomalley/monkey_lang/code"
	"github.com/lukeomalley/monkey_lang/compiler"
	"github.com/lukeomalley/monkey_lang/object"
)

// True is the global true value used throughout the vm
var True = &object.Boolean{Value: true}

// False is the global false value used throughout the vm
var False = &object.Boolean{Value: false}

// StackSize sets the maximum size of the stack
const StackSize = 2048

// VM takes bytecode instrutions and evaluates them
type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // Always points to the next value, top of stack is (sp - 1)
}

// New constructs a VM
func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

// Run executes the bytecode operations
func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err

			}

		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.executeBianryOperation(op)
			if err != nil {

			}

		case code.OpPop:
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

		}

	}

	return nil
}

// =============================================================================
// Helper Methods
// =============================================================================

func (vm *VM) executeBianryOperation(op code.Opcode) error {
	// Pop two elements off the stack
	right := vm.pop()
	left := vm.pop()
	leftType := left.Type()
	rightType := right.Type()

	// Check types to see which operation to perform
	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error {
	// Convert to integer values
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	// Perform the operation
	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	// Push the result onto the stack
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

// LastPoppedStackElem Used only for tests to examine the values popped off the stack
func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}
