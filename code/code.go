package code

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte

// Opcode is an 8-bit encoding of an operation
type Opcode byte

// OpConstant encodes a constant
const (
	OpConstant = iota
)

// Definition of the Opcodes used within the virtual stack machine
type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
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
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}

		offset += width
	}

	return instruction
}
