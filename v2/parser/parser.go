package parser

import (
	"fmt"
	"../lexer"
)

type NodeType string

const (
	ProgramNode      NodeType = "PROGRAM"
	ExpressionNode   NodeType = "EXPRESSION"
	NumberNode       NodeType = "NUMBER"
	IdentifierNode   NodeType = "IDENTIFIER"
	BinaryOpNode     NodeType = "BINARY_OP"
)

type Node struct {
	Type     NodeType
	Value    string
	Children []*Node
}

type Parser struct {
	lexer      *lexer.Lexer
	currentToken lexer.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.nextToken() 
	return p
}

func (p *Parser) Parse() *Node {
	return p.parseProgram()
}

func (p *Parser) parseProgram() *Node {
	node := &Node{Type: ProgramNode}

	for p.currentToken.Type != lexer.EOF {
		node.Children = append(node.Children, p.parseExpression())
	}

	return node
}

func (p *Parser) parseExpression() *Node {
	leftNode := p.parseTerm()

	for p.currentToken.Type == lexer.PLUS || p.currentToken.Type == lexer.MINUS {
		opNode := &Node{
			Type:  BinaryOpNode,
			Value: p.currentToken.Lexeme,
		}
		p.nextToken() 

		rightNode := p.parseTerm()
		opNode.Children = []*Node{leftNode, rightNode}

		leftNode = opNode
	}

	return leftNode
}

func (p *Parser) parseTerm() *Node {
	var node *Node

	switch p.currentToken.Type {
	case lexer.NUMBER:
		node = &Node{
			Type:  NumberNode,
			Value: p.currentToken.Lexeme,
		}
		p.nextToken()
	case lexer.IDENTIFIER:
		node = &Node{
			Type:  IdentifierNode,
			Value: p.currentToken.Lexeme,
		}
		p.nextToken()
	case lexer.LPAREN:
		p.nextToken() 
		node = p.parseExpression()
		if p.currentToken.Type != lexer.RPAREN {
			panic("Expected closing parenthesis")
		}
		p.nextToken() 
	default:
		panic(fmt.Sprintf("Unexpected token: %s", p.currentToken.Lexeme))
	}

	return node
}

func (p *Parser) nextToken() {
	p.currentToken = p.lexer.NextToken()
}