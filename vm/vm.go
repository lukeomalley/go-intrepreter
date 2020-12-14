package vm

import (
	"fmt"

	"github.com/lukeomalley/monkey_lang/code"
	"github.com/lukeomalley/monkey_lang/compiler"
	"github.com/lukeomalley/monkey_lang/object"
)

// True is the global true value referenced throughout the vm
var True = &object.Boolean{Value: true}

// False is the global false value referenced throughout the vm
var False = &object.Boolean{Value: false}

// Null is the global null value referenced through the vm
var Null = &object.Null{}

// StackSize sets the maximum size of the stack
const StackSize = 2048

// GlobalsSize sets the maximum size of the global variable store
const GlobalsSize = 65536

// VM takes bytecode instrutions and evaluates them
type VM struct {
	constants    []object.Object
	instructions code.Instructions
	globals      []object.Object
	stack        []object.Object
	sp           int // Always points to the next value, top of stack is (sp - 1)
}

// New constructs a VM
func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		globals:      make([]object.Object, GlobalsSize),
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

// NewWithGlobalsStore constructs a new VM with the globals from a previous instance of a VM
func NewWithGlobalsStore(bytecode *compiler.Bytecode, globals []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = globals
	return vm
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
				return err
			}
		case code.OpEqual, code.OpGreaterThan, code.OpNotEqual:
			err := vm.executeComparison(op)
			if err != nil {
				return err
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
		case code.OpMinus:
			err := vm.executeMinusOperator(op)
			if err != nil {
				return err
			}
		case code.OpBang:
			err := vm.executeBangOperator(op)
			if err != nil {
				return err
			}
		case code.OpJump:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1 // Since the loop will increment the ip
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2 // Jump over the operand to the conditional jump

			condition := vm.pop()
			if !isTruthy(condition) {
				ip = pos - 1
			}
		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			// Decode the index operand
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			// Set the index in the globals slice to the value on the stack
			vm.globals[globalIndex] = vm.pop()

		case code.OpGetGlobal:
			// Decode the index operand
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			// Push the value onto the stack
			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err

			}

		case code.OpArray:
			numElements := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements

			err := vm.push(array)
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

func (vm *VM) buildArray(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return &object.Array{Elements: elements}
}

func (vm *VM) executeMinusOperator(op code.Opcode) error {
	operand := vm.pop()

	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}

	value := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
}

func (vm *VM) executeBangOperator(op code.Opcode) error {
	operand := vm.pop()
	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

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

	if leftType == object.STRING_OBJ && rightType == object.STRING_OBJ {
		return vm.executeBinaryStringOperation(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

func (vm *VM) executeBinaryStringOperation(op code.Opcode, left, right object.Object) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operation %d", op)
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return vm.push(&object.String{Value: leftValue + rightValue})
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

func (vm *VM) executeComparison(op code.Opcode) error {
	// Pop two values off of the stack
	right := vm.pop()
	left := vm.pop()

	// Execute the comparison
	if left.Type() == object.INTEGER_OBJ || right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(left == right))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(left != right))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(leftValue == rightValue))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(leftValue != rightValue))
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
