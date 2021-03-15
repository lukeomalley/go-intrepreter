package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Instructions is an alias for a byte slice of instructions
type Instructions []byte

// Opcode is an 8-bit encoding of an operation
type Opcode byte

// Enum for OpCodes
const (
	OpConstant = iota
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpPop
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpMinus
	OpBang
	OpJump
	OpJumpNotTruthy
	OpNull
	OpSetGlobal
	OpGetGlobal
	OpArray
	OpHash
	OpIndex
	OpCall
	OpReturnValue
	OpReturn
	OpSetLocal
	OpGetLocal
	OpGetBuiltin
	OpClosure
	OpGetFree
	OpCurrentClosure
)

// Definition of the Opcodes used within the virtual stack machine
type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant:       {"OpConstant", []int{2}},
	OpPop:            {"OpPop", []int{}},
	OpAdd:            {"OpAdd", []int{}},
	OpSub:            {"OpSub", []int{}},
	OpMul:            {"OpMul", []int{}},
	OpDiv:            {"OpDiv", []int{}},
	OpTrue:           {"OpTrue", []int{}},
	OpFalse:          {"OpFalse", []int{}},
	OpEqual:          {"OpFalse", []int{}},
	OpNotEqual:       {"OpFalse", []int{}},
	OpGreaterThan:    {"OpFalse", []int{}},
	OpMinus:          {"OpMinus", []int{}},
	OpBang:           {"OpBang", []int{}},
	OpJump:           {"OpJump", []int{2}},
	OpJumpNotTruthy:  {"OpJumpNotTruthy", []int{2}},
	OpNull:           {"OpNull", []int{}},
	OpSetGlobal:      {"OpSetGlobal", []int{2}},
	OpGetGlobal:      {"OpGetGlobal", []int{2}},
	OpArray:          {"OpArray", []int{2}},
	OpHash:           {"OpHash", []int{2}},
	OpIndex:          {"OpIndex", []int{}},
	OpCall:           {"OpCall", []int{1}},
	OpReturnValue:    {"OpReturnValue", []int{}},
	OpReturn:         {"OpReturn", []int{}},
	OpSetLocal:       {"OpSetLocal", []int{1}},
	OpGetLocal:       {"OpGetLocal", []int{1}},
	OpGetBuiltin:     {"OpGetBuiltin", []int{1}},
	OpClosure:        {"OpClosure", []int{2, 1}},
	OpGetFree:        {"OpGetFree", []int{1}},
	OpCurrentClosure: {"OpCurrentClosure", []int{}},
}

// Lookup returns the corresponding Opcode for a given byte
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

// Make constructs a btye slice for specified operation
func Make(op Opcode, operands ...int) []byte {
	// Lookup the OpCode
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	// Count the total number of bytes required for the instruction
	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	// Make a byte slice of the required length
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	// Load the operands into the instruction byte slice
	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 1:
			instruction[offset] = byte(o)
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}

		offset += width
	}
	return instruction
}

// ReadOperands transforms bytecodes into their integer representations
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func ReadUint8(ins Instructions) uint8 {
	return uint8(ins[0])
}

// ReadUint16 reads a unsigned 16-bit binary integer
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "Error: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("Error: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("Error: unhandled operand count for %s\n", def.Name)
}
