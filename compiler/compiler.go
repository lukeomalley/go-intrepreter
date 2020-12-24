package compiler

import (
	"fmt"
	"sort"

	"github.com/lukeomalley/monkey_lang/ast"
	"github.com/lukeomalley/monkey_lang/code"
	"github.com/lukeomalley/monkey_lang/object"
)

// Bytecode output of compiler used by the VM
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

// Bytecode constructs a new bytecode object
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

// EmittedInstruction used to store the last two instructions emitted by the compiler
type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

// Compiler structure used to store state of compiled code
type Compiler struct {
	instructions        code.Instructions
	constants           []object.Object
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
	symbolTable         *SymbolTable
}

// New cnstructs a new compiler
func New() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		symbolTable:         NewSymbolTable(),
	}
}

// NewWithState constructs a compiler object with state from a previous compiler
func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = s
	compiler.constants = constants
	return compiler
}

// Compile traverses the AST and emits bytecode
func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		// Define the name within the symbol table
		symbol := c.symbolTable.Define(node.Name.Value)

		// emit the op code for setting
		c.emit(code.OpSetGlobal, symbol.Index)

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable: %s", symbol.Name)
		}

		c.emit(code.OpGetGlobal, symbol.Index)

	case *ast.InfixExpression:
		// "Rewrite" code for less than to reduce instruction set
		if node.Operator == "<" {
			// Flip the order of the right and left nodes
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}

			c.emit(code.OpGreaterThan)
			return nil
		}

		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		case ">":
			c.emit(code.OpGreaterThan)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown prefix operator %s", node.Operator)
		}
	case *ast.IfExpression:
		/* Example w/ Bytecode:
		if (true) { 10 } else { 20 }; 3333;
			OpTrue
			OpJumpNotTruthy [pos]
			OpConstant        |
			OpJump [pos]      |
			OpConstant  <-----'
			OpPop
			OpConstant
			OpPop
		*/

		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// Create a conditional jump with a dummy location and store the position to be updated later
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999) // 9999 is a dummy value

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		// Remove the additional pop is present
		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}

		// Create a jump with a dummy location and store the position to be updated later
		jumpPos := c.emit(code.OpJump, 9999)

		// Update the conditional jump location to after the jump
		afterConsequencePos := len(c.instructions)
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			// Fill the empty space with a null return value
			c.emit(code.OpNull)
		} else {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			// Remove any additional pops
			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}
		}

		// Update the location of the jump to after the alternative
		afterAlternativePos := len(c.instructions)
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.IntegerLiteral:
		// Create an integer object
		integer := &object.Integer{Value: node.Value}

		// Append integer to the constants slice and emit the instruction
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			err := c.Compile(el)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpArray, len(node.Elements))

	case *ast.HashLiteral:
		keys := []ast.Expression{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			err := c.Compile(k)
			if err != nil {
				return err
			}

			err = c.Compile(node.Pairs[k])
			if err != nil {
				return err
			}
		}

		c.emit(code.OpHash, len(node.Pairs)*2)

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	}

	return nil
}

// =============================================================================
// Helper Methods
// =============================================================================

// emit constructs a bytecode object, adds it to the instructions slice and
// returns the position of the newly added instruction
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)
	c.setLastInstruction(op, pos)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.previousInstruction = previous
	c.lastInstruction = last
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
}

func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.instructions[opPos])
	newInstruction := code.Make(op, operand)
	c.replaceInstruction(opPos, newInstruction)
}
