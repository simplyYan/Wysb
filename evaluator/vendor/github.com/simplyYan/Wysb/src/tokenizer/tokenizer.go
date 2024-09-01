// tokenizer/tokenizer.go
package tokenizer

import (
	"unicode"
)

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	IDENT
	INT
	ASSIGN
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	SEMICOLON
	LPAREN
	RPAREN
	LCURLY
	RCURLY
	IF
	ELSE
	LET
	COLON
	EQUAL
	PRINT
	STRING
)

type Token struct {
	Type    TokenType
	Literal string
}

func Tokenize(input string) []Token {
	var tokens []Token
	for i := 0; i < len(input); i++ {
		char := input[i]

		switch {
		case unicode.IsLetter(rune(char)):
			start := i
			for i < len(input) && unicode.IsLetter(rune(input[i])) {
				i++
			}
			literal := input[start:i]
			switch literal {
			case "let":
				tokens = append(tokens, Token{Type: LET, Literal: literal})
			case "if":
				tokens = append(tokens, Token{Type: IF, Literal: literal})
			case "else":
				tokens = append(tokens, Token{Type: ELSE, Literal: literal})
			case "print":
				tokens = append(tokens, Token{Type: PRINT, Literal: literal})
			default:
				tokens = append(tokens, Token{Type: IDENT, Literal: literal})
			}
			i--

		case unicode.IsDigit(rune(char)):
			start := i
			for i < len(input) && unicode.IsDigit(rune(input[i])) {
				i++
			}
			tokens = append(tokens, Token{Type: INT, Literal: input[start:i]})
			i--

		case char == '=':
			if i+1 < len(input) && input[i+1] == '=' {
				tokens = append(tokens, Token{Type: EQUAL, Literal: "=="})
				i++
			} else {
				tokens = append(tokens, Token{Type: ASSIGN, Literal: string(char)})
			}

		case char == '+':
			tokens = append(tokens, Token{Type: PLUS, Literal: string(char)})
		case char == '-':
			tokens = append(tokens, Token{Type: MINUS, Literal: string(char)})
		case char == '*':
			tokens = append(tokens, Token{Type: MULTIPLY, Literal: string(char)})
		case char == '/':
			tokens = append(tokens, Token{Type: DIVIDE, Literal: string(char)})
		case char == ';':
			tokens = append(tokens, Token{Type: SEMICOLON, Literal: string(char)})
		case char == '(':
			tokens = append(tokens, Token{Type: LPAREN, Literal: string(char)})
		case char == ')':
			tokens = append(tokens, Token{Type: RPAREN, Literal: string(char)})
		case char == '{':
			tokens = append(tokens, Token{Type: LCURLY, Literal: string(char)})
		case char == '}':
			tokens = append(tokens, Token{Type: RCURLY, Literal: string(char)})
		case char == ':':
			tokens = append(tokens, Token{Type: COLON, Literal: string(char)})
		case char == '"':
			start := i + 1
			i++
			for i < len(input) && input[i] != '"' {
				i++
			}
			tokens = append(tokens, Token{Type: STRING, Literal: input[start:i]})

		case unicode.IsSpace(rune(char)):
			continue

		default:
			tokens = append(tokens, Token{Type: ILLEGAL, Literal: string(char)})
		}
	}

	tokens = append(tokens, Token{Type: EOF, Literal: ""})
	return tokens
}
