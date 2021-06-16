package evaluator

import (
	"github.com/GenericEntity/interpreter-go/monkey/ast"
	"github.com/GenericEntity/interpreter-go/monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalStatements(node.Statements)

	case *ast.IfExpression:
		return evalIfExpression(node)
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, operand object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(operand)
	case "-":
		return evalMinusPrefixOperatorExpression(operand)
	default:
		return NULL
	}
}

func evalBangOperatorExpression(operand object.Object) *object.Boolean {
	switch operand {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(operand object.Object) object.Object {
	if operand.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := operand.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	// Uses pointer comparison to handle == and != for booleans
	//  works because there are two singleton boolean objects, TRUE and FALSE
	//  and requires all other overloads of == and != to be handled above these cases
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	lValue := left.(*object.Integer).Value
	rValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: lValue + rValue}
	case "-":
		return &object.Integer{Value: lValue - rValue}
	case "*":
		return &object.Integer{Value: lValue * rValue}
	case "/":
		return &object.Integer{Value: lValue / rValue}
	case ">":
		return nativeBoolToBooleanObject(lValue > rValue)
	case "<":
		return nativeBoolToBooleanObject(lValue < rValue)
	case "==":
		return nativeBoolToBooleanObject(lValue == rValue)
	case "!=":
		return nativeBoolToBooleanObject(lValue != rValue)
	default:
		return NULL
	}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	return obj != NULL && obj != FALSE
}
