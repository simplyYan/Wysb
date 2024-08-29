package codegen

import (
	"fmt"
	"github.com/simplyYan/Wysb/parser"
)

type Instrucao struct {
	Opcode string
	Arg    string
}

type CodeGenerator struct {
	ast      *parser.Node
	intermediario []Instrucao
}

func New(ast *parser.Node) *CodeGenerator {
	return &CodeGenerator{ast: ast}
}

func (cg *CodeGenerator) Generate() ([]Instrucao, error) {
	cg.intermediario = []Instrucao{}
	if err := cg.gerarInstrucoes(cg.ast); err != nil {
		return nil, err
	}
	return cg.intermediario, nil
}

func (cg *CodeGenerator) gerarInstrucoes(node *parser.Node) error {
	switch node.Type {
	case parser.ProgramNode:
		for _, child := range node.Children {
			if err := cg.gerarInstrucoes(child); err != nil {
				return err
			}
		}
	case parser.ExpressionNode:
		if len(node.Children) == 0 {
			return fmt.Errorf("expressão vazia")
		}
		left := node.Children[0]
		if err := cg.gerarInstrucoes(left); err != nil {
			return err
		}
		if len(node.Children) > 1 {
			op := node.Children[1]
			if err := cg.gerarInstrucoes(op); err != nil {
				return err
			}
			if len(node.Children) > 2 {
				right := node.Children[2]
				if err := cg.gerarInstrucoes(right); err != nil {
					return err
				}
				cg.intermediario = append(cg.intermediario, Instrucao{
					Opcode: "ADD",
					Arg:    left.Value + " " + right.Value,
				})
			}
		}
	case parser.NumberNode:
		cg.intermediario = append(cg.intermediario, Instrucao{
			Opcode: "PUSH",
			Arg:    node.Value,
		})
	case parser.IdentifierNode:
		cg.intermediario = append(cg.intermediario, Instrucao{
			Opcode: "LOAD",
			Arg:    node.Value,
		})
	case parser.BinaryOpNode:
		if len(node.Children) != 2 {
			return fmt.Errorf("número incorreto de filhos para BinaryOpNode")
		}
		left := node.Children[0]
		right := node.Children[1]
		if err := cg.gerarInstrucoes(left); err != nil {
			return err
		}
		if err := cg.gerarInstrucoes(right); err != nil {
			return err
		}
		cg.intermediario = append(cg.intermediario, Instrucao{
			Opcode: "BINOP",
			Arg:    left.Value + " " + right.Value,
		})
	default:
		return fmt.Errorf("tipo de nó desconhecido: %s", node.Type)
	}
	return nil
}
