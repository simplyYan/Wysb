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
	keywords_float64 := []string{"$[float32]"}

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
	keywords_float128 := []string{"$[float32]"}

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

	keywords_int64 := []string{"$[float32]"}

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
	keywords_int128 := []string{"$[float32]"}

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

	keywords_bool := []string{"$[string]"}

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

}

func main() {

	ReadWysb("test.wys")
	idade, err := gc.Get("idade")
	if err != nil {
		panic(err)
	}
	fmt.Println(idade)
	peso, err := gc.Get("peso")
	if err != nil {
		panic(err)
	}
	fmt.Println(peso)

}
