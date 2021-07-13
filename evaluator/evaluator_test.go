package evaluator

import (
	"errors"
	"testing"

	"github.com/GenericEntity/interpreter-go/monkey/lexer"
	"github.com/GenericEntity/interpreter-go/monkey/object"
	"github.com/GenericEntity/interpreter-go/monkey/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	lex := lexer.New(input)
	p := parser.New(lex)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. expected=%d, got=%d", expected, result.Value)
		return false
	}

	return true
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. expected=%t, got=%t", expected, result.Value)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		expectedInt, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(expectedInt))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 10; 9;", 10},
		{`
		if (10 > 1) {
			if (10 > 1) {
				return 10;
			}

			return 1;
		}`, 10},
		{
			`let add = fn(x, y) {x + y};
			fn(){
				let x = add(5, 2);
				let y = add(-2, 5);
				let z = add(1000, 2);
				return x + y;
			}()`,
			10,
		},
		{
			`
			let f = fn(x) {
			let result = x + 10;
			return result;
			return 10;
			};
			f(10);`,
			20,
		},
		{
			`
			let f = fn(x) {
			return x;
			x + 10;
			};
			f(10);`,
			10,
		},
		{
			// ensure that the return statement in foo does not
			// cause the outer function call to return
			`
			let foo = fn(){ return 10; }
			fn(){
				foo()
				5;
			}()`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if ( 10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar;",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not *object.Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "{(x + 2)}"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
		{"fn(x){ fn(y){x + y} }(2)(3)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not string. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. expected=%q. got=%q", "Hello World!", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// len(String)
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, errors.New("type of 1st argument to `len` not supported, got INTEGER")},
		{`len("one", "two")`, errors.New("wrong number of arguments. got=2, want=1")},

		// len(Array)
		{`len([])`, 0},
		{`len([3])`, 1},
		{`len(["string", true])`, 2},
		{`let arr = [0, 0, 0]; len(arr)`, 3},
		{`let arr = fn(){ [1,2] }(); len(arr)`, 2},

		// first(Array)
		{`first([3])`, 3},
		{`first(["string", true])`, "string"},
		{`let arr = [0, 1, 2]; first(arr)`, 0},
		{`let arr = fn(){ [1,2] }(); first(arr)`, 1},
		{`first("asd")`, errors.New("type of 1st argument to `first` not supported, got STRING")},
		{`first(5)`, errors.New("type of 1st argument to `first` not supported, got INTEGER")},
		{`first(true)`, errors.New("type of 1st argument to `first` not supported, got BOOLEAN")},
		{`first([])`, errors.New("`first` should not be called on empty array")},
		{`first(["one"], ["two"])`, errors.New("wrong number of arguments. got=2, want=1")},

		// last(Array)
		{`last([3])`, 3},
		{`last(["string", true])`, true},
		{`let arr = [0, 1, 2]; last(arr)`, 2},
		{`let arr = fn(){ [1,2] }(); last(arr)`, 2},
		{`last("asd")`, errors.New("type of 1st argument to `last` not supported, got STRING")},
		{`last(5)`, errors.New("type of 1st argument to `last` not supported, got INTEGER")},
		{`last(true)`, errors.New("type of 1st argument to `last` not supported, got BOOLEAN")},
		{`last([])`, errors.New("`last` should not be called on empty array")},
		{`last(["one"], ["two"])`, errors.New("wrong number of arguments. got=2, want=1")},

		// rest(Array)
		{`rest([3])`, []interface{}{}},
		{`rest(["string", true])`, []interface{}{true}},
		{`let arr = [0, 1, 2]; rest(arr)`, []interface{}{1, 2}},
		{`let arr = fn(){ [1,2] }(); rest(arr)`, []interface{}{2}},
		{`rest("asd")`, errors.New("type of 1st argument to `rest` not supported, got STRING")},
		{`rest(5)`, errors.New("type of 1st argument to `rest` not supported, got INTEGER")},
		{`rest(true)`, errors.New("type of 1st argument to `rest` not supported, got BOOLEAN")},
		{`rest([])`, errors.New("`rest` should not be called on empty array")},
		{`rest(["one"], ["two"])`, errors.New("wrong number of arguments. got=2, want=1")},

		// push(Array, elem)
		{`push([3], 2)`, []interface{}{3, 2}},
		{`push(["string", true], "hi")`, []interface{}{"string", true, "hi"}},
		{`let arr = [0, 1, 2]; push(arr, 3)`, []interface{}{0, 1, 2, 3}},
		{`let arr = fn(){ [1,2] }(); push(arr, false)`, []interface{}{1, 2, false}},
		{`push([], 1)`, []interface{}{1}},
		{`push(["one"], ["two"])`, []interface{}{"one", []interface{}{"two"}}},
		{`push("asd", 2)`, errors.New("type of 1st argument to `push` not supported, got STRING")},
		{`push(5, 1)`, errors.New("type of 1st argument to `push` not supported, got INTEGER")},
		{`push(true, 8)`, errors.New("type of 1st argument to `push` not supported, got BOOLEAN")},
		{`push(["one"], "two", "three")`, errors.New("wrong number of arguments. got=3, want=2")},
		{`push(["one"])`, errors.New("wrong number of arguments. got=1, want=2")},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testObject(t, evaluated, tt.expected)
	}
}

func TestArrayLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{"[1,2, 3]", []interface{}{1, 2, 3}},
		{"[true, true, false]", []interface{}{true, true, false}},
		{"[1, true]", []interface{}{1, true}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		arr, ok := evaluated.(*object.Array)
		if !ok {
			t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
		}

		if len(arr.Elements) != len(tt.expected) {
			t.Fatalf("array has wrong number of elements. expected=%d, got=%d", len(tt.expected), len(arr.Elements))
		}

		for i, elem := range tt.expected {
			switch elem := elem.(type) {
			case int:
				testIntegerObject(t, arr.Elements[i], int64(elem))
			case bool:
				testBooleanObject(t, arr.Elements[i], elem)
			default:
				t.Fatalf("unsupported type in expected: %T", elem)
			}
		}
	}
}

func testObject(t *testing.T, obj object.Object, expected interface{}) bool {
	if expected == nil {
		return testNullObject(t, obj)
	}

	switch expected := expected.(type) {
	case int:
		return testIntegerObject(t, obj, int64(expected))
	case int64:
		return testIntegerObject(t, obj, expected)
	case bool:
		return testBooleanObject(t, obj, expected)
	case string:
		return testStringObject(t, obj, expected)
	case error:
		return testErrorObject(t, obj, expected.Error())
	case []interface{}:
		return testArray(t, obj, expected)
	case map[interface{}]interface{}:
		return testHash(t, obj, expected)

	default:
		t.Fatalf("unsupported type in expected: %T", expected)
		return false
	}
}

func testArray(t *testing.T, obj object.Object, expected []interface{}) bool {
	arr, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("object is not Array. got=%T (%+v)", obj, obj)
		return false
	}

	if len(arr.Elements) != len(expected) {
		t.Errorf("Array has wrong length. expected=%d. got=%d", len(expected), len(arr.Elements))
		return false
	}

	for i, val := range expected {
		if !testObject(t, arr.Elements[i], val) {
			t.Errorf("Array has wrong value at index %d. expected=%v. got=%v", i, val, arr.Elements[i])
			return false
		}
	}

	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	str, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not string. got=%T (%+v)", obj, obj)
		return false
	}

	if str.Value != expected {
		t.Errorf("String has wrong value. expected=%q. got=%q", expected, str.Value)
		return false
	}

	return true
}

func testErrorObject(t *testing.T, obj object.Object, expectedMsg string) bool {
	errObj, ok := obj.(*object.Error)
	if !ok {
		t.Errorf("object is not Error. got=%T (%+v)", obj, obj)
		return false
	}

	if errObj.Message != expectedMsg {
		t.Errorf("wrong error message. expected=%q, got=%q", expectedMsg, errObj.Message)
		return false
	}

	return true
}

func TestArraySubscriptExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1,2, 3][1];", 2},
		{"[true, false][0]", true},
		{`[1, true][fn(x,y){x+y}(-1, 9-8)]`, 1},
		{`[1, true][5*2 - 9]`, true},
		{`["hi", "there"][5*2 - 9]`, "there"},
		{`[fn(x){x + 1}, fn(x){x + 2}][1](5)`, 7},
		{
			`let arr = [3, 2, 1, 0]
			arr[0]`,
			3,
		},
		{
			`let arr = [3, 2, 1, 0]
			arr[0] + arr[2] + arr[1] + arr[3]`,
			6,
		},

		{"[1, true][2]", errors.New("index out of range: 2. array length: 2")},
		{"[1, true][-1]", errors.New("index out of range: -1. array length: 2")},
		{`[][true]`, errors.New("non-integer argument to array subscript not supported, got BOOLEAN")},
		{`[][fn(x){x}]`, errors.New("non-integer argument to array subscript not supported, got FUNCTION")},
		{`[]["hi"]`, errors.New("non-integer argument to array subscript not supported, got STRING")},
		{`fn(x){x}["hi"]`, errors.New("subscript operator not supported for type: FUNCTION")},
		{`"I'm not an array"[0]`, errors.New("subscript operator not supported for type: STRING")},
		{`1[0]`, errors.New("subscript operator not supported for type: INTEGER")},

		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{"let myArray = [1, 2, 3]; myArray[2];", 3},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]", 2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testObject(t, evaluated, tt.expected)
	}
}

func testHash(t *testing.T, obj object.Object, expected map[interface{}]interface{}) bool {
	hash, ok := obj.(*object.Hash)
	if !ok {
		t.Errorf("object is not Hash. got=%T (%+v)", obj, obj)
		return false
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash has wrong number of pairs. expected=%d, got=%d", len(expected), len(hash.Pairs))
		return false
	}

	for _, pair := range hash.Pairs {
		var expectedVal interface{}
		var ok bool

		switch key := pair.Key.(type) {
		case *object.Integer:
			expectedVal, ok = expected[int(key.Value)]
			if !ok {
				expectedVal, ok = expected[key.Value]
			}

		case *object.Boolean:
			expectedVal, ok = expected[key.Value]

		case *object.String:
			expectedVal, ok = expected[key.Value]

		default:
			t.Errorf("unsupported key type to test: %T (%+v)", key, key)
			return false
		}

		if !ok || !testObject(t, pair.Value, expectedVal) {
			t.Errorf("hash.Pairs contains unexpected pair: <%+v, %+v>.", pair.Key, pair.Value)
			return false
		}
	}

	return true
}

func TestHashLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{}`, map[interface{}]interface{}{}},
		{`{1: 2, 2: 4, 3: 6};`, map[interface{}]interface{}{1: 2, 2: 4, 3: 6}},
		{`{true: 1, 2: false};`, map[interface{}]interface{}{true: 1, 2: false}},
		{`{"hi": 2, "there": 1}`, map[interface{}]interface{}{"hi": 2, "there": 1}},
		{`let x = "x"; let y = "not y"; {x: 2, y: 1}`, map[interface{}]interface{}{"x": 2, "not y": 1}},
		{`fn(){ {2 * 2: 2, 5 == 10: 1} }();`, map[interface{}]interface{}{4: 2, false: 1}},

		{`{[1]: 2}`, errors.New("invalid key type: ARRAY")},
		{`{{1: 2}: 2}`, errors.New("invalid key type: HASH")},
		{`{fn(){}: 2}`, errors.New("invalid key type: FUNCTION")},
		{`{1: 2, 1: 3}`, errors.New("duplicate key: 1")},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testObject(t, evaluated, tt.expected)
	}
}

func TestHashSubscriptExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{1: 2, 2: 4, 3: 6}[1];`, 2},
		{`{true: 1, 2: false}[true];`, 1},
		{`{"hi": 2, "there": 1}["there"]`, 1},
		{`let x = "x"; let y = "not y"; {x: 2, y: 1}["not y"]`, 1},
		{`fn(){ {2 * 2: 2, 5 == 10: 1} }()[8/2];`, 2},

		{`{}[1]`, nil},
		{`{"hi": 2, "there": 1}["asd"]`, nil},

		{`{"hi": 2, "there": 1}[[1,2]]`, errors.New("invalid key type: ARRAY")},
		{`{"hi": 2, "there": 1}[{1: 2}]`, errors.New("invalid key type: HASH")},
		{`{"hi": 2, "there": 1}[fn(){}]`, errors.New("invalid key type: FUNCTION")},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testObject(t, evaluated, tt.expected)
	}
}
