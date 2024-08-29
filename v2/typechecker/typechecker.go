package typechecker

import (
	"fmt"
	"../parser"
)

type Tipo string

const (
	NumberTipo     Tipo = "NUMBER"
	IdentifierTipo Tipo = "IDENTIFIER"
)

type TipoChecker struct {
	ast *parser.Node
}

func New(ast *parser.Node) *TipoChecker {
	return &TipoChecker{ast: ast}
}

func (tc *TipoChecker) Check() error {
	return tc.checkNode(tc.ast)
}

func (tc *TipoChecker) checkNode(node *parser.Node) error {
	switch node.Type {
	case parser.ProgramNode:
		for _, child := range node.Children {
			if err := tc.checkNode(child); err != nil {
				return err
			}
		}
	case parser.ExpressionNode:
		if len(node.Children) == 0 {
			return fmt.Errorf("expressão vazia")
		}
		left := node.Children[0]
		if err := tc.checkNode(left); err != nil {
			return err
		}
		if len(node.Children) > 1 {
			op := node.Children[1]
			if err := tc.checkNode(op); err != nil {
				return err
			}
			if len(node.Children) > 2 {
				right := node.Children[2]
				if err := tc.checkNode(right); err != nil {
					return err
				}
				if left.Type != right.Type {
					return fmt.Errorf("tipos incompatíveis: %s e %s", left.Type, right.Type)
				}
			}
		}
	case parser.NumberNode:
		node.Type = NumberTipo
	case parser.IdentifierNode:
		node.Type = IdentifierTipo
	case parser.BinaryOpNode:

		if len(node.Children) != 2 {
			return fmt.Errorf("número incorreto de filhos para BinaryOpNode")
		}
		left := node.Children[0]
		right := node.Children[1]
		if err := tc.checkNode(left); err != nil {
			return err
		}
		if err := tc.checkNode(right); err != nil {
			return err
		}
		if left.Type != right.Type {
			return fmt.Errorf("tipos incompatíveis para operação binária: %s e %s", left.Type, right.Type)
		}
	default:
		return fmt.Errorf("tipo de nó desconhecido: %s", node.Type)
	}
	return nil
}