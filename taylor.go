package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type LineResult int

const (
	outputLineResult   LineResult = 0
	funcDefLineResult  LineResult = 1
	constDefLineResult LineResult = 2
)

var debugMode = false

var operators = map[string]int{
	"|":  0,
	"&":  1,
	"=":  2,
	"<":  2,
	">":  2,
	"<=": 2,
	">=": 2,
	"+":  3,
	"-":  3,
	"*":  4,
	"/":  4,
	"^":  5,
}

var tokenizerRegExp *regexp.Regexp
var symbols map[string]*Symbol
var scope *Scope = nil
var execute = true

func main() {
	files := []string{}
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-debug" {
			debugMode = true
		} else {
			files = append(files, arg)
		}
	}
	transcendentals = map[string]Transcendental{
		"sin":     {1, transSin, functionVariant},
		"cos":     {1, transCos, functionVariant},
		"tan":     {1, transTan, functionVariant},
		"csc":     {1, transCsc, functionVariant},
		"sec":     {1, transSec, functionVariant},
		"cot":     {1, transCot, functionVariant},
		"sinh":    {1, transSinh, functionVariant},
		"cosh":    {1, transCosh, functionVariant},
		"tanh":    {1, transTanh, functionVariant},
		"csch":    {1, transCsch, functionVariant},
		"sech":    {1, transSech, functionVariant},
		"coth":    {1, transCoth, functionVariant},
		"arcsin":  {1, transArcsin, functionVariant},
		"arccos":  {1, transArccos, functionVariant},
		"arctan":  {1, transArctan, functionVariant},
		"arccsc":  {1, transArccsc, functionVariant},
		"arcsec":  {1, transArcsec, functionVariant},
		"arccot":  {1, transArccot, functionVariant},
		"arcsinh": {1, transArcsinh, functionVariant},
		"arccosh": {1, transArccosh, functionVariant},
		"arctanh": {1, transArctanh, functionVariant},
		"arccsch": {1, transArccsch, functionVariant},
		"arcsech": {1, transArcsech, functionVariant},
		"arccoth": {1, transArccoth, functionVariant},
		"log":     {2, transLog, functionVariant},
		"lg":      {1, transLg, functionVariant},
		"ln":      {1, transLn, functionVariant},
		"exp":     {1, transExp, functionVariant},
		// determinant when?
		"sqrt": {1, transSqrt, functionVariant},
		"pi":   {0, transPi, constantVariant},
		"e":    {0, transE, constantVariant},
		"rad":  {1, transRad, functionVariant},
		"deg":  {1, transDeg, functionVariant},
		"abs":  {1, transAbs, functionVariant},
		"inc":  {1, transInc, functionVariant},
		"hirt": {2, transHirt, functionVariant},
		"read": {0, transRead, functionVariant},
		"exit": {0, transExit, functionVariant},
	}
	tokenizerRegExp, _ = regexp.Compile("radians|degrees|sinh|cosh|tanh|csch|sech|coth|sin|cos|tan|csc|sec|cot|arcsinh|arccosh|arctanh|arccsch|arcsech|arccoth|arcsin|arccos|arctan|arccsc|arcsec|arccot|log|lg|ln|exp|det|sqrt|pi|exit|e|rad|deg|abs|inc|hirt|read|<=|>=|<|>|=|\\\"[^\\\"\\\\\\\\]*(\\\\\\\\.[^\\\"\\\\\\\\]*)*\\\"|[-.0-9]+|\\S")
	debug("using tokenizer regex:", tokenizerRegExp.String())
	if len(files) == 0 {
		symbols = make(map[string]*Symbol)
		info("entered line-by-line mode. use the exit function to terminate the program.")
		for reader := bufio.NewReader(os.Stdin); true; {
			fmt.Print("> ")
			line, _ := reader.ReadString('\n')
			line = line[:len(line)-2]
			interpretLine(line, -1, true)
		}
		return
	}
	multi := len(files) > 1
	for _, filename := range files {
		symbols = make(map[string]*Symbol)
		if multi {
			info(filename + ":")
		}
		ok, file := interpret(filename)
		if file != nil {
			file.Close()
		}
		if !ok {
			return
		}
		for k := range symbols {
			delete(symbols, k)
		}
	}
}

