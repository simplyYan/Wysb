package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"sort"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

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
var ci = cutinfo.New()
var vm = otto.New()

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
	defer file.Close()

	var content = []byte("0");
	content, err = io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	processKeywords(content)
}

func processKeywords(content []byte) {
	keywords := map[string]struct {
		handler   func()
		priority  float64
	}{
		"$[int32]":               {func() { handleInt32(content) }, 0.9},
		"$[float32]":             {func() { handleFloat32(content) }, 0.9},
		"$[string]":              {func() { handleString(content) }, 0.9},
		"$[bool]":                {func() { handleBool(content) }, 0.9},
		"$[array]":               {func() { handleArray(content) }, 0.9},
		"$[static]":              {func() { handleStatic(content) }, 0.8},
		"$[float64]":             {func() { handleFloat64(content) }, 0.9},
		"$[float128]":            {func() { handleFloat128(content) }, 0.9},
		"$[int64]":               {func() { handleInt64(content) }, 0.9},
		"$[int128]":              {func() { handleInt128(content) }, 0.9},
		"<extends":               {func() { handleExtends(content) }, 0.7},
		"<to.string":             {func() { handleToString(content) }, 0.7},
		"<math.sum":              {func() { handleMathSum(content) }, 0.9},
		"<math.sub":              {func() { handleMathSub(content) }, 0.9},
		"<math.div":              {func() { handleMathDiv(content) }, 0.9},
		"<math.mult":             {func() { handleMathMult(content) }, 0.9},
		"<math.pi":               {func() { handleMathPi(content) }, 0.9},
		"<math.pow":              {func() { handleMathPow(content) }, 0.9},
		"<math.sqrt":             {func() { handleMathSqrt(content) }, 0.5},
		"<io.input":              {func() { handleInput(content) }, 0.4},
		"<time.measure>":         {func() { handleTimeMeasure(content) }, 0.4},
		"<os.ReadFile":           {func() { handleReadFile(content) }, 0.4},
		"<os.WriteFile":          {func() { handleWriteFile(content) }, 0.4},
		"<cryptography.encrypt":  {func() { handleEncrypt(content) }, 0.3},
		"<cryptography.decrypt":  {func() { handleDecrypt(content) }, 0.3},
		"<logger.Send":           {func() { handleLogger(content) }, 0.3},
		"<wysb.Allocate":         {func() { handleAllocate(content) }, 0.3},
		"<wysb.Free":             {func() { handleFree(content) }, 0.3},
		"<os.RenameFile":         {func() { handleRenameFile(content) }, 0.3},
		"<os.RemoveFile":         {func() { handleRemoveFile(content) }, 0.3},
		"<strings.Replace":       {func() { handleStringsReplace(content) }, 0.3},
		"<rand.Num":              {func() { handleRandNum(content) }, 0.3},
		"<wysb.UseArg":           {func() { handleUseArg(content) }, 0.3},
		"<os.ShellScript":        {func() { handleShellScript(content) }, 0.3},
		"<os.BatchScript":        {func() { handleBatchScript(content) }, 0.3},
		"<wysb.Addon":            {func() { handleAddon(content) }, 0.3},
		"<web.get":               {func() { handleWebGet(content) }, 0.2},
		"<web.listen":            {func() { handleWebListen(content) }, 0.2},
		"<web.static":            {func() { handleWebStatic(content) }, 0.2},
		"<web.post":              {func() { handleWebPost(content) }, 0.2},
		"${":                     {func() { handleFindvar(content) }, 0.1},
		"println(":               {func() { handlePrintln(content) }, 0.1},
	}

	var sortedKeywords []struct {
		keyword  string
		handler  func()
		priority float64
	}
	for keyword, entry := range keywords {
		sortedKeywords = append(sortedKeywords, struct {
			keyword  string
			handler  func()
			priority float64
		}{keyword, entry.handler, entry.priority})
	}

	sort.Slice(sortedKeywords, func(i, j int) bool {
		return sortedKeywords[i].priority > sortedKeywords[j].priority
	})

	for _, entry := range sortedKeywords {
		count := countKeywords(string(content), []string{entry.keyword})
		for i := 1; i <= count; i++ {
			entry.handler()
		}
	}
}

