package vm

import (
	"github.com/lukeomalley/monkey_lang/code"
	"github.com/lukeomalley/monkey_lang/object"
)

// Frame represents a call frame in the virtual machine
type Frame struct {
	fn          *object.CompiledFunction
	ip          int
	basePointer int
}

// NewFrame constructs a new vm frame
func NewFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	return &Frame{fn: fn, ip: -1, basePointer: basePointer}
}

// Instructions returns the bytecode instructions for the function of a partirular vm frame
func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