func interpret(filename string) (bool, *os.File) {
	file, e := os.Open(filename)
	if e != nil {
		err("could not open file:", filename)
		return false, file
	}
	scn := bufio.NewScanner(file)
	for line := 1; scn.Scan(); line++ {
		text := scn.Text()
		ok, _ := interpretLine(text, line, true)
		if !ok {
			return false, file
		}
	}
	file.Close()
	return true, nil
}

func interpretLine(data string, line int, print bool) (bool, string) {
	comment := strings.Index(data, "//")
	if comment != -1 {
		data = data[0:comment]
	}
	trimmed := strings.TrimSpace(data)
	if len(trimmed) == 0 { // blank line
		return true, ""
	}
	if !execute && print && !strings.HasPrefix(trimmed, ":break") && !strings.HasPrefix(trimmed, ":if") {
		return true, ""
	}
	if isStringLiteral(trimmed) {
		if !print {
			return true, ""
		}
		unwrapped := unwrapStringLiteral(trimmed)
		if trimmed[0] == '"' {
			fmt.Println(unwrapped)
		} else {
			fmt.Print(unwrapped)
		}
		return true, unwrapped
	}
	if strings.HasPrefix(trimmed, ":") { // directive
		brokenData := strings.Split(trimmed, " ")
		directive := strings.TrimSpace(brokenData[0][1:])
		args := []string{}
		if len(brokenData) >= 2 {
			args = brokenData[1:]
		}
		switch directive {
		case "radians":
			{
				if len(args) != 0 {
					errLine("directive radians takes 0 arguments", line)
					return false, ""
				}
				trigMode = Radians
				return true, ""
			}
		case "degrees":
			{
				if len(args) != 0 {
					errLine("directive degrees takes 0 arguments", line)
					return false, ""
				}
				trigMode = Degrees
				return true, ""
			}
		case "delete":
			{
				if len(args) != 1 {
					errLine("directive delete takes 1 argument", line)
					return false, ""
				}
				delete(symbols, args[0])
				return true, ""
			}
		case "if":
			{
				if !execute {
					scope = &Scope{scope, ifVariant, []string{}}
					return true, ""
				}
				if len(args) == 0 {
					errLine("directive if takes an expression as an argument", line)
					return false, ""
				}
				ok, result := interpretLine(strings.Join(args, " "), line, false)
				if !ok {
					return false, ""
				}
				scope = &Scope{scope, ifVariant, []string{}}
				if atof(result) == 0 {
					execute = false
				}
				return true, ""
			}
		case "break":
			{
				if scope == nil {
					errLine("nothing to break out of", line)
					return false, ""
				}
				execute = true
				scope = scope.scope
				return true, ""
			}
		default:
			{
				errLine("unknown directive", line)
				return false, ""
			}
		}
	}
	tokens := tokenizerRegExp.FindAllString(data, -1)
	for i := 0; i < len(tokens)-1; i++ { // we love implicit multiplication!
		t1, t2 := tokens[i], tokens[i+1]
		if isOperator(t1) || isOperator(t2) {
			continue
		}
		if t1 == "(" || t2 == ")" || t1 == "," || t2 == "," {
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
	result := outputLineResult
	var fSymbol *Symbol = nil
	var cSymbol *Symbol = nil
	debug(strconv.Itoa(int(result)))
	if len(tokens) >= 2 && isAlpha(tokens[0]) && tokens[1] == "=" && print {
		result = constDefLineResult
		if len(tokens) == 2 {
			errLine("expression expected after equal (=) sign", line)
			return false, ""
		}
		cName := tokens[0]
		tokens = tokens[2:]
		_, constExists := symbols[cName]
		if !constExists {
			symbols[cName] = &Symbol{cName, constantVariant, []*Symbol{}, nil, Deque[string]{}}
		}
		cSymbol = symbols[cName]
	}
	if len(tokens) >= 4 && isAlpha(tokens[0]) && tokens[1] == "(" {
		next := indexOfSA("(", tokens, 2)
		end := indexOfSA(")", tokens, 2)
		if (next == -1 || next > end) && end+1 < len(tokens) && tokens[end+1] == "=" {
			funcDef := tokens[:end+1]
			result = funcDefLineResult
			if end+2 >= len(tokens) {
				errLine("expression expected after equal (=) sign", line)
				return false, ""
			}
			tokens = tokens[end+2:]
			defIt := 0
			fName := funcDef[defIt]
			_, funcExists := symbols[fName]
			if funcExists {
				errLine("a constant or function is already declared with the name '"+funcDef[defIt]+"'", line)
				return false, ""
			}
			defIt++
			if funcDef[defIt] != "(" {
				errLine("expected left parenthesis to start function parameter list", line)
				return false, ""
			}
			defIt++
			symbols[fName] = &Symbol{fName, functionVariant, []*Symbol{}, nil, Deque[string]{}}
			fSymbol = symbols[fName]
			for ; defIt < len(funcDef); defIt++ {
				current := funcDef[defIt]
				if !isAlpha(current) {
					errLine("a function parameter's name may only contain alphabetical characters", line)
					return false, ""
				}
				foundShade, fsExists := symbols[current]
				var shaded *Symbol = nil
				if fsExists {
					shaded = foundShade
				}
				symbols[current] = &Symbol{current, functionParameterVariant, []*Symbol{}, shaded, Deque[string]{}}
				fSymbol.children = append(fSymbol.children, symbols[current])
				defIt++
				if defIt >= len(funcDef) || (funcDef[defIt] != "," && funcDef[defIt] != ")") {
					errLine("expected right parenthesis or comma", line)
				}
			}
		}
	}
	debug(strings.Join(tokens, ", "))
	output := Deque[string]{}
	ops := Stack[string]{}
	for i := 0; i < len(tokens); i++ {
		tp := &(tokens[i])
		fSymbol, fSymbolExists := symbols[*tp]
		if fSymbolExists && fSymbol.variant != functionVariant {
			fSymbolExists = false
		}
		trans, transExists := transcendentals[*tp]
		if isNumeric(*tp) {
			output.push(tp)
		} else if fSymbolExists || (transExists && trans.variant == functionVariant) {
			ops.push(tp)
		} else if isOperator(*tp) {
			for !ops.empty() && *(ops.peek()) != "(" && (operators[*(ops.peek())] > operators[*tp] || operators[*(ops.peek())] == operators[*tp] && *tp != "^") {
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
				return false, ""
			}
			ops.pop()
			if !ops.empty() && isAlpha(*(ops.peek())) {
				output.push(ops.pop())
			}
		} else if *tp != "," {
			output.push(tp)
		}
		debug(output.string())
	}
	for !ops.empty() {
		if *(ops.peek()) == "(" {
			err("unexpected left parenthesis")
			return false, ""
		}
		output.push(ops.pop())
		debug(output.string())
	}
	debug(output.string())

	full := ""
	if result == funcDefLineResult {
		fSymbol.data = output
	} else {
		value, no, ok := resolveExpression(&output, nil, nil, line)
		if !ok {
			return false, ""
		}
		if result == constDefLineResult {
			cSymbol.data = Deque[string]{}
			cSymbol.data.push(&value)
		} else {
			result := atof(value)
			if result >= -0.00001 && result <= 0.00001 {
				result = 0
			}
			full = ftoa(result)
			if !no && print {
				if !math.IsNaN(result) && !math.IsInf(result, 0) {
					fmt.Println(result)
				} else {
					fmt.Println("undefined")
				}
			}
		}
	}
	// clean up
	if fSymbol != nil {
		// delete all children (function parameters) from symbol table or replace with their shaded versions
		for _, child := range fSymbol.children {
			if child.shaded != nil {
				symbols[child.name] = child.shaded
				child.shaded = nil
			} else {
				delete(symbols, child.name)
			}
		}
	}
	for k := range symbols {
		debug(k, symbols[k].string())
	}
	return true, full
}

func resolveExpression(data *Deque[string], function *Symbol, operations *Stack[string], line int) (string, bool, bool) {
	valueTable := make(map[string]string)
	if function != nil && operations != nil {
		for i := len(function.children) - 1; i >= 0; i-- {
			value := operations.pop()
			if value == nil {
				errLine("function "+function.name+" expects "+strconv.Itoa(len(function.children))+" argument(s), got "+strconv.Itoa(len(function.children)-i-1), line)
				return "", false, false
			}
			valueTable[function.children[i].name] = *value
		}
	}
	localOperations := Stack[string]{}
	noOutput := false
	for current := data.shift(); current != nil; current = data.shift() {
		symbol, fSymbolExists := symbols[*current]
		cSymbolExists := fSymbolExists
		if fSymbolExists && symbol.variant != functionVariant {
			fSymbolExists = false
		}
		value, valueExists := valueTable[*current]
		trans, transExists := transcendentals[*current]
		if fSymbolExists {
			data := symbol.data
			result, nol, ok := resolveExpression(&data, symbol, &localOperations, line)
			if !ok {
				return "", false, false
			}
			if nol {
				noOutput = true
			}
			localOperations.push(&result)
		} else if transExists {
			args := make([]string, 0, trans.argCount)
			for i := 0; i < trans.argCount; i++ {
				arg := localOperations.pop()
				if arg == nil {
					errLine("function "+*current+" expects "+strconv.Itoa(trans.argCount)+" argument(s), got "+strconv.Itoa(i), line)
					return "", false, false
				}
				args = append(args, *arg)
			}
			result, ok := trans.operation(args)
			if !ok {
				return "", false, false
			}
			if result == "" {
				result = "0"
				noOutput = true
			}
			localOperations.push(&result)
		} else if valueExists {
			localOperations.push(&value)
		} else if cSymbolExists {
			localOperations.push(symbol.data.front())
		} else if isNumeric(*current) {
			localOperations.push(current)
		} else if isStringLiteral(*current) {
			value := unwrapStringLiteral(*current)
			localOperations.push(&value)
		} else if isOperator(*current) {
			rhs := localOperations.pop()
			lhs := localOperations.pop()
			if lhs == nil || rhs == nil {
				errLine("'"+*current+"' operator expects 2 operands", line)
				return "", false, false
			}
			lhsv, lhsve := strconv.ParseFloat(*lhs, 64)
			rhsv, rhsve := strconv.ParseFloat(*rhs, 64)
			if lhsve != nil || rhsve != nil {
				errLine("'"+*current+"' operator expects 2 operands", line)
				return "", false, false
			}
			value := ""
			switch *current {
			case "|":
				value = btois(itob(lhsv) || itob(rhsv))
			case "&":
				value = btois(itob(lhsv) && itob(rhsv))
			case "=":
				value = btois(lhsv == rhsv)
			case "<":
				value = btois(lhsv < rhsv)
			case ">":
				value = btois(lhsv > rhsv)
			case "<=":
				value = btois(lhsv <= rhsv)
			case ">=":
				value = btois(lhsv >= rhsv)
			case "+":
				value = ftoa(lhsv + rhsv)
			case "-":
				value = ftoa(lhsv - rhsv)
			case "*":
				value = ftoa(lhsv * rhsv)
			case "/":
				value = ftoa(lhsv / rhsv)
			case "^":
				value = ftoa(math.Pow(lhsv, rhsv))
			default:
				{
					errLine("'"+*current+"' operator has not been implemented yet", line)
					return "", false, false
				}
			}
			localOperations.push(&value)
		} else {
			errLine("unknown symbol encountered: "+*current, line)
			return "", false, false
		}
	}
	if localOperations.empty() {
		errLine("there was an issue while resolving an expression", line)
		return "", false, false
	}
	return *(localOperations.pop()), noOutput, true
}
