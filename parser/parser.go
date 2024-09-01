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

type IfStatement struct {
	Condition Expression
	Consequence []Node
	Alternative []Node
}

func (is *IfStatement) String() string {
	return fmt.Sprintf("if %s { ... } else { ... }", is.Condition.String())
}

type ForStatement struct {
	Identifier string
	Range      Expression
	Body       []Node
}

func (fs *ForStatement) String() string {
	return fmt.Sprintf("for %s in %s { ... }", fs.Identifier, fs.Range.String())
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

		case tokenizer.IF:
			i++
			condition := parseExpression(tokens, &i)
			i++ // Skip '{'
			consequence := parseBlock(tokens, &i)
			var alternative []Node
			if tokens[i].Type == tokenizer.ELSE {
				i++
				i++ // Skip '{'
				alternative = parseBlock(tokens, &i)
			}
			statements = append(statements, &IfStatement{Condition: condition, Consequence: consequence, Alternative: alternative})

		case tokenizer.FOR:
			i++
			identifier := tokens[i].Literal
			i++ // Skip 'in'
			rangeExpr := parseExpression(tokens, &i)
			i++ // Skip '{'
			body := parseBlock(tokens, &i)
			statements = append(statements, &ForStatement{Identifier: identifier, Range: rangeExpr, Body: body})
		}
		i++
	}

	return statements
}

func parseBlock(tokens []tokenizer.Token, i *int) []Node {
	var block []Node
	for tokens[*i].Type != tokenizer.RCURLY {
		block = append(block, Parse(tokens)[*i])
		*i++
	}
	return block
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
