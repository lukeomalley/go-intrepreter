package vm

import (
	"fmt"

	"github.com/lukeomalley/monkey_lang/code"
	"github.com/lukeomalley/monkey_lang/compiler"
	"github.com/lukeomalley/monkey_lang/object"
)

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

		case code.OpAdd:
			// Pop two elements off the stack
			right := vm.pop()
			left := vm.pop()

			// Add the elements
			leftValue := left.(*object.Integer).Value
			rightValue := right.(*object.Integer).Value
			result := leftValue + rightValue

			// Push the result onto the stack
			vm.push(&object.Integer{Value: result})

		case code.OpPop:
			vm.pop()
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

// LastPoppedStackElem Used only for tests to examine the values popped off the stack
func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}
