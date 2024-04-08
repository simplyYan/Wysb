package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/simplyYan/GalaktaGlare"

	"github.com/robertkrimen/otto"

	"math"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/bluele/gcache"
	lagra "github.com/simplyYan/LAGRA"
	w7 "github.com/simplyYan/W7DTH"
	"github.com/simplyYan/cutinfo"
)

var gc = gcache.New(20).
	LRU().
	Build()

var logger, err = lagra.New(`
	log_file = "output.log"
	`)

var vm = otto.New()

var glare = galaktaglare.New()
var web = fiber.New()

func countKeywords(text string, keywords []string) int {
	total := 0

	for _, keyword := range keywords {
		total += strings.Count(text, keyword)
	}

	return total
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func checkInternetConnection() bool {
	_, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		return false
	}
	return true
}

func execShellScript(script string) error {
	var comando string

	switch runtime.GOOS {
	case "windows":

		comando = "cmd /c " + script
	default:

		comando = "sh -c '" + script + "'"
	}

	cmd := exec.Command(comando)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func execBatchScript(script string) error {
	var comando string

	switch runtime.GOOS {
	case "windows":

		comando = "cmd /c " + script
	default:

		return fmt.Errorf("BatchScripts are not supported on Unix systems")
	}

	cmd := exec.Command(comando)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
    fileToZip, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer fileToZip.Close()

    fileInfo, err := fileToZip.Stat()
    if err != nil {
        return err
    }

    header, err := zip.FileInfoHeader(fileInfo)
    if err != nil {
        return err
    }

    header.Name = filepath.Base(filename)

    writer, err := zipWriter.CreateHeader(header)
    if err != nil {
        return err
    }

    _, err = io.Copy(writer, fileToZip)
    return err
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	randomBytes := make([]byte, length)
	for i := range randomBytes {

		randomBytes[i] = charset[rand.Intn(len(charset))]
	}

	return string(randomBytes)
}

func ReadWysb(filename string) {
	file, err := os.Open(filename) 
	if err != nil {
		log.Fatal(err)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	keywords := []string{"$[int32]"}

	counts := countKeywords(string(content), keywords)

	keywords_float := []string{"$[float32]"}

	counts_float := countKeywords(string(content), keywords_float)

	ci := cutinfo.New()

	keywords_comments := []string{"!:!"}

	counts_comments := countKeywords(string(content), keywords_comments)
	for i := 1; i == counts_comments; i++ {
		ccontent := ci.Target(string(content), "!:!", "!:!")
		replacer := "!:!" + ccontent + "!:!"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}
	keywords_extends := []string{"<extends"}

	counts_extends := countKeywords(string(content), keywords_extends)
	for i := 1; i == counts_extends; i++ {
		extender := ci.Target(string(content), "<extends ", "!>")
		ReadWysb(extender)
		replacer := "<extends " + extender + "!>"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_convertStr := []string{"<to.string"}

	counts_convertStr := countKeywords(string(content), keywords_convertStr)
	for i := 1; i <= counts_convertStr; i++ {
		convert := ci.Target(string(content), "<to.string ", "!>")

		replacer := "<to.string " + convert + "!>"
		toString := "'"+convert+"'"
		content = []byte(strings.Replace(string(content), replacer, toString, i))

	}

	keywords_sum := []string{"<math.sum"}

	counts_sum := countKeywords(string(content), keywords_sum)
	for i := 1; i <= counts_sum; i++ {
		var_sum := ci.Target(string(content), "<math.sum", ">")
		var_num1 := ci.Target(var_sum, " ", " +")
		var_num2 := ci.Target(var_sum, "+ ", "!")
		num1, err := strconv.Atoi(var_num1)
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			return
		}
		num2, err := strconv.Atoi(var_num2)
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			return
		}

		replacer := "<math.sum " + var_num1 + " + " + var_num2 + "!>"
		rsult := num1 + num2
		rsult_str := strconv.Itoa(rsult)
		content = []byte(strings.Replace(string(content), replacer, rsult_str, i))

	}

	keywords_sub := []string{"<math.sub"}

	counts_sub := countKeywords(string(content), keywords_sub)
	for i := 1; i <= counts_sub; i++ {
		var_sum := ci.Target(string(content), "<math.sub", ">")
		var_num1 := ci.Target(var_sum, " ", " -")
		var_num2 := ci.Target(var_sum, "- ", "!")
		num1, err := strconv.Atoi(var_num1)
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			return
		}
		num2, err := strconv.Atoi(var_num2)
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			return
		}

		replacer := "<math.sub " + var_num1 + " - " + var_num2 + "!>"
		rsult := num1 - num2
		rsult_str := strconv.Itoa(rsult)
		content = []byte(strings.Replace(string(content), replacer, rsult_str, i))

	}

	keywords_div := []string{"<math.div"}

	counts_div := countKeywords(string(content), keywords_div)
	for i := 1; i <= counts_div; i++ {
		var_sum := ci.Target(string(content), "<math.div", ">")
		var_num1 := ci.Target(var_sum, " ", " /")
		var_num2 := ci.Target(var_sum, "/ ", "!")
		num1, err := strconv.Atoi(var_num1)
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			return
		}
		num2, err := strconv.Atoi(var_num2)
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			return
		}

		replacer := "<math.div " + var_num1 + " / " + var_num2 + "!>"
		rsult := num1 / num2
		rsult_str := strconv.Itoa(rsult)
		content = []byte(strings.Replace(string(content), replacer, rsult_str, i))

	}

	keywords_mult := []string{"<math.mult"}

	counts_mult := countKeywords(string(content), keywords_mult)
	for i := 1; i <= counts_mult; i++ {
		var_sum := ci.Target(string(content), "<math.mult", ">")
		var_num1 := ci.Target(var_sum, " ", " *")
		var_num2 := ci.Target(var_sum, "* ", "!")
		num1, err := strconv.Atoi(var_num1)
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			return
		}
		num2, err := strconv.Atoi(var_num2)
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			return
		}

		replacer := "<math.mult " + var_num1 + " * " + var_num2 + "!>"
		rsult := num1 * num2
		rsult_str := strconv.Itoa(rsult)
		content = []byte(strings.Replace(string(content), replacer, rsult_str, i))

	}

	keywords_pi := []string{"<math.pi"}

	counts_pi := countKeywords(string(content), keywords_pi)
	for i := 1; i <= counts_pi; i++ {
		replacer := "<math.pi!>"
		content = []byte(strings.Replace(string(content), replacer, "3.14159265358979323846", i))

	}

	keywords_pow := []string{"<math.pow"}

	counts_pow := countKeywords(string(content), keywords_pow)
	for i := 1; i <= counts_pow; i++ {
		var_sum := ci.Target(string(content), "<math.pow", ">")
		var_num1 := ci.Target(var_sum, " ", " ::")
		var_num2 := ci.Target(var_sum, ":: ", "!")
		num1, err := strconv.Atoi(var_num1)
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			return
		}
		num2, err := strconv.Atoi(var_num2)
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			return
		}

		replacer := "<math.pow " + var_num1 + " :: " + var_num2 + "!>"
		rsult := math.Pow(float64(num1), float64(num2))
		rsult_str := strconv.FormatFloat(rsult, 'f', -1, 64)
		content = []byte(strings.Replace(string(content), replacer, rsult_str, i))

	}

	keywords_sqrt := []string{"<math.sqrt"}

	counts_sqrt := countKeywords(string(content), keywords_sqrt)
	for i := 1; i <= counts_sqrt; i++ {
		var_sum := ci.Target(string(content), "<math.sqrt", ">")
		var_num1 := ci.Target(var_sum, " ", "!")
		num1, err := strconv.Atoi(var_num1)
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			return
		}

		replacer := "<math.sqrt " + var_num1 + "!>"
		rsult := math.Sqrt(float64(num1))
		rsult_str := strconv.FormatFloat(rsult, 'f', -1, 64)
		content = []byte(strings.Replace(string(content), replacer, rsult_str, i))

	}

	keywords_input := []string{"<io.input"}

	counts_input := countKeywords(string(content), keywords_input)
	for i := 1; i <= counts_input; i++ {
		input := ci.Target(string(content), "<io.input", ">")
		allocate := ci.Target(input, " ", "!")
		var inp string
		fmt.Scanln(&inp)
		gc.Set(allocate, inp)
		replacer := "<io.input " + allocate + "!>"
		content = []byte(strings.Replace(string(content), replacer, inp, i))
	}

	keywords_openFile := []string{"<os.ReadFile"}

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

	keywords_gglareimg := []string{"<gglare.imagescan"}

	counts_gglareimg := countKeywords(string(content), keywords_gglareimg)
	for i := 1; i == counts_gglareimg; i++ {
		webstatic := ci.Target(string(content), "<gglare.imagescan ", ">")
		db := ci.Target(webstatic, " ", " ::")
		imgtoscan := ci.Target(webstatic, ":: ", "!")
		err := glare.ImageDB(db)
		if err != nil {
			fmt.Println("Error loading the image database:", err)
			return
		}
		similarity, err := glare.ImageScan(imgtoscan)
		if err != nil {
			fmt.Println("Error when comparing images:", err)
			return
		}
		replacer := "<gglare.imagescan " + db + " :: " + imgtoscan + "!>"

		// Converte o float para uma string
		s := strconv.FormatFloat(similarity, 'f', -1, 64)
		content = []byte(strings.Replace(string(content), replacer, string(s), i))

	}	

	keywords_gglaretext := []string{"<gglare.textclassifier"}

	counts_gglaretext := countKeywords(string(content), keywords_gglaretext)
	for i := 1; i == counts_gglaretext; i++ {
		textclassifier := ci.Target(string(content), "<gglare.textclassifier ", ">")
		text := ci.Target(textclassifier, " ", " ::")
		toml := ci.Target(textclassifier, ":: ", "!")

		result, err := glare.TextClassifier(text, toml)
		if err != nil {
			panic(err)
		}


		replacer := "<gglare.textclassifier " + text + " :: " + toml + "!>"


		content = []byte(strings.Replace(string(content), replacer, result, i))

	}	



	for i := 1; i <= counts; i++ {
		var_i32 := ci.Target(string(content), "$[int32]", ";")
		var_i32_name := ci.Target(var_i32, " ", " =")
		var_i32_value := ci.Target(var_i32, "= ", ">")
		gc.Set(var_i32_name, var_i32_value)
		vm.Set(var_i32_name, var_i32_value)
		replacer := "$[int32] " + var_i32_name + " = " + var_i32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	for i := 1; i <= counts_float; i++ {
		var_f32 := ci.Target(string(content), "$[float32]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		vm.Set(var_f32_name, var_f32_value)
		replacer := "$[float32] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}
	keywords_float64 := []string{"$[float64]"}

	counts_float64 := countKeywords(string(content), keywords_float64)
	for i := 1; i <= counts_float64; i++ {
		var_f32 := ci.Target(string(content), "$[float64]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		vm.Set(var_f32_name, var_f32_value)
		replacer := "$[float64] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}
	keywords_static := []string{"$[static]"}

	counts_static := countKeywords(string(content), keywords_static)
	for i := 1; i <= counts_static; i++ {
		var_f32 := ci.Target(string(content), "$[static]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		vm.Set(var_f32_name, var_f32_value)
		replacer := "$[static] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}
	keywords_float128 := []string{"$[float128]"}

	counts_float128 := countKeywords(string(content), keywords_float128)
	for i := 1; i <= counts_float128; i++ {
		var_f32 := ci.Target(string(content), "$[float128]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		vm.Set(var_f32_name, var_f32_value)
		replacer := "$[float128] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_int64 := []string{"$[int64]"}

	counts_int64 := countKeywords(string(content), keywords_int64)
	for i := 1; i <= counts_int64; i++ {
		var_f32 := ci.Target(string(content), "$[int64]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		replacer := "$[int64] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}
	keywords_int128 := []string{"$[int128]"}

	counts_int128 := countKeywords(string(content), keywords_int128)
	for i := 1; i <= counts_int128; i++ {
		var_f32 := ci.Target(string(content), "$[int128]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		vm.Set(var_f32_name, var_f32_value)
		replacer := "$[int128] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}
	keywords_string := []string{"$[string]"}

	counts_string := countKeywords(string(content), keywords_string)
	for i := 1; i <= counts_string; i++ {
		var_f32 := ci.Target(string(content), "$[string]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		vm.Set(var_f32_name, var_f32_value)
		replacer := "$[string] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_bool := []string{"$[bool]"}

	counts_bool := countKeywords(string(content), keywords_bool)
	for i := 1; i <= counts_bool; i++ {
		var_f32 := ci.Target(string(content), "$[bool]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		vm.Set(var_f32_name, var_f32_value)
		replacer := "$[bool] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_array := []string{"$[array]"}

	counts_array := countKeywords(string(content), keywords_array)
	for i := 1; i <= counts_array; i++ {
		var_f32 := ci.Target(string(content), "$[array]", ";")
		var_f32_name := ci.Target(var_f32, " ", " =")
		var_f32_value := ci.Target(var_f32, "= ", ">")
		gc.Set(var_f32_name, var_f32_value)
		vm.Set(var_f32_name, var_f32_value)
		replacer := "$[array] " + var_f32_name + " = " + var_f32_value
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_setarg := []string{"<wysb.SetArg"}

	counts_setarg := countKeywords(string(content), keywords_setarg)
	for i := 1; i <= counts_setarg; i++ {
		var_f32 := ci.Target(string(content), "<wysb.SetArg", ">")
		argname := ci.Target(var_f32, " ", " ::")
		argvalue := ci.Target(var_f32, ":: ", "!")
		gc.Set(argname, argvalue)
		replacer := "<wysb.UseArg " + argname + " :: " + argvalue + "!>"
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_getvar := []string{"${"}

	counts_getvar := countKeywords(string(content), keywords_getvar)
	for i := 1; i <= counts_getvar; i++ {
		var_f32 := ci.Target(string(content), "${", "}")
		findVar, err := gc.Get(var_f32)
		if err != nil {
			panic(err)
		}

		findVarStr, ok := findVar.(string)
		if !ok {
			panic("findVar is not a string")
		}

		replacer := "${" + var_f32 + "}"
		content = []byte(strings.Replace(string(content), replacer, findVarStr, i))
	}
	keywords_fun := []string{"!fun"}

	counts_fun := countKeywords(string(content), keywords_fun)
	for i := 1; i == counts_fun; i++ {
		var_fun := ci.Target(string(content), "!fun", "{")
		var_fun_untyped_str := strings.Replace(var_fun, "[string]", "", -1)
		var_fun_untyped_int := strings.Replace(var_fun_untyped_str, "[int]", "", -1)
		var_fun_untyped_float := strings.Replace(var_fun_untyped_int, "[float]", "", -1)
		var_fun_untyped_bool := strings.Replace(var_fun_untyped_float, "[bool]", "", -1)
		var_fun_name := ci.Target(var_fun_untyped_bool, " ", "(")
		var_fun_args := ci.Target(var_fun_untyped_bool, "(", ")")
		var_fun_argsTYPED := ci.Target(var_fun, "(", ")")
		var_fun_value := ci.Target(string(content), "{", "}")

		jsCode := fmt.Sprintf(`
			function %s(%s) {
				%s
			}
		`, var_fun_name, var_fun_args, var_fun_value)

		_, err := vm.Run(jsCode)
		if err != nil {
			panic(err)
		}
		gc.Set(var_fun_name, var_fun_value)
		replacer := "!fun " + var_fun_name + "(" + var_fun_argsTYPED + ")" + "{" + `
		` + var_fun_value + `
	}
		`
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_callfn := []string{"@:"}

	counts_callfn := countKeywords(string(content), keywords_callfn)
	for i := 1; i == counts_callfn; i++ {
		var_Callfun := ci.Target(string(content), "@:", "?")
		fnname := ci.Target(var_Callfun, ":", "(")
		fnargs := ci.Target(var_Callfun, "(", ")")
		jsCode := fmt.Sprintf(`
			%s(%s);
		`, fnname, fnargs)

		result, err := vm.Run(jsCode)
		if err != nil {
			panic(err)
		}
		if value, err := result.ToString(); err == nil {
			finalname := "exec::" + fnname
			gc.Set(finalname, value)
		} else {
			fmt.Println("Error getting the result of JavaScript execution:", err)
		}
		replacer := "@::" + fnname + "(" + fnargs + ")?"
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_if := []string{"$if[>]"}

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
			panic("target is not a string")
		}

		compare, err := gc.Get(var_compare)
		if err != nil {
			panic(err)
		}
		compareStr, ok := compare.(string)
		if !ok {
			panic("compare is not a string")
		}

		toExec, err := gc.Get(var_toExec)
		if err != nil {
			panic(err)
		}
		toExecStr, ok := toExec.(string)
		if !ok {
			panic("compare is not a string")
		}
		if targetStr > compareStr {
			runAddon, err := vm.Run(string(toExecStr))
			if err != nil {
				panic(err)
			}
			if value, err := runAddon.ToString(); err == nil {
			gc.Set("exec::"+targetStr+""+compareStr+"", value)
			} else {
				fmt.Println("Error getting the result of JavaScript execution:", err)
			}
		}
		replacer := "$if[>] " + targetStr + " :: " + compareStr + " ! " + toExecStr + "();"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_ifLess := []string{"$if[<]"}

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
			panic("target is not a string")
		}

		compare, err := gc.Get(var_compare)
		if err != nil {
			panic(err)
		}
		compareStr, ok := compare.(string)
		if !ok {
			panic("compare is not a string")
		}

		toExec, err := gc.Get(var_toExec)
		if err != nil {
			panic(err)
		}
		toExecStr, ok := toExec.(string)
		if !ok {
			panic("compare is not a string")
		}
		if targetStr < compareStr {
			runAddon, err := vm.Run(string(toExecStr))
			if err != nil {
				panic(err)
			}
			if value, err := runAddon.ToString(); err == nil {
			gc.Set("exec::"+targetStr+""+compareStr+"", value)
			} else {
				fmt.Println("Error getting the result of JavaScript execution:", err)
			}
		}
		replacer := "$if[<] " + targetStr + " :: " + compareStr + " ! " + toExecStr + "();"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_ifEqual := []string{"$if[==]"}

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
			panic("target is not a string")
		}

		compare, err := gc.Get(var_compare)
		if err != nil {
			panic(err)
		}
		compareStr, ok := compare.(string)
		if !ok {
			panic("compare is not a string")
		}

		toExec, err := gc.Get(var_toExec)
		if err != nil {
			panic(err)
		}
		toExecStr, ok := toExec.(string)
		if !ok {
			panic("compare is not a string")
		}
		if targetStr == compareStr {
			runAddon, err := vm.Run(string(toExecStr))
			if err != nil {
				panic(err)
			}
			if value, err := runAddon.ToString(); err == nil {
			gc.Set("exec::"+targetStr+""+compareStr+"", value)
			} else {
				fmt.Println("Error getting the result of JavaScript execution:", err)
			}
		}
		replacer := "$if[==] " + targetStr + " :: " + compareStr + " ! " + toExecStr + "();"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_ifLessEqual := []string{"$if[<=]"}

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
			panic("target is not a string")
		}

		compare, err := gc.Get(var_compare)
		if err != nil {
			panic(err)
		}
		compareStr, ok := compare.(string)
		if !ok {
			panic("compare is not a string")
		}

		toExec, err := gc.Get(var_toExec)
		if err != nil {
			panic(err)
		}
		toExecStr, ok := toExec.(string)
		if !ok {
			panic("compare is not a string")
		}
		if targetStr <= compareStr {
			runAddon, err := vm.Run(string(toExecStr))
			if err != nil {
				panic(err)
			}
			if value, err := runAddon.ToString(); err == nil {
			gc.Set("exec::"+targetStr+""+compareStr+"", value)
			} else {
				fmt.Println("Error getting the result of JavaScript execution:", err)
			}
		}
		replacer := "$if[<=] " + targetStr + " :: " + compareStr + " ! " + toExecStr + "();"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_getvar2 := []string{"&{"}

	counts_getvar2 := countKeywords(string(content), keywords_getvar2)
	for i := 1; i <= counts_getvar2; i++ {
		var_f32 := ci.Target(string(content), "&{", "}")
		findVar, err := gc.Get(var_f32)
		if err != nil {
			panic(err)
		}

		findVarStr, ok := findVar.(string)
		if !ok {
			panic("findVar não é uma string")
		}

		replacer := "&{" + var_f32 + "}"
		content = []byte(strings.Replace(string(content), replacer, findVarStr, i))
	}

	keywords_ifBiggerEqual := []string{"$if[>=]"}

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
			panic("target is not a string")
		}

		compare, err := gc.Get(var_compare)
		if err != nil {
			panic(err)
		}
		compareStr, ok := compare.(string)
		if !ok {
			panic("compare is not a string")
		}

		toExec, err := gc.Get(var_toExec)
		if err != nil {
			panic(err)
		}
		toExecStr, ok := toExec.(string)
		if !ok {
			panic("compare is not a string")
		}
		if targetStr >= compareStr {
			runAddon, err := vm.Run(string(toExecStr))
			if err != nil {
				panic(err)
			}
			if value, err := runAddon.ToString(); err == nil {
			gc.Set("exec::"+targetStr+""+compareStr+"", value)
			} else {
				fmt.Println("Error getting the result of JavaScript execution:", err)
			}
		}
		replacer := "$if[>=] " + targetStr + " :: " + compareStr + " ! " + toExecStr + "();"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_notequal := []string{"$if[!=]"}

	counts_notequal := countKeywords(string(content), keywords_notequal)
	for i := 1; i == counts_notequal; i++ {
		var_if := ci.Target(string(content), "$if[!=]", ";")
		var_target := ci.Target(var_if, " ", " ::")
		var_compare := ci.Target(var_if, ":: ", " !")
		var_toExec := ci.Target(var_if, "! ", "(")
		target, err := gc.Get(var_target)
		if err != nil {
			panic(err)
		}

		targetStr, ok := target.(string)
		if !ok {
			panic("target is not a string")
		}

		compare, err := gc.Get(var_compare)
		if err != nil {
			panic(err)
		}
		compareStr, ok := compare.(string)
		if !ok {
			panic("compare is not a string")
		}

		toExec, err := gc.Get(var_toExec)
		if err != nil {
			panic(err)
		}
		toExecStr, ok := toExec.(string)
		if !ok {
			panic("compare is not a string")
		}
		if targetStr != compareStr {
			runAddon, err := vm.Run(string(toExecStr))
			if err != nil {
				panic(err)
			}
			if value, err := runAddon.ToString(); err == nil {
			gc.Set("exec::"+targetStr+""+compareStr+"", value)
			} else {
				fmt.Println("Error getting the result of JavaScript execution:", err)
			}
		}
		replacer := "$if[!=] " + targetStr + " :: " + compareStr + " ! " + toExecStr + "();"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_for := []string{"!for[num++]"}

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
		for i := init; i < condition; i++ {
			execFunc(var_code)
		}
		replacer := `!for[num++] (` + var_init + ` :: ` + var_condition + `) {
` + var_code + `
			`
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_writefile := []string{"<os.WriteFile"}

	counts_writefile := countKeywords(string(content), keywords_writefile)
	for i := 1; i <= counts_writefile; i++ {
		fileopen := ci.Target(string(content), "<os.WriteFile", ">")
		filedir := ci.Target(fileopen, " ", " ::")
		filecontent := ci.Target(fileopen, ":: ", "!")
		data := []byte(filecontent)
		if err := os.WriteFile(filedir, data, 0644); err != nil {
			panic(err)
		}
		replacer := "<os.WriteFile " + filedir + " :: " + filecontent + "!>"
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_encrypt := []string{"<cryptography.encrypt"}

	counts_encrypt := countKeywords(string(content), keywords_encrypt)
	for i := 1; i <= counts_encrypt; i++ {
		cryptography := ci.Target(string(content), "<cryptography.encrypt", ">")
		targetstring := ci.Target(cryptography, " ", " ::")
		targetkey := ci.Target(cryptography, ":: ", "!")
		w := w7.New()
		key, err := w.Key(targetkey)
		if err != nil {
			panic(err)
		}
		gc.Set("encryptKey", key)
		encrypted, err := w.Encrypt(targetstring)
		if err != nil {
			panic(err)
		}

		replacer := "<cryptography.encrypt " + targetkey + " :: " + targetstring + "!>"
		content = []byte(strings.Replace(string(content), replacer, encrypted, i))
	}

	keywords_decrypt := []string{"<cryptography.decrypt"}

	counts_decrypt := countKeywords(string(content), keywords_decrypt)
	for i := 1; i <= counts_decrypt; i++ {
		cryptography := ci.Target(string(content), "<cryptography.decrypt", ">")
		targetstring := ci.Target(cryptography, " ", " ::")
		targetkey := ci.Target(cryptography, ":: ", "!")
		w := w7.New()
		key, err := w.Key(targetkey)
		if err != nil {
			panic(err)
		}
		gc.Set("encryptKey", key)
		encrypted, err := w.Decrypt(targetstring)
		if err != nil {
			panic(err)
		}

		replacer := "<cryptography.decrypt " + targetkey + " :: " + targetstring + "!>"
		content = []byte(strings.Replace(string(content), replacer, encrypted, i))
	}

	keywords_logger := []string{"<logger.Send"}

	counts_logger := countKeywords(string(content), keywords_logger)
	for i := 1; i <= counts_logger; i++ {
		lgr := ci.Target(string(content), "<logger.Send", ">")
		lgrtype := ci.Target(lgr, " ", " ::")
		lgrcontent := ci.Target(lgr, ":: ", "!")
		if lgrtype == "Info" {
			logger.Info(context.Background(), lgrcontent)
		} else if lgrtype == "Warn" {
			logger.Warn(context.Background(), lgrcontent)
		} else if lgrtype == "Error" {
			logger.Error(context.Background(), lgrcontent)
		} else {
			panic("Undefined: " + lgrtype)
		}
		replacer := "<logger.Send " + lgrtype + " :: " + lgrcontent + "!>"
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_allocate := []string{"<wysb.Allocate"}

	counts_allocate := countKeywords(string(content), keywords_allocate)
	for i := 1; i <= counts_allocate; i++ {
		fileopen := ci.Target(string(content), "<wysb.Allocate", ">")
		altype := ci.Target(fileopen, " ", " ::")
		alvarname := ci.Target(fileopen, ":: ", "!")
		if altype == "string" {
			replacer := "<wysb.Allocate " + altype + " :: " + alvarname + "!>"
			content = []byte(strings.Replace(string(content), replacer, "-_A string has been allocated here_-", i))
			gc.Set(alvarname, altype)
		} else if altype == "int" {
			replacer := "<wysb.Allocate " + altype + " :: " + alvarname + "!>"
			content = []byte(strings.Replace(string(content), replacer, "-_A int has been allocated here_-", i))
			gc.Set(alvarname, altype)
		} else if altype == "float" {
			replacer := "<wysb.Allocate " + altype + " :: " + alvarname + "!>"
			content = []byte(strings.Replace(string(content), replacer, "-_A float has been allocated here_-", i))
			gc.Set(alvarname, altype)
		} else if altype == "bool" {
			replacer := "<wysb.Allocate " + altype + " :: " + alvarname + "!>"
			content = []byte(strings.Replace(string(content), replacer, "-_A bool has been allocated here_-", i))
			gc.Set(alvarname, altype)
		} else {
			panic("Undefined: " + altype)
		}

	}

	keywords_free := []string{"<wysb.Free"}

	counts_free := countKeywords(string(content), keywords_free)
	for i := 1; i <= counts_free; i++ {
		free := ci.Target(string(content), "<wysb.Free", ">")
		vartarget := ci.Target(free, " ", "!")
		gc.Remove(vartarget)
		replacer := "<wysb.Free " + vartarget + "!>"
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_renamefile := []string{"<os.RenameFile"}

	counts_renamefile := countKeywords(string(content), keywords_renamefile)
	for i := 1; i <= counts_renamefile; i++ {
		fileopen := ci.Target(string(content), "<os.RenameFile", ">")
		filedir := ci.Target(fileopen, " ", " ::")
		newname := ci.Target(fileopen, ":: ", "!")
		e := os.Rename(filedir, newname)
		if e != nil {
			panic(e)
		}
		replacer := "<os.RenameFile " + filedir + " :: " + newname + "!>"
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_removefile := []string{"<os.RemoveFile"}

	counts_removefile := countKeywords(string(content), keywords_removefile)
	for i := 1; i <= counts_removefile; i++ {
		fileopen := ci.Target(string(content), "<os.RemoveFile", ">")
		filedir := ci.Target(fileopen, " ", "!")
		err := os.Remove(filedir)
		if err != nil {
			panic(err)
		}
		replacer := "<os.RemoveFile " + filedir + "!>"
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_stringsreplace := []string{"<strings.Replace"}

	counts_stringsreplace := countKeywords(string(content), keywords_stringsreplace)
	for i := 1; i <= counts_stringsreplace; i++ {
		fileopen := ci.Target(string(content), "<strings.Replace", ">")
		targetvar := ci.Target(fileopen, " ", " ::")
		targetselection := ci.Target(fileopen, ":: ", " ~")
		newselection := ci.Target(fileopen, "~ ", "!")
		targetvar_content, err := gc.Get(targetvar)
		if err != nil {
			panic(err)
		}
		tvStr, ok := targetvar_content.(string)
		if !ok {
			panic("tvStr is not a string")
		}
		edited := strings.Replace(tvStr, targetselection, string(newselection), -1)
		gc.Set(targetvar, edited)

		replacer := "<strings.Replace " + targetvar + " :: " + targetselection + " ~ " + newselection + "!>"
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_rand := []string{"<rand.Num"}

	counts_rand := countKeywords(string(content), keywords_rand)
	for i := 1; i <= counts_rand; i++ {
		fileopen := ci.Target(string(content), "<rand.Num", ">")
		randtype := ci.Target(fileopen, " ", " ::")
		randsize := ci.Target(fileopen, ":: ", "!")
		rand.Seed(time.Now().UnixNano())

		if randtype == "float" {
			randsize_int, err := strconv.Atoi(randsize)
			if err != nil {
				panic(err)
			}
			minFloat := float64(randsize_int) * 0.01
			maxFloat := float64(randsize_int)
			randomInRange := minFloat + rand.Float64()*(maxFloat-minFloat)
			replacer := "<rand.Num " + randtype + " :: " + randsize + "!>"
			randomInRange_float := strconv.FormatFloat(float64(randomInRange), 'f', -1, 64)
			content = []byte(strings.Replace(string(content), replacer, randomInRange_float, i))
		} else if randtype == "int" {
			randsize_int, err := strconv.Atoi(randsize)
			if err != nil {
				panic(err)
			}
			randomInt := rand.Intn(randsize_int)
			replacer := "<rand.Num " + randtype + " :: " + randsize + "!>"
			randomIntStr := strconv.Itoa(randomInt)
			content = []byte(strings.Replace(string(content), replacer, randomIntStr, i))
		} else if randtype == "string" {
			randsize_int, err := strconv.Atoi(randsize)
			if err != nil {
				panic(err)
			}
			randomStr := randomString(randsize_int)
			replacer := "<rand.Num " + randtype + " :: " + randsize + "!>"
			content = []byte(strings.Replace(string(content), replacer, randomStr, i))
		} else {
			panic("Undefined: " + randtype)
		}

	}

	keywords_usearg := []string{"<wysb.UseArg"}

	counts_usearg := countKeywords(string(content), keywords_usearg)
	for i := 1; i <= counts_usearg; i++ {
		var_f32 := ci.Target(string(content), "<wysb.UseArg", ">")
		var_f32_name := ci.Target(var_f32, " ", "!")
		targetvar, err := gc.Get(var_f32_name)
		if err != nil {
			fmt.Println("")
		}
		str := targetvar.(string);
		replacer := "<wysb.UseArg " + var_f32_name + "!>"
		content = []byte(strings.Replace(string(content), replacer, str, i))
	}

	keywords_execShellScript := []string{"<os.ShellScript"}

	counts_execShellScript := countKeywords(string(content), keywords_execShellScript)
	for i := 1; i <= counts_execShellScript; i++ {
		var_f32 := ci.Target(string(content), "<os.ShellScript", ">")
		var_f32_name := ci.Target(var_f32, "[] ", "!")
		err := execShellScript(var_f32_name)
		if err != nil {
			fmt.Println("Erro ao executar ShellScript:", err)
		}		
		replacer := "<os.ShellScript [] " + var_f32_name + "!>"
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_execBatchScript := []string{"<os.BatchScript"}

	counts_execBatchScript := countKeywords(string(content), keywords_execBatchScript)
	for i := 1; i <= counts_execBatchScript; i++ {
		var_f32 := ci.Target(string(content), "<os.BatchScript", ">")
		var_f32_name := ci.Target(var_f32, "[] ", "!")
		err := execBatchScript(var_f32_name)
		if err != nil {
			fmt.Println("Erro ao executar BatchScript:", err)
		}		
		replacer := "<os.BatchScript [] " + var_f32_name + "!>"
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_addon := []string{"<wysb.Addon"}

	counts_addon := countKeywords(string(content), keywords_addon)
	for i := 1; i <= counts_addon; i++ {
		addon := ci.Target(string(content), "<wysb.Addon", ">")
		funcname := ci.Target(addon, " ", "!")

		replacer := "<wysb.Addon " + funcname + "!>"
		content = []byte(strings.Replace(string(content), replacer, "", i))
	}

	keywords_println := []string{"println("}

	counts_println := countKeywords(string(content), keywords_println)
	for i := 1; i == counts_println; i++ {
		var_fun := ci.Target(string(content), "println(", ")")
		logger.Info(context.Background(), var_fun)
		replacer := "println(" + var_fun + ")"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_webget := []string{"<web.get"}

	counts_webget := countKeywords(string(content), keywords_webget)
	for i := 1; i == counts_webget; i++ {
		webget := ci.Target(string(content), "<web.get", ">")
		route := ci.Target(webget, " ", " ::")
		returnn := ci.Target(webget, ":: ", "!")
		web.Get(route, func(c *fiber.Ctx) error {
			c.Type("html")
			return c.SendString(returnn)
			
		})
		replacer := "<web.get " + route + " :: " + returnn + "!>"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_weblisten := []string{"<web.listen"}

	counts_weblisten := countKeywords(string(content), keywords_weblisten)
	for i := 1; i == counts_weblisten; i++ {
		weblisten := ci.Target(string(content), "<web.listen ", "!>")
		web.Listen(weblisten)
		replacer := "<web.listen " + weblisten + "!>"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_webstatic := []string{"<web.static"}

	counts_webstatic := countKeywords(string(content), keywords_webstatic)
	for i := 1; i == counts_webstatic; i++ {
		webstatic := ci.Target(string(content), "<web.static ", ">")
		route := ci.Target(webstatic, " ", " ::")
		dir := ci.Target(webstatic, ":: ", "!")
		web.Static(route, dir)
		replacer := "<web.static " + route + " :: " + dir + "!>"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}

	keywords_webpost := []string{"<web.post"}

	counts_webpost := countKeywords(string(content), keywords_webpost)
	for i := 1; i == counts_webpost; i++ {
		webpost := ci.Target(string(content), "<web.post", ">")
		route := ci.Target(webpost, " ", " ::")
		returnn := ci.Target(webpost, ":: ", "!")
		web.Post(route, func(c *fiber.Ctx) error {
			c.Type("html")
			return c.SendString(returnn)
			
		})
		replacer := "<web.post " + route + " :: " + returnn + "!>"
		content = []byte(strings.Replace(string(content), replacer, "", i))

	}



}

func execFunc(data string) {
	logger.Info(context.Background(), "Executing function: "+data)
}

func main() {
	url := "https://raw.githubusercontent.com/simplyYan/Wysb/main/add-ons/wysb-addon.js"
	fileName := "wysb-addon.js"

	if checkInternetConnection() {
		content, err := downloadFile(url)
		if err != nil {
			fmt.Println("Error downloading the file:", err)
			return
		}

		err = os.WriteFile(fileName, content, 0644)
		if err != nil {
			fmt.Println("Error writing the file:", err)
			return
		}

	} else {

		_, err := os.Stat(fileName)
		if os.IsNotExist(err) {
			fmt.Println("Error: File does not exist locally.")
			return
		}

		addoncontent, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Println("Error reading the file:", err)
			return
		}

		_, err = vm.Run(string(addoncontent))
		if err != nil {
			panic(err)
		}
		

	}

    if _, err := os.Stat("main.wys"); err == nil {

        compiledcontent, err := os.ReadFile("main.wys")
        if err != nil {
            fmt.Println("Error reading the file:", err)
            return
        }
        ReadWysb(string(compiledcontent))
    } else if os.IsNotExist(err) {

		fmt.Println("")
    } else {

        fmt.Println("Error verifying the existence of the file:", err)
    }

	var input string
	fmt.Println(`

 _       __           __  
| |     / /_  _______/ /_ 
| | /| / / / / / ___/ __ \
| |/ |/ / /_/ (__  ) /_/ /
|__/|__/\__, /____/_____/ 
       /____/

Welcome to Wysb. To learn how to use the commands, you can use the "wysb help" command.`)
	fmt.Scanln(&input)

	finalinput := strings.Replace(input, "wysb ", "", 1)
	if (finalinput == "help") {
		fmt.Println(`
		Complete list of commands:
		- run: Used to run/test a Wysb file. After using it, you must pass the name of the file.

		- compile: To convert a Wysb file into an executable. After using it, you must pass the file name.

		- cardwmy: Used to initialize "cardwmy", the Wysb package manager. After using it, you must enter the URL of the package to be downloaded.
		`)
	} else if (finalinput == "run") {
		fmt.Println("Enter the name of the Wysb file you want to run: ")
		var runinput string
		fmt.Scanln(&runinput)
		fmt.Println("Running the "+runinput+" file")
		ReadWysb(runinput)

	} else if (finalinput == "compile") {
		fmt.Println("Enter the name of the Wysb file you want to compile: ")
		var programname string
		fmt.Scan(&programname)
		fmt.Println("Now type in the name of the executable (including the file extension) of the Wysb compiler: ")
		var compilerexec string
		fmt.Scan(&compilerexec)
		files := []string{"main.wys", compilerexec}

    zipFilename := programname

    newZipFile, err := os.Create(zipFilename)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer newZipFile.Close()

    zipWriter := zip.NewWriter(newZipFile)
    defer zipWriter.Close()

    for _, file := range files {
        if err := addFileToZip(zipWriter, file); err != nil {
            fmt.Println(err)
            return
        }
    }

    fmt.Println("Arquivo zip criado com sucesso:", zipFilename)

	} else if (finalinput == "cardwmy") {
		fmt.Println("Enter the ID of the Cardwmy package you want to download (example: '!PackageAuthor/PackageName!'): ")
		var carinput string 
		fmt.Scanln(&carinput)
		ci := cutinfo.New()
		username := ci.Target(carinput, "!", "/")
		pkg := ci.Target(carinput, "/", "!")
		result := "https://raw.githubusercontent.com/"+username+"/"+pkg+"/main/main.wys"
		dlfile, err := downloadFile(result)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(username, dlfile, 0644)
		if err != nil {
			fmt.Println("Error writing the file:", err)
			return
		}
	} else {
		panic("unknown command: "+finalinput+"")
	}
}
