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
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
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
		{"!(if (false) { 5; })", true},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
	}

	runVMTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{input: "if (true) { 10 }", expected: 10},
		{input: "if (true) { 10 } else { 20 }", expected: 10},
		{input: "if (false) { 10 } else { 20 }", expected: 20},
		{input: "if (1) { 10 }", expected: 10},
		{input: "if (1 < 2) { 10 }", expected: 10},
		{input: "if (1 < 2) { 10 } else { 20 }", expected: 10},
		{input: "if (1 > 2) { 10 } else { 20 }", expected: 20},
		{input: "if (false) { 10 }", expected: Null},
		{input: "if (1 > 2) { 10 }", expected: Null},
	}

	runVMTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{input: "let one = 1; one", expected: 1},
		{input: "let one = 1; let two = 2; one + two", expected: 3},
		{input: "let one = 1; let two = one + one; one + two", expected: 3},
	}

	runVMTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{input: `"monkey"`, expected: "monkey"},
		{input: `"mon" + "key"`, expected: "monkey"},
		{input: `"mon" + "key" + "banana"`, expected: "monkeybanana"},
	}

	runVMTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{input: "[]", expected: []int{}},
		{input: "[1, 2, 3]", expected: []int{1, 2, 3}},
		{input: "[1 + 2, 3 * 4, 5 + 6]", expected: []int{3, 12, 11}},
	}

	runVMTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    "{}",
			expected: map[object.HashKey]int64{},
		},
		{
			input: "{1: 2, 2: 3}",
			expected: map[object.HashKey]int64{
				(&object.Integer{Value: 1}).HashKey(): 2,
				(&object.Integer{Value: 2}).HashKey(): 3,
			},
		},
		{
			input: "{1 + 1: 2 * 2, 3 + 3: 4 * 4}",
			expected: map[object.HashKey]int64{
				(&object.Integer{Value: 2}).HashKey(): 4,
				(&object.Integer{Value: 6}).HashKey(): 16,
			},
		},
	}

	runVMTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{input: "[1, 2, 3][1]", expected: 2},
		{input: "[1, 2, 3][0 + 2]", expected: 3},
		{input: "[[1, 1, 1]][0][0]", expected: 1},
		{input: "[][0]", expected: Null},
		{input: "[1, 2, 9][99]", expected: Null},
		{input: "{1: 1, 2: 2}[1]", expected: 1},
		{input: "{1: 1, 2: 2}[2]", expected: 2},
		{input: "{1: 1}[0]", expected: Null},
		{input: "{}[0]", expected: Null},
	}

	runVMTests(t, tests)
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			  let fivePlusTen = fn() { 5 + 10 };
				fivePlusTen();
			`,
			expected: 15,
		},
		{
			input: `
				let one = fn() { 1; };
				let two = fn() { 2; };
				one() + two()
			`,
			expected: 3,
		},
		{
			input: `
				let a = fn() { 1 };
				let b = fn() { a() + 1 };
				let c = fn() { b() + 1 };
				c();
			`,
			expected: 3,
		},
	}

	runVMTests(t, tests)
}

func TestFunctionsWithReturnStatement(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			  let earlyExit = fn() { return 99; 100; };
				earlyExit();
			`,
			expected: 99,
		},
		{
			input: `
			  let earlyExit = fn() { return 99; return 100; };
				earlyExit();
			`,
			expected: 99,
		},
	}

	runVMTests(t, tests)
}

func TestFunctionsWithoutReturnValue(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			  let noReturn = fn() { };
				noReturn();
			`,
			expected: Null,
		},
		{
			input: `
			  let noReturnOne = fn() { };
			  let noReturnTwo = fn() { noReturnOne(); };
				noReturnOne();
				noReturnTwo();
			`,
			expected: Null,
		},
	}

	runVMTests(t, tests)
}

func TestFirstClassFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			  let returnsOne = fn() { 1; };
			  let returnsOneReturner = fn() { returnsOne; };
				returnsOneReturner()();
			`,
			expected: 1,
		},
		{
			input: `
			  let returnsOneReturner = fn() { 
			    let returnsOne = fn() { 1; };
					returnsOne;
				};

				returnsOneReturner()();
			`,
			expected: 1,
		},
	}

	runVMTests(t, tests)
}

