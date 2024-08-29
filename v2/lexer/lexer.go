package lexer

import (
	"unicode"
	"unicode/utf8"
)

type TokenType string

const (
	EOF          TokenType = "EOF"
	NUMBER        TokenType = "NUMBER"
	IDENTIFIER    TokenType = "IDENTIFIER"
	PLUS          TokenType = "PLUS"
	MINUS         TokenType = "MINUS"
	MULTIPLY      TokenType = "MULTIPLY"
	DIVIDE        TokenType = "DIVIDE"
	LPAREN        TokenType = "LPAREN"
	RPAREN        TokenType = "RPAREN"
	ASSIGN        TokenType = "ASSIGN"
	SEMICOLON     TokenType = "SEMICOLON"
	WHITESPACE    TokenType = "WHITESPACE"
	UNKNOWN       TokenType = "UNKNOWN"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal string
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch, _ = utf8.DecodeRuneInString(l.input[l.readPosition:])
	}
	l.position = l.readPosition
	l.readPosition += utf8.RuneLen(l.ch)
}

func (l *Lexer) NextToken() Token {
	var tok Token

	switch l.ch {
	case '+':
		tok = newToken(PLUS, l.ch)
	case '-':
		tok = newToken(MINUS, l.ch)
	case '*':
		tok = newToken(MULTIPLY, l.ch)
	case '/':
		tok = newToken(DIVIDE, l.ch)
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case '=':
		tok = newToken(ASSIGN, l.ch)
	case ';':
		tok = newToken(SEMICOLON, l.ch)
	case 0:
		tok.Lexeme = ""
		tok.Type = EOF
	default:
		if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Lexeme = l.readNumber()
			tok.Literal = tok.Lexeme
			return tok
		} else if isLetter(l.ch) {
			tok.Type = IDENTIFIER
			tok.Lexeme = l.readIdentifier()
			tok.Literal = tok.Lexeme
			return tok
		} else {
			tok = newToken(UNKNOWN, l.ch)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType TokenType, ch rune) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  string(ch),
		Literal: string(ch),
	}
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}