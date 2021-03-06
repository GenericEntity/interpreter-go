package evaluator

import (
	"fmt"

	"github.com/GenericEntity/interpreter-go/monkey/ast"
	"github.com/GenericEntity/interpreter-go/monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{
			Parameters: params,
			Body:       body,
			Env:        env,
		}

	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(fn, args)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)

	case *ast.SubscriptExpression:
		return evalSubscriptExpression(node, env)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	}

	return nil
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch function := fn.(type) {
	case *object.Function:
		if len(function.Parameters) != len(args) {
			return newError("wrong number of arguments to function. expected=%d, got=%d",
				len(function.Parameters),
				len(args))
		}

		fnCallEnv := extendFunctionEnv(function, args)
		evaluated := Eval(function.Body, fnCallEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		// no need to unwrap return value because we never return an object.ReturnValue
		return function.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(function *object.Function, args []object.Object) *object.Environment {
	env := object.ExtendEnvironment(function.Env)
	for i, param := range function.Parameters {
		env.Set(param.Value, args[i])
	}
	return env
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		// Eval could return nil. e.g. for LET statements
		if result == nil {
			continue
		}

		rt := result.Type()

		// Don't unwrap the return value yet (in case we're in a nested statement)
		// This statement stops evaluation of later statements
		if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
			return result
		}
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
		return newError("unknown operator: %s%s", operator, operand.Type())
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
		return newError("unknown operator: -%s", operand.Type())
	}

	value := operand.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)

	// Uses pointer comparison to handle == and != for booleans
	//  works because there are two singleton boolean objects, TRUE and FALSE
	//  and requires all other overloads of == and != to be handled above these cases
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
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
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	return obj != NULL && obj != FALSE
}

func newError(format string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}

func isError(obj object.Object) bool {
	if obj == nil {
		return false
	}

	return obj.Type() == object.ERROR_OBJ
}

func evalIdentifier(id *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(id.Value); ok {
		return val
	}

	if builtin, ok := builtins[id.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", id.Value)
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	lValue := left.(*object.String).Value
	rValue := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: lValue + rValue}
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalArrayLiteral(arr *ast.ArrayLiteral, env *object.Environment) object.Object {
	exprs := evalExpressions(arr.Elements, env)
	if len(exprs) == 1 && isError(exprs[0]) {
		return exprs[0]
	}

	return &object.Array{Elements: exprs}
}

func evalSubscriptExpression(subscriptExpr *ast.SubscriptExpression, env *object.Environment) object.Object {
	left := Eval(subscriptExpr.Left, env)
	if isError(left) {
		return left
	}

	indexObj := Eval(subscriptExpr.Index, env)
	if isError(indexObj) {
		return indexObj
	}

	switch left := left.(type) {
	case *object.Array:
		return evalArraySubscriptExpression(left, indexObj)

	case *object.Hash:
		return evalHashSubscriptExpression(left, indexObj)

	default:
		return newError("subscript operator not supported for type: %s", left.Type())
	}
}

func evalPairExpressions(exprs map[ast.Expression]ast.Expression, env *object.Environment) (map[object.HashKey]object.HashPair, object.Object) {
	result := make(map[object.HashKey]object.HashPair, len(exprs))
	for key, val := range exprs {
		evaluatedKey := Eval(key, env)
		if isError(evaluatedKey) {
			return nil, evaluatedKey
		}

		hashable, ok := evaluatedKey.(object.Hashable)
		if !ok {
			return nil, newError("invalid key type: %s", evaluatedKey.Type())
		}

		evaluatedVal := Eval(val, env)
		if isError(evaluatedVal) {
			return nil, evaluatedVal
		}

		_, duplicateKey := result[hashable.HashKey()]
		if duplicateKey {
			return nil, newError("duplicate key: %s", evaluatedKey.Inspect())
		}

		result[hashable.HashKey()] = object.HashPair{Key: evaluatedKey, Value: evaluatedVal}
	}

	return result, nil
}

func evalHashLiteral(hash *ast.HashLiteral, env *object.Environment) object.Object {
	pairs, err := evalPairExpressions(hash.Pairs, env)
	if err != nil {
		return err
	}
	return &object.Hash{Pairs: pairs}
}

func evalHashSubscriptExpression(hash *object.Hash, index object.Object) object.Object {
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("invalid key type: %s", index.Type())
	}
	pair, ok := hash.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

func evalArraySubscriptExpression(array *object.Array, index object.Object) object.Object {
	idx, ok := index.(*object.Integer)
	if !ok {
		return newError("non-integer argument to array subscript not supported, got %s", index.Type())
	}
	if idx.Value < 0 || idx.Value >= int64(len(array.Elements)) {
		return newError("index out of range: %d. array length: %d", idx.Value, len(array.Elements))
	}
	return array.Elements[idx.Value]
}
