package github.com/simplyYan/Wysb/errorhandling

import (
	"fmt"
)

type Error struct {
	Stage   string 
	Message string 
	Details string 
}

func New(stage, message, details string) *Error {
	return &Error{
		Stage:   stage,
		Message: message,
		Details: details,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("Erro na %s: %s. Detalhes: %s", e.Stage, e.Message, e.Details)
}

func Wrap(err error, stage, details string) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		Stage:   stage,
		Message: err.Error(),
		Details: details,
	}
}

func ExampleUsage() {

	err1 := New("Análise Léxica", "Erro ao tokenizar o código", "Token inesperado encontrado")
	fmt.Println(err1.Error())

	originalErr := fmt.Errorf("erro original")
	err2 := Wrap(originalErr, "Análise Sintática", "Erro ao analisar a árvore de sintaxe")
	fmt.Println(err2.Error())
}
