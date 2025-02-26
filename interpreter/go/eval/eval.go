package eval

import (
	"fmt"
	"kaze/ast"
	"kaze/object"
	"reflect"
)

var (
	TRUE     = &object.Boolean{Value: true}
	FALSE    = &object.Boolean{Value: false}
	NULL     = &object.Null{}
	NAN      = &object.NaN{}
	BREAK    = &object.Break{}
	CONTINUE = &object.Continue{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		evaluated := evalStatements(node.Statements, env)
		if returnValue, ok := evaluated.(*object.ReturnValue); ok {
			return returnValue.Value
		}
		return evaluated
	case *ast.VarStatement:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}

		env.Create(node.Name.Value, value)
	case *ast.ReturnStatement:
		if node.ReturnValue == nil {
			return &object.ReturnValue{Value: NULL}
		}
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.FunctionDefinitionStatement:
		fn := &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}
		env.Create(node.Name.Value, fn)
	case *ast.WhileStatement:
		return evalWhileStatement(node, env)
	case *ast.BreakStatement:
		return BREAK
	case *ast.ContinueStatement:
		return CONTINUE
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
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
	case *ast.AssignExpression:
		lvalue, err := evalLValue(node.Left, env)
		if err != nil {
			return err
		}

		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}

		if result, ok := lvalue.Update(value); ok {
			return result
		}
		return newError("assignment failed")
	case *ast.BlockExpression:
		return evalStatements(node.Statements, object.NewEnclosedEnvironment(env))
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
	case *ast.IfExpression:
		condition := Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			return Eval(node.Consequence, env)
		}

		if node.Alternative != nil {
			return Eval(node.Alternative, env)
		}

		return NULL
	case *ast.IndexExpression:
		array := Eval(node.Left, env)
		if isError(array) {
			return array
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(array, index)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.StringLiteral:
		return newString(node.Value)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)
	case *ast.NullLiteral:
		return NULL
	}

	return nil
}

func evalLValue(node ast.Expression, env *object.Environment) (object.LValue, *object.Error) {
	switch node := node.(type) {
	case *ast.Identifier:
		return &object.Variable{Name: node.Value, Env: env}, nil
	case *ast.IndexExpression:
		array, err := evalLValue(node.Left, env)
		if err != nil {
			return nil, err
		}
		index := Eval(node.Index, env)
		return &object.IndexRef{Left: array, Index: index}, nil
	}
	return nil, newError("not a lvalue: %s", node.String())
}

func evalArrayLiteral(node *ast.ArrayLiteral, env *object.Environment) object.Object {
	elements := evalExpressions(node.Elements, env)
	if len(elements) == 1 && isError(elements[0]) {
		return elements[0]
	}
	return &object.Array{Elements: elements}
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch {
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(left object.Object, index object.Object) object.Object {
	array := left.(*object.Array)
	idx := index.(*object.Integer).Value

	if idx < 0 || idx >= int64(len(array.Elements)) {
		return newError("index out of range: %d", idx)
	}

	return array.Elements[idx]
}

func evalHashIndexExpression(hash object.Object, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalStringIndexExpression(stringObj object.Object, indexObj object.Object) object.Object {
	str := stringObj.(*object.String).Value
	index := indexObj.(*object.Integer).Value

	if index < 0 || index >= int64(len(str)) {
		return newError("index out of range: %d", index)
	}

	return newString(string(str[index]))
}

func evalWhileStatement(node *ast.WhileStatement, env *object.Environment) object.Object {
	for isTruthy(Eval(node.Condition, env)) {
		result := Eval(node.Body, env)
		if isError(result) {
			return result
		}
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue
		}
		if result == BREAK {
			break
		}
		if result == CONTINUE {
			continue
		}
	}
	return NULL
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if val, ok := builtins[node.Value]; ok {
		return val
	}

	return newError("identifier not found: " + node.Value)
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)
		if isError(result) {
			return result
		}
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue
		}
		if result == BREAK || result == CONTINUE {
			return result
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
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

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(reflect.DeepEqual(left, right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!reflect.DeepEqual(left, right))
	case operator == "&&":
		return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))
	case operator == "||":
		return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return newString(leftVal + rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if left == NAN || right == NAN {
		return NAN
	}

	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expression := range expressions {
		expr := Eval(expression, env)
		if isError(expr) {
			return []object.Object{expr}
		}
		result = append(result, expr)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	if fn.Type() == object.BUILTIN_OBJ {
		return fn.(*object.Builtin).Fn(args...)
	}

	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(function *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(function.Env)

	for i, param := range function.Parameters {
		env.Create(param.Value, args[i])
	}

	return env
}

func unwrapReturnValue(evaluated object.Object) object.Object {
	if returnValue, ok := evaluated.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return evaluated
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newString(value string) *object.String {
	return &object.String{Value: value}
}
