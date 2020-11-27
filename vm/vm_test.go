package vm

import (
	"fmt"
	"testing"

	"github.com/lukeomalley/monkey_lang/ast"
	"github.com/lukeomalley/monkey_lang/compiler"
	"github.com/lukeomalley/monkey_lang/lexer"
	"github.com/lukeomalley/monkey_lang/object"
	"github.com/lukeomalley/monkey_lang/parser"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"2 - 1", 1},
		{"2 * 2", 4},
		{"6 / 3", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
	}

	runVMTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"false == false", true},
		{"1 > 2", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"1 < 2", true},
		{"true != true", false},
	}

	runVMTests(t, tests)
}

// =============================================================================
// Helper Functions
// =============================================================================

func runVMTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for i, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.LastPoppedStackElem()
		testExpectedObject(t, i, tt.expected, stackElem)
	}
}

func testExpectedObject(t *testing.T, testIndex int, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("[index - %d] testIntegerObject failed: %s", testIndex, err)
		}

	case bool:
		err := testBooleanObject(bool(expected), actual)
		if err != nil {
			t.Errorf("[index - %d] testBooleanObject failed: %s", testIndex, err)
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	// Check that the result is a int
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T, (%+v)", actual, actual)
	}

	// Check that the result is the correct value
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
	// Check if the actual value is a boolean
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not a boolean. got=%t, (%+v)", result.Value, expected)
	}

	// Check that the result is the correct value
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}

	return nil
}
