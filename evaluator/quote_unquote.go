package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
	"monkey/token"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{Node: node}
}

func evalUnquoteCalls(quoted ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}

		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if len(call.Arguments) != 1 {
			return node
		}

        unquoted := Eval(call.Arguments[0], env)
        return convertObjectToASTNode(unquoted)
	})
}

func isUnquoteCall(node ast.Node) bool {
	CallExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}

	return CallExpression.Function.TokenLiteral() == "unquote"
}

func convertObjectToASTNode(obj object.Object) ast.Node {
    switch obj := obj.(type) {
    case *object.Integer:
        t := token.Token{
            Type:    token.INT,
            Literal: fmt.Sprintf("%d", obj.Value),
        }
        return &ast.IntegerLiteral{Token: t, Value: obj.Value}
    default:
        return nil
    }
}
