package standardlibrary

import (
	"fmt"
	"math"
	"strings"
)

var FuncMap = map[string]func([]interface{}) (interface{}, error){
	"print":    print,
	"sqrt":     sqrt,
	"concat":   concat,
	"length":   length,
	"toupper":  toupper,
	"tolower":  tolower,
	"max":      max,
	"min":      min,
}

func print(args []interface{}) (interface{}, error) {
	for _, arg := range args {
		fmt.Print(arg)
	}
	fmt.Println()
	return nil, nil
}

func sqrt(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("a função sqrt espera 1 argumento")
	}
	arg, ok := args[0].(float64)
	if !ok {
		return nil, fmt.Errorf("o argumento deve ser um número")
	}
	return math.Sqrt(arg), nil
}

func concat(args []interface{}) (interface{}, error) {
	var result strings.Builder
	for _, arg := range args {
		str, ok := arg.(string)
		if !ok {
			return nil, fmt.Errorf("todos os argumentos devem ser strings")
		}
		result.WriteString(str)
	}
	return result.String(), nil
}

func length(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("a função length espera 1 argumento")
	}
	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("o argumento deve ser uma string")
	}
	return len(str), nil
}

func toupper(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("a função toupper espera 1 argumento")
	}
	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("o argumento deve ser uma string")
	}
	return strings.ToUpper(str), nil
}

func tolower(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("a função tolower espera 1 argumento")
	}
	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("o argumento deve ser uma string")
	}
	return strings.ToLower(str), nil
}

func max(args []interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("a função max espera pelo menos 1 argumento")
	}
	var max float64
	for i, arg := range args {
		num, ok := arg.(float64)
		if !ok {
			return nil, fmt.Errorf("todos os argumentos devem ser números")
		}
		if i == 0 || num > max {
			max = num
		}
	}
	return max, nil
}

func min(args []interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("a função min espera pelo menos 1 argumento")
	}
	var min float64
	for i, arg := range args {
		num, ok := arg.(float64)
		if !ok {
			return nil, fmt.Errorf("todos os argumentos devem ser números")
		}
		if i == 0 || num < min {
			min = num
		}
	}
	return min, nil
}