package evaluator

import (
	"sunbird/ast"
	"sunbird/object"
)

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
		

	return newError("identifier not found: %s", node.Value)
}