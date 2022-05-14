package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var keywords = []string{"sin", "cos", "tan", "csc", "sec", "cot", "sinh", "cosh", "tanh", "csch", "sech", "coth", "arcsin", "arccos", "arctan", "arccsc", "arcsec", "arccot",
	"arcsinh", "arccosh", "arctanh", "arccsch", "arcsech", "arccoth", "log", "lg", "ln", "exp", "det"}

func main() {
	if len(os.Args) <= 1 {
		err("no input files")
		return
	}
	tokenizeRegExp, _ := regexp.Compile(strings.Join(keywords, "|") + "|[0-9.]+|\\S")
	info(tokenizeRegExp.String())
	file, e := os.Open(os.Args[1])
	checkError(e)
	defer file.Close()
	scn := bufio.NewScanner(file)
	for scn.Scan() {
		text := scn.Text()
		comment := strings.Index(text, "//")
		if comment != -1 {
			text = text[0:comment]
		}
		info(strings.Join(tokenizeRegExp.FindAllString(text, -1), ", "))
	}
	file.Close()
}
