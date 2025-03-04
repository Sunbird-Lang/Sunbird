package evaluator

import (
	"sunbird/ast"
	"sunbird/object"
)

var (
	NULL  = &object.Null{} 
	TRUE 	= &object.Boolean{Value: true}
  FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}	

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	
	case *ast.BlockStatement:
		return evalStatements(node.Statements)

	case *ast.IfExpression:
		return evalIfExpression(node)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: val}
	}
	
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	
	for _, statement := range stmts {
		result = Eval(statement)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
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

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperator(right)
	
	case "-":
		return evalMinusPrefixOperator(right)
	default:
		return NULL
	}
}

func evalBangOperator(right object.Object) object.Object {
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

func evalMinusPrefixOperator(right object.Object) object.Object {
	if right.Type() == object.INTEGER_OBJ {
		value := right.(*object.Integer).Value
		return &object.Integer{ Value: -value }
	}

	if right.Type() == object.FLOAT_OBJ {
		value := right.(*object.Float).Value
		return &object.Float{ Value: -value }
	}

	return NULL
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	
	// TODO: handle string here

	case left.Type() == object.FLOAT_OBJ || right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right) 

	case operator == "==":
		return nativeBoolToBooleanObject(left == right)

	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	default:
		return NULL
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
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
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return NULL
}
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	var leftVal, rightVal float64

	if left.Type() == object.INTEGER_OBJ {
			leftVal = float64(left.(*object.Integer).Value)
	} else {
			leftVal = left.(*object.Float).Value
	}

	if right.Type() == object.INTEGER_OBJ {
			rightVal = float64(right.(*object.Integer).Value)
	} else {
			rightVal = right.(*object.Float).Value
	}

	switch operator {
	case "+":
			return &object.Float{Value: leftVal + rightVal}
	case "-":
			return &object.Float{Value: leftVal - rightVal}
	case "*":
			return &object.Float{Value: leftVal * rightVal}
	case "/":
			return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)	
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