// evaluator/evaluator.go
package evaluator

import (
	"github.com/simplyYan/Wysb/src/environment"
	"github.com/simplyYan/Wysb/src/parser"
)

func Eval(node parser.Node, env *environment.Environment) interface{} {
	switch node := node.(type) {
	case *parser.LetStatement:
		val := Eval(node.Value, env)
		env.Set(node.Name, val)
	case *parser.IntegerLiteral:
		return node.Value
	case *parser.InfixExpression:
		leftVal := Eval(node.Left, env)
		rightVal := Eval(node.Right, env)
		switch node.Operator {
		case "+":
			return leftVal.(int) + rightVal.(int)
		case "-":
			return leftVal.(int) - rightVal.(int)
		case "*":
			return leftVal.(int) * rightVal.(int)
		case "/":
			return leftVal.(int) / rightVal.(int)
		}
	case *parser.Identifier:
		val, _ := env.Get(node.Name)
		return val
	}

	return nil
}