func handleInt32(content []byte) {
	var_i32 := ci.Target(string(content), "$[int32]", ";")
	var_i32_name := ci.Target(var_i32, " ", " =")
	var_i32_value := ci.Target(var_i32, "= ", ">")
	gc.Set(var_i32_name, var_i32_value)
	vm.Set(var_i32_name, var_i32_value)
	replacer := "$[int32] " + var_i32_name + " = " + var_i32_value
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleFindvar(content []byte) {
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
	content = []byte(strings.Replace(string(content), replacer, findVarStr, -1))
	processKeywords(content)
}

func handleFloat32(content []byte) {
	var_f32 := ci.Target(string(content), "$[float32]", ";")
	var_f32_name := ci.Target(var_f32, " ", " =")
	var_f32_value := ci.Target(var_f32, "= ", ">")
	gc.Set(var_f32_name, var_f32_value)
	vm.Set(var_f32_name, var_f32_value)
	replacer := "$[float32] " + var_f32_name + " = " + var_f32_value
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleString(content []byte) {
	var_f32 := ci.Target(string(content), "$[string]", ";")
	var_f32_name := ci.Target(var_f32, " ", " =")
	var_f32_value := ci.Target(var_f32, "= ", ">")
	gc.Set(var_f32_name, var_f32_value)
	vm.Set(var_f32_name, var_f32_value)
	replacer := "$[string] " + var_f32_name + " = " + var_f32_value
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleBool(content []byte) {
	var_f32 := ci.Target(string(content), "$[bool]", ";")
	var_f32_name := ci.Target(var_f32, " ", " =")
	var_f32_value := ci.Target(var_f32, "= ", ">")
	gc.Set(var_f32_name, var_f32_value)
	vm.Set(var_f32_name, var_f32_value)
	replacer := "$[bool] " + var_f32_name + " = " + var_f32_value
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleArray(content []byte) {
	var_f32 := ci.Target(string(content), "$[array]", ";")
	var_f32_name := ci.Target(var_f32, " ", " =")
	var_f32_value := ci.Target(var_f32, "= ", ">")
	gc.Set(var_f32_name, var_f32_value)
	vm.Set(var_f32_name, var_f32_value)
	replacer := "$[array] " + var_f32_name + " = " + var_f32_value
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleStatic(content []byte) {

}

func handleFloat64(content []byte) {
	var_f32 := ci.Target(string(content), "$[float64]", ";")
	var_f32_name := ci.Target(var_f32, " ", " =")
	var_f32_value := ci.Target(var_f32, "= ", ">")
	gc.Set(var_f32_name, var_f32_value)
	vm.Set(var_f32_name, var_f32_value)
	replacer := "$[float64] " + var_f32_name + " = " + var_f32_value
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleFloat128(content []byte) {
	var_f32 := ci.Target(string(content), "$[float128]", ";")
	var_f32_name := ci.Target(var_f32, " ", " =")
	var_f32_value := ci.Target(var_f32, "= ", ">")
	gc.Set(var_f32_name, var_f32_value)
	vm.Set(var_f32_name, var_f32_value)
	replacer := "$[float128] " + var_f32_name + " = " + var_f32_value
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleInt64(content []byte) {
	var_f32 := ci.Target(string(content), "$[int64]", ";")
	var_f32_name := ci.Target(var_f32, " ", " =")
	var_f32_value := ci.Target(var_f32, "= ", ">")
	gc.Set(var_f32_name, var_f32_value)
	replacer := "$[int64] " + var_f32_name + " = " + var_f32_value
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleInt128(content []byte) {
	var_f32 := ci.Target(string(content), "$[int128]", ";")
	var_f32_name := ci.Target(var_f32, " ", " =")
	var_f32_value := ci.Target(var_f32, "= ", ">")
	gc.Set(var_f32_name, var_f32_value)
	vm.Set(var_f32_name, var_f32_value)
	replacer := "$[int128] " + var_f32_name + " = " + var_f32_value
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleExtends(content []byte) {
	extender := ci.Target(string(content), "<extends ", "!>")
	ReadWysb(extender)
	replacer := "<extends " + extender + "!>"
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)

}

func handleToString(content []byte) {
	convert := ci.Target(string(content), "<to.string ", "!>")

	replacer := "<to.string " + convert + "!>"
	toString := "'" + convert + "'"
	content = []byte(strings.Replace(string(content), replacer, toString, -1))
	processKeywords(content)

}

func handleMathSum(content []byte) {
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
	content = []byte(strings.Replace(string(content), replacer, rsult_str, -1))
	processKeywords(content)

}

func handleMathSub(content []byte) {
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
	content = []byte(strings.Replace(string(content), replacer, rsult_str, -1))
	processKeywords(content)

}

func handleMathDiv(content []byte) {
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
	content = []byte(strings.Replace(string(content), replacer, rsult_str, -1))
	processKeywords(content)

}

func handleMathMult(content []byte) {
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
	content = []byte(strings.Replace(string(content), replacer, rsult_str, -1))
	processKeywords(content)

}

func handleMathPi(content []byte) {
	replacer := "<math.pi!>"
	content = []byte(strings.Replace(string(content), replacer, "3.14159265358979323846", -1))
	processKeywords(content)
}

func handleMathPow(content []byte) {
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
	content = []byte(strings.Replace(string(content), replacer, rsult_str, -1))
	processKeywords(content)
}

func handleMathSqrt(content []byte) {
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
	content = []byte(strings.Replace(string(content), replacer, rsult_str, -1))
	processKeywords(content)

}

func handleInput(content []byte) {
	input := ci.Target(string(content), "<io.input", ">")
	allocate := ci.Target(input, " ", "!")
	var inp string
	fmt.Scanln(&inp)
	gc.Set(allocate, inp)
	replacer := "<io.input " + allocate + "!>"
	content = []byte(strings.Replace(string(content), replacer, inp, -1))
	processKeywords(content)
}

func handleReadFile(content []byte) {
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

	content = []byte(strings.Replace(string(content), replacer, string(thiscontent), -1))
	processKeywords(content)
}

func handleWriteFile(content []byte) {
	fileopen := ci.Target(string(content), "<os.WriteFile", ">")
	filedir := ci.Target(fileopen, " ", " ::")
	filecontent := ci.Target(fileopen, ":: ", "!")
	data := []byte(filecontent)
	if err := os.WriteFile(filedir, data, 0644); err != nil {
		panic(err)
	}
	replacer := "<os.WriteFile " + filedir + " :: " + filecontent + "!>"
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleEncrypt(content []byte) {
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
	content = []byte(strings.Replace(string(content), replacer, encrypted, -1))
	processKeywords(content)
}

func handleDecrypt(content []byte) {
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
	content = []byte(strings.Replace(string(content), replacer, encrypted, -1))
	processKeywords(content)
}

func handleLogger(content []byte) {
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
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleAllocate(content []byte) {
	fileopen := ci.Target(string(content), "<wysb.Allocate", ">")
	altype := ci.Target(fileopen, " ", " ::")
	alvarname := ci.Target(fileopen, ":: ", "!")
	if altype == "string" {
		replacer := "<wysb.Allocate " + altype + " :: " + alvarname + "!>"
		content = []byte(strings.Replace(string(content), replacer, "-_A string has been allocated here_-", -1))
		processKeywords(content)
		gc.Set(alvarname, altype)
	} else if altype == "int" {
		replacer := "<wysb.Allocate " + altype + " :: " + alvarname + "!>"
		content = []byte(strings.Replace(string(content), replacer, "-_A int has been allocated here_-", -1))
		processKeywords(content)
		gc.Set(alvarname, altype)
	} else if altype == "float" {
		replacer := "<wysb.Allocate " + altype + " :: " + alvarname + "!>"
		content = []byte(strings.Replace(string(content), replacer, "-_A float has been allocated here_-", -1))
		processKeywords(content)
		gc.Set(alvarname, altype)
	} else if altype == "bool" {
		replacer := "<wysb.Allocate " + altype + " :: " + alvarname + "!>"
		content = []byte(strings.Replace(string(content), replacer, "-_A bool has been allocated here_-", -1))
		processKeywords(content)
		gc.Set(alvarname, altype)
	} else {
		panic("Undefined: " + altype)
	}

}

func handleFree(content []byte) {
	free := ci.Target(string(content), "<wysb.Free", ">")
	vartarget := ci.Target(free, " ", "!")
	gc.Remove(vartarget)
	replacer := "<wysb.Free " + vartarget + "!>"
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleRenameFile(content []byte) {
	fileopen := ci.Target(string(content), "<os.RenameFile", ">")
	filedir := ci.Target(fileopen, " ", " ::")
	newname := ci.Target(fileopen, ":: ", "!")
	e := os.Rename(filedir, newname)
	if e != nil {
		panic(e)
	}
	replacer := "<os.RenameFile " + filedir + " :: " + newname + "!>"
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleRemoveFile(content []byte) {
	fileopen := ci.Target(string(content), "<os.RemoveFile", ">")
	filedir := ci.Target(fileopen, " ", "!")
	err := os.Remove(filedir)
	if err != nil {
		panic(err)
	}
	replacer := "<os.RemoveFile " + filedir + "!>"
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleStringsReplace(content []byte) {
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
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleRandNum(content []byte) {
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
		content = []byte(strings.Replace(string(content), replacer, randomInRange_float, -1))
		processKeywords(content)
	} else if randtype == "int" {
		randsize_int, err := strconv.Atoi(randsize)
		if err != nil {
			panic(err)
		}
		randomInt := rand.Intn(randsize_int)
		replacer := "<rand.Num " + randtype + " :: " + randsize + "!>"
		randomIntStr := strconv.Itoa(randomInt)
		content = []byte(strings.Replace(string(content), replacer, randomIntStr, -1))
		processKeywords(content)
	} else if randtype == "string" {
		randsize_int, err := strconv.Atoi(randsize)
		if err != nil {
			panic(err)
		}
		randomStr := randomString(randsize_int)
		replacer := "<rand.Num " + randtype + " :: " + randsize + "!>"
		content = []byte(strings.Replace(string(content), replacer, randomStr, -1))
		processKeywords(content)
	} else {
		panic("Undefined: " + randtype)
	}
}

func handleUseArg(content []byte) {
	var_f32 := ci.Target(string(content), "<wysb.UseArg", ">")
	var_f32_name := ci.Target(var_f32, " ", "!")
	targetvar, err := gc.Get(var_f32_name)
	if err != nil {
		fmt.Println("")
	}
	str := targetvar.(string)
	replacer := "<wysb.UseArg " + var_f32_name + "!>"
	content = []byte(strings.Replace(string(content), replacer, str, -1))
	processKeywords(content)
}

func handleShellScript(content []byte) {
	var_f32 := ci.Target(string(content), "<os.ShellScript", ">")
	var_f32_name := ci.Target(var_f32, "[] ", "!")
	err := execShellScript(var_f32_name)
	if err != nil {
		fmt.Println("Erro ao executar ShellScript:", err)
	}
	replacer := "<os.ShellScript [] " + var_f32_name + "!>"
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleBatchScript(content []byte) {
	var_f32 := ci.Target(string(content), "<os.BatchScript", ">")
	var_f32_name := ci.Target(var_f32, "[] ", "!")
	err := execBatchScript(var_f32_name)
	if err != nil {
		fmt.Println("Erro ao executar BatchScript:", err)
	}
	replacer := "<os.BatchScript [] " + var_f32_name + "!>"
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleAddon(content []byte) {
	addon := ci.Target(string(content), "<wysb.Addon", ">")
	funcname := ci.Target(addon, " ", "!")

	replacer := "<wysb.Addon " + funcname + "!>"
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handlePrintln(content []byte) {
	var_fun := ci.Target(string(content), "println(", ")")
	result := evaluateExpression(var_fun) 
	logger.Info(context.Background(), result) 
	replacer := "println(" + var_fun + ")"
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func evaluateExpression(expression string) string {

	expression = strings.TrimSpace(expression)

	if strings.HasPrefix(expression, "<math.sum") {
		return handleMathSumExpression(expression)
	} else if strings.HasPrefix(expression, "<math.sub") {
		return handleMathSubExpression(expression)
	} else if strings.HasPrefix(expression, "<math.mult") {
		return handleMathMultExpression(expression)
	} else if strings.HasPrefix(expression, "<math.div") {
		return handleMathDivExpression(expression)
	}

	return "WARNING: EXPRESSION NOT RECOGNIZED (IGNORE THIS IF THE CODE HAS NO ERRORS, THE COMPILER MAY INTERPRET SOME THINGS REPEATEDLY, CAUSING THIS ERROR)"
}

func handleMathSumExpression(expression string) string {

	var_sum := ci.Target(expression, "<math.sum", ">")
	var_num1 := ci.Target(var_sum, " ", " +")
	var_num2 := ci.Target(var_sum, "+ ", "!")

	num1, err1 := strconv.Atoi(var_num1)
	num2, err2 := strconv.Atoi(var_num2)

	if err1 != nil || err2 != nil {
		return "Error: numbers could not be converted"
	}

	result := num1 + num2
	return strconv.Itoa(result)
}

func handleMathSubExpression(expression string) string {
	var_sub := ci.Target(expression, "<math.sub", ">")
	var_num1 := ci.Target(var_sub, " ", " -")
	var_num2 := ci.Target(var_sub, "- ", "!")

	num1, err1 := strconv.Atoi(var_num1)
	num2, err2 := strconv.Atoi(var_num2)

	if err1 != nil || err2 != nil {
		return "Error: numbers could not be converted"
	}

	result := num1 - num2
	return strconv.Itoa(result)
}

func handleMathMultExpression(expression string) string {
	var_mult := ci.Target(expression, "<math.mult", ">")
	var_num1 := ci.Target(var_mult, " ", " *")
	var_num2 := ci.Target(var_mult, "* ", "!")

	num1, err1 := strconv.Atoi(var_num1)
	num2, err2 := strconv.Atoi(var_num2)

	if err1 != nil || err2 != nil {
		return "Error: numbers could not be converted"
	}

	result := num1 * num2
	return strconv.Itoa(result)
}

func handleMathDivExpression(expression string) string {
	var_div := ci.Target(expression, "<math.div", ">")
	var_num1 := ci.Target(var_div, " ", " /")
	var_num2 := ci.Target(var_div, "/ ", "!")

	num1, err1 := strconv.Atoi(var_num1)
	num2, err2 := strconv.Atoi(var_num2)

	if err1 != nil || err2 != nil {
		return "Error: numbers could not be converted"
	}

	if num2 == 0 {
		return "Error: division by zero"
	}

	result := num1 / num2
	return strconv.Itoa(result)
}

func handleWebGet(content []byte) {
	webget := ci.Target(string(content), "<web.get", ">")
	route := ci.Target(webget, " ", " ::")
	returnn := ci.Target(webget, ":: ", "!")
	web.Get(route, func(c *fiber.Ctx) error {
		c.Type("html")
		return c.SendString(returnn)

	})
	replacer := "<web.get " + route + " :: " + returnn + "!>"
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleTimeMeasure(content []byte) {
	start := time.Now()
	measureContent := ci.Target(string(content), "<time.measure>", "</time.measure>")
	processKeywords([]byte(measureContent))
	duration := time.Since(start)
	replacer := "<time.measure>" + measureContent + "</time.measure>"
	content = []byte(strings.Replace(string(content), replacer, fmt.Sprintf("Time: %s", duration), -1))
	processKeywords(content)
}

func handleWebListen(content []byte) {
	weblisten := ci.Target(string(content), "<web.listen ", "!>")
	web.Listen(weblisten)
	replacer := "<web.listen " + weblisten + "!>"
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleWebStatic(content []byte) {
	webstatic := ci.Target(string(content), "<web.static ", ">")
	route := ci.Target(webstatic, " ", " ::")
	dir := ci.Target(webstatic, ":: ", "!")
	web.Static(route, dir)
	replacer := "<web.static " + route + " :: " + dir + "!>"
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
}

func handleWebPost(content []byte) {
	webpost := ci.Target(string(content), "<web.post", ">")
	route := ci.Target(webpost, " ", " ::")
	returnn := ci.Target(webpost, ":: ", "!")
	web.Post(route, func(c *fiber.Ctx) error {
		c.Type("html")
		return c.SendString(returnn)

	})
	replacer := "<web.post " + route + " :: " + returnn + "!>"
	content = []byte(strings.Replace(string(content), replacer, "", -1))
	processKeywords(content)
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
	if finalinput == "help" {
		fmt.Println(`
		Complete list of commands:
		- run: Used to run/test a Wysb file. After using it, you must pass the name of the file.

		- compile: To convert a Wysb file into an executable. After using it, you must pass the file name.

		- cardwmy: Used to initialize "cardwmy", the Wysb package manager. After using it, you must enter the URL of the package to be downloaded.
		`)
	} else if finalinput == "run" {
		fmt.Println("Enter the name of the Wysb file you want to run: ")
		var runinput string
		fmt.Scanln(&runinput)
		fmt.Println("Running the " + runinput + " file")
		ReadWysb(runinput)

	} else if finalinput == "compile" {
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

		fmt.Println("Zip file created successfully:", zipFilename)

	} else if finalinput == "cardwmy" {
		fmt.Println("Enter the ID of the Cardwmy package you want to download (example: '!PackageAuthor/PackageName!'): ")
		var carinput string
		fmt.Scanln(&carinput)
		ci := cutinfo.New()
		username := ci.Target(carinput, "!", "/")
		pkg := ci.Target(carinput, "/", "!")
		result := "https://raw.githubusercontent.com/" + username + "/" + pkg + "/main/main.wys"
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
		panic("unknown command: " + finalinput + "")
	}
}
