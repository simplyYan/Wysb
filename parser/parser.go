// parser/parser.go
package parser

import (
	"fmt"
	"strconv"

	"github.com/simplyYan/Wysb/src/tokenizer"
)

type Node interface {
	String() string
}

type LetStatement struct {
	Name  string
	Type  string
	Value Expression
}

func (ls *LetStatement) String() string {
	return fmt.Sprintf("let %s: %s = %s", ls.Name, ls.Type, ls.Value.String())
}

type Expression interface {
	Node
}

type IntegerLiteral struct {
	Value int
}

func (il *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", il.Value)
}

type Identifier struct {
	Name string
}

func (id *Identifier) String() string {
	return id.Name
}

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}

func Parse(tokens []tokenizer.Token) []Node {
	var statements []Node
	i := 0

	for i < len(tokens) {
		switch tokens[i].Type {
		case tokenizer.LET:
			i++
			if tokens[i].Type != tokenizer.IDENT {
				panic("expected identifier after 'let'")
			}
			name := tokens[i].Literal
			i++
			if tokens[i].Type != tokenizer.COLON {
				panic("expected ':' after identifier")
			}
			i++
			if tokens[i].Type != tokenizer.IDENT {
				panic("expected type after ':'")
			}
			varType := tokens[i].Literal
			i++
			if tokens[i].Type != tokenizer.ASSIGN {
				panic("expected '=' after type")
			}
			i++
			value := parseExpression(tokens, &i)
			statements = append(statements, &LetStatement{Name: name, Type: varType, Value: value})
		}
		i++
	}

	return statements
}

func parseExpression(tokens []tokenizer.Token, i *int) Expression {
	token := tokens[*i]

	switch token.Type {
	case tokenizer.INT:
		return &IntegerLiteral{Value: atoi(token.Literal)}
	case tokenizer.IDENT:
		left := &Identifier{Name: token.Literal}
		*i++
		operator := tokens[*i].Literal
		*i++
		right := parseExpression(tokens, i)
		return &InfixExpression{Left: left, Operator: operator, Right: right}
	}
	return nil
}

func atoi(str string) int {
	num, _ := strconv.Atoi(str)
	return num
}
