package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var keywords = []string{"sin", "cos", "tan", "csc", "sec", "cot", "sinh", "cosh", "tanh", "csch", "sech", "coth", "arcsin", "arccos", "arctan", "arccsc", "arcsec", "arccot",
	"arcsinh", "arccosh", "arctanh", "arccsch", "arcsech", "arccoth", "log", "lg", "ln", "exp", "det", "sqrt"}

var operators = map[string]int{
	"=":  0,
	"<":  0,
	">":  0,
	"<=": 0,
	">=": 0,
	"+":  1,
	"-":  1,
	"*":  2,
	"/":  2,
	"^":  3,
}

func main() {
	if len(os.Args) <= 1 {
		err("no input files")
		return
	}
	tokenizeRegExp, _ := regexp.Compile(strings.Join(keywords, "|") + "|<=|>=|[0-9.]+|\\S")
	info(tokenizeRegExp.String())
	file, e := os.Open(os.Args[1])
	checkError(e)
	defer file.Close()
	scn := bufio.NewScanner(file)
	symbols := make(map[string]Symbol)
	for scn.Scan() {
		text := scn.Text()
		comment := strings.Index(text, "//")
		if comment != -1 {
			text = text[0:comment]
		}
		if len(strings.TrimSpace(text)) == 0 { // blank line
			continue
		}
		tokens := tokenizeRegExp.FindAllString(text, -1)
		for i := 0; i < len(tokens)-1; i++ { // we love implicit multiplication!
			t1, t2 := tokens[i], tokens[i+1]
			if isOperator(t1) || isOperator(t2) {
				continue
			}
			if t1 == "(" || t2 == ")" {
				continue
			}
			_, st1e := symbols[t1]
			if isAlpha(t1) && t2 == "(" && (!st1e || symbols[t1].variant != constantVariant) {
				continue
			}
			tokens = append(tokens[:i+1], tokens[i:]...)
			tokens[i+1] = "*"
			i++
		}
		if len(tokens) >= 4 && isAlpha(tokens[0]) && tokens[1] == "(" {
			next := indexOfSA("(", tokens, 2)
			end := indexOfSA(")", tokens, 2)
			if (next == -1 || next > end) && end+1 < len(tokens) && tokens[end+1] == "=" {
				fdef := strings.Join(tokens[:end+1], "")
				tokens = tokens[end:]
				tokens[0] = fdef
			}
		}
		info(strings.Join(tokens, ", "))
		output := Deque[string]{}
		ops := Stack[string]{}
		for i := 0; i < len(tokens); i++ {
			tp := &(tokens[i])
			if isNumeric(*tp) {
				output.push(tp)
			} else if isAlpha(*tp) {
				ops.push(tp)
			} else if isOperator(*tp) {
				for !ops.empty() && *(ops.peek()) != "(" && operators[*(ops.peek())] >= operators[*tp] {
					output.push(ops.pop())
				}
				ops.push(tp)
			} else if *tp == "(" {
				ops.push(tp)
			} else if *tp == ")" {
				for !ops.empty() && *(ops.peek()) != "(" {
					output.push(ops.pop())
				}
				if ops.empty() || *(ops.peek()) != "(" {
					err("expected left parenthesis")
					return
				}
				ops.pop()
				if !ops.empty() && isAlpha(*(ops.peek())) {
					output.push(ops.pop())
				}
			} else {
				output.push(tp)
			}
			info(output.string())
		}
		for !ops.empty() {
			if *(ops.peek()) == "(" {
				err("unexpected left parenthesis")
				return
			}
			output.push(ops.pop())
			info(output.string())
		}
		info(output.string())
	}
	file.Close()
}
