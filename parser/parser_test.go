package parser

import (
	"fmt"
	"testing"

	"github.com/GenericEntity/interpreter-go/monkey/ast"
	"github.com/GenericEntity/interpreter-go/monkey/lexer"
)

type id struct {
	name string
}

func TestLetStatements(t *testing.T) {
	// proper testing would mock the lexer
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", id{"y"}},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statement(s). got=%d",
				1, len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", stmt)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got='%s'", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got='%s'", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return a;", id{"a"}},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statement(s). got=%d",
				1, len(program.Statements))
		}

		returnStmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", program.Statements[0])
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral is not 'return'. got=%q",
				returnStmt.TokenLiteral())
		}

		if !testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	program := testParse(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testIdentifier(t, stmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	program := testParse(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	testIntegerLiteral(t, stmt.Expression, 5)
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input         string
		operator      string
		expectedValue interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		program := testParse(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %q. got=%q", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.expectedValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, expectedValue int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != expectedValue {
		t.Errorf("integ.Value not %d. got=%d.", expectedValue, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", expectedValue) {
		t.Errorf("integ.TokenLiteral() not %d. got=%s", expectedValue, integ.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, expr ast.Expression, value string) bool {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Errorf("expr is not *ast.Identifier. got=%T", expr)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value is not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() is not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case int64:
		return testIntegerLiteral(t, expr, v)
	case string:
		return testStringLiteral(t, expr, v)
	case id:
		return testIdentifier(t, expr, v.name)
	case bool:
		return testBooleanLiteral(t, expr, v)
	}
	t.Errorf("type of expr not handled. got=%T", expr)
	return false
}

func testInfixExpression(t *testing.T, expr ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExpr, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expr is not *ast.InfixExpression. got=%T(%s)", expr, expr)
		return false
	}

	if !testLiteralExpression(t, opExpr.Left, left) {
		return false
	}

	if opExpr.Operator != operator {
		t.Errorf("expr.Operator is not %s. got=%s", operator, opExpr.Operator)
		return false
	}

	if !testLiteralExpression(t, opExpr.Right, right) {
		return false
	}

	return true
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 6", 5, "+", 6},
		{"5 - 6", 5, "-", 6},
		{"5 * 6", 5, "*", 6},
		{"5 / 6", 5, "/", 6},
		{"5 > 6", 5, ">", 6},
		{"5 < 6", 5, "<", 6},
		{"5 == 6", 5, "==", 6},
		{"5 != 6", 5, "!=", 6},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		program := testParse(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not %d statement(s). got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			continue
		}
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"(1 + (2 - 3) * (5 - 2) + (((1))))/2",
			"(((1 + ((2 - 3) * (5 - 2))) + 1) / 2)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"[fn(x){x + 1}, fn(x){x + 2}][1](5)",
			"([fn(x){(x + 1)}, fn(x){(x + 2)}][1])(5)",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testBooleanLiteral(t *testing.T, expr ast.Expression, value bool) bool {
	b, ok := expr.(*ast.Boolean)
	if !ok {
		t.Errorf("expr is not *ast.Boolean. got=%T", expr)
		return false
	}

	if b.Value != value {
		t.Errorf("b.Value is not %t. got=%t", value, b.Value)
		return false
	}

	if b.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("b.TokenLiteral() is not %t. got=%s", value, b.TokenLiteral())
		return false
	}
	return true
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program does not have %d statement(s). got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		testBooleanLiteral(t, stmt.Expression, tt.expected)
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	program := testParse(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have %d statement(s). got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, expr.Condition, id{"x"}, "<", id{"y"}) {
		return
	}

	if len(expr.Consequence.Statements) != 1 {
		t.Errorf("expr has more than %d statement(s). got=%d", 1, len(expr.Consequence.Statements))
	}

	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T",
			expr.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expr.Alternative != nil {
		t.Errorf("expr.Alternative is not nil. got=%+v", expr.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	program := testParse(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have %d statement(s). got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, expr.Condition, id{"x"}, "<", id{"y"}) {
		return
	}

	if len(expr.Consequence.Statements) != 1 {
		t.Errorf("expr does not have %d statement(s). got=%d", 1, len(expr.Consequence.Statements))
	}

	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T",
			expr.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expr.Alternative == nil {
		t.Fatalf("expr.Alternative is nil")
	}

	if len(expr.Alternative.Statements) != 1 {
		t.Errorf("expr.Alternative does not have %d statement(s). got%d", 1, len(expr.Alternative.Statements))
	}

	alternative, ok := expr.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative.Statements[0] is not *ast.ExpressionStatement. got=%T",
			expr.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	program := testParse(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have %d statement(s). got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function does not have %d parameter(s). got=%d", 2, len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], id{"x"})
	testLiteralExpression(t, function.Parameters[1], id{"y"})

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function body does not have %d statement(s). got=%d", 1, len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not *ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, id{"x"}, "+", id{"y"})
}

func TestFunctionLiteralParameters(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []id
	}{
		{input: "fn(){};", expectedParams: []id{}},
		{input: "fn(x){};", expectedParams: []id{{"x"}}},
		{input: "fn(x, y, z){};", expectedParams: []id{{"x"}, {"y"}, {"z"}}},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("wrong number of parameters. expected=%d. got=%d",
				len(tt.expectedParams), len(function.Parameters))
		}

		for i, id := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], id)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`

	program := testParse(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have %d statement(s). got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	callExpr, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, callExpr.Function, "add") {
		return
	}

	if len(callExpr.Arguments) != 3 {
		t.Fatalf("callExpr has wrong number of arguments. expected=%d, got=%d",
			3, len(callExpr.Arguments))
	}

	testLiteralExpression(t, callExpr.Arguments[0], 1)
	testInfixExpression(t, callExpr.Arguments[1], 2, "*", 3)
	testInfixExpression(t, callExpr.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameters(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []interface{}
	}{
		{input: "nullary();", expectedParams: []interface{}{}},
		{input: "unary(x);", expectedParams: []interface{}{id{"x"}}},
		{input: "binary(a, b);", expectedParams: []interface{}{id{"a"}, id{"b"}}},
		{input: "ternary(x, y, z);", expectedParams: []interface{}{id{"x"}, id{"y"}, id{"z"}}},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		callExpr := stmt.Expression.(*ast.CallExpression)

		if len(callExpr.Arguments) != len(tt.expectedParams) {
			t.Errorf("wrong number of parameters. expected=%d. got=%d",
				len(tt.expectedParams), len(callExpr.Arguments))
		}

		for i, id := range tt.expectedParams {
			testLiteralExpression(t, callExpr.Arguments[i], id)
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	program := testParse(t, input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestArrayLiteralExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedElems []interface{}
	}{
		{input: `[1,2, 3];`, expectedElems: []interface{}{1, 2, 3}},
		{input: `[true, false];`, expectedElems: []interface{}{true, false}},
		{input: `[1, true]`, expectedElems: []interface{}{1, true}},
		{input: `[]`, expectedElems: []interface{}{}},
		{input: `[1]`, expectedElems: []interface{}{1}},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		arr, ok := stmt.Expression.(*ast.ArrayLiteral)
		if !ok {
			t.Fatalf("expr not *ast.ArrayLiteral. got=%T", stmt.Expression)
		}

		if len(arr.Elements) != len(tt.expectedElems) {
			t.Fatalf("literal.Elements is of wrong length. got=%d, want=%d", len(arr.Elements), len(tt.expectedElems))
		}

		for i, elem := range tt.expectedElems {
			testLiteralExpression(t, arr.Elements[i], elem)
		}
	}
}

func testArrayLiteral(t *testing.T, arr ast.Expression, expectedElems []interface{}) bool {
	array, ok := arr.(*ast.ArrayLiteral)
	if !ok {
		t.Errorf("arr not *ast.ArrayLiteral. got=%T", arr)
		return false
	}

	if len(array.Elements) != len(expectedElems) {
		t.Errorf("array.Elements has wrong length. expected=%d, got=%d", len(expectedElems), len(array.Elements))
		return false
	}

	for i, e := range expectedElems {
		testLiteralExpression(t, array.Elements[i], e)
	}

	// empty array
	return true
}

func TestArraySubscriptExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedElems []interface{}
		expectedIndex int
	}{
		{input: `[1,2, 3][1];`, expectedElems: []interface{}{1, 2, 3}, expectedIndex: 1},
		{input: `[true, false][0];`, expectedElems: []interface{}{true, false}, expectedIndex: 0},
		{input: `[1, true][59]`, expectedElems: []interface{}{1, true}, expectedIndex: 59},
		{input: `[][10]`, expectedElems: []interface{}{}, expectedIndex: 10},
		{input: `[1][1]`, expectedElems: []interface{}{1}, expectedIndex: 1},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		subscriptExpr, ok := stmt.Expression.(*ast.SubscriptExpression)
		if !ok {
			t.Fatalf("expr not *ast.SubscriptExpression. got=%T", stmt.Expression)
		}

		testArrayLiteral(t, subscriptExpr.Left, tt.expectedElems)

		testIntegerLiteral(t, subscriptExpr.Index, int64(tt.expectedIndex))
	}
}

func testParse(t *testing.T, input string) *ast.Program {
	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	return program
}

func TestHashLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected map[interface{}]interface{}
	}{
		{`{1: 2, 2: 4, 3: 6};`, map[interface{}]interface{}{1: 2, 2: 4, 3: 6}},
		{`{true: 1, 2: false};`, map[interface{}]interface{}{true: 1, 2: false}},
		{`{x: 2, y: 1}`, map[interface{}]interface{}{id{"x"}: 2, id{"y"}: 1}},
		{`{"hi": 2, "there": 1}`, map[interface{}]interface{}{"hi": 2, "there": 1}},
		{`{}`, map[interface{}]interface{}{}},
	}

	for _, tt := range tests {
		program := testParse(t, tt.input)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		hash, ok := stmt.Expression.(*ast.HashLiteral)
		if !ok {
			t.Fatalf("expr not *ast.HashLiteral. got=%T", stmt.Expression)
		}

		if len(hash.Pairs) != len(tt.expected) {
			t.Fatalf("hash.Pairs is of wrong length. got=%d, want=%d", len(hash.Pairs), len(tt.expected))
		}

		for key, val := range hash.Pairs {
			var expectedVal interface{}
			var ok bool

			switch key := key.(type) {
			case *ast.IntegerLiteral:
				expectedVal, ok = tt.expected[int(key.Value)]

			case *ast.Boolean:
				expectedVal, ok = tt.expected[key.Value]

			case *ast.StringLiteral:
				expectedVal, ok = tt.expected[key.Value]

			case *ast.Identifier:
				expectedVal, ok = tt.expected[id{key.Value}]

			default:
				t.Fatalf("unsupported key type to test: %T (%+v)", key, key)
			}

			if !ok || !testLiteralExpression(t, val, expectedVal) {
				t.Fatalf("hash.Pairs contains unexpected pair: <%+v, %+v>.", key, val)
			}
		}
	}
}

func testStringLiteral(t *testing.T, expr ast.Expression, expectedValue string) bool {
	str, ok := expr.(*ast.StringLiteral)
	if !ok {
		t.Errorf("expr not *ast.StringLiteral. got=%T", expr)
		return false
	}

	if str.Value != expectedValue {
		t.Errorf("str.Value not %s. got=%s.", expectedValue, str.Value)
		return false
	}

	return true
}
