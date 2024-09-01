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
