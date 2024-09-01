package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/simplyYan/Wysb/src/environment"
	"github.com/simplyYan/Wysb/src/evaluator"
	"github.com/simplyYan/Wysb/src/parser" 
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wysb <file>")
		return
	}

	filename := os.Args[1]
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	tokens := parser.Tokenize(string(source))
	ast := parser.Parse(tokens)

	env := environment.NewEnvironment()
	for _, stmt := range ast {
		evaluator.Eval(stmt, env)
	}

	fmt.Println("Execution finished. Environment state:", env.Variables())
}
