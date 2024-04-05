package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/bluele/gcache"
	lagra "github.com/simplyYan/LAGRA"
	"github.com/simplyYan/cutinfo"
	"strconv"
	"math"
)

var gc = gcache.New(20).
	LRU().
	Build()

var logger, err = lagra.New(`
	log_file = "output.log"
	`)

func countKeywords(text string, keywords []string) int {
	total := 0

	for _, keyword := range keywords {
		total += strings.Count(text, keyword)
	}

	return total
}

func ReadWysb(filename string) {
	file, err := os.Open(filename) // For read access.
	if err != nil {
		log.Fatal(err)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// Defina as palavras-chave que você deseja contar
	keywords := []string{"$[int32]"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts := countKeywords(string(content), keywords)
	// Defina as palavras-chave que você deseja contar
	keywords_float := []string{"$[float32]"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_float := countKeywords(string(content), keywords_float)

	ci := cutinfo.New()

	keywords_sum := []string{"<math.sum"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_sum := countKeywords(string(content), keywords_sum)
	for i := 1; i <= counts_sum; i++ {
		var_sum := ci.Target(string(content), "<math.sum", ">")
		var_num1 := ci.Target(var_sum, " ", " +")
		var_num2 := ci.Target(var_sum, "+ ", "!")
		num1, err := strconv.Atoi(var_num1)
		if err != nil {
			fmt.Println("Erro ao converter string para inteiro:", err)
			return
		}
		num2, err := strconv.Atoi(var_num2)
		if err != nil {
			fmt.Println("Erro ao converter string para inteiro:", err)
			return
		}
	
		replacer := "<math.sum " + var_num1 + " + " + var_num2 + "!>"
		rsult := num1 + num2
		rsult_str := strconv.Itoa(rsult)
		content = []byte(strings.Replace(string(content), replacer, rsult_str, i))

	}

	keywords_sub := []string{"<math.sub"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_sub := countKeywords(string(content), keywords_sub)
	for i := 1; i <= counts_sub; i++ {
		var_sum := ci.Target(string(content), "<math.sub", ">")
		var_num1 := ci.Target(var_sum, " ", " -")
		var_num2 := ci.Target(var_sum, "- ", "!")
		num1, err := strconv.Atoi(var_num1)
		if err != nil {
			fmt.Println("Erro ao converter string para inteiro:", err)
			return
		}
		num2, err := strconv.Atoi(var_num2)
		if err != nil {
			fmt.Println("Erro ao converter string para inteiro:", err)
			return
		}
	
		replacer := "<math.sub " + var_num1 + " - " + var_num2 + "!>"
		rsult := num1 - num2
		rsult_str := strconv.Itoa(rsult)
		content = []byte(strings.Replace(string(content), replacer, rsult_str, i))

	}

	keywords_div := []string{"<math.div"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_div := countKeywords(string(content), keywords_div)
	for i := 1; i <= counts_div; i++ {
		var_sum := ci.Target(string(content), "<math.div", ">")
		var_num1 := ci.Target(var_sum, " ", " /")
		var_num2 := ci.Target(var_sum, "/ ", "!")
		num1, err := strconv.Atoi(var_num1)
		if err != nil {
			fmt.Println("Erro ao converter string para inteiro:", err)
			return
		}
		num2, err := strconv.Atoi(var_num2)
		if err != nil {
			fmt.Println("Erro ao converter string para inteiro:", err)
			return
		}
	
		replacer := "<math.div " + var_num1 + " / " + var_num2 + "!>"
		rsult := num1 / num2
		rsult_str := strconv.Itoa(rsult)
		content = []byte(strings.Replace(string(content), replacer, rsult_str, i))

	}

	keywords_mult := []string{"<math.mult"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_mult := countKeywords(string(content), keywords_mult)
	for i := 1; i <= counts_mult; i++ {
		var_sum := ci.Target(string(content), "<math.mult", ">")
		var_num1 := ci.Target(var_sum, " ", " *")
		var_num2 := ci.Target(var_sum, "* ", "!")
		num1, err := strconv.Atoi(var_num1)
		if err != nil {
			fmt.Println("Erro ao converter string para inteiro:", err)
			return
		}
		num2, err := strconv.Atoi(var_num2)
		if err != nil {
			fmt.Println("Erro ao converter string para inteiro:", err)
			return
		}
	
		replacer := "<math.mult " + var_num1 + " * " + var_num2 + "!>"
		rsult := num1 * num2
		rsult_str := strconv.Itoa(rsult)
		content = []byte(strings.Replace(string(content), replacer, rsult_str, i))

	}
	
	keywords_pi := []string{"<math.pi"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_pi := countKeywords(string(content), keywords_pi)
	for i := 1; i <= counts_pi; i++ {
		replacer := "<math.pi!>"
		content = []byte(strings.Replace(string(content), replacer, "3.14159265358979323846", i))

	}

	keywords_pow := []string{"<math.pow"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_pow := countKeywords(string(content), keywords_pow)
	for i := 1; i <= counts_pow; i++ {
		var_sum := ci.Target(string(content), "<math.pow", ">")
		var_num1 := ci.Target(var_sum, " ", " ::")
		var_num2 := ci.Target(var_sum, ":: ", "!")
		num1, err := strconv.Atoi(var_num1)
		if err != nil {
			fmt.Println("Erro ao converter string para inteiro:", err)
			return
		}
		num2, err := strconv.Atoi(var_num2)
		if err != nil {
			fmt.Println("Erro ao converter string para inteiro:", err)
			return
		}
	
		replacer := "<math.pow " + var_num1 + " :: " + var_num2 + "!>"
		rsult := math.Pow(float64(num1), float64(num2))
		rsult_str := strconv.FormatFloat(rsult, 'f', -1, 64)
		content = []byte(strings.Replace(string(content), replacer, rsult_str, i))

	}

	keywords_sqrt := []string{"<math.sqrt"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_sqrt := countKeywords(string(content), keywords_sqrt)
	for i := 1; i <= counts_sqrt; i++ {
		var_sum := ci.Target(string(content), "<math.sqrt", ">")
		var_num1 := ci.Target(var_sum, " ", "!")
		num1, err := strconv.Atoi(var_num1)
		if err != nil {
			fmt.Println("Erro ao converter string para inteiro:", err)
			return
		}
	
		replacer := "<math.sqrt " + var_num1 + "!>"
		rsult := math.Sqrt(float64(num1))
		rsult_str := strconv.FormatFloat(rsult, 'f', -1, 64)
		content = []byte(strings.Replace(string(content), replacer, rsult_str, i))

	}

	keywords_openFile := []string{"<os.ReadFile"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_openFile := countKeywords(string(content), keywords_openFile)
	for i := 1; i <= counts_openFile; i++ {
		fileopen := ci.Target(string(content), "<os.ReadFile", ">")
		target := ci.Target(fileopen, " ", "!")
		thisfile, err := os.Open(target) 
		if err != nil {
			log.Fatal(err)
		}
	
		thiscontent, err := io.ReadAll(thisfile)
		if err != nil {
			panic(err)
		}


		replacer := "<os.ReadFile " + target + "!>"

		content = []byte(strings.Replace(string(content), replacer, string(thiscontent), i))

	}
	//Variables

	for i := 1; i <= counts; i++ {
		var_i32 := ci.Target(string(content), "$[int32]", ";")
		var_i32_name := ci.Target(var_i32, " ", " =")
		var_i32_value := ci.Target(var_i32, "= ", ">")
		gc.Set(var_i32_name, var_i32_value)
		fmt.Println("A chave ", var_i32_name, " com o valor ", var_i32_value, " foi registrada.")
		replacer := "$[int32] " + var_i32_name + " = " + var_i32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	for i := 1; i <= counts_float; i++ {
		var_f32 := ci.Target(string(content), "$[float32]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		fmt.Println("A chave ", var_f32_name, " com o valor ", var_f32_value, " foi registrada.")
		replacer := "$[float32] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}
	keywords_float64 := []string{"$[float64]"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_float64 := countKeywords(string(content), keywords_float64)
	for i := 1; i <= counts_float64; i++ {
		var_f32 := ci.Target(string(content), "$[float64]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		fmt.Println("A chave ", var_f32_name, " com o valor ", var_f32_value, " foi registrada.")
		replacer := "$[float64] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}
	keywords_float128 := []string{"$[float128]"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_float128 := countKeywords(string(content), keywords_float128)
	for i := 1; i <= counts_float128; i++ {
		var_f32 := ci.Target(string(content), "$[float128]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		fmt.Println("A chave ", var_f32_name, " com o valor ", var_f32_value, " foi registrada.")
		replacer := "$[float128] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_int64 := []string{"$[int64]"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_int64 := countKeywords(string(content), keywords_int64)
	for i := 1; i <= counts_int64; i++ {
		var_f32 := ci.Target(string(content), "$[int64]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		fmt.Println("A chave ", var_f32_name, " com o valor ", var_f32_value, " foi registrada.")
		replacer := "$[int64] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}
	keywords_int128 := []string{"$[int128]"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_int128 := countKeywords(string(content), keywords_int128)
	for i := 1; i <= counts_int128; i++ {
		var_f32 := ci.Target(string(content), "$[int128]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		fmt.Println("A chave ", var_f32_name, " com o valor ", var_f32_value, " foi registrada.")
		replacer := "$[int128] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}
	keywords_string := []string{"$[string]"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_string := countKeywords(string(content), keywords_string)
	for i := 1; i <= counts_string; i++ {
		var_f32 := ci.Target(string(content), "$[string]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		fmt.Println("A chave ", var_f32_name, " com o valor ", var_f32_value, " foi registrada.")
		replacer := "$[string] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_bool := []string{"$[bool]"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_bool := countKeywords(string(content), keywords_bool)
	for i := 1; i <= counts_bool; i++ {
		var_f32 := ci.Target(string(content), "$[bool]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		fmt.Println("A chave ", var_f32_name, " com o valor ", var_f32_value, " foi registrada.")
		replacer := "$[bool] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_array := []string{"$[array]"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_array := countKeywords(string(content), keywords_array)
	for i := 1; i <= counts_array; i++ {
		var_f32 := ci.Target(string(content), "$[bool]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		fmt.Println("A chave ", var_f32_name, " com o valor ", var_f32_value, " foi registrada.")
		replacer := "$[array] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_getvar := []string{"${"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_getvar := countKeywords(string(content), keywords_getvar)
	for i := 1; i <= counts_getvar; i++ {
		var_f32 := ci.Target(string(content), "${", "}")
		findVar, err := gc.Get(var_f32)
		if err != nil {
			panic(err)
		}

		// Convertendo findVar para uma string usando type assertion
		findVarStr, ok := findVar.(string)
		if !ok {
			panic("findVar não é uma string")
		}

		replacer := "${" + var_f32 + "}"
		content = []byte(strings.Replace(string(content), replacer, findVarStr, i))
	}

	

	keywords_fun := []string{"!fun"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_fun := countKeywords(string(content), keywords_fun)
	for i := 1; i == counts_fun; i++ {
		var_fun := ci.Target(string(content), "!fun", "{")
		var_fun_name := ci.Target(var_fun, " ", "(")
		var_fun_value := ci.Target(string(content), "{", "}")
		gc.Set(var_fun_name, var_fun_value)
		fmt.Println("A chave ", var_fun_name, " com o valor ", var_fun_value, " foi registrada.")
		replacer := "!fun " + var_fun_name + " = " + var_fun_value
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}
	keywords_println := []string{"println("}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_println := countKeywords(string(content), keywords_println)
	for i := 1; i == counts_println; i++ {
		var_fun := ci.Target(string(content), "println(", ")")
		logger.Info(context.Background(), var_fun)
		fmt.Println("O valor ", var_fun, " foi registrado.")
		replacer := "println(" + var_fun + ")"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_callfn := []string{"@:"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_callfn := countKeywords(string(content), keywords_callfn)
	for i := 1; i == counts_callfn; i++ {
		var_Callfun := ci.Target(string(content), "@:", "(")
		callFunc, err := gc.Get(var_Callfun)
		if err != nil {
			panic(err)
		}
		findVarStr, ok := callFunc.(string)
		if !ok {
			panic("findVar não é uma string")
		}
		execFunc(findVarStr)

		replacer := "@:" + var_Callfun + "("
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_if := []string{"$if[>]"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_if := countKeywords(string(content), keywords_if)
	for i := 1; i == counts_if; i++ {
		var_if := ci.Target(string(content), "$if[>]", ";")
		var_target := ci.Target(var_if, " ", " ::")
		var_compare := ci.Target(var_if, ":: ", " !")
		var_toExec := ci.Target(var_if, "! ", "(")
		target, err := gc.Get(var_target)
		if err != nil {
			panic(err)
		}

		targetStr, ok := target.(string)
		if !ok {
			panic("target não é uma string")
		}

		compare, err := gc.Get(var_compare)
		if err != nil {
			panic(err)
		}
		compareStr, ok := compare.(string)
		if !ok {
			panic("compare não é uma string")
		}

		toExec, err := gc.Get(var_toExec)
		if err != nil {
			panic(err)
		}
		toExecStr, ok := toExec.(string)
		if !ok {
			panic("compare não é uma string")
		}
		if targetStr > compareStr {
			execFunc(toExecStr)
		}
		replacer := "$if[>] " + targetStr + " :: " + compareStr + " ! " + toExecStr + "();"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_ifLess := []string{"$if[<]"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_ifLess := countKeywords(string(content), keywords_ifLess)
	for i := 1; i == counts_ifLess; i++ {
		var_if := ci.Target(string(content), "$if[<]", ";")
		var_target := ci.Target(var_if, " ", " ::")
		var_compare := ci.Target(var_if, ":: ", " !")
		var_toExec := ci.Target(var_if, "! ", "(")
		target, err := gc.Get(var_target)
		if err != nil {
			panic(err)
		}

		targetStr, ok := target.(string)
		if !ok {
			panic("target não é uma string")
		}

		compare, err := gc.Get(var_compare)
		if err != nil {
			panic(err)
		}
		compareStr, ok := compare.(string)
		if !ok {
			panic("compare não é uma string")
		}

		toExec, err := gc.Get(var_toExec)
		if err != nil {
			panic(err)
		}
		toExecStr, ok := toExec.(string)
		if !ok {
			panic("compare não é uma string")
		}
		if targetStr < compareStr {
			execFunc(toExecStr)
		}
		replacer := "$if[<] " + targetStr + " :: " + compareStr + " ! " + toExecStr + "();"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_ifEqual := []string{"$if[==]"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_ifEqual := countKeywords(string(content), keywords_ifEqual)
	for i := 1; i == counts_ifEqual; i++ {
		var_if := ci.Target(string(content), "$if[==]", ";")
		var_target := ci.Target(var_if, " ", " ::")
		var_compare := ci.Target(var_if, ":: ", " !")
		var_toExec := ci.Target(var_if, "! ", "(")
		target, err := gc.Get(var_target)
		if err != nil {
			panic(err)
		}

		targetStr, ok := target.(string)
		if !ok {
			panic("target não é uma string")
		}

		compare, err := gc.Get(var_compare)
		if err != nil {
			panic(err)
		}
		compareStr, ok := compare.(string)
		if !ok {
			panic("compare não é uma string")
		}

		toExec, err := gc.Get(var_toExec)
		if err != nil {
			panic(err)
		}
		toExecStr, ok := toExec.(string)
		if !ok {
			panic("compare não é uma string")
		}
		if targetStr == compareStr {
			execFunc(toExecStr)
		}
		replacer := "$if[==] " + targetStr + " :: " + compareStr + " ! " + toExecStr + "();"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_ifLessEqual := []string{"$if[<=]"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_ifLessEqual := countKeywords(string(content), keywords_ifLessEqual)
	for i := 1; i == counts_ifLessEqual; i++ {
		var_if := ci.Target(string(content), "$if[<=]", ";")
		var_target := ci.Target(var_if, " ", " ::")
		var_compare := ci.Target(var_if, ":: ", " !")
		var_toExec := ci.Target(var_if, "! ", "(")
		target, err := gc.Get(var_target)
		if err != nil {
			panic(err)
		}

		targetStr, ok := target.(string)
		if !ok {
			panic("target não é uma string")
		}

		compare, err := gc.Get(var_compare)
		if err != nil {
			panic(err)
		}
		compareStr, ok := compare.(string)
		if !ok {
			panic("compare não é uma string")
		}

		toExec, err := gc.Get(var_toExec)
		if err != nil {
			panic(err)
		}
		toExecStr, ok := toExec.(string)
		if !ok {
			panic("compare não é uma string")
		}
		if targetStr <= compareStr {
			execFunc(toExecStr)
		}
		replacer := "$if[<=] " + targetStr + " :: " + compareStr + " ! " + toExecStr + "();"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_ifBiggerEqual := []string{"$if[>=]"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_ifBiggerEqual := countKeywords(string(content), keywords_ifBiggerEqual)
	for i := 1; i == counts_ifBiggerEqual; i++ {
		var_if := ci.Target(string(content), "$if[>=]", ";")
		var_target := ci.Target(var_if, " ", " ::")
		var_compare := ci.Target(var_if, ":: ", " !")
		var_toExec := ci.Target(var_if, "! ", "(")
		target, err := gc.Get(var_target)
		if err != nil {
			panic(err)
		}

		targetStr, ok := target.(string)
		if !ok {
			panic("target não é uma string")
		}

		compare, err := gc.Get(var_compare)
		if err != nil {
			panic(err)
		}
		compareStr, ok := compare.(string)
		if !ok {
			panic("compare não é uma string")
		}

		toExec, err := gc.Get(var_toExec)
		if err != nil {
			panic(err)
		}
		toExecStr, ok := toExec.(string)
		if !ok {
			panic("compare não é uma string")
		}
		if targetStr >= compareStr {
			execFunc(toExecStr)
		}
		replacer := "$if[>=] " + targetStr + " :: " + compareStr + " ! " + toExecStr + "();"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_for := []string{"!for[num++]"}

	// Conte as palavras-chave no conteúdo do arquivo
	counts_for := countKeywords(string(content), keywords_for)
	for i := 1; i <= counts_for; i++ {
		var_for := ci.Target(string(content), "!for[num++]", "}")
		var_init := ci.Target(var_for, "(", " ::")
		var_condition := ci.Target(var_for, ":: ", ")")
		var_code := ci.Target(var_for, "{", "endloop")
		init, err := strconv.Atoi(var_init)
		if err != nil {
			panic(err)
		}
		condition, err := strconv.Atoi(var_condition)
		if err != nil {
			panic(err)
		}
		for i:=init; i<condition; i++ {              
			execFunc(var_code)       
		} 
		replacer := `!for[num++] (`+var_init+` :: `+var_condition+`) {
`+var_code+`
			`
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}




}

func execFunc(data string) {
	logger.Info(context.Background(), "Executing function: "+data)
}

func main() {

	ReadWysb("test.wys")

}
