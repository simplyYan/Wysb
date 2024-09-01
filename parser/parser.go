package parser

import (
	"strings"
	"unicode"
)

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	IDENTIFIER
	INT
	ASSIGN
	PLUS
	MINUS
	ASTERISK
	SLASH
	LPAREN
	RPAREN
	SEMICOLON

)

type Token struct {
	Type    TokenType
	Literal string
}

func Tokenize(input string) []Token {
	var tokens []Token
	runes := []rune(input)

	for i := 0; i < len(runes); {
		ch := runes[i]

		if unicode.IsSpace(ch) {
			i++
			continue
		}

		switch ch {
		case '=':
			tokens = append(tokens, Token{Type: ASSIGN, Literal: string(ch)})
		case '+':
			tokens = append(tokens, Token{Type: PLUS, Literal: string(ch)})
		case '-':
			tokens = append(tokens, Token{Type: MINUS, Literal: string(ch)})
		case '*':
			tokens = append(tokens, Token{Type: ASTERISK, Literal: string(ch)})
		case '/':
			tokens = append(tokens, Token{Type: SLASH, Literal: string(ch)})
		case '(':
			tokens = append(tokens, Token{Type: LPAREN, Literal: string(ch)})
		case ')':
			tokens = append(tokens, Token{Type: RPAREN, Literal: string(ch)})
		case ';':
			tokens = append(tokens, Token{Type: SEMICOLON, Literal: string(ch)})
		default:
			if isLetter(ch) {
				identifier := readIdentifier(runes, &i)
				tokens = append(tokens, Token{Type: IDENTIFIER, Literal: identifier})
				continue
			} else if isDigit(ch) {
				number := readNumber(runes, &i)
				tokens = append(tokens, Token{Type: INT, Literal: number})
				continue
			} else {
				tokens = append(tokens, Token{Type: ILLEGAL, Literal: string(ch)})
			}
		}

		i++
	}

	tokens = append(tokens, Token{Type: EOF, Literal: ""})
	return tokens
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func readIdentifier(input []rune, start *int) string {
	var sb strings.Builder
	for *start < len(input) && isLetter(input[*start]) {
		sb.WriteRune(input[*start])
		*start++
	}
	return sb.String()
}

func readNumber(input []rune, start *int) string {
	var sb strings.Builder
	for *start < len(input) && isDigit(input[*start]) {
		sb.WriteRune(input[*start])
		*start++
	}
	return sb.String()
}

type NodeType int

const (
	VARIABLE NodeType = iota
	NUMBER
	EXPRESSION
	ASSIGNMENT
	STATEMENT
)

type Node struct {
	Type  NodeType
	Value string
	Left  *Node
	Right *Node
}

type Statement struct {
	Node *Node
}

func Parse(tokens []Token) []Statement {
	var statements []Statement
	var current *Node
	var stack []*Node
	var node *Node

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token.Type {
		case IDENTIFIER:
			node = &Node{Type: VARIABLE, Value: token.Literal}
			if current != nil {
				current.Left = node
			}
			stack = append(stack, node)
			current = nil
		case INT:
			node = &Node{Type: NUMBER, Value: token.Literal}
			if current != nil {
				current.Right = node
			}
			stack = append(stack, node)
			current = nil
		case ASSIGN:
			node = &Node{Type: ASSIGNMENT, Value: token.Literal}
			if len(stack) > 0 {
				node.Left = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, node)
		case PLUS, MINUS, ASTERISK, SLASH:
			node = &Node{Type: EXPRESSION, Value: token.Literal}
			if len(stack) > 0 {
				node.Left = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 {
				node.Right = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, node)
			current = node
		case SEMICOLON:
			if len(stack) > 0 {
				statements = append(statements, Statement{Node: stack[0]})
				stack = nil
			}
		case LPAREN, RPAREN:

		case EOF:
			if len(stack) > 0 {
				statements = append(statements, Statement{Node: stack[0]})
			}
		default:

		}
	}

	return statements
}