func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			  let one = fn() { let one = 1; one };
				one();
			`,
			expected: 1,
		},
		{
			input: `
			  let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
				oneAndTwo();
			`,
			expected: 3,
		},
		{
			input: `
			  let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
			  let threeAndFour = fn() { let three = 3; let four = 4; three + four; };
				oneAndTwo() + threeAndFour();
			`,
			expected: 10,
		},
		{
			input: `
			  let firstFoobar = fn() { let foobar = 50; foobar; };
			  let secondFoobar = fn() { let foobar = 100; foobar; };
				firstFoobar() + secondFoobar();
			`,
			expected: 150,
		},
		{
			input: `
			  let globalSeed = 50;
				let minusOne = fn() {
					let num = 1;
					globalSeed - num;
				};


				let minusTwo = fn() {
					let num = 2;
					globalSeed - num;
				};

				minusOne() + minusTwo();
			`,
			expected: 97,
		},
	}

	runVMTests(t, tests)
}

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			  let identity = fn(a) { a; };
				identity(4);
			`,
			expected: 4,
		},
		{
			input: `
			  let sum = fn(a, b) { a + b; };
				sum(1, 2);
			`,
			expected: 3,
		},
		{
			input: `
			  let sum = fn(a, b) { 
					let c = a + b; 
					c;
			  };
				sum(1, 2);
			`,
			expected: 3,
		},
		{
			input: `
			  let sum = fn(a, b) { 
					let c = a + b; 
					c;
			  };

				let outer = fn() {
				  sum(1, 2) + sum(3, 4);
				};

				outer();
			`,
			expected: 10,
		},
		{
			input: `
			  let globalNum = 10;

				let sum = fn(a, b) {
					let c = a + b;
					c + globalNum;
				};

				let outer = fn() {
					return sum(1, 2) + sum(3, 4) + globalNum;
				};

				outer() + globalNum;
			`,
			expected: 50,
		},
	}

	runVMTests(t, tests)
}

func TestCallingFunctionsWithWrongArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `fn() { 1; }(1);`,
			expected: `wrong number of arguments: want=0, got=1`,
		},
		{
			input:    `fn(a) { a; }();`,
			expected: `wrong number of arguments: want=1, got=0`,
		},
		{
			input:    `fn(a, b) { a + b; }(1);`,
			expected: `wrong number of arguments: want=2, got=1`,
		},
	}

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err == nil {
			t.Fatalf("expected VM error, but resulted in none.")
		}

		if err.Error() != tt.expected {
			t.Fatalf("wrong VM error. want=%q, but got=%q", tt.expected, err)
		}
	}
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

	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}

	case []int:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("object not Array, %T (%+v)", actual, actual)
			return
		}

		if len(array.Elements) != len(expected) {
			t.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(array.Elements))
			return
		}

		for i, expectedElem := range expected {
			err := testIntegerObject(int64(expectedElem), array.Elements[i])
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}

	case map[object.HashKey]int64:
		hash, ok := actual.(*object.Hash)
		if !ok {
			t.Errorf("object is not Hash. got=%T (%+v)", actual, actual)
			return
		}

		if len(hash.Pairs) != len(expected) {
			t.Errorf("hash has wrong number of Pairs. want=%d, got=%d", len(expected), len(hash.Pairs))
		}

		for expectedKey, expectedValue := range expected {
			pair, ok := hash.Pairs[expectedKey]
			if !ok {
				t.Errorf("no pair for given key in Pairs")
			}

			err := testIntegerObject(expectedValue, pair.Value)
			if err != nil {
				t.Errorf("testIntegerObject failes: %s", err)
			}
		}

	case *object.Null:
		if actual != Null {
			t.Errorf("object is not Null: %T (%+v)", actual, actual)
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

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%t+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
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
